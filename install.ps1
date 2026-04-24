# PowerShell script para instalar shopify-cli en Windows
# Ejecución: powershell -ExecutionPolicy Bypass -File install.ps1

param(
    [switch]$Force = $false
)

# Configuración
$REPO = "JacuXx/shopify-cli"
$BINARY_NAME = "shopify-cli.exe"
$INSTALL_DIR = "$env:PROGRAMFILES\shopify-cli"
$VERSION = "latest"

Write-Host "🛒 Instalando shopify-cli..." -ForegroundColor Cyan

# Verificar permisos de administrador
$isAdmin = ([Security.Principal.WindowsPrincipal] [Security.Principal.WindowsIdentity]::GetCurrent()).IsInRole([Security.Principal.WindowsBuiltInRole] 'Administrator')
if (-not $isAdmin) {
    Write-Host "❌ Este script requiere ejecutarse como Administrador" -ForegroundColor Red
    Write-Host "   Abre PowerShell como administrador y vuelve a ejecutar:" -ForegroundColor Yellow
    Write-Host "   powershell -ExecutionPolicy Bypass -File install.ps1" -ForegroundColor Yellow
    exit 1
}

# Detectar arquitectura
$ARCH = if ([System.Environment]::Is64BitProcess) { "amd64" } else { "386" }
Write-Host "📦 Sistema: windows/$ARCH" -ForegroundColor Green

# Verificar si Go está instalado (opcional, pero recomendado)
$goInstalled = $null -ne (Get-Command go -ErrorAction SilentlyContinue)
if (-not $goInstalled) {
    Write-Host "⚠️  Go no está instalado. Para compilar desde fuente:" -ForegroundColor Yellow
    Write-Host "   https://go.dev/dl/" -ForegroundColor Yellow
    Write-Host "   O descarga un binario pre-compilado desde:" -ForegroundColor Yellow
    Write-Host "   https://github.com/$REPO/releases" -ForegroundColor Yellow
}

# Verificar si Shopify CLI está instalado
$shopifyInstalled = $null -ne (Get-Command shopify -ErrorAction SilentlyContinue)
if (-not $shopifyInstalled) {
    Write-Host "⚠️  Shopify CLI no encontrado. Instálalo con:" -ForegroundColor Yellow
    Write-Host "   npm install -g @shopify/cli" -ForegroundColor Yellow
}

# Crear directorio de instalación si no existe
if (-not (Test-Path $INSTALL_DIR)) {
    New-Item -ItemType Directory -Path $INSTALL_DIR -Force | Out-Null
    Write-Host "📁 Creado directorio: $INSTALL_DIR" -ForegroundColor Green
}

# Opción 1: Descargar binario pre-compilado
Write-Host "`n📥 Buscando binarios pre-compilados..." -ForegroundColor Cyan

try {
    # Obtener la última release
    $releases = Invoke-RestMethod -Uri "https://api.github.com/repos/$REPO/releases/latest" -ErrorAction Stop
    $downloadUrl = $releases.assets | Where-Object { $_.name -match "windows-amd64" } | Select-Object -First 1 -ExpandProperty browser_download_url

    if ($downloadUrl) {
        Write-Host "✅ Encontrado binario: $downloadUrl" -ForegroundColor Green

        $tempFile = "$env:TEMP\shopify-cli.exe"
        Write-Host "📥 Descargando..." -ForegroundColor Cyan

        Invoke-WebRequest -Uri $downloadUrl -OutFile $tempFile -UseBasicParsing

        $finalPath = "$INSTALL_DIR\$BINARY_NAME"

        # Si ya existe, hacer backup
        if (Test-Path $finalPath) {
            $backupPath = "$INSTALL_DIR\shopify-cli.exe.bak"
            Copy-Item -Path $finalPath -Destination $backupPath -Force
            Write-Host "💾 Backup anterior guardado en: $backupPath" -ForegroundColor Yellow
        }

        # Mover el binario
        Move-Item -Path $tempFile -Destination $finalPath -Force
        Write-Host "✅ Binario instalado en: $finalPath" -ForegroundColor Green
    }
    else {
        Write-Host "❌ No se encontró binario para windows/amd64" -ForegroundColor Red
        Write-Host "   Compila manualmente con: .\build.sh" -ForegroundColor Yellow
        exit 1
    }
}
catch {
    Write-Host "❌ Error descargando: $_" -ForegroundColor Red
    Write-Host "   Intenta compilar manualmente:" -ForegroundColor Yellow
    Write-Host "   1. go build -o shopify-cli.exe ." -ForegroundColor Yellow
    Write-Host "   2. Copia el .exe a: $INSTALL_DIR" -ForegroundColor Yellow
    exit 1
}

# Agregar al PATH si no está
$pathEnvVar = [Environment]::GetEnvironmentVariable("Path", "Machine")
if ($pathEnvVar -notlike "*$INSTALL_DIR*") {
    Write-Host "`n🔧 Agregando $INSTALL_DIR al PATH..." -ForegroundColor Cyan
    $newPath = "$pathEnvVar;$INSTALL_DIR"
    [Environment]::SetEnvironmentVariable("Path", $newPath, "Machine")
    Write-Host "✅ PATH actualizado" -ForegroundColor Green
    Write-Host "⚠️  Reinicia PowerShell o CMD para que los cambios surtan efecto" -ForegroundColor Yellow
}
else {
    Write-Host "✅ $INSTALL_DIR ya está en el PATH" -ForegroundColor Green
}

Write-Host "`n✅ ¡Instalación completada!" -ForegroundColor Green
Write-Host ""
Write-Host "🚀 Para ejecutar:" -ForegroundColor Cyan
Write-Host "   shopify-cli" -ForegroundColor White
Write-Host ""
Write-Host "📖 Más info: https://github.com/$REPO" -ForegroundColor Blue
Write-Host ""
