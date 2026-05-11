#!/bin/sh
set -e

REPO="JacuXx/shopify-tui"
BINARY="sho"
INSTALL_DIR="/usr/local/bin"

# Detect OS
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
case "$OS" in
  darwin) PLATFORM="darwin" ;;
  linux)  PLATFORM="linux" ;;
  *)
    echo "Error: Sistema operativo no soportado: $OS"
    echo "   Plataformas soportadas: darwin (macOS), linux"
    exit 1
    ;;
esac

# Detect architecture
ARCH=$(uname -m)
case "$ARCH" in
  x86_64)  ARCH="amd64" ;;
  aarch64) ARCH="arm64" ;;
  arm64)   ARCH="arm64" ;;
  *)
    echo "Error: Arquitectura no soportada: $ARCH"
    echo "   Arquitecturas soportadas: x86_64, arm64"
    exit 1
    ;;
esac

# Get latest version from GitHub API
echo "Buscando ultima version..."
VERSION=$(curl -fsSL "https://api.github.com/repos/${REPO}/releases/latest" \
  | grep '"tag_name"' \
  | sed 's/.*"tag_name": *"v\([^"]*\)".*/\1/')

if [ -z "$VERSION" ]; then
  echo "Error: No se pudo obtener la version mas reciente."
  exit 1
fi

BINARY_NAME="shopify-tui-${PLATFORM}-${ARCH}"
DOWNLOAD_URL="https://github.com/${REPO}/releases/download/v${VERSION}/${BINARY_NAME}"

echo "Descargando sho v${VERSION} para ${PLATFORM}/${ARCH}..."

TMP_FILE=$(mktemp)
curl -fsSL "$DOWNLOAD_URL" -o "$TMP_FILE"
chmod +x "$TMP_FILE"

if [ -w "$INSTALL_DIR" ]; then
  mv "$TMP_FILE" "${INSTALL_DIR}/${BINARY}"
else
  echo "Se requieren permisos de administrador para instalar en ${INSTALL_DIR}..."
  sudo mv "$TMP_FILE" "${INSTALL_DIR}/${BINARY}"
fi

echo "sho v${VERSION} instalado en ${INSTALL_DIR}/${BINARY}"

BUN_BIN="${HOME}/.bun/bin/${BINARY}"
if [ -f "$BUN_BIN" ]; then
  cp "${INSTALL_DIR}/${BINARY}" "$BUN_BIN"
  echo "sho actualizado tambien en ${BUN_BIN}"
fi

echo ""
echo "Ejecuta: sho"
