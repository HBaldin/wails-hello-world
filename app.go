package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"

	wailsruntime "github.com/wailsapp/wails/v2/pkg/runtime"

	"wails-hello-world/internal/updater"
)

// Estrutura da App
type App struct {
	ctx        context.Context
	updateCfg  updater.Config
	pendingRel *updater.Release
}

// NewApp cria uma nova estrutura da aplicação App
func NewApp() *App {
	return &App{
		updateCfg: updater.Config{
			RepoOwner:      updateRepoOwner,
			RepoName:       updateRepoName,
			CurrentVersion: version,
		},
	}
}

// startup é chamado quando o app inicia. O contexto é salvo
// para que possamos chamar os métodos de runtime
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// Greet retorna uma saudação para o nome dado
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Olá %s, é hora do show!", name)
}

// AppVersion retorna a versão da compilação em execução atualmente.
func (a *App) AppVersion() string {
	return version
}

// UpdateInfo é o que o frontend recebe quando uma atualização está disponível.
type UpdateInfo struct {
	Available bool   `json:"available"`
	Version   string `json:"version"`
	Notes     string `json:"notes"`
}

// CheckForUpdate entra em contato com o GitHub e reporta se um lançamento mais recente existe.
func (a *App) CheckForUpdate() (UpdateInfo, error) {
	rel, err := updater.CheckForUpdate(a.ctx, a.updateCfg)
	if err != nil {
		return UpdateInfo{}, err
	}
	if rel == nil {
		return UpdateInfo{Available: false}, nil
	}
	a.pendingRel = rel
	return UpdateInfo{Available: true, Version: rel.Version, Notes: rel.Notes}, nil
}

// InstallUpdate faz download e aplica o lançamento encontrado pela chamada
// CheckForUpdate mais recente, depois reinicia a aplicação.
func (a *App) InstallUpdate() error {
	if a.pendingRel == nil {
		return fmt.Errorf("nenhuma atualização foi verificada ainda")
	}

	if err := updater.Apply(a.ctx, a.updateCfg, a.pendingRel); err != nil {
		return err
	}

	exe, err := os.Executable()
	if err != nil {
		// A atualização foi aplicada mas não podemos reiniciar automaticamente;
		// a nova versão executará na próxima vez que o usuário iniciar o app.
		wailsruntime.Quit(a.ctx)
		return nil
	}

	cmd := exec.Command(exe, os.Args[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("atualização aplicada, mas falhou ao reiniciar: %w", err)
	}

	wailsruntime.Quit(a.ctx)
	return nil
}
