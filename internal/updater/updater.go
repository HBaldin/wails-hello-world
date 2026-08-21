// Package updater implements self-update support for the application.
//
// It checks the GitHub "latest release" of a repository for a newer
// version than the one currently running, and — if the user confirms —
// downloads the matching platform asset, verifies its SHA-256 checksum
// and replaces the currently running executable with it (via
// github.com/minio/selfupdate).
//
// Release assets are expected to follow the naming convention:
//
//	<repo>_<GOOS>_<GOARCH>[.exe]        the binary for the platform
//	<repo>_<GOOS>_<GOARCH>[.exe].sha256 a text file containing its hex SHA-256
//
// The workflow in .github/workflows/build-release.yml produces assets in this
// exact format automatically.
//
// # macOS note
//
// On macOS the running executable is inside a .app bundle and may be code
// signed. minio/selfupdate replaces the binary in-place, which invalidates
// the ad-hoc signature Apple applies to Wails binaries on first launch.
// The binary still runs, but macOS may warn the user. For distribution,
// sign the binary with a Developer ID and notarize it. See the RELEASE.md
// guide for details.
package updater

import (
	"context"
	"crypto"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"runtime"
	"strings"
	"time"

	"github.com/blang/semver/v4"
	"github.com/minio/selfupdate"
)

// Config describes where to look for updates and what is currently running.
type Config struct {
	// RepoOwner and RepoName identify the GitHub repository that publishes releases.
	RepoOwner string
	RepoName  string

	// CurrentVersion is the version of the running binary, e.g. "1.2.3" or "v1.2.3".
	CurrentVersion string

	// HTTPClient is used for all network calls. If nil, a client with a
	// reasonable timeout is created.
	HTTPClient *http.Client
}

func (c Config) httpClient() *http.Client {
	if c.HTTPClient != nil {
		return c.HTTPClient
	}
	return &http.Client{Timeout: 30 * time.Second}
}

// Release describes an available update found by CheckForUpdate.
type Release struct {
	Version     string `json:"version"`
	Notes       string `json:"notes"`
	assetURL    string
	checksumURL string
}

type ghAsset struct {
	Name               string `json:"name"`
	BrowserDownloadURL string `json:"browser_download_url"`
}

type ghRelease struct {
	TagName string    `json:"tag_name"`
	Body    string    `json:"body"`
	Assets  []ghAsset `json:"assets"`
}

// assetName returns the expected release asset name for the current platform.
func assetName(repo string) string {
	name := fmt.Sprintf("%s_%s_%s", repo, runtime.GOOS, runtime.GOARCH)
	if runtime.GOOS == "windows" {
		name += ".exe"
	}
	return name
}

// CheckForUpdate queries the latest GitHub release and returns a Release
// describing it if it is newer than cfg.CurrentVersion. It returns
// (nil, nil) when the running binary is already up to date.
func CheckForUpdate(ctx context.Context, cfg Config) (*Release, error) {
	current, err := semver.ParseTolerant(cfg.CurrentVersion)
	if err != nil {
		return nil, fmt.Errorf("parse current version %q: %w", cfg.CurrentVersion, err)
	}

	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", cfg.RepoOwner, cfg.RepoName)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	// GitHub API requires a User-Agent header.
	req.Header.Set("User-Agent", fmt.Sprintf("%s-updater/1.0", cfg.RepoName))

	resp, err := cfg.httpClient().Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetch latest release: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 4<<10))
		return nil, fmt.Errorf("fetch latest release: unexpected status %s: %s", resp.Status, body)
	}

	var rel ghRelease
	if err := json.NewDecoder(resp.Body).Decode(&rel); err != nil {
		return nil, fmt.Errorf("decode release: %w", err)
	}

	latest, err := semver.ParseTolerant(rel.TagName)
	if err != nil {
		return nil, fmt.Errorf("parse release tag %q: %w", rel.TagName, err)
	}

	if !latest.GT(current) {
		return nil, nil
	}

	wantBin := assetName(cfg.RepoName)
	var binURL, sumURL string
	for _, a := range rel.Assets {
		switch a.Name {
		case wantBin:
			binURL = a.BrowserDownloadURL
		case wantBin + ".sha256":
			sumURL = a.BrowserDownloadURL
		}
	}
	if binURL == "" {
		return nil, fmt.Errorf("release %s has no asset named %q for this platform", rel.TagName, wantBin)
	}

	return &Release{
		Version:     latest.String(),
		Notes:       rel.Body,
		assetURL:    binURL,
		checksumURL: sumURL,
	}, nil
}

// Apply downloads the release binary, verifies its checksum (when a
// checksum asset was published) and replaces the currently running
// executable with it. On success, the calling process must restart itself
// for the new version to take effect.
func Apply(ctx context.Context, cfg Config, rel *Release) error {
	client := cfg.httpClient()

	var checksum []byte
	if rel.checksumURL != "" {
		sum, err := fetchChecksum(ctx, client, rel.checksumURL)
		if err != nil {
			return fmt.Errorf("fetch checksum: %w", err)
		}
		checksum = sum
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, rel.assetURL, nil)
	if err != nil {
		return err
	}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("download update: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download update: unexpected status %s", resp.Status)
	}

	opts := selfupdate.Options{
		Checksum: checksum,
		// Use SHA-256 explicitly for hash verification.
		Hash: crypto.SHA256,
	}
	if err := selfupdate.Apply(resp.Body, opts); err != nil {
		if rerr := selfupdate.RollbackError(err); rerr != nil {
			return fmt.Errorf("update failed and rollback failed: %v (original error: %w)", rerr, err)
		}
		return fmt.Errorf("apply update: %w", err)
	}

	return nil
}

// fetchChecksum downloads a "<hex-digest>  <filename>"-style sha256 file
// (the format produced by `sha256sum`) and returns the decoded digest.
func fetchChecksum(ctx context.Context, client *http.Client, url string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status %s", resp.Status)
	}

	body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<10))
	if err != nil {
		return nil, err
	}

	hexDigest := strings.Fields(string(body))
	if len(hexDigest) == 0 {
		return nil, fmt.Errorf("empty checksum file")
	}

	digest, err := hex.DecodeString(hexDigest[0])
	if err != nil {
		return nil, fmt.Errorf("decode checksum: %w", err)
	}
	if len(digest) != sha256.Size {
		return nil, fmt.Errorf("checksum has unexpected length %d", len(digest))
	}
	return digest, nil
}