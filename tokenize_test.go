package textmate

import (
	"reflect"
	"testing"
	"time"

	"github.com/eugenioenko/textmate-go/oniguruma"
)

type tokenizerTestGrammar struct {
	rules      map[ruleID]rule
	injections []injection
	attributes map[string]uint32
}

func (g *tokenizerTestGrammar) getRule(id ruleID) rule { return g.rules[id] }

func (g *tokenizerTestGrammar) getInjections() []injection { return g.injections }

func (g *tokenizerTestGrammar) scopeAttributes(scopeName string) uint32 {
	return g.attributes[scopeName]
}

type recordedToken struct {
	start  int
	end    int
	scopes []string
}

type recordingTokenHandler struct {
	lastEnd int
	tokens  []recordedToken
}

func (h *recordingTokenHandler) handle(scopes *attributedScopeStack, endIndex int) {
	if h.lastEnd >= endIndex {
		return
	}
	h.tokens = append(h.tokens, recordedToken{
		start:  h.lastEnd,
		end:    endIndex,
		scopes: append([]string(nil), scopes.scopeNames()...),
	})
	h.lastEnd = endIndex
}

func tokenizerTestRootStack(id ruleID) *StateStack {
	scopes := newAttributedScopeRoot("source.test", 0)
	return newStateStack(nil, id, -1, -1, false, nil, scopes, scopes)
}

func tokenizerTestString(value string) *string { return &value }

func TestTokenizeStringUsesRuneOffsetsAndSyntheticNewline(t *testing.T) {
	rootID, matchID := ruleID(1), ruleID(2)
	match := newMatchRule(nil, matchID, tokenizerTestString("constant.numeric"), `(\d)(\d)`, []*captureRule{
		nil,
		newCaptureRule(nil, 3, tokenizerTestString("digit.first"), nil, 0),
		newCaptureRule(nil, 4, tokenizerTestString("digit.second"), nil, 0),
	})
	root := newIncludeOnlyRule(nil, rootID, nil, nil, compilePatternsResult{patterns: []ruleID{matchID}})
	grammar := &tokenizerTestGrammar{rules: map[ruleID]rule{rootID: root, matchID: match}}
	handler := &recordingTokenHandler{}

	// Grammar.TokenizeLine supplies this trailing newline. Its rune is visible
	// to the core handler and removed by the public token adapter.
	result := tokenizeString(
		grammar,
		oniguruma.NewString("é12\n"),
		true,
		0,
		tokenizerTestRootStack(rootID),
		handler,
		true,
		0,
	)

	if result.stoppedEarly {
		t.Fatal("tokenization stopped early")
	}
	want := []recordedToken{
		{start: 0, end: 1, scopes: []string{"source.test"}},
		{start: 1, end: 2, scopes: []string{"source.test", "constant.numeric", "digit.first"}},
		{start: 2, end: 3, scopes: []string{"source.test", "constant.numeric", "digit.second"}},
		{start: 3, end: 4, scopes: []string{"source.test"}},
	}
	if !reflect.DeepEqual(handler.tokens, want) {
		t.Fatalf("tokens mismatch:\n got %#v\nwant %#v", handler.tokens, want)
	}
	if result.stack == nil {
		t.Fatal("ended with a nil stack")
	}
	if result.stack.ruleID != rootID {
		t.Fatalf("ended on rule %d, want root %d", result.stack.ruleID, rootID)
	}
}

