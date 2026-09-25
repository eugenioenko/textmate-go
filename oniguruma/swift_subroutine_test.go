package oniguruma

import "testing"

func TestSwiftRightRecursiveNumericSubroutines(t *testing.T) {
	const numericLiteral = `(?x)
		(?<decimal-literal>            \g<decimal-digit> \g<decimal-characters>? ){0}
		(?<decimal-digit>              \d ){0}
		(?<decimal-character>          \g<decimal-digit> | _ ){0}
		(?<decimal-characters>         \g<decimal-character> \g<decimal-characters>? ){0}
		\g<decimal-literal>`

	scanner := NewScanner([]string{numericLiteral})
	if diagnostics := scanner.Diagnostics(); len(diagnostics) != 0 {
		t.Fatalf("compile diagnostics: %+v", diagnostics)
	}
	match := scanner.FindNextMatch(NewString("0"), 0, FindOptionNone)
	if match == nil || match.Captures[0] != (Capture{Start: 0, End: 1}) {
		t.Fatalf("match = %#v, want [0,1)", match)
	}
}
