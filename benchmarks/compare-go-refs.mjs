#!/usr/bin/env node
// Compares the current checkout with another textmate-go checkout using only
// the lightweight Go benchmark. Intended for stable, same-runner PR reporting.

import { execFileSync } from 'node:child_process';
import fs from 'node:fs';
import os from 'node:os';
import path from 'node:path';
import { fileURLToPath } from 'node:url';

const root = path.resolve(path.dirname(fileURLToPath(import.meta.url)), '..');
const benchmarkDir = path.join(root, 'benchmarks/textmate');
const baseDir = path.resolve(process.env.BASE_DIR ?? '');
const count = Number(process.env.COUNT ?? 3);
const benchTime = process.env.BENCH_TIME ?? '500ms';
const marker = '<!-- textmate-go-performance -->';

if (!process.env.BASE_DIR) throw new Error('BASE_DIR must name the base checkout');
if (!fs.existsSync(path.join(baseDir, 'go.mod'))) throw new Error(`${baseDir} is not a Go module`);
if (!Number.isInteger(count) || count < 1) throw new Error(`COUNT must be a positive integer, got ${count}`);

const manifest = JSON.parse(fs.readFileSync(path.join(root, 'benchmarks/corpus.json'), 'utf8'));
const cases = manifest.cases;
const names = cases.map((test) => test.name);
const escapedNames = names.map((name) => name.replace(/[.*+?^${}()|[\]\\]/g, '\\$&'));
const temporaryDir = fs.mkdtempSync(path.join(os.tmpdir(), 'textmate-go-pr-bench-'));
const headBinary = path.join(temporaryDir, 'head.test');
const baseBinary = path.join(temporaryDir, 'base.test');
const workspace = path.join(temporaryDir, 'base.work');

function compile(binary, environment) {
  execFileSync('go', ['test', '-c', '-o', binary, '.'], {
    cwd: benchmarkDir,
    env: { ...process.env, CGO_ENABLED: '0', ...environment },
    stdio: ['ignore', 'inherit', 'inherit'],
  });
}

function run(binary) {
  return execFileSync(
    binary,
    [
      '-test.run', '^$',
      '-test.bench', `^BenchmarkWarmLineByLine$/^(${escapedNames.join('|')})$`,
      '-test.benchmem',
      '-test.benchtime', benchTime,
      '-test.count', '1',
    ],
    { cwd: benchmarkDir, encoding: 'utf8', maxBuffer: 64 << 20 },
  );
}

function collect(output, results) {
  for (const line of output.split('\n')) {
    const name = names.find((candidate) => line.startsWith(`BenchmarkWarmLineByLine/${candidate}-`));
    if (!name) continue;
    const nsPerLine = Number(line.match(/([\d.]+) ns\/line/)?.[1]);
    const allocsPerOp = Number(line.match(/([\d.]+) allocs\/op/)?.[1]);
    if (!Number.isFinite(nsPerLine) || !Number.isFinite(allocsPerOp)) continue;
    if (!results[name] || nsPerLine < results[name].nsPerLine) {
      results[name] = { nsPerLine, allocsPerOp };
    }
  }
}

const short = (value, fallback) => value ? value.slice(0, 12) : fallback;
const percent = (next, previous) => {
  if (!Number.isFinite(next) || !Number.isFinite(previous) || previous === 0) return 'n/a';
  const change = ((next / previous) - 1) * 100;
  if (Math.abs(change) < 0.05) return '0.0%';
  return `${change >= 0 ? '+' : ''}${change.toFixed(1)}%`;
};

try {
  fs.writeFileSync(
    workspace,
    `go 1.25\n\nuse ${JSON.stringify(benchmarkDir)}\n\n` +
      `replace github.com/eugenioenko/textmate-go => ${JSON.stringify(baseDir)}\n`,
  );
  console.error('building PR and base benchmark binaries...');
  compile(headBinary, { GOWORK: 'off' });
  compile(baseBinary, { GOWORK: workspace });

  const head = {};
  const base = {};
  for (let runIndex = 0; runIndex < count; runIndex++) {
    console.error(`benchmark pass ${runIndex + 1}/${count}...`);
    const order = runIndex % 2 === 0
      ? [[baseBinary, base], [headBinary, head]]
      : [[headBinary, head], [baseBinary, base]];
    for (const [binary, results] of order) collect(run(binary), results);
  }

  for (const test of cases) {
    if (!base[test.name] || !head[test.name]) throw new Error(`missing benchmark result for ${test.name}`);
  }

  const baseLabel = short(process.env.BASE_SHA, 'base');
  const headLabel = short(process.env.HEAD_SHA, 'PR');
  const ratios = cases.map((test) => head[test.name].nsPerLine / base[test.name].nsPerLine);
  const geometricChange = (Math.exp(ratios.reduce((sum, ratio) => sum + Math.log(ratio), 0) / ratios.length) - 1) * 100;
  const lines = [
    marker,
    '## TextMate Go performance',
    '',
    `All ${cases.length} core + extended cases, comparing \`${baseLabel}\` with \`${headLabel}\` on the same runner. ` +
      `Values are the fastest of ${count} alternating passes at ${benchTime} per case.`,
    '',
    `Geometric-mean time change: **${geometricChange >= 0 ? '+' : ''}${geometricChange.toFixed(1)}%**.`,
    '',
    '| Case | Base us/line | PR us/line | Time change | Base allocs/line | PR allocs/line | Allocation change |',
    '| --- | ---: | ---: | ---: | ---: | ---: | ---: |',
  ];
  for (const test of cases) {
    const before = base[test.name];
    const after = head[test.name];
    const baseAllocs = before.allocsPerOp / test.lines;
    const headAllocs = after.allocsPerOp / test.lines;
    lines.push(
      `| ${test.name} | ${(before.nsPerLine / 1000).toFixed(1)} | ${(after.nsPerLine / 1000).toFixed(1)} | ` +
      `${percent(after.nsPerLine, before.nsPerLine)} | ${baseAllocs.toFixed(1)} | ${headAllocs.toFixed(1)} | ` +
      `${percent(headAllocs, baseAllocs)} |`,
    );
  }
  lines.push('', '_Timing changes are informational because shared CI runners are noisy; a benchmark error still fails the job._', '');
  const report = lines.join('\n');
  console.log(report);
  if (process.env.REPORT_FILE) fs.writeFileSync(process.env.REPORT_FILE, report);
  if (process.env.GITHUB_STEP_SUMMARY) fs.appendFileSync(process.env.GITHUB_STEP_SUMMARY, report);
} finally {
  fs.rmSync(temporaryDir, { recursive: true, force: true });
}
