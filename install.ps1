$ErrorActionPreference = "Stop"

$REPO = "JacuXx/shopify-tui"
$BINARY = "sho.exe"
$INSTALL_DIR = "$env:USERPROFILE\.local\bin"

# Detect architecture
$ARCH = if ([System.Environment]::Is64BitOperatingSystem) {
  $cpu = (Get-CimInstance Win32_Processor).Architecture
  if ($cpu -eq 12) { "arm64" } else { "amd64" }
} else {
  Write-Host "ERROR: Arquitectura no soportada. Se requiere 64-bit."
  exit 1
}

# Get latest version from GitHub API
Write-Host "Buscando ultima version..."
$release = Invoke-RestMethod -Uri "https://api.github.com/repos/$REPO/releases/latest" -UseBasicParsing
$VERSION = $release.tag_name -replace '^v', ''

if (-not $VERSION) {
  Write-Host "ERROR: No se pudo obtener la version mas reciente."
  exit 1
}

$BINARY_NAME = "shopify-tui-windows-$ARCH.exe"
$DOWNLOAD_URL = "https://github.com/$REPO/releases/download/v$VERSION/$BINARY_NAME"

Write-Host "Descargando sho v$VERSION para windows/$ARCH..."

if (-not (Test-Path $INSTALL_DIR)) {
  New-Item -ItemType Directory -Path $INSTALL_DIR -Force | Out-Null
}

$DEST = Join-Path $INSTALL_DIR $BINARY
Invoke-WebRequest -Uri $DOWNLOAD_URL -OutFile $DEST -UseBasicParsing

$CHECKSUM_URL = "https://github.com/$REPO/releases/download/v$VERSION/$BINARY_NAME.sha256"
try {
  $EXPECTED = (Invoke-WebRequest -Uri $CHECKSUM_URL -UseBasicParsing).Content.Trim()
} catch {
  Write-Host "ERROR: No se pudo obtener el checksum de verificacion."
  Remove-Item -Path $DEST -Force -ErrorAction SilentlyContinue
  exit 1
}
$COMPUTED = (Get-FileHash -Path $DEST -Algorithm SHA256).Hash.ToLower()
if ($COMPUTED -ne $EXPECTED) {
  Write-Host "ERROR: Verificacion de integridad fallida. El binario puede estar corrupto."
  Remove-Item -Path $DEST -Force
  exit 1
}

$USER_PATH = [System.Environment]::GetEnvironmentVariable("PATH", "User")
if ($USER_PATH -notlike "*$INSTALL_DIR*") {
  [System.Environment]::SetEnvironmentVariable("PATH", "$INSTALL_DIR;$USER_PATH", "User")
  Write-Host ""
  Write-Host "PATH configurado: $INSTALL_DIR"
  Write-Host "   Reinicia tu terminal para que tome efecto."
}

Write-Host "OK: sho v$VERSION instalado en $DEST"

$BUN_BIN = Join-Path $env:USERPROFILE ".bun\bin\sho.exe"
if (Test-Path $BUN_BIN) {
  Copy-Item -Path $DEST -Destination $BUN_BIN -Force
  Write-Host "sho actualizado tambien en $BUN_BIN"
}

Write-Host ""
Write-Host "Ejecuta: sho"
