# Performance baseline

Performance optimization is intentionally deferred to a separate pass. This
document preserves the Phase 4 baseline so that pass starts with reproducible
latency and memory measurements rather than anecdotal results.

The benchmark corpus repeats checked-in realistic source files to 2,000 lines
for TSX, HTML, and Go, and 1,000 lines for Markdown. It primes the grammar and
tokenizer before timing the warm path. Run the baseline on one logical CPU to
reduce scheduler noise:

```sh
GOMAXPROCS=1 CGO_ENABLED=0 go test -run '^$' \
  -bench '^(BenchmarkTokenizeWarm|BenchmarkFirstUseTSX)$' \
  -benchmem -benchtime=3x -count=1 -cpu=1 ./...
```

Measurements below were taken on Linux/amd64 with Go 1.25.1 on an AMD Ryzen 7
6800H. They are a local comparison point, not portable absolute guarantees.

| Benchmark | Average | Slowest observed line | Bytes/op | Allocs/op |
| --- | ---: | ---: | ---: | ---: |
| TSX, 2,000 lines | 391.5 us/line | 7.67 ms | 46,866,410 | 485,655 |
| HTML, 2,000 lines | 131.7 us/line | 4.19 ms | 42,818,258 | 484,802 |
| Go, 2,000 lines | 50.2 us/line | 1.55 ms | 26,896,370 | 339,155 |
| Markdown, 1,000 lines | 40.6 us/line | 7.93 ms | 10,230,648 | 128,650 |
| TSX first use | 7.83 ms/op | n/a | 2,803,434 | 27,623 |

The first-use target is met. Go and Markdown meet the average warm target; TSX
and HTML do not, and isolated lines can exceed 5 ms when garbage collection or
scheduler work lands in the per-line timing window. The plan's performance
checkbox therefore remains open.

The transient allocation rate is substantial: approximately 23.4 KiB and 243
allocations per TSX line, 20.9 KiB and 242 allocations per HTML line, 13.1 KiB
and 170 allocations per Go line, and 10.0 KiB and 129 allocations per Markdown
line. An allocation profile attributes most bytes to regexp2 match/runner
creation and capture materialization. A sampled post-benchmark heap retained
about 19 MiB in this process, principally compiled regexp2 programs and runner
pools; that number includes benchmark/runtime state and is not an API memory
reservation guarantee.

CPU profiling likewise places the hot path in regexp2 evaluation. Instrumented
scans observed about 127 regexp searches per TSX line and 88 per HTML line,
with no regex compilation in the measured warm region. The separate
performance pass should measure both latency and retained/transient memory,
keep match timeouts for backtracking patterns, and rerun the conformance and
differential gates after every scanner change.
