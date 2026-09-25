package oniguruma

import (
	"testing"
	"time"
)

func TestFindNextMatch(t *testing.T) {
	tests := []struct {
		name     string
		patterns []string
		text     string
		start    int
		want     *Match
	}{
		{
			name: "earliest beats pattern order", patterns: []string{`bar`, `foo`},
			text: "xx foo bar", want: &Match{Index: 1, Captures: []Capture{{Start: 3, End: 6}}},
		},
		{
			name: "lowest index breaks tie", patterns: []string{`f..`, `foo`},
			text: "foo", want: &Match{Index: 0, Captures: []Capture{{Start: 0, End: 3}}},
		},
		{
			name: "start position", patterns: []string{`foo`}, text: "foo foo", start: 1,
			want: &Match{Index: 0, Captures: []Capture{{Start: 4, End: 7}}},
		},
		{name: "no match", patterns: []string{`z`}, text: "abc"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := NewScanner(test.patterns).FindNextMatch(NewString(test.text), test.start, FindOptionNone)
			assertMatch(t, got, test.want)
		})
	}
}

func TestFindNextMatchWinnerBoundPreservesFullInput(t *testing.T) {
	tests := []struct {
		name     string
		patterns []string
		text     string
		want     *Match
	}{
		{
			name: "later winner consumes beyond earlier start", patterns: []string{`a`, `x.*a`}, text: "x---a",
			want: &Match{Index: 1, Captures: []Capture{{Start: 0, End: 5}}},
		},
		{
			name: "later winner looks beyond earlier start", patterns: []string{`z`, `x(?=---z)`}, text: "x---z",
			want: &Match{Index: 1, Captures: []Capture{{Start: 0, End: 1}}},
		},
		{
			name: "match at bound loses tie", patterns: []string{`z`, `(?=z)`}, text: "x---z",
			want: &Match{Index: 0, Captures: []Capture{{Start: 4, End: 5}}},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := NewScanner(test.patterns).FindNextMatch(NewString(test.text), 0, FindOptionNone)
			assertMatch(t, got, test.want)
		})
	}
}

func TestFindNextMatchBoundedMissDoesNotPoisonCache(t *testing.T) {
	scanner := NewScanner([]string{`a`, `z`})
	input := NewString("xayz")
	first := scanner.FindNextMatch(input, 0, FindOptionNone)
	assertMatch(t, first, &Match{Index: 0, Captures: []Capture{{Start: 1, End: 2}}})

	// The z search in the first call was bounded before the a match. Its miss
	// must not be remembered as an unbounded failure for this later start.
	second := scanner.FindNextMatch(input, 2, FindOptionNone)
	assertMatch(t, second, &Match{Index: 1, Captures: []Capture{{Start: 3, End: 4}}})
}

func TestStringAndRuneOffsets(t *testing.T) {
	input := NewString("a😀é")
	if input.Content() != "a😀é" || input.Len() != 3 {
		t.Fatalf("String content/length = %q/%d", input.Content(), input.Len())
	}
	match := NewScanner([]string{`(😀)(x)?(é)`}).FindNextMatch(input, 0, FindOptionNone)
	want := &Match{Index: 0, Captures: []Capture{
		{Start: 1, End: 3}, {Start: 1, End: 2}, {Start: -1, End: -1}, {Start: 2, End: 3},
	}}
	assertMatch(t, match, want)
}

