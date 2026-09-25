package textmate

import (
	"sync"
	"testing"
	"time"
	"weak"
)

func canonicalStateTestGrammar() *Grammar {
	begin, end := `"`, `"`
	name, contentName := "string.quoted.test", "string.quoted.content.test"
	return newGrammar("source.test", &RawGrammar{
		ScopeName: "source.test",
		Patterns: []*RawRule{{
			Begin:       &begin,
			End:         &end,
			Name:        &name,
			ContentName: &contentName,
		}},
	}, nil)
}

func TestGrammarCanonicalizesRepeatedEndStates(t *testing.T) {
	grammar := canonicalStateTestGrammar()
	first := grammar.TokenizeLine("plain", nil).RuleStack
	second := grammar.TokenizeLine("different", nil).RuleStack
	third := grammar.TokenizeLine("plain", InitialState).RuleStack

	if first == nil || first != second || first != third {
		t.Fatalf("equal root states were not canonical: %p %p %p", first, second, third)
	}
	if InitialState == first {
		t.Fatal("grammar root state unexpectedly reused InitialState")
	}
}

func TestGrammarCanonicalStatesConvergeAcrossChangedLines(t *testing.T) {
	grammar := canonicalStateTestGrammar()
	root := grammar.TokenizeLine("plain", nil).RuleStack

	before := grammar.TokenizeLine(`"before`, root).RuleStack
	after := grammar.TokenizeLine(`"after`, root).RuleStack
	if before == nil || before != after {
		t.Fatalf("changed lines with equal end states did not converge: %p != %p", before, after)
	}

	closedBefore := grammar.TokenizeLine(`"`, before).RuleStack
	closedAfter := grammar.TokenizeLine(`"`, after).RuleStack
	if closedBefore != root || closedAfter != root {
		t.Fatalf("closed states did not converge to root: %p %p, want %p", closedBefore, closedAfter, root)
	}
}

func TestGrammarCanonicalizesCarriedMultilineState(t *testing.T) {
	grammar := canonicalStateTestGrammar()
	opened := grammar.TokenizeLine(`"open`, nil).RuleStack
	first := grammar.TokenizeLine("continued", opened).RuleStack
	second := grammar.TokenizeLine("changed continuation", opened).RuleStack

	if opened == nil || opened != first || opened != second {
		t.Fatalf("equal carried states were not canonical: opened=%p first=%p second=%p", opened, first, second)
	}
}

func TestGrammarCanonicalizesStoppedStates(t *testing.T) {
	grammar := canonicalStateTestGrammar()
	opened := grammar.TokenizeLine(`"open`, nil).RuleStack

	lineStop := grammar.TokenizeLineWithOptions(
		"too long",
		opened,
		TokenizeOptions{MaxLineBytes: 1},
	)
	if !lineStop.Stopped || lineStop.StoppedReason != StopReasonLineLimit {
		t.Fatalf("line stop = (%v, %v), want line limit", lineStop.Stopped, lineStop.StoppedReason)
	}
	if lineStop.RuleStack != opened {
		t.Fatalf("line-limit state = %p, want canonical %p", lineStop.RuleStack, opened)
	}

	timeStop := grammar.TokenizeLineWithOptions(
		"continued",
		opened,
		TokenizeOptions{TimeLimit: time.Nanosecond},
	)
	if !timeStop.Stopped || timeStop.StoppedReason != StopReasonTimeLimit {
		t.Fatalf("time stop = (%v, %v), want time limit", timeStop.Stopped, timeStop.StoppedReason)
	}
	if timeStop.RuleStack != opened {
		t.Fatalf("time-limit state = %p, want canonical %p", timeStop.RuleStack, opened)
	}

	firstLineStop := grammar.TokenizeLineWithOptions(
		"too long",
		nil,
		TokenizeOptions{MaxLineBytes: 1},
	).RuleStack
	secondLineStop := grammar.TokenizeLineWithOptions(
		"also too long",
		InitialState,
		TokenizeOptions{MaxLineBytes: 1},
	).RuleStack
	if firstLineStop == nil || firstLineStop != secondLineStop {
		t.Fatalf("first-line stopped states were not canonical: %p != %p", firstLineStop, secondLineStop)
	}
}

