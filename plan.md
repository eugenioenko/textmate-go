# textmate-go: implementation plan

A pure-Go port of Microsoft's [vscode-textmate](https://github.com/microsoft/vscode-textmate), the TextMate grammar engine that powers VS Code's syntax highlighting. Module `github.com/eugenioenko/textmate-go`, package `textmate`.

This document is written so an agent can pick the work up cold. Read it fully before starting.

---

## 1. Why this exists

[ttt](https://github.com/eugenioenko/ttt) is a terminal text editor in Go. It highlights with [chroma](https://github.com/alecthomas/chroma), one line at a time. Chroma has no way to carry lexer state from one line to the next, so every construct that spans lines breaks: block comments, template literals, docstrings, and HTML/JSX tags whose attributes are split over several lines ([ttt#670](https://github.com/eugenioenko/ttt/issues/670)). ttt works around this with probing hacks (`internal/highlight` in ttt, PRs #441 and #600), one hack per construct.

TextMate grammars are designed around line-by-line tokenization with a carried state stack, which is exactly what an editor needs. VS Code, Sublime, Shiki and GitHub all use them, and about 260 maintained grammars exist.

Options already evaluated and rejected:
- **Tree-sitter** (odvcencio/gotreesitter, pure Go): 206 grammars are 15MB compressed; a young, very large dependency.
- **Forking chroma to carry state**: works, but 44+ lexers need their rules rewritten, and we would carry a fork.
- **micro's highlighter**: too coarse (doesn't colour tag names).
- **gopher-textmate** (andersonpem, pure Go, regexp2): same idea, but 4 stars and 7 commits, no conformance suite. Useful as a reference for the regex adapter, not as a dependency.
- **go-textmate** (friedelschoen): needs the Oniguruma C library. We ship with `CGO_ENABLED=0`.

## 2. Goal and definition of done

A pure-Go library (no cgo) that:
1. Loads TextMate grammars in JSON form (`.tmLanguage.json`).
2. Tokenizes one line at a time given the previous line's state, returning tokens with their scope stacks and the state for the next line.
3. Produces the **same tokens as vscode-textmate** on its own conformance fixtures.
4. Is fast enough for an editor (see section 8).

Done means:
- [x] vscode-textmate's 95 tokenization cases and 72 golden theme-file cases pass against this engine through the conformance harness (section 9). The 55 other tests in `themes.test.ts` exercise the deliberately unported upstream theme implementation or unbounded performance smoke paths and are documented exclusions.
- [x] The differential test against the real vscode-textmate shows no token differences on the real-file corpus (section 9.3), except documented deviations.
- [ ] Benchmarks meet the targets in section 8.
- [x] `go vet`, `gofmt`, `golangci-lint` clean; `CGO_ENABLED=0 go build ./...` works.
- [x] README with usage, attribution and license notes.

ttt integration is **not** part of this repo's done criteria. See section 10.

## 3. Reference material

Clone these next to the repo; do not vendor them:

```sh
git clone https://github.com/microsoft/vscode-textmate ../vscode-textmate          # MIT, the source to port
git clone https://github.com/shikijs/textmate-grammars-themes ../tm-grammars       # 260 JSON grammars, mixed per-grammar licenses
git clone https://github.com/andersonpem/gopher-textmate ../gopher-textmate        # MIT, reference for Oniguruma -> regexp2
```

Pin the vscode-textmate commit you port from in the README (at planning time: `fbe4996`, 2026-08-25).

vscode-textmate source, 6,232 lines of TypeScript excluding tests:

