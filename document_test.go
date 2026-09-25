package textmate

import (
	"errors"
	"fmt"
	"math/rand"
	"reflect"
	"sync"
	"testing"
)

func newDocumentTestGrammar() *Grammar {
	begin, end := `"`, `"`
	stringName := "string.quoted.test"
	contentName := "string.quoted.content.test"
	word, wordName := `\p{L}+`, "word.test"
	return newGrammar("source.test", &RawGrammar{
		ScopeName: "source.test",
		Patterns: []*RawRule{
			{
				Begin:       &begin,
				End:         &end,
				Name:        &stringName,
				ContentName: &contentName,
			},
			{Match: &word, Name: &wordName},
		},
	}, nil)
}

func TestDocumentRandomAccessCarriesStateAndOwnsLines(t *testing.T) {
	grammar := newDocumentTestGrammar()
	document := NewDocument(grammar, DocumentOptions{})
	lines := []string{`"open`, "héllo😀", `"`, "tail"}
	document.SetLines(lines)
	lines[0] = "caller mutation"

	inside, ok := document.Line(1)
	if !ok {
		t.Fatal("random-access line was not found")
	}
	wantScopes := []string{
		"source.test",
		"string.quoted.test",
		"string.quoted.content.test",
	}
	if got := inside.Tokens; len(got) != 1 || got[0].Start != 0 || got[0].End != 6 ||
		!reflect.DeepEqual(got[0].Scopes, wantScopes) {
		t.Fatalf("inside tokens = %#v, want one six-rune string token with scopes %#v", got, wantScopes)
	}
	if got := len(document.states); got != 3 {
		t.Fatalf("materialized states = %d, want start states through line 2", got)
	}

	closing, ok := document.Line(2)
	if !ok || closing.RuleStack == nil || closing.RuleStack.parent != nil {
		t.Fatalf("closing result = (%#v, %v), want root state", closing, ok)
	}
	end, ok := document.StateAt(document.Len())
	if !ok || end == nil || end.parent != nil {
		t.Fatalf("EOF state = (%v, %v), want root", end, ok)
	}
}

func TestDocumentNilGrammarAndBounds(t *testing.T) {
	var nilDocument *Document
	if nilDocument.Len() != 0 {
		t.Fatal("nil document length was nonzero")
	}
	if _, ok := nilDocument.Line(0); ok {
		t.Fatal("nil document returned a line")
	}
	if _, ok := nilDocument.StateAt(0); ok {
		t.Fatal("nil document returned a state")
	}
	if err := nilDocument.ReplaceLines(0, 0, nil); !errors.Is(err, ErrInvalidLineRange) {
		t.Fatalf("nil ReplaceLines error = %v, want ErrInvalidLineRange", err)
	}

	document := NewDocument(nil, DocumentOptions{})
	document.SetLines([]string{"one", "two"})
	for _, index := range []int{-1, 2} {
		if _, ok := document.Line(index); ok {
			t.Fatalf("Line(%d) unexpectedly succeeded", index)
		}
	}
	result, ok := document.Line(1)
	if !ok || len(result.Tokens) != 0 || result.RuleStack != nil || result.Stopped {
		t.Fatalf("nil-grammar result = (%#v, %v), want valid empty result", result, ok)
	}
	for _, index := range []int{0, 2} {
		state, exists := document.StateAt(index)
		if !exists || state != nil {
			t.Fatalf("StateAt(%d) = (%v, %v), want (nil, true)", index, state, exists)
		}
	}
	for _, index := range []int{-1, 3} {
		if _, exists := document.StateAt(index); exists {
			t.Fatalf("StateAt(%d) unexpectedly succeeded", index)
		}
	}
}

