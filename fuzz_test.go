package textmate

import (
	"strings"
	"testing"
	"unicode/utf8"
)

// FuzzTokenizeLine exercises the public line-by-line API with carried state.
// Besides panics, it catches gaps, overlaps, and byte/rune offset confusion in
// the returned token stream.
func FuzzTokenizeLine(f *testing.F) {
	for _, seed := range []string{
		"",
		"plain text",
		"before /* open\nstill commented */ after",
		"const message = \"hello 🙂\";",
		"/* 日本語 */\nidentifier_2",
		"\xff\xfeinvalid UTF-8",
	} {
		f.Add(seed)
	}

	value := func(text string) *string { return &text }
	raw := &RawGrammar{
		ScopeName: "source.fuzz",
		Patterns: []*RawRule{
			{
				Begin: value(`/\*`),
				End:   value(`\*/`),
				Name:  value("comment.block.fuzz"),
			},
			{
				Begin: value(`"`),
				End:   value(`"`),
				Name:  value("string.quoted.double.fuzz"),
				Patterns: []*RawRule{{
					Match: value(`\\.`),
					Name:  value("constant.character.escape.fuzz"),
				}},
			},
			{
				Match: value(`\b[[:alpha:]_][[:alnum:]_]*\b`),
				Name:  value("identifier.fuzz"),
			},
		},
	}
	grammar := newGrammar(raw.ScopeName, raw, nil)

	f.Fuzz(func(t *testing.T, contents string) {
		if len(contents) > 32<<10 {
			t.Skip("bounded fuzz input")
		}

		var state *StateStack
		for lineNumber, line := range strings.Split(contents, "\n") {
			result := grammar.TokenizeLine(line, state)
			if result.RuleStack == nil {
				t.Fatalf("line %d returned nil state", lineNumber+1)
			}
			if result.Stopped {
				t.Fatalf("line %d stopped without a public time limit", lineNumber+1)
			}

			lineLength := utf8.RuneCountInString(line)
			position := 0
			for tokenIndex, token := range result.Tokens {
				if token.Start != position {
					t.Fatalf(
						"line %d token %d starts at %d after %d: %#v",
						lineNumber+1, tokenIndex, token.Start, position, result.Tokens,
					)
				}
				if token.End < token.Start || token.End > lineLength {
					t.Fatalf(
						"line %d token %d has invalid rune span [%d,%d) for length %d",
						lineNumber+1, tokenIndex, token.Start, token.End, lineLength,
					)
				}
				if len(token.Scopes) == 0 {
					t.Fatalf("line %d token %d has no scopes", lineNumber+1, tokenIndex)
				}
				position = token.End
			}
			if position != lineLength {
				t.Fatalf(
					"line %d tokens end at %d, want rune length %d: %#v",
					lineNumber+1, position, lineLength, result.Tokens,
				)
			}
			state = result.RuleStack
		}
	})
}