| File | Lines | Port |
|---|---|---|
| `src/grammar/grammar.ts` | 1,222 | Yes. Grammar, `StateStack`, `AttributedScopeStack`, `LineTokens` |
| `src/rule.ts` | 899 | Yes. Rule compilation: match, begin/end, begin/while, captures, `RegExpSource`, `RegExpSourceList` |
| `src/grammar/tokenizeString.ts` | 641 | Yes. The tokenizer loop, `\G` anchor handling, while-condition checks |
| `src/grammar/grammarDependencies.ts` | 295 | Yes. Resolving `include` across grammars (`source.js`, `#repo`, `$self`, `$base`) |
| `src/main.ts` | 305 | Yes, as the Go API (`Registry`, `Grammar`, `TokenizeLine`) |
| `src/encodedTokenAttributes.ts` | 246 | Only if implementing `tokenizeLine2` (binary tokens). Optional, see section 5 |
| `src/grammar/basicScopesAttributeProvider.ts` | 111 | No. Binary token metadata is omitted; the JS theme harness reconstructs it from scope stacks |
| `src/utils.ts` | 205 | Partly: capture-reference substitution (`$1`, `${1:/downcase}`), escaping |
| `src/matcher.ts` | 106 | Yes. Scope selectors, used by injections |
| `src/registry.ts` | 102 | Yes |
| `src/rawGrammar.ts` | 76 | Yes, as Go structs with JSON tags |
| `src/diffStateStacks.ts` | 47 | Optional |
| `src/onigLib.ts` | 56 | Yes, as the Go regex interface (section 6) |
| `src/theme.ts` | 787 | **No.** ttt maps scopes to its own theme (section 7), and the theme conformance tests apply themes on the JS side (section 9) |
| `src/json.ts`, `src/plist.ts`, `src/parseRawGrammar.ts` | 1,046 | No. Use `encoding/json`; support JSON grammars only. Several test grammars are XML plist (`.plist`, `.tmLanguage`); the conformance harness converts them to JSON with vscode-textmate's own `parseRawGrammar` before sending them to Go |
| `src/typings`, `src/debug.ts` | 84 | No |

Expect about 3,500 to 4,500 lines of Go, plus the regex adapter and tests.

## 4. Porting rules

- **Port faithfully first, optimise later.** Keep the structure, names and control flow of vscode-textmate recognisable (Go-cased: `StateStack`, `tokenizeString`, `RuleFactory`), so any divergence can be diffed against the TS. Don't redesign while porting.
- **Offsets are rune indices**, not UTF-16 code units (JS) and not bytes. `regexp2` matches over `[]rune` and reports rune indices, and ttt's cursor column is a rune index. This is the one deliberate semantic change. Wherever the TS does UTF-16 arithmetic, convert. Fixtures containing characters outside the BMP (emoji, some CJK) will have UTF-16 offsets or lengths; the test runner compares token **values and scopes**, not offsets, so this mostly doesn't matter. Note any case where it does.
- Pure Go only. `CGO_ENABLED=0` must build.
- Go 1.25+. Dependency: `github.com/dlclark/regexp2/v2` v2.8.0. The v2
  line provides ordered mixed captures, current Unicode tables, bounded
  backtracking, and the Go 1.25 API; the reason is also recorded in the README.
  Nothing else outside the standard library without a reason in the README.
