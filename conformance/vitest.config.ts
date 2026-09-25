import { timingSafeEqual } from 'node:crypto';
import fs from 'node:fs';
import os from 'node:os';
import path from 'node:path';
import { fileURLToPath } from 'node:url';
import { defineConfig, type Plugin } from 'vitest/config';

const here = path.dirname(fileURLToPath(import.meta.url));
const vscodeTextmateDir = path.resolve(
  process.env.VSCODE_TEXTMATE_DIR ?? path.join(here, '../../vscode-textmate'),
);

const backend = process.env.TEXTMATE_BACKEND ?? 'reference';
if (backend !== 'reference' && backend !== 'go') {
  throw new Error(`Invalid TEXTMATE_BACKEND ${JSON.stringify(backend)}; expected "reference" or "go"`);
}

type SuiteMode = 'all' | 'tokenization' | 'themes';
const requestedSuite = process.env.TEXTMATE_SUITE;
const defaultSuite: SuiteMode =
  backend === 'go' || process.env.TEXTMATE_TOKENIZATION_ONLY === '1' ? 'tokenization' : 'all';
const suiteMode = requestedSuite ?? defaultSuite;
if (suiteMode !== 'all' && suiteMode !== 'tokenization' && suiteMode !== 'themes') {
  throw new Error(
    `Invalid TEXTMATE_SUITE ${JSON.stringify(suiteMode)}; expected "all", "tokenization", or "themes"`,
  );
}
if (backend === 'go' && suiteMode === 'all') {
  throw new Error('The Go backend must use the tokenization or isolated themes suite');
}
const includeAllThemeTests = process.env.TEXTMATE_THEME_ALL_TESTS === '1';
if (includeAllThemeTests && (suiteMode !== 'themes' || backend !== 'reference')) {
  throw new Error('TEXTMATE_THEME_ALL_TESTS is only supported by the isolated reference theme suite');
}

const pinnedRevision = 'fbe49961ab8077e587fdf5282019655ae69e5f9e';
const isolationMarkerName = '.textmate-go-theme-isolation.json';

function validateThemeIsolation(): void {
  const configuredRoot = process.env.TEXTMATE_THEME_ISOLATION_ROOT;
  const configuredToken = process.env.TEXTMATE_THEME_ISOLATION_TOKEN;
  if (!configuredRoot || !configuredToken || !/^[a-f0-9]{64}$/.test(configuredToken)) {
    throw new Error('Theme conformance requires the isolated runner and its one-time marker');
  }

  const root = fs.realpathSync.native(configuredRoot);
  const checkout = fs.realpathSync.native(vscodeTextmateDir);
  const temporaryParent = fs.realpathSync.native(os.tmpdir());
  if (root !== checkout) {
    throw new Error('VSCODE_TEXTMATE_DIR must be the marked theme-isolation root');
  }
  if (path.dirname(root) !== temporaryParent || !path.basename(root).startsWith('textmate-go-theme-')) {
    throw new Error('Theme-isolation root is not an exact mkdtemp child of the system temporary directory');
  }

  const markerPath = path.join(root, isolationMarkerName);
  const markerStat = fs.lstatSync(markerPath);
  if (!markerStat.isFile() || markerStat.isSymbolicLink()) {
    throw new Error('Theme-isolation marker must be a regular file');
  }
  const marker = JSON.parse(fs.readFileSync(markerPath, 'utf8')) as Record<string, unknown>;
  const markerToken = typeof marker.token === 'string' ? marker.token : '';
  const tokensMatch =
    markerToken.length === configuredToken.length &&
    timingSafeEqual(Buffer.from(markerToken), Buffer.from(configuredToken));
  if (
    marker.schema !== 1 ||
    marker.root !== root ||
    marker.revision !== pinnedRevision ||
    !tokensMatch
  ) {
    throw new Error('Theme-isolation marker failed validation');
  }

  for (const relative of ['src/tests/themes.test.ts', 'test-cases/themes/tests']) {
    const candidate = fs.realpathSync.native(path.join(root, relative));
    if (candidate !== path.join(root, relative) || !candidate.startsWith(`${root}${path.sep}`)) {
      throw new Error(`Theme-isolation content must not escape through a symlink: ${relative}`);
    }
  }
}

if (suiteMode === 'themes') validateThemeIsolation();

