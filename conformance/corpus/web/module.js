/* A multiline comment keeps lexical state
 * until this closing delimiter. */
export function greeting(name) {
  const prefix = "hello";
  return `${prefix}, ${name ?? "reader"}!`;
}