func TestDocumentReplaceLinesValidatesAndCopiesReplacement(t *testing.T) {
	document := NewDocument(newDocumentTestGrammar(), DocumentOptions{})
	document.SetLines([]string{"zero", "one", "two"})

	invalid := [][2]int{{-1, 0}, {2, 1}, {0, 4}}
	for _, lineRange := range invalid {
		if err := document.ReplaceLines(lineRange[0], lineRange[1], nil); !errors.Is(err, ErrInvalidLineRange) {
			t.Fatalf("ReplaceLines(%d, %d) error = %v", lineRange[0], lineRange[1], err)
		}
	}
	if document.Len() != 3 {
		t.Fatalf("invalid edits changed length to %d", document.Len())
	}

	replacement := []string{"ONE", "inserted"}
	if err := document.ReplaceLines(1, 2, replacement); err != nil {
		t.Fatal(err)
	}
	replacement[0] = "caller mutation"
	if document.Len() != 4 {
		t.Fatalf("length after replacement = %d, want 4", document.Len())
	}
	result, ok := document.Line(1)
	if !ok || len(result.Tokens) == 0 || result.Tokens[0].End != 3 {
		t.Fatalf("replacement line = (%#v, %v), want owned text ONE", result, ok)
	}
}

func TestDocumentNoOpUpdatePreservesTablesAndCache(t *testing.T) {
	document := NewDocument(newDocumentTestGrammar(), DocumentOptions{})
	lines := []string{"one", "two", "three"}
	document.SetLines(lines)
	if _, ok := document.Line(2); !ok {
		t.Fatal("line 2 was not found")
	}
	states := document.states
	stateSrc := document.stateSrc
	cacheLength := len(document.cache)

	document.SetLines(append([]string(nil), lines...))
	if len(document.states) != len(states) || &document.states[0] != &states[0] ||
		len(document.stateSrc) != len(stateSrc) || &document.stateSrc[0] != &stateSrc[0] {
		t.Fatal("no-op SetLines replaced the materialized state table")
	}
	if len(document.cache) != cacheLength || document.tail != nil {
		t.Fatal("no-op SetLines disturbed the cache or created a tail")
	}
	if err := document.ReplaceLines(1, 2, []string{"two"}); err != nil {
		t.Fatal(err)
	}
	if len(document.states) != len(states) || &document.states[0] != &states[0] {
		t.Fatal("no-op ReplaceLines replaced the materialized state table")
	}
}

func TestDocumentSplicesUnchangedTailAtCanonicalState(t *testing.T) {
	grammar := newDocumentTestGrammar()
	document := NewDocument(grammar, DocumentOptions{})
	lines := make([]string, 100)
	for index := range lines {
		lines[index] = fmt.Sprintf("plain %d", index)
	}
	document.SetLines(lines)
	if _, ok := document.StateAt(len(lines)); !ok {
		t.Fatal("failed to materialize original table")
	}
	oldStates := append([]*StateStack(nil), document.states...)

	changed := append([]string(nil), lines...)
	changed[0] = "changed first line"
	document.SetLines(changed)
	if len(document.states) != 1 || document.tail == nil {
		t.Fatalf("edit retained %d states with tail %v, want initial state and pending tail", len(document.states), document.tail)
	}
	if _, ok := document.Line(0); !ok {
		t.Fatal("changed line was not found")
	}
	if document.tail != nil || len(document.states) != len(lines)+1 {
		t.Fatalf("tail was not spliced: tail=%v states=%d", document.tail, len(document.states))
	}
	for _, index := range []int{1, 50, 100} {
		if document.states[index] != oldStates[index] {
			t.Fatalf("state %d = %p, want reused %p", index, document.states[index], oldStates[index])
		}
	}
}

func TestDocumentDeletionCanSpliceTailImmediately(t *testing.T) {
	document := NewDocument(newDocumentTestGrammar(), DocumentOptions{})
	document.SetLines([]string{"head", "delete me", "tail one", "tail two"})
	if _, ok := document.StateAt(document.Len()); !ok {
		t.Fatal("failed to materialize original table")
	}
	oldTailState := document.states[2]

	if err := document.ReplaceLines(1, 2, nil); err != nil {
		t.Fatal(err)
	}
	if document.tail != nil {
		t.Fatal("equal deletion boundary did not splice the pending tail immediately")
	}
	if len(document.states) != 4 || document.states[1] != oldTailState {
		t.Fatalf("states after deletion = %d, boundary %p; want 4 and %p", len(document.states), document.states[1], oldTailState)
	}
	if got := document.stateSrc; !reflect.DeepEqual(got, []string{"head", "tail one", "tail two"}) {
		t.Fatalf("state sources after deletion = %#v", got)
	}
}

