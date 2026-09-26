#!/usr/bin/env node
// Compares textmate-go, Chroma, and vscode-textmate on benchmarks/corpus.json.
// Every engine tokenizes each case line by line after five warm-up passes;
// the reported figure is the fastest of COUNT runs, in microseconds per line.

import { execFileSync } from 'node:child_process';
import fs from 'node:fs';
import os from 'node:os';
import path from 'node:path';
import zlib from 'node:zlib';
import { createRequire } from 'node:module';
import { fileURLToPath } from 'node:url';

const root = path.resolve(path.dirname(fileURLToPath(import.meta.url)), '..');
const count = Number(process.env.COUNT ?? 5);
const benchTime = process.env.BENCH_TIME ?? '1s';
const filter = new RegExp(process.env.FILTER ?? '.');
const skipJS = process.env.SKIP_JS === '1';
const vscodeTextmateDir = path.resolve(
  process.env.VSCODE_TEXTMATE_DIR ?? path.join(root, '../vscode-textmate'),
);

const manifest = JSON.parse(fs.readFileSync(path.join(root, 'benchmarks/corpus.json'), 'utf8'));
const cases = manifest.cases.filter((c) => filter.test(c.name));
if (cases.length === 0) throw new Error(`FILTER=${filter} matches no case`);

function repeatedLines({ fixture, lines }) {
  const seed = fs.readFileSync(path.join(root, fixture), 'utf8').replace(/\n$/, '').split('\n');
  return Array.from({ length: lines }, (_, i) => seed[i % seed.length]);
}

// Build first, then run the binary, so compilation does not share the timed
// window with the measurements.
function runGo() {
  const dir = path.join(root, 'benchmarks/chroma');
  const binary = path.join(fs.mkdtempSync(path.join(os.tmpdir(), 'textmate-bench-')), 'bench.test');
  console.error('building Go benchmark...');
  execFileSync('go', ['test', '-c', '-o', binary, '.'], {
    cwd: dir,
    env: { ...process.env, CGO_ENABLED: '0' },
    stdio: ['ignore', 'inherit', 'inherit'],
  });
  const names = cases.map((c) => c.name).join('|');
  console.error('running textmate-go and Chroma...');
  const output = execFileSync(
    binary,
    [
      '-test.run', '^$',
      '-test.bench', `^BenchmarkWarmLineByLine$/^(${names})$`,
      '-test.benchmem',
      '-test.benchtime', benchTime,
      '-test.count', String(count),
    ],
    { cwd: dir, encoding: 'utf8', maxBuffer: 64 << 20 },
  );
  fs.rmSync(path.dirname(binary), { recursive: true, force: true });

  const results = {};
  for (const line of output.split('\n')) {
    const match = line.match(/^BenchmarkWarmLineByLine\/([^/]+)\/(textmate-go|chroma)\S*\s/);
    if (!match) continue;
    const [, name, engine] = match;
    const nsPerLine = Number(line.match(/([\d.]+) ns\/line/)?.[1]);
    const allocsPerOp = Number(line.match(/([\d.]+) allocs\/op/)?.[1]);
    const key = `${name}/${engine}`;
    if (!results[key] || nsPerLine < results[key].nsPerLine) {
      results[key] = { nsPerLine, allocsPerOp };
    }
  }
  return results;
}

