# textmate-go

`textmate-go` is a pure-Go port of Microsoft's
[`vscode-textmate`](https://github.com/microsoft/vscode-textmate). It provides
line-oriented TextMate grammar tokenization with an immutable state stack that
can be carried between lines.

The port is under active development, so its API may still evolve before a
stable release. It targets Go 1.25 or newer and builds with `CGO_ENABLED=0`.

## Reference revisions

Development and conformance testing use sibling checkouts at these revisions:

- `vscode-textmate`: `fbe49961ab8077e587fdf5282019655ae69e5f9e`
- `textmate-grammars-themes`: `37edd1b26f18838050661d912334aba0ca7f4931`
- `gopher-textmate`: `ce2b42e5386c93ae781add9df2e6328338b06f9e`
- `regexp2`: v2.8.0 base `9d0d2ffe88a8b90012f7979ec85424e46d5ef48f`,
  local capture-index branch `49b524a3791dd2f86ec56b81bb7acf9626c1ded4`

The implementation is ported from `vscode-textmate`; its Microsoft MIT license
is reproduced in [`THIRD_PARTY_NOTICES`](THIRD_PARTY_NOTICES).

## Regex engine

TextMate grammars require lookbehind, backreferences, and position-sensitive
anchors that Go's standard RE2-based `regexp` package cannot provide. This port
uses `github.com/dlclark/regexp2/v2` v2.8.0 with RE2 compatibility mode disabled.
The v2 line was selected for ordered mixed captures, current Unicode tables,
bounded backtracking, and its Go 1.25 API. Oniguruma-only syntax is translated
before compilation; unsupported constructs are recorded as diagnostics and the
affected construct or pattern is safely degraded.

The current performance branch uses a sibling `../regexp2` checkout with a
text-free capture-index API. This avoids constructing public match/group data
that the tokenizer immediately discarded. The local replacement is temporary:
it must be upstreamed or moved to a durable fork before a module release.

## Development

```sh
make test
make bench
make lint
CGO_ENABLED=0 go build ./...
```

## Embedded grammars

The optional `grammars` subpackage embeds 40 license-reviewed common grammar
roots plus one MIT support grammar used for Markdown's inline HTML. Assets are
individually compressed and loaded once on demand:

```go
import (
	"log"
	"time"

	textmate "github.com/eugenioenko/textmate-go"
	"github.com/eugenioenko/textmate-go/grammars"
)

registry := textmate.NewRegistry(textmate.RegistryOptions{
	LoadGrammar: grammars.Load,
})
grammar, err := registry.LoadGrammar(grammars.ScopeForFilename("main.go"))
if err != nil {
	log.Fatal(err)
}
if grammar == nil {
	log.Fatal("no embedded grammar for main.go")
}
defer registry.Dispose()

var state *textmate.StateStack
for _, line := range []string{"package main", `const message = "hello"`} {
	result := grammar.TokenizeLine(line, state)
	for _, token := range result.Tokens {
		// token.Start and token.End are rune offsets; token.Scopes is ordered
		// from the root grammar to the most specific matched scope.
		// ScopeStack and its ID are comparable cache keys for resolved styles.
		styleCacheKey := token.ScopeStack.ID()
		_ = styleCacheKey
		_ = token
	}
	state = result.RuleStack
}

// Unsupported Oniguruma constructs, compile failures, and match timeouts are
// non-fatal and can be inspected after lazy compilation/tokenization.
for _, diagnostic := range grammar.Diagnostics() {
	_ = diagnostic
}
```

A state stack is immutable and reusable, but belongs to the grammar that
created it. Pass `nil` or `textmate.InitialState` for the first line. Equal
states returned by the same grammar have the same pointer while both are live,
so consumers may use `*StateStack` as a short-lived cache key and stop
incremental highlighting with a pointer comparison. The grammar holds only a
weak reference to canonical states: if every returned copy is discarded, a
later equal state may receive a new pointer. Pointer identity is not shared
between grammars, and states must not be carried from one grammar into another.
Use `StateStack.Equal` when a previously returned equal pointer is no longer
being retained.

### Interned token scopes

Every token also carries an immutable `*ScopeStack`. Equal scope sequences
produced by one grammar reuse the same pointer and `ScopeStackID`, even when
the tokenizer reconstructed its internal state. Consumers can therefore cache
resolved styles by `token.ScopeStack` or `token.ScopeStack.ID()` instead of
hashing strings for every token. The handle remains valid independently of the
grammar and its accessors are safe for concurrent reads. `Names` returns a
defensive copy; `Len`, `At`, and `Range` inspect names without allocating.

The legacy `Token.Scopes` field remains available and preserves existing JSON
and conformance output. Equal stacks share its backing array, so it must be
treated as read-only. Copy the slice before modifying it.

### Token categories

For consumers that need semantic classes rather than full TextMate theme
resolution, each token and scope stack exposes a cached coarse category:

```go
switch token.Category() {
case textmate.TokenCategoryComment:
	// Apply the application's comment style.
case textmate.TokenCategoryInserted:
	// Apply the application's added-line style.
}
```

`textmate.ClassifyScopes(scopes)` provides the same classification for a raw
scope-name slice. It recognizes common comments, strings, regular expressions,
numbers, keywords, operators, functions, tags, attributes, types, builtins,
variables, markup, diff, punctuation, and invalid scopes. The innermost
recognized scope wins, while reset scopes such as `meta.embedded` prevent an
outer string from coloring embedded source. This is deliberately a stable,
theme-neutral heuristic; it is not a replacement for TextMate selector and
theme resolution.

### Tokenization limits

Editors can bound work on untrusted, generated, or minified input without
changing the behavior of `TokenizeLine`:

```go
result := grammar.TokenizeLineWithOptions(line, state, textmate.TokenizeOptions{
	MaxLineBytes: 20_000,
	MaxLineRunes: 10_000,
	TimeLimit:    5 * time.Millisecond,
})
if result.Stopped {
	log.Printf("tokenization stopped at rune %d: %s", result.StoppedAt, result.StoppedReason)
}
```

Zero or negative limits are unlimited. A line that exceeds either size cap is
not parsed: it receives one fallback token with the incoming scopes, its reset
incoming state is returned unchanged, and `StoppedAt` is zero. Size caps are
checked before allocating the newline-appended regexp input and rune buffer.

`TimeLimit` is deliberately a soft overall budget. It includes lock wait, lazy
root compilation, and begin/while setup, and is checked between regexp
searches. The library cannot interrupt an in-flight regexp search, so a call
can return after the requested duration. A partial result reports the next
unparsed rune in `StoppedAt`, carries the state reached at the last completed
match boundary, and gives the unparsed tail a fallback token so the returned
tokens still cover the complete line.

See [`grammars/README.md`](grammars/README.md) for the curated selection,
regeneration instructions, source/license manifest, and all-versus-curated size
measurements.

The current benchmark and allocation baseline is recorded in
[`docs/performance-baseline.md`](docs/performance-baseline.md). Performance
optimization is tracked as a separate follow-up; the documented warm TSX and
HTML targets are not yet met.

See [`plan.md`](plan.md) for the implementation phases and conformance goals.

## License

MIT. See [`LICENSE`](LICENSE).
