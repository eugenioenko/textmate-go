package textmate

import (
	"reflect"
	"testing"
	"time"

	"github.com/eugenioenko/textmate-go/oniguruma"
)

func TestTokenizeLineWithOptionsSkipsLongFirstLineBeforeParsing(t *testing.T) {
	match, name := `x`, "letter.x"
	grammar := newGrammar("source.test", &RawGrammar{
		ScopeName: "source.test",
		Patterns:  []*RawRule{{Match: &match, Name: &name}},
	}, nil)

	result := grammar.TokenizeLineWithOptions("ééx", nil, TokenizeOptions{MaxLineRunes: 2})
	if !result.Stopped || result.StoppedReason != StopReasonLineLimit || result.StoppedAt != 0 {
		t.Fatalf("stop = (%v, %v, %d), want line limit at rune 0", result.Stopped, result.StoppedReason, result.StoppedAt)
	}
	if result.RuleStack == nil || result.RuleStack.ruleID != grammar.rootID {
		t.Fatalf("rule stack = %v, want initialized root", result.RuleStack)
	}
	want := []Token{{Start: 0, End: 3, Scopes: []string{"source.test"}}}
	if !tokenValuesEqual(result.Tokens, want) {
		t.Fatalf("tokens = %#v, want fallback %#v", result.Tokens, want)
	}

	// The cap is inclusive, and the unrestricted TokenizeLine API remains
	// equivalent to all-zero options.
	exact := grammar.TokenizeLineWithOptions("ééx", nil, TokenizeOptions{MaxLineRunes: 3})
	unlimited := grammar.TokenizeLine("ééx", nil)
	if exact.Stopped || !reflect.DeepEqual(exact.Tokens, unlimited.Tokens) {
		t.Fatalf("exact-cap result = %#v, unlimited = %#v", exact, unlimited)
	}
}

func TestTokenizeLineWithOptionsByteAndRuneLimitsAreIndependent(t *testing.T) {
	grammar := newGrammar("source.test", &RawGrammar{ScopeName: "source.test"}, nil)
	tests := []struct {
		name    string
		options TokenizeOptions
		stopped bool
	}{
		{name: "byte cap exceeded", options: TokenizeOptions{MaxLineBytes: 3}, stopped: true},
		{name: "byte cap exact", options: TokenizeOptions{MaxLineBytes: 4}},
		{name: "rune cap exceeded", options: TokenizeOptions{MaxLineRunes: 1}, stopped: true},
		{name: "rune cap exact", options: TokenizeOptions{MaxLineRunes: 2}},
		{name: "zero unlimited", options: TokenizeOptions{}},
		{name: "negative unlimited", options: TokenizeOptions{MaxLineBytes: -1, MaxLineRunes: -1, TimeLimit: -1}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := grammar.TokenizeLineWithOptions("éé", nil, test.options)
			if result.Stopped != test.stopped {
				t.Fatalf("Stopped = %v, want %v", result.Stopped, test.stopped)
			}
			if result.Tokens[0].End != 2 {
				t.Fatalf("token end = %d, want rune length 2", result.Tokens[0].End)
			}
		})
	}
}

func TestTokenizeLineWithOptionsLongLinePreservesCarriedState(t *testing.T) {
	begin, end := `"`, `"`
	name, contentName := "string.quoted.test", "string.quoted.content.test"
	grammar := newGrammar("source.test", &RawGrammar{
		ScopeName: "source.test",
		Patterns: []*RawRule{{
			Begin: &begin, End: &end, Name: &name, ContentName: &contentName,
		}},
	}, nil)

	opened := grammar.TokenizeLine(`"open`, nil)
	if opened.RuleStack == nil || opened.RuleStack.parent == nil {
		t.Fatalf("opening line did not carry string state: %v", opened.RuleStack)
	}
	skipped := grammar.TokenizeLineWithOptions("ééé", opened.RuleStack, TokenizeOptions{MaxLineRunes: 2})
	if !skipped.Stopped || skipped.StoppedReason != StopReasonLineLimit || !skipped.RuleStack.Equal(opened.RuleStack) {
		t.Fatalf("skipped state = %v, want state equal to %v", skipped.RuleStack, opened.RuleStack)
	}
	wantScopes := []string{"source.test", name, contentName}
	if got := skipped.Tokens; len(got) != 1 || got[0].Start != 0 || got[0].End != 3 || !reflect.DeepEqual(got[0].Scopes, wantScopes) {
		t.Fatalf("fallback tokens = %#v, want one active-string token through rune 3", got)
	}

	closed := grammar.TokenizeLine(`"`, skipped.RuleStack)
	if closed.RuleStack == nil || closed.RuleStack.parent != nil {
		t.Fatalf("closing line stack = %v, want root", closed.RuleStack)
	}
}

