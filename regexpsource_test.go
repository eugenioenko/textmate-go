package textmate

import (
	"testing"

	"github.com/eugenioenko/textmate-go/oniguruma"
)

func TestRegExpSourceNormalizesEndOfStringAndFindsFlags(t *testing.T) {
	tests := []struct {
		name        string
		source      string
		wantSource  string
		wantAnchor  bool
		wantBackRef bool
	}{
		{name: "empty", source: ""},
		{name: "end of string", source: `foo\z`, wantSource: `foo$(?!\n)(?<!\n)`},
		{name: "anchors", source: `\Afoo\G`, wantSource: `\Afoo\G`, wantAnchor: true},
		{name: "back reference", source: `^\12$`, wantSource: `^\12$`, wantBackRef: true},
		{name: "escaped slash is skipped", source: `\\A\\1`, wantSource: `\\A\\1`, wantBackRef: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			source := newRegExpSource(test.source, 3)
			wantSource := test.wantSource
			if wantSource == "" {
				wantSource = test.source
			}
			if source.source != wantSource {
				t.Fatalf("source = %q, want %q", source.source, wantSource)
			}
			if source.hasAnchor != test.wantAnchor {
				t.Fatalf("hasAnchor = %v, want %v", source.hasAnchor, test.wantAnchor)
			}
			if source.hasBackReferences != test.wantBackRef {
				t.Fatalf("hasBackReferences = %v, want %v", source.hasBackReferences, test.wantBackRef)
			}
		})
	}
}

