package textmate

import (
	"sort"
	"sync"
	"time"
	"unicode/utf8"

	"github.com/eugenioenko/textmate-go/oniguruma"
)

// Token is one contiguous span of a tokenized line. Start and End are rune
// offsets into the original line, and Scopes is ordered outermost first.
type Token struct {
	Start  int
	End    int
	Scopes []string
}

// LineResult is the result of tokenizing one line.
type LineResult struct {
	Tokens    []Token
	RuleStack *StateStack
	Stopped   bool
}

// Grammar is a compiled TextMate grammar. Compilation remains lazy: rules and
// their regular expressions are created only when the grammar is first used.
// A Grammar serializes tokenization because compiled rules contain mutable
// scanner caches.
type Grammar struct {
	mu sync.Mutex

	diagnosticsMu sync.Mutex
	diagnostics   []oniguruma.Diagnostic
	diagnosticSet map[grammarDiagnosticKey]struct{}

	rootScopeName string
	rootID        ruleID
	lastRuleID    ruleID
	rules         []rule

	includedGrammars map[string]*RawGrammar
	repository       grammarRepository
	raw              *RawGrammar
	ruleFactory      *RuleFactory

	injections      []injection
	injectionsReady bool
}

// newGrammar constructs the integration object used by Registry. The raw
// grammar is cloned by initGrammar, so compilation never annotates or extends
// caller-owned maps.
func newGrammar(
	rootScopeName string,
	raw *RawGrammar,
	repository grammarRepository,
) *Grammar {
	if raw == nil {
		return nil
	}
	if rootScopeName == "" {
		rootScopeName = raw.ScopeName
	}
	return &Grammar{
		rootScopeName:    rootScopeName,
		rules:            []rule{nil}, // rule IDs are positive and index this slice
		includedGrammars: make(map[string]*RawGrammar),
		repository:       repository,
		raw:              initGrammar(raw, nil),
		ruleFactory:      newRuleFactory(),
		diagnosticSet:    make(map[grammarDiagnosticKey]struct{}),
	}
}

// TokenizeLine tokenizes line and returns immutable state to pass to the next
// call. A nil state and InitialState both begin a new document.
func (g *Grammar) TokenizeLine(line string, prev *StateStack) LineResult {
	return g.tokenizeLine(line, prev, 0)
}

// Diagnostics returns regex translation, compilation, and match diagnostics
// from rules compiled or exercised so far. Grammars compile lazily, so callers
// should take another snapshot after tokenizing new syntax paths. Equivalent
// diagnostics produced by anchor variants or shared rule lists are returned
// once.
func (g *Grammar) Diagnostics() []oniguruma.Diagnostic {
	if g == nil {
		return nil
	}
	g.diagnosticsMu.Lock()
	defer g.diagnosticsMu.Unlock()
	return append([]oniguruma.Diagnostic(nil), g.diagnostics...)
}

type grammarDiagnosticKey struct {
	kind       oniguruma.DiagnosticKind
	pattern    string
	translated string
	message    string
}

func (g *Grammar) addRegexDiagnostic(diagnostic oniguruma.Diagnostic) {
	key := grammarDiagnosticKey{
		kind:       diagnostic.Kind,
		pattern:    diagnostic.Pattern,
		translated: diagnostic.Translated,
		message:    diagnostic.Message,
	}
	g.diagnosticsMu.Lock()
	defer g.diagnosticsMu.Unlock()
	if _, exists := g.diagnosticSet[key]; exists {
		return
	}
	g.diagnosticSet[key] = struct{}{}
	g.diagnostics = append(g.diagnostics, diagnostic)
}

