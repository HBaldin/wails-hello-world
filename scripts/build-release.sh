#!/bin/bash
# =============================================================================
# build-release.sh — Script de build local para releases
# =============================================================================
# Uso: ./scripts/build-release.sh [versão] [plataformas]
#
# Exemplos:
#   ./scripts/build-release.sh 1.0.0                                     # Todas as plataformas
#   ./scripts/build-release.sh 1.0.0 darwin/arm64                        # Apenas macOS ARM64
#   ./scripts/build-release.sh 1.0.0 "darwin/arm64 linux/amd64"          # Múltiplas
#
# O script gera binários e checksums SHA-256 no diretório build/releases/.
# Use o nome do repositório do seu go.mod module.
# =============================================================================

set -euo pipefail

APP_NAME="${GO_APP_NAME:-$(grep '^module ' go.mod | awk '{print $2}')}"
VERSION="${1:-0.0.0}"
PLATFORMS="${2:-darwin/arm64 darwin/amd64 linux/amd64 windows/amd64}"

GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m'

BUILD_DIR="build/releases"
mkdir -p "$BUILD_DIR"

echo -e "${BLUE}=== Compilando ${APP_NAME} v${VERSION} ===${NC}"
echo "Plataformas: $PLATFORMS"
echo ""

for PLATFORM in $PLATFORMS; do
  IFS='/' read -r GOOS GOARCH <<< "$PLATFORM"

  echo -e "${BLUE}Compilando para $GOOS/$GOARCH...${NC}"

  # Nome do binário seguindo a convenção do updater
  if [ "$GOOS" = "windows" ]; then
    BINARY_NAME="${APP_NAME}_${GOOS}_${GOARCH}.exe"
  else
    BINARY_NAME="${APP_NAME}_${GOOS}_${GOARCH}"
  fi

  # Gera as bindings Go → JS antes do build (necessário sempre)
  wails generate module

  # Build com Wails
  wails build \
    -ldflags "-X main.version=${VERSION}" \
    -platform "${GOOS}/${GOARCH}" \
    -o "${BINARY_NAME}"

  # Localiza o binário gerado (o Wails coloca em build/bin/)
  SRC=""
  if [ "$GOOS" = "darwin" ] || [ "$GOOS" = "macos" ]; then
    BUNDLE=$(find build/bin -name '*.app' -type d -maxdepth 2 2>/dev/null | head -1)
    if [ -n "$BUNDLE" ]; then
      SRC="$BUNDLE/Contents/MacOS/$(basename "$BUNDLE" .app)"
    fi
  else
    CANDIDATE="build/bin/${APP_NAME}"
    [ "$GOOS" = "windows" ] && CANDIDATE="build/bin/${APP_NAME}.exe"
    [ -f "$CANDIDATE" ] && SRC="$CANDIDATE"
  fi

  if [ -z "$SRC" ] || [ ! -f "$SRC" ]; then
    echo "  ⚠ Binário não encontrado em build/bin/. Tentando path alternativo..."
    # Fallback: procurar no diretório de output
    SRC=$(find build -name "$BINARY_NAME" -type f 2>/dev/null | head -1)
  fi

  if [ -z "$SRC" ] || [ ! -f "$SRC" ]; then
    echo "  ❌ Binário não encontrado para $GOOS/$GOARCH"
    exit 1
  fi

  cp "$SRC" "$BUILD_DIR/$BINARY_NAME"
  chmod +x "$BUILD_DIR/$BINARY_NAME"

  # Gerar checksum SHA-256
  sha256sum "$BUILD_DIR/$BINARY_NAME" > "$BUILD_DIR/$BINARY_NAME.sha256"

  echo -e "${GREEN}✓ $BUILD_DIR/$BINARY_NAME${NC}"
  cat "$BUILD_DIR/$BINARY_NAME.sha256"
  echo ""
done

echo -e "${GREEN}=== Build completo! ===${NC}"
echo "Arquivos em: $BUILD_DIR"
echo ""
echo "Para criar um release no GitHub:"
echo "  git tag -a v${VERSION} -m 'Release v${VERSION}'"
echo "  git push origin v${VERSION}"