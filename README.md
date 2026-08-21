# Wails + Auto-Update Template

![Go](https://img.shields.io/badge/Go-1.23+-00ADD8?logo=go)
![Wails](https://img.shields.io/badge/Wails-v2.15+-red)
![License](https://img.shields.io/badge/License-MIT-green)

Template de aplicação desktop **cross-platform** com **Go + Wails v2** e
**auto-update via GitHub Releases**.

## ✨ Funcionalidades

- **Desktop nativo** — Aplicações para Windows, macOS e Linux com HTML/CSS/React
- **Auto-update integrado** — Verifica, baixa e instala atualizações automaticamente
  via GitHub Releases
- **Pipeline CI/CD** — GitHub Actions que compila para 4 plataformas e publica
  releases ao empurrar tags
- **Verificação de checksum** — SHA-256 de cada binário antes de aplicar a
  atualização
- **Template pronto** — Um script inicializa tudo para o seu repositório

## 🚀 Começando

### 1. Clonar e inicializar

```bash
git clone <seu-fork-deste-template> meu-app
cd meu-app
./scripts/init-template.sh
```

O script pede o nome do seu módulo Go, nome do app e dados do GitHub e
renomeia todos os arquivos automaticamente.

### 2. Desenvolvimento local

```bash
# Instalar Wails CLI
go install github.com/wailsapp/wails/v2/cmd/wails@latest

# Desenvolvimento (com hot-reload do frontend)
wails dev

# Build para produção
wails build -ldflags "-X main.version=0.1.0" -platform darwin/arm64
```

## 📦 Pipeline de Release

### Como publicar uma versão

```bash
# 1. Atualizar version.go com a nova versão (opcional, a tag prevalece)
# 2. Commit
git add version.go && git commit -m "chore: bump to 1.0.0"
git push origin main

# 3. Criar tag (dispara o pipeline de release)
git tag -a v1.0.0 -m "Release 1.0.0"
git push origin v1.0.0
```

O GitHub Actions vai:

1. Compilar para **macOS ARM64**, **macOS AMD64**, **Linux AMD64**, **Windows AMD64**
2. Gerar checksums SHA-256
3. Publicar um GitHub Release com os binários

### Build local (para testes)

```bash
./scripts/build-release.sh 1.0.0
# Binários em: build/releases/
```

## 🔄 Auto-Update (Como Funciona)

Quando o aplicativo está em execução:

1. **CheckForUpdate()** — Consulta a API do GitHub (`/releases/latest`) e compara
   a versão remota com a versão compilada no binário (via `-ldflags`).
2. Se há uma versão mais recente, o frontend mostra um botão "Instalar".
3. **InstallUpdate()** — Baixa o binário da plataforma correspondente, verifica o
   checksum SHA-256 e substitui o executável em execução via
   [minio/selfupdate](https://github.com/minio/selfupdate).
4. O app reinicia automaticamente com a nova versão.

### Arquitetura

```
┌─────────────────────────────────────────┐
│  Frontend (React / TypeScript / Vite)   │
│  ┌─────────────────────────────────────┐│
│  │  UpdatePanel.tsx                    ││
│  │  ├─ "Check for Update" →           ││
│  │  │   window.go.main.App.           ││
│  │  │     CheckForUpdate()            ││
│  │  ├─ "Install Update" →            ││
│  │  │   window.go.main.App.           ││
│  │  │     InstallUpdate()             ││
│  └─────────────────────────────────────┘│
└──────────────┬──────────────────────────┘
               │ (Wails bindings)
┌──────────────▼──────────────────────────┐
│  Backend (Go)                          │
│  ┌────────────────────────────────────┐ │
│  │  app.go                            │ │
│  │  ├─ CheckForUpdate() → GitHub API │ │
│  │  └─ InstallUpdate() → selfupdate  │ │
│  └────────────────────────────────────┘ │
│  ┌────────────────────────────────────┐ │
│  │  internal/updater/updater.go      │ │
│  │  ├─ CheckForUpdate()             │ │
│  │  ├─ Apply()                      │ │
│  │  └─ fetchChecksum()              │ │
│  └────────────────────────────────────┘ │
└─────────────────────────────────────────┘
```

### Convenção de nomenclatura dos assets

O updater procura no GitHub Release assets com o padrão:

```
<nome-repo>_<GOOS>_<GOARCH>[.exe]
<nome-repo>_<GOOS>_<GOARCH>[.exe].sha256
```

Exemplos gerados pelo pipeline:

| Plataforma | Binário | Checksum |
|---|---|---|
| macOS ARM64 | `meu-app_darwin_arm64` | `meu-app_darwin_arm64.sha256` |
| macOS AMD64 | `meu-app_darwin_amd64` | `meu-app_darwin_amd64.sha256` |
| Linux AMD64 | `meu-app_linux_amd64` | `meu-app_linux_amd64.sha256` |
| Windows AMD64 | `meu-app_windows_amd64.exe` | `meu-app_windows_amd64.exe.sha256` |

## 📁 Estrutura do Projeto

```
├── .github/workflows/
│   └── build-release.yml      # Pipeline CI/CD
├── internal/updater/
│   └── updater.go              # Lógica de auto-update
├── scripts/
│   ├── build-release.sh        # Build local de releases
│   └── init-template.sh        # Inicialização do template
├── frontend/
│   ├── src/
│   │   ├── features/
│   │   │   ├── update/         # Componente de atualização (React)
│   │   │   └── greeting/       # Exemplo de página
│   │   ├── assets/             # Ícones, fontes e imagens globais
│   │   ├── components/ui/      # Componentes reutilizáveis
│   │   ├── services/           # Ponte genérica Go/IPC (wailsClient.ts)
│   │   ├── App.tsx
│   │   └── main.tsx
│   ├── wailsjs/                # Bindings Go → JS (gerado automaticamente)
│   └── package.json
├── app.go                      # Estrutura principal do app
├── main.go                     # Ponto de entrada
├── version.go                  # Versão e config do repositório
├── wails.json                  # Configuração Wails
└── go.mod
```

### Organização do Frontend

O frontend segue um layout modificado de Atomic Design + co-localização:

- **Aninhamento local** — um subcomponente usado por exatamente um pai vive
  dentro da pasta `components/` daquele pai, não em `components/ui/` global.
- **Promoção** — se um segundo consumidor precisar, mova sua pasta para
  `src/components/ui/`.
- **Encapsulamento via `index.ts`** — o `index.ts` de cada pasta de componente
  exporta apenas o componente público.
- **Burro vs. container** — subcomponentes são UI pura controlada por props;
  o componente de nível superior da pasta é a única coisa que fala com o backend Go.

Execute a suite de testes do frontend em `frontend/`: `npm test`.

## 🛠️ Requisitos

- Go **1.23+**
- Node.js **18+**
- Wails CLI v2
- Dependências do sistema por plataforma

### macOS

```bash
xcode-select --install
brew install pkg-config
```

### Linux (Ubuntu/Debian)

```bash
sudo apt-get install -y libgtk-3-dev libwebkit2gtk-4.1-dev build-essential pkg-config
```

### Windows

Veja o guia oficial do Wails: https://wails.io/docs/next/guides/installation

## ⚠️ Notas sobre macOS

No macOS, o binário fica dentro de um bundle `.app`. A substituição via
`selfupdate` invalida a assinatura ad-hoc que a Apple aplica automaticamente.
Para distribuição oficial:

1. **Assine com Apple Developer ID**: `codesign -s "Developer ID Application: Seu Nome"`
2. **Notarize**: Envie para a Apple para notarização
3. O app ainda funciona sem assinatura, mas pode mostrar avisos no Gatekeeper

## 📄 Licença

MIT — Use livremente como base para seus projetos.

---

**Feito com ❤️ usando [Wails](https://wails.io/), [Go](https://go.dev/) e [React](https://react.dev/)**