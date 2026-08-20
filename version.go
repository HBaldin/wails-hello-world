package main

// version é a versão semântica da aplicação. É sobrescrita no tempo de
// compilação via:
//
//	go build -ldflags "-X main.version=1.2.3"
//
// então o binário em execução pode se comparar com o lançamento mais recente
// do GitHub ao verificar atualizações.
var version = "0.0.0"

// Defina estes para o repositório GitHub que publica lançamentos para este app.
const (
	updateRepoOwner = "HBaldin"
	updateRepoName  = "wails-hello-world"
)
