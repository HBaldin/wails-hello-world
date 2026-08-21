# Guia de Release e Auto-Update

## Visão Geral

Este projeto usa GitHub Actions para compilar automaticamente e criar releases
no padrão correto para o auto-updater. Quando você faz push de uma tag `v*`,
o workflow:

1. Compila a aplicação para 4 plataformas (macOS ARM64, macOS AMD64,
   Linux AMD64, Windows AMD64)
2. Gera checksums SHA256 para cada binário
3. Cria um release no GitHub com todos os arquivos

## 🔧 Configuração Inicial

### 1. Executar o script de inicialização

```bash
./scripts/init-template.sh
```

Isso atualiza: módulo Go, nome do app, versão, repositório GitHub, e imports.

### 2. Alternativa manual

Edite `version.go` com seu repositório GitHub real:

```go
const (
    updateRepoOwner = "seu-usuario"
    updateRepoName  = "meu-app"
)
```

### 3. Verificar permissões do GitHub Actions

O workflow usa `GITHUB_TOKEN` que já vem pré-configurado — nenhum setup adicional
é necessário.

## 🚀 Como Publicar uma Nova Versão

### Via Git (Recomendado)

```bash
# 1. Atualizar version.go com a nova versão (opcional)
# Edite a linha: var version = "1.0.0"

# 2. Commit das alterações
git add version.go
git commit -m "chore: bump version to 1.0.0"
git push origin main

# 3. Criar uma tag (dispara o workflow de release)
git tag -a v1.0.0 -m "Release 1.0.0"
git push origin v1.0.0
```

### Script Local (Para Testes)

```bash
./scripts/build-release.sh 1.0.0
# Binários em: build/releases/
```

## 📦 Padrão de Nomenclatura

O workflow cria automaticamente arquivos com o padrão:

```
<nome-repo>_<GOOS>_<GOARCH>[.exe]
<nome-repo>_<GOOS>_<GOARCH>[.exe].sha256
```

### Exemplos:

- **macOS ARM64**: `meu-app_darwin_arm64` + `meu-app_darwin_arm64.sha256`
- **macOS Intel**: `meu-app_darwin_amd64` + `meu-app_darwin_amd64.sha256`
- **Linux**: `meu-app_linux_amd64` + `meu-app_linux_amd64.sha256`
- **Windows**: `meu-app_windows_amd64.exe` + `meu-app_windows_amd64.exe.sha256`

## ⚙️ Workflow do GitHub Actions

### Gatilhos

| Evento | Ação |
|---|---|
| Push de tag `v*` | Compila + publica release |
| Push para `main` | Apenas compila (validação, sem release) |

### Plataformas Compiladas

| OS | Arquitetura | Runner |
|---|---|---|
| macOS | ARM64 | `macos-latest` |
| macOS | x86_64 | `macos-13` |
| Linux | x86_64 | `ubuntu-22.04` |
| Windows | x86_64 | `windows-latest` |

## 🧪 Testando o Auto-Update

### 1. Compilar versão local

```bash
./scripts/build-release.sh 0.1.0
```

### 2. Subir como release de teste

```bash
gh release create v0.1.0 build/releases/* --draft
```

### 3. Configurar version.go para seu repositório (se ainda não fez)

```go
const (
    updateRepoOwner = "seu-usuario"
    updateRepoName  = "meu-app"
)
```

### 4. Executar a versão antiga

```bash
./build/releases/meu-app_darwin_arm64
```

### 5. Testar no app

1. Clique em "Check for Update"
2. Deve detectar a versão 0.1.0 (ou a versão do release)
3. Clique em "Install Update"
4. O app fará download e reiniciará com a nova versão

## 🐛 Resolução de Problemas

### Release não foi criado

- Verificar se a tag tem formato `v*` (ex: `v1.0.0`)
- Ver logs em "Actions" no GitHub
- Verificar se o workflow `fail_on_unmatched_files: true` não falhou

### Auto-update não encontra a versão

- Verificar se `updateRepoOwner` e `updateRepoName` em `version.go` estão corretos
- Verificar se os binários no release seguem o padrão de nomenclatura
- Verificar se há arquivo `.sha256` para cada binário

### Erro no build Linux: "libwebkit2gtk-4.0-dev not found"

**Causa**: Ubuntu 24.04+ usa `libwebkit2gtk-4.1-dev` (versão 4.1, não 4.0).

**Solução**: O workflow já usa `libwebkit2gtk-4.1-dev`. Para build local:

```bash
sudo apt-get install libwebkit2gtk-4.1-dev
```

### Erro de checksum no auto-update

O checksum SHA-256 é gerado com `sha256sum`. Verifique:

- O formato do arquivo `.sha256` é `<hex>  <nome-do-arquivo>`
- O checksum corresponde exatamente ao binário baixado

## Versionamento Semântico

Sempre use [Semantic Versioning](https://semver.org/) para as tags:

```
v MAJOR . MINOR . PATCH
v 1     . 2      . 3

v1.0.0  ✅ Correto
1.0.0   ✅ Aceito (mas prefira com 'v')
v1.0    ❌ Sem patch
```

## Customizações Futuras

Para adicionar novas plataformas, edite `.github/workflows/build-release.yml`:

```yaml
strategy:
  matrix:
    include:
      - os: [novo-runner]
        goos: [novo-os]
        goarch: [nova-arch]
```

## Referências

- [Documentação Wails - Building](https://wails.io/docs/guides/building/)
- [Documentação Wails](https://wails.io/docs/next/)
- [minio/selfupdate](https://github.com/minio/selfupdate)
- [GitHub Actions](https://docs.github.com/en/actions)
- [Semantic Versioning](https://semver.org/)