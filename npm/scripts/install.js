const https = require('https');
const fs = require('fs');
const path = require('path');
const os = require('os');
const { execSync } = require('child_process');

const VERSION = '1.0.0';
const REPO = 'JacuXx/shopify-cli';

// Mapeo de plataforma/arquitectura a nombre de binario
function getBinaryName() {
  const platform = os.platform();
  const arch = os.arch();
  
  const platformMap = {
    'darwin': 'darwin',
    'linux': 'linux',
    'win32': 'windows'
  };
  
  const archMap = {
    'x64': 'amd64',
    'arm64': 'arm64'
  };
  
  const p = platformMap[platform];
  const a = archMap[arch];
  
  if (!p || !a) {
    console.error(`❌ Plataforma no soportada: ${platform}/${arch}`);
    process.exit(1);
  }
  
  const ext = platform === 'win32' ? '.exe' : '';
  return `shopify-cli-${p}-${a}${ext}`;
}

// Descargar archivo
function download(url, dest) {
  return new Promise((resolve, reject) => {
    const file = fs.createWriteStream(dest);
    
    const request = (url) => {
      https.get(url, (response) => {
        // Manejar redirecciones
        if (response.statusCode === 302 || response.statusCode === 301) {
          request(response.headers.location);
          return;
        }
        
        if (response.statusCode !== 200) {
          reject(new Error(`HTTP ${response.statusCode}`));
          return;
        }
        
        response.pipe(file);
        file.on('finish', () => {
          file.close();
          resolve();
        });
      }).on('error', (err) => {
        fs.unlink(dest, () => {});
        reject(err);
      });
    };
    
    request(url);
  });
}

async function install() {
  const binaryName = getBinaryName();
  const binDir = path.join(__dirname, '..', 'bin');
  const destName = os.platform() === 'win32' ? 'shopify-cli.exe' : 'shopify-cli';
  const destPath = path.join(binDir, destName);
  
  // Si ya existe el binario, no hacer nada
  if (fs.existsSync(destPath)) {
    console.log('✅ shopify-cli ya está instalado');
    return;
  }
  
  const url = `https://github.com/${REPO}/releases/download/v${VERSION}/${binaryName}`;
  
  console.log('📦 Descargando shopify-cli...');
  console.log(`   ${url}`);
  
  try {
    await download(url, destPath);
    
    // Hacer ejecutable en Unix
    if (os.platform() !== 'win32') {
      fs.chmodSync(destPath, '755');
    }
    
    console.log('✅ shopify-cli instalado correctamente!');
    console.log('');
    console.log('🚀 Ejecuta: shopify-cli');
    
  } catch (err) {
    console.error('❌ Error descargando binario:', err.message);
    console.error('');
    console.error('💡 Alternativa: instala con Go:');
    console.error('   go install github.com/JacuXx/shopify-cli@latest');
    process.exit(1);
  }
}

install();