func TestCanonicalStatePointerWorksAsCacheKey(t *testing.T) {
	grammar := canonicalStateTestGrammar()
	first := grammar.TokenizeLine(`"one`, nil).RuleStack
	cache := map[*StateStack]string{first: "style"}

	recomputed := grammar.TokenizeLine(`"two`, nil).RuleStack
	if got, ok := cache[recomputed]; !ok || got != "style" {
		t.Fatalf("pointer-key cache miss for recomputed equal state %p (original %p)", recomputed, first)
	}
}

func TestStateStackInternerChecksHashCollisions(t *testing.T) {
	scopes := newAttributedScopeRoot("source.test", 0)
	first := newStateStack(nil, 1, -1, -1, false, nil, scopes, scopes)
	second := newStateStack(nil, 2, -1, -1, false, nil, scopes, scopes)
	firstCopy := newStateStack(nil, 1, 99, 100, true, nil, scopes, scopes)

	const collisionHash = 42
	first.equalityHash = collisionHash
	second.equalityHash = collisionHash
	firstCopy.equalityHash = collisionHash

	var interner stateStackInterner
	canonicalFirst := interner.intern(first)
	canonicalSecond := interner.intern(second)
	if canonicalFirst == canonicalSecond {
		t.Fatal("unequal states sharing a hash were merged")
	}
	if got := interner.intern(firstCopy); got != canonicalFirst {
		t.Fatalf("equal state in collision bucket returned %p, want %p", got, canonicalFirst)
	}
	if got := interner.intern(second); got != canonicalSecond {
		t.Fatalf("second collision entry was lost: got %p, want %p", got, canonicalSecond)
	}
	interner.byHash[collisionHash] = append(
		interner.byHash[collisionHash],
		weak.Pointer[StateStack]{},
	)
	if got := interner.intern(firstCopy); got != canonicalFirst {
		t.Fatalf("canonical state changed while compacting dead entry: %p", got)
	}
	if got := len(interner.byHash[collisionHash]); got != 2 {
		t.Fatalf("collision bucket retained dead entry: length %d, want 2", got)
	}

	const staleHash = 99
	interner.byHash[staleHash] = []weak.Pointer[StateStack]{{}}
	interner.calls = stateStackInternerSweepInterval - 1
	interner.intern(firstCopy)
	if _, exists := interner.byHash[staleHash]; exists {
		t.Fatal("periodic sweep retained a dead bucket")
	}
	if interner.intern(nil) != nil || interner.intern(InitialState) != InitialState {
		t.Fatal("interner did not preserve nil or InitialState")
	}
}

func TestCanonicalStateIdentityIsPerGrammar(t *testing.T) {
	first := canonicalStateTestGrammar().TokenizeLine("plain", nil).RuleStack
	second := canonicalStateTestGrammar().TokenizeLine("plain", nil).RuleStack
	if first == second {
		t.Fatal("separate grammars unexpectedly shared a canonical state pointer")
	}
}

func TestGrammarCanonicalStateConcurrentTokenization(t *testing.T) {
	grammar := canonicalStateTestGrammar()
	want := grammar.TokenizeLine(`"held`, nil).RuleStack

	const goroutines = 32
	const iterations = 50
	results := make(chan *StateStack, goroutines*iterations)
	var wait sync.WaitGroup
	wait.Add(goroutines)
	for index := 0; index < goroutines; index++ {
		go func() {
			defer wait.Done()
			for iteration := 0; iteration < iterations; iteration++ {
				results <- grammar.TokenizeLine(`"recomputed`, nil).RuleStack
			}
		}()
	}
	wait.Wait()
	close(results)

	for got := range results {
		if got != want {
			t.Fatalf("concurrent canonical state = %p, want %p", got, want)
		}
	}
}