func TestCaptureOrderAndRegexp2Features(t *testing.T) {
	tests := []struct {
		name    string
		pattern string
		text    string
		want    []Capture
	}{
		{
			name: "mixed named and unnamed capture order", pattern: `(?<first>a)(b)(?<third>c)`, text: "abc",
			want: []Capture{{0, 3}, {0, 1}, {1, 2}, {2, 3}},
		},
		{
			name: "lookbehind", pattern: `(?<=a)b`, text: "ab",
			want: []Capture{{1, 2}},
		},
		{
			name: "backreference", pattern: `(a)\1`, text: "aa",
			want: []Capture{{0, 2}, {0, 1}},
		},
		{
			name: "participating empty and unmatched captures", pattern: `(a*)b(c)?`, text: "b",
			want: []Capture{{0, 1}, {0, 0}, {-1, -1}},
		},
		{
			name: "repeated group returns last capture", pattern: `(a)+`, text: "aaa",
			want: []Capture{{0, 3}, {2, 3}},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			match := NewScanner([]string{test.pattern}).FindNextMatch(NewString(test.text), 0, FindOptionNone)
			if match == nil {
				t.Fatal("expected match")
			}
			assertCaptures(t, match.Captures, test.want)
		})
	}
}

func TestFindOptionsAndAnchors(t *testing.T) {
	tests := []struct {
		name    string
		pattern string
		text    string
		start   int
		opts    FindOption
		match   bool
	}{
		{name: "A allowed", pattern: `\Afoo`, text: "foo", match: true},
		{name: "A disabled", pattern: `\Afoo`, text: "foo", opts: FindOptionNotBeginString},
		{name: "G uses search origin", pattern: `\Gfoo`, text: "xfoo", start: 1, match: true},
		{name: "G disabled", pattern: `\Gfoo`, text: "xfoo", start: 1, opts: FindOptionNotBeginPosition},
		{name: "Z allowed", pattern: `x\Z`, text: "x\n", match: true},
		{name: "Z disabled", pattern: `x\Z`, text: "x\n", opts: FindOptionNotEndString},
		{name: "z allowed", pattern: `x\z`, text: "x", match: true},
		{name: "z disabled", pattern: `x\z`, text: "x", opts: FindOptionNotEndString},
		{name: "anchor escapes in classes are literals", pattern: `[\A\G\Z\z]+`, text: "AGZz", opts: FindOptionNotBeginString | FindOptionNotBeginPosition | FindOptionNotEndString, match: true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := NewScanner([]string{test.pattern}).FindNextMatch(NewString(test.text), test.start, test.opts)
			if (got != nil) != test.match {
				t.Fatalf("match = %v, want %v", got, test.match)
			}
		})
	}
	if FindOptionNotBeginString != 1 || FindOptionNotEndString != 2 || FindOptionNotBeginPosition != 4 || FindOptionDebugCall != 8 {
		t.Fatal("FindOption values no longer mirror vscode-textmate")
	}
}

func TestLineStartDoesNotMatchAfterFinalNewline(t *testing.T) {
	tests := []struct {
		name    string
		pattern string
		text    string
		start   int
		want    *Capture
	}{
		{name: "empty input", pattern: `^$`, text: "", want: &Capture{Start: 0, End: 0}},
		{name: "empty first line", pattern: `^$`, text: "\n", want: &Capture{Start: 0, End: 0}},
		{name: "final newline is not another line", pattern: `^$`, text: "last\n"},
		{name: "final newline excluded from later search", pattern: `^`, text: "last\n", start: 4},
		{name: "real empty interior line", pattern: `^$`, text: "last\n\n", want: &Capture{Start: 5, End: 5}},
		{name: "nonempty interior line", pattern: `^next`, text: "last\nnext", want: &Capture{Start: 5, End: 9}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			match := NewScanner([]string{test.pattern}).FindNextMatch(NewString(test.text), test.start, FindOptionNone)
			if test.want == nil {
				if match != nil {
					t.Fatalf("match = %+v, want nil", match)
				}
				return
			}
			if match == nil || len(match.Captures) == 0 || match.Captures[0] != *test.want {
				t.Fatalf("match = %+v, want span %+v", match, *test.want)
			}
		})
	}

	// The internal beginning-of-input assertion must not make an ordinary ^
	// sensitive to the Oniguruma option that disables only an explicit \A.
	match := NewScanner([]string{`^a`}).FindNextMatch(NewString("a"), 0, FindOptionNotBeginString)
	if match == nil || match.Captures[0] != (Capture{Start: 0, End: 1}) {
		t.Fatalf("^ with FindOptionNotBeginString = %+v, want [0,1)", match)
	}
}

