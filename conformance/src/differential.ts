import fs from 'node:fs';
import path from 'node:path';
import { execFileSync } from 'node:child_process';
import { pathToFileURL } from 'node:url';
import { GoRegistry, ReferenceRegistry, type IRawGrammar } from './adapter';
import { getOniguruma } from './onigLibs';

export interface CorpusEntry {
  file: string;
  grammar: string;
  features: string[];
  source?: string;
}

interface CorpusSource {
  repository: string;
  revision: string;
  license: string;
}

export interface ScopeToken {
  text: string;
  scopes: string[];
  start: number;
  end: number;
}

export interface Difference {
  line: number;
  column: number;
  token: number;
  reference?: ScopeToken;
  go?: ScopeToken;
}

interface GrammarMetadata {
  name: string;
  scopeName: string;
  injectTo?: string[];
}

interface TokenizeResult {
  tokens: Array<{ startIndex: number; endIndex: number; scopes: string[] }>;
  ruleStack: unknown;
}

interface ScopeGrammar {
  tokenizeLine(line: string, state: unknown | null): TokenizeResult;
}

export interface DifferentialSession {
  compare(contents: string): Difference | undefined;
  close(): void;
}

export interface GeneratedCase {
  iteration: number;
  operation: 'splice' | 'edit' | 'utf8';
  grammar: string;
  sourceFile: string;
  contents: string;
}

function readJSON<T>(filename: string): T {
  return JSON.parse(fs.readFileSync(filename, 'utf8')) as T;
}

export function loadCorpus(
  corpusDir: string,
  externalRoots: Record<string, string> = {},
): Array<CorpusEntry & { contents: string }> {
  const manifest = readJSON<{
    sources?: Record<string, CorpusSource>;
    files: CorpusEntry[];
  }>(path.join(corpusDir, 'manifest.json'));
  for (const [name, source] of Object.entries(manifest.sources ?? {})) {
    const root = externalRoots[name];
    if (!root) throw new Error(`corpus source ${JSON.stringify(name)} has no configured checkout`);
    const revision = execFileSync('git', ['-C', root, 'rev-parse', 'HEAD'], { encoding: 'utf8' }).trim();
    if (revision !== source.revision) {
      throw new Error(`corpus source ${name} is at ${revision}, want ${source.revision}`);
    }
    const status = execFileSync('git', ['-C', root, 'status', '--porcelain'], { encoding: 'utf8' });
    if (status !== '') throw new Error(`corpus source ${name} checkout is dirty`);
  }
  return manifest.files.map((entry) => ({
    ...entry,
    contents: fs.readFileSync(
      entry.source
        ? path.join(externalRoots[entry.source]!, entry.file)
        : path.join(corpusDir, entry.file),
      'utf8',
    ),
  }));
}

export class GrammarRepository {
  private constructor(
    private readonly root: string,
    private readonly metadata: GrammarMetadata[],
    private readonly byScope: Map<string, GrammarMetadata>,
    private readonly byName: Map<string, GrammarMetadata>,
  ) {}

  static async open(root: string): Promise<GrammarRepository> {
    const modulePath = path.join(root, 'index.js');
    if (!fs.existsSync(modulePath)) {
      throw new Error(`tm-grammars was not found at ${root}`);
    }
    const loaded = (await import(pathToFileURL(modulePath).href)) as {
      grammars: GrammarMetadata[];
    };
    const metadata = loaded.grammars;
    return new GrammarRepository(
      root,
      metadata,
      new Map(metadata.map((grammar) => [grammar.scopeName, grammar])),
      new Map(metadata.map((grammar) => [grammar.name, grammar])),
    );
  }

  scopeForName(name: string): string {
    const metadata = this.byName.get(name);
    if (!metadata) throw new Error(`unknown tm-grammars grammar ${JSON.stringify(name)}`);
    return metadata.scopeName;
  }

  loadGrammar = async (scopeName: string): Promise<IRawGrammar | null> => {
    const metadata = this.byScope.get(scopeName);
    if (!metadata) return null;
    return readJSON<IRawGrammar>(path.join(this.root, 'grammars', `${metadata.name}.json`));
  };

  getInjections = (scopeName: string): string[] =>
    this.metadata
      .filter((grammar) => grammar.injectTo?.includes(scopeName))
      .map((grammar) => grammar.scopeName);

  async createSession(grammarName: string): Promise<DifferentialSession> {
    const scopeName = this.scopeForName(grammarName);
    const options = {
      onigLib: getOniguruma(),
      loadGrammar: this.loadGrammar,
      getInjections: this.getInjections,
    };
    const referenceRegistry = new ReferenceRegistry(options);
    const goRegistry = new GoRegistry(options);
    try {
      const [reference, go] = await Promise.all([
        referenceRegistry.loadGrammar(scopeName),
        goRegistry.loadGrammar(scopeName),
      ]);
      if (!reference || !go) throw new Error(`could not load ${scopeName}`);
      return {
        compare(contents: string): Difference | undefined {
          return firstDifference(
            tokenize(contents, reference as unknown as ScopeGrammar),
            tokenize(contents, go as unknown as ScopeGrammar),
          );
        },
        close(): void {
          referenceRegistry.dispose();
          goRegistry.dispose();
        },
      };
    } catch (error) {
      referenceRegistry.dispose();
      goRegistry.dispose();
      throw error;
    }
  }
}

function tokenize(contents: string, grammar: ScopeGrammar): ScopeToken[][] {
  let state: unknown | null = null;
  return contents.split(/\r\n|\r|\n/).map((line) => {
    const result = grammar.tokenizeLine(line, state);
    state = result.ruleStack;
    return result.tokens.map((token) => ({
      text: line.slice(token.startIndex, token.endIndex),
      scopes: token.scopes,
      start: token.startIndex,
      end: token.endIndex,
    }));
  });
}