- Code comments: only where missing them would cause a bug or misuse (a hidden invariant, a deliberate deviation from vscode-textmate, an Oniguruma quirk). Don't narrate the port. Do reference the TS function when the Go is a direct translation of something subtle.
- Licensing: vscode-textmate is MIT (Copyright Microsoft). Keep its `LICENSE` text in `THIRD_PARTY_NOTICES` and state in the README that this is a port. Grammars keep their own licenses (Shiki's `NOTICE` file lists them).

## 5. Package layout and API

```
textmate-go/
  go.mod                      module github.com/eugenioenko/textmate-go
  registry.go                 Registry: grammar loading, lookup by scope name, injection map
  grammar.go                  Grammar, TokenizeLine
  statestack.go               StateStack (immutable, comparable), AttributedScopeStack
  rule.go                     rule types + RuleFactory
  regexpsource.go             RegExpSource, RegExpSourceList, anchor (\A, \G) variants
  tokenize.go                 tokenizeString loop, while-condition check
  dependencies.go             include resolution across grammars
  matcher.go                  scope selector matcher (for injections)
  captures.go                 capture handling, $n / ${n:/downcase} substitution
  rawgrammar.go               JSON grammar structs
  oniguruma/                  regex adapter: Scanner interface + regexp2 implementation (section 6)
  cmd/textmate-cli/           JSON-lines server over stdin/stdout, used by the conformance harness (section 9.1)
  conformance/                Node/TypeScript harness: adapter + vitest config that runs vscode-textmate's tests against the Go engine (section 9)
  internal/testsuite/         Go-only quick runner for the JSON fixtures, so basic parity runs without Node (section 9.4)
  grammars/                   embedded, compressed grammar set (section 7)
  THIRD_PARTY_NOTICES
  README.md
```

Public API, mirroring `main.ts`:

```go
type Registry struct { ... }
func NewRegistry(opts RegistryOptions) *Registry
type RegistryOptions struct {
    LoadGrammar func(scopeName string) (*RawGrammar, error)  // lazy loader, like onLoadGrammar
    GetInjections func(scopeName string) []string
}
func (r *Registry) LoadGrammar(scopeName string) (*Grammar, error)

type Grammar struct { ... }
func (g *Grammar) TokenizeLine(line string, prev *StateStack) LineResult

type LineResult struct {
    Tokens    []Token
    RuleStack *StateStack   // pass to the next line
    Stopped   bool          // time limit hit, if a limit is supported
}
type Token struct {
    Start, End int        // rune indices into line
    Scopes     []string   // outermost first, e.g. ["source.js", "string.quoted.double.js"]
}

var InitialState *StateStack   // INITIAL in vscode-textmate
func (s *StateStack) Equal(o *StateStack) bool
```

Requirements on `StateStack`:
- **Immutable** once returned, so an editor can store one per line and reuse it.
- **`Equal`** must be cheap and correct. An editor re-tokenizes from an edited line down and stops as soon as a line's end state equals the one it had before. This is the key to incremental highlighting; vscode-textmate relies on the same property.

`tokenizeLine2` (the packed binary token format with theme metadata) is optional. Only add it if the ttt integration needs it.

## 6. The regex adapter: the riskiest part

TextMate grammars use Oniguruma regex syntax. Go's `regexp` (RE2) can't run them: no lookbehind, no backreferences, no `\G`. `regexp2` is a .NET-style backtracking engine that covers most of it.

Interface, mirroring `onigLib.ts`:

```go
type Scanner interface {
    // FindNextMatch returns the earliest match of any pattern at or after start;
    // ties go to the lowest pattern index (Oniguruma's OnigScanner semantics).
    FindNextMatch(s *String, start int, opts FindOption) *Match
}
type Match struct { Index int; Captures []Capture }  // Capture{Start, End} in runes; -1 if unmatched
```

`OnigScanner` semantics matter: it searches all patterns and returns the one that matches **earliest**, with ties broken by **lowest index**. Implement it by running each regexp2 regex from `start` and keeping the best result, then optimise (see section 8).

Oniguruma → regexp2 translation. Rewrite patterns at compile time and cache the result. Check each item against gopher-textmate's `oniglib`:
- `\G` (anchor to the end of the previous match, i.e. the search start position): regexp2 supports `\G` when matching from a start position; verify it. vscode-textmate also substitutes `\G` itself in `RegExpSource` (the `allowG` variants), so port that logic exactly.
- `\A`, `\z`, `\Z`: supported by regexp2. vscode-textmate's `allowA` handling for `\A` must be ported as-is.
- `\h` (hex digit) → `[0-9a-fA-F]`; `\H` → `[^0-9a-fA-F]`.
- Possessive quantifiers (`*+`, `++`, `?+`) → atomic groups. In Oniguruma's
  default syntax, interval-plus (`{n,m}+`) repeats the interval rather than
  making it possessive, while a reversed interval such as `{3,2}` means an
  atomic `{2,3}`; translate these separately.
- `\x{HHHH}` is accepted by regexp2 v2.8.0. Preserve astral escapes rather
  than narrowing them to four-digit `\u` escapes.
- `\p{...}` Unicode property names: map Oniguruma names to .NET names where they differ.
- POSIX bracket classes `[[:alpha:]]` etc. → explicit classes.
- `(?i)` and other inline options, `(?x)` extended mode with `#` comments:
  regexp2 supports these. Oniguruma's inline `m` means dot-all, so translate it
  to regexp2's `s`; compile with regexp2 `Multiline` to preserve Oniguruma's
  default `^`/`$` behavior. Verify comments and options inside character
  classes.
- Named groups `(?<name>...)`: supported. regexp2 does not natively support
  subroutine calls `\g<name>`. Expand safe acyclic named calls and the common
  generated right-recursive form; diagnose and locally neutralize unresolved
  recursive or otherwise unsafe calls so supported alternatives remain usable.
  The initial corpus scan found calls in 23 of 260 Shiki grammars.
- Character-class intersection and duplicate named captures also differ from
  regexp2. Detect and degrade them rather than silently returning wrong
  matches/capture slots. The corpus scan found class intersection in 58
  patterns across four grammars.
- Backreferences in `end` patterns (`\1`) are **not** regex backreferences: vscode-textmate substitutes the begin captures into the end pattern text (`RegExpSource.resolveBackReferences`). Port that.
- Set `regexp2.RE2` off, and use the options that match Oniguruma defaults (dot doesn't match newline, `^`/`$` per line). Lines are tokenized with a trailing `\n` appended, exactly as vscode-textmate does. Keep that.
- A pattern that fails to compile must not crash the grammar: treat the rule as never matching and collect the error in `Grammar.Diagnostics`. The CLI logs newly observed diagnostics to stderr; the library itself does not write to process-global output.
- **Timeouts:** set `regexp2`'s `MatchTimeout` so catastrophic backtracking in one pattern can't freeze an editor. Treat a timeout as no match, and record it.

Write a table-driven test file for the adapter with one case per construct above. Add cases whenever a grammar exposes a new one.

## 7. Grammars and themes

**Grammars.** Source: `shikijs/textmate-grammars-themes`, `packages/tm-grammars/grammars/*.json` (260 grammars, 12MB raw, about 1.4MB gzipped). Its `NOTICE` lists each grammar's origin and license.
- Phase 1–3: load grammars from disk in tests and benchmarks.
- Phase 4: a `grammars` subpackage embeds 40 license-reviewed roots plus one support grammar with `embed` and exposes `grammars.Load(scopeName)` and `grammars.ForFilename(name)` (extension → scope map). It decompresses lazily per grammar on first use and regenerates deterministically from the pinned tm-grammars commit.

**Themes.** Not needed for the port to be correct: tokens carry scope names. ttt will map scopes to its own theme styles itself, with a small prefix table (`comment` → comment style, `string` → string style, `entity.name.tag` → tag style, `entity.other.attribute-name` → attribute style, and so on; longest prefix wins). Porting `theme.ts` (VS Code theme files and scope-selector matching) is optional, later work, and only worth doing if we want to load VS Code themes directly.

## 8. Performance targets

Reference: chroma costs about 95µs per line on ttt's benchmarks, and ttt tokenizes only the visible lines plus what it needs to carry state down to them.

Targets, measured with `go test -bench` on realistic files (a 2k-line TSX file, a 2k-line HTML file, a 2k-line Go file, a 1k-line Markdown file with fenced code):
- **Per line, warm grammar:** ≤ 100µs average, and no line over 5ms except pathological cases, which the regex timeout must cap.
- **Grammar first use** (decompress, parse JSON, create rules; regex compilation is lazy): ≤ 30ms for TypeScript/TSX, the largest common grammar. Compile regexes lazily per rule the way vscode-textmate does, so a grammar costs only what the text actually exercises. gopher-textmate reports a 350ms warm-up, which is what we must avoid.
- **Allocations:** report them and keep an eye on them; not a hard target at first.

Optimisations to try only after the fixtures pass, keeping the fixtures green after each one:
1. Cache compiled `RegExpSourceList` scanners per rule and per anchor variant (vscode-textmate already does this; port it faithfully).
2. In `FindNextMatch`, stop early once a pattern matches exactly at `start`, since nothing can match earlier and the lowest index wins ties.
3. Remember, per pattern, the last failed search position on the current line: if a pattern found no match from position p, it won't match from any later position on the same string, so skip it (Oniguruma's scanners do a similar cache).
4. Intern scope strings.

## 9. Testing

Parity with VS Code is the goal, so VS Code's own tests are the reference, and **they are not ported**. The vitest suites stay in TypeScript and run against the Go engine through a thin adapter. Node (with pnpm or npm) is a development-only dependency; the library itself stays pure Go.

What vscode-textmate ships (at `fbe4996`):
- `src/tests/tokenization.test.ts` (113 lines) over `test-cases/first-mate/tests.json` (64 tests, 97 lines), `test-cases/suite1/tests.json` (22 tests, 51 lines) and `test-cases/suite1/whileTests.json` (9 tests, 18 lines). Each entry lists grammar files, input lines, and the expected tokens (`value` + `scopes`). 13 of the 62 fixture grammars are XML plist.
- `src/tests/themes.test.ts` (827 lines) with `themeTest.ts`, `themedTokenizer.ts`, `resolver.ts` over `test-cases/themes/`: **72 real sample files** (HTML, Java, LESS, TypeScript, Dockerfile, commit messages, …), **48 real grammars**, and the expected colour of every token under 14 VS Code themes. This is the strongest real-world check.
- `grammar.test.ts` (197), `matcher.test.ts` (58): small unit tests. `json.test.ts`: tests their JSON parser; not relevant.

The `first-mate` fixtures come from Atom's first-mate engine, which vscode-textmate adopted. The original TextMate editor's parser tests (`textmate/textmate`, GPL-3) are not needed: the target is VS Code's behaviour.

### 9.1 The CLI (`cmd/textmate-cli`)

A long-running process speaking JSON lines on stdin/stdout, one request and one response per line, so a test run doesn't spawn a process per call. Minimal protocol (refine as needed; keep it documented in the CLI's README):

