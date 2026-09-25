package textmate

import (
	"sync"
	"testing"
)

func TestTokenCategoryRules(t *testing.T) {
	want := map[string]TokenCategory{
		"comment":                               TokenCategoryComment,
		"punctuation.definition.comment":        TokenCategoryComment,
		"string.regexp":                         TokenCategoryRegexp,
		"constant.other.character-class.regexp": TokenCategoryRegexp,
		"keyword.control.anchor.regexp":         TokenCategoryRegexp,
		"keyword.operator.quantifier.regexp":    TokenCategoryRegexp,
		"string":                                TokenCategoryString,
		"punctuation.definition.string":         TokenCategoryString,
		"constant.character.escape":             TokenCategoryString,
		"markup.inline.raw":                     TokenCategoryString,
		"markup.raw":                            TokenCategoryString,
		"constant.numeric":                      TokenCategoryNumber,
		"constant.language":                     TokenCategoryKeyword,
		"constant":                              TokenCategoryNumber,
		"keyword.operator":                      TokenCategoryOperator,
		"keyword":                               TokenCategoryKeyword,
		"storage.type":                          TokenCategoryKeyword,
		"storage.modifier":                      TokenCategoryKeyword,
		"storage":                               TokenCategoryKeyword,
		"entity.name.function":                  TokenCategoryFunction,
		"support.function":                      TokenCategoryFunction,
		"entity.name.tag":                       TokenCategoryTag,
		"entity.other.attribute-name":           TokenCategoryAttribute,
		"support.type.property-name":            TokenCategoryAttribute,
		"entity.name.type":                      TokenCategoryType,
		"entity.name.class":                     TokenCategoryType,
		"entity.other.inherited-class":          TokenCategoryType,
		"entity.name":                           TokenCategoryType,
		"support.type":                          TokenCategoryType,
		"support.class":                         TokenCategoryType,
		"support.constant":                      TokenCategoryNumber,
		"support.variable":                      TokenCategoryBuiltin,
		"support":                               TokenCategoryBuiltin,
		"variable.language":                     TokenCategoryKeyword,
		"variable":                              TokenCategoryVariable,
		"entity.name.section":                   TokenCategoryHeading,
		"markup.heading":                        TokenCategoryHeading,
		"markup.bold":                           TokenCategoryBold,
		"markup.italic":                         TokenCategoryItalic,
		"markup.quote":                          TokenCategoryQuote,
		"markup.inserted":                       TokenCategoryInserted,
		"markup.deleted":                        TokenCategoryDeleted,
		"punctuation":                           TokenCategoryPunctuation,
		"invalid":                               TokenCategoryInvalid,
		"meta.embedded":                         TokenCategoryPlain,
		"meta.template.expression":              TokenCategoryPlain,
	}

	if len(want) != len(categoryRules) {
		t.Fatalf("test covers %d prefixes, rules contain %d", len(want), len(categoryRules))
	}
	seen := make(map[string]struct{}, len(categoryRules))
	for _, rule := range categoryRules {
		if _, exists := seen[rule.prefix]; exists {
			t.Errorf("duplicate category prefix %q", rule.prefix)
		}
		seen[rule.prefix] = struct{}{}
		category, exists := want[rule.prefix]
		if !exists {
			t.Errorf("category prefix %q is missing from exhaustive test", rule.prefix)
			continue
		}
		for _, scope := range []string{rule.prefix, rule.prefix + ".go"} {
			if got := ClassifyScopes([]string{scope}); got != category {
				t.Errorf("ClassifyScopes(%q) = %v, want %v", scope, got, category)
			}
		}
	}
}

func TestTokenCategoryNames(t *testing.T) {
	want := []string{
		"plain", "comment", "string", "regexp", "number", "keyword",
		"operator", "function", "tag", "attribute", "type", "builtin",
		"variable", "heading", "bold", "italic", "quote", "inserted",
		"deleted", "punctuation", "invalid",
	}
	for category, name := range want {
		if got := TokenCategory(category).String(); got != name {
			t.Errorf("TokenCategory(%d).String() = %q, want %q", category, got, name)
		}
	}
	if got := TokenCategory(255).String(); got != "unknown" {
		t.Fatalf("unknown category String() = %q", got)
	}
}