func TestDocumentRepeatedEditsBeforeConvergenceDiscardStaleTail(t *testing.T) {
	grammar := newDocumentTestGrammar()
	document := NewDocument(grammar, DocumentOptions{})
	document.SetLines([]string{`"open`, "inside one", "inside two", `"`, "tail"})
	if _, ok := document.StateAt(document.Len()); !ok {
		t.Fatal("failed to materialize original table")
	}

	if err := document.ReplaceLines(1, 2, []string{`" closes`}); err != nil {
		t.Fatal(err)
	}
	if document.tail == nil {
		t.Fatal("first edit did not retain a tail candidate")
	}
	if err := document.ReplaceLines(2, 3, []string{`"reopens`}); err != nil {
		t.Fatal(err)
	}

	wantLines := []string{`"open`, `" closes`, `"reopens`, `"`, "tail"}
	assertDocumentMatchesSequential(t, document, grammar, wantLines)
}

func TestDocumentEditsMatchSequentialTokenization(t *testing.T) {
	grammar := newDocumentTestGrammar()
	document := NewDocument(grammar, DocumentOptions{})
	tests := []struct {
		name  string
		lines []string
	}{
		{name: "empty", lines: nil},
		{name: "grow", lines: []string{"α", `"open`, "β", `"`, "ω"}},
		{name: "insert", lines: []string{"header", "α", `"open`, "β", `"`, "ω"}},
		{name: "delete", lines: []string{"header", "α", "β", `"`, "ω"}},
		{name: "shrink", lines: []string{"header", "α"}},
		{name: "whole rewrite", lines: []string{`"new`, "😀 content", `"`}},
		{name: "empty again", lines: []string{}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			document.SetLines(test.lines)
			assertDocumentMatchesSequential(t, document, grammar, test.lines)
		})
	}
}

func TestDocumentRandomEditsMatchSequentialTokenization(t *testing.T) {
	random := rand.New(rand.NewSource(42))
	grammar := newDocumentTestGrammar()
	document := NewDocument(grammar, DocumentOptions{CacheCapacity: 32})
	lines := []string{"alpha", `"open`, "inside", `"`, "omega"}
	document.SetLines(lines)
	choices := []string{"plain", `"open`, `close"`, "éclair", "😀", "", "tail"}

	for iteration := range 200 {
		start := random.Intn(len(lines) + 1)
		end := start + random.Intn(len(lines)-start+1)
		replacement := make([]string, random.Intn(4))
		for index := range replacement {
			replacement[index] = choices[random.Intn(len(choices))]
		}
		if err := document.ReplaceLines(start, end, replacement); err != nil {
			t.Fatalf("edit %d [%d,%d): %v", iteration, start, end, err)
		}
		next := make([]string, 0, len(lines)-(end-start)+len(replacement))
		next = append(next, lines[:start]...)
		next = append(next, replacement...)
		next = append(next, lines[end:]...)
		lines = next

		assertDocumentMatchesSequential(t, document, grammar, lines)
	}
}

func TestDocumentCacheIsBoundedAndCanBeDisabled(t *testing.T) {
	grammar := newDocumentTestGrammar()
	document := NewDocument(grammar, DocumentOptions{CacheCapacity: 2})
	document.SetLines([]string{"one", "two", "three"})
	for index := range 3 {
		if _, ok := document.Line(index); !ok {
			t.Fatalf("line %d was not found", index)
		}
	}
	if len(document.cache) != 2 || len(document.cacheOrder) != 2 {
		t.Fatalf("cache sizes = (%d, %d), want (2, 2)", len(document.cache), len(document.cacheOrder))
	}
	if _, exists := document.cache[documentCacheKey{line: "one", state: nil}]; exists {
		t.Fatal("oldest cache entry was not evicted")
	}

	disabled := NewDocument(grammar, DocumentOptions{CacheCapacity: -1})
	disabled.SetLines([]string{"one"})
	if _, ok := disabled.Line(0); !ok {
		t.Fatal("disabled-cache document did not return its line")
	}
	if disabled.cache != nil || disabled.cacheOrder != nil {
		t.Fatal("disabled cache retained a result")
	}
}

var documentAllocationResult LineResult

