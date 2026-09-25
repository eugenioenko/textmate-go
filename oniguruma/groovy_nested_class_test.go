package oniguruma

import "testing"

func TestGroovyGenericEndNestedNegatedClassMatchesBrace(t *testing.T) {
	const pattern = `[>[^],<?\[\w\s]]`
	scanner := NewScanner([]string{pattern})
	if diagnostics := scanner.Diagnostics(); len(diagnostics) != 0 {
		t.Fatalf("compile diagnostics: %+v", diagnostics)
	}
	match := scanner.FindNextMatch(NewString("{ "), 0, FindOptionNone)
	if match == nil || match.Captures[0] != (Capture{Start: 0, End: 1}) {
		t.Fatalf("match = %+v, want brace at [0,1)", match)
	}
}