```jsonc
// → create a registry from raw grammars (already JSON; the harness converts plist)
{"id":1,"op":"newRegistry","grammars":[{...raw grammar...}, ...],"injections":{"source.x":["scope.injected"]}}
// ← {"id":1,"registry":"r1"}

{"id":2,"op":"loadGrammar","registry":"r1","scopeName":"source.js","embeddedLanguages":{"source.css":2}}
// ← {"id":2,"grammar":"g1"}          or {"id":2,"error":"..."}

{"id":3,"op":"tokenizeLine","grammar":"g1","line":"const a = 1;","state":null}
// ← {"id":3,"tokens":[{"start":0,"end":5,"scopes":["source.js","storage.type.js"]}, ...],"state":"s17"}

{"id":4,"op":"dispose","registry":"r1"}   // frees grammars and states of that registry
```

- `state` is an opaque handle: the CLI keeps a map from handle to `*StateStack` per registry and frees it on `dispose`.
- Offsets in the protocol are **UTF-16 code units**, converted from runes at the CLI boundary, so the JS side's `line.substring(start, end)` works unchanged. The library API itself stays rune-based (section 4).
- Request and grammar-loading failures are returned as `error` fields and also
  written to stderr. Non-fatal regex translation, compilation, and match
  failures are returned as structured `diagnostics`, written to stderr once per
  grammar handle, and degraded to no-match rather than crashing tokenization.