func (g *Grammar) tokenizeLine(
	line string,
	prev *StateStack,
	timeLimit time.Duration,
) LineResult {
	if g == nil {
		return LineResult{}
	}

	g.mu.Lock()
	defer g.mu.Unlock()

	g.ensureRootCompiled()
	if g.rootID == 0 {
		return LineResult{}
	}

	isFirstLine := prev == nil || prev.ruleID == 0
	if isFirstLine {
		rootScopeName := "unknown"
		if rootRule := g.getRule(g.rootID); rootRule != nil {
			if name := rootRule.getName("", nil); name != "" {
				rootScopeName = name
			}
		}
		scopes := newAttributedScopeRoot(rootScopeName, g.scopeAttributes(rootScopeName))
		prev = newStateStack(nil, g.rootID, -1, -1, false, nil, scopes, scopes)
	} else {
		prev = prev.reset()
	}

	lineLength := utf8.RuneCountInString(line)
	input := oniguruma.NewString(line + "\n")
	handler := &lineTokenHandler{}
	result := tokenizeString(
		g,
		input,
		isFirstLine,
		0,
		prev,
		handler,
		true,
		timeLimit,
	)

	return LineResult{
		Tokens:    handler.result(result.stack, lineLength),
		RuleStack: result.stack,
		Stopped:   result.stoppedEarly,
	}
}

func (g *Grammar) ensureRootCompiled() {
	if g.rootID != 0 || g.raw == nil {
		return
	}
	self := g.raw.Repository["$self"]
	if self == nil {
		return
	}
	g.rootID = g.ruleFactory.getCompiledRuleID(self, g, g.raw.Repository)
	// Compile injections immediately after the root. Their IDs are therefore
	// deterministic across the registry and any worker using the same grammar.
	g.getInjections()
}

func (g *Grammar) registerRule(factory func(id ruleID) rule) rule {
	g.lastRuleID++
	id := g.lastRuleID
	g.rules = append(g.rules, nil)
	compiled := factory(id)
	g.rules[int(id)] = compiled
	return compiled
}

func (g *Grammar) getRule(id ruleID) rule {
	if id <= 0 || int(id) >= len(g.rules) {
		return nil
	}
	return g.rules[int(id)]
}

// getExternalGrammar uses the first base context with which a scope is
// requested. This is vscode-textmate's observable cache behavior and matters
// when the same external grammar is included through different repositories.
func (g *Grammar) getExternalGrammar(
	scopeName string,
	repository RawRepository,
) *RawGrammar {
	if included := g.includedGrammars[scopeName]; included != nil {
		return included
	}
	if g.repository == nil {
		return nil
	}
	raw := g.repository.lookup(scopeName)
	if raw == nil {
		return nil
	}
	var base *RawRule
	if repository != nil {
		base = repository["$base"]
	}
	included := initGrammar(raw, base)
	g.includedGrammars[scopeName] = included
	return included
}

func (g *Grammar) scopeAttributes(string) uint32 { return 0 }

func (g *Grammar) getInjections() []injection {
	if g.injectionsReady {
		return g.injections
	}
	g.injectionsReady = true

	var result []injection
	if g.raw != nil {
		selectors := make([]string, 0, len(g.raw.Injections))
		for selector := range g.raw.Injections {
			selectors = append(selectors, selector)
		}
		sort.Strings(selectors)
		for _, selector := range selectors {
			g.collectInjections(&result, selector, g.raw.Injections[selector], g.raw)
		}

		if g.repository != nil {
			for _, injectionScopeName := range g.repository.injections(g.rootScopeName) {
				external := g.getExternalGrammar(injectionScopeName, g.raw.Repository)
				if external == nil || external.InjectionSelector == "" {
					continue
				}
				g.collectInjections(
					&result,
					external.InjectionSelector,
					external.Repository["$self"],
					external,
				)
			}
		}
	}

	// SliceStable preserves selector and registry order within each priority.
	sort.SliceStable(result, func(i, j int) bool {
		return result[i].priority < result[j].priority
	})
	g.injections = result
	return g.injections
}

func (g *Grammar) collectInjections(
	result *[]injection,
	selector string,
	rawRule *RawRule,
	rawGrammar *RawGrammar,
) {
	if rawRule == nil || rawGrammar == nil {
		return
	}
	id := g.ruleFactory.getCompiledRuleID(rawRule, g, rawGrammar.Repository)
	for _, candidate := range createScopeMatchers(selector) {
		*result = append(*result, injection{
			debugSelector: selector,
			matcher:       candidate.matcher,
			priority:      candidate.priority,
			ruleID:        id,
		})
	}
}

