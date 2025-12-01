#!/bin/bash

# Script de instalación para shopify-cli
# Uso: curl -sSL <url>/install.sh | bash

set -e

REPO="JacuXx/shopify-cli"
BINARY_NAME="shopify-cli"
INSTALL_DIR="/usr/local/bin"

echo "🛒 Instalando $BINARY_NAME..."

# Detectar sistema operativo y arquitectura
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case $ARCH in
    x86_64) ARCH="amd64" ;;
    aarch64|arm64) ARCH="arm64" ;;
    *) echo "❌ Arquitectura no soportada: $ARCH"; exit 1 ;;
esac

echo "📦 Sistema: $OS/$ARCH"

# Verificar dependencias
if ! command -v go &> /dev/null; then
    echo "❌ Go no está instalado. Instálalo primero:"
    echo "   https://go.dev/dl/"
    exit 1
fi

if ! command -v shopify &> /dev/null; then
    echo "⚠️  Shopify CLI no encontrado. Instálalo con:"
    echo "   npm install -g @shopify/cli"
fi

# Clonar y compilar
TEMP_DIR=$(mktemp -d)
cd "$TEMP_DIR"

echo "📥 Descargando código fuente..."
git clone --depth 1 "https://github.com/$REPO.git" .

echo "🔨 Compilando..."
go build -o "$BINARY_NAME" .

echo "📦 Instalando en $INSTALL_DIR..."
sudo mv "$BINARY_NAME" "$INSTALL_DIR/"
sudo chmod +x "$INSTALL_DIR/$BINARY_NAME"

# Limpiar
cd -
rm -rf "$TEMP_DIR"

echo ""
echo "✅ ¡Instalación completada!"
echo ""
echo "🚀 Ejecuta: $BINARY_NAME"
echo ""
