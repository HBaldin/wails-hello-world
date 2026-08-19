#!/bin/bash
set -e

# Script para compilar a aplicação para múltiplas plataformas
# Uso: ./scripts/build-release.sh [versão] [plataformas]
#
# Exemplos:
#   ./scripts/build-release.sh 1.0.0              # Compila tudo com versão 1.0.0
#   ./scripts/build-release.sh 1.0.0 darwin/arm64 # Apenas macOS ARM64

VERSION="${1:-0.0.0}"
PLATFORMS="${2:-darwin/arm64 darwin/amd64 linux/amd64 windows/amd64}"

# Cores para output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m'

BUILD_DIR="build/releases"
mkdir -p "$BUILD_DIR"

echo -e "${BLUE}=== Compilando Wails Hello World v${VERSION} ===${NC}"
echo "Plataformas: $PLATFORMS"
echo ""

for PLATFORM in $PLATFORMS; do
  IFS='/' read -r GOOS GOARCH <<< "$PLATFORM"

  echo -e "${BLUE}Compilando para $GOOS/$GOARCH...${NC}"

  # Determinar nome do arquivo
  if [ "$GOOS" = "windows" ]; then
    BINARY_NAME="wails-hello-world_${GOOS}_${GOARCH}.exe"
  else
    BINARY_NAME="wails-hello-world_${GOOS}_${GOARCH}"
  fi

  # Build
  wails build \
    -ldflags "-X main.version=${VERSION}" \
    -platform "${GOOS}/${GOARCH}" \
    -o "${BINARY_NAME}" \
    -skipbindings

  # Mover binário para pasta de releases
  if [ "$GOOS" = "macos" ] || [ "$GOOS" = "darwin" ]; then
    # No macOS, o binário fica dentro do .app
    APP_PATH="build/bin/wails-hello-world.app/Contents/MacOS/wails-hello-world"
    if [ -f "$APP_PATH" ]; then
      cp "$APP_PATH" "$BUILD_DIR/$BINARY_NAME"
    fi
  elif [ "$GOOS" = "linux" ]; then
    if [ -f "build/bin/wails-hello-world" ]; then
      cp "build/bin/wails-hello-world" "$BUILD_DIR/$BINARY_NAME"
      chmod +x "$BUILD_DIR/$BINARY_NAME"
    fi
  elif [ "$GOOS" = "windows" ]; then
    if [ -f "build/bin/wails-hello-world.exe" ]; then
      cp "build/bin/wails-hello-world.exe" "$BUILD_DIR/$BINARY_NAME"
    fi
  fi

  # Gerar checksum
  if command -v shasum &> /dev/null; then
    shasum -a 256 "$BUILD_DIR/$BINARY_NAME" > "$BUILD_DIR/$BINARY_NAME.sha256"
  else
    sha256sum "$BUILD_DIR/$BINARY_NAME" > "$BUILD_DIR/$BINARY_NAME.sha256"
  fi

  echo -e "${GREEN}✓ Compilado: $BUILD_DIR/$BINARY_NAME${NC}"
  cat "$BUILD_DIR/$BINARY_NAME.sha256"
  echo ""
done

echo -e "${GREEN}=== Build completo! ===${NC}"
echo "Arquivos em: $BUILD_DIR"
echo ""
echo "Para criar um release no GitHub:"
echo "  git tag -a v${VERSION} -m 'Release v${VERSION}'"
echo "  git push origin v${VERSION}"
