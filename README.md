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
created it. Pass `nil` or `textmate.InitialState` for the first line. Use
`StateStack.Equal` when incremental highlighting reaches a line whose end
state may already be current.

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
