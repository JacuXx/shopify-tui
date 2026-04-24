const fs = require('fs');
const path = require('path');
const os = require('os');
const https = require('https');

function getPackageVersion() {
  const packageJsonPath = path.join(__dirname, '..', 'package.json');
  const packageJson = JSON.parse(fs.readFileSync(packageJsonPath, 'utf8'));
  return packageJson.version;
}

function getBinaryName() {
  const platform = os.platform();
  const arch = os.arch();

  // Mapear arquitectura Node.js a Go
  const archMap = {
    'x64': 'amd64',
    'arm64': 'arm64',
    'ia32': '386'
  };

  const platformMap = {
    'darwin': 'darwin',
    'linux': 'linux',
    'win32': 'windows'
  };

  const goArch = archMap[arch];
  const goOS = platformMap[platform];

  if (!goOS || !goArch) {
    console.error(`❌ Plataforma no soportada: ${platform}/${arch}`);
    console.error('   Plataformas soportadas: linux, darwin (macOS), windows');
    console.error('   Arquitecturas soportadas: x64, arm64, ia32');
    process.exit(1);
  }

  const ext = platform === 'win32' ? '.exe' : '';
  return `shopify-cli-${goOS}-${goArch}${ext}`;
}

function downloadBinary(downloadUrl, destPath) {
  return new Promise((resolve, reject) => {
    const file = fs.createWriteStream(destPath);
    https.get(downloadUrl, (response) => {
      if (response.statusCode === 302 || response.statusCode === 301) {
        // Seguir redirecciones
        downloadBinary(response.headers.location, destPath).then(resolve).catch(reject);
        return;
      }
      if (response.statusCode !== 200) {
        reject(new Error(`HTTP ${response.statusCode}: ${response.statusMessage}`));
        return;
      }
      response.pipe(file);
      file.on('finish', () => {
        file.close();
        resolve();
      });
    }).on('error', (err) => {
      fs.unlink(destPath, () => {});
      reject(err);
    });
  });
}

function cleanOldBinaries(binDir) {
  const oldNames = ['shopify-cli', 'shopify-cli.exe', 'sho.', 'sho..exe'];
  oldNames.forEach(name => {
    const oldPath = path.join(binDir, name);
    if (fs.existsSync(oldPath)) {
      try {
        fs.unlinkSync(oldPath);
        console.log(`🧹 Eliminado binario viejo: ${name}`);
      } catch (e) {}
    }
  });
}

function getNpmGlobalBin() {
  try {
    const prefix = execSync('npm config get prefix', { encoding: 'utf8' }).trim();
    return path.join(prefix, 'bin');
  } catch (e) {
    return null;
  }
}

function isInPath(dir) {
  const pathEnv = process.env.PATH || '';
  return pathEnv.split(path.delimiter).includes(dir);
}

function getShellConfigFile() {
  const shell = process.env.SHELL || '';
  const home = os.homedir();
  
  if (shell.includes('zsh')) {
    return path.join(home, '.zshrc');
  } else if (shell.includes('bash')) {
    if (os.platform() === 'darwin') {
      return path.join(home, '.bash_profile');
    }
    return path.join(home, '.bashrc');
  }
  return null;
}

function setupPath() {
  if (os.platform() === 'win32') return;
  
  const npmBin = getNpmGlobalBin();
  if (!npmBin) return;
  
  if (isInPath(npmBin)) {
    return;
  }
  
  const configFile = getShellConfigFile();
  if (!configFile) {
    console.log('');
    console.log('⚠️  El directorio de npm no está en tu PATH.');
    console.log(`   Agrega esto a tu archivo de configuración del shell:`);
    console.log(`   export PATH="${npmBin}:$PATH"`);
    return;
  }
  
  const exportLine = `export PATH="${npmBin}:$PATH"`;
  
  try {
    let configContent = '';
    if (fs.existsSync(configFile)) {
      configContent = fs.readFileSync(configFile, 'utf8');
    }
    
    if (configContent.includes(npmBin)) {
      return;
    }
    
    fs.appendFileSync(configFile, `\n# Agregado por shopify-cli-tui\n${exportLine}\n`);
    console.log('');
    console.log(`✅ PATH configurado automáticamente en ${path.basename(configFile)}`);
    console.log('   Reinicia tu terminal o ejecuta:');
    console.log(`   source ${configFile}`);
    
  } catch (err) {
    console.log('');
    console.log('⚠️  No se pudo configurar el PATH automáticamente.');
    console.log(`   Agrega esta línea a ${configFile}:`);
    console.log(`   ${exportLine}`);
  }
}

async function install() {
  const binaryName = getBinaryName();
  const binDir = path.join(__dirname, '..', 'bin');

  // Crear directorio bin si no existe
  if (!fs.existsSync(binDir)) {
    fs.mkdirSync(binDir, { recursive: true });
  }

  const destName = os.platform() === 'win32' ? 'sho.exe' : 'sho';
  const destPath = path.join(binDir, destName);

  cleanOldBinaries(binDir);

  if (fs.existsSync(destPath)) {
    fs.unlinkSync(destPath);
  }

  console.log(`📦 Configurando sho para ${os.platform()}/${os.arch()}...`);

  try {
    // Intentar descargar desde GitHub Releases
    console.log(`📥 Descargando ${binaryName} desde GitHub...`);
    const version = getPackageVersion();
    const downloadUrl = `https://github.com/JacuXx/shopify-tui/releases/download/v${version}/${binaryName}`;

    const tempPath = path.join(binDir, `${binaryName}.tmp`);

    await downloadBinary(downloadUrl, tempPath);
    fs.renameSync(tempPath, destPath);

    if (os.platform() !== 'win32') {
      fs.chmodSync(destPath, 0o755);
    }

    console.log('✅ sho instalado correctamente!');

    setupPath();

    console.log('');
    console.log('🚀 Ejecuta: sho');

  } catch (err) {
    console.error('❌ Error instalando binario:', err.message);
    console.error('');
    console.error('💡 Solución alternativa - Compilar desde fuente:');
    console.error('   1. git clone https://github.com/JacuXx/shopify-tui.git');
    console.error('   2. cd shopify-tui');
    console.error('   3. go build -o sho .');
    console.error('   4. Copia el binario a:', destPath);
    process.exit(1);
  }
}

install();
