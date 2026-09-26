# Embedded grammars

This package embeds 124 TextMate grammar roots plus the small MIT
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

Accepted licenses are MIT, Apache-2.0, BSD-3-Clause, ISC, MPL-2.0, and the
TextMate bundle grant ("Permission to copy, use, modify, sell and distribute
this software is granted..."). A grammar whose upstream metadata states no
license or `NOASSERTION` can be embedded only with an entry in
`license-reviews.json`, which records the license, a link to the evidence at the
pinned upstream commit, and the license text; the generator appends that text
to `NOTICE`. It currently covers YAML (MIT, per its `YAML-license.txt`), TOML
(the TextMate bundle grant), Elixir (Apache-2.0 header), and Sass (MIT). An
entry for a grammar that already states a license, or that is not selected, is
an error. GPL-licensed grammars are deliberately excluded.

The selection covers the original 40 mainstream languages plus every
permissively licensed grammar for a language that Chroma highlights, so
switching from Chroma loses little coverage. The all-260 corpus is not suitable
for automatic redistribution without a separate license review: its pinned
metadata contains 194 MIT, 22 Apache-2.0, 21 unknown, 12 NOASSERTION, five
GPL-3.0, three MPL-2.0, and one each GNU, ISC, and BSD-3-Clause entries. Exact
comparative measurements are recorded below; `SOURCE` records the reproducible
current curated totals.

| Pinned set | Grammars | Compacted JSON | Per-file gzip | Stripped probe binary |
| --- | ---: | ---: | ---: | ---: |
| Curated | 125 | 4,648,925 B | 785,956 B | 3,211,412 B |
| Original 40 | 41 | 2,485,124 B | 331,365 B | 2,740,372 B |
| All | 260 | 7,654,872 B | 1,318,324 B | 3,752,084 B |

The binary figures were measured on Linux/amd64 with Go 1.25.1 using a minimal
tracked `cmd/grammar-size-probe`, built with
`CGO_ENABLED=0 go build -trimpath -ldflags='-s -w -buildid=' -o <output>
./cmd/grammar-size-probe`. The current curated figure is directly reproducible;
the all-260 figure is the Phase 4 measurement made with the same minimal API
calls before the fail-closed generator prohibited staging unreviewed licenses.
Extending the original 40 to 124 roots adds 454,591 compressed-data bytes and
471,040 binary bytes; the all-260 set would add another 532,368 and 540,672 and
carries unresolved license metadata, so the curated set is the production
default.

Some curated roots contain optional includes for other languages, such as
Markdown fenced-code grammars, Vue preprocessors, C/C++ assembly dialects, and
Ruby template variants. The registry deliberately treats an unavailable
external include as non-fatal. The exact sorted set is checked in at
`OPTIONAL_DEPENDENCIES` and guarded by a test, so adding or removing an
upstream dependency cannot silently change the standalone package. The opt-in
`TestEmbeddedCorpus` gate tokenizes the 42 focused regression files and the
pinned upstream sample of every embedded grammar, 166 files in all, using only
`grammars.Load`, and checks each against the complete source repository. The
same corpus drives the vscode-textmate differential test. CI enables the gate
with `TEXTMATE_GO_REQUIRE_EMBEDDED_CORPUS=1`.

## Regex compatibility golden

The grammar-report test parses and registry-loads all 260 grammars in the
pinned source corpus, then scans all 34,698 regex fields through the same
translation and compilation path used at runtime. The current report has no
grammar load failures, 159 expected unsupported-syntax diagnostics, four
compile failures, and no other diagnostics. The four compile failures are two
patterns each in C# and the unbundled Razor grammar.

Among the bundled roots, C++, C#, Go, Haskell, Kotlin, Markdown, PowerShell,
and Swift contain at least one explicitly degraded unsupported pattern; all
still load, and the embedded differential corpus remains at exact token and
scope parity. The complete per-pattern output lives beside its generator as
[`cmd/grammar-report/testdata/report.golden.md`](../cmd/grammar-report/testdata/report.golden.md).

After an intentional translator or pinned-corpus change, update the golden file
from the repository root and review the resulting diagnostic diff:

```sh
TM_GRAMMARS_DIR=../tm-grammars/packages/tm-grammars \
  go test -count=1 ./cmd/grammar-report \
  -run '^TestGrammarReportGolden$' -update
```

Without `-update`, CI regenerates the report in memory and compares it with the
checked file. The test refuses a dirty checkout or a revision other than the
one recorded by the generated grammar package.