async function runVscodeTextmate() {
  const mainJS = path.join(vscodeTextmateDir, 'out/src/main.js');
  if (!fs.existsSync(mainJS)) {
    throw new Error(
      `${mainJS} not found. Clone microsoft/vscode-textmate next to this repository ` +
        '(or set VSCODE_TEXTMATE_DIR) and run `npm ci && npm run compile` in it, or set SKIP_JS=1.',
    );
  }
  const require = createRequire(import.meta.url);
  const vsctm = require(mainJS);
  const onigRequire = createRequire(path.join(root, 'conformance/package.json'));
  const oniguruma = onigRequire('vscode-oniguruma');
  const wasm = fs.readFileSync(
    path.join(path.dirname(onigRequire.resolve('vscode-oniguruma/package.json')), 'release/onig.wasm'),
  );
  await oniguruma.loadWASM(wasm.buffer);

  const grammarsDir = path.join(root, 'grammars/data');
  const byScope = new Map();
  for (const file of fs.readdirSync(grammarsDir)) {
    if (!file.endsWith('.json.gz')) continue;
    const raw = zlib.gunzipSync(fs.readFileSync(path.join(grammarsDir, file))).toString();
    byScope.set(JSON.parse(raw).scopeName, { raw, file: file.replace(/\.gz$/, '') });
  }
  const registry = new vsctm.Registry({
    onigLib: Promise.resolve({
      createOnigScanner: (sources) => new oniguruma.OnigScanner(sources),
      createOnigString: (text) => new oniguruma.OnigString(text),
    }),
    loadGrammar: async (scopeName) => {
      const entry = byScope.get(scopeName);
      return entry ? vsctm.parseRawGrammar(entry.raw, entry.file) : null;
    },
  });

  const budgetNs = BigInt(Math.round(parseDuration(benchTime) * 1e9));
  const results = {};
  for (const test of cases) {
    console.error(`running vscode-textmate ${test.name}...`);
    const lines = repeatedLines(test);
    const grammar = await registry.loadGrammar(test.scopeName);
    if (!grammar) throw new Error(`vscode-textmate could not load ${test.scopeName}`);
    const pass = () => {
      let state = vsctm.INITIAL;
      for (const line of lines) state = grammar.tokenizeLine(line, state).ruleStack;
    };
    for (let i = 0; i < 5; i++) pass();
    let best = Infinity;
    for (let run = 0; run < count; run++) {
      let iterations = 0;
      const started = process.hrtime.bigint();
      let elapsed = 0n;
      while (elapsed < budgetNs) {
        pass();
        iterations++;
        elapsed = process.hrtime.bigint() - started;
      }
      best = Math.min(best, Number(elapsed) / (iterations * lines.length));
    }
    results[test.name] = { nsPerLine: best };
  }
  return results;
}

function parseDuration(value) {
  const match = String(value).match(/^([\d.]+)(ms|s)$/);
  if (!match) throw new Error(`BENCH_TIME must look like 500ms or 2s, got ${value}`);
  return Number(match[1]) / (match[2] === 'ms' ? 1000 : 1);
}

const us = (ns) => (ns === undefined || Number.isNaN(ns) ? 'n/a' : (ns / 1000).toFixed(1));
const ratio = (a, b) => (a && b ? `${(a / b).toFixed(2)}x` : 'n/a');

const goResults = runGo();
const jsResults = skipJS ? {} : await runVscodeTextmate();

const header = ['Case', 'textmate-go', 'allocs/line', 'vscode-textmate', 'chroma', 'vs JS', 'vs chroma'];
console.log(`\nus/line, fastest of ${count} runs; ratios are textmate-go time over the other engine.\n`);
console.log(`| ${header.join(' | ')} |`);
console.log(`|${header.map((_, i) => (i === 0 ? ' --- ' : ' ---: ')).join('|')}|`);
for (const test of cases) {
  const tm = goResults[`${test.name}/textmate-go`];
  const chroma = goResults[`${test.name}/chroma`];
  const js = jsResults[test.name];
  console.log(
    `| ${[
      test.name,
      us(tm?.nsPerLine),
      tm ? (tm.allocsPerOp / test.lines).toFixed(1) : 'n/a',
      us(js?.nsPerLine),
      us(chroma?.nsPerLine),
      ratio(tm?.nsPerLine, js?.nsPerLine),
      ratio(tm?.nsPerLine, chroma?.nsPerLine),
    ].join(' | ')} |`,
  );
}
