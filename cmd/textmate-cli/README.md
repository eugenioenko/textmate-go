# textmate-cli

`textmate-cli` is a JSON-lines server used by the Node conformance harness. It
reads one request per line from standard input and writes exactly one JSON
response per request to standard output. Diagnostics go to standard error.

The server parses and preloads the supplied raw grammars into a real
`textmate.Registry`, then tokenizes through the compiled `textmate.Grammar`.
Library token offsets are rune offsets; responses convert them to UTF-16 code
units at this JavaScript-facing protocol boundary.

Supported operations are `newRegistry`, `loadGrammar`, `tokenizeLine`, and
`dispose`, with the request fields described in `plan.md` section 9.1. Registry,
grammar, and state values are opaque handles. Every tokenization returns a new
immutable state handle, so callers may retain and branch from earlier states. A
state can only be used with the grammar handle that created it; `dispose`
invalidates the registry and all of its grammar and state handles.

Request failures have a human-readable `error` string and a machine-readable
`code`. Non-fatal regex translation, compilation, and match failures appear in
a tokenization response's `diagnostics` array; a bounded summary is also
written to standard error. Each diagnostic is reported once per grammar
handle. The request `id` is echoed in every response.
