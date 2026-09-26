# Usage guide

This guide covers the APIs beyond basic line tokenization. See the project
[`README`](../README.md) for installation and a minimal example.

## State stacks

A state stack is immutable and reusable, but belongs to the grammar that
created it. Pass `nil` or `textmate.InitialState` for the first line.

Equal states returned by the same grammar have the same pointer while both are
live. Consumers can use `*StateStack` as a short-lived cache key and stop
incremental highlighting with a pointer comparison. The grammar holds only a
weak reference to canonical states, so use `StateStack.Equal` when an earlier
equal pointer is no longer retained. Pointer identity is not shared between
grammars, and states must not be carried from one grammar into another.

## Editable documents

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

line, ok := document.Line(1)
if ok {
	for _, token := range line.Tokens {
		_ = token
	}
}

// Replace the half-open range [1, 2). Common-prefix states are retained.
if err := document.ReplaceLines(1, 2, []string{`const message = "goodbye"`}); err != nil {
	log.Fatal(err)
}

// Discard state from a line onward before retrying transiently stopped work.
if err := document.InvalidateFrom(1); err != nil {
	log.Fatal(err)
}
```

`SetLines` and `ReplaceLines` shallow-copy their input slices. `Line`,
`StateAt`, `Len`, and updates are safe to call concurrently. `StateAt` also
accepts `document.Len()` to obtain the state after the final line.

Line results are immutable, library-owned views and remain valid after later
calls or cache eviction. This makes warmed cache hits allocation-free. Copy
`LineResult.Tokens` before changing token fields. `Token.Scopes` and
`Token.ScopeStack` remain read-only.

Cache keys combine line text with the canonical start-state pointer, so
identical lines in identical states reuse a result after edits. `CacheCapacity`
bounds retained results. Zero selects the 20,000-entry default and a negative
value disables the cache.

Line-length stops are deterministic and cached normally. Time-limit stops are
transient. `Document` materializes their outgoing state, but does not cache the
line result. Calling `Line` again retries that line. If a successful retry
changes its outgoing state, the document truncates stale downstream states and
reconciles the unchanged tail again.

## Interned token scopes

Every token carries an immutable `*ScopeStack`. Equal scope sequences produced
by one grammar reuse the same pointer and `ScopeStackID`, even when the
tokenizer reconstructed its internal state. Consumers can cache resolved styles
by `token.ScopeStack` or `token.ScopeStack.ID()` instead of hashing strings for
every token.

The handle remains valid independently of the grammar and its accessors are
safe for concurrent reads. `Names` returns a defensive copy. `Len`, `At`, and
`Range` inspect names without allocating.

The legacy `Token.Scopes` field remains available and preserves existing JSON
and conformance output. Equal stacks share its backing array, so it must be
treated as read-only. Copy the slice before modifying it.

## Token categories

For consumers that need semantic classes instead of full TextMate theme
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
variables, markup, diff, punctuation, and invalid scopes.

The innermost recognized scope wins, while reset scopes such as
`meta.embedded` prevent an outer string from coloring embedded source. This is
a stable, theme-neutral heuristic, not a replacement for TextMate selector and
theme resolution.

## Tokenization limits

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
not parsed. It receives one fallback token with the incoming scopes, its reset
incoming state is returned unchanged, and `StoppedAt` is zero. Size caps are
checked before allocating the newline-appended regexp input and rune buffer.

`TimeLimit` is a soft overall budget. It includes lock wait, lazy root
compilation, and begin/while setup, and is checked between regexp searches. The
library cannot interrupt an in-flight regexp search, so a call can return after
the requested duration. A partial result reports the next unparsed rune in
`StoppedAt`, carries the state reached at the last completed match boundary,
and gives the unparsed tail a fallback token so returned tokens still cover the
complete line.
