package oniguruma

import "testing"

func TestClojureKeywordLookaround(t *testing.T) {
	const pattern = `(?<=(\s|\(|\[|\{)):[a-zA-Z0-9\#\.\-\_\:\+\=\>\<\/\!\?\*]+(?=(\s|\)|\]|\}))`
	scanner := NewScanner([]string{pattern})
	if diagnostics := scanner.Diagnostics(); len(diagnostics) != 0 {
		t.Fatalf("unexpected diagnostics: %#v", diagnostics)
	}
	match := scanner.FindNextMatch(NewString("[:two]"), 0, FindOptionNone)
	if match == nil {
		t.Fatal("Clojure keyword did not match")
	}
	if got, want := match.Captures[0], (Capture{Start: 1, End: 5}); got != want {
		t.Fatalf("match = %#v, want %#v", got, want)
	}
}