func TestClassifyScopesPrecedenceAndBoundaries(t *testing.T) {
	tests := []struct {
		name   string
		scopes []string
		want   TokenCategory
	}{
		{name: "empty", want: TokenCategoryPlain},
		{name: "unknown", scopes: []string{"source.go", "meta.block.go"}, want: TokenCategoryPlain},
		{name: "dot boundary", scopes: []string{"commentary.go"}, want: TokenCategoryPlain},
		{name: "longest prefix", scopes: []string{"string.regexp.go"}, want: TokenCategoryRegexp},
		{name: "property before type", scopes: []string{"support.type.property-name.css"}, want: TokenCategoryAttribute},
		{name: "function before support", scopes: []string{"support.function.builtin.go"}, want: TokenCategoryFunction},
		{name: "operator before keyword", scopes: []string{"keyword.operator.assignment.go"}, want: TokenCategoryOperator},
		{name: "innermost recognized", scopes: []string{"comment.block", "meta.block", "variable.other"}, want: TokenCategoryVariable},
		{name: "skip unknown inner", scopes: []string{"comment.block", "meta.block", "source.go"}, want: TokenCategoryComment},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := ClassifyScopes(test.scopes); got != test.want {
				t.Fatalf("ClassifyScopes(%q) = %v, want %v", test.scopes, got, test.want)
			}
		})
	}
}

func TestClassifyScopesPlainResetStopsOuterCategory(t *testing.T) {
	for _, reset := range []string{"meta.embedded.block.go", "meta.template.expression.js"} {
		scopes := []string{"source.js", "string.quoted.template.js", reset, "source.go"}
		if got := ClassifyScopes(scopes); got != TokenCategoryPlain {
			t.Errorf("ClassifyScopes(%q) = %v, want plain reset", scopes, got)
		}
	}

	scopes := []string{
		"source.js",
		"string.quoted.template.js",
		"meta.embedded.block.go",
		"source.go",
		"keyword.control.go",
	}
	if got := ClassifyScopes(scopes); got != TokenCategoryKeyword {
		t.Fatalf("recognized scope inside reset = %v, want keyword", got)
	}
}

func TestScopeStackAndTokenCategoryCache(t *testing.T) {
	stack := &ScopeStack{names: []string{"source.go", "comment.line.go"}}
	if cached := stack.category.Load(); cached != 0 {
		t.Fatalf("new cache = %d, want zero", cached)
	}
	if got := stack.Category(); got != TokenCategoryComment {
		t.Fatalf("ScopeStack.Category() = %v, want comment", got)
	}
	if cached := stack.category.Load(); cached != uint32(TokenCategoryComment)+1 {
		t.Fatalf("populated cache = %d", cached)
	}
	if allocations := testing.AllocsPerRun(1_000, func() {
		benchmarkTokenCategory = stack.Category()
	}); allocations != 0 {
		t.Fatalf("cached classification allocations = %.2f, want zero", allocations)
	}

	token := Token{Scopes: []string{"string.quoted.go"}, ScopeStack: stack}
	if got := token.Category(); got != TokenCategoryComment {
		t.Fatalf("Token.Category() did not prefer ScopeStack: got %v", got)
	}
	legacy := Token{Scopes: []string{"source.go", "constant.numeric.go"}}
	if got := legacy.Category(); got != TokenCategoryNumber {
		t.Fatalf("legacy Token.Category() = %v, want number", got)
	}
	if got := (*ScopeStack)(nil).Category(); got != TokenCategoryPlain {
		t.Fatalf("nil ScopeStack.Category() = %v, want plain", got)
	}
}

func TestScopeStackCategoryIsConcurrent(t *testing.T) {
	stack := &ScopeStack{names: []string{"source.ts", "meta.function.ts", "entity.name.function.ts"}}
	const goroutines = 32
	var wait sync.WaitGroup
	wait.Add(goroutines)
	for range goroutines {
		go func() {
			defer wait.Done()
			for range 1_000 {
				if got := stack.Category(); got != TokenCategoryFunction {
					t.Errorf("concurrent Category() = %v, want function", got)
					return
				}
			}
		}()
	}
	wait.Wait()
}

var benchmarkTokenCategory TokenCategory

func BenchmarkTokenCategory(b *testing.B) {
	scopes := []string{
		"source.tsx",
		"meta.function.tsx",
		"meta.block.tsx",
		"variable.other.readwrite.tsx",
	}
	stack := &ScopeStack{names: scopes}
	stack.Category()

	b.Run("scope-stack-cached", func(b *testing.B) {
		b.ReportAllocs()
		for range b.N {
			benchmarkTokenCategory = stack.Category()
		}
	})
	b.Run("raw-scopes", func(b *testing.B) {
		b.ReportAllocs()
		for range b.N {
			benchmarkTokenCategory = ClassifyScopes(scopes)
		}
	})
}
