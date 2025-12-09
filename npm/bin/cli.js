#!/usr/bin/env node

const { spawn } = require('child_process');
const path = require('path');
const os = require('os');
const fs = require('fs');

function getBinaryPath() {
  const platform = os.platform();
  const arch = os.arch();
  
  let binaryName = 'sho.';
  
  if (platform === 'win32') {
    binaryName = 'sho..exe';
  }
  
  const binaryPath = path.join(__dirname, binaryName);
  
  if (!fs.existsSync(binaryPath)) {
    console.error('❌ Binario no encontrado:', binaryPath);
    console.error('   Por favor reinstala el paquete: npm install -g shopify-cli-tui');
    process.exit(1);
  }
  
  return binaryPath;
}

const binaryPath = getBinaryPath();
const child = spawn(binaryPath, process.argv.slice(2), {
  stdio: 'inherit',
  env: process.env
});

child.on('error', (err) => {
  console.error('❌ Error al ejecutar sho.:', err.message);
  process.exit(1);
});

child.on('exit', (code) => {
  process.exit(code || 0);
});