### 9.2 The conformance harness (`conformance/`)

A small Node/TypeScript package that:
1. Depends on a pinned vscode-textmate checkout (git submodule or a path via `VSCODE_TEXTMATE_DIR`). Don't copy its 13MB of fixtures into this repo.
2. Provides `GoRegistry` and `GoGrammar` classes implementing the parts of vscode-textmate's `Registry` and `IGrammar` that the tests use (`loadGrammar`, `loadGrammarWithEmbeddedLanguages`, `setTheme`, `getColorMap`, `tokenizeLine`, `tokenizeLine2`) by calling the CLI.
3. Runs vscode-textmate's `tokenization.test.ts` and `themes.test.ts` **as-is**, with the `../main` import redirected to the adapter (a vitest alias, not an edited copy of the tests). Fall back to a patched copy only if an alias can't work, and keep the diff minimal and documented.

Adapter details that matter:
- **Grammar loading:** the tests load grammars through a resolver (`RegistryOptions.loadGrammar`), often on demand when one grammar includes another. The simplest approach is to have the adapter send every grammar the test's resolver can provide when it creates the registry. Parse plist grammars with vscode-textmate's `parseRawGrammar` and send JSON.
- **`tokenizeLine2` and themes:** the theme tests read colours out of the packed `tokenizeLine2` metadata. The Go engine returns scopes only. The adapter rebuilds `tokenizeLine2` in JS: call the CLI's `tokenizeLine`, resolve each token's scope stack with vscode-textmate's own `Theme` (`theme.match`) and `EncodedTokenAttributes`, and **merge adjacent tokens with identical metadata**, because the real `tokenizeLine2` does and the expected results rely on it.
- **State-stack diffing:** the tests call `diffStateStacksRefEq`/`applyStateStackDiff` on the returned state. That tests a VS Code internal, not tokenization. Make these pass through the opaque handle unchanged, and note it as a known deviation.
- **Validate the adapter first:** before pointing it at Go, run the harness with the **real** vscode-textmate as the backend (the same adapter path, backend swapped). Every test must pass. Then any failure against Go is the Go engine's fault, not the adapter's. Keep this "reference backend" mode as a switch.
- Expected failures go in a checked-in `known_failures.txt`, which starts with everything and should shrink to near zero, so progress is visible.

### 9.3 Differential testing: the real vscode-textmate as oracle

