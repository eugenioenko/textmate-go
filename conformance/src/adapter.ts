import * as Reference from 'vscode-textmate-reference';
import { EncodedTokenAttributes, OptionalStandardTokenType } from 'vscode-textmate-attributes';
import { BasicScopeAttributesProvider } from 'vscode-textmate-basic-scopes';
import { ScopeStack, Theme } from 'vscode-textmate-theme';
import { GoProcessClient } from './client';
import type { Token, TokenizeLineResponse } from './protocol';

const sharedGoClient = new GoProcessClient();

export async function closeSharedGoClient(): Promise<void> {
  await sharedGoClient.close();
}

export const parseRawGrammar = Reference.parseRawGrammar;
export const INITIAL = Reference.INITIAL;
export type IRawGrammar = Reference.IRawGrammar;
export type IRawTheme = Reference.IRawTheme;
export type IEmbeddedLanguagesMap = Reference.IEmbeddedLanguagesMap;
export type RegistryOptions = Reference.RegistryOptions;
export type StateStack = Reference.StateStack;
export type IGrammar = Reference.IGrammar;

interface OpaqueState {
  readonly handle: string | null;
}

interface LoadSettings {
  initialLanguage: number;
  embeddedLanguages: IEmbeddedLanguagesMap;
}

class MetadataResolver {
  private readonly basicScopes: InstanceType<typeof BasicScopeAttributesProvider>;

  constructor(
    private readonly settings: LoadSettings,
    private readonly getTheme: () => InstanceType<typeof Theme>,
  ) {
    this.basicScopes = new BasicScopeAttributesProvider(
      settings.initialLanguage,
      settings.embeddedLanguages,
    );
  }

  pack(tokens: Array<{ startIndex: number; scopes: string[] }>): Uint32Array {
    const packed: number[] = [];
    let previousMetadata = -1;
    for (const token of tokens) {
      const metadata = this.metadataFor(token.scopes);
      if (metadata !== previousMetadata) {
        packed.push(token.startIndex, metadata);
        previousMetadata = metadata;
      }
    }
    if (packed.length === 0) packed.push(0, this.metadataFor([]));
    return Uint32Array.from(packed);
  }

  private metadataFor(scopes: string[]): number {
    const theme = this.getTheme();
    const defaults = theme.getDefaults();
    let metadata = EncodedTokenAttributes.set(
      0,
      this.settings.initialLanguage,
      OptionalStandardTokenType.NotSet,
      false,
      defaults.fontStyle,
      defaults.foregroundId,
      defaults.backgroundId,
    );
    let stack: InstanceType<typeof ScopeStack> | null = null;
    for (const scope of scopes) {
      stack = ScopeStack.push(stack, [scope]);
      const basic = this.basicScopes.getBasicScopeAttributes(scope);
      const style = theme.match(stack);
      metadata = EncodedTokenAttributes.set(
        metadata,
        basic.languageId,
        basic.tokenType,
        null,
        style?.fontStyle ?? -1,
        style?.foregroundId ?? 0,
        style?.backgroundId ?? 0,
      );
    }
    return metadata >>> 0;
  }
}

export class ReferenceRegistry {
  private readonly inner: InstanceType<typeof Reference.Registry>;
  private theme: InstanceType<typeof Theme>;

  constructor(options: RegistryOptions) {
    this.inner = new Reference.Registry(options);
    this.theme = Theme.createFromRawTheme(options.theme, options.colorMap);
  }

  setTheme(theme: IRawTheme, colorMap?: string[]): void {
    this.inner.setTheme(theme, colorMap);
    this.theme = Theme.createFromRawTheme(theme, colorMap);
  }

  getColorMap(): string[] {
    return this.theme.getColorMap();
  }

  async loadGrammar(scopeName: string): Promise<ReferenceGrammar | null> {
    const grammar = await this.inner.loadGrammar(scopeName);
    return grammar
      ? new ReferenceGrammar(grammar, { initialLanguage: 0, embeddedLanguages: {} }, () => this.theme)
      : null;
  }