const tokenizationTest = path.join(vscodeTextmateDir, 'src/tests/tokenization.test.ts');
const themesTest = path.join(vscodeTextmateDir, 'src/tests/themes.test.ts');
const requiredTest = suiteMode === 'themes' ? themesTest : tokenizationTest;
if (!fs.existsSync(requiredTest)) {
  throw new Error(
    `vscode-textmate was not found at ${vscodeTextmateDir}. ` +
      'Clone the pinned checkout there or set VSCODE_TEXTMATE_DIR.',
  );
}

// Deliberately static: upstream discovery must not silently expand the safe
// gate. These are the 72 golden fixture test names in the pinned checkout.
const themeGoldenTestNames = [
  '12750.html',
  '13448.html',
  '14119.less',
  'COMMIT_EDITMSG',
  'Dockerfile',
  'basic.java',
  'git-rebase-todo',
  'issue-1550.yaml',
  'issue-4008.yaml',
  'issue-6303.yaml',
  'makefile',
  'test-13777.go',
  'test-4287.jade',
  'test-6611.rs',
  'test-7115.xml',
  'test-brackets.tsx',
  'test-cssvariables.less',
  'test-cssvariables.scss',
  'test-function-inv.ts',
  'test-issue11.ts',
  'test-issue5431.ts',
  'test-issue5465.ts',
  'test-issue5566.ts',
  'test-keywords.ts',
  'test-members.ts',
  'test-object-literals.ts',
  'test-regex.coffee',
  'test-strings.ts',
  'test-this.ts',
  'test-variables.css',
  'test.bat',
  'test.c',
  'test.cc',
  'test.clj',
  'test.coffee',
  'test.cpp',
  'test.cshtml',
  'test.css',
  'test.diff',
  'test.fs',
  'test.go',
  'test.groovy',
  'test.handlebars',
  'test.hbs',
  'test.html',
  'test.ini',
  'test.js',
  'test.json',
  'test.jsx',
  'test.less',
  'test.lua',
  'test.m',
  'test.md',
  'test.php',
  'test.pl',
  'test.ps1',
  'test.py',
  'test.r',
  'test.rb',
  'test.rs',
  'test.scss',
  'test.sh',
  'test.shader',
  'test.sql',
  'test.swift',
  'test.ts',
  'test.vb',
  'test.xml',
  'test.yaml',
  'test2.pl',
  'test6916.js',
  'tsconfig.json',
] as const;

function escapeRegExp(value: string): string {
  return value.replace(/[.*+?^${}()|[\]\\]/g, '\\$&');
}

const themeGoldenPattern = new RegExp(
  `^(?:${themeGoldenTestNames.map(escapeRegExp).join('|')})$`,
);

function conformanceRedirects(): Plugin {
  const testsDir = `${path.join(vscodeTextmateDir, 'src/tests')}${path.sep}`;
  return {
    name: 'textmate-go-conformance-redirects',
    enforce: 'pre',
    resolveId(source, importer) {
      if (!importer || !path.resolve(importer).startsWith(testsDir)) return null;
      if (source === '../main') return path.join(here, 'src/adapter.ts');
      if (source === './onigLibs') return path.join(here, 'src/onigLibs.ts');
      if (source === '../diffStateStacks') return path.join(here, 'src/stateDiff.ts');
      return null;
    },
  };
}

const upstreamTests =
  suiteMode === 'themes'
    ? [themesTest]
    : suiteMode === 'tokenization'
      ? [tokenizationTest, path.join(here, 'src/**/*.test.ts')]
      : [tokenizationTest, themesTest, path.join(here, 'src/**/*.test.ts')];

export default defineConfig({
  plugins: [conformanceRedirects()],
  resolve: {
    alias: {
      'vscode-textmate-reference': path.join(vscodeTextmateDir, 'src/main.ts'),
      'vscode-textmate-theme': path.join(vscodeTextmateDir, 'src/theme.ts'),
      'vscode-textmate-attributes': path.join(vscodeTextmateDir, 'src/encodedTokenAttributes.ts'),
      'vscode-textmate-basic-scopes': path.join(
        vscodeTextmateDir,
        'src/grammar/basicScopesAttributeProvider.ts',
      ),
      'vscode-textmate-onig-types': path.join(vscodeTextmateDir, 'src/onigLib.ts'),
    },
  },
  test: {
    include: upstreamTests,
    setupFiles: [path.join(here, 'src/setup.ts')],
    testNamePattern: suiteMode === 'themes' && !includeAllThemeTests ? themeGoldenPattern : undefined,
    testTimeout: 30_000,
    hookTimeout: 30_000,
    pool: 'forks',
    maxWorkers: backend === 'go' || suiteMode === 'themes' ? 1 : undefined,
  },
});
