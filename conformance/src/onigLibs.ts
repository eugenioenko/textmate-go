import fs from 'node:fs';
import path from 'node:path';
import { createRequire } from 'node:module';
import { fileURLToPath } from 'node:url';
import type { IOnigLib } from 'vscode-textmate-onig-types';

const require = createRequire(import.meta.url);
let onigurumaLib: Promise<IOnigLib> | null = null;

export function getOniguruma(): Promise<IOnigLib> {
  if (!onigurumaLib) {
    const oniguruma = require('vscode-oniguruma') as typeof import('vscode-oniguruma');
    const packageDir = path.dirname(
      fileURLToPath(import.meta.resolve('vscode-oniguruma/package.json')),
    );
    const wasm = Uint8Array.from(fs.readFileSync(path.join(packageDir, 'release/onig.wasm'))).buffer;
    onigurumaLib = oniguruma.loadWASM(wasm).then(() => ({
      createOnigScanner(patterns: string[]) {
        return new oniguruma.OnigScanner(patterns);
      },
      createOnigString(value: string) {
        return new oniguruma.OnigString(value);
      },
    }));
  }
  return onigurumaLib;
}
