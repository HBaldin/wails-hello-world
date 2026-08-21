package main

// =============================================================================
// version — Configuração de versão e auto-update
// =============================================================================
//
// Para usar este template no SEU repositório:
//   1. Altere updateRepoOwner e updateRepoName abaixo para os valores do
//      seu repositório GitHub (ex: "seu-user" / "meu-app").
//   2. O workflow .github/workflows/build-release.yml compila automaticamente
//      quando uma tag v* é empurrada.
//   3. O updater em internal/updater/updater.go baixa a release mais recente
//      e substitui o binário em execução.
//
// version é sobrescrita no tempo de compilação via ldflags. Exemplo:
//
//	go build -ldflags "-X main.version=1.2.3"
//
// O valor default "0.0.0" é usado apenas em desenvolvimento local.
// =============================================================================

// version é a versão semântica da aplicação. Sobrescrita via ldflags no build.
var version = "0.0.0"

// Defina com os dados do SEU repositório GitHub.
const (
	updateRepoOwner = "HBaldin"
	updateRepoName  = "wails-hello-world"
)