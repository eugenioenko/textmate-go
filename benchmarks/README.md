# Engine comparison benchmarks

`make bench-compare` measures textmate-go against Chroma and vscode-textmate
on the same corpus and prints one table:

```
| Case | textmate-go | allocs/line | vscode-textmate | chroma | vs JS | vs chroma |
| --- | ---: | ---: | ---: | ---: | ---: | ---: |
| Go_2000 | 9.0 | 31.3 | 13.4 | 19.9 | 0.67x | 0.45x |
```

Times are microseconds per line, the fastest of `COUNT` runs. Ratios are
textmate-go's time over the other engine's, so values below 1 mean
textmate-go is faster.

## Setup

- Go 1.25 and Node 20 or later.
- The tagged `github.com/eugenioenko/regexp2/v2` fork is downloaded normally
  by the Go toolchain; no sibling checkout or `replace` directive is needed.
- vscode-oniguruma, installed by the conformance harness:
  `cd conformance && pnpm install`.
- A compiled [vscode-textmate](https://github.com/microsoft/vscode-textmate)
  checkout next to this repository (the same one the conformance harness uses):

  ```sh
  git clone https://github.com/microsoft/vscode-textmate ../vscode-textmate
  (cd ../vscode-textmate && npm ci && npm run compile)
  ```

## Running

```sh
make bench-compare
make bench-compare-extended
```

The default command runs the 14-case `core` group. The extended command runs
the complete 20-case corpus: core plus C#, PHP, SQL, YAML, Haskell, and Vue.

| Variable | Default | Meaning |
| --- | --- | --- |
| `COUNT` | `5` | Runs per engine and case; the fastest is reported. |
| `BENCH_TIME` | `1s` | Minimum duration of each run (`500ms`, `2s`, ...). |
| `FILTER` | all | Regular expression over case names, e.g. `FILTER='TSX\|CPP'`. |
| `GROUP` | `core` | `core`, `extended`, or `all`; the extended Make target sets `all`. |
| `SKIP_JS` | unset | `1` skips vscode-textmate when it is not set up. |
| `VSCODE_TEXTMATE_DIR` | `../vscode-textmate` | vscode-textmate checkout. |

```sh
FILTER='TSX|HTML' COUNT=3 make bench-compare
GROUP=extended make bench-compare
SKIP_JS=1 make bench-compare
```

## What is measured

All three engines tokenize each case line by line after five untimed warm-up
passes, loading the same embedded grammars from `grammars/data`.
textmate-go and vscode-textmate carry their rule stack between lines. Chroma
tokenizes each line plus a newline from its root state, which mirrors ttt's
lexical hot path but is not multiline-correct, so its column is a throughput
reference rather than an equivalent result.

vscode-textmate runs the same algorithm over Oniguruma compiled to
WebAssembly, so the `vs JS` column is the meaningful target for the port.

## Corpus

[`corpus.json`](corpus.json) lists every case: its group, fixture, the file name
used to select a Chroma lexer, the TextMate scope, and how many lines the fixture
is repeated to. The Go side and the JavaScript side both read it;
`TestCorpusManifest` in `chroma/` checks that each scope matches textmate-go's
filename table and that Chroma has a lexer for it.

TSX, HTML, Go, and Markdown reuse the conformance corpus so their numbers stay
comparable with [`docs/performance-baseline.md`](../docs/performance-baseline.md).
The other languages use the realistic fixtures in [`corpus/`](corpus/).

To add a language, drop a fixture in `corpus/`, add an entry to `corpus.json`,
and run `go test ./...` in `chroma/`.

## Pull-request comparison

On pull requests, CI checks out the exact base and head commits and benchmarks
all 20 cases with the lightweight harness in [`textmate/`](textmate/). It runs
only textmate-go, alternating base and head on the same runner; Chroma,
vscode-textmate, Node tokenization, and Oniguruma are not part of this job.

The workflow writes the full table to the job summary and to one persistent PR
comment. Later pushes update that comment instead of adding another. Benchmark
or corpus errors fail the job, but timing changes are informational because
shared-runner timing is noisy. GitHub gives fork-originated PRs a read-only
token, so those runs keep the report in the job summary without attempting a
comment. To reproduce the comparison locally:

```sh
BASE_DIR=/path/to/base-checkout node benchmarks/compare-go-refs.mjs
```

## Noise

Laptops with boost clocks throttle after a compile, and the first cases in a
run can come out two to three times slower than steady state. The script builds
the Go benchmark binary before timing anything, and reports the fastest run,
but re-run a case with `FILTER` before trusting an outlier.

## Go-only benchmarks

The Go harness in [`chroma/`](chroma/) can also be run directly. It keeps
Chroma, pinned to ttt's v2.24.1, out of textmate-go's production dependency
graph and has two benchmarks:

- `BenchmarkWarmLineByLine`: the comparison above.
- `BenchmarkWarmWholeDocument`: each engine in its correctness-preserving
  mode. TextMate carries state between lines; Chroma lexes the whole document.

```sh
cd benchmarks/chroma
go test -run '^$' -bench '^BenchmarkWarm' -benchmem -count=5
```
