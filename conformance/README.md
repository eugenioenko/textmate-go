# Conformance harness

This package runs upstream `vscode-textmate` fixtures without copying or editing
them. The default checkout is `../../vscode-textmate`; override it with an
absolute or relative `VSCODE_TEXTMATE_DIR` path. The supported reference commit
is `fbe49961ab8077e587fdf5282019655ae69e5f9e`.
The harness pins `vscode-oniguruma` 1.7.0, matching that checkout's numeric
find-option API.

```sh
pnpm install
pnpm test:reference
pnpm test:reference:tokenization
pnpm test:go
pnpm test:go:theme-scopes
pnpm test:go:themes
pnpm test:differential
pnpm test:differential:fuzz
pnpm test:differential:ci
```

`TEXTMATE_BACKEND=reference` routes the upstream tests through the adapter and
then into the real TypeScript engine; `tokenizeLine2` is still reconstructed by
the adapter rather than delegated. This must be green before Go results are
trusted. `TEXTMATE_BACKEND=go` starts one long-running JSON-lines CLI process.
By default that process is `go run ./cmd/textmate-cli`, with the repository root
as its working directory. Set `TEXTMATE_GO_CLI` to a prebuilt executable, and
optionally set `TEXTMATE_GO_CLI_ARGS` to a JSON array of arguments.

Because vscode-textmate's tokenization API is synchronous, the Node 24 harness
uses synchronous reads and writes on Node's internal child-pipe file descriptor
(`_handle.fd`), retrying nonblocking `EAGAIN` results. This intentionally small
development-only bridge is Node-specific; a worker-thread bridge is the
portable fallback if a future Node release removes that descriptor.

All registries in a Vitest worker share that single CLI process. The harness
teardown closes it after the worker's test file completes; disposing a registry
only frees its `rN`/`gN`/`sN` handles and does not stop the shared process.

The tests' `../main`, `./onigLibs`, and `../diffStateStacks` imports are resolved
to harness adapters by the Vitest config. Upstream test and fixture files remain
unchanged. State-stack diffing intentionally passes opaque Go state handles
through; it does not test VS Code's private `StateStackImpl` representation.

The CLI adapter gathers external `include` dependencies through the supplied
`RegistryOptions.loadGrammar`, sends raw grammars in `newRegistry`, and then
uses opaque `rN`, `gN`, and `sN` handles. CLI replies shaped as
`{error:string, code?:string}` become rejected `GoProcessError`s. Protocol token
offsets are UTF-16 indices, matching JavaScript's `substring` behavior.

`tokenizeLine2` is reconstructed on the JavaScript side from Go's scope stacks,
the upstream theme matcher, embedded-language matching, and packed metadata.
Adjacent runs with equal metadata are merged.

Go mode's normal command is restricted to tokenization because
`themes.test.ts` writes `.result` files on assertion failures. The dedicated
theme command verifies the exact clean pinned checkout, copies only `src/` and
`test-cases/` into a marked temporary directory, and runs a static list of the
72 golden fixture tests there. Direct or unmarked Go theme runs are rejected;
the runner enforces an outer process-group timeout and validates the exact
temporary root before removing it.

Current conformance gates at the pinned checkout: full reference mode passes
all 222 upstream tests (228 total with two client and four infrastructure
tests); the reference-tokenization and Go-tokenization commands each pass 101
tests (95 upstream fixtures, two client tests, and four infrastructure tests),
with three opt-in tests skipped. Go also passes all 72 isolated theme goldens
under the checkout's 14 themes. `src/themeScopeDiff.test.ts` is an
opt-in, non-mutating source/scope differential exposed as
`pnpm test:go:theme-scopes`; any difference fails the command.

The GitHub Actions merge gate runs typechecking, reference and Go tokenization,
all 127 upstream theme tests through the isolated reference adapter, the
72-file scope differential, the real-file corpus, the bounded 48-case fuzzer,
and all 72 isolated Go golden-file cases. Of the 127 reference theme tests, 53
exercise vscode-textmate's deliberately unported Theme parser/matcher and two
are unbounded tokenization-time smoke tests; their reference-adapter path
remains enforced by `pnpm test:reference:themes`. The Go golden suite takes
roughly 82 seconds and runs on every pull request, `main` push, version tag, and
manual workflow dispatch so a green CI result covers the complete required
suite.

The differential commands use the real vscode-textmate/vscode-oniguruma engine
as an oracle and the pinned sibling `tm-grammars` checkout as the grammar
source. Override that source with `TM_GRAMMARS_DIR`. The corpus command carries
each backend's state line by line and compares token text plus the full scope
stack. Its 42 focused regression files are original MIT-licensed test material;
40 additional real language samples are read from the exact clean pinned
`textmate-grammars-themes` checkout under its root MIT license. The manifest
records source, revision, license, grammar, and feature coverage. Mismatch lines
are one-based; columns are zero-based UTF-16 offsets, matching the upstream JS
API.

The bounded fuzz command defaults to 48 deterministic cases covering line
splices, edits, and arbitrary valid UTF-8. Set `TEXTMATE_DIFF_SEED` and
`TEXTMATE_DIFF_FUZZ_CASES` to replay or extend a run. A mismatch is delta-debugged
within a fixed minimization budget, then reported with its seed, iteration,
grammar, first differing line/column/token, and minimized input. The `:ci`
command runs both the corpus and the explicitly bounded 48-case fuzz gate. The
82-entry corpus currently reaches exact text-and-scope parity for all files,
covering every embedded root with both a focused regression and a real upstream
sample, plus YAML and TOML; the extended Phase 4 gate also passed 500
deterministic fuzz cases.
