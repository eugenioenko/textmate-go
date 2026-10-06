# textmate-go

[![CI](https://github.com/eugenioenko/textmate-go/actions/workflows/ci.yml/badge.svg)](https://github.com/eugenioenko/textmate-go/actions/workflows/ci.yml)

`textmate-go` is a pure-Go port of Microsoft's
[`vscode-textmate`](https://github.com/microsoft/vscode-textmate). It provides
line-oriented TextMate grammar tokenization, reusable state stacks, and
incremental document highlighting without CGO.

- Exact token-text and scope-stack parity is continuously checked against
  vscode-textmate.
- The optional `grammars` package provides 124 bundled languages.
- Immutable states and interned scope stacks support efficient editor caches.
- Size and time limits bound work on untrusted or minified input.
- Go 1.25 or newer is required. `CGO_ENABLED=0` is supported.

The project is under active development, so its API may evolve before a stable
release.

## Installation

```sh
go get github.com/eugenioenko/textmate-go@latest
```

## Quick start

```go
package main

import (
	"log"

	textmate "github.com/eugenioenko/textmate-go"
	"github.com/eugenioenko/textmate-go/grammars"
)

func main() {
	registry := textmate.NewRegistry(textmate.RegistryOptions{
		LoadGrammar: grammars.Load,
	})
	defer registry.Dispose()

	grammar, err := registry.LoadGrammar(grammars.ScopeForFilename("main.go"))
	if err != nil {
		log.Fatal(err)
	}
	if grammar == nil {
		log.Fatal("no embedded grammar for main.go")
	}

	var state *textmate.StateStack
	for _, line := range []string{"package main", `const message = "hello"`} {
		result := grammar.TokenizeLine(line, state)
		for _, token := range result.Tokens {
			log.Printf("%d:%d %v", token.Start, token.End, token.Scopes)
		}
		state = result.RuleStack
	}
}
```

Token offsets are rune offsets. Scopes are ordered from the root grammar to the
most specific matched scope. Pass `nil` or `textmate.InitialState` for the first
line, then carry `RuleStack` to the next line.

See the [usage guide](docs/usage.md) for editable documents, state identity,
interned scope stacks, token categories, and tokenization limits.

## Embedded grammars

The optional `grammars` package embeds 124 license-reviewed grammar roots plus
one MIT support grammar used for Markdown's inline HTML. Assets are individually
compressed and loaded once on demand.

- **A-D:** ActionScript, AutoHotkey, Angular HTML, Assembly, AWK, Ballerina,
  Batch File, BibTeX, Bicep, C, C3, Chapel, Clojure, CMake, COBOL,
  CoffeeScript, Common Lisp, Rocq, C++, Crystal, C#, CSS, CSV, D, Dart,
  Desktop, Diff, Dockerfile, and dotEnv.
- **E-H:** Elixir, Elm, Emacs Lisp, Erlang, Fennel, Fish, F#, GDScript,
  Gherkin, Gleam, Go, GraphQL, Groovy, Handlebars, Haskell, Haxe, HashiCorp
  HCL, HLSL, HTML, HTTP, and Hy.
- **I-M:** INI, Java, JavaScript, Jinja, JSON, JSON with Comments, Jsonnet,
  JSX, Julia, KDL, Kotlin, Lean 4, Less, Lua, Makefile, Markdown, Mojo, and
  MoonBit.
- **N-R:** Nix, nushell, Objective-C, Objective-C++, OCaml, Odin, OpenSCAD,
  Pascal, Perl, PHP, PL/SQL, PowerQuery, PowerShell, Protocol Buffer 3, Puppet,
  Python, QML, R, Windows Registry Script, reStructuredText, Ruby, and Rust.
- **S-Z:** SAS, Sass, Scala, Scheme, SCSS, Shell, Shell Session, GNU Smalltalk,
  Solidity, SQL, Stylus, Svelte, Swift, SystemVerilog, Systemd Units,
  Terraform, TeX, TOML, TSX, Twig, TypeScript, Typst, V, Vala, Visual Basic,
  Verilog, VHDL, Vim Script, Vue, WebAssembly, WGSL, XML, YAML, and Zig.

The internal `html-derivative` grammar supporting Markdown is not counted as a
separate language. Canonical IDs, aliases, filename mappings, scopes, source
revisions, and licenses are recorded in
[`grammars/MANIFEST.md`](grammars/MANIFEST.md). Regeneration instructions and
all-versus-curated size measurements are in
[`grammars/README.md`](grammars/README.md).

## Regex compatibility and diagnostics

TextMate grammars require lookbehind, backreferences, and position-sensitive
anchors that Go's standard RE2-based `regexp` package cannot provide. This port
uses the maintained `github.com/eugenioenko/regexp2/v2` fork without RE2
compatibility mode. The fork adds a text-free capture-index API and
TextMate-oriented search optimizations.

Oniguruma-only syntax is translated before compilation. Unsupported constructs,
compile failures, match panics, and timeouts are non-fatal and available through
`grammar.Diagnostics()` so one bad pattern does not stop the tokenizer.

The fork has its own module path, so downstream users receive it as a normal
transitive dependency without a `replace` directive.

## Verification

The GitHub Actions merge gate runs on pull requests, `main`, version tags, and
manual dispatches. A green required check covers:

- formatting, vet, lint, generated-file cleanliness, pure-Go builds, unit
  tests, race tests, and a bounded fuzz run;
- all 101 upstream tokenization tests against vscode-textmate and textmate-go;
- exact differential tokenization across the real-file corpus and deterministic
  edit, splice, and UTF-8 fuzz cases;
- upstream theme tests, 72 theme scope differentials, and 72 isolated golden
  fixtures across 14 themes;
- standalone tokenization of all 166 embedded roots, plus parsing and loading
  all 260 pinned source grammars and scanning their 34,698 regex fields;
- base-versus-PR performance results for all 20 benchmark cases, posted as one
  updated PR comment.

Benchmark timing is informational because shared runners are noisy. Correctness,
corpus, and benchmark execution failures block the merge.

See [`conformance/README.md`](conformance/README.md) for commands and oracle
behavior, [`grammars/README.md`](grammars/README.md) for updating the grammar
report golden file, and [`benchmarks/README.md`](benchmarks/README.md) for local
benchmark reproduction and CI comment behavior.

## Performance

Warm line-by-line tokenization from `make bench-compare-extended`, in
microseconds per line, using the fastest of five Linux/amd64 runs with Go
1.25.1 on an AMD Ryzen 7 6800H. The textmate-go column was measured on
2026-10-06 with regexp2 v2.8.2; the vscode-textmate and Chroma columns are from
2026-09-25 on the same machine. Ratios below 1 mean textmate-go is faster.

Chroma starts each line from its root state, so its column is a throughput
reference rather than an equivalent incremental result.

| Case | textmate-go | allocs/line | vscode-textmate | Chroma | vs JS | vs Chroma |
| --- | ---: | ---: | ---: | ---: | ---: | ---: |
| TSX | 19.8 | 24.0 | 20.8 | 17.1 | 0.95x | 1.16x |
| HTML | 19.7 | 37.4 | 26.0 | 5.2 | 0.76x | 3.78x |
| Go | 8.0 | 17.9 | 14.0 | 20.3 | 0.57x | 0.39x |
| Markdown | 15.0 | 19.3 | 15.3 | 19.1 | 0.98x | 0.79x |
| TypeScript | 60.1 | 36.6 | 58.2 | 26.3 | 1.03x | 2.29x |
| JavaScript | 52.2 | 36.2 | 64.3 | 26.7 | 0.81x | 1.95x |
| CSS | 23.3 | 31.9 | 74.0 | 11.2 | 0.31x | 2.08x |
| JSON | 7.9 | 29.4 | 7.6 | 9.3 | 1.04x | 0.85x |
| Python | 34.1 | 34.2 | 48.4 | 53.6 | 0.70x | 0.64x |
| Rust | 17.2 | 26.5 | 31.2 | 29.7 | 0.55x | 0.58x |
| Java | 42.0 | 30.8 | 38.9 | 35.5 | 1.08x | 1.18x |
| C++ | 164.1 | 67.8 | 158.0 | 45.0 | 1.04x | 3.65x |
| Ruby | 25.2 | 25.9 | 59.0 | 62.6 | 0.43x | 0.40x |
| Shell | 25.1 | 42.7 | 31.8 | 23.0 | 0.79x | 1.09x |
| C# | 24.3 | 24.7 | 24.9 | 14.1 | 0.98x | 1.73x |
| PHP | 13.1 | 18.6 | 24.6 | 10.9 | 0.53x | 1.20x |
| SQL | 11.8 | 15.2 | 41.9 | 25.5 | 0.28x | 0.46x |
| YAML | 8.2 | 17.8 | 11.2 | 7.7 | 0.73x | 1.07x |
| Haskell | 6.3 | 15.2 | 9.0 | 7.2 | 0.70x | 0.88x |
| Vue | 21.0 | 31.2 | 28.6 | 14.9 | 0.74x | 1.41x |

textmate-go is faster than vscode-textmate in sixteen of the twenty cases and
within 8% in the other four: TypeScript, JSON, Java, and C++. See
[`benchmarks/README.md`](benchmarks/README.md) for the method and corpus, and
[`docs/performance-baseline.md`](docs/performance-baseline.md) for the earlier
Phase 4 baseline.

## Development

```sh
make test
make bench
make bench-compare
make bench-compare-extended
make lint
CGO_ENABLED=0 go build ./...
```

`make bench-compare` needs a compiled vscode-textmate checkout and the
conformance harness dependencies. See
[`benchmarks/README.md`](benchmarks/README.md) for setup and options.

Pinned revisions used by development and conformance testing:

- `vscode-textmate`: `fbe49961ab8077e587fdf5282019655ae69e5f9e`
- `textmate-grammars-themes`: `37edd1b26f18838050661d912334aba0ca7f4931`
- `gopher-textmate`: `ce2b42e5386c93ae781add9df2e6328338b06f9e`
- `regexp2`: maintained fork `github.com/eugenioenko/regexp2/v2`, release
  `v2.8.1`, based on upstream v2.8.0 commit
  `9d0d2ffe88a8b90012f7979ec85424e46d5ef48f`

See [`plan.md`](plan.md) for implementation phases and conformance goals.

## Acknowledgments and thanks

This project would not exist without Microsoft's
[`vscode-textmate`](https://github.com/microsoft/vscode-textmate), whose
implementation defines the behavior that this Go port follows, and Shiki's
[`textmate-grammars-themes`](https://github.com/shikijs/textmate-grammars-themes),
which provides the pinned grammar corpus, source provenance, and license
metadata used to build the optional embedded grammar package. Thanks also to
the language-extension and grammar authors whose work is collected there.

The tokenizer relies on the work behind
[`regexp2`](https://github.com/dlclark/regexp2) for the regular-expression
features TextMate grammars require. Microsoft's
[`vscode-oniguruma`](https://github.com/microsoft/vscode-oniguruma) provides an
essential compatibility oracle. We are grateful to all of these projects and
their contributors.

The Microsoft MIT license for vscode-textmate is reproduced in
[`THIRD_PARTY_NOTICES`](THIRD_PARTY_NOTICES).

## License

MIT. See [`LICENSE`](LICENSE).
