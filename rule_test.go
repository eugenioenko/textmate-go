package textmate

import (
	"reflect"
	"testing"

	"github.com/eugenioenko/textmate-go/oniguruma"
)

type testRuleFactoryHelper struct {
	rules      map[ruleID]rule
	externals  map[string]*RawGrammar
	lastRuleID ruleID
}

func newTestRuleFactoryHelper() *testRuleFactoryHelper {
	return &testRuleFactoryHelper{
		rules:     make(map[ruleID]rule),
		externals: make(map[string]*RawGrammar),
	}
}

func (h *testRuleFactoryHelper) getRule(id ruleID) rule { return h.rules[id] }

func (h *testRuleFactoryHelper) registerRule(factory func(id ruleID) rule) rule {
	h.lastRuleID++
	id := h.lastRuleID
	result := factory(id)
	h.rules[id] = result
	return result
}

func (h *testRuleFactoryHelper) getExternalGrammar(
	scopeName string,
	_ RawRepository,
) *RawGrammar {
	return h.externals[scopeName]
}

func TestRuleFactoryClassifiesRawRules(t *testing.T) {
	match := RegExpString("match")
	emptyMatch := RegExpString("")
	begin := RegExpString("begin")
	while := RegExpString("while")
	emptyWhile := RegExpString("")

	tests := []struct {
		name string
		raw  *RawRule
		want any
	}{
		{
			name: "truthy match wins",
			raw:  &RawRule{Match: &match, Begin: &begin, While: &while},
			want: (*matchRule)(nil),
		},
		{
			name: "empty match is not truthy",
			raw:  &RawRule{Match: &emptyMatch, Begin: &begin},
			want: (*beginEndRule)(nil),
		},
		{
			name: "absent begin is include only",
			raw:  &RawRule{Patterns: []*RawRule{}},
			want: (*includeOnlyRule)(nil),
		},
		{
			name: "truthy while",
			raw:  &RawRule{Begin: &begin, While: &while},
			want: (*beginWhileRule)(nil),
		},
		{
			name: "empty while is begin end",
			raw:  &RawRule{Begin: &begin, While: &emptyWhile},
			want: (*beginEndRule)(nil),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			helper := newTestRuleFactoryHelper()
			factory := newRuleFactory()
			id := factory.getCompiledRuleID(test.raw, helper, nil)
			if got := reflect.TypeOf(helper.getRule(id)); got != reflect.TypeOf(test.want) {
				t.Fatalf("compiled type = %v, want %v", got, reflect.TypeOf(test.want))
			}
		})
	}
}

func TestRuleFactoryReservesRecursiveRuleID(t *testing.T) {
	include := IncludeString("#loop")
	raw := &RawRule{Patterns: []*RawRule{{Include: &include}}}
	repository := RawRepository{"loop": raw}
	helper := newTestRuleFactoryHelper()

	id := newRuleFactory().getCompiledRuleID(raw, helper, repository)
	compiled, ok := helper.getRule(id).(*includeOnlyRule)
	if !ok {
		t.Fatalf("compiled rule = %T, want *includeOnlyRule", helper.getRule(id))
	}
	if helper.lastRuleID != 1 {
		t.Fatalf("registered %d rules for recursive node, want 1", helper.lastRuleID)
	}
	if !reflect.DeepEqual(compiled.patterns, []ruleID{id}) {
		t.Fatalf("recursive patterns = %v, want [%d]", compiled.patterns, id)
	}
}

func TestRuleFactoryCaptureFallbackDistinguishesNilAndEmpty(t *testing.T) {
	begin := RegExpString("begin")
	while := RegExpString("while")
	captureName := "capture.$1"
	common := RawCaptures{
		"1": {Name: &captureName},
	}
	raw := &RawRule{
		Begin:         &begin,
		While:         &while,
		Captures:      common,
		BeginCaptures: RawCaptures{}, // Present and empty suppresses fallback.
		WhileCaptures: nil,           // Absent falls back to captures.
	}
	help := newTestRuleFactoryHelper()

	id := newRuleFactory().getCompiledRuleID(raw, help, nil)
	compiled := help.getRule(id).(*beginWhileRule)
	if len(compiled.beginCaptures) != 1 || compiled.beginCaptures[0] != nil {
		t.Fatalf("present-empty begin captures = %#v, want one nil slot", compiled.beginCaptures)
	}
	if len(compiled.whileCaptures) != 2 || compiled.whileCaptures[1] == nil {
		t.Fatalf("fallback while captures = %#v, want capture at index 1", compiled.whileCaptures)
	}
}

