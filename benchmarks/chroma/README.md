# Chroma comparison benchmark

This nested module keeps Chroma out of textmate-go's production dependency
graph. It pins Chroma v2.24.1, the version in ttt's committed `go.mod`, and
compares both engines on the same repeated TSX, HTML, Go, and Markdown corpus.

The harness exposes two comparisons:

- `BenchmarkWarmWholeDocument` uses each engine's correctness-preserving mode.
  TextMate carries state between lines; Chroma receives the entire document.
- `BenchmarkWarmLineByLine` mirrors ttt's lexical hot path: Chroma tokenizes
  each line plus a newline from the root state, while TextMate carries state.

Both sides consume every token and run five untimed warm-up passes. The
line-by-line result is intentionally a throughput and allocation comparison,
not a claim of equivalent multiline correctness: released Chroma cannot return
its lexer stack, so ttt supplements it with region-probing logic that is
outside this benchmark.

Run on one logical CPU to reduce scheduler noise:

```sh
cd benchmarks/chroma
GOMAXPROCS=1 go test -run '^$' -bench '^BenchmarkWarm' \
  -benchmem -benchtime=3x -count=5 -cpu=1
```
