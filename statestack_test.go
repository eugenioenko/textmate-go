package textmate

import (
	"reflect"
	"testing"
)

func TestInitialStateAndStackOperations(t *testing.T) {
	if endRuleID != -1 || whileRuleID != -2 {
		t.Fatalf("reserved rule IDs changed: end=%d while=%d", endRuleID, whileRuleID)
	}
	if InitialState == nil {
		t.Fatal("InitialState is nil")
	}
	if InitialState.pop() != nil {
		t.Fatal("InitialState.pop() should return nil")
	}
	if InitialState.safePop() != InitialState {
		t.Fatal("InitialState.safePop() should preserve the root")
	}
	if InitialState.clone() != InitialState {
		t.Fatal("clone should reuse an immutable stack")
	}

	scopes := newAttributedScopeRoot("source.go", 1)
	child := InitialState.push(7, 3, 2, false, nil, scopes, scopes)
	if got := child.getRuleID(); got != 7 {
		t.Fatalf("getRuleID() = %d, want 7", got)
	}
	if got := child.getEnterPos(); got != 3 {
		t.Fatalf("getEnterPos() = %d, want 3", got)
	}
	if got := child.getAnchorPos(); got != 2 {
		t.Fatalf("getAnchorPos() = %d, want 2", got)
	}
	if child.pop() != InitialState || child.safePop() != InitialState {
		t.Fatal("pop operations did not return the parent")
	}
	if got, want := child.String(), "[(0, , ),(7, source.go, source.go)]"; got != want {
		t.Fatalf("String() = %q, want %q", got, want)
	}
}

func TestStateStackEqual(t *testing.T) {
	rootA := newAttributedScopeRoot("source.go", 10)
	rootB := newAttributedScopeRoot("source.go", 10)
	contentA := rootA.pushAttributed("comment.block.go", 20)
	contentB := rootB.pushAttributed("comment.block.go", 20)
	end := `\*/`

	a := newStateStack(nil, 1, 4, 8, false, nil, rootA, rootA).
		push(2, 12, 9, true, &end, rootA, contentA)
	b := newStateStack(nil, 1, -1, -1, true, nil, rootB, rootB).
		push(2, 99, 100, false, &end, rootB, contentB)

	if !a.Equal(a) {
		t.Fatal("stack should equal itself")
	}
	if !a.Equal(b) || !b.Equal(a) {
		t.Fatal("structurally equal stacks with different transient fields are unequal")
	}
	if a.Equal(nil) || (*StateStack)(nil).Equal(a) {
		t.Fatal("non-nil and nil stacks compare equal")
	}
	if !(*StateStack)(nil).Equal(nil) {
		t.Fatal("two nil stacks should compare equal")
	}

	otherRule := b.parent.push(3, 99, 100, true, &end, rootB, contentB)
	if a.Equal(otherRule) {
		t.Fatal("stacks with different rule IDs compare equal")
	}
	otherEnd := b.withEndRule("END")
	if a.Equal(otherEnd) {
		t.Fatal("stacks with different end rules compare equal")
	}
	otherContent := rootB.pushAttributed("string.quoted.go", 20)
	if a.Equal(b.withContentNameScopesList(otherContent)) {
		t.Fatal("stacks with different content scopes compare equal")
	}
	otherAttributes := rootB.pushAttributed("comment.block.go", 21)
	if a.Equal(b.withContentNameScopesList(otherAttributes)) {
		t.Fatal("stacks with different scope attributes compare equal")
	}
}

func TestStateStackResetIsImmutable(t *testing.T) {
	scopes := newAttributedScopeRoot("source.go", 1)
	root := newStateStack(nil, 1, 5, 6, false, nil, scopes, scopes)
	original := root.push(2, 7, 8, false, nil, scopes, scopes)

	reset := original.reset()
	if reset == original || reset.parent == root {
		t.Fatal("reset reused a node containing transient positions")
	}
	if got := [2]int{reset.enterPos, reset.anchorPos}; got != [2]int{-1, -1} {
		t.Fatalf("reset top positions = %v, want [-1 -1]", got)
	}
	if got := [2]int{reset.parent.enterPos, reset.parent.anchorPos}; got != [2]int{-1, -1} {
		t.Fatalf("reset parent positions = %v, want [-1 -1]", got)
	}
	if got := [2]int{original.enterPos, original.anchorPos}; got != [2]int{7, 8} {
		t.Fatalf("original top was mutated: %v", got)
	}
	if got := [2]int{root.enterPos, root.anchorPos}; got != [2]int{5, 6} {
		t.Fatalf("original parent was mutated: %v", got)
	}
	if !original.Equal(reset) {
		t.Fatal("reset changed persistent tokenizer state")
	}
	if reset.reset() != reset {
		t.Fatal("reset should reuse an already-reset immutable stack")
	}
}

