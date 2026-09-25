export interface ProtocolError {
  error: string;
  code?: string;
}

export interface NewRegistryRequest {
  grammars: unknown[];
  injections: Record<string, string[]>;
}

export interface Token {
  start: number;
  end: number;
  scopes: string[];
}

export interface TokenizeLineResponse {
  tokens: Token[];
  state: string;
  stopped?: boolean;
}