func TestDocumentWarmCacheHitDoesNotAllocate(t *testing.T) {
	document := NewDocument(newDocumentTestGrammar(), DocumentOptions{})
	document.SetLines([]string{"cached"})
	if _, ok := document.Line(0); !ok {
		t.Fatal("line was not found")
	}
	allocations := testing.AllocsPerRun(1_000, func() {
		var ok bool
		documentAllocationResult, ok = document.Line(0)
		if !ok {
			panic("cached line disappeared")
		}
	})
	if allocations != 0 {
		t.Fatalf("warm cache hit allocations = %.2f, want 0", allocations)
	}
}

func TestDocumentAppliesTokenizeOptionsToEveryLine(t *testing.T) {
	document := NewDocument(newDocumentTestGrammar(), DocumentOptions{
		TokenizeOptions: TokenizeOptions{MaxLineRunes: 5},
	})
	document.SetLines([]string{`"open`, "😀😀😀😀😀😀", `"`})

	stopped, ok := document.Line(1)
	if !ok || !stopped.Stopped || stopped.StoppedReason != StopReasonLineLimit || stopped.StoppedAt != 0 {
		t.Fatalf("stopped line = (%#v, %v), want line limit", stopped, ok)
	}
	if len(stopped.Tokens) != 1 || stopped.Tokens[0].End != 6 {
		t.Fatalf("stopped Unicode tokens = %#v, want full six-rune fallback", stopped.Tokens)
	}
	closed, ok := document.Line(2)
	if !ok || closed.RuleStack == nil || closed.RuleStack.parent != nil {
		t.Fatalf("line after stopped line = (%#v, %v), want carried string to close", closed, ok)
	}
}

func TestDocumentDoesNotCacheTimeLimitStops(t *testing.T) {
	document := NewDocument(newDocumentTestGrammar(), DocumentOptions{
		TokenizeOptions: TokenizeOptions{TimeLimit: 1},
	})
	document.SetLines([]string{"one", "two", "three"})

	for attempt := range 2 {
		result, ok := document.Line(0)
		if !ok || !result.Stopped || result.StoppedReason != StopReasonTimeLimit {
			t.Fatalf("attempt %d = (%#v, %v), want transient time stop", attempt, result, ok)
		}
		if len(document.cache) != 0 {
			t.Fatalf("attempt %d cached a time-limited result", attempt)
		}
	}

	// A consistently tiny budget must still make finite progress through a
	// random-access state request rather than repeatedly retrying one line.
	state, ok := document.StateAt(document.Len())
	if !ok || state == nil || len(document.states) != document.Len()+1 {
		t.Fatalf("EOF after persistent stops = (%v, %v), states=%d", state, ok, len(document.states))
	}
}

func TestDocumentSuccessfulRetryRepairsMaterializedTail(t *testing.T) {
	document := NewDocument(newDocumentTestGrammar(), DocumentOptions{
		TokenizeOptions: TokenizeOptions{TimeLimit: 1},
	})
	lines := []string{`"open`, "inside", `"`, "tail"}
	document.SetLines(lines)
	if _, ok := document.StateAt(document.Len()); !ok {
		t.Fatal("failed to materialize time-stopped document")
	}
	if len(document.states) != len(lines)+1 {
		t.Fatalf("materialized state count = %d", len(document.states))
	}
	stoppedState := document.states[1]

	// Options are fixed in the public API. Changing the internal copy models a
	// later call succeeding after transient pressure without a timing-sensitive
	// test.
	document.options.TimeLimit = 0
	retried, ok := document.Line(0)
	if !ok || retried.Stopped || retried.RuleStack == stoppedState {
		t.Fatalf("retry = (%#v, %v), want a newly opened string state", retried, ok)
	}
	if len(document.states) != 2 {
		t.Fatalf("successful retry retained %d stale downstream states, want 2", len(document.states))
	}
	if document.states[1] != retried.RuleStack {
		t.Fatal("successful retry did not install its outgoing state")
	}

	assertDocumentMatchesSequential(t, document, document.grammar, lines)
}

