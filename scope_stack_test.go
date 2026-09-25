package textmate

import (
	"encoding/json"
	"reflect"
	"sync"
	"testing"
)

func TestGrammarInternsEqualTokenScopeStacks(t *testing.T) {
	match, name := "a", "constant.a.test"
	grammar := newGrammar("source.test", &RawGrammar{
		ScopeName: "source.test",
		Patterns:  []*RawRule{{Match: &match, Name: &name}},
	}, nil)

	first := grammar.TokenizeLine("aba", nil)
	if len(first.Tokens) != 3 {
		t.Fatalf("tokens = %#v, want three", first.Tokens)
	}
	matchedLeft, root, matchedRight := first.Tokens[0], first.Tokens[1], first.Tokens[2]
	if matchedLeft.ScopeStack == nil || root.ScopeStack == nil {
		t.Fatal("token did not carry an interned scope stack")
	}
	if matchedLeft.ScopeStack != matchedRight.ScopeStack || matchedLeft.ScopeStack.ID() != matchedRight.ScopeStack.ID() {
		t.Fatal("equal token scopes did not reuse identity")
	}
	if matchedLeft.ScopeStack == root.ScopeStack || matchedLeft.ScopeStack.ID() == root.ScopeStack.ID() {
		t.Fatal("different token scopes reused identity")
	}
	if &matchedLeft.Scopes[0] != &matchedRight.Scopes[0] {
		t.Fatal("equal legacy Scopes did not share backing storage")
	}

	next := grammar.TokenizeLine("plain", first.RuleStack)
	if next.Tokens[0].ScopeStack != root.ScopeStack {
		t.Fatal("equal scopes on a later line did not reuse identity")
	}
	if &next.Tokens[0].Scopes[0] != &root.Scopes[0] {
		t.Fatal("equal scopes on a later line did not reuse legacy storage")
	}
}

func TestScopeStackInternsReconstructedEqualPathsAndChecksHashCollisions(t *testing.T) {
	first := newAttributedScopeRoot("source.test", 1).pushAttributed("meta.one", 2)
	second := newAttributedScopeRoot("source.test", 99).pushAttributed("meta.one", 100)
	different := newAttributedScopeRoot("source.test", 1).pushAttributed("meta.two", 2)
	var interner scopeStackInterner

	firstStack := interner.intern(first)
	if got := interner.intern(second); got != firstStack {
		t.Fatal("reconstructed equal scope paths did not intern together")
	}

	const forcedHash = 1
	one := interner.internPathWithHash(first.scopePath, forcedHash)
	two := interner.internPathWithHash(different.scopePath, forcedHash)
	if one == two || reflect.DeepEqual(one.names, two.names) {
		t.Fatal("hash collision merged unequal scope paths")
	}
}

func TestScopeStackAccessorsAreImmutableAndConcurrent(t *testing.T) {
	scopes := newAttributedScopeRoot("source.test", 0).pushAttributed("string.test", 0)
	var interner scopeStackInterner
	stack := interner.intern(scopes)
	if stack.ID().IsZero() || stack.Len() != 2 {
		t.Fatalf("ID/Len = (%v, %d), want nonzero/2", stack.ID(), stack.Len())
	}
	if got, ok := stack.At(1); !ok || got != "string.test" {
		t.Fatalf("At(1) = (%q, %v)", got, ok)
	}
	if _, ok := stack.At(-1); ok {
		t.Fatal("At(-1) unexpectedly succeeded")
	}
	names := stack.Names()
	names[0] = "mutated"
	if got, _ := stack.At(0); got != "source.test" {
		t.Fatalf("Names exposed mutable storage: At(0) = %q", got)
	}

	want := []string{"source.test", "string.test"}
	var wait sync.WaitGroup
	for range 16 {
		wait.Add(1)
		go func() {
			defer wait.Done()
			for range 100 {
				var got []string
				stack.Range(func(name string) bool {
					got = append(got, name)
					return true
				})
				if !reflect.DeepEqual(got, want) || stack.ID().IsZero() {
					t.Errorf("concurrent read = %v, ID = %v", got, stack.ID())
					return
				}
			}
		}()
	}
	wait.Wait()
}

func TestFallbackAndEmptyLineTokensCarryInternedScopeStack(t *testing.T) {
	grammar := newGrammar("source.test", &RawGrammar{ScopeName: "source.test"}, nil)
	empty := grammar.TokenizeLine("", nil)
	limited := grammar.TokenizeLineWithOptions("long", empty.RuleStack, TokenizeOptions{MaxLineBytes: 1})
	if len(empty.Tokens) != 1 || len(limited.Tokens) != 1 {
		t.Fatalf("empty/limited tokens = %#v / %#v", empty.Tokens, limited.Tokens)
	}
	if empty.Tokens[0].ScopeStack == nil || limited.Tokens[0].ScopeStack != empty.Tokens[0].ScopeStack {
		t.Fatal("empty and fallback tokens did not carry the canonical root scopes")
	}

	handler := &lineTokenHandler{scopeStacks: &grammar.scopeStacks}
	root := newAttributedScopeRoot("source.test", 0)
	handler.handle(root, 1)
	partial := handler.result(newStateStack(nil, 1, -1, -1, false, nil, root, root), 3)
	if len(partial) != 2 || partial[0].ScopeStack != partial[1].ScopeStack {
		t.Fatalf("partial fallback scopes are not interned: %#v", partial)
	}
}

func TestTokenScopeStackDoesNotChangeJSON(t *testing.T) {
	stack := (&scopeStackInterner{}).intern(newAttributedScopeRoot("source.test", 0))
	token := Token{Start: 1, End: 2, Scopes: stack.compatibilityNames, ScopeStack: stack}
	got, err := json.Marshal(token)
	if err != nil {
		t.Fatal(err)
	}
	const want = `{"Start":1,"End":2,"Scopes":["source.test"]}`
	if string(got) != want {
		t.Fatalf("JSON = %s, want %s", got, want)
	}
}

var benchmarkScopeStack *ScopeStack
var benchmarkScopeNames []string

func TestScopeStackInternHitAllocatesLessThanLegacyScopeCopy(t *testing.T) {
	scopes := newAttributedScopeRoot("source.test", 0).pushAttributed("meta.test", 0)
	var interner scopeStackInterner
	interner.intern(scopes)
	interned := testing.AllocsPerRun(1_000, func() {
		benchmarkScopeStack = interner.intern(scopes)
	})
	legacy := testing.AllocsPerRun(1_000, func() {
		benchmarkScopeNames = append([]string(nil), scopes.scopeNames()...)
	})
	if interned >= legacy {
		t.Fatalf("intern hit allocations = %.2f, legacy copies = %.2f", interned, legacy)
	}
}

func BenchmarkRepeatedEqualTokenScopes(b *testing.B) {
	scopes := newAttributedScopeRoot("source.test", 0).pushAttributed("meta.test", 0)
	b.Run("interned", func(b *testing.B) {
		var interner scopeStackInterner
		interner.intern(scopes)
		b.ReportAllocs()
		for range b.N {
			benchmarkScopeStack = interner.intern(scopes)
		}
	})
	b.Run("legacy-copy", func(b *testing.B) {
		b.ReportAllocs()
		for range b.N {
			benchmarkScopeNames = append([]string(nil), scopes.scopeNames()...)
		}
	})
}