func TestTokenizeStringCarriesBeginEndStateAndResolvesBackReference(t *testing.T) {
	rootID, tagID, textID := ruleID(1), ruleID(2), ruleID(3)
	text := newMatchRule(nil, textID, tokenizerTestString("text.tag"), `x`, nil)
	tag := newBeginEndRule(
		nil,
		tagID,
		tokenizerTestString("meta.tag"),
		tokenizerTestString("meta.tag.content"),
		`<(\w+)>`,
		nil,
		`</\1>`,
		nil,
		false,
		compilePatternsResult{patterns: []ruleID{textID}},
	)
	root := newIncludeOnlyRule(nil, rootID, nil, nil, compilePatternsResult{patterns: []ruleID{tagID}})
	grammar := &tokenizerTestGrammar{rules: map[ruleID]rule{
		rootID: root,
		tagID:  tag,
		textID: text,
	}}

	firstHandler := &recordingTokenHandler{}
	first := tokenizeString(
		grammar,
		oniguruma.NewString("<tag>x\n"),
		true,
		0,
		tokenizerTestRootStack(rootID),
		firstHandler,
		true,
		0,
	)
	if first.stack.ruleID != tagID || first.stack.parent == nil || first.stack.parent.ruleID != rootID {
		t.Fatalf("first line stack = %v, want open tag over root", first.stack)
	}
	if first.stack.endRule == nil || *first.stack.endRule != `</tag>` {
		t.Fatalf("resolved end = %v, want </tag>", first.stack.endRule)
	}
	if got := first.stack.contentNameScopesList.scopeNames(); !reflect.DeepEqual(got, []string{"source.test", "meta.tag", "meta.tag.content"}) {
		t.Fatalf("content scopes = %v", got)
	}

	secondHandler := &recordingTokenHandler{}
	second := tokenizeString(
		grammar,
		oniguruma.NewString("</tag>\n"),
		false,
		0,
		first.stack.reset(),
		secondHandler,
		true,
		0,
	)
	if second.stack.ruleID != rootID || second.stack.parent != nil {
		t.Fatalf("second line stack = %v, want root", second.stack)
	}
	if got := secondHandler.tokens[0].scopes; !reflect.DeepEqual(got, []string{"source.test", "meta.tag"}) {
		t.Fatalf("end delimiter scopes = %v", got)
	}
}

func TestCheckWhileRunsBottomToTopAndAdvancesAnchor(t *testing.T) {
	rootID, lowerID, upperID := ruleID(1), ruleID(2), ruleID(3)
	root := newIncludeOnlyRule(nil, rootID, nil, nil, compilePatternsResult{})
	lower := newBeginWhileRule(nil, lowerID, nil, nil, ``, nil, `a`, nil, compilePatternsResult{})
	upper := newBeginWhileRule(nil, upperID, nil, nil, ``, nil, `\Gb`, nil, compilePatternsResult{})
	grammar := &tokenizerTestGrammar{rules: map[ruleID]rule{
		rootID:  root,
		lowerID: lower,
		upperID: upper,
	}}
	rootStack := tokenizerTestRootStack(rootID)
	lowerStack := rootStack.push(lowerID, -1, -1, false, nil, rootStack.nameScopesList, rootStack.contentNameScopesList)
	upperStack := lowerStack.push(upperID, -1, -1, false, nil, lowerStack.nameScopesList, lowerStack.contentNameScopesList)

	checked := checkWhile(grammar, oniguruma.NewString("ab\n"), true, 0, upperStack, &recordingTokenHandler{})
	if checked.stack != upperStack {
		t.Fatal("both while conditions matched, but the stack changed")
	}
	if checked.linePos != 2 || checked.anchorPosition != 2 || checked.isFirstLine {
		t.Fatalf("while result = pos %d anchor %d first %v", checked.linePos, checked.anchorPosition, checked.isFirstLine)
	}
}

func TestCheckWhileFailureCutsStackAboveFailedRule(t *testing.T) {
	rootID, lowerID, upperID := ruleID(1), ruleID(2), ruleID(3)
	root := newIncludeOnlyRule(nil, rootID, nil, nil, compilePatternsResult{})
	lower := newBeginWhileRule(nil, lowerID, nil, nil, ``, nil, `a`, nil, compilePatternsResult{})
	upper := newBeginWhileRule(nil, upperID, nil, nil, ``, nil, `z`, nil, compilePatternsResult{})
	grammar := &tokenizerTestGrammar{rules: map[ruleID]rule{
		rootID:  root,
		lowerID: lower,
		upperID: upper,
	}}
	rootStack := tokenizerTestRootStack(rootID)
	lowerStack := rootStack.push(lowerID, -1, -1, false, nil, rootStack.nameScopesList, rootStack.contentNameScopesList)
	upperStack := lowerStack.push(upperID, -1, -1, false, nil, lowerStack.nameScopesList, lowerStack.contentNameScopesList)

	checked := checkWhile(grammar, oniguruma.NewString("ab\n"), true, 0, upperStack, &recordingTokenHandler{})
	if checked.stack != lowerStack {
		t.Fatalf("stack = %v, want lower while stack", checked.stack)
	}
}