// initGrammar returns a structurally independent grammar with synthetic
// $self and $base entries. The supplied base is already owned by this Grammar
// and intentionally remains shared across an external include boundary.
func initGrammar(raw *RawGrammar, base *RawRule) *RawGrammar {
	if raw == nil {
		return nil
	}
	grammar := cloneRawGrammar(raw)
	if grammar.Repository == nil {
		grammar.Repository = make(RawRepository)
	}
	self := &RawRule{
		Location: grammar.Location,
		Patterns: grammar.Patterns,
	}
	if grammar.ScopeName != "" {
		name := grammar.ScopeName
		self.Name = &name
	}
	grammar.Repository["$self"] = self
	if base == nil {
		base = self
	}
	grammar.Repository["$base"] = base
	return grammar
}

func cloneRawGrammar(raw *RawGrammar) *RawGrammar {
	memo := make(map[*RawRule]*RawRule)
	cloneRule := func(*RawRule) *RawRule { return nil }
	cloneRule = func(source *RawRule) *RawRule {
		if source == nil {
			return nil
		}
		if result := memo[source]; result != nil {
			return result
		}
		result := &RawRule{
			Include:             cloneString(source.Include),
			Name:                cloneString(source.Name),
			ContentName:         cloneString(source.ContentName),
			Match:               cloneString(source.Match),
			Begin:               cloneString(source.Begin),
			End:                 cloneString(source.End),
			While:               cloneString(source.While),
			ApplyEndPatternLast: source.ApplyEndPatternLast,
			Location:            cloneLocation(source.Location),
		}
		memo[source] = result
		result.Captures = cloneRawRuleMap(source.Captures, cloneRule)
		result.BeginCaptures = cloneRawRuleMap(source.BeginCaptures, cloneRule)
		result.EndCaptures = cloneRawRuleMap(source.EndCaptures, cloneRule)
		result.WhileCaptures = cloneRawRuleMap(source.WhileCaptures, cloneRule)
		result.Patterns = cloneRawRules(source.Patterns, cloneRule)
		result.Repository = cloneRawRuleMap(source.Repository, cloneRule)
		return result
	}

	result := &RawGrammar{
		ScopeName:         raw.ScopeName,
		InjectionSelector: raw.InjectionSelector,
		FileTypes:         append([]string(nil), raw.FileTypes...),
		Name:              raw.Name,
		FirstLineMatch:    raw.FirstLineMatch,
		Location:          cloneLocation(raw.Location),
	}
	result.Patterns = cloneRawRules(raw.Patterns, cloneRule)
	result.Repository = cloneRawRuleMap(raw.Repository, cloneRule)
	result.Injections = RawInjections(cloneRawRuleMap(raw.Injections, cloneRule))
	return result
}

func cloneRawRules(
	source []*RawRule,
	cloneRule func(*RawRule) *RawRule,
) []*RawRule {
	if source == nil {
		return nil
	}
	result := make([]*RawRule, len(source))
	for index, rawRule := range source {
		result[index] = cloneRule(rawRule)
	}
	return result
}

func cloneRawRuleMap[M ~map[string]*RawRule](
	source M,
	cloneRule func(*RawRule) *RawRule,
) M {
	if source == nil {
		return nil
	}
	result := make(M, len(source))
	for name, rawRule := range source {
		result[name] = cloneRule(rawRule)
	}
	return result
}

func cloneLocation(source *Location) *Location {
	if source == nil {
		return nil
	}
	result := *source
	return &result
}

type lineTokenHandler struct {
	tokens  []Token
	lastEnd int
}

func (h *lineTokenHandler) handle(scopes *attributedScopeStack, end int) {
	if end <= h.lastEnd {
		return
	}
	h.tokens = append(h.tokens, Token{
		Start:  h.lastEnd,
		End:    end,
		Scopes: append([]string(nil), scopes.scopeNames()...),
	})
	h.lastEnd = end
}

func (h *lineTokenHandler) result(stack *StateStack, lineLength int) []Token {
	result := make([]Token, 0, len(h.tokens))
	for _, token := range h.tokens {
		if token.Start >= lineLength {
			break
		}
		if token.End > lineLength {
			token.End = lineLength
		}
		if token.End > token.Start {
			result = append(result, token)
		}
	}
	if len(result) != 0 {
		return result
	}

	var scopes []string
	if stack != nil && stack.contentNameScopesList != nil {
		scopes = append(scopes, stack.contentNameScopesList.scopeNames()...)
	}
	return []Token{{Start: 0, End: lineLength, Scopes: scopes}}
}
