import path from 'node:path';
import { fileURLToPath } from 'node:url';
import { describe, expect, test } from 'vitest';
import {
  formatDifference,
  generateCases,
  GrammarRepository,
  loadCorpus,
  minimizeFailure,
  type Difference,
} from './differential';

const enabled = process.env.TEXTMATE_DIFFERENTIAL === '1';
const mode = process.env.TEXTMATE_DIFFERENTIAL_MODE ?? 'corpus';
const here = path.dirname(fileURLToPath(import.meta.url));
const repositoryRoot = path.resolve(here, '../..');
const corpusDir = path.join(repositoryRoot, 'conformance/corpus');
const grammarsDir = path.resolve(
  process.env.TM_GRAMMARS_DIR ?? path.join(repositoryRoot, '../tm-grammars/packages/tm-grammars'),
);
const tmGrammarsCheckout = path.resolve(grammarsDir, '../..');

describe.skipIf(!enabled)('vscode-textmate differential', () => {
  test.skipIf(mode !== 'corpus' && mode !== 'all')(
    'matches the checked-in real-file corpus',
    { timeout: 180_000 },
    async () => {
      const started = performance.now();
      const repository = await GrammarRepository.open(grammarsDir);
      const corpus = loadCorpus(corpusDir, { 'tm-grammars': tmGrammarsCheckout });
      const failures: string[] = [];

      for (const entry of corpus) {
        const session = await repository.createSession(entry.grammar);
        try {
          const difference = session.compare(entry.contents);
          if (difference) failures.push(`${entry.file}: ${formatDifference(difference)}`);
        } finally {
          session.close();
        }
      }

      const seconds = ((performance.now() - started) / 1000).toFixed(2);
      console.log(`[differential] corpus ${corpus.length - failures.length}/${corpus.length}; ${seconds}s`);
      expect(failures, failures.join('\n')).toEqual([]);
    },
  );

  test.skipIf(mode !== 'fuzz' && mode !== 'all')(
    'matches seeded splice/edit/valid-UTF-8 fuzz cases',
    { timeout: 180_000 },
    async () => {
      const started = performance.now();
      const seed = parseInteger('TEXTMATE_DIFF_SEED', 0x5eed_c0de);
      const count = parseInteger('TEXTMATE_DIFF_FUZZ_CASES', 48);
      const repository = await GrammarRepository.open(grammarsDir);
      const corpus = loadCorpus(corpusDir, { 'tm-grammars': tmGrammarsCheckout });
      const generated = generateCases(corpus, seed, count);
      const byGrammar = new Map<string, typeof generated>();
      for (const entry of generated) {
        const group = byGrammar.get(entry.grammar) ?? [];
        group.push(entry);
        byGrammar.set(entry.grammar, group);
      }

      for (const [grammar, cases] of byGrammar) {
        const session = await repository.createSession(grammar);
        try {
          for (const generatedCase of cases) {
            const difference = session.compare(generatedCase.contents);
            if (!difference) continue;
            const minimized = minimizeFailure(
              generatedCase.contents,
              (candidate) => session.compare(candidate) !== undefined,
            );
            const minimizedDifference = session.compare(minimized.contents) as Difference;
            throw new Error(
              [
                `differential fuzz failure (seed=${seed}, iteration=${generatedCase.iteration}, ` +
                  `operation=${generatedCase.operation}, grammar=${grammar}, source=${generatedCase.sourceFile})`,
                formatDifference(minimizedDifference),
                `minimized after ${minimized.attempts} attempts: ${JSON.stringify(minimized.contents)}`,
                `replay: TEXTMATE_DIFF_SEED=${seed} TEXTMATE_DIFF_FUZZ_CASES=${generatedCase.iteration + 1} pnpm test:differential:fuzz`,
              ].join('\n'),
            );
          }
        } finally {
          session.close();
        }
      }

      const seconds = ((performance.now() - started) / 1000).toFixed(2);
      console.log(`[differential] fuzz ${count}/${count}; seed ${seed}; ${seconds}s`);
    },
  );
});

function parseInteger(name: string, fallback: number): number {
  const raw = process.env[name];
  if (raw === undefined) return fallback;
  if (!/^(?:0x[\da-f]+|\d+)$/i.test(raw)) throw new Error(`${name} must be a non-negative integer`);
  const value = Number.parseInt(raw, raw.toLowerCase().startsWith('0x') ? 16 : 10);
  if (!Number.isSafeInteger(value) || value < 0) throw new Error(`${name} is out of range`);
  return value;
}
