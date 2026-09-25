package textmate

import "regexp"

type scopeMatcher func([]string) bool

type scopeMatcherWithPriority struct {
	matcher  scopeMatcher
	priority int
}

var selectorTokenPattern = regexp.MustCompile(`([LR]:|[\w\.:][\w\.:\-]*|[,|\-()])`)

func createScopeMatchers(selector string) []scopeMatcherWithPriority {
	tokens := selectorTokenPattern.FindAllString(selector, -1)
	position := 0
	next := func() string {
		if position >= len(tokens) {
			return ""
		}
		result := tokens[position]
		position++
		return result
	}
	token := next()
	result := make([]scopeMatcherWithPriority, 0, 1)

	var parseOperand func() scopeMatcher
	var parseConjunction func() scopeMatcher
	var parseInnerExpression func() scopeMatcher

	parseOperand = func() scopeMatcher {
		if token == "-" {
			token = next()
			negated := parseOperand()
			return func(scopes []string) bool {
				return negated != nil && !negated(scopes)
			}
		}
		if token == "(" {
			token = next()
			inner := parseInnerExpression()
			if token == ")" {
				token = next()
			}
			return inner
		}
		if isSelectorIdentifier(token) {
			identifiers := make([]string, 0, 1)
			for isSelectorIdentifier(token) {
				identifiers = append(identifiers, token)
				token = next()
			}
			return func(scopes []string) bool {
				return matchScopeIdentifiers(identifiers, scopes)
			}
		}
		return nil
	}

	parseConjunction = func() scopeMatcher {
		matchers := make([]scopeMatcher, 0, 1)
		for matcher := parseOperand(); matcher != nil; matcher = parseOperand() {
			matchers = append(matchers, matcher)
		}
		return func(scopes []string) bool {
			for _, matcher := range matchers {
				if !matcher(scopes) {
					return false
				}
			}
			return true
		}
	}

	parseInnerExpression = func() scopeMatcher {
		matchers := make([]scopeMatcher, 0, 1)
		for {
			matchers = append(matchers, parseConjunction())
			if token != "|" && token != "," {
				break
			}
			for token == "|" || token == "," {
				token = next()
			}
		}
		return func(scopes []string) bool {
			for _, matcher := range matchers {
				if matcher(scopes) {
					return true
				}
			}
			return false
		}
	}

	for token != "" {
		priority := 0
		if len(token) == 2 && token[1] == ':' {
			switch token[0] {
			case 'L':
				priority = -1
			case 'R':
				priority = 1
			}
			token = next()
		}
		result = append(result, scopeMatcherWithPriority{
			matcher:  parseConjunction(),
			priority: priority,
		})
		if token != "," {
			break
		}
		token = next()
	}
	return result
}

func isSelectorIdentifier(token string) bool {
	if token == "" {
		return false
	}
	first := token[0]
	return first == '_' || first == ':' || first == '.' ||
		(first >= '0' && first <= '9') ||
		(first >= 'A' && first <= 'Z') ||
		(first >= 'a' && first <= 'z')
}

func matchScopeIdentifiers(identifiers, scopes []string) bool {
	lastIndex := 0
	for _, identifier := range identifiers {
		found := false
		for index := lastIndex; index < len(scopes); index++ {
			if scopesMatch(scopes[index], identifier) {
				lastIndex = index + 1
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

func scopesMatch(scope, selector string) bool {
	if scope == selector {
		return true
	}
	return len(scope) > len(selector) && scope[:len(selector)] == selector &&
		scope[len(selector)] == '.'
}