  async loadGrammarWithEmbeddedLanguages(
    scopeName: string,
    initialLanguage: number,
    embeddedLanguages: IEmbeddedLanguagesMap,
  ): Promise<ReferenceGrammar | null> {
    const grammar = await this.inner.loadGrammarWithEmbeddedLanguages(
      scopeName,
      initialLanguage,
      embeddedLanguages,
    );
    return grammar
      ? new ReferenceGrammar(grammar, { initialLanguage, embeddedLanguages }, () => this.theme)
      : null;
  }

  async loadGrammarWithConfiguration(
    scopeName: string,
    initialLanguage: number,
    configuration: { embeddedLanguages?: IEmbeddedLanguagesMap },
  ): Promise<ReferenceGrammar | null> {
    const grammar = await this.inner.loadGrammarWithConfiguration(
      scopeName,
      initialLanguage,
      configuration,
    );
    return grammar
      ? new ReferenceGrammar(
          grammar,
          { initialLanguage, embeddedLanguages: configuration.embeddedLanguages ?? {} },
          () => this.theme,
        )
      : null;
  }

  async addGrammar(
    rawGrammar: IRawGrammar,
    injections: string[] = [],
    initialLanguage = 0,
    embeddedLanguages: IEmbeddedLanguagesMap = {},
  ): Promise<ReferenceGrammar> {
    const grammar = await this.inner.addGrammar(
      rawGrammar,
      injections,
      initialLanguage,
      embeddedLanguages,
    );
    return new ReferenceGrammar(grammar, { initialLanguage, embeddedLanguages }, () => this.theme);
  }

  dispose(): void {
    this.inner.dispose();
  }
}

class ReferenceGrammar {
  private readonly metadata: MetadataResolver;

  constructor(
    private readonly inner: InstanceType<typeof Reference.Registry> extends never
      ? never
      : Reference.IGrammar,
    settings: LoadSettings,
    getTheme: () => InstanceType<typeof Theme>,
  ) {
    this.metadata = new MetadataResolver(settings, getTheme);
  }

  tokenizeLine(line: string, state: Reference.StateStack | null, timeLimit?: number): Reference.ITokenizeLineResult {
    return this.inner.tokenizeLine(line, state, timeLimit);
  }

  tokenizeLine2(line: string, state: Reference.StateStack | null, timeLimit?: number): Reference.ITokenizeLineResult2 {
    const result = this.inner.tokenizeLine(line, state, timeLimit);
    return { ...result, tokens: this.metadata.pack(result.tokens) };
  }
}

export class GoRegistry {
  private readonly client = sharedGoClient;
  private readonly grammars = new Map<string, IRawGrammar>();
  private registryHandle: string | null = null;
  private theme: InstanceType<typeof Theme>;

  constructor(private readonly options: RegistryOptions) {
    this.theme = Theme.createFromRawTheme(options.theme, options.colorMap);
  }

  setTheme(theme: IRawTheme, colorMap?: string[]): void {
    this.theme = Theme.createFromRawTheme(theme, colorMap);
  }

  getColorMap(): string[] {
    return this.theme.getColorMap();
  }

  async loadGrammar(scopeName: string): Promise<GoGrammar | null> {
    return this.load(scopeName, { initialLanguage: 0, embeddedLanguages: {} });
  }

  async loadGrammarWithEmbeddedLanguages(
    scopeName: string,
    initialLanguage: number,
    embeddedLanguages: IEmbeddedLanguagesMap,
  ): Promise<GoGrammar | null> {
    return this.load(scopeName, { initialLanguage, embeddedLanguages });
  }

  async loadGrammarWithConfiguration(
    scopeName: string,
    initialLanguage: number,
    configuration: { embeddedLanguages?: IEmbeddedLanguagesMap },
  ): Promise<GoGrammar | null> {
    return this.load(scopeName, {
      initialLanguage,
      embeddedLanguages: configuration.embeddedLanguages ?? {},
    });
  }

  async addGrammar(
    rawGrammar: IRawGrammar,
    _injections: string[] = [],
    initialLanguage = 0,
    embeddedLanguages: IEmbeddedLanguagesMap = {},
  ): Promise<GoGrammar> {
    this.grammars.set(rawGrammar.scopeName, rawGrammar);
    await this.recreateRegistry();
    return (await this.load(rawGrammar.scopeName, { initialLanguage, embeddedLanguages }))!;
  }

  dispose(): void {
    if (this.registryHandle) {
      void this.client.request('dispose', { registry: this.registryHandle }).catch(() => undefined);
    }
    this.registryHandle = null;
  }