The fixtures are small. Real parity comes from comparing against the actual VS Code engine (vscode-textmate + vscode-oniguruma) on real input. This is **required**, not optional.
- **Corpus:** real files for each language in the embedded grammar set (take them from Shiki's sample files, chroma's `lexers/testdata`, or small open-source repos, with licenses noted), including the ttt#670 cases: HTML/TSX tags whose attributes span lines, multi-line JSX attribute expressions, block comments, template literals, docstrings.
- **Harness:** the same `conformance/` package tokenizes each file with both backends in one process, using the same grammars and carrying state line by line, and reports the first differing token per file (text, scopes, line, column).
- **Fuzzing:** generate inputs (random line splices of corpus files, random edits, random byte strings in valid UTF-8) and compare both backends. The oracle gives the expected output for any input, so this needs no hand-written expectations. Seeded and reproducible; minimise failing inputs before reporting them.
- Run the differential corpus in CI; run the fuzzer for a fixed budget in CI and longer locally.

### 9.4 Go-side tests

- **Quick runner** (`internal/testsuite`, about 100 lines): runs the JSON (non-plist) entries of the three tokenization fixture files straight in Go, so `go test ./...` gives basic parity without Node. Reads fixtures via `VSCODE_TEXTMATE_DIR`, skipping if unset.
- **Unit tests:** the regex adapter table (section 6), capture substitution, the scope matcher (port `matcher.test.ts`, 58 lines), `StateStack.Equal`.
- **Benchmarks** per section 8, with realistic files in `testdata/`.
- **Fuzz** `TokenizeLine` natively (`go test -fuzz`): it must never panic, and tokens must cover the line contiguously from 0 to its length.

## 10. Phases

Do them in order. Each phase ends with its tests green and a short progress note appended to this file (section 12).

**Phase 0: bootstrap.**
- `go mod init github.com/eugenioenko/textmate-go`, add a MIT `LICENSE` (owner: Eugenio Enko), `THIRD_PARTY_NOTICES` with vscode-textmate's license, `.gitignore`, and a `Makefile` with `test`, `bench`, `lint` targets.
- Clone the references (section 3).
- Build the test infrastructure before the engine: `cmd/textmate-cli` with a stub engine that returns one token per line, the `conformance/` harness with its adapter, and the reference-backend switch. **Validate the adapter against the real vscode-textmate first (section 9.2): all tests green in reference mode.** Then switch to Go: nearly everything fails, which is the progress metric.

**Phase 1: regex adapter.** The `oniguruma` package with the `Scanner` interface and regexp2 implementation, the pattern translation, and its unit table. Nothing else depends on grammars yet.

**Phase 2: core port.** `rawgrammar`, `rule`, `regexpsource`, `grammar`/`statestack`, `tokenize`, `dependencies`, `captures`, `registry`. Target: the first-mate fixtures pass through the harness, and the Go quick runner agrees.

**Phase 3: complete conformance.** Injections (`injections`, `injectionSelector`, the matcher), begin/while rules, embedded languages, `$base`/`$self`, suite1 and whileTests, then the themes suite (72 real files). Target: the known-failures list is empty or each remaining entry is documented.

**Phase 4: real grammars, differential testing, packaging, and a performance baseline.** Load all 260 Shiki grammars and report any that fail to load or have patterns that fail to compile (write the list to `docs/grammar-report.md`). Build the differential corpus and fuzzer (section 9.3) and drive the differences to zero. Build the embedded `grammars` package and measure binary size for "all" versus "curated". Record benchmarks and profiles; performance optimization is a separate follow-up, per the user's direction.

**Phase 5 (in the ttt repo, separate work): integration.** Replace the chroma-based internals of ttt's `internal/highlight` behind its existing API (`HighlightLineAt(lines, idx) []Span`). Store a `*StateStack` per line, re-tokenize from the edited line down and stop when the end state equals the stored one, map scopes to `term.Style`, and delete the region-probing code. Do not start this without the maintainer's go-ahead.

## 11. Risks and open questions

