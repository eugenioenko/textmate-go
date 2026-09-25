# Embedded grammars

This package embeds 40 mainstream TextMate grammar roots plus the small MIT
`html-derivative` support grammar required for Markdown inline HTML, all as
individually compressed JSON assets. `Load` decompresses and parses a grammar
only on first use, then caches the result. `ForFilename` maps conventional
filenames and extensions to the root scopes.

`Infos` exposes the pinned catalog metadata for the embedded set: each
grammar's canonical language ID, display name, scope, aliases, and declared
file types. `InfoForFilename`, `InfoForID`, `InfoForAlias`, and `InfoForScope`
keep the different lookup namespaces explicit; in particular, Markdown code
fences can try the canonical ID first and then the alias without treating an
alias as a filename. ID and alias matching is case-insensitive. Returned
metadata owns its slices and is safe for callers to modify.

The selection is listed in `curated.txt`. Regenerate it from the pinned sibling
checkout with:

```sh
go generate ./grammars
```

Generation verifies the checkout revision and its clean status before writing
output. It validates every selected grammar, metadata record, scope name,
license, and provenance field, then builds the complete result in a staging
directory before replacing generated files. `SOURCE` records the exact source
and revision, `MANIFEST.md` records each embedded grammar's source commit and
declared license, and `NOTICE` preserves the upstream per-grammar license
notices. Any selection containing a missing or unreviewed license fails without
changing the existing generated package; this also means `-selection all`
deliberately fails at the pinned revision.

The curated set is recommended over embedding all 260 grammars: it covers the
common web, Go, Python, Rust, JVM, C-family, shell, configuration, and Markdown
use cases while keeping the application payload substantially smaller. The
all-260 corpus is not suitable for automatic redistribution without a separate
license review: its pinned metadata contains 194 MIT, 22 Apache-2.0, 21 unknown,
12 NOASSERTION, five GPL-3.0, three MPL-2.0, and one each GNU, ISC, and
BSD-3-Clause entries. Exact comparative measurements are recorded below;
`SOURCE` records the reproducible current curated totals.

| Pinned set | Grammars | Compacted JSON | Per-file gzip | Stripped probe binary |
| --- | ---: | ---: | ---: | ---: |
| Curated | 41 | 2,485,124 B | 331,365 B | 2,740,372 B |
| All | 260 | 7,654,872 B | 1,318,324 B | 3,752,084 B |

The binary figures were measured on Linux/amd64 with Go 1.25.1 using a minimal
tracked `cmd/grammar-size-probe`, built with
`CGO_ENABLED=0 go build -trimpath -ldflags='-s -w -buildid=' -o <output>
./cmd/grammar-size-probe`. The current curated figure is directly reproducible;
the all-260 figure is the Phase 4 measurement made with the same minimal API
calls before the fail-closed generator prohibited staging unreviewed licenses.
The curated build saves 986,959 compressed-data bytes and 1,011,712 binary
bytes. Because the all-260 set is roughly four times the compressed payload and
carries unresolved license metadata, the curated set is the production default.

Some curated roots contain optional includes for other languages, such as
Markdown fenced-code grammars, Vue preprocessors, C/C++ assembly dialects, and
Ruby template variants. The registry deliberately treats an unavailable
external include as non-fatal. The exact sorted set is checked in at
`OPTIONAL_DEPENDENCIES` and guarded by a test, so adding or removing an
upstream dependency cannot silently change the standalone package. The opt-in
`TestEmbeddedCorpus` gate tokenizes the 40 focused regression files whose
grammars are embedded and all 40 pinned upstream samples using only
`grammars.Load`; the two supplemental YAML/TOML regressions are explicitly
skipped because those grammars await license review. CI enables the gate with
`TEXTMATE_GO_REQUIRE_EMBEDDED_CORPUS=1`.