  private async load(scopeName: string, settings: LoadSettings): Promise<GoGrammar | null> {
    const raw = await this.options.loadGrammar(scopeName);
    if (!raw) return null;
    const before = this.grammars.size;
    await this.collect(raw);
    if (!this.registryHandle || before !== this.grammars.size) await this.recreateRegistry();
    const response = await this.client.request<{ grammar: string }>('loadGrammar', {
      registry: this.registryHandle,
      scopeName,
      initialLanguage: settings.initialLanguage,
      embeddedLanguages: settings.embeddedLanguages,
    });
    return new GoGrammar(this.client, response.grammar, settings, () => this.theme);
  }

  private async collect(raw: IRawGrammar): Promise<void> {
    if (this.grammars.has(raw.scopeName)) return;
    this.grammars.set(raw.scopeName, raw);
    const dependencies = externalIncludes(raw);
    for (const injected of this.options.getInjections?.(raw.scopeName) ?? []) {
      dependencies.add(injected);
    }
    for (const scopeName of dependencies) {
      const dependency = await this.options.loadGrammar(scopeName);
      if (dependency) await this.collect(dependency);
    }
  }

  private async recreateRegistry(): Promise<void> {
    if (this.registryHandle) {
      await this.client.request('dispose', { registry: this.registryHandle });
    }
    const injections: Record<string, string[]> = {};
    for (const scopeName of this.grammars.keys()) {
      const values = this.options.getInjections?.(scopeName);
      if (values?.length) injections[scopeName] = values;
    }
    const response = await this.client.request<{ registry: string }>('newRegistry', {
      grammars: [...this.grammars.values()],
      injections,
    });
    this.registryHandle = response.registry;
  }
}

export class GoGrammar {
  private readonly metadata: MetadataResolver;

  constructor(
    private readonly client: GoProcessClient,
    private readonly grammarHandle: string,
    private readonly settings: LoadSettings,
    private readonly getTheme: () => InstanceType<typeof Theme>,
  ) {
    this.metadata = new MetadataResolver(settings, getTheme);
  }

  tokenizeLine(line: string, prevState: OpaqueState | null): {
    tokens: Token[];
    fonts: unknown[];
    ruleStack: OpaqueState;
    stoppedEarly: boolean;
  } {
    const response = this.client.requestSync<TokenizeLineResponse>('tokenizeLine', {
          grammar: this.grammarHandle,
          line,
          state: prevState?.handle ?? null,
        });
    return {
          tokens: response.tokens.map((token) => ({
            startIndex: token.start,
            endIndex: token.end,
            scopes: token.scopes,
          })) as unknown as Token[],
          fonts: [],
          ruleStack: { handle: response.state },
          stoppedEarly: response.stopped ?? false,
        };
  }

  tokenizeLine2(line: string, prevState: OpaqueState | null): {
    tokens: Uint32Array;
    fonts: unknown[];
    ruleStack: OpaqueState;
    stoppedEarly: boolean;
  } {
    const result = this.tokenizeLine(line, prevState);
    return {
      ...result,
      tokens: this.metadata.pack(
        result.tokens as unknown as Array<{ startIndex: number; scopes: string[] }>,
      ),
    };
  }
}

function externalIncludes(value: unknown, found = new Set<string>()): Set<string> {
  if (Array.isArray(value)) {
    for (const child of value) externalIncludes(child, found);
  } else if (value && typeof value === 'object') {
    for (const [key, child] of Object.entries(value)) {
      if (key === 'include' && typeof child === 'string') {
        const scope = child.split('#', 1)[0];
        if (scope && !scope.startsWith('$') && !scope.startsWith('#')) found.add(scope);
      } else {
        externalIncludes(child, found);
      }
    }
  }
  return found;
}

const backend = process.env.TEXTMATE_BACKEND ?? 'reference';
if (backend !== 'reference' && backend !== 'go') {
  throw new Error(`Invalid TEXTMATE_BACKEND ${JSON.stringify(backend)}; expected "reference" or "go"`);
}

export const Registry: typeof Reference.Registry =
  backend === 'reference'
    ? (ReferenceRegistry as unknown as typeof Reference.Registry)
    : (GoRegistry as unknown as typeof Reference.Registry);
