# textmate-go

[![CI](https://github.com/eugenioenko/textmate-go/actions/workflows/ci.yml/badge.svg)](https://github.com/eugenioenko/textmate-go/actions/workflows/ci.yml)

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
- `regexp2`: maintained fork `github.com/eugenioenko/regexp2/v2`, release
  `v2.8.1` (upstream v2.8.0 base `9d0d2ffe88a8b90012f7979ec85424e46d5ef48f`)

The implementation is ported from `vscode-textmate`; its Microsoft MIT license
is reproduced in [`THIRD_PARTY_NOTICES`](THIRD_PARTY_NOTICES).

## Regex engine

TextMate grammars require lookbehind, backreferences, and position-sensitive
anchors that Go's standard RE2-based `regexp` package cannot provide. This port
uses the maintained `github.com/eugenioenko/regexp2/v2` v2.8.1 fork with RE2
compatibility mode disabled. The v2 line was selected for ordered mixed
captures, current Unicode tables, bounded backtracking, and its Go 1.25 API.
Oniguruma-only syntax is translated before compilation; unsupported constructs
are recorded as diagnostics and the affected construct or pattern is safely
degraded.

The fork adds a text-free capture-index API and TextMate-oriented search
optimizations. This avoids constructing public match/group data that the
tokenizer immediately discarded. It has its own module path, so downstream
users receive it as a normal transitive dependency without a `replace`
directive.

## Development

```sh
make test
make bench
make bench-compare   # textmate-go vs vscode-textmate vs Chroma
make bench-compare-extended  # all core + extended language cases
make lint
CGO_ENABLED=0 go build ./...
```

`make bench-compare` needs a compiled vscode-textmate checkout and the
conformance harness's `pnpm install`; see
[`benchmarks/README.md`](benchmarks/README.md) for setup, options, and the
corpus.

## Verification and conformance

The GitHub Actions workflow pins Go, Node, vscode-textmate,
vscode-oniguruma, and the grammar corpus. Every pull request, `main` push,
version tag, and manual dispatch runs the complete merge gate below, including
the 72-file Go theme-golden suite. A green CI check therefore covers every
required test rather than deferring part of the suite to a manual release step.

| Gate | Coverage |
| --- | --- |
| Go correctness | Pure-Go tests and builds for every package, followed by the race detector across every package |
| Static checks | `gofmt`, `go vet`, golangci-lint, TypeScript type-checking, and generated-file cleanliness |
| Native robustness | A bounded 10-second `FuzzTokenizeLine` run on every routine CI execution |
| Upstream tokenization | 101 tests against the real vscode-textmate backend and the same 101 tests against Go: 95 upstream fixtures, two client tests, and four harness tests |
| Themes | All 127 upstream theme tests against the reference adapter, a 72-file Go-vs-reference scope differential, and 72 isolated Go golden fixtures across 14 themes |
| Differential tokenization | Exact token-text and full-scope-stack parity with vscode-textmate on an 82-file corpus, plus 48 deterministic edit/splice/UTF-8 fuzz cases; an extended 500-case gate has also passed |
| Grammar corpus | Standalone tokenization of 166 files covering every embedded root, plus parsing and registry-loading all 260 pinned source grammars and scanning their 34,698 regex fields |
| PR performance | All 20 core + extended cases compare the exact base and PR commits with textmate-go alone; CI updates one persistent PR comment with timings and allocations |

See [`conformance/README.md`](conformance/README.md) for the harness commands,
oracle behavior, fixture counts, and reproducibility details. The checked
[`grammar-report` golden file](cmd/grammar-report/testdata/report.golden.md)
records every known regex translation diagnostic rather than silently treating
unsupported Oniguruma syntax as compatible.

The PR performance comparison is a separate, parallel job. Benchmark and corpus
failures are merge failures; timing deltas are reported but remain informational
because shared GitHub runners do not provide stable enough timing for a hard
regression threshold. See [`benchmarks/README.md`](benchmarks/README.md) for the
groups, local reproduction command, and sticky-comment behavior.

After an intentional regex translator or pinned-grammar change, review the
compatibility difference and refresh that golden file from the repository root:

```sh
TM_GRAMMARS_DIR=../tm-grammars/packages/tm-grammars \
  go test -count=1 ./cmd/grammar-report \
  -run '^TestGrammarReportGolden$' -update
```

The test verifies the pinned checkout revision and cleanliness before writing.
CI sets `TEXTMATE_GO_REQUIRE_GRAMMAR_REPORT=1`, so a missing checkout or stale
golden file fails the merge gate; ordinary `go test ./...` runs skip this
external-corpus test.

## Embedded grammars

The optional `grammars` subpackage embeds 124 license-reviewed grammar roots
plus one MIT support grammar used for Markdown's inline HTML. Assets are
individually compressed and loaded once on demand:

- **A–D:** ActionScript, AutoHotkey, Angular HTML, Assembly, AWK, Ballerina,
  Batch File, BibTeX, Bicep, C, C3, Chapel, Clojure, CMake, COBOL,
  CoffeeScript, Common Lisp, Rocq, C++, Crystal, C#, CSS, CSV, D, Dart,
  Desktop, Diff, Dockerfile, and dotEnv.
- **E–H:** Elixir, Elm, Emacs Lisp, Erlang, Fennel, Fish, F#, GDScript,
  Gherkin, Gleam, Go, GraphQL, Groovy, Handlebars, Haskell, Haxe, HashiCorp
  HCL, HLSL, HTML, HTTP, and Hy.
- **I–M:** INI, Java, JavaScript, Jinja, JSON, JSON with Comments, Jsonnet,
  JSX, Julia, KDL, Kotlin, Lean 4, Less, Lua, Makefile, Markdown, Mojo, and
  MoonBit.
- **N–R:** Nix, nushell, Objective-C, Objective-C++, OCaml, Odin, OpenSCAD,
  Pascal, Perl, PHP, PL/SQL, PowerQuery, PowerShell, Protocol Buffer 3, Puppet,
  Python, QML, R, Windows Registry Script, reStructuredText, Ruby, and Rust.
- **S–Z:** SAS, Sass, Scala, Scheme, SCSS, Shell, Shell Session, GNU Smalltalk,
  Solidity, SQL, Stylus, Svelte, Swift, SystemVerilog, Systemd Units,
  Terraform, TeX, TOML, TSX, Twig, TypeScript, Typst, V, Vala, Visual Basic,
  Verilog, VHDL, Vim Script, Vue, WebAssembly, WGSL, XML, YAML, and Zig.

The internal `html-derivative` grammar supporting Markdown is not counted as a
separate language. Canonical IDs, aliases, filename mappings, scopes, source
revisions, and licenses are recorded in
[`grammars/MANIFEST.md`](grammars/MANIFEST.md).

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

### Editable documents

`Document` owns the state table and token cache needed to highlight an editable
buffer without replaying it from the beginning for every visible line:

```go
document := textmate.NewDocument(grammar, textmate.DocumentOptions{
	TokenizeOptions: textmate.TokenizeOptions{
		MaxLineBytes: 20_000,
		TimeLimit:    5 * time.Millisecond,
	},
	CacheCapacity: 20_000,
})
document.SetLines([]string{"package main", `const message = "hello"`})

line, ok := document.Line(1) // computes any missing states above line 1
if ok {
	for _, token := range line.Tokens {
		_ = token
	}
}

// Replace the half-open range [1, 2). Common-prefix states are retained.
// The unchanged tail is reattached as soon as its canonical state converges.
if err := document.ReplaceLines(1, 2, []string{`const message = "goodbye"`}); err != nil {
	log.Fatal(err)
}

// Explicitly discard state from a line onward when retrying work that stopped
// under a transient time budget.
if err := document.InvalidateFrom(1); err != nil {
	log.Fatal(err)
}
```

`SetLines` and `ReplaceLines` shallow-copy their input slices. `Line`,
`StateAt`, `Len`, and updates are safe to call concurrently. `StateAt` also
accepts `document.Len()` to obtain the state after the final line.

Line results are immutable, library-owned views and remain valid after later
calls or cache eviction. This makes warmed cache hits allocation-free. Copy
`LineResult.Tokens` before changing token fields; `Token.Scopes` and
`Token.ScopeStack` remain read-only. Cache keys combine line text with the
canonical start-state pointer, so identical lines in identical states reuse a
result even after edits. `CacheCapacity` bounds retained results, with zero
selecting the 20,000-entry default and a negative value disabling the cache.

Line-length stops are deterministic for a document's fixed options and are
cached normally. Time-limit stops are transient: `Document` materializes their
outgoing state so a tiny budget cannot trap a random-access request in a retry
loop, but it never caches their line result. Calling `Line` again retries that
line. If a successful retry changes its outgoing state, the document truncates
stale downstream states and reconciles the unchanged tail again.
`InvalidateFrom(i)` explicitly discards materialized state from line `i`
onward and clears the result cache, which is useful before retrying a viewport
whose earlier state was produced by a time-limited line. Passing `Len()` is
valid and clears the cache without discarding an existing line state.

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

## Performance

Warm line-by-line tokenization from `make bench-compare`, in microseconds per
line (2026-09-25, fastest of five runs, Linux/amd64, Go 1.25.1, AMD Ryzen 7
6800H, regexp2 v2.8.1). Each ratio is textmate-go's time over the comparison
engine, so below 1 is faster. Chroma lexes each line from its root state, so its
column is a throughput reference rather than an equivalent result. See
[`benchmarks/README.md`](benchmarks/README.md) for the method and corpus.

| Case | textmate-go | allocs/line | vscode-textmate | Chroma | vs JS | vs Chroma |
| --- | ---: | ---: | ---: | ---: | ---: | ---: |
| TSX | 35.9 | 42.4 | 21.1 | 17.7 | 1.70x | 2.02x |
| HTML | 23.3 | 62.2 | 27.1 | 5.4 | 0.86x | 4.35x |
| Go | 9.8 | 31.3 | 14.2 | 20.4 | 0.69x | 0.48x |
| Markdown | 17.3 | 29.3 | 15.7 | 19.7 | 1.10x | 0.88x |
| TypeScript | 79.7 | 74.1 | 59.7 | 27.2 | 1.33x | 2.93x |
| JavaScript | 69.7 | 75.2 | 68.2 | 27.8 | 1.02x | 2.51x |
| CSS | 25.9 | 53.7 | 78.5 | 11.6 | 0.33x | 2.23x |
| JSON | 8.8 | 52.8 | 7.9 | 9.8 | 1.12x | 0.90x |
| Python | 46.6 | 71.1 | 48.9 | 55.4 | 0.95x | 0.84x |
| Rust | 28.4 | 48.6 | 31.5 | 30.6 | 0.90x | 0.93x |
| Java | 59.9 | 57.1 | 39.8 | 36.3 | 1.50x | 1.65x |
| C++ | 265.5 | 106.3 | 162.9 | 46.2 | 1.63x | 5.75x |
| Ruby | 41.6 | 51.1 | 58.7 | 63.3 | 0.71x | 0.66x |
| Shell | 28.4 | 86.7 | 32.0 | 23.8 | 0.89x | 1.20x |

textmate-go is faster than vscode-textmate in seven of the fourteen cases. C++
is the slowest case: its grammar's very large patterns backtrack heavily in
regexp2's interpreter. The earlier Phase 4 baseline is kept in
[`docs/performance-baseline.md`](docs/performance-baseline.md).

See [`plan.md`](plan.md) for the implementation phases and conformance goals.

## License

MIT. See [`LICENSE`](LICENSE).
