# 📦 Guía de Deployment - npm install -g

Esta guía explica cómo todo funciona cuando alguien ejecuta `npm install -g shopify-cli-tui` en cualquier plataforma (Windows, macOS, Linux).

---

## 🔄 Flujo Completo

```
Usuario ejecuta:
  npm install -g shopify-cli-tui
         ↓
npm descarga el paquete desde npmjs.com
         ↓
npm ejecuta automáticamente: npm run postinstall
         ↓
Se ejecuta: node scripts/install.js
         ↓
El script detecta tu SO y arquitectura
         ↓
Descarga el binario desde GitHub Releases
         ↓
Lo coloca en: npm/bin/sho (o sho.exe en Windows)
         ↓
Agrega al PATH (automático en Windows)
         ↓
✅ Comando "sho" disponible globalmente
```

---

## 🚀 Para Liberar una Nueva Versión

### 1. Compilar binarios para todas las plataformas

Ejecuta el workflow de GitHub Actions:

```bash
# Opción A: Crear un tag (dispara el workflow automáticamente)
git tag v1.2.0
git push --tags

# Opción B: Disparar manualmente desde GitHub
# Ve a: Actions → Build and Release → Run workflow
```

### 2. Esperar a que GitHub Actions compile

El workflow `.github/workflows/build-release.yml` hace esto:

1. ✅ Compila para Linux (amd64, arm64)
2. ✅ Compila para macOS (amd64, arm64)
3. ✅ Compila para Windows (amd64, 386)
4. ✅ Crea un GitHub Release con todos los binarios

**Verificar:** https://github.com/JacuXx/shopify-tui/releases

### 3. Actualizar versión en npm/package.json

```json
{
  "version": "1.2.0"
}
```

### 4. Publicar en npm

```bash
cd npm

# Cambiar versión
npm version patch

# Publicar
npm publish

# Verificar
npm view shopify-cli-tui
```

---

## 📋 Cómo Funciona npm/scripts/install.js

Este script se ejecuta automáticamente después de `npm install -g`:

### 1. Detecta tu sistema

```javascript
const platform = os.platform();  // 'win32', 'darwin', 'linux'
const arch = os.arch();          // 'x64', 'arm64', 'ia32'
```

### 2. Mapea a nombres de Go

```javascript
'win32'  → 'windows'
'x64'    → 'amd64'
'arm64'  → 'arm64'
'ia32'   → '386'
```

### 3. Descarga el binario desde GitHub Releases

```javascript
https://github.com/JacuXx/shopify-tui/releases/download/latest/
  ├── shopify-cli-windows-amd64.exe
  ├── shopify-cli-linux-amd64
  ├── shopify-cli-darwin-amd64
  └── shopify-cli-darwin-arm64
```

### 4. Lo coloca en npm/bin/

```
npm/bin/
├── sho          (Linux/macOS)
└── sho.exe      (Windows)
```

### 5. Agrega al PATH

- **Windows:** Automático (npm lo maneja)
- **macOS/Linux:** Actualiza `.bashrc` o `.zshrc`

---

## ✅ Testing Local

Para probar que funciona antes de publicar:

```bash
# 1. Compilar localmente
./build.sh

# 2. Crear un tag para disparar GitHub Actions
git tag v1.0.0-test
git push --tags

# 3. Esperar a que compile en GitHub
# https://github.com/JacuXx/shopify-tui/actions

# 4. Instalar localmente desde npm
npm install -g ./npm

# 5. Verificar
sho --version
```

---

## 🐛 Solución de Problemas

### Script no encuentra el binario en GitHub

```bash
# Verificar que existe la release:
https://github.com/JacuXx/shopify-tui/releases

# Los nombres deben ser exactos:
shopify-cli-windows-amd64.exe
shopify-cli-linux-amd64
shopify-cli-darwin-amd64
shopify-cli-darwin-arm64
```

### "command not found: sho"

```bash
# Windows: Reinicia PowerShell
# macOS/Linux:
source ~/.bashrc  # o ~/.zshrc
```

---

## 📚 Referencias

- npm package: https://www.npmjs.com/package/shopify-cli-tui
- GitHub: https://github.com/JacuXx/shopify-tui
- Releases: https://github.com/JacuXx/shopify-tui/releases