- **regexp2 vs Oniguruma gaps** (subroutine calls, some Unicode property names, `\G` edge cases). Mitigation: the adapter table, the grammar report from phase 4, and graceful degradation (an unsupported construct or pattern is neutralized safely; it does not crash).
- **Performance of backtracking regexes in Go.** Mitigation: lazy compilation, scanner caching, match timeouts, benchmarks from phase 4. If per-line cost is far above target, profile before redesigning.
- **UTF-16 vs rune offsets.** Deliberate, documented in section 4. The CLI converts to UTF-16 at the boundary so the JS tests compare like with like.
- **Adapter fidelity.** A wrong adapter would hide or invent engine bugs. Mitigation: the reference-backend mode (section 9.2) must pass 100% before any Go result is trusted, and must stay green whenever the adapter changes.
- **Grammar licenses:** the all-260 set has mixed and sometimes unresolved metadata. Only reviewed compatible grammars may be embedded. `THIRD_PARTY_NOTICES` points to the generated `grammars/MANIFEST.md` and `grammars/NOTICE`, which list the selected grammars, provenance, and license texts without duplicating a drift-prone inventory.
- **Resolved:** embed 40 curated, license-reviewed grammar roots plus one MIT
  support grammar. It adds 331,365 bytes of compressed grammar data versus
  1,318,324 bytes for all 260, saves 1,011,712 bytes in the stripped probe
  binary, and avoids automatically
  redistributing grammars with unknown or incompatible license metadata.

## 12. Progress log

Append an entry after each phase: date, what was done, fixture pass counts (for example `first-mate 36/37`), benchmark numbers, and any decisions or deviations.

