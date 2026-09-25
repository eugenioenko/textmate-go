package oniguruma

import "testing"

func TestNestedCharacterClassUnionMatchesYAMLPlainKey(t *testing.T) {
	pattern := `(?x)
        (?=
            (?x:
                  [^\s[-?:,\[\]{}#&*!|>'"%@` + "`" + `]]
                | [?:-] \S
            )
            (
                  [^\s:]
                | : \S
                | \s+ (?![#\s])
            )*
            \s*
            :
            (\s|$)
        )`
	scanner := NewScanner([]string{pattern})
	if diagnostics := scanner.Diagnostics(); len(diagnostics) != 0 {
		t.Fatalf("unexpected scanner diagnostics: %#v", diagnostics)
	}
	match := scanner.FindNextMatch(NewString("- run:\n"), 1, FindOptionNone)
	if match == nil {
		t.Fatal("YAML plain-key lookahead did not match")
	}
	if got, want := match.Captures[0], (Capture{Start: 2, End: 2}); got != want {
		t.Fatalf("match = %#v, want %#v", got, want)
	}
}

func TestNestedCharacterClassUnionSupportsComplementMembers(t *testing.T) {
	pattern := `(?:^|\G)(#{1,6})\s*(?=[\S[^#]])`
	scanner := NewScanner([]string{pattern})
	if diagnostics := scanner.Diagnostics(); len(diagnostics) != 0 {
		t.Fatalf("unexpected scanner diagnostics: %#v", diagnostics)
	}
	match := scanner.FindNextMatch(NewString("## This is *great* stuff\n"), 0, FindOptionNone)
	if match == nil {
		t.Fatal("Markdown heading with a nested complemented class did not match")
	}
	if got, want := match.Captures[0], (Capture{Start: 0, End: 3}); got != want {
		t.Fatalf("match = %#v, want %#v", got, want)
	}
}

func TestQuantifierBindsWholeNestedNegatedClass(t *testing.T) {
	scanner := NewScanner([]string{`[^\s[-?:]]+`})
	if diagnostics := scanner.Diagnostics(); len(diagnostics) != 0 {
		t.Fatalf("unexpected scanner diagnostics: %#v", diagnostics)
	}
	match := scanner.FindNextMatch(NewString("run:value"), 0, FindOptionNone)
	if match == nil {
		t.Fatal("nested negated class did not match")
	}
	if got, want := match.Captures[0], (Capture{Start: 0, End: 3}); got != want {
		t.Fatalf("match = %#v, want %#v", got, want)
	}
}