func TestOnigurumaTranslationsMatch(t *testing.T) {
	tests := []struct {
		name    string
		pattern string
		text    string
		start   int
		end     int
	}{
		{name: "hex digit", pattern: `\h+`, text: "z09Af!", start: 1, end: 5},
		{name: "not hex digit", pattern: `\H+`, text: "Aff-!", start: 3, end: 5},
		{name: "BMP x brace", pattern: `\x{00E9}`, text: "é", start: 0, end: 1},
		{name: "astral x brace", pattern: `\x{1F600}`, text: "😀", start: 0, end: 1},
		{name: "Unicode alpha POSIX", pattern: `[[:alpha:]]+`, text: "1éλ!", start: 1, end: 3},
		{name: "nested POSIX identifier", pattern: `[$_[:alpha:]][$_[:alnum:]]*`, text: "9_é2!", start: 1, end: 4},
		{name: "negated ASCII POSIX", pattern: `[[:^ascii:]]+`, text: "Aéλ", start: 1, end: 3},
		{name: "mixed negated POSIX includes explicit digit", pattern: `[0[:^digit:]]`, text: "0", start: 0, end: 1},
		{name: "mixed negated POSIX matches nondigit", pattern: `[0[:^digit:]]`, text: "a", start: 0, end: 1},
		{name: "Oniguruma word property", pattern: `\p{word}+`, text: "λ_2!", start: 0, end: 3},
		{name: "Oniguruma XIDS property", pattern: `\p{XIDS}+`, text: "λ2", start: 0, end: 1},
		{name: "Oniguruma identifier properties", pattern: `\p{XIDS}\p{XIDC}*`, text: "α2_", start: 0, end: 3},
		{name: "Oniguruma blank property", pattern: `\p{blank}+`, text: "\t ", start: 0, end: 2},
		{name: "Oniguruma upper property", pattern: `\p{upper}+`, text: "aA", start: 1, end: 2},
		{name: "Oniguruma lower property", pattern: `\p{lower}+`, text: "Aa", start: 1, end: 2},
		{name: "property name normalization", pattern: `\p{Uppercase_Letter}+`, text: "aAZ", start: 1, end: 3},
		{name: "script canonicalization", pattern: `\p{greek}+`, text: "xλ", start: 1, end: 2},
		{name: "inline ignore case", pattern: `(?i:abc)`, text: "ABC", start: 0, end: 3},
		{name: "Oniguruma m is dotall", pattern: `(?m:a.b)`, text: "a\nb", start: 0, end: 3},
		{name: "extended mode comment", pattern: "(?x)a # ignored\n b", text: "ab", start: 0, end: 2},
		{name: "extended mode hash in class", pattern: `(?x)[#]+`, text: "##", start: 0, end: 2},
		{name: "global line anchors", pattern: `^b$`, text: "a\nb\nc", start: 2, end: 3},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			scanner := NewScanner([]string{test.pattern})
			if diagnostics := scanner.Diagnostics(); len(diagnostics) != 0 {
				t.Fatalf("unexpected diagnostics: %+v", diagnostics)
			}
			match := scanner.FindNextMatch(NewString(test.text), 0, FindOptionNone)
			if match == nil {
				t.Fatal("expected match")
			}
			got := match.Captures[0]
			if got.Start != test.start || got.End != test.end {
				t.Fatalf("span = %+v, want {%d %d}", got, test.start, test.end)
			}
		})
	}
}