func TestInjectionWinsEarlierOrLeftPriorityTie(t *testing.T) {
	for _, test := range []struct {
		name            string
		normalPattern   string
		priority        int
		wantMatchedRule ruleID
	}{
		{name: "earlier", normalPattern: `y`, priority: 0, wantMatchedRule: 3},
		{name: "left priority tie", normalPattern: `x`, priority: -1, wantMatchedRule: 3},
		{name: "normal wins unprioritized tie", normalPattern: `x`, priority: 0, wantMatchedRule: 2},
	} {
		t.Run(test.name, func(t *testing.T) {
			rootID, normalID, injectionID := ruleID(1), ruleID(2), ruleID(3)
			normal := newMatchRule(nil, normalID, nil, test.normalPattern, nil)
			injected := newMatchRule(nil, injectionID, nil, `x`, nil)
			root := newIncludeOnlyRule(nil, rootID, nil, nil, compilePatternsResult{patterns: []ruleID{normalID}})
			grammar := &tokenizerTestGrammar{
				rules: map[ruleID]rule{rootID: root, normalID: normal, injectionID: injected},
				injections: []injection{{
					matcher:  func([]string) bool { return true },
					priority: test.priority,
					ruleID:   injectionID,
				}},
			}

			matched := matchRuleOrInjections(
				grammar,
				oniguruma.NewString("xy\n"),
				true,
				0,
				tokenizerTestRootStack(rootID),
				-1,
			)
			if matched == nil || matched.matchedRuleID != test.wantMatchedRule {
				t.Fatalf("matched %#v, want rule %d", matched, test.wantMatchedRule)
			}
		})
	}
}

func TestMatchRuleAtResolvesAAndGFromTokenizerPosition(t *testing.T) {
	t.Run("A only on first line", func(t *testing.T) {
		rootID, anchoredID := ruleID(1), ruleID(2)
		anchored := newMatchRule(nil, anchoredID, nil, `\Afirst`, nil)
		root := newIncludeOnlyRule(nil, rootID, nil, nil, compilePatternsResult{patterns: []ruleID{anchoredID}})
		grammar := &tokenizerTestGrammar{rules: map[ruleID]rule{rootID: root, anchoredID: anchored}}
		line := oniguruma.NewString("first\n")
		stack := tokenizerTestRootStack(rootID)

		if got := matchRuleAt(grammar, line, true, 0, stack, -1); got == nil || got.matchedRuleID != anchoredID {
			t.Fatalf("first-line match = %#v, want rule %d", got, anchoredID)
		}
		if got := matchRuleAt(grammar, line, false, 0, stack, -1); got != nil {
			t.Fatalf("non-first-line \\A unexpectedly matched: %#v", got)
		}
	})

	t.Run("G only at anchor position", func(t *testing.T) {
		rootID, anchoredID := ruleID(1), ruleID(2)
		anchored := newMatchRule(nil, anchoredID, nil, `\Gx`, nil)
		root := newIncludeOnlyRule(nil, rootID, nil, nil, compilePatternsResult{patterns: []ruleID{anchoredID}})
		grammar := &tokenizerTestGrammar{rules: map[ruleID]rule{rootID: root, anchoredID: anchored}}
		line := oniguruma.NewString("ax\n")
		stack := tokenizerTestRootStack(rootID)

		if got := matchRuleAt(grammar, line, false, 1, stack, 1); got == nil || got.matchedRuleID != anchoredID {
			t.Fatalf("anchor-position match = %#v, want rule %d", got, anchoredID)
		}
		if got := matchRuleAt(grammar, line, false, 1, stack, -1); got != nil {
			t.Fatalf("disabled \\G unexpectedly matched: %#v", got)
		}
	})
}