- **2026-09-25 — Phase 0 complete.** Bootstrapped the Go 1.25 module,
  licensing/notices, Make targets, pure-Go build, stub JSON-lines CLI, and the
  Node/Vitest conformance adapter. Cloned and pinned `vscode-textmate` at
  `fbe49961ab8077e587fdf5282019655ae69e5f9e`,
  `textmate-grammars-themes` at
  `37edd1b26f18838050661d912334aba0ca7f4931`, and `gopher-textmate` at
  `ce2b42e5386c93ae781add9df2e6328338b06f9e`. The selected upstream
  tokenization-and-theme baseline is 222/222; reference mode through the
  adapter is 222/222 (224/224 including two protocol tests). The one-token Go
  stub is 3/95 on the upstream
  tokenization fixtures (first-mate #13, #65, and #74 pass; 92 fail). Go mode
  excludes the theme suite while failures are expected because upstream's
  failure handler rewrites golden `.result` files. No benchmarks yet. Corrected
  the plan's theme count from 10 to the 14 present at the pinned revision.
- **2026-09-25 — Phase 1 complete.** Added the pure-Go `oniguruma` adapter on
  `regexp2/v2` v2.8.0 with rune offsets, ordered multi-pattern scanning,
  captures, A/G/z/Z find options, syntax translation, per-pattern timeouts,
  and deduplicated diagnostics. The table suite covers offsets, earliest/tie
  selection, zero-width/unmatched/repeated captures, lookaround/backreferences,
  POSIX and Unicode classes, inline/extended options, possessive and interval
  semantics, anchors, compile failures, and timeout continuation. Quality
  gates: `go test ./...`, repeated focused tests, `go test -race ./...`,
  `make lint`, and `CGO_ENABLED=0 go build ./...` all pass; adapter statement
  coverage is 86.7%. Safe named subroutine calls are expanded; unresolved
  recursion, character-class intersection, and duplicate named captures are
  deliberately diagnosed and degraded. A scan of 260
  Shiki grammars (34,686 regex fields) and 110 VS Code fixture grammars (8,762
  fields) corrected the original plan's interval-plus and inline-`m`
  assumptions. The conformance harness was also corrected to pin
  `vscode-oniguruma` 1.7.0, matching the reference checkout's find-option API.
- **2026-09-25 — Phase 2 complete.** Ported the raw grammar model, immutable
  state and attributed-scope stacks, rule factory and scanner caches,
  dependency/include resolution, captures, selector matching, tokenizer loop,
  grammar, registry, and the real CLI/state bridge. Added the Go fixture
  runner and preserved rune offsets in the library with UTF-16 conversion only
  at the JSON-lines boundary. The first-mate gate passes 64/64 both through
  the unchanged Vitest suite and the strict Go runner; the full tokenization
  harness is 93/95 upstream fixtures (95/97 including two protocol tests),
  leaving two suite1 cases for phase 3. Quality gates pass under the race
  detector, `golangci-lint`, `go vet`, and `CGO_ENABLED=0`. Two oracle-driven
  regex corrections were required: a final newline is not a second `^` line,
  and unsupported `\g<...>` subroutine atoms are neutralized locally so valid
  alternatives and capture numbering remain usable. No phase-2 benchmarks;
  performance remains a phase-4 gate.
- **2026-09-25 — Phase 3 complete.** Completed parity for injections,
  begin/while state, embedded languages, external includes, `$base`/`$self`,
  and all upstream tokenization fixtures: 95/95 upstream cases pass (97/97
  including the protocol tests). A read-only differential across the 72 theme
  fixture sources reports 72/72 exact text-and-scope parity. The protected
  golden runner copies the pinned clean reference into a marked temporary
  directory before invoking the unchanged upstream theme suite, so failure
  rewrites cannot touch the canonical checkout; all 72 files pass under all 14
  themes (55 non-golden utility/performance tests are deliberately excluded
  from this Go gate). Regex compatibility fixes covered nested Oniguruma
  character classes, literal punctuation escapes in classes, safe named
  subexpression-call expansion (including a generated right-recursive form),
  and anchor text inside inline regex comments. The known-failures list is
  empty. Race, lint, typecheck, reference-tokenization, and `CGO_ENABLED=0`
  build gates pass. Performance remains a phase-4 gate.
- **2026-09-25 — Phase 4 functional scope complete; performance deferred.**
  Added a deterministic all-grammar report: all 260 Shiki grammars parse and
  registry-load, while 34,698 reachable regex fields produce 244 documented
  unsupported-syntax diagnostics and 31 regexp2 compile failures. Added a
  82-entry MIT corpus: 42 focused regressions plus 40 pinned real Shiki samples
  covering every embedded root and YAML/TOML, including the ttt#670 multiline
  cases. It has exact token-text and full-scope parity with vscode-textmate on
  82/82 files; the deterministic
  fuzzer passes its 48-case CI budget and a 500-case extended run. Oracle-found
  fixes cover nested repeat operators, leading `]` in nested classes,
  escaped-hyphen range endpoints, and `\G`-sensitive scanner caching. The
  production `grammars` package embeds 40 roots plus one support grammar as
  individually gzipped, MIT-licensed assets with lazy concurrent loading,
  filename lookup, pinned deterministic
  regeneration, notices, and a fail-closed license check. The curated payload
  is 331,365 bytes and its stripped probe binary is 1,011,712 bytes smaller
  than the all-260 build. TSX first use is 7.83 ms, below the 30 ms target.
  Warm TSX (391.5 us/line) and HTML (131.7 us/line) remain above target, with
  substantial regexp2-driven allocation rates; the user explicitly deferred
  latency and memory optimization to a separate pass, so the benchmark
  definition-of-done item remains open. The reproducible baseline and memory
  profile summary are in `docs/performance-baseline.md`. Phase 5 remains out of
  scope until the maintainer gives the required go-ahead.
- **2026-09-25 — Phase 4 completion audit passed.** Closed the remaining
  correctness and reproducibility gaps: added native state-carrying fuzzing,
  public and CLI regex diagnostics, complete registry disposal, fail-closed
  staged grammar generation, exact clean-revision report guards, enforceable
  theme-scope parity, and pinned CI. The 82-entry manifest now structurally
  requires exactly one real pinned sample for each of the 40 embedded roots.
  A standalone gate tokenizes the 80 applicable corpus entries through the
  embedded package and compares every token against the full source checkout;
  YAML and TOML are the two explicit license-review skips. Markdown parity
  required one additional MIT `html-derivative` support asset, bringing the
  package to 40 roots plus one support grammar. CI regenerates and checks the
  grammar assets and compatibility report, then runs the bounded corpus/fuzz
  gates. Final evidence: 101/101 reference-tokenization and 101/101 Go
  tokenization tests, 72/72 scope differentials, 72/72 isolated theme goldens,
  82/82 corpus parity, 500 deterministic differential fuzz cases, a 10-second
  native fuzz run, standalone embedded parity, race tests, lint, typecheck,
  pure-Go build, byte-identical regeneration/report reproduction, and clean
  exact reference revisions all pass. A subsequent requirement-level audit
  made the Go-only fixture runner fail hard on any detected mismatch, fixed
  extended-mode comment false positives in regex compatibility checks,
  restored the pinned Microsoft license text byte-for-byte, and added an
  isolated 127/127 reference-theme CI gate for the complete adapter path. CI
  also runs the race detector and a bounded native `TokenizeLine` fuzz session,
  making those correctness claims reproducible rather than progress-log-only.
  Enabling the durable race gate exposed and removed a test-only mutation of
  regexp2's unsafe process-global timeout-clock period; ten repeated timeout
  tests and the complete race suite pass with the default clock.
  Performance remains the deliberately separate follow-up, and Phase 5 remains
  gated on maintainer approval.
