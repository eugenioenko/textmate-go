package oniguruma

import "testing"

func TestPythonContinuationEndMatchesAtStartOfCodeLine(t *testing.T) {
	const end = `(?x)
  (?=^\s*$)
  |
  (?! (\s* [rR]? (\'\'\'|\"\"\"|\'|\"))
      |
      (\G $)  (?# '\G' is necessary for ST)
  )
`
	scanner := NewScanner([]string{end})
	if diagnostics := scanner.Diagnostics(); len(diagnostics) != 0 {
		t.Fatalf("compile diagnostics: %+v", diagnostics)
	}
	match := scanner.FindNextMatch(NewString("   and 1 <= day <= 31 and 0 <= hour < 24 \\\n"), 0, FindOptionNone)
	if match == nil || match.Captures[0].Start != 0 || match.Captures[0].End != 0 {
		t.Fatalf("match = %#v, want empty match at start", match)
	}
}