function equalToken(left: ScopeToken | undefined, right: ScopeToken | undefined): boolean {
  if (!left || !right) return left === right;
  return left.text === right.text && arraysEqual(left.scopes, right.scopes);
}

function arraysEqual(left: string[], right: string[]): boolean {
  return left.length === right.length && left.every((value, index) => value === right[index]);
}

export function firstDifference(reference: ScopeToken[][], go: ScopeToken[][]): Difference | undefined {
  const lineCount = Math.max(reference.length, go.length);
  for (let lineIndex = 0; lineIndex < lineCount; lineIndex++) {
    const referenceTokens = reference[lineIndex] ?? [];
    const goTokens = go[lineIndex] ?? [];
    const tokenCount = Math.max(referenceTokens.length, goTokens.length);
    for (let tokenIndex = 0; tokenIndex < tokenCount; tokenIndex++) {
      const referenceToken = referenceTokens[tokenIndex];
      const goToken = goTokens[tokenIndex];
      if (!equalToken(referenceToken, goToken)) {
        return {
          line: lineIndex + 1,
          column: (referenceToken ?? goToken)?.start ?? 0,
          token: tokenIndex + 1,
          reference: referenceToken,
          go: goToken,
        };
      }
    }
  }
  return undefined;
}

export function formatDifference(difference: Difference): string {
  return (
    `line ${difference.line}, column ${difference.column}, token ${difference.token}: ` +
    `reference=${JSON.stringify(difference.reference ?? null)} ` +
    `go=${JSON.stringify(difference.go ?? null)}`
  );
}

class Random {
  constructor(private state: number) {}

  next(): number {
    let value = (this.state += 0x6d2b79f5);
    value = Math.imul(value ^ (value >>> 15), value | 1);
    value ^= value + Math.imul(value ^ (value >>> 7), value | 61);
    return ((value ^ (value >>> 14)) >>> 0) / 0x1_0000_0000;
  }

  int(limit: number): number {
    return Math.floor(this.next() * limit);
  }

  pick<T>(values: readonly T[]): T {
    return values[this.int(values.length)]!;
  }
}

export function generateCases(
  corpus: Array<CorpusEntry & { contents: string }>,
  seed: number,
  count: number,
): GeneratedCase[] {
  const random = new Random(seed >>> 0);
  const result: GeneratedCase[] = [];
  for (let iteration = 0; iteration < count; iteration++) {
    const source = random.pick(corpus);
    const operation = random.pick(['splice', 'edit', 'utf8'] as const);
    const contents =
      operation === 'splice'
        ? spliceLines(source.contents, random.pick(corpus).contents, random)
        : operation === 'edit'
          ? editText(source.contents, random)
          : randomUTF8(random);
    result.push({
      iteration,
      operation,
      grammar: source.grammar,
      sourceFile: source.file,
      contents,
    });
  }
  return result;
}

function spliceLines(contents: string, donorContents: string, random: Random): string {
  const lines = contents.split(/\r\n|\r|\n/);
  const donor = donorContents.split(/\r\n|\r|\n/);
  if (donor.length < 2) return contents;
  const first = random.int(donor.length);
  const second = random.int(donor.length);
  const start = Math.min(first, second);
  const end = Math.max(first, second) + 1;
  const insertAt = random.int(lines.length + 1);
  const chunk = donor.slice(start, end);
  return [...lines.slice(0, insertAt), ...chunk, ...lines.slice(insertAt)].join('\n');
}

function editText(contents: string, random: Random): string {
  const points = Array.from(contents);
  const edits = 1 + random.int(4);
  const alphabet = ['a', 'Z', '0', ' ', '\n', '/', '*', '`', '"', "'", '{', '}', '<', '>', 'λ', '🙂'];
  for (let index = 0; index < edits; index++) {
    const at = random.int(points.length + 1);
    const operation = random.int(3);
    if (operation === 0 && points.length > 0) points.splice(Math.min(at, points.length - 1), 1);
    else if (operation === 1 && points.length > 0) points[Math.min(at, points.length - 1)] = random.pick(alphabet);
    else points.splice(at, 0, random.pick(alphabet));
  }
  return points.join('');
}

function randomUTF8(random: Random): string {
  const alphabet = [
    ...Array.from('abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789 \t\n/*`\'"{}[]()<>=:;,.#-$'),
    'é',
    'λ',
    '中',
    '🙂',
    '\u0301',
  ];
  const length = 1 + random.int(256);
  return Array.from({ length }, () => random.pick(alphabet)).join('');
}

export function minimizeFailure(
  contents: string,
  stillFails: (candidate: string) => boolean,
  budget = 128,
): { contents: string; attempts: number } {
  let attempts = 0;
  const reduce = (units: string[], joiner: string): string[] => {
    let current = units;
    let granularity = 2;
    while (current.length > 1 && attempts < budget) {
      const chunkSize = Math.ceil(current.length / granularity);
      let reduced = false;
      for (let start = 0; start < current.length && attempts < budget; start += chunkSize) {
        const candidate = [...current.slice(0, start), ...current.slice(start + chunkSize)];
        attempts++;
        if (stillFails(candidate.join(joiner))) {
          current = candidate;
          granularity = Math.max(2, granularity - 1);
          reduced = true;
          break;
        }
      }
      if (!reduced) {
        if (granularity >= current.length) break;
        granularity = Math.min(current.length, granularity * 2);
      }
    }
    return current;
  };

  const lines = reduce(contents.split('\n'), '\n');
  const points = reduce(Array.from(lines.join('\n')), '');
  return { contents: points.join(''), attempts };
}
