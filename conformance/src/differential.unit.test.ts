import fs from 'node:fs';
import path from 'node:path';
import { fileURLToPath } from 'node:url';
import { describe, expect, test } from 'vitest';
import {
  firstDifference,
  generateCases,
  minimizeFailure,
  type CorpusEntry,
  type ScopeToken,
} from './differential';

const corpus: Array<CorpusEntry & { contents: string }> = [
  {
    file: 'sample.ts',
    grammar: 'typescript',
    features: ['test'],
    contents: 'const value = `hello ${name}`;\n',
  },
];
const repositoryRoot = path.resolve(path.dirname(fileURLToPath(import.meta.url)), '../..');
const embeddedSupportGrammars = new Set(['html-derivative']);

describe('differential infrastructure', () => {
  test('seeded generation is reproducible and remains valid UTF-8', () => {
    const left = generateCases(corpus, 12345, 24);
    const right = generateCases(corpus, 12345, 24);
    expect(left).toEqual(right);
    expect(left.some((entry) => entry.operation === 'splice')).toBe(true);
    expect(left.some((entry) => entry.operation === 'edit')).toBe(true);
    expect(left.some((entry) => entry.operation === 'utf8')).toBe(true);
    for (const entry of left) {
      expect(Buffer.from(entry.contents, 'utf8').toString('utf8')).toBe(entry.contents);
    }
  });

  test('reports the first differing full scope stack with its column', () => {
    const base: ScopeToken = { text: 'value', scopes: ['source.ts'], start: 4, end: 9 };
    const difference = firstDifference(
      [[base]],
      [[{ ...base, scopes: ['source.ts', 'variable.other.ts'] }]],
    );
    expect(difference).toMatchObject({ line: 1, column: 4, token: 1 });
  });

  test('minimizes a reproducible failure within its attempt budget', () => {
    const result = minimizeFailure('noise\nleft BUG right\nnoise', (value) => value.includes('BUG'), 64);
    expect(result.contents).toBe('BUG');
    expect(result.attempts).toBeLessThanOrEqual(64);
  });

  test('the checked-in corpus covers every embedded curated grammar', () => {
    const curated = fs
      .readFileSync(path.join(repositoryRoot, 'grammars/curated.txt'), 'utf8')
      .split(/\r?\n/)
      .map((line) => line.trim())
      .filter((line) => line !== '' && !line.startsWith('#'));
    const manifest = JSON.parse(
      fs.readFileSync(path.join(repositoryRoot, 'conformance/corpus/manifest.json'), 'utf8'),
    ) as { files: CorpusEntry[] };
    const roots = curated.filter((grammar) => !embeddedSupportGrammars.has(grammar));
    const covered = new Set(manifest.files.map((entry) => entry.grammar));
    expect(roots.filter((grammar) => !covered.has(grammar))).toEqual([]);

    const upstreamSamples = manifest.files.filter((entry) => entry.source === 'tm-grammars');
    expect(upstreamSamples).toHaveLength(roots.length);
    expect(new Set(upstreamSamples.map((entry) => entry.file)).size).toBe(roots.length);
    expect(upstreamSamples.map((entry) => entry.grammar).sort()).toEqual(roots.sort());
  });
});
