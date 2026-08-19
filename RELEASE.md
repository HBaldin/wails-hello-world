# Guia de Release e Auto-Update

## Visão Geral

Este projeto usa GitHub Actions para compilar automaticamente e criar releases no padrão correto para o auto-updater. Quando você faz push de uma tag para o GitHub, o workflow:

1. Compila a aplicação para múltiplas plataformas
2. Gera checksums SHA256 para cada binário
3. Cria um release no GitHub com todos os arquivos

## Configuração Inicial

### 1. Atualizar o repositório em version.go

Edite `version.go` com seu repositório GitHub real:

```go
const (
	updateRepoOwner = "seu-usuario-github"
	updateRepoName  = "wails-hello-world"
)
```

### 2. Criar um Personal Access Token (PAT)

O workflow usa `GITHUB_TOKEN` que já vem pré-configurado, então não precisa de setup adicional.

## Como Publicar uma Nova Versão

### Opção 1: Via Git (Recomendado)

```bash
# 1. Atualizar version.go com a nova versão
# Edite a linha: var version = "1.0.0"

# 2. Commit das alterações
git add version.go
git commit -m "chore: bump version to 1.0.0"
git push origin main

# 3. Criar uma tag (dispara o workflow de release)
git tag -a v1.0.0 -m "Release 1.0.0"
git push origin v1.0.0
```

### Opção 2: Script Local (Para Testes)

```bash
# Compilar localmente para todas as plataformas
chmod +x scripts/build-release.sh
./scripts/build-release.sh 1.0.0

# Os binários estarão em: build/releases/
ls -la build/releases/
```

## Padrão de Nomenclatura

O workflow cria automaticamente arquivos com o padrão:

```
wails-hello-world_<GOOS>_<GOARCH>[.exe]
wails-hello-world_<GOOS>_<GOARCH>[.exe].sha256
```

### Exemplos:
- **macOS ARM64**: `wails-hello-world_darwin_arm64` + `wails-hello-world_darwin_arm64.sha256`
- **macOS Intel**: `wails-hello-world_darwin_amd64` + `wails-hello-world_darwin_amd64.sha256`
- **Linux**: `wails-hello-world_linux_amd64` + `wails-hello-world_linux_amd64.sha256`
- **Windows**: `wails-hello-world_windows_amd64.exe` + `wails-hello-world_windows_amd64.exe.sha256`

## Workflow do GitHub Actions

### Quando é Disparado?

- ✅ **Push de tag** (`v*`): Cria o release automaticamente
- ℹ️ **Push para main**: Apenas compila (sem criar release)

### Plataformas Compiladas

| OS | Arquitetura | Runner | Binário |
|---|---|---|---|
| macOS | ARM64 | `macos-latest` | `darwin_arm64` |
| macOS | x86_64 | `macos-13` | `darwin_amd64` |
| Linux | x86_64 | `ubuntu-latest` | `linux_amd64` |
| Windows | x86_64 | `windows-latest` | `windows_amd64` |

## Testando o Auto-Update

### 1. Compilar versão 1.0.0 Localmente

```bash
./scripts/build-release.sh 1.0.0
```

### 2. Criar um Release de Teste

Se quiser testar sem publicar um release real:

```bash
# Criar release local
gh release create v1.0.0 build/releases/* --draft
```

### 3. Configurar version.go para Seu Repositório

```go
const (
	updateRepoOwner = "seu-usuario"
	updateRepoName  = "wails-hello-world"
)
```

### 4. Executar a Versão Antiga

```bash
./build/releases/wails-hello-world_darwin_arm64
```

### 5. Testar no App

1. Clique em "Check for Update"
2. Deve detectar a versão 1.0.0 (ou a versão que subiu no release)
3. Clique em "Install Update"
4. O app fará download e reiniciará com a nova versão

## Resolução de Problemas

### Release não foi criado
- Verificar se a tag foi feita no formato `v*` (ex: `v1.0.0`)
- Ver logs em "Actions" no GitHub

### Auto-update não encontra a versão
- Verificar se `updateRepoOwner` e `updateRepoName` em `version.go` estão corretos
- Verificar se os binários no release seguem o padrão de nomenclatura
- Verificar se há arquivo `.sha256` para o binário

### Binário não funciona em produção
- Garantir que `wails build` com `-ldflags` gera a versão correta
- Testar com: `./binario --help` ou verificar em "About" no app

## Versionamento Semântico

Sempre use [Semantic Versioning](https://semver.org/) para as tags:

```
v MAJOR . MINOR . PATCH
v 1     . 2      . 3

v1.0.0  ✅ Correto
1.0.0   ❌ Sem 'v' (o workflow ainda funciona, mas use 'v')
v1.0    ❌ Sem patch
v1      ❌ Sem minor
```

## Customizações Futuras

Para adicionar novas plataformas, edite `.github/workflows/build-release.yml`:

```yaml
strategy:
  matrix:
    include:
      # Adicione novos items aqui
      - os: self-hosted
        goos: linux
        goarch: arm64
```

## Referências

- [Documentação Wails - Building](https://wails.io/docs/guides/building/)
- [Documentação Wails - Self Update](https://wails.io/docs/guides/development/)
- [GitHub Actions Documentation](https://docs.github.com/en/actions)
- [Semantic Versioning](https://semver.org/)
