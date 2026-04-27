#!/usr/bin/env bun

import { execFileSync } from 'child_process';
import { fileURLToPath } from 'url';
import path from 'path';
import os from 'os';
import fs from 'fs';

const __dirname: string = path.dirname(fileURLToPath(import.meta.url));

function getBinaryPath(): string {
  const platform: string = os.platform();

  let binaryName: string = 'sho';

  if (platform === 'win32') {
    binaryName = 'sho.exe';
  }

  const binaryPath: string = path.join(__dirname, binaryName);

  if (!fs.existsSync(binaryPath)) {
    console.error('❌ Binario no encontrado:', binaryPath);
    console.error('   Por favor reinstala el paquete: bun add -g @jsr/jacuxx__shopify-tui');
    process.exit(1);
  }

  return binaryPath;
}

const binaryPath: string = getBinaryPath();

try {
  execFileSync(binaryPath, process.argv.slice(2), {
    stdio: 'inherit',
    env: process.env
  });
} catch (err: unknown) {
  const error = err as { status?: number | null; message?: string };
  if (error.status !== null && error.status !== undefined) {
    process.exit(error.status);
  }
  console.error('❌ Error al ejecutar sho:', error.message);
  process.exit(1);
}