func TestRuleFactoryCompilesSparseAndRetokenizedCaptures(t *testing.T) {
	match := RegExpString("x")
	name := "entity.$2"
	retokenizedName := "retokenized"
	captures := RawCaptures{
		"2-extra": {
			Name:     &name,
			Patterns: []*RawRule{{Name: &retokenizedName, Match: &match}},
		},
	}
	help := newTestRuleFactoryHelper()
	factory := newRuleFactory()

	compiled := factory.compileCaptures(captures, help, nil)
	if len(compiled) != 3 || compiled[0] != nil || compiled[1] != nil || compiled[2] == nil {
		t.Fatalf("sparse captures = %#v", compiled)
	}
	if compiled[2].retokenizeCapturedWithRuleID == 0 {
		t.Fatal("capture patterns did not create a retokenization rule")
	}
	retokenize := help.getRule(compiled[2].retokenizeCapturedWithRuleID)
	if _, ok := retokenize.(*includeOnlyRule); !ok {
		t.Fatalf("retokenization rule = %T, want *includeOnlyRule", retokenize)
	}
}

func TestRuleFactoryCompilesEveryIncludeKindAndPrunesMissingShells(t *testing.T) {
	matchBase := RegExpString("base")
	matchSelf := RegExpString("self")
	matchLocal := RegExpString("local")
	matchExternalSelf := RegExpString("external-self")
	matchExternalNamed := RegExpString("external-named")
	matchInline := RegExpString("inline")
	base := &RawRule{Match: &matchBase}
	self := &RawRule{Match: &matchSelf}
	local := &RawRule{Match: &matchLocal}
	externalSelf := &RawRule{Match: &matchExternalSelf}
	externalNamed := &RawRule{Match: &matchExternalNamed}
	repository := RawRepository{
		"$base": base,
		"$self": self,
		"local": local,
	}
	help := newTestRuleFactoryHelper()
	help.externals["source.external"] = &RawGrammar{
		ScopeName: "source.external",
		Repository: RawRepository{
			"$self": externalSelf,
			"named": externalNamed,
		},
	}
	factory := newRuleFactory()

	patterns := []*RawRule{
		includePattern("$base"),
		includePattern("$self"),
		includePattern("#local"),
		includePattern("source.external"),
		includePattern("source.external#named"),
		{Name: stringPointer("inline"), Match: &matchInline},
		includePattern("#missing"),
		{Patterns: []*RawRule{includePattern("#also-missing")}},
	}
	result := factory.compilePatterns(patterns, help, repository)
	if !result.hasMissingPatterns {
		t.Fatal("missing includes were not reported")
	}
	if len(result.patterns) != 6 {
		t.Fatalf("compiled %d patterns, want 6: %v", len(result.patterns), result.patterns)
	}

	wantSources := []string{"base", "self", "local", "external-self", "external-named", "inline"}
	for i, id := range result.patterns {
		compiled, ok := help.getRule(id).(*matchRule)
		if !ok || compiled.match.source != wantSources[i] {
			t.Fatalf("pattern %d = %T/%v, want match %q", i, help.getRule(id), compiled, wantSources[i])
		}
	}
}

func TestIncludeOnlyRuleMergesLocalRepository(t *testing.T) {
	outerMatch := RegExpString("outer")
	innerMatch := RegExpString("inner")
	include := IncludeString("#value")
	raw := &RawRule{
		Patterns: []*RawRule{{Include: &include}},
		Repository: RawRepository{
			"value": {Match: &innerMatch},
		},
	}
	repository := RawRepository{"value": {Match: &outerMatch}}
	help := newTestRuleFactoryHelper()

	id := newRuleFactory().getCompiledRuleID(raw, help, repository)
	compiled := help.getRule(id).(*includeOnlyRule)
	if len(compiled.patterns) != 1 {
		t.Fatalf("patterns = %v, want one", compiled.patterns)
	}
	nested := help.getRule(compiled.patterns[0]).(*matchRule)
	if nested.match.source != "inner" {
		t.Fatalf("local repository resolved %q, want inner", nested.match.source)
	}
}