func TestDocumentCachesDeterministicLineLimitStops(t *testing.T) {
	document := NewDocument(newDocumentTestGrammar(), DocumentOptions{
		TokenizeOptions: TokenizeOptions{MaxLineRunes: 2},
	})
	document.SetLines([]string{"too long"})
	first, ok := document.Line(0)
	if !ok || !first.Stopped || first.StoppedReason != StopReasonLineLimit {
		t.Fatalf("first result = (%#v, %v), want line stop", first, ok)
	}
	if len(document.cache) != 1 {
		t.Fatalf("line-limit cache size = %d, want 1", len(document.cache))
	}
	second, ok := document.Line(0)
	if !ok || len(second.Tokens) == 0 || &second.Tokens[0] != &first.Tokens[0] {
		t.Fatal("deterministic line stop did not reuse its cached result")
	}
}

func TestDocumentInvalidateFromRetriesTail(t *testing.T) {
	document := NewDocument(newDocumentTestGrammar(), DocumentOptions{})
	lines := []string{`"open`, "inside", `"`, "tail"}
	document.SetLines(lines)
	if _, ok := document.StateAt(document.Len()); !ok {
		t.Fatal("failed to materialize document")
	}
	if len(document.cache) == 0 {
		t.Fatal("expected populated result cache")
	}

	for _, index := range []int{-1, document.Len() + 1} {
		if err := document.InvalidateFrom(index); !errors.Is(err, ErrInvalidLineRange) {
			t.Fatalf("InvalidateFrom(%d) error = %v", index, err)
		}
	}
	if err := document.InvalidateFrom(1); err != nil {
		t.Fatal(err)
	}
	if len(document.states) != 2 || len(document.stateSrc) != 1 {
		t.Fatalf("invalidated tables = (%d states, %d sources), want (2, 1)", len(document.states), len(document.stateSrc))
	}
	if len(document.cache) != 0 || len(document.cacheOrder) != 0 || document.tail != nil {
		t.Fatal("invalidation retained cache entries or a pending tail")
	}
	assertDocumentMatchesSequential(t, document, document.grammar, lines)

	if err := document.InvalidateFrom(document.Len()); err != nil {
		t.Fatalf("EOF invalidation: %v", err)
	}
	var nilDocument *Document
	if err := nilDocument.InvalidateFrom(0); !errors.Is(err, ErrInvalidLineRange) {
		t.Fatalf("nil invalidation error = %v", err)
	}
}

func TestDocumentConcurrentReadsAndUpdates(t *testing.T) {
	grammar := newDocumentTestGrammar()
	document := NewDocument(grammar, DocumentOptions{CacheCapacity: 16})
	document.SetLines([]string{"one", `"open`, "inside", `"`, "tail"})

	var wait sync.WaitGroup
	for worker := range 8 {
		wait.Add(1)
		go func() {
			defer wait.Done()
			for iteration := range 100 {
				length := document.Len()
				if length > 0 {
					_, _ = document.Line((worker + iteration) % length)
				}
				_, _ = document.StateAt(0)
			}
		}()
	}
	wait.Add(1)
	go func() {
		defer wait.Done()
		for iteration := range 100 {
			if iteration%2 == 0 {
				document.SetLines([]string{"one", `"open`, "inside", `"`, "tail"})
			} else {
				document.SetLines([]string{"one", "plain", "tail"})
			}
		}
	}()
	wait.Wait()

	document.SetLines([]string{"final", "😀"})
	assertDocumentMatchesSequential(t, document, grammar, []string{"final", "😀"})
}

func assertDocumentMatchesSequential(
	t *testing.T,
	document *Document,
	grammar *Grammar,
	lines []string,
) {
	t.Helper()
	if got := document.Len(); got != len(lines) {
		t.Fatalf("document length = %d, want %d", got, len(lines))
	}
	var state *StateStack
	for index, line := range lines {
		want := grammar.TokenizeLine(line, state)
		got, ok := document.Line(index)
		if !ok {
			t.Fatalf("Line(%d) was not found", index)
		}
		if got.Stopped != want.Stopped || got.StoppedReason != want.StoppedReason ||
			got.StoppedAt != want.StoppedAt || !tokenValuesEqual(got.Tokens, want.Tokens) ||
			!got.RuleStack.Equal(want.RuleStack) {
			t.Fatalf("Line(%d) = %#v, want sequential %#v", index, got, want)
		}
		state = want.RuleStack
	}
	gotEOF, ok := document.StateAt(len(lines))
	if !ok || !gotEOF.Equal(state) {
		t.Fatalf("EOF state = (%v, %v), want %v", gotEOF, ok, state)
	}
}
