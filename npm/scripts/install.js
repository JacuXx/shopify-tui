const fs = require('fs');
const path = require('path');
const os = require('os');
const { execSync } = require('child_process');

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

function install() {
  const binaryName = getBinaryName();
  const binDir = path.join(__dirname, '..', 'bin');
  const sourcePath = path.join(binDir, binaryName);
  const destName = os.platform() === 'win32' ? 'sho.exe' : 'sho';
  const destPath = path.join(binDir, destName);
  
  if (!fs.existsSync(sourcePath)) {
    console.error(`❌ Binario no encontrado: ${binaryName}`);
    console.error('   Los binarios incluidos son:');
    fs.readdirSync(binDir).filter(f => f.startsWith('shopify-tui-')).forEach(f => {
      console.error(`   - ${f}`);
    });
    process.exit(1);
  }
  
  cleanOldBinaries(binDir);
  
  if (fs.existsSync(destPath)) {
    fs.unlinkSync(destPath);
  }
  
  console.log(`📦 Configurando sho para ${os.platform()}/${os.arch()}...`);
  
  try {
    fs.copyFileSync(sourcePath, destPath);
    
    if (os.platform() !== 'win32') {
      fs.chmodSync(destPath, 0o755);
    }
    
    console.log('✅ sho instalado correctamente!');
    
    setupPath();
    
    console.log('');
    console.log('🚀 Ejecuta: sho');
    
  } catch (err) {
    console.error('❌ Error configurando binario:', err.message);
    process.exit(1);
  }
}

install();