func TestTokenizeLineWithOptionsReportsImmediateSoftTimeStop(t *testing.T) {
	grammar := newGrammar("source.test", &RawGrammar{ScopeName: "source.test"}, nil)
	result := grammar.TokenizeLineWithOptions("a😀b", nil, TokenizeOptions{TimeLimit: time.Nanosecond})
	if !result.Stopped || result.StoppedReason != StopReasonTimeLimit || result.StoppedAt != 0 {
		t.Fatalf("stop = (%v, %v, %d), want time limit at rune 0", result.Stopped, result.StoppedReason, result.StoppedAt)
	}
	if result.RuleStack == nil {
		t.Fatal("time-limited first line did not return initialized state")
	}
	if got := result.Tokens; len(got) != 1 || got[0].End != 3 {
		t.Fatalf("fallback tokens = %#v, want full rune coverage", got)
	}
}

func TestTokenizeLineWithOptionsTimeStopPreservesCarriedState(t *testing.T) {
	begin, end := `"`, `"`
	name, contentName := "string.quoted.test", "string.quoted.content.test"
	grammar := newGrammar("source.test", &RawGrammar{
		ScopeName: "source.test",
		Patterns: []*RawRule{{
			Begin: &begin, End: &end, Name: &name, ContentName: &contentName,
		}},
	}, nil)
	opened := grammar.TokenizeLine(`"open`, nil)

	result := grammar.TokenizeLineWithOptions("a😀b", opened.RuleStack, TokenizeOptions{TimeLimit: time.Nanosecond})
	if !result.Stopped || result.StoppedReason != StopReasonTimeLimit || result.StoppedAt != 0 {
		t.Fatalf("stop = (%v, %v, %d), want time limit at rune 0", result.Stopped, result.StoppedReason, result.StoppedAt)
	}
	if result.RuleStack == nil || !result.RuleStack.Equal(opened.RuleStack) {
		t.Fatalf("time-limited state = %v, want state equal to %v", result.RuleStack, opened.RuleStack)
	}
	wantScopes := []string{"source.test", name, contentName}
	if got := result.Tokens; len(got) != 1 || got[0].End != 3 || !reflect.DeepEqual(got[0].Scopes, wantScopes) {
		t.Fatalf("fallback tokens = %#v, want carried scopes through rune 3", got)
	}
}

func TestTokenizeStringBudgetReportsPartialRuneOffsetAndCoversTail(t *testing.T) {
	rootID, matchID := ruleID(1), ruleID(2)
	match := newMatchRule(nil, matchID, tokenizerTestString("letter.test"), `[a-z]`, nil)
	root := newIncludeOnlyRule(nil, rootID, nil, nil, compilePatternsResult{patterns: []ruleID{matchID}})
	grammar := &tokenizerTestGrammar{rules: map[ruleID]rule{rootID: root, matchID: match}}
	handler := &lineTokenHandler{}
	stack := tokenizerTestRootStack(rootID)

	result := tokenizeStringWithBudget(
		grammar,
		oniguruma.NewString("éa界\n"),
		true,
		0,
		stack,
		handler,
		true,
		tokenizationBudget{stop: func() bool { return handler.lastEnd >= 2 }},
	)
	if !result.stoppedEarly || result.stoppedAt != 2 || result.stack != stack {
		t.Fatalf("result = %+v, want stop at rune 2 on root stack", result)
	}
	want := []Token{
		{Start: 0, End: 1, Scopes: []string{"source.test"}},
		{Start: 1, End: 2, Scopes: []string{"source.test", "letter.test"}},
		{Start: 2, End: 3, Scopes: []string{"source.test"}},
	}
	if got := handler.result(result.stack, 3); !tokenValuesEqual(got, want) {
		t.Fatalf("partial tokens = %#v, want %#v", got, want)
	}
}

func TestTokenizeStringBudgetIncludesWhileConditionSetup(t *testing.T) {
	rootID, whileID := ruleID(1), ruleID(2)
	root := newIncludeOnlyRule(nil, rootID, nil, nil, compilePatternsResult{})
	whileRule := newBeginWhileRule(nil, whileID, nil, nil, ``, nil, `a`, nil, compilePatternsResult{})
	grammar := &tokenizerTestGrammar{rules: map[ruleID]rule{rootID: root, whileID: whileRule}}
	rootStack := tokenizerTestRootStack(rootID)
	active := rootStack.push(whileID, -1, -1, false, nil, rootStack.nameScopesList, rootStack.contentNameScopesList)
	handler := &recordingTokenHandler{}
	checks := 0

	result := tokenizeStringWithBudget(
		grammar,
		oniguruma.NewString("a\n"),
		false,
		0,
		active,
		handler,
		true,
		tokenizationBudget{stop: func() bool {
			checks++
			return checks == 2
		}},
	)
	if !result.stoppedEarly || result.stoppedAt != 0 || result.stack != active {
		t.Fatalf("result = %+v, want unchanged active stack stopped before while search", result)
	}
	if len(handler.tokens) != 0 {
		t.Fatalf("while-budget stop emitted tokens: %#v", handler.tokens)
	}
}

func TestStopReasonString(t *testing.T) {
	want := map[StopReason]string{
		StopReasonNone:      "none",
		StopReasonTimeLimit: "time_limit",
		StopReasonLineLimit: "line_limit",
		StopReason(255):     "unknown",
	}
	for reason, value := range want {
		if got := reason.String(); got != value {
			t.Fatalf("%d.String() = %q, want %q", reason, got, value)
		}
	}
}
