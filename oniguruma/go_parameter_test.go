package oniguruma

import "testing"

func TestGoFunctionParameterCapturesNameAndType(t *testing.T) {
	const pattern = `(\w+\s+)?([*.\w]+(?:\[(?:[*.\w]+(?:,\s+)?)+{0,1}])?)`
	scanner := NewScanner([]string{pattern})
	if diagnostics := scanner.Diagnostics(); len(diagnostics) != 0 {
		t.Fatalf("compile diagnostics: %+v", diagnostics)
	}
	match := scanner.FindNextMatch(NewString("name string"), 0, FindOptionNone)
	if match == nil {
		t.Fatal("parameter did not match")
	}
	want := []Capture{{Start: 0, End: 11}, {Start: 0, End: 5}, {Start: 5, End: 11}}
	if len(match.Captures) != len(want) {
		t.Fatalf("captures = %+v, want %+v", match.Captures, want)
	}
	for index := range want {
		if match.Captures[index] != want[index] {
			t.Fatalf("capture %d = %+v, want %+v", index, match.Captures[index], want[index])
		}
	}
}
