#!/bin/bash

# Script para compilar binarios para todas las plataformas
# Ejecutar desde el directorio raíz del proyecto

set -e

VERSION="1.0.0"
OUTPUT_DIR="./releases"
BINARY_NAME="shopify-cli"

echo "🔨 Compilando shopify-cli v${VERSION} para todas las plataformas..."

# Crear directorio de salida
mkdir -p "$OUTPUT_DIR"

# Plataformas a compilar
PLATFORMS=(
  "linux/amd64"
  "linux/arm64"
  "darwin/amd64"
  "darwin/arm64"
  "windows/amd64"
)

for PLATFORM in "${PLATFORMS[@]}"; do
  OS="${PLATFORM%/*}"
  ARCH="${PLATFORM#*/}"
  
  OUTPUT_NAME="${BINARY_NAME}-${OS}-${ARCH}"
  
  if [ "$OS" = "windows" ]; then
    OUTPUT_NAME="${OUTPUT_NAME}.exe"
  fi
  
  echo "📦 Compilando para ${OS}/${ARCH}..."
  
  GOOS=$OS GOARCH=$ARCH go build -ldflags="-s -w" -o "${OUTPUT_DIR}/${OUTPUT_NAME}" .
  
  echo "   ✅ ${OUTPUT_NAME}"
done

echo ""
echo "🎉 Compilación completada!"
echo "   Binarios en: ${OUTPUT_DIR}/"
ls -la "$OUTPUT_DIR"
