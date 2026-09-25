package textmate

// TokenCategory is a coarse semantic classification derived from TextMate
// scope names. It is intentionally independent of colors and font styles.
// Applications can map categories to their own presentation or use a full
// TextMate theme resolver when selector-accurate styling is required.
type TokenCategory uint8

const (
	TokenCategoryPlain TokenCategory = iota
	TokenCategoryComment
	TokenCategoryString
	TokenCategoryRegexp
	TokenCategoryNumber
	TokenCategoryKeyword
	TokenCategoryOperator
	TokenCategoryFunction
	TokenCategoryTag
	TokenCategoryAttribute
	TokenCategoryType
	TokenCategoryBuiltin
	TokenCategoryVariable
	TokenCategoryHeading
	TokenCategoryBold
	TokenCategoryItalic
	TokenCategoryQuote
	TokenCategoryInserted
	TokenCategoryDeleted
	TokenCategoryPunctuation
	TokenCategoryInvalid
)

// String returns the stable, lower-case name of category.
func (category TokenCategory) String() string {
	switch category {
	case TokenCategoryPlain:
		return "plain"
	case TokenCategoryComment:
		return "comment"
	case TokenCategoryString:
		return "string"
	case TokenCategoryRegexp:
		return "regexp"
	case TokenCategoryNumber:
		return "number"
	case TokenCategoryKeyword:
		return "keyword"
	case TokenCategoryOperator:
		return "operator"
	case TokenCategoryFunction:
		return "function"
	case TokenCategoryTag:
		return "tag"
	case TokenCategoryAttribute:
		return "attribute"
	case TokenCategoryType:
		return "type"
	case TokenCategoryBuiltin:
		return "builtin"
	case TokenCategoryVariable:
		return "variable"
	case TokenCategoryHeading:
		return "heading"
	case TokenCategoryBold:
		return "bold"
	case TokenCategoryItalic:
		return "italic"
	case TokenCategoryQuote:
		return "quote"
	case TokenCategoryInserted:
		return "inserted"
	case TokenCategoryDeleted:
		return "deleted"
	case TokenCategoryPunctuation:
		return "punctuation"
	case TokenCategoryInvalid:
		return "invalid"
	default:
		return "unknown"
	}
}

type categoryRule struct {
	prefix   string
	category TokenCategory
}

