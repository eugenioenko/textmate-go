import { randomBytes } from 'node:crypto';
import { execFileSync, spawn } from 'node:child_process';
import fs from 'node:fs/promises';
import os from 'node:os';
import path from 'node:path';
import { fileURLToPath } from 'node:url';

const here = path.dirname(fileURLToPath(import.meta.url));
const conformanceDir = path.resolve(here, '..');
const pinnedRevision = 'fbe49961ab8077e587fdf5282019655ae69e5f9e';
const markerName = '.textmate-go-theme-isolation.json';
const temporaryPrefix = 'textmate-go-theme-';
const defaultTimeoutMs = 20 * 60 * 1000;

function parseTimeout(value) {
  if (value === undefined) return defaultTimeoutMs;
  if (!/^[0-9]+$/.test(value)) throw new Error('TEXTMATE_THEME_TIMEOUT_MS must be an integer');
  const timeout = Number(value);
  if (!Number.isSafeInteger(timeout) || timeout < 1000 || timeout > 60 * 60 * 1000) {
    throw new Error('TEXTMATE_THEME_TIMEOUT_MS must be between 1000 and 3600000');
  }
  return timeout;
}

async function validateSourceCheckout(sourceRoot) {
  const resolved = await fs.realpath(sourceRoot);
  const revision = execFileSync('git', ['-C', resolved, 'rev-parse', 'HEAD'], {
    encoding: 'utf8',
    stdio: ['ignore', 'pipe', 'pipe'],
  }).trim();
  if (revision !== pinnedRevision) {
    throw new Error(`vscode-textmate revision ${revision} does not match pinned ${pinnedRevision}`);
  }
  const status = execFileSync('git', ['-C', resolved, 'status', '--porcelain', '--untracked-files=all'], {
    encoding: 'utf8',
    stdio: ['ignore', 'pipe', 'pipe'],
  });
  if (status !== '') {
    throw new Error('vscode-textmate checkout must be clean before creating the theme baseline copy');
  }
  for (const relative of ['src/tests/themes.test.ts', 'test-cases/themes/tests']) {
    const candidate = await fs.realpath(path.join(resolved, relative));
    if (!candidate.startsWith(`${resolved}${path.sep}`)) {
      throw new Error(`Source checkout content escapes through a symlink: ${relative}`);
    }
  }
  return resolved;
}

async function createIsolation(sourceRoot) {
  const temporaryParent = await fs.realpath(os.tmpdir());
  const root = await fs.mkdtemp(path.join(temporaryParent, temporaryPrefix));
  const token = randomBytes(32).toString('hex');
  try {
    await Promise.all(
      ['src', 'test-cases'].map((relative) =>
        fs.cp(path.join(sourceRoot, relative), path.join(root, relative), {
          recursive: true,
          dereference: true,
          preserveTimestamps: true,
        }),
      ),
    );
    const marker = { schema: 1, token, root, source: sourceRoot, revision: pinnedRevision };
    await fs.writeFile(path.join(root, markerName), `${JSON.stringify(marker)}\n`, {
      encoding: 'utf8',
      flag: 'wx',
      mode: 0o600,
    });
    return { root, token, temporaryParent };
  } catch (error) {
    await fs.rm(root, { recursive: true, force: true });
    throw error;
  }
}

async function cleanupIsolation(isolation) {
  const root = await fs.realpath(isolation.root);
  if (
    root !== isolation.root ||
    path.dirname(root) !== isolation.temporaryParent ||
    !path.basename(root).startsWith(temporaryPrefix)
  ) {
    throw new Error(`Refusing to clean unvalidated path ${JSON.stringify(isolation.root)}`);
  }
  const markerPath = path.join(root, markerName);
  const markerStat = await fs.lstat(markerPath);
  const marker = JSON.parse(await fs.readFile(markerPath, 'utf8'));
  if (
    !markerStat.isFile() ||
    markerStat.isSymbolicLink() ||
    marker.schema !== 1 ||
    marker.root !== root ||
    marker.token !== isolation.token ||
    marker.revision !== pinnedRevision
  ) {
    throw new Error(`Refusing to clean theme isolation with an invalid marker: ${root}`);
  }
  await fs.rm(root, { recursive: true, force: false, maxRetries: 3 });
}

function terminateProcessGroup(child, signal) {
  if (!child.pid) return;
  try {
    if (process.platform === 'win32') {
      spawn('taskkill', ['/pid', String(child.pid), '/t', '/f'], { stdio: 'ignore' });
    } else {
      process.kill(-child.pid, signal);
    }
  } catch (error) {
    if (error?.code !== 'ESRCH') throw error;
  }
}

async function runVitest(isolation, backend, timeoutMs) {
  const vitest = path.join(conformanceDir, 'node_modules/vitest/vitest.mjs');
  const child = spawn(process.execPath, [vitest, 'run', '--config', 'vitest.config.ts'], {
    cwd: conformanceDir,
    detached: process.platform !== 'win32',
    stdio: 'inherit',
    env: {
      ...process.env,
      TEXTMATE_BACKEND: backend,
      TEXTMATE_SUITE: 'themes',
      TEXTMATE_THEME_ISOLATION_ROOT: isolation.root,
      TEXTMATE_THEME_ISOLATION_TOKEN: isolation.token,
      VSCODE_TEXTMATE_DIR: isolation.root,
    },
  });

  let timedOut = false;
  let hardKillTimer;
  const timer = setTimeout(() => {
    timedOut = true;
    process.stderr.write(`Theme conformance exceeded ${timeoutMs}ms; terminating its process group.\n`);
    terminateProcessGroup(child, 'SIGTERM');
    if (process.platform !== 'win32') {
      hardKillTimer = setTimeout(() => terminateProcessGroup(child, 'SIGKILL'), 2000);
    }
  }, timeoutMs);

  let exit;
  try {
    exit = await new Promise((resolve, reject) => {
      child.once('error', reject);
      child.once('close', (code, signal) => resolve({ code, signal }));
    });
  } finally {
    clearTimeout(timer);
    if (hardKillTimer) clearTimeout(hardKillTimer);
  }
  if (timedOut) return 124;
  if (exit.signal) {
    process.stderr.write(`Theme conformance exited from signal ${exit.signal}.\n`);
    return 1;
  }
  return exit.code ?? 1;
}

async function main() {
  const backend = process.env.TEXTMATE_THEME_BACKEND ?? 'go';
  if (backend !== 'reference' && backend !== 'go') {
    throw new Error('TEXTMATE_THEME_BACKEND must be "reference" or "go"');
  }
  const timeoutMs = parseTimeout(process.env.TEXTMATE_THEME_TIMEOUT_MS);
  const sourceRoot = await validateSourceCheckout(
    path.resolve(process.env.VSCODE_TEXTMATE_DIR ?? path.join(conformanceDir, '../../vscode-textmate')),
  );
  const isolation = await createIsolation(sourceRoot);
  process.stderr.write(`Theme conformance is isolated in ${isolation.root}\n`);

  let exitCode;
  try {
    exitCode = await runVitest(isolation, backend, timeoutMs);
  } finally {
    await cleanupIsolation(isolation);
  }
  process.exitCode = exitCode;
}

main().catch((error) => {
  process.stderr.write(`${error?.stack ?? error}\n`);
  process.exitCode = 1;
});
