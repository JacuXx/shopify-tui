import fs from 'fs';
import { fileURLToPath } from 'url';
import path from 'path';
import os from 'os';
import { execSync } from 'child_process';

const __dirname = path.dirname(fileURLToPath(import.meta.url));

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

function getBunGlobalBin() {
  try {
    const result = execSync('bun pm bin -g', { encoding: 'utf8' }).trim();
    return result || null;
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
  const bunBin = getBunGlobalBin();
  if (!bunBin) return;

  if (os.platform() === 'win32') {
    if (isInPath(bunBin)) return;
    try {
      execSync(
        `powershell -Command "[System.Environment]::SetEnvironmentVariable('PATH', [System.Environment]::GetEnvironmentVariable('PATH','User') + ';${bunBin}', 'User')"`,
        { encoding: 'utf8' }
      );
      console.log('');
      console.log(`✅ PATH configurado: ${bunBin}`);
      console.log('   Reinicia tu terminal para que tome efecto.');
    } catch (e) {
      console.log('');
      console.log(`⚠️  Agrega manualmente a tu PATH de usuario: ${bunBin}`);
    }
    return;
  }

  if (isInPath(bunBin)) {
    return;
  }

  const configFile = getShellConfigFile();
  if (!configFile) {
    console.log('');
    console.log('⚠️  El directorio de bun no está en tu PATH.');
    console.log(`   Agrega esto a tu archivo de configuración del shell:`);
    console.log(`   export PATH="${bunBin}:$PATH"`);
    return;
  }

  const exportLine = `export PATH="${bunBin}:$PATH"`;

  try {
    let configContent = '';
    if (fs.existsSync(configFile)) {
      configContent = fs.readFileSync(configFile, 'utf8');
    }

    if (configContent.includes(bunBin)) {
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
    try {
      fs.unlinkSync(destPath);
    } catch (e) {
      if (os.platform() === 'win32' && e.code === 'EPERM') {
        const oldPath = destPath + '.old';
        try { fs.unlinkSync(oldPath); } catch (_) {}
        fs.renameSync(destPath, oldPath);
      } else {
        throw e;
      }
    }
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