func TestRuleNamesResolveCaptureReferences(t *testing.T) {
	name := ".entity.$1.${2:/upcase}"
	contentName := "meta.${1:/downcase}"
	rule := newCaptureRule(nil, 1, &name, &contentName, 0)
	line := " xx Word "
	captures := []oniguruma.Capture{
		{Start: 0, End: len([]rune(line))},
		{Start: 1, End: 3},
		{Start: 4, End: 8},
	}

	if got := rule.getName(line, captures); got != ".entity.xx.WORD" {
		t.Fatalf("resolved name = %q", got)
	}
	if got := rule.getContentName(line, captures); got != "meta.xx" {
		t.Fatalf("resolved content name = %q", got)
	}
	if got := rule.getName(line, nil); got != name {
		t.Fatalf("name without captures = %q, want raw %q", got, name)
	}
}

func TestBeginEndCompileOrderBackReferencesAndCacheInvalidation(t *testing.T) {
	child := newMatchRule(nil, 2, nil, "child", nil)
	help := newTestRuleFactoryHelper()
	help.rules[2] = child
	patterns := compilePatternsResult{patterns: []ruleID{2}}
	first := newBeginEndRule(nil, 1, nil, nil, "begin", nil, `end-\1`, nil, false, patterns)
	last := newBeginEndRule(nil, 3, nil, nil, "begin", nil, `end-\1`, nil, true, patterns)
	resolvedA := "end-a"
	resolvedB := "end-b"

	firstCompiled := first.compile(help, &resolvedA)
	if !reflect.DeepEqual(firstCompiled.regexps, []string{"end-a", "child"}) ||
		!reflect.DeepEqual(firstCompiled.rules, []ruleID{endRuleID, 2}) {
		t.Fatalf("end-first compile = %v / %v", firstCompiled.regexps, firstCompiled.rules)
	}
	if same := first.compile(help, &resolvedA); same != firstCompiled {
		t.Fatal("unchanged dynamic end did not reuse compiled scanner")
	}
	if changed := first.compile(help, &resolvedB); changed == firstCompiled {
		t.Fatal("changed dynamic end did not invalidate compiled scanner")
	}

	lastCompiled := last.compile(help, &resolvedA)
	if !reflect.DeepEqual(lastCompiled.regexps, []string{"child", "end-a"}) ||
		!reflect.DeepEqual(lastCompiled.rules, []ruleID{2, endRuleID}) {
		t.Fatalf("end-last compile = %v / %v", lastCompiled.regexps, lastCompiled.rules)
	}

	line := "begin-tag"
	captures := []oniguruma.Capture{{Start: 0, End: 9}, {Start: 6, End: 9}}
	if got := first.getEndWithResolvedBackReferences(line, captures); got != "end-tag" {
		t.Fatalf("resolved end = %q", got)
	}
}

func TestBeginWhileUsesSeparateContentAndWhileScanners(t *testing.T) {
	child := newMatchRule(nil, 2, nil, "child", nil)
	help := newTestRuleFactoryHelper()
	help.rules[2] = child
	compiled := newBeginWhileRule(
		nil,
		1,
		nil,
		nil,
		"begin",
		nil,
		`while-\1`,
		nil,
		compilePatternsResult{patterns: []ruleID{2}},
	)

	contentScanner := compiled.compile(help, nil)
	if !reflect.DeepEqual(contentScanner.regexps, []string{"child"}) ||
		!reflect.DeepEqual(contentScanner.rules, []ruleID{2}) {
		t.Fatalf("content scanner = %v / %v", contentScanner.regexps, contentScanner.rules)
	}
	resolved := "while-value"
	whileScanner := compiled.compileWhile(help, &resolved)
	if !reflect.DeepEqual(whileScanner.regexps, []string{"while-value"}) ||
		!reflect.DeepEqual(whileScanner.rules, []ruleID{whileRuleID}) {
		t.Fatalf("while scanner = %v / %v", whileScanner.regexps, whileScanner.rules)
	}
}

func TestRuleDisposeDropsMutableCompiledCaches(t *testing.T) {
	rule := newMatchRule(nil, 1, nil, "x", nil)
	help := newTestRuleFactoryHelper()
	first := rule.compile(help, nil)
	rule.dispose()
	second := rule.compile(help, nil)
	if first == second {
		t.Fatal("dispose retained compiled scanner cache")
	}
}

func includePattern(include string) *RawRule {
	return &RawRule{Include: &include}
}

func stringPointer(value string) *string { return &value }
