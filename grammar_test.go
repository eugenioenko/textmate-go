package textmate

import (
	"reflect"
	"testing"

	"github.com/eugenioenko/textmate-go/oniguruma"
)

type testGrammarRepository struct {
	grammars      map[string]*RawGrammar
	injectionMap  map[string][]string
	lookupHistory []string
}

func (r *testGrammarRepository) lookup(scopeName string) *RawGrammar {
	r.lookupHistory = append(r.lookupHistory, scopeName)
	return r.grammars[scopeName]
}

func (r *testGrammarRepository) injections(scopeName string) []string {
	return append([]string(nil), r.injectionMap[scopeName]...)
}

func TestInitGrammarDoesNotMutateInput(t *testing.T) {
	match := "x"
	nestedName := "nested"
	nested := &RawRule{Name: &nestedName, Match: &match}
	raw := &RawGrammar{
		ScopeName:  "source.test",
		Patterns:   []*RawRule{nested},
		Repository: RawRepository{"nested": nested},
		Injections: RawInjections{"L:source.test": nested},
		FileTypes:  []string{"test"},
	}

	initialized := initGrammar(raw, nil)
	if initialized == raw {
		t.Fatal("initGrammar returned its input")
	}
	if raw.Repository["$self"] != nil || raw.Repository["$base"] != nil {
		t.Fatal("initGrammar mutated the input repository")
	}
	if initialized.Repository["nested"] == nested || initialized.Patterns[0] == nested {
		t.Fatal("nested rules were not cloned")
	}
	if initialized.Repository["nested"] != initialized.Patterns[0] ||
		initialized.Injections["L:source.test"] != initialized.Patterns[0] {
		t.Fatal("clone did not preserve shared raw-rule identity")
	}
	if initialized.Repository["$base"] != initialized.Repository["$self"] {
		t.Fatal("root grammar $base should point at $self")
	}
	if got := initialized.Repository["$self"].Patterns[0]; got != initialized.Patterns[0] {
		t.Fatal("$self does not use the cloned top-level patterns")
	}
	if got := *initialized.Repository["$self"].Name; got != raw.ScopeName {
		t.Fatalf("$self name = %q, want %q", got, raw.ScopeName)
	}

	*initialized.Patterns[0].Name = "changed"
	initialized.FileTypes[0] = "changed"
	if *nested.Name != "nested" || raw.FileTypes[0] != "test" {
		t.Fatal("mutating initialized grammar changed caller-owned data")
	}
}

func TestInitGrammarClonesCyclicRawRuleGraph(t *testing.T) {
	cycle := &RawRule{}
	cycle.Patterns = []*RawRule{cycle}
	raw := &RawGrammar{ScopeName: "source.cycle", Patterns: []*RawRule{cycle}}

	initialized := initGrammar(raw, nil)
	if initialized.Patterns[0] == cycle {
		t.Fatal("cyclic rule was not cloned")
	}
	if initialized.Patterns[0].Patterns[0] != initialized.Patterns[0] {
		t.Fatal("cycle was not preserved in clone")
	}
}

func TestExternalGrammarCacheUsesFirstBase(t *testing.T) {
	external := &RawGrammar{ScopeName: "source.external"}
	repository := &testGrammarRepository{
		grammars: map[string]*RawGrammar{"source.external": external},
	}
	g := newGrammar("source.root", &RawGrammar{ScopeName: "source.root"}, repository)
	firstBase := &RawRule{}
	secondBase := &RawRule{}

	first := g.getExternalGrammar("source.external", RawRepository{"$base": firstBase})
	second := g.getExternalGrammar("source.external", RawRepository{"$base": secondBase})
	if first == nil || second != first {
		t.Fatal("external grammar was not cached by scope name")
	}
	if got := first.Repository["$base"]; got != firstBase {
		t.Fatalf("cached $base = %p, want first base %p", got, firstBase)
	}
	if len(repository.lookupHistory) != 1 {
		t.Fatalf("lookup calls = %v, want one", repository.lookupHistory)
	}
	if external.Repository != nil {
		t.Fatal("external source grammar was mutated")
	}
}