func TestHandleCapturesRetokenizesCapturedText(t *testing.T) {
	rootID, outerID, retokenizeID, innerID := ruleID(1), ruleID(2), ruleID(3), ruleID(4)
	inner := newMatchRule(nil, innerID, tokenizerTestString("inner.a"), `a`, nil)
	retokenize := newIncludeOnlyRule(nil, retokenizeID, nil, nil, compilePatternsResult{patterns: []ruleID{innerID}})
	outer := newMatchRule(nil, outerID, tokenizerTestString("outer"), `(ab)`, []*captureRule{
		nil,
		newCaptureRule(nil, 5, tokenizerTestString("meta.capture"), nil, retokenizeID),
	})
	root := newIncludeOnlyRule(nil, rootID, nil, nil, compilePatternsResult{patterns: []ruleID{outerID}})
	grammar := &tokenizerTestGrammar{rules: map[ruleID]rule{
		rootID:       root,
		outerID:      outer,
		retokenizeID: retokenize,
		innerID:      inner,
	}}
	handler := &recordingTokenHandler{}

	tokenizeString(grammar, oniguruma.NewString("ab\n"), true, 0, tokenizerTestRootStack(rootID), handler, true, 0)
	want := []recordedToken{
		{start: 0, end: 1, scopes: []string{"source.test", "outer", "meta.capture", "inner.a"}},
		{start: 1, end: 2, scopes: []string{"source.test", "outer", "meta.capture"}},
		{start: 2, end: 3, scopes: []string{"source.test"}},
	}
	if !reflect.DeepEqual(handler.tokens, want) {
		t.Fatalf("tokens mismatch:\n got %#v\nwant %#v", handler.tokens, want)
	}
}

func TestHandleCapturesNestsAndClosesScopesInCaptureOrder(t *testing.T) {
	rootID := ruleID(1)
	root := newIncludeOnlyRule(nil, rootID, nil, nil, compilePatternsResult{})
	grammar := &tokenizerTestGrammar{rules: map[ruleID]rule{rootID: root}}
	stack := tokenizerTestRootStack(rootID)
	handler := &recordingTokenHandler{}
	captureRules := []*captureRule{
		nil,
		newCaptureRule(nil, 2, tokenizerTestString("outer"), nil, 0),
		newCaptureRule(nil, 3, tokenizerTestString("first"), nil, 0),
		newCaptureRule(nil, 4, tokenizerTestString("second"), nil, 0),
	}
	captures := []oniguruma.Capture{
		{Start: 0, End: 6},
		{Start: 0, End: 6},
		{Start: 1, End: 2},
		{Start: 4, End: 5},
	}

	handleCaptures(grammar, oniguruma.NewString("abcdef"), true, stack, handler, captureRules, captures)
	want := []recordedToken{
		{start: 0, end: 1, scopes: []string{"source.test", "outer"}},
		{start: 1, end: 2, scopes: []string{"source.test", "outer", "first"}},
		{start: 2, end: 4, scopes: []string{"source.test", "outer"}},
		{start: 4, end: 5, scopes: []string{"source.test", "outer", "second"}},
		{start: 5, end: 6, scopes: []string{"source.test", "outer"}},
	}
	if !reflect.DeepEqual(handler.tokens, want) {
		t.Fatalf("tokens mismatch:\n got %#v\nwant %#v", handler.tokens, want)
	}
}

