const fs = require('fs');
const path = require('path');
const os = require('os');

// Mapeo de plataforma/arquitectura a nombre de binario incluido
function getBinaryName() {
  const platform = os.platform();
  const arch = os.arch();
  
  const platformMap = {
    'darwin': 'darwin',
    'linux': 'linux',
    'win32': 'win32'
  };
  
  const archMap = {
    'x64': 'x64',
    'arm64': 'arm64'
  };
  
  const p = platformMap[platform];
  const a = archMap[arch];
  
  if (!p || !a) {
    console.error(`❌ Plataforma no soportada: ${platform}/${arch}`);
    console.error('   Plataformas soportadas: linux, darwin (macOS), win32');
    console.error('   Arquitecturas soportadas: x64, arm64');
    process.exit(1);
  }
  
  const ext = platform === 'win32' ? '.exe' : '';
  return `shopify-tui-${p}-${a}${ext}`;
}

function install() {
  const binaryName = getBinaryName();
  const binDir = path.join(__dirname, '..', 'bin');
  const sourcePath = path.join(binDir, binaryName);
  const destName = os.platform() === 'win32' ? 'shopify-cli.exe' : 'shopify-cli';
  const destPath = path.join(binDir, destName);
  
  // Verificar que existe el binario para esta plataforma
  if (!fs.existsSync(sourcePath)) {
    console.error(`❌ Binario no encontrado: ${binaryName}`);
    console.error('   Los binarios incluidos son:');
    fs.readdirSync(binDir).filter(f => f.startsWith('shopify-tui-')).forEach(f => {
      console.error(`   - ${f}`);
    });
    process.exit(1);
  }
  
  // Si ya existe el destino, verificar si es el mismo
  if (fs.existsSync(destPath)) {
    const sourceStats = fs.statSync(sourcePath);
    const destStats = fs.statSync(destPath);
    if (sourceStats.size === destStats.size) {
      console.log('✅ shopify-cli ya está instalado');
      return;
    }
    // Si son diferentes, eliminar el viejo
    fs.unlinkSync(destPath);
  }
  
  console.log(`📦 Configurando shopify-cli para ${os.platform()}/${os.arch()}...`);
  
  try {
    // Copiar el binario al nombre final
    fs.copyFileSync(sourcePath, destPath);
    
    // Hacer ejecutable en Unix
    if (os.platform() !== 'win32') {
      fs.chmodSync(destPath, 0o755);
    }
    
    console.log('✅ shopify-cli instalado correctamente!');
    console.log('');
    console.log('🚀 Ejecuta: shopify-cli');
    
  } catch (err) {
    console.error('❌ Error configurando binario:', err.message);
    process.exit(1);
  }
}

install();
