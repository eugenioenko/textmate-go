# textmate-go: decision notes

A running log of what we looked at, what we measured, and why we decided what we did. It's the raw material for an article later, so it keeps the numbers, the dead ends and the reasoning, not just the conclusions. Add to it as decisions are made; don't rewrite history.

---

## 2026-09-25: how we got here

### The trigger: ttt#670

A user reported ([ttt#670](https://github.com/eugenioenko/ttt/issues/670)) that in HTML, Angular templates and JSX/TSX, when a tag's attributes are split over several lines (Prettier's default for long tags), only the first line is coloured:

```html
<button
  type="button"        ← plain text
  class="a"            ← plain text
  aria-hidden="true"   ← plain text
>
```

They guessed the cause correctly: ttt tokenizes each line separately, so the lexer starts every line from scratch and doesn't know it's still inside a tag. They suggested carrying lexer state between lines, mapping file patterns to a lexer (chroma has an Angular2 lexer), and "tree-sitter someday".

### How ttt highlighted at that point

- [chroma](https://github.com/alecthomas/chroma) v2 (pure Go, regex lexers, ~250 languages, a Pygments port).
- One line at a time, with a span cache. About 95µs per line to lex.
- Chroma has no API to resume a lexer mid-state or to report where a line ended. `TokeniseOptions.State` accepts a single starting state name, and the final state stack is never returned.
- History of multi-line fixes in ttt:
  - **#439 (closed):** re-lex the whole buffer on every edit. 68ms per keystroke at 1k lines, 670ms at 10k, 3.26s at 50k. Too slow.
  - **#440 (closed):** carry a per-line "in comment" flag, with a naive string pre-scan. Two different detectors disagreed, so colours depended on scroll position. It also hard-coded `/* */` for every language, so `rm -rf build/*` in a shell script greyed out the rest of the file.
  - **#441 (merged):** block comments. Ask chroma itself where a comment opens (by appending the closer and lexing), memoize per line, and keep a per-line state table cut from the edited line down. Typing cost stayed close to baseline.
  - **#600 (merged):** extended #441 to a list of "regions" found by probing each lexer: block comments, template literals, raw strings, docstrings. It needed several more probes: one block comment per language (otherwise TypeScript got an HTML `<!-- -->` region, and `a <!--b` greyed out the file), whether a backslash escapes the closer (true for JS template literals, false for Go raw strings), and whether the delimiter is its own token (to reject Go/Rust `"""`, which is three ordinary string literals). Swift, Scala and Java lost their `"""` support because their token streams look identical to Rust's.
  - **#599 (closed draft):** loading external chroma lexers from `~/.config/ttt/lexers`; the regions half became #600.
- Every construct needed its own hack, and each hack needed probes to avoid false positives.

### Option 1: one more chroma hack (built, then shelved)

The idea: a "tag" state next to the regions. Detect from chroma's tokens that a line ends inside an open tag (a `<` followed by a `NameTag` token, with no `>` after it). Lex the next line behind a fake `<x ` prefix, so chroma starts in its tag state, and drop the prefix's columns.

What we learned while building it:
- `TokeniseOptions{State: "tag"}` works for attribute lines, but a line consisting only of `>` pops the stack empty and the newline after it comes out as an `Error` token. The prefix trick avoids that, because it gives chroma a real stack to pop.
- Chroma quirks we had to work around:
  - `</span` on its own line lexes as `Error "<"` + `Text "/span"` in HTML, and as `Operator "<"` + `Error "/Button"` in TSX. Needed a regex fallback.
  - TypeScript generics lex as tags: `Map<string, number` → `NameTag "string"`, `NameAttribute "number"`. So does `i <n`. We ignored a `<` directly after an identifier.
  - XML emits `<node` as a single `NameTag` token and `b=` as the attribute token.
  - PHP and Markdown don't treat tags as tags at all.
- Performance, first version: typing cost unchanged, but a cold state-table build for a 10k-line HTML file went from 10ms to **115ms**, because every line containing `<` was lexed. Skipping lines that end in `>` (almost every line of formatted markup) brought it back to **12.7ms**.
- It worked for the reported cases. It also **regressed** one common case: a multi-line JSX handler

  ```tsx
  <Button
    onClick={() => {
      const next = count + 1;   ← coloured as attributes
    }}
  >
  ```

  because the code inside `{…}` was lexed as more attributes. Fixable with brace-depth tracking, which would be yet another piece of state.
- `form.submit()` inside JSX was plain text too, but that's unrelated: chroma's TypeScript lexer emits `NameOther` for every identifier and never marks function calls.

Status: stashed in ttt as "chroma-hack: multiline tags (#670)". It might still ship as a stopgap.

### Option 2: tree-sitter

- The official Go bindings need cgo. ttt's releases are built with `CGO_ENABLED=0` for five targets, and we want to keep static binaries.
- Pure-Go options found:
  - **[odvcencio/gotreesitter](https://github.com/odvcencio/gotreesitter):** a pure-Go reimplementation of the tree-sitter runtime. 206 grammars, incremental reparsing, highlight queries. About 568 stars, 19 contributors, very active, created February 2026. It claims incremental edits are ~90× faster than the C runtime (unverified).
  - [malivvan/tree-sitter](https://github.com/malivvan/tree-sitter): tree-sitter compiled to WebAssembly and run with wazero. Dead since January 2025.
- **Binary size killed it.** gotreesitter's grammar blobs are **15.1MB compressed** for 207 grammars, and the default build is ~24MB. The largest: Nim 646K, Verilog 637K, COBOL 614K, SQL 568K, C++ 406K; TSX 124K, TypeScript 121K. ttt's whole release binary is 15.6MB.
- This isn't specific to gotreesitter: every tree-sitter grammar is a large generated LR parse table. Neovim installs parsers on demand, Helix ships a big runtime folder, Zed bundles a curated set.
- A curated subset (build tags) would work, but it was still the biggest and youngest dependency of all the options, and a rewrite of the highlight package.
- It remains the best option if ttt ever wants structure-aware features (syntax-node selection, outlines, better folding).

### Option 3: fork chroma to carry state

The most direct fix of the root cause.

- The fork change is tiny: `TokeniseOptions.Stack` (the stack to start with) and `TokeniseOptions.EndStack` (receives the stack at the end). 19 lines in `lexer.go` and `regexp.go`. One subtlety: chroma's "unmatched newline resets to the initial state" rule must reset to the root state instead. Nested (`using`) lexers get fresh options, so they can't clobber the output.
- Tested line by line with a carried stack:
  - Work: HTML attributes and comments; TSX attributes; **multi-line JSX handlers** (the case the hack broke); template literals with `${}`; Python docstrings.
  - Don't work: TS/JS, Go and CSS `/* */` comments, Go raw strings, and HTML's split closing tag. These lexers match the construct with **one regex that spans lines**, so the rule itself has to be rewritten as push/pop states (the way chroma's Rust lexer already does comments).
  - TS generics still lex as tags, and the stack now carries that mistake onto the next line.
- **Measured across chroma's own sample files** (whole-file lexing vs line-by-line with a carried stack): of 153 lexers with samples, **109 identical, 44 differ**. Worst: OpenEdge ABL, YAML (multi-line quoted scalars), Markdown (front matter, fences), Org, PowerShell here-strings, OCaml, Svelte, Object Pascal, Terraform heredocs, GraphQL, TSX, Kotlin, SCSS. It undercounts: the samples don't exercise everything (Go and JavaScript showed as clean but break on block comments), and about 125 lexers have no samples. Structurally, **41 lexers** match `/* */` with a single regex.
- Upstream status: nobody has asked chroma for a state API; we searched their issues and PRs. The maintainer is active but selective: small lexer fixes merge within days, while issues and test-only PRs sit. Our earlier issue about an infinite loop in `LexerState.Iterator` ([chroma#1377](https://github.com/alecthomas/chroma/issues/1377)) and its repro PR (#1378) got no reply in two weeks; [#1382](https://github.com/alecthomas/chroma/pull/1382), by someone else, looks like a fix for it. Chroma v3 is in alpha (v3.0.0-alpha.5), with no state API there either. Its tracing feature does record the state per rule, but only as JSON on stderr.
- Verdict: feasible, but "fix about a dozen lexers now and up to 44 eventually, and maintain a fork" turned a bug fix into a project. The exploratory fork is at `~/Documents/chroma`, branch `ttt/stateful`.

### Option 4: other Go highlighters

We searched for Go libraries:
- **[zyedidia/highlight](https://github.com/zyedidia/highlight)** (micro's engine): YAML grammars with native multi-line regions, pure Go. Rejected: too coarse. It doesn't even colour tag names, so it's a downgrade from chroma. [amrnt/highlight](https://pkg.go.dev/github.com/amrnt/highlight) is an abandoned 2017 fork of it.
- **[d4l3k/go-highlight](https://github.com/d4l3k/go-highlight):** highlight.js definitions converted to Go. Old and narrow.
- **[jackielii/go-tree-sitter-highlight](https://pkg.go.dev/github.com/jackielii/go-tree-sitter-highlight):** cgo.
- **[bethropolis/tide](https://github.com/bethropolis/tide):** a Go terminal editor on tree-sitter with 5 languages; an application, not a library.
- Chroma is effectively the only mature pure-Go highlighter; its GitHub "alternatives" are its own forks.

### Option 5: TextMate grammars, the eventual choice

What TextMate grammars are: the grammar format of the 2004 macOS editor TextMate, now the industry standard (VS Code, Sublime's derived format, Shiki, GitHub, Atom). A grammar is JSON (originally plist) with:
- `match` rules for single patterns,
- `begin`/`end` (and `begin`/`while`) rules for regions that can span lines and nest other rules,
- scope names on every token (`entity.name.tag.html`, `comment.block.js`); themes colour by scope.

The engine tokenizes **line by line and carries a stack of open rules** to the next line. That's exactly what chroma lacks and what an editor needs, so #670, block comments, template strings and docstrings all "just work" when the grammar is right.

Go implementations found:
- **[andersonpem/gopher-textmate](https://github.com/andersonpem/gopher-textmate):** pure Go on `regexp2`, with `TokenizeLine(prevState)`. Claims ~20µs per keystroke, ~44ms for a 500-line re-tokenize, and a ~350ms one-time grammar warm-up. 4 stars, 7 commits, created May 2026, ships only HTML and PHP grammars. Too young to depend on, but a useful reference for translating Oniguruma regexes to regexp2.
- **[friedelschoen/go-textmate](https://pkg.go.dev/github.com/friedelschoen/go-textmate):** needs the Oniguruma C library. Out, because of cgo.

The idea that settled it: **VS Code is open source.** Its TextMate engine, [microsoft/vscode-textmate](https://github.com/microsoft/vscode-textmate), is MIT, about 6.2k lines of TypeScript, and comes with conformance tests. Porting it gives:
- highlighting identical to VS Code, which is what most users already know,
- a reference implementation to diff against whenever something's wrong,
- ~260 maintained grammars with mixed per-grammar licenses via [shikijs/textmate-grammars-themes](https://github.com/shikijs/textmate-grammars-themes) (the same grammars VS Code and Shiki use),
- no dependence on a young third-party engine.

Why not TextMate's own engine: the original [textmate/textmate](https://github.com/textmate/textmate) is C++ under GPL-3, and it's not what defines VS Code's behaviour.

### Sizes: why TextMate beats tree-sitter here

| | Raw | Compressed |
|---|---|---|
| tree-sitter, 207 grammars (gotreesitter) | — | **15.1MB** |
| TextMate, 260 grammars (Shiki) | 12MB | **1.4MB** gzip, 0.9MB xz |
| TextMate, curated ~40 grammars | 3.7MB | **0.3MB** gzip |
| chroma's share of ttt today | 2.5MB of lexer XML + 0.4MB of code | — |

TextMate grammars are repetitive JSON and compress very well; tree-sitter grammars are dense parse tables.

Estimated ttt release binary (15.6MB today) if the port replaces chroma: 15.6 − 2.9 (chroma) + 0.4 (engine) + 1.4 (all grammars) ≈ **14.5MB**, or ≈ **13.4MB** with the curated set. It shrinks.

### How big the port is

vscode-textmate at commit `fbe4996` (2026-08-25): 6,232 lines of TypeScript excluding tests.
- To port: about 4,300 lines (`grammar.ts` 1,222, `rule.ts` 899, `tokenizeString.ts` 641, `grammarDependencies.ts` 295, `main.ts` 305, `utils.ts` 205, `matcher.ts` 106, `registry.ts` 102, and so on).
- Skip: `json.ts`/`plist.ts`/`parseRawGrammar.ts` (1,046 lines: Go's `encoding/json` replaces them, and we load JSON grammars only), `theme.ts` (787: ttt maps scopes to its own theme), typings and debug code.
- Expect about 3,500–4,500 lines of Go, plus a regex adapter.
- The hard part is regexes. Grammars use Oniguruma syntax, which Go's RE2 `regexp` can't run (no lookbehind, backreferences or `\G`). `github.com/dlclark/regexp2` (already a chroma dependency) covers most of it. Gaps to handle: `\h`, possessive quantifiers, POSIX classes, some Unicode property names, and subroutine calls `\g<name>` (unsupported). Match timeouts are needed so a catastrophic backtrack can't freeze the editor.
- Deliberate semantic change: offsets are **runes** instead of JavaScript's UTF-16 code units, because regexp2 works on runes and ttt's cursor column is a rune index.

### Naming

We considered `go-textmate` (taken by friedelschoen's cgo version), `gotm`, `scopes`, `tmgo`, `mate`, and colour-themed names. We picked **`textmate-go`** (`github.com/eugenioenko/textmate-go`, package `textmate`): it says what it is, it's searchable, and code reads `textmate.NewRegistry(...)`.

### Testing strategy: don't port the tests, run them

- vscode-textmate's tests are TypeScript (vitest). Much of the value is data:
  - `first-mate/tests.json` (64 tests, 97 lines), `suite1/tests.json` (22 tests, 51 lines), `suite1/whileTests.json` (9 tests, 18 lines): grammar paths, input lines and expected tokens with scopes. The `first-mate` suite comes from Atom's first-mate engine, which vscode-textmate adopted.
  - `themes/`: 72 real sample files, 48 real grammars, and the expected colour of every token under 10 VS Code themes. The best real-world check, but it compares colours.
  - 13 of the 62 fixture grammars are XML plist, not JSON.
- Decision: **don't port vitest or the tests.** Build a Go CLI (a long-running JSON-lines process) and a thin JS adapter implementing vscode-textmate's `Registry`/`IGrammar` on top of it, and run VS Code's own suites unchanged against the Go engine. Benefits:
  - the tests stay VS Code's, with no translation errors and easy updates;
  - plist grammars are converted to JSON on the JS side with vscode-textmate's own parser, so Go stays JSON-only;
  - themes are applied on the JS side with VS Code's own theme code, so there's no need to port `theme.ts` just for testing.
- Safeguard: the adapter has a **reference mode** where the backend is the real vscode-textmate. It must pass 100% before any Go result is trusted, so adapter bugs can't pass for engine bugs.
- **Differential testing with the real vscode-textmate as oracle** (the same idea as the Vim differential fuzzer for ttt-vim): run both engines over a corpus of real files and fuzzed inputs, and treat any token difference as a bug. The oracle provides the expected output for any input, so fuzzing needs no hand-written expectations. This is required for parity, because the fixtures are small.

### Where things were left

- `plan.md` in this folder: the agent-ready implementation plan (phases 0–5; the ttt integration, phase 5, needs the maintainer's go-ahead).
- ttt: the tag hack is stashed. Branch `feat/chroma-stateful-lexing` exists only for the exploration, with an uncommitted `go.mod` `replace` pointing at `../chroma`.
- `~/Documents/chroma`, branch `ttt/stateful`: the Stack/EndStack fork change, kept for reference.
- Open decisions: ship the chroma tag hack for #670 as a stopgap? Embed all 260 grammars or a curated set (decide in phase 4 with measurements)?

---

## Log

Add dated entries below as the port progresses: decisions, surprises, numbers.

### 2026-09-25: Phase 0 — reproducible red baseline

- Cloned the three development references beside the repository and recorded
  their full revisions in `README.md` and `plan.md`.
- The untouched `vscode-textmate` tokenization and theme suites pass 222/222 at
  `fbe4996`. Running the same suites through the reference-backed adapter also
  passes 222/222, plus two JSON-lines transport tests. This validates the
  adapter's import redirects, state pass-through, and reconstructed
  `tokenizeLine2` metadata rather than merely bypassing to the reference API.
- `tokenizeLine` is synchronous, while normal Node child-process pipes are
  asynchronous. The development harness keeps one Go CLI process per Vitest
  worker and performs synchronous reads/writes on Node 24's internal pipe file
  descriptor, retrying `EAGAIN`. This is intentionally isolated to the test
  harness; a worker-thread bridge is the fallback if Node removes `_handle.fd`.
- A failed upstream theme golden test rewrites its `.result` file. Until the Go
  engine is expected to pass, Go mode is restricted to the tokenization suite
  so an expected failure cannot mutate the canonical reference checkout.
- The one-token Go stub passes 3/95 tokenization fixtures coincidentally and
  fails 92. This is the starting conformance metric for the engine port.
- The pinned theme suite contains 14 themes, not the 10 estimated during
  planning.

### 2026-09-25: Phase 1 — regex compatibility is syntax, not substitutions

- Chose `regexp2/v2` v2.8.0. Its ordered mixed-capture option is required to
  preserve Oniguruma's numeric capture order, and the v2 line adds newer
  Unicode tables plus a bounded backtracking stack. RE2 and ECMAScript
  compatibility modes stay off because they change `\w`, `\d`, `\s`, `$`,
  unknown escapes, and POSIX-class behavior.
- Scanned 34,686 regex fields in all 260 Shiki grammars and 8,762 fields in 110
  VS Code fixture grammars. POSIX fragments are the broadest translation need:
  6,434 occurrences in the Shiki set, mostly `alpha` and `alnum`. `\h` appears
  3,509 times. Subroutine calls occur in 23 grammars; class intersections occur
  in 58 patterns across four grammars.
- Oracle probes against the pinned WASM corrected two assumptions in the plan.
  Oniguruma's default `^`/`$` behavior is multiline while its inline `m` flag
  controls dot-all. Also, `{n,m}+` repeats an interval in the default syntax;
  it is not the possessive interval used by Perl/Java syntax. Reversed ranges
  such as `{3,2}` are the possessive form.
- `regexp2` merges duplicate named groups into one slot, unlike Oniguruma, and
  cannot execute recursive `\g<...>` calls or Oniguruma class intersections.
  These patterns are detected, diagnosed, and made non-matching rather than
  allowed to produce plausible but incorrect tokens.
- Each regex evaluation has a 20ms configured timeout, with failures treated
  as no-match and diagnostic emission deduplicated per pattern. regexp2's
  global timeout clock is approximate; the performance phase must measure the
  real pathological cap before claiming a hard latency bound.

### 2026-09-25: Phase 2 — the core port reaches a green first-mate gate

- Ported the grammar object model, rule compiler, scanner caches, immutable
  line state, dependency walker, captures, selector matcher, tokenization
  loop, registry, and CLI integration. Raw grammars are cloned before adding
  synthetic `$self`/`$base` rules, so registry-owned input is never mutated.
- State handles in the CLI retain immutable `*StateStack` values and support
  branching from an older line state. Library token offsets are runes; the
  protocol converts boundaries to UTF-16, including astral characters.
- The Go-only fixture runner reads the pinned upstream manifests directly.
  In strict mode first-mate is 64/64; the unchanged Vitest tests independently
  report the same 64/64. The complete tokenization run is 93/95 upstream
  fixtures, plus 2/2 protocol tests. The two remaining suite1 differences are
  the starting point for phase 3.
- Oracle probes found that Oniguruma does not expose a phantom `^` match after
  the final newline, while regexp2 multiline mode does. The translation now
  preserves interior line starts but rejects that terminal position.
- A Thrift fixture uses recursive `\g<ft>` calls in some alternatives and a
  simple identifier in another. Degrading the whole pattern hid the valid
  branch. Unsupported subroutine atoms are now replaced locally with a
  never-match expression, retaining supported alternatives, capture slots,
  and an explicit diagnostic.
- The race suite, pure-Go build, vet, formatting, and golangci-lint gates pass.
  Performance measurements are intentionally deferred to phase 4.

### 2026-09-25: Phase 3 — complete upstream conformance

- All 95 upstream tokenization fixtures now match vscode-textmate, with the two
  protocol tests bringing the normal Go adapter run to 97/97.
- A separate in-memory differential tokenizes each of the 72 theme fixture
  files once with each engine and compares token text plus full scope stacks.
  It reaches 72/72 parity in about eight seconds without reading or writing
  golden results.
- The upstream theme failure path overwrites expected `.result` files. The Go
  golden gate therefore refuses direct execution: it verifies the exact clean
  pinned revision, copies `src/` and `test-cases/` to a marked `mkdtemp`
  directory, runs only a static list of the 72 golden filenames, enforces a
  process-group timeout, and validates that exact temporary root before
  cleanup. The complete isolated run passes 72/72 across all 14 themes in
  83.24 seconds; the canonical checkout remains clean.
- The last real-world gaps all came from Oniguruma syntax: nested character
  classes, identity-escaped punctuation in classes, named subexpression calls
  used by PHP and generated Swift patterns, and `\G` text inside an inline
  regex comment. Safe acyclic named calls and the common right-recursive form
  are expanded; unresolved complex recursion still degrades with a diagnostic.
- `go test -race ./...`, formatting, vet, golangci-lint, TypeScript typecheck,
  reference-tokenization parity, and the pure-Go build are green. There are no
  known Phase 3 conformance failures.

### 2026-09-25: Phase 4 — real grammars and differential parity

- All 260 pinned Shiki grammars parse and registry-load. The generated report
  scans 34,698 reachable regex fields and records 244 expected unsupported
  syntax diagnostics plus 31 regexp2 compile failures; no load failure is
  hidden.
- An 82-entry MIT-licensed corpus combines 42 focused local regressions with 40
  real samples from the exact pinned Shiki checkout. It covers every curated
  embedded root plus YAML and TOML and reaches 82/82 exact token-text and
  scope-stack parity against vscode-textmate. Both the 48-case CI fuzz budget
  and a 500-case extended run pass.
- The embedded package uses 40 curated root grammars plus one MIT support
  grammar as individually gzipped assets. Deterministic regeneration verifies
  the clean pinned source and fails closed on unreviewed licenses. Curated data
  is 331,365 bytes; compared with all 260 grammars it saves 986,959
  compressed-data bytes and 1,011,712 bytes in a stripped probe binary.
- TSX first use is about 7.83 ms and meets its target. Warm TSX and HTML remain
  above the planned per-line target, and allocation profiling identifies
  regexp2 matches, captures, and runner state as the dominant transient memory
  cost. Performance and memory optimization were explicitly deferred to a
  separate pass; the reproducible baseline is in
  `docs/performance-baseline.md`, and the performance done-checkbox stays open.
- Phase 5 changes the separate ttt repository and still requires explicit
  maintainer approval, so it was not started.

### 2026-09-25: condensed implementation arc

- Began with a one-token stub and pinned clean checkouts of vscode-textmate,
  Shiki's grammar collection, and gopher-textmate. A JSON-lines Go process and
  a thin TypeScript adapter let us run VS Code's tests unchanged instead of
  translating their expected results.
- Built the Oniguruma compatibility layer on regexp2 first, then ported the
  rule compiler, tokenizer, immutable line-state stack, captures, injections,
  dependency loading, registry, and public API. Unsupported regex constructs
  fail locally with diagnostics rather than taking down an entire grammar.
- Closed the gaps incrementally: 3/95 fixtures with the stub, 93/95 after the
  core port, then 95/95 tokenization fixtures and 72/72 real theme files. The
  broader gate finished at 82/82 real and focused files, with a 500-case fuzz
  run also matching vscode-textmate.
- Packaged 40 reviewed root grammars plus one support grammar as independent
  gzip assets. The embedded set is 331,365 bytes and saves about 1.0MB in a
  stripped binary versus embedding all 260 grammars. The finished library is
  pure Go and passes race, vet, lint, and `CGO_ENABLED=0` build gates.

### 2026-09-25: preliminary Chroma performance comparison

The first direct comparison used Chroma v2.24.1 and one logical CPU on the
Ryzen 7 6800H host. Both engines were warmed, every token was consumed, and the
reported values are medians of five samples. This test measures raw line
tokenization; it excludes ttt's result cache and style conversion.

| Input | textmate-go | Chroma | Chroma lead | Allocated per line, TM / Chroma |
|---|---:|---:|---:|---:|
| TSX | 170µs/line | 18.4µs/line | 9.2x | 23.4KB / 8.2KB |
| HTML | 134µs/line | 6.0µs/line | 22.3x | 21.4KB / 3.8KB |
| Go | 48.9µs/line | 23.1µs/line | 2.1x | 13.4KB / 6.3KB |
| Markdown | 59.6µs/line | 25.3µs/line | 2.4x | 12.0KB / 9.2KB |

The corresponding throughput was roughly 5.9k/54.3k TSX lines per second,
7.5k/166.7k HTML, 20.4k/43.3k Go, and 16.8k/39.5k Markdown. TSX had a visible
warm-up split: settled samples were about 165–170µs/line, while early samples
were 315–475µs/line.

These inputs were deliberately small regression fixtures repeated in memory:
428-byte TSX and 444-byte HTML seeds became about 39KB/2,000-line documents;
the 225-byte Go seed became 32KB/2,000 lines; and the 205-byte Markdown seed
became 17KB/1,000 lines. The result is useful as an initial engine baseline,
but the next benchmark should use non-repeated 100KB–5MB Go, JavaScript, TSX,
HTML, C, and Python files and report MB/s plus retained memory.

### 2026-09-25: profiler findings

- The bottleneck is regexp2 plus scanner fan-out: 94.5% of TSX and 91.3% of
  HTML CPU passes through `OnigScanner.FindNextMatch`, which performs about 127
  and 88 independent regex searches per line respectively. The regexp2
  interpreter itself accounts for about 61% of CPU; timeout checks are only
  about 5%, so removing the safety limit would not solve the problem.
- A 2,000-line pass allocates 46.9MB/485,651 objects for TSX and
  42.8MB/484,818 for HTML. Regexp2 is responsible for 82% and 76% of those
  bytes, primarily because every attempted pattern creates public match-text,
  match, group, and capture objects even though textmate-go keeps only rune
  offsets. Token and scope construction is a secondary 5–8% cost.
- Sampled live heap after warming was about 22MB for TSX and 16MB for HTML,
  mostly regexp2 runner stacks, scratch space, and compiled regex programs held
  by per-pattern pools.
- First optimization to try: add a text-free/index-oriented regexp2 API and
  reuse match storage. Then remove redundant token/scope copies, replace the
  hot regex-variant mutex/map with fixed storage, bound later searches after an
  earlier candidate is found, and cache a small number of resolved TSX end
  patterns. Replacing the regex engine remains a later, higher-risk option.

### 2026-09-25: capture-index optimization

A local regexp2 v2.8.0 branch added a capture-index API that keeps match state
on its pooled runner and returns only caller-owned rune ranges. textmate-go
reuses temporary index buffers and no longer constructs regexp2's public text,
match, group, and repeated-capture objects. Full Go/race/vet/pure-Go gates and
the 48-case differential fuzz gate remain green.

| Input | Before | After | Time | Allocated bytes | Allocation count |
|---|---:|---:|---:|---:|---:|
| TSX | 185.2µs/line | 150.3µs/line | -18.8% | -82.1% | -72.1% |
| HTML | 147.3µs/line | 116.5µs/line | -20.9% | -77.2% | -61.9% |
| Go | 48.5µs/line | 37.5µs/line | -22.6% | -83.7% | -73.9% |
| Markdown | 57.9µs/line | 40.3µs/line | -30.4% | -81.7% | -71.7% |

The change removes most transient allocation pressure, but regexp execution
still dominates CPU and TSX still shows one slow early sample. The next CPU
target is reducing the 88–127 separate pattern searches per line. The current
module uses a sibling `../regexp2` replacement; publishing requires upstreaming
the API or hosting a durable fork.
