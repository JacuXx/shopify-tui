#!/usr/bin/env bun

import { execFileSync } from 'child_process';
import { fileURLToPath } from 'url';
import path from 'path';
import os from 'os';
import fs from 'fs';

const __dirname = path.dirname(fileURLToPath(import.meta.url));

function getBinaryPath() {
  const platform = os.platform();

  let binaryName = 'sho';

  if (platform === 'win32') {
    binaryName = 'sho.exe';
  }

  const binaryPath = path.join(__dirname, binaryName);

  if (!fs.existsSync(binaryPath)) {
    console.error('❌ Binario no encontrado:', binaryPath);
    console.error('   Por favor reinstala el paquete: bun add -g shopify-cli-tui');
    process.exit(1);
  }

  return binaryPath;
}

const binaryPath = getBinaryPath();

try {
  execFileSync(binaryPath, process.argv.slice(2), {
    stdio: 'inherit',
    env: process.env
  });
} catch (err) {
  if (err.status !== null) {
    process.exit(err.status);
  }
  console.error('❌ Error al ejecutar sho:', err.message);
  process.exit(1);
}