func TestPossessiveAndIntervalTranslation(t *testing.T) {
	translationTests := []struct{ source, want string }{
		{`a*+`, `(?>a*)`},
		{`a++`, `(?>a+)`},
		{`a?+`, `(?>a?)`},
		{`[ab]++`, `(?>[ab]+)`},
		{`(?:ab)++`, `(?>(?:ab)+)`},
		{`(a++b)`, `((?>a+)b)`},
		{`a{2,3}+`, `(?:a{2,3})+`},
		{`a+{0,1}`, `(?:a+){0,1}`},
		{`a?{2}`, `(?:a?){2}`},
		{`a{,3}`, `a{0,3}`},
		{`a{3,2}`, `(?>a{2,3})`},
	}
	for _, test := range translationTests {
		got, err := translatePattern(test.source)
		if err != nil {
			t.Errorf("translatePattern(%q): %v", test.source, err)
		} else if got != test.want {
			t.Errorf("translatePattern(%q) = %q, want %q", test.source, got, test.want)
		}
	}

	if match := NewScanner([]string{`a*+a`}).FindNextMatch(NewString("aaa"), 0, FindOptionNone); match != nil {
		t.Fatalf("possessive quantifier backtracked: %+v", match)
	}
	match := NewScanner([]string{`a{2}+`}).FindNextMatch(NewString("aaaaa"), 0, FindOptionNone)
	if match == nil || match.Captures[0] != (Capture{Start: 0, End: 4}) {
		t.Fatalf("interval-plus match = %+v, want four runes", match)
	}
	match = NewScanner([]string{`a{3,2}`}).FindNextMatch(NewString("aaaa"), 0, FindOptionNone)
	if match == nil || match.Captures[0] != (Capture{Start: 0, End: 3}) {
		t.Fatalf("reversed interval match = %+v, want three runes", match)
	}
}

func TestDegradedPatternsProduceDiagnosticsAndNeverMatch(t *testing.T) {
	tests := []struct {
		name    string
		pattern string
		kind    DiagnosticKind
	}{
		{name: "compile error", pattern: `(`, kind: DiagnosticCompileError},
		{name: "subroutine by number", pattern: `(a)\g<1>`, kind: DiagnosticUnsupportedSyntax},
		{name: "duplicate named captures", pattern: `(?<n>a)(?<n>b)`, kind: DiagnosticUnsupportedSyntax},
		{name: "class intersection", pattern: `[a-z&&[^aeiou]]`, kind: DiagnosticUnsupportedSyntax},
		{name: "nested class intersection", pattern: `[[\p{S}\p{P}]&&[^()]]+`, kind: DiagnosticUnsupportedSyntax},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			scanner := NewScanner([]string{test.pattern, `a`})
			diagnostics := scanner.Diagnostics()
			if len(diagnostics) != 1 || diagnostics[0].Kind != test.kind || diagnostics[0].PatternIndex != 0 {
				t.Fatalf("diagnostics = %+v", diagnostics)
			}
			match := scanner.FindNextMatch(NewString("a"), 0, FindOptionNone)
			if match == nil || match.Index != 1 {
				t.Fatalf("degraded pattern interfered with valid one: %+v", match)
			}
		})
	}
}

func TestInlineOptionTranslationIsSyntaxAware(t *testing.T) {
	tests := []struct{ source, want string }{
		{`\(?m)`, `\(?m)`},
		{`[(?m)]`, `[(?m)]`},
		{`(?# (?m))(?m:a)`, `(?# (?m))(?s:a)`},
		{"(?x)# (?m)\n(?m:a)", "(?x)# (?m)\n(?s:a)"},
	}
	for _, test := range tests {
		if got := translateInlineOptions(test.source); got != test.want {
			t.Errorf("translateInlineOptions(%q) = %q, want %q", test.source, got, test.want)
		}
	}
	pattern := "(?x)a # ) b++ \\g<n>\n c++"
	translated, err := translatePattern(translateInlineOptions(pattern))
	if err != nil {
		t.Fatalf("extended-mode comment triggered compatibility error: %v", err)
	}
	if want := "(?x)a # ) b++ \\g<n>\n (?>c+)"; translated != want {
		t.Fatalf("extended-mode translation = %q, want %q", translated, want)
	}

	for _, pattern := range []string{
		"(?x)# [a&&b] (?<name>ignored)\n(?<name>a)",
		"(?x:(?# [a&&b])# (?<name>ignored)\na)",
	} {
		scanner := NewScanner([]string{pattern})
		if diagnostics := scanner.Diagnostics(); len(diagnostics) != 0 {
			t.Errorf("comment-only constructs in %q produced diagnostics: %+v", pattern, diagnostics)
			continue
		}
		if match := scanner.FindNextMatch(NewString("a"), 0, FindOptionNone); match == nil {
			t.Errorf("comment-safe pattern %q did not match", pattern)
		}
	}
}

