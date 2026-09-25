export interface OpaqueStackDiff<T = unknown> {
  readonly next: T;
}

export function diffStateStacksRefEq<T>(_first: T, second: T): OpaqueStackDiff<T> {
  return { next: second };
}

export function applyStateStackDiff<T>(_stack: T | null, diff: OpaqueStackDiff<T>): T {
  return diff.next;
}

export type StackDiff = OpaqueStackDiff;
