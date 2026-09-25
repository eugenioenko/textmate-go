package oniguruma

import "testing"

func TestEscapedHyphenCanStartCharacterClassRange(t *testing.T) {
	scanner := NewScanner([]string{`[\--9]+`})
	if diagnostics := scanner.Diagnostics(); len(diagnostics) != 0 {
		t.Fatalf("compile diagnostics: %+v", diagnostics)
	}
	for _, text := range []string{"-", ".", "/", "0", "1", "9"} {
		if match := scanner.FindNextMatch(NewString(text), 0, FindOptionNone); match == nil {
			t.Errorf("%q did not match", text)
		}
	}
	for _, text := range []string{",", ":"} {
		if match := scanner.FindNextMatch(NewString(text), 0, FindOptionNone); match != nil {
			t.Errorf("%q unexpectedly matched: %+v", text, match)
		}
	}
}
