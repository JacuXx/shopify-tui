# 📦 Guía de Instalación - shopify-cli

Esta guía cubre cómo instalar `shopify-cli` en diferentes sistemas operativos.

---

## 🪟 Windows (PowerShell)

### Opción 1: Script Automático (Recomendado)

```powershell
# 1. Abre PowerShell como Administrador
# Presiona Win + X y selecciona "Windows PowerShell (Admin)"

# 2. Ejecuta:
powershell -ExecutionPolicy Bypass -File https://raw.githubusercontent.com/JacuXx/shopify-tui/main/install.ps1
```

**Lo que hace el script:**
- ✅ Detecta tu arquitectura (amd64 o 386)
- ✅ Descarga el binario más reciente desde GitHub Releases
- ✅ Lo instala en `C:\Program Files\shopify-cli\`
- ✅ Agrega el directorio al PATH del sistema
- ✅ Crea un backup si ya existe una versión anterior

### Opción 2: Descarga Manual

1. Ve a [GitHub Releases](https://github.com/JacuXx/shopify-tui/releases)
2. Descarga `shopify-cli-windows-amd64.exe`
3. Colócalo en: `C:\Program Files\shopify-cli\`
4. Agrega `C:\Program Files\shopify-cli\` al PATH del sistema

### Opción 3: Compilar desde Fuente

```powershell
# Requiere Git y Go 1.22+

git clone https://github.com/JacuXx/shopify-tui.git
cd shopify-tui
go build -o shopify-cli.exe .

# Copia a un directorio en el PATH:
move shopify-cli.exe "C:\Program Files\shopify-cli\"
```

---

## 🍎 macOS

### Opción 1: Script de Instalación

```bash
curl -fsSL https://raw.githubusercontent.com/JacuXx/shopify-tui/main/install.sh | bash
```

### Opción 2: Homebrew (Si lo agregas después)

```bash
brew install JacuXx/shopify-tui/shopify-cli
```

### Opción 3: Compilar desde Fuente

```bash
git clone https://github.com/JacuXx/shopify-tui.git
cd shopify-tui
go build -o shopify-cli .
sudo mv shopify-cli /usr/local/bin/
```

---

## 🐧 Linux

### Opción 1: Script de Instalación

```bash
curl -fsSL https://raw.githubusercontent.com/JacuXx/shopify-tui/main/install.sh | bash
```

Alternativamente:
```bash
wget -qO- https://raw.githubusercontent.com/JacuXx/shopify-tui/main/install.sh | bash
```

### Opción 2: Compilar desde Fuente

```bash
git clone https://github.com/JacuXx/shopify-tui.git
cd shopify-tui
go build -o shopify-cli .
sudo mv shopify-cli /usr/local/bin/
```

---

## ✅ Verificar la Instalación

Después de instalar, verifica que funcione:

```bash
# Mostrar versión
shopify-cli --version

# Mostrar ayuda
shopify-cli --help
```

Si el comando no se encuentra:
- **Windows:** Reinicia PowerShell o CMD
- **macOS/Linux:** Reinicia la terminal o ejecuta: `source ~/.bashrc` (o `~/.zshrc`)

---

## 📋 Requisitos Previos

Antes de instalar `shopify-cli`, necesitas:

1. **Shopify CLI** - [Instalación](https://shopify.dev/docs/api/shopify-cli)
   ```bash
   npm install -g @shopify/cli
   ```

2. **Node.js & npm** (para Shopify CLI)
   - Descargable desde: https://nodejs.org/

3. **Go** (solo si compiles desde fuente)
   - Versión mínima: 1.22
   - Descargable desde: https://go.dev/dl/

---

## 🔄 Actualización

### Windows (PowerShell)
```powershell
powershell -ExecutionPolicy Bypass -File https://raw.githubusercontent.com/JacuXx/shopify-tui/main/install.ps1
```

### macOS / Linux
```bash
curl -fsSL https://raw.githubusercontent.com/JacuXx/shopify-tui/main/install.sh | bash
```

---

## 🗑️ Desinstalación

### Windows
```powershell
# Opción 1: Borrar el directorio
Remove-Item "C:\Program Files\shopify-cli" -Recurse -Force

# Opción 2: Quitar del PATH manualmente
# Busca en: Settings > System > About > Advanced System Settings > Environment Variables
```

### macOS / Linux
```bash
sudo rm /usr/local/bin/shopify-cli
```

---

## 🐛 Solución de Problemas

### "Command not found: shopify-cli"

**Windows:**
- Abre PowerShell como Administrador
- Ejecuta: `$env:Path`
- Verifica que `C:\Program Files\shopify-cli` esté en la lista
- Si no está, vuelve a ejecutar el instalador

**macOS/Linux:**
```bash
# Verifica dónde se instaló
which shopify-cli

# Si no existe, instálalo nuevamente:
curl -fsSL https://raw.githubusercontent.com/JacuXx/shopify-tui/main/install.sh | bash
```

### "Permission denied" (Linux/macOS)

```bash
# Asegúrate de tener permisos:
chmod +x /usr/local/bin/shopify-cli

# O reinstala con permisos correctos:
sudo curl -fsSL https://raw.githubusercontent.com/JacuXx/shopify-tui/main/install.sh | bash
```

### El script PowerShell no ejecuta en Windows

Si ves un error sobre políticas de ejecución:

```powershell
# Abre PowerShell como Administrador y ejecuta:
Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser

# Luego intenta de nuevo:
powershell -ExecutionPolicy Bypass -File https://raw.githubusercontent.com/JacuXx/shopify-tui/main/install.ps1
```

---

## 📚 Más Información

- **GitHub:** https://github.com/JacuXx/shopify-tui
- **Releases:** https://github.com/JacuXx/shopify-tui/releases
- **Issues:** https://github.com/JacuXx/shopify-tui/issues