func TestOrdinaryDotDoesNotMatchNewline(t *testing.T) {
	if match := NewScanner([]string{`a.b`}).FindNextMatch(NewString("a\nb"), 0, FindOptionNone); match != nil {
		t.Fatalf("ordinary dot unexpectedly matched newline: %+v", match)
	}
}

func TestPerPatternTimeoutIsRecordedAndScanningContinues(t *testing.T) {
	// regexp2's clock-period test hook is process-global and unsafe to mutate
	// once another timed match has started its background clock.
	scanner := NewScanner([]string{`^(a|aa)+$`, `X`}, WithMatchTimeout(time.Millisecond))
	input := NewString("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaX")
	match := scanner.FindNextMatch(input, 0, FindOptionNone)
	if match == nil || match.Index != 1 {
		t.Fatalf("scanner did not continue after pathological pattern: %+v", match)
	}
	diagnostics := scanner.Diagnostics()
	if len(diagnostics) != 1 || diagnostics[0].Kind != DiagnosticMatchTimeout {
		t.Fatalf("timeout diagnostics = %+v", diagnostics)
	}
	if diagnostics[0].Message == "" || diagnostics[0].Message == input.Content() {
		t.Fatalf("timeout diagnostic was empty or leaked full input: %+v", diagnostics[0])
	}

	_ = scanner.FindNextMatch(input, 0, FindOptionNone)
	if got := len(scanner.Diagnostics()); got != 1 {
		t.Fatalf("timeout diagnostic count = %d, want deduplicated count 1", got)
	}
}

func TestSearchBounds(t *testing.T) {
	scanner := NewScanner([]string{`.`})
	if match := scanner.FindNextMatch(NewString("a"), -1, FindOptionNone); match == nil || match.Captures[0] != (Capture{Start: 0, End: 1}) {
		t.Errorf("negative start was not clamped: %+v", match)
	}
	if match := scanner.FindNextMatch(NewString("a"), 2, FindOptionNone); match != nil {
		t.Errorf("out-of-range start returned %+v", match)
	}
	if match := scanner.FindNextMatch(nil, 0, FindOptionNone); match != nil {
		t.Errorf("nil input returned %+v", match)
	}
	match := NewScanner([]string{`$`}).FindNextMatch(NewString("a"), 1, FindOptionNone)
	if match == nil || match.Captures[0] != (Capture{Start: 1, End: 1}) {
		t.Errorf("zero-width end match = %+v", match)
	}
}

func assertMatch(t *testing.T, got, want *Match) {
	t.Helper()
	if got == nil || want == nil {
		if got != want {
			t.Fatalf("match = %+v, want %+v", got, want)
		}
		return
	}
	if got.Index != want.Index {
		t.Fatalf("pattern index = %d, want %d", got.Index, want.Index)
	}
	assertCaptures(t, got.Captures, want.Captures)
}

func assertCaptures(t *testing.T, got, want []Capture) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("captures = %+v, want %+v", got, want)
	}
	for index := range got {
		if got[index] != want[index] {
			t.Fatalf("capture %d = %+v, want %+v", index, got[index], want[index])
		}
	}
}