func TestRegExpSourceResolvesAnchorsLikeVSCodeTextmate(t *testing.T) {
	source := newRegExpSource(`[\A]\A-\G-\\A`, 1)
	tests := []struct {
		name           string
		allowA, allowG bool
		want           string
	}{
		{name: "neither", want: "[\\\uffff]\\\uffff-\\\uffff-\\\\A"},
		{name: "A", allowA: true, want: `[\A]\A-\` + "\uffff" + `-\\A`},
		{name: "G", allowG: true, want: "[\\\uffff]\\\uffff-\\G-\\\\A"},
		{name: "both", allowA: true, allowG: true, want: `[\A]\A-\G-\\A`},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := source.resolveAnchors(test.allowA, test.allowG); got != test.want {
				t.Fatalf("resolveAnchors() = %q, want %q", got, test.want)
			}
		})
	}

}

func TestRegExpSourceResolvesBackReferencesInRuneOffsets(t *testing.T) {
	source := newRegExpSource(`^\1/\2/\99$`, endRuleID)
	line := "a😀.[x]"
	captures := []oniguruma.Capture{
		{Start: 0, End: 6},
		{Start: 1, End: 6},
		{Start: -1, End: -1},
	}
	want := `^😀\.\[x\]//$`
	if got := source.resolveBackReferences(line, captures); got != want {
		t.Fatalf("resolveBackReferences() = %q, want %q", got, want)
	}
}

func TestRegExpSourceCloneKeepsIndependentDynamicSource(t *testing.T) {
	original := newRegExpSource(`\A\1`, endRuleID)
	clone := original.clone()
	clone.setSource(`\Gresolved`)

	if original.source != `\A\1` {
		t.Fatalf("original source changed to %q", original.source)
	}
	if got := clone.resolveAnchors(true, false); got != "\\\uffffresolved" {
		t.Fatalf("clone anchor cache was not rebuilt: %q", got)
	}
	if !clone.hasBackReferences {
		t.Fatal("dynamic source must retain the original back-reference flag")
	}
}

func TestRegExpSourceListCachesCompiledRulesAndInvalidatesOnChange(t *testing.T) {
	list := newRegExpSourceList()
	list.push(newRegExpSource(`\Ax`, 10))
	list.push(newRegExpSource(`x`, 11))

	plain := list.compile()
	if plain != list.compile() {
		t.Fatal("plain compiled rule was not cached")
	}
	a1g1 := list.compileAG(true, true)
	if a1g1 != list.compileAG(true, true) {
		t.Fatal("anchor variant was not cached")
	}
	if a1g1 == list.compileAG(false, true) {
		t.Fatal("different anchor variants share a compiled rule")
	}
	unchanged := list.compileAG(true, true)
	list.setSource(0, `\Ax`)
	if unchanged != list.compileAG(true, true) {
		t.Fatal("setting an unchanged source invalidated its cache")
	}

	list.setSource(0, `\Ay`)
	if plain == list.compile() {
		t.Fatal("plain cache survived a source change")
	}
	if a1g1 == list.compileAG(true, true) {
		t.Fatal("anchor cache survived a source change")
	}

	list.compileAG(false, false)
	list.compileAG(false, true)
	list.compileAG(true, false)
	list.dispose()
	if list.cached != nil || list.anchorCache.a0g0 != nil || list.anchorCache.a0g1 != nil ||
		list.anchorCache.a1g0 != nil || list.anchorCache.a1g1 != nil {
		t.Fatal("dispose did not clear all compiled caches")
	}
}

func TestCompiledRulePreservesSourceOrderForMatchPrecedence(t *testing.T) {
	list := newRegExpSourceList()
	list.push(newRegExpSource(`x`, 20))
	list.push(newRegExpSource(`x`, 21))
	list.unshift(newRegExpSource(`z`, endRuleID))

	if list.length() != 3 {
		t.Fatalf("length = %d, want 3", list.length())
	}
	compiled := list.compile()
	match := compiled.findNextMatch(oniguruma.NewString("zx"), 0, oniguruma.FindOptionNone)
	if match == nil || match.ruleID != endRuleID {
		t.Fatalf("first match = %#v, want end rule", match)
	}
	match = compiled.findNextMatch(oniguruma.NewString("x"), 0, oniguruma.FindOptionNone)
	if match == nil || match.ruleID != 20 {
		t.Fatalf("tie winner = %#v, want rule 20", match)
	}
	if match.captureIndices[0] != (oniguruma.Capture{Start: 0, End: 1}) {
		t.Fatalf("capture = %#v, want [0,1]", match.captureIndices[0])
	}

	earliest := newCompiledRule([]string{`z`, `x`}, []ruleID{40, 41})
	match = earliest.findNextMatch(oniguruma.NewString("xz"), 0, oniguruma.FindOptionNone)
	if match == nil || match.ruleID != 41 {
		t.Fatalf("earliest match = %#v, want rule 41", match)
	}
}

func TestCompiledAnchorVariantsControlMatches(t *testing.T) {
	list := newRegExpSourceList()
	list.push(newRegExpSource(`\Ax`, 30))
	list.push(newRegExpSource(`\Gx`, 31))
	input := oniguruma.NewString("x")

	tests := []struct {
		name           string
		allowA, allowG bool
		want           ruleID
		wantMatch      bool
	}{
		{name: "both", allowA: true, allowG: true, want: 30, wantMatch: true},
		{name: "only A", allowA: true, want: 30, wantMatch: true},
		{name: "only G", allowG: true, want: 31, wantMatch: true},
		{name: "neither"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			match := list.compileAG(test.allowA, test.allowG).findNextMatch(input, 0, oniguruma.FindOptionNone)
			if !test.wantMatch {
				if match != nil {
					t.Fatalf("match = %#v, want nil", match)
				}
				return
			}
			if match == nil || match.ruleID != test.want {
				t.Fatalf("match = %#v, want rule %d", match, test.want)
			}
		})
	}
}

func TestCompiledRuleString(t *testing.T) {
	compiled := newCompiledRule([]string{"a", "b"}, []ruleID{2, endRuleID})
	if got, want := compiled.String(), "   - 2: a\n   - -1: b"; got != want {
		t.Fatalf("String() = %q, want %q", got, want)
	}
}
