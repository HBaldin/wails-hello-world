#!/usr/bin/env bash
# =============================================================================
# init-template.sh — Inicializa este template para um novo repositório
# =============================================================================
# Uso: ./scripts/init-template.sh
#
# Este script renomeia o módulo Go, o app no wails.json e os imports para
# o seu próprio repositório. Execute UMA VEZ ao clonar/forkar este template.
#
# Requer: git, go, sed
# =============================================================================

set -euo pipefail

OLD_MODULE="wails-hello-world"
OLD_NAME="wails-hello-world"

echo "=== Inicialização do Template Wails + Auto-Update ==="
echo ""

# --- Coletar informações do usuário ---
read -rp "Novo nome do módulo Go (ex: github.com/seu-user/meu-app): " NEW_MODULE
read -rp "Novo nome do app (ex: MeuApp): " NEW_NAME
read -rp "Seu usuário GitHub: " GH_USER
read -rp "Nome do repositório GitHub: " GH_REPO

if [ -z "$NEW_MODULE" ] || [ -z "$NEW_NAME" ]; then
  echo "❌ Nome do módulo e nome do app são obrigatórios."
  exit 1
fi

echo ""
echo "Resumo das alterações:"
echo "  Módulo Go:  $OLD_MODULE → $NEW_MODULE"
echo "  App name:   $OLD_NAME → $NEW_NAME"
echo "  GitHub:     $GH_USER / $GH_REPO"
echo ""

read -rp "Continuar? (s/N): " CONFIRM
if [ "$CONFIRM" != "s" ] && [ "$CONFIRM" != "S" ]; then
  echo "Cancelado."
  exit 0
fi

# --- Renomear módulo Go ---
echo "✏️ Atualizando go.mod..."
sed -i '' "s|^module $OLD_MODULE|module $NEW_MODULE|g" go.mod

# --- Atualizar imports nos arquivos .go ---
echo "✏️ Atualizando imports Go..."
find . -name '*.go' -not -path './vendor/*' -not -path './node_modules/*' -not -path './deepseek-harness/*' \
  -exec sed -i '' "s|$OLD_MODULE|$NEW_MODULE|g" {} +

# --- Atualizar wails.json ---
echo "✏️ Atualizando wails.json..."
sed -i '' "s|\"name\": \"$OLD_NAME\"|\"name\": \"$NEW_NAME\"|g" wails.json
sed -i '' "s|\"outputfilename\": \"$OLD_NAME\"|\"outputfilename\": \"$NEW_NAME\"|g" wails.json

# --- Atualizar version.go ---
echo "✏️ Atualizando version.go..."
sed -i '' "s|updateRepoOwner = \".*\"|updateRepoOwner = \"$GH_USER\"|g" version.go
sed -i '' "s|updateRepoName = \".*\"|updateRepoName = \"$GH_REPO\"|g" version.go

# --- Regenerar bindings Wails ---
echo "🔗 Regenerando bindings Go → JS..."
wails generate module 2>/dev/null || echo "  ⚠ wails generate module falhou (instale o Wails CLI primeiro)"

echo ""
echo "✅ Template inicializado!"
echo ""
echo "Próximos passos:"
echo "  1. git add -A"
echo "  2. git commit -m 'feat: initialize from wails-hello-world template'"
echo "  3. Verifique o arquivo version.go"
echo "  4. git tag -a v0.1.0 -m 'First release' && git push origin v0.1.0"
echo ""