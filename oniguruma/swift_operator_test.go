package oniguruma

import "testing"

func TestSwiftPostfixOperatorWinsBeforeGenericInfix(t *testing.T) {
	const operatorClass = `[-!%\&*+/<-?^|~¡-§©«¬®°±¶»¿×÷̀-ͯ᷀-᷿‖‗†-‧‰-‾⁁-⁓⁕-⁞⃐-⃿←-⏿─-❵➔-⯿⸀-⹿、。〃〈-〰︀-️︠-︯\x{E0100}-\x{E01EF}]`
	postfix := `\G(?<!^|[(,:;\[{\s])((?!(//|/\*|\*/))(` + operatorClass + `))++(?=[]),:;}\s]|\z)`
	infix := `\G((?!(//|/\*|\*/))(` + operatorClass + `))++`
	scanner := NewScanner([]string{postfix, infix})
	if diagnostics := scanner.Diagnostics(); len(diagnostics) != 0 {
		t.Fatalf("compile diagnostics: %+v", diagnostics)
	}
	match := scanner.FindNextMatch(NewString("x=U>\n"), 3, FindOptionNone)
	if match == nil || match.Index != 0 || match.Captures[0] != (Capture{Start: 3, End: 4}) {
		t.Fatalf("match = %+v, want postfix pattern at [3,4)", match)
	}
}

func TestGAnchorRetriesWhenSearchStartAdvances(t *testing.T) {
	scanner := NewScanner([]string{`\G>`})
	input := NewString("x=U>")
	if match := scanner.FindNextMatch(input, 1, FindOptionNone); match != nil {
		t.Fatalf("match at 1 = %+v, want nil", match)
	}
	match := scanner.FindNextMatch(input, 3, FindOptionNone)
	if match == nil || match.Captures[0] != (Capture{Start: 3, End: 4}) {
		t.Fatalf("match at 3 = %+v, want [3,4)", match)
	}
}