func TestRegisterRuleReservesIDBeforeFactory(t *testing.T) {
	g := newGrammar("source.test", &RawGrammar{ScopeName: "source.test"}, nil)
	var allocated ruleID
	registered := g.registerRule(func(id ruleID) rule {
		allocated = id
		if got := g.getRule(id); got != nil {
			t.Fatalf("reserved slot = %T, want nil during construction", got)
		}
		return newCaptureRule(nil, id, nil, nil, 0)
	})
	if allocated != 1 || registered.getID() != 1 || g.getRule(1) != registered {
		t.Fatalf("unexpected registration: allocated=%d result=%v slot=%v", allocated, registered, g.getRule(1))
	}
}

func TestInjectionCollectionHasStablePriorityOrder(t *testing.T) {
	match := "x"
	raw := &RawGrammar{
		ScopeName: "source.test",
		Injections: RawInjections{
			"R:source.test": {Match: &match},
			"source.test":   {Match: &match},
			"L:source.test": {Match: &match},
		},
	}
	external := &RawGrammar{
		ScopeName:         "source.external",
		InjectionSelector: "source.test",
		Patterns:          []*RawRule{{Match: &match}},
	}
	repository := &testGrammarRepository{
		grammars:     map[string]*RawGrammar{"source.external": external},
		injectionMap: map[string][]string{"source.test": {"source.external"}},
	}

	first := newGrammar("source.test", raw, repository).getInjections()
	second := newGrammar("source.test", raw, repository).getInjections()
	priorities := func(injections []injection) []int {
		result := make([]int, len(injections))
		for index, item := range injections {
			result[index] = item.priority
		}
		return result
	}
	if got, want := priorities(first), []int{-1, 0, 0, 1}; !reflect.DeepEqual(got, want) {
		t.Fatalf("priorities = %v, want %v", got, want)
	}
	if got, want := priorities(second), priorities(first); !reflect.DeepEqual(got, want) {
		t.Fatalf("second collection = %v, want %v", got, want)
	}
	for index := range first {
		if first[index].ruleID != second[index].ruleID {
			t.Fatalf("injection %d rule ID = %d, want deterministic %d", index, first[index].ruleID, second[index].ruleID)
		}
	}
}

func TestLineTokenHandlerDropsSyntheticNewline(t *testing.T) {
	scopes := newAttributedScopeRoot("source.test", 0)
	handler := &lineTokenHandler{}
	handler.handle(scopes, 2)
	handler.handle(scopes.pushAttributed("string.test", 0), 3)

	got := handler.result(nil, 2)
	want := []Token{{Start: 0, End: 2, Scopes: []string{"source.test"}}}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("tokens = %#v, want %#v", got, want)
	}
}

func TestLineTokenHandlerProducesEmptyLineToken(t *testing.T) {
	scopes := newAttributedScopeRoot("source.test", 0)
	stack := newStateStack(nil, 1, -1, -1, false, nil, scopes, scopes)
	handler := &lineTokenHandler{}
	handler.handle(scopes, 1)

	got := handler.result(stack, 0)
	want := []Token{{Start: 0, End: 0, Scopes: []string{"source.test"}}}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("tokens = %#v, want %#v", got, want)
	}
}

func TestGrammarTokenizeLineUsesRuneOffsetsAndReusableState(t *testing.T) {
	match := "🙂"
	name := "constant.emoji.test"
	raw := &RawGrammar{
		ScopeName: "source.test",
		Patterns: []*RawRule{{
			Match: &match,
			Name:  &name,
		}},
	}
	g := newGrammar(raw.ScopeName, raw, nil)
	if g.lastRuleID != 0 {
		t.Fatalf("grammar compiled eagerly: last rule ID = %d", g.lastRuleID)
	}

	result := g.TokenizeLine("a🙂b", InitialState)
	want := []Token{
		{Start: 0, End: 1, Scopes: []string{"source.test"}},
		{Start: 1, End: 2, Scopes: []string{"source.test", "constant.emoji.test"}},
		{Start: 2, End: 3, Scopes: []string{"source.test"}},
	}
	if !reflect.DeepEqual(result.Tokens, want) {
		t.Fatalf("tokens = %#v, want %#v", result.Tokens, want)
	}
	if result.RuleStack == nil || result.Stopped {
		t.Fatalf("unexpected result state: stack=%v stopped=%v", result.RuleStack, result.Stopped)
	}
	if g.lastRuleID == 0 {
		t.Fatal("first tokenization did not compile the root")
	}

	next := g.TokenizeLine("plain", result.RuleStack)
	if got, want := next.Tokens, []Token{{Start: 0, End: 5, Scopes: []string{"source.test"}}}; !reflect.DeepEqual(got, want) {
		t.Fatalf("next-line tokens = %#v, want %#v", got, want)
	}
}