// categoryRules is deliberately a compact semantic heuristic, not a TextMate
// selector or theme implementation. classifyScope chooses the longest matching
// prefix, so this table does not depend on declaration order.
var categoryRules = [...]categoryRule{
	{prefix: "comment", category: TokenCategoryComment},
	{prefix: "punctuation.definition.comment", category: TokenCategoryComment},
	{prefix: "string.regexp", category: TokenCategoryRegexp},
	{prefix: "constant.other.character-class.regexp", category: TokenCategoryRegexp},
	{prefix: "keyword.control.anchor.regexp", category: TokenCategoryRegexp},
	{prefix: "keyword.operator.quantifier.regexp", category: TokenCategoryRegexp},
	{prefix: "string", category: TokenCategoryString},
	{prefix: "punctuation.definition.string", category: TokenCategoryString},
	{prefix: "constant.character.escape", category: TokenCategoryString},
	{prefix: "markup.inline.raw", category: TokenCategoryString},
	{prefix: "markup.raw", category: TokenCategoryString},
	{prefix: "constant.numeric", category: TokenCategoryNumber},
	{prefix: "constant.language", category: TokenCategoryKeyword},
	{prefix: "constant", category: TokenCategoryNumber},
	{prefix: "keyword.operator", category: TokenCategoryOperator},
	{prefix: "keyword", category: TokenCategoryKeyword},
	{prefix: "storage.type", category: TokenCategoryKeyword},
	{prefix: "storage.modifier", category: TokenCategoryKeyword},
	{prefix: "storage", category: TokenCategoryKeyword},
	{prefix: "entity.name.function", category: TokenCategoryFunction},
	{prefix: "support.function", category: TokenCategoryFunction},
	{prefix: "entity.name.tag", category: TokenCategoryTag},
	{prefix: "entity.other.attribute-name", category: TokenCategoryAttribute},
	{prefix: "support.type.property-name", category: TokenCategoryAttribute},
	{prefix: "entity.name.type", category: TokenCategoryType},
	{prefix: "entity.name.class", category: TokenCategoryType},
	{prefix: "entity.other.inherited-class", category: TokenCategoryType},
	{prefix: "entity.name", category: TokenCategoryType},
	{prefix: "support.type", category: TokenCategoryType},
	{prefix: "support.class", category: TokenCategoryType},
	{prefix: "support.constant", category: TokenCategoryNumber},
	{prefix: "support.variable", category: TokenCategoryBuiltin},
	{prefix: "support", category: TokenCategoryBuiltin},
	{prefix: "variable.language", category: TokenCategoryKeyword},
	{prefix: "variable", category: TokenCategoryVariable},
	{prefix: "entity.name.section", category: TokenCategoryHeading},
	{prefix: "markup.heading", category: TokenCategoryHeading},
	{prefix: "markup.bold", category: TokenCategoryBold},
	{prefix: "markup.italic", category: TokenCategoryItalic},
	{prefix: "markup.quote", category: TokenCategoryQuote},
	{prefix: "markup.inserted", category: TokenCategoryInserted},
	{prefix: "markup.deleted", category: TokenCategoryDeleted},
	{prefix: "punctuation", category: TokenCategoryPunctuation},
	{prefix: "invalid", category: TokenCategoryInvalid},

	// A matched plain category is a deliberate reset. It prevents an outer
	// string or markup scope from coloring embedded language tokens that have
	// no more-specific recognized category of their own.
	{prefix: "meta.embedded", category: TokenCategoryPlain},
	{prefix: "meta.template.expression", category: TokenCategoryPlain},
}

type categoryMatch struct {
	category TokenCategory
	matched  bool
}

// ClassifyScopes returns a coarse category for a TextMate scope stack ordered
// outermost first. The innermost recognized scope wins. Within one scope, the
// longest prefix at a dot boundary wins. Plain is returned when no scope is
// recognized and for deliberate reset scopes such as meta.embedded.
//
// Classification is a heuristic semantic layer, not TextMate theme
// resolution. The function does not retain or modify scopes and allocates no
// memory.
func ClassifyScopes(scopes []string) TokenCategory {
	for index := len(scopes) - 1; index >= 0; index-- {
		if match := classifyScope(scopes[index]); match.matched {
			return match.category
		}
	}
	return TokenCategoryPlain
}

func classifyScope(scope string) categoryMatch {
	bestLength := -1
	best := TokenCategoryPlain
	for _, rule := range categoryRules {
		if len(rule.prefix) <= bestLength || !scopePrefixMatches(scope, rule.prefix) {
			continue
		}
		bestLength = len(rule.prefix)
		best = rule.category
	}
	return categoryMatch{category: best, matched: bestLength >= 0}
}

func scopePrefixMatches(scope, prefix string) bool {
	if len(scope) < len(prefix) || scope[:len(prefix)] != prefix {
		return false
	}
	return len(scope) == len(prefix) || scope[len(prefix)] == '.'
}

// Category returns the cached coarse category of s. A nil ScopeStack is
// plain. The first call computes the result; subsequent calls are lock-free
// and allocate no memory.
func (s *ScopeStack) Category() TokenCategory {
	if s == nil {
		return TokenCategoryPlain
	}
	if cached := s.category.Load(); cached != 0 {
		return TokenCategory(cached - 1)
	}
	category := ClassifyScopes(s.names)
	s.category.CompareAndSwap(0, uint32(category)+1)
	return category
}

// Category returns the token's coarse semantic category. It uses the cached
// ScopeStack classification when available and falls back to Scopes for tokens
// constructed by callers or older integrations.
func (token Token) Category() TokenCategory {
	if token.ScopeStack != nil {
		return token.ScopeStack.Category()
	}
	return ClassifyScopes(token.Scopes)
}
