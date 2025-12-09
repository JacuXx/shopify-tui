const fs = require('fs');
const path = require('path');
const os = require('os');

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
  const destName = os.platform() === 'win32' ? 'sho..exe' : 'sho.';
  const destPath = path.join(binDir, destName);
  
  if (!fs.existsSync(sourcePath)) {
    console.error(`❌ Binario no encontrado: ${binaryName}`);
    console.error('   Los binarios incluidos son:');
    fs.readdirSync(binDir).filter(f => f.startsWith('shopify-tui-')).forEach(f => {
      console.error(`   - ${f}`);
    });
    process.exit(1);
  }
  
  if (fs.existsSync(destPath)) {
    const sourceStats = fs.statSync(sourcePath);
    const destStats = fs.statSync(destPath);
    if (sourceStats.size === destStats.size) {
      console.log('✅ sho. ya está instalado');
      return;
    }
    fs.unlinkSync(destPath);
  }
  
  console.log(`📦 Configurando sho. para ${os.platform()}/${os.arch()}...`);
  
  try {
    fs.copyFileSync(sourcePath, destPath);
    
    if (os.platform() !== 'win32') {
      fs.chmodSync(destPath, 0o755);
    }
    
    console.log('✅ sho. instalado correctamente!');
    console.log('');
    console.log('🚀 Ejecuta: sho.');
    
  } catch (err) {
    console.error('❌ Error configurando binario:', err.message);
    process.exit(1);
  }
}

install();
