# LEIA-ME

## Sobre

Aplicativo desktop Wails com frontend em React + TypeScript e módulo de 
auto-atualização em Go.

Você pode configurar o projeto editando `wails.json`. Mais informações sobre 
as configurações do projeto podem ser encontradas aqui:
https://wails.io/docs/reference/project-config

## Arquitetura do Frontend

`frontend/src` segue um layout modificado de Atomic Design + co-localização: 
cada componente possui uma pasta com tudo o que precisa, e partes que 
pertencem apenas a um pai vivem aninhadas nela, em vez de estar na árvore 
global.

```
frontend/src/
├── assets/                     # ícones, fontes e imagens globais
├── components/
│   └── ui/                     # primitivas reutilizáveis (ex: Button)
│       └── Button/
│           ├── Button.tsx
│           ├── Button.module.css
│           ├── Button.types.ts
│           ├── Button.test.tsx
│           └── index.ts
├── features/                   # módulos de negócio
│   ├── greeting/
│   │   ├── api/                # chamadas para o backend Go desta feature
│   │   ├── GreetingPage.tsx    # container: busca dados, conecta subcomponentes
│   │   └── components/
│   │       └── GreetForm/      # componente burro, usado apenas por GreetingPage
│   └── update/
│       ├── api/
│       ├── hooks/               # useUpdater — gerencia a máquina de estados
│       ├── types/
│       ├── UpdatePanel.tsx      # container
│       └── components/
│           └── UpdateStatus/    # componente burro, usado apenas por UpdatePanel
├── services/                    # ponte genérica Go/IPC (wailsClient.ts)
└── App.tsx
```

Regras práticas:

- **Aninhamento local** — um subcomponente usado por exatamente um pai vive 
  dentro da pasta `components/` daquele pai, não em `components/ui/` global.
- **Promoção** — se um segundo consumidor precisar, mova sua pasta para 
  `src/components/ui/`.
- **Encapsulamento via `index.ts`** — o `index.ts` de cada pasta de componente 
  exporta apenas o componente público; subcomponentes internos não são 
  acessíveis de fora.
- **Burro vs. container** — subcomponentes (`GreetForm`, `UpdateStatus`) são 
  UI pura controlada por props; o componente de nível superior da pasta 
  (`GreetingPage`, `UpdatePanel`) é a única coisa que fala com o backend Go, 
  via seu módulo `api/` local e a ponte compartilhada 
  `services/wailsClient.ts`. Isso mantém IPC fora da árvore visual, então um 
  componente burro renderiza igual no Storybook/testes que no app real.

Execute a suite de testes do frontend em `frontend/`: `npm test`.

## Desenvolvimento em Tempo Real

Para executar em modo de desenvolvimento em tempo real, execute `wails dev` no 
diretório do projeto. Isso executará um servidor de desenvolvimento Vite que 
fornecerá recarga a quente muito rápida de suas alterações no frontend. Se você 
quiser desenvolver em um navegador e ter acesso aos seus métodos Go, há também 
um servidor de desenvolvimento que roda em http://localhost:34115. Conecte-se 
a isso em seu navegador e você pode chamar seu código Go a partir do devtools.

## Compilação

Para compilar um pacote redistributível em modo de produção, use `wails build`.

Para gravar uma versão real no binário (usado pelo auto-atualizador para saber 
qual versão está em execução), compile com:

```
wails build -ldflags "-X main.version=1.2.3"
```

Sem essa flag o app se reporta como versão `0.0.0`.

## Auto-atualização

O pacote [internal/updater](internal/updater/updater.go) verifica o 
**lançamento mais recente do GitHub** do repositório configurado em 
[version.go](version.go) (`updateRepoOwner` / `updateRepoName`) e, se uma 
tag semver mais recente for publicada, faz o download e a aplica no local 
usando [minio/selfupdate](https://github.com/minio/selfupdate).

Está conectado à UI via três métodos Go vinculados (veja [app.go](app.go)):

- `AppVersion()` — a versão em execução atualmente.
- `CheckForUpdate()` — consulta o GitHub, retorna se uma atualização está disponível.
- `InstallUpdate()` — faz download, verifica e aplica o lançamento verificado 
  por último, depois reinicia o app.

### Publicando um lançamento que o atualizador possa encontrar

Para cada lançamento do GitHub, faça upload de dois arquivos por plataforma de 
destino, seguindo esta convenção de nomenclatura:

```
wails-hello-world_<GOOS>_<GOARCH>[.exe]          o binário em si
wails-hello-world_<GOOS>_<GOARCH>[.exe].sha256   seu sha256sum (formato de saída sha256sum)
```

Por exemplo, para macOS no Apple Silicon:

```
wails-hello-world_darwin_arm64
wails-hello-world_darwin_arm64.sha256
```

A tag git do lançamento deve ser um semver válido (um `v` inicial é aceitável, 
ex: `v1.2.3`) — é isso que é comparado com `main.version` incorporado no 
binário. Esta convenção corresponde ao que ferramentas como 
[GoReleaser](https://goreleaser.com/) produzem com configuração mínima.

Antes de enviar, defina `updateRepoOwner`/`updateRepoName` em 
[version.go](version.go) para seu repositório GitHub real.
