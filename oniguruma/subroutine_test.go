package oniguruma

import "testing"

func TestResolvableNamedSubroutineExpands(t *testing.T) {
	scanner := NewScanner([]string{`(?<name>[a-z]+):\g<name>`})
	if diagnostics := scanner.Diagnostics(); len(diagnostics) != 0 {
		t.Fatalf("unexpected diagnostics: %+v", diagnostics)
	}
	match := scanner.FindNextMatch(NewString("left:right"), 0, FindOptionNone)
	if match == nil || match.Captures[0] != (Capture{Start: 0, End: 10}) {
		t.Fatalf("match = %+v, want [0,10)", match)
	}
}

func TestUnsupportedSubroutineOnlyNeutralizesItsAlternative(t *testing.T) {
	const thriftField = `(?x)
		(?<ft>
			map\s*<\s*\g<ft>\s*,\s*\g<ft>\s*> |
			set\s*<\s*\g<ft>\s*> |
			list\s*<\s*\g<ft>\s*>\s*(cpp_type(?!\S))? |
			[a-zA-Z_][\w.]*
		)[ \t]*
		(?:([a-zA-Z_][\w.]*)[ \t]*)?`

	scanner := NewScanner([]string{thriftField})
	diagnostics := scanner.Diagnostics()
	if len(diagnostics) != 1 || diagnostics[0].Kind != DiagnosticUnsupportedSyntax {
		t.Fatalf("diagnostics = %+v, want one unsupported-syntax diagnostic", diagnostics)
	}
	match := scanner.FindNextMatch(NewString("string message"), 0, FindOptionNone)
	if match == nil {
		t.Fatal("primitive Thrift field type should still match")
	}
	want := []Capture{
		{Start: 0, End: 14},
		{Start: 0, End: 6},
		{Start: -1, End: -1},
		{Start: 7, End: 14},
	}
	if len(match.Captures) != len(want) {
		t.Fatalf("captures = %+v, want %+v", match.Captures, want)
	}
	for index := range want {
		if match.Captures[index] != want[index] {
			t.Fatalf("capture %d = %+v, want %+v", index, match.Captures[index], want[index])
		}
	}
}

func TestSubroutineDetectionIsSyntaxAware(t *testing.T) {
	scanner := NewScanner([]string{"(?x)a # \\g<n> is only a comment\n b"})
	if diagnostics := scanner.Diagnostics(); len(diagnostics) != 0 {
		t.Fatalf("comment produced diagnostics: %+v", diagnostics)
	}
	match := scanner.FindNextMatch(NewString("ab"), 0, FindOptionNone)
	if match == nil || match.Captures[0] != (Capture{Start: 0, End: 2}) {
		t.Fatalf("match = %+v, want [0,2)", match)
	}
}