func TestTokenizeStringZeroWidthLoopGuards(t *testing.T) {
	t.Run("pushed and popped begin-end", func(t *testing.T) {
		rootID, beginID := ruleID(1), ruleID(2)
		root := newIncludeOnlyRule(nil, rootID, nil, nil, compilePatternsResult{})
		begin := newBeginEndRule(nil, beginID, nil, nil, `x`, nil, ``, nil, false, compilePatternsResult{})
		grammar := &tokenizerTestGrammar{rules: map[ruleID]rule{rootID: root, beginID: begin}}
		rootStack := tokenizerTestRootStack(rootID)
		active := rootStack.push(beginID, 0, -1, false, nil, rootStack.nameScopesList, rootStack.contentNameScopesList)

		result := tokenizeString(grammar, oniguruma.NewString("x\n"), true, 0, active, &recordingTokenHandler{}, false, 0)
		if result.stack != active {
			t.Fatalf("guard 1 stack = %v, want restored active stack", result.stack)
		}
	})

	t.Run("repeated begin-end push", func(t *testing.T) {
		rootID, beginID := ruleID(1), ruleID(2)
		root := newIncludeOnlyRule(nil, rootID, nil, nil, compilePatternsResult{})
		begin := newBeginEndRule(nil, beginID, nil, nil, ``, nil, `(?!)`, nil, false, compilePatternsResult{patterns: []ruleID{beginID}})
		grammar := &tokenizerTestGrammar{rules: map[ruleID]rule{rootID: root, beginID: begin}}
		rootStack := tokenizerTestRootStack(rootID)
		active := rootStack.push(beginID, 0, -1, false, nil, rootStack.nameScopesList, rootStack.contentNameScopesList)

		result := tokenizeString(grammar, oniguruma.NewString("x\n"), true, 0, active, &recordingTokenHandler{}, false, 0)
		if result.stack != active {
			t.Fatalf("guard 2 stack = %v, want active stack", result.stack)
		}
	})

	t.Run("repeated begin-while push", func(t *testing.T) {
		rootID, beginID := ruleID(1), ruleID(2)
		root := newIncludeOnlyRule(nil, rootID, nil, nil, compilePatternsResult{})
		begin := newBeginWhileRule(nil, beginID, nil, nil, ``, nil, `x`, nil, compilePatternsResult{patterns: []ruleID{beginID}})
		grammar := &tokenizerTestGrammar{rules: map[ruleID]rule{rootID: root, beginID: begin}}
		rootStack := tokenizerTestRootStack(rootID)
		active := rootStack.push(beginID, 0, -1, false, nil, rootStack.nameScopesList, rootStack.contentNameScopesList)

		result := tokenizeString(grammar, oniguruma.NewString("x\n"), true, 0, active, &recordingTokenHandler{}, false, 0)
		if result.stack != active {
			t.Fatalf("guard 3 stack = %v, want active stack", result.stack)
		}
	})

	t.Run("zero-width match pops enclosing state", func(t *testing.T) {
		rootID, matchID := ruleID(1), ruleID(2)
		root := newIncludeOnlyRule(nil, rootID, nil, nil, compilePatternsResult{})
		match := newMatchRule(nil, matchID, nil, ``, nil)
		grammar := &tokenizerTestGrammar{rules: map[ruleID]rule{rootID: root, matchID: match}}
		rootStack := tokenizerTestRootStack(rootID)
		active := rootStack.push(matchID, 0, -1, false, nil, rootStack.nameScopesList, rootStack.contentNameScopesList)

		result := tokenizeString(grammar, oniguruma.NewString("x\n"), true, 0, active, &recordingTokenHandler{}, false, 0)
		if result.stack != rootStack {
			t.Fatalf("guard 4 stack = %v, want enclosing root", result.stack)
		}
	})
}

func TestTokenizeStringHonorsTimeLimit(t *testing.T) {
	rootID := ruleID(1)
	root := newIncludeOnlyRule(nil, rootID, nil, nil, compilePatternsResult{})
	grammar := &tokenizerTestGrammar{rules: map[ruleID]rule{rootID: root}}
	stack := tokenizerTestRootStack(rootID)
	handler := &recordingTokenHandler{}

	result := tokenizeString(grammar, oniguruma.NewString("text\n"), true, 0, stack, handler, true, time.Nanosecond)
	if !result.stoppedEarly {
		t.Fatal("expected the one-nanosecond time limit to stop tokenization")
	}
	if result.stack != stack || len(handler.tokens) != 0 {
		t.Fatalf("time-limit result changed state or emitted tokens: %#v %#v", result.stack, handler.tokens)
	}
}

func TestPushAttributedScopeLooksUpEachSpaceSeparatedScope(t *testing.T) {
	grammar := &tokenizerTestGrammar{attributes: map[string]uint32{"meta.outer": 7, "meta.inner": 11}}
	root := newAttributedScopeRoot("source.test", 1)
	got := pushAttributedScope(root, "meta.outer meta.inner", grammar)

	if names := got.scopeNames(); !reflect.DeepEqual(names, []string{"source.test", "meta.outer", "meta.inner"}) {
		t.Fatalf("scope names = %v", names)
	}
	if got.tokenAttributes != 11 || got.parent.tokenAttributes != 7 {
		t.Fatalf("attributes = inner %d outer %d", got.tokenAttributes, got.parent.tokenAttributes)
	}
}