func TestStateStackUpdatesAreImmutable(t *testing.T) {
	rootScopes := newAttributedScopeRoot("source.go", 1)
	commentScopes := rootScopes.pushAttributed("comment.block.go", 2)
	stack := newStateStack(nil, 1, -1, -1, false, nil, rootScopes, rootScopes)

	withContent := stack.withContentNameScopesList(commentScopes)
	if stack.contentNameScopesList != rootScopes {
		t.Fatal("withContentNameScopesList mutated the original")
	}
	if withContent.contentNameScopesList != commentScopes {
		t.Fatal("withContentNameScopesList did not install the new scopes")
	}
	if withContent.withContentNameScopesList(commentScopes) != withContent {
		t.Fatal("withContentNameScopesList should reuse an unchanged state")
	}

	withEnd := stack.withEndRule("stop")
	if stack.endRule != nil {
		t.Fatal("withEndRule mutated the original")
	}
	if withEnd.endRule == nil || *withEnd.endRule != "stop" {
		t.Fatal("withEndRule did not install the end rule")
	}
	if withEnd.withEndRule("stop") != withEnd {
		t.Fatal("withEndRule should reuse an unchanged state")
	}

	end := "before"
	defensive := newStateStack(nil, 1, -1, -1, false, &end, rootScopes, rootScopes)
	end = "after"
	if got := *defensive.endRule; got != "before" {
		t.Fatalf("state retained a mutable end-rule pointer: got %q", got)
	}
}

func TestAttributedScopeStackExtensions(t *testing.T) {
	root := newAttributedScopeRoot("source.go", 1)
	stack := root.pushAttributed("meta.function.go entity.name.function.go", 2)
	if !stack.scopePath.extends(root.scopePath) || root.scopePath.extends(stack.scopePath) {
		t.Fatal("scope-path ancestry is incorrect")
	}
	if got, want := stack.scopeNames(), []string{
		"source.go",
		"meta.function.go",
		"entity.name.function.go",
	}; !reflect.DeepEqual(got, want) {
		t.Fatalf("scopeNames() = %#v, want %#v", got, want)
	}

	extension, ok := stack.extensionFrom(root)
	if !ok {
		t.Fatal("extensionFrom(root) was not defined")
	}
	rebuilt := attributedScopeStackFromExtension(root, extension)
	if !attributedScopeStacksEqual(stack, rebuilt) {
		t.Fatal("extension round trip changed attributed scopes")
	}
	if got := rebuilt.scopeNames(); !reflect.DeepEqual(got, stack.scopeNames()) {
		t.Fatalf("rebuilt scopes = %#v, want %#v", got, stack.scopeNames())
	}

	unrelated := newAttributedScopeRoot("text.html", 1)
	if _, ok := stack.extensionFrom(unrelated); ok {
		t.Fatal("extensionFrom(unrelated) unexpectedly succeeded")
	}
}

func TestStateStackFrameRoundTrip(t *testing.T) {
	rootScopes := newAttributedScopeRoot("source.go", 1)
	nameScopes := rootScopes.pushAttributed("string.quoted.double.go", 2)
	contentScopes := nameScopes.pushAttributed("meta.embedded.go", 3)
	end := `"`
	parent := newStateStack(nil, 1, -1, -1, false, nil, rootScopes, rootScopes)
	stack := parent.push(2, 4, 3, true, &end, nameScopes, contentScopes)

	rebuilt := pushStateStackFrame(parent, stack.toStateStackFrame())
	if !stack.Equal(rebuilt) {
		t.Fatalf("frame round trip changed stack:\noriginal %s\nrebuilt  %s", stack, rebuilt)
	}
	if rebuilt.enterPos != 4 || rebuilt.anchorPos != 3 || !rebuilt.beginRuleCapturedEOL {
		t.Fatalf("frame round trip lost transient fields: %#v", rebuilt)
	}
}

func TestHasSameRuleAs(t *testing.T) {
	scopes := newAttributedScopeRoot("source.go", 1)
	root := newStateStack(nil, 1, 0, -1, false, nil, scopes, scopes)
	stack := root.push(2, 5, -1, false, nil, scopes, scopes).
		push(3, 5, -1, false, nil, scopes, scopes)

	if !stack.hasSameRuleAs(newStateStack(nil, 2, 5, -1, false, nil, scopes, scopes)) {
		t.Fatal("same rule at the same enter position was not found")
	}
	if stack.hasSameRuleAs(newStateStack(nil, 2, 4, -1, false, nil, scopes, scopes)) {
		t.Fatal("search crossed an enter-position boundary")
	}
}