func TestGrammarDiagnosticsReportsLazyRegexFailures(t *testing.T) {
	invalid := "("
	raw := &RawGrammar{
		ScopeName: "source.test",
		Patterns: []*RawRule{{
			Match: &invalid,
		}},
	}
	grammar := newGrammar(raw.ScopeName, raw, nil)
	if diagnostics := grammar.Diagnostics(); len(diagnostics) != 0 {
		t.Fatalf("diagnostics before lazy compilation = %+v", diagnostics)
	}

	result := grammar.TokenizeLine("text", nil)
	if len(result.Tokens) != 1 || result.Tokens[0].End != 4 {
		t.Fatalf("degraded grammar tokens = %#v", result.Tokens)
	}
	diagnostics := grammar.Diagnostics()
	if len(diagnostics) != 1 || diagnostics[0].Kind != oniguruma.DiagnosticCompileError || diagnostics[0].Pattern != invalid {
		t.Fatalf("diagnostics = %+v, want one compile error for %q", diagnostics, invalid)
	}
	if again := grammar.Diagnostics(); !reflect.DeepEqual(again, diagnostics) {
		t.Fatalf("second diagnostic snapshot = %+v, want %+v", again, diagnostics)
	}
}

func TestGrammarTokenizeLineUsesNestedRepositoryHeadingBeforeParagraph(t *testing.T) {
	include := func(value string) *RawRule { return &RawRule{Include: &value} }
	value := func(value string) *string { return &value }

	heading := &RawRule{
		Begin:       value(`(?:^|\G)(#{1,6})\s*(?=[\S[^#]])`),
		End:         value(`\s*(#{1,6})?$\n?`),
		Name:        value("markup.heading.markdown"),
		ContentName: value("entity.name.section.markdown"),
		Captures: RawCaptures{
			"1": {Name: value("punctuation.definition.heading.markdown")},
		},
	}
	paragraph := &RawRule{
		Begin: value(`(^|\G)(?=\S)`),
		While: value(`(^|\G)(?!\s*$|#)`),
		Name:  value("meta.paragraph.markdown"),
	}
	block := &RawRule{
		Patterns: []*RawRule{include("#heading"), include("#paragraph")},
		Repository: RawRepository{
			"heading":   heading,
			"paragraph": paragraph,
		},
	}
	raw := &RawGrammar{
		ScopeName: "text.html.markdown",
		Patterns:  []*RawRule{include("#block")},
		Repository: RawRepository{
			"block": block,
		},
	}

	grammar := newGrammar(raw.ScopeName, raw, nil)
	first := grammar.TokenizeLine("This is a paragraph", InitialState)
	if got, want := first.Tokens, []Token{{
		Start: 0, End: 19,
		Scopes: []string{"text.html.markdown", "meta.paragraph.markdown"},
	}}; !reflect.DeepEqual(got, want) {
		t.Fatalf("paragraph tokens = %#v, want %#v", got, want)
	}

	second := grammar.TokenizeLine("## This is great", first.RuleStack)
	want := []Token{
		{Start: 0, End: 2, Scopes: []string{
			"text.html.markdown",
			"markup.heading.markdown",
			"punctuation.definition.heading.markdown",
		}},
		{Start: 2, End: 3, Scopes: []string{
			"text.html.markdown",
			"markup.heading.markdown",
		}},
		{Start: 3, End: 16, Scopes: []string{
			"text.html.markdown",
			"markup.heading.markdown",
			"entity.name.section.markdown",
		}},
	}
	if !reflect.DeepEqual(second.Tokens, want) {
		t.Fatalf("heading tokens = %#v, want %#v", second.Tokens, want)
	}
}
