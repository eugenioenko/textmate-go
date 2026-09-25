import fs from 'node:fs';
import path from 'node:path';
import { fileURLToPath } from 'node:url';
import { expect, test } from 'vitest';
import {
  GoRegistry,
  ReferenceRegistry,
  type IEmbeddedLanguagesMap,
} from './adapter';
import { getOniguruma } from './onigLibs';
import {
  Resolver,
  type IGrammarRegistration,
  type ILanguageRegistration,
} from '../../../vscode-textmate/src/tests/resolver';

const enabled = process.env.TEXTMATE_THEME_SCOPE_DIFF === '1';
const here = path.dirname(fileURLToPath(import.meta.url));
const vscodeTextmateDir = path.resolve(
  process.env.VSCODE_TEXTMATE_DIR ?? path.join(here, '../../../vscode-textmate'),
);
const themesDir = path.join(vscodeTextmateDir, 'test-cases/themes');

interface ScopeToken {
  text: string;
  scopes: string[];
}

interface ScopeLine {
  tokens: ScopeToken[];
}

interface TokenizeResult {
  tokens: Array<{ startIndex: number; endIndex: number; scopes: string[] }>;
  ruleStack: unknown;
}

interface ScopeGrammar {
  tokenizeLine(line: string, state: unknown | null): TokenizeResult;
}

interface Difference {
  file: string;
  line?: number;
  token?: number;
  reference?: ScopeToken;
  go?: ScopeToken;
  error?: string;
}

function readJSON<T>(filename: string): T {
  return JSON.parse(fs.readFileSync(filename, 'utf8')) as T;
}

function createResolver(): Resolver {
  const grammars = readJSON<IGrammarRegistration[]>(path.join(themesDir, 'grammars.json'));
  for (const grammar of grammars) grammar.path = path.join(themesDir, grammar.path);
  const languages = readJSON<ILanguageRegistration[]>(path.join(themesDir, 'languages.json'));
  return new Resolver(grammars, languages, getOniguruma());
}

function sourceFiles(): string[] {
  return fs
    .readdirSync(path.join(themesDir, 'tests'))
    .filter((filename) => !/\.(?:result(?:\.patch)?|actual|diff\.html)$/.test(filename))
    .sort();
}

function languageSettings(
  resolver: Resolver,
  filename: string,
): { scopeName: string; languageID: number; embeddedLanguages: IEmbeddedLanguagesMap } {
  const language =
    resolver.findLanguageByExtension(path.extname(filename)) ??
    resolver.findLanguageByFilename(filename);
  if (!language) throw new Error(`could not determine language for ${filename}`);

  const registration = resolver.findGrammarByLanguage(language);
  const embeddedLanguages: IEmbeddedLanguagesMap = Object.create(null) as IEmbeddedLanguagesMap;
  for (const [scopeName, embeddedLanguage] of Object.entries(
    registration.embeddedLanguages ?? {},
  )) {
    embeddedLanguages[scopeName] = resolver.language2id[embeddedLanguage];
  }
  return {
    scopeName: registration.scopeName,
    languageID: resolver.language2id[language],
    embeddedLanguages,
  };
}

function tokenize(contents: string, grammar: ScopeGrammar): ScopeLine[] {
  let state: unknown | null = null;
  return contents.split(/\r\n|\r|\n/).map((line) => {
    const result = grammar.tokenizeLine(line, state);
    state = result.ruleStack;
    return {
      tokens: result.tokens.map((token) => ({
        text: line.slice(token.startIndex, token.endIndex),
        scopes: token.scopes,
      })),
    };
  });
}

function firstDifference(
  file: string,
  reference: ScopeLine[],
  go: ScopeLine[],
): Difference | undefined {
  const lineCount = Math.max(reference.length, go.length);
  for (let lineIndex = 0; lineIndex < lineCount; lineIndex++) {
    const referenceTokens = reference[lineIndex]?.tokens ?? [];
    const goTokens = go[lineIndex]?.tokens ?? [];
    const tokenCount = Math.max(referenceTokens.length, goTokens.length);
    for (let tokenIndex = 0; tokenIndex < tokenCount; tokenIndex++) {
      const referenceToken = referenceTokens[tokenIndex];
      const goToken = goTokens[tokenIndex];
      if (
        referenceToken?.text !== goToken?.text ||
        JSON.stringify(referenceToken?.scopes) !== JSON.stringify(goToken?.scopes)
      ) {
        return {
          file,
          line: lineIndex + 1,
          token: tokenIndex + 1,
          reference: referenceToken,
          go: goToken,
        };
      }
    }
  }
  return undefined;
}

function display(token: ScopeToken | undefined): string {
  return token ? JSON.stringify(token) : '<missing>';
}

test.skipIf(!enabled)(
  'matches reference scopes across theme fixtures',
  { timeout: 0 },
  async () => {
    const started = performance.now();
    const referenceResolver = createResolver();
    const goResolver = createResolver();
    const referenceRegistry = new ReferenceRegistry(referenceResolver);
    const goRegistry = new GoRegistry(goResolver);
    const files = sourceFiles();
    const differences: Difference[] = [];

    try {
      for (const file of files) {
        try {
          const settings = languageSettings(referenceResolver, file);
          const [referenceGrammar, goGrammar] = await Promise.all([
            referenceRegistry.loadGrammarWithEmbeddedLanguages(
              settings.scopeName,
              settings.languageID,
              settings.embeddedLanguages,
            ),
            goRegistry.loadGrammarWithEmbeddedLanguages(
              settings.scopeName,
              settings.languageID,
              settings.embeddedLanguages,
            ),
          ]);
          if (!referenceGrammar || !goGrammar) {
            throw new Error(`could not load ${settings.scopeName}`);
          }

          const contents = fs.readFileSync(path.join(themesDir, 'tests', file), 'utf8');
          const difference = firstDifference(
            file,
            tokenize(contents, referenceGrammar as unknown as ScopeGrammar),
            tokenize(contents, goGrammar as unknown as ScopeGrammar),
          );
          if (difference) differences.push(difference);
        } catch (error) {
          differences.push({
            file,
            error: error instanceof Error ? error.message : String(error),
          });
        }
      }
    } finally {
      referenceRegistry.dispose();
      goRegistry.dispose();
    }

    const parity = files.length - differences.length;
    const elapsed = ((performance.now() - started) / 1000).toFixed(2);
    console.log(
      `[theme-scope-diff] parity ${parity}/${files.length}; differences ${differences.length}; ${elapsed}s`,
    );
    for (const difference of differences) {
      if (difference.error) {
        console.log(`[theme-scope-diff] ${difference.file}: ERROR ${difference.error}`);
        continue;
      }
      console.log(
        `[theme-scope-diff] ${difference.file}:${difference.line}:${difference.token} ` +
        `reference=${display(difference.reference)} go=${display(difference.go)}`,
      );
    }
    expect(
      differences,
      differences
        .map((difference) =>
          difference.error
            ? `${difference.file}: ${difference.error}`
            : `${difference.file}:${difference.line}:${difference.token} ` +
              `reference=${display(difference.reference)} go=${display(difference.go)}`,
        )
        .join('\n'),
    ).toEqual([]);
  },
);
