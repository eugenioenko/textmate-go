package textmate

import (
	"strconv"
	"strings"

	"github.com/eugenioenko/textmate-go/oniguruma"
)

// ruleRegistry is the portion of Grammar needed while flattening rule
// patterns. Rule implementations deliberately depend on this small interface
// so their compilation and caching can be tested independently.
type ruleRegistry interface {
	getRule(id ruleID) rule
}

// ruleFactoryHelper supplies rule IDs and external grammars. registerRule must
// choose the ID before invoking factory: getCompiledRuleID records that ID
// before compiling children, which is what makes recursive repositories
// finite.
type ruleFactoryHelper interface {
	ruleRegistry
	registerRule(factory func(id ruleID) rule) rule
	getExternalGrammar(scopeName string, repository RawRepository) *RawGrammar
}

type regexDiagnosticCollector interface {
	addRegexDiagnostic(oniguruma.Diagnostic)
}

func newRuleRegExpSourceList(grammar ruleRegistry) *regexpSourceList {
	if collector, ok := grammar.(regexDiagnosticCollector); ok {
		return newRegExpSourceList(collector.addRegexDiagnostic)
	}
	return newRegExpSourceList()
}

type rule interface {
	dispose()
	getID() ruleID
	getName(lineText string, captureIndices []oniguruma.Capture) string
	getContentName(lineText string, captureIndices []oniguruma.Capture) string
	collectPatterns(grammar ruleRegistry, out *regexpSourceList)
	compile(grammar ruleRegistry, endRegexSource *string) *compiledRule
	compileAG(grammar ruleRegistry, endRegexSource *string, allowA, allowG bool) *compiledRule
}

type baseRule struct {
	location *Location
	id       ruleID

	nameIsCapturing        bool
	name                   string
	contentNameIsCapturing bool
	contentName            string
}

func newBaseRule(location *Location, id ruleID, name, contentName *string) baseRule {
	result := baseRule{location: location, id: id}
	if name != nil {
		result.name = *name
	}
	if contentName != nil {
		result.contentName = *contentName
	}
	// JavaScript's `name || null` turns the empty string into no name.
	result.nameIsCapturing = result.name != "" && hasCaptureReferences(result.name)
	result.contentNameIsCapturing = result.contentName != "" && hasCaptureReferences(result.contentName)
	return result
}

func (r *baseRule) getID() ruleID { return r.id }

func (r *baseRule) getName(lineText string, captureIndices []oniguruma.Capture) string {
	if !r.nameIsCapturing || r.name == "" || captureIndices == nil {
		return r.name
	}
	return replaceCaptureReferences(r.name, []rune(lineText), captureIndices)
}

func (r *baseRule) getContentName(lineText string, captureIndices []oniguruma.Capture) string {
	if !r.contentNameIsCapturing || r.contentName == "" {
		return r.contentName
	}
	return replaceCaptureReferences(r.contentName, []rune(lineText), captureIndices)
}

type captureRule struct {
	baseRule
	retokenizeCapturedWithRuleID ruleID
}

func newCaptureRule(
	location *Location,
	id ruleID,
	name, contentName *string,
	retokenizeCapturedWithRuleID ruleID,
) *captureRule {
	return &captureRule{
		baseRule:                     newBaseRule(location, id, name, contentName),
		retokenizeCapturedWithRuleID: retokenizeCapturedWithRuleID,
	}
}

func (r *captureRule) dispose() {}

func (r *captureRule) collectPatterns(ruleRegistry, *regexpSourceList) {
	panic("textmate: capture rules cannot collect patterns")
}

func (r *captureRule) compile(ruleRegistry, *string) *compiledRule {
	panic("textmate: capture rules cannot be compiled")
}

func (r *captureRule) compileAG(ruleRegistry, *string, bool, bool) *compiledRule {
	panic("textmate: capture rules cannot be compiled")
}

type matchRule struct {
	baseRule
	match                  *regexpSource
	captures               []*captureRule
	cachedCompiledPatterns *regexpSourceList
}

func newMatchRule(
	location *Location,
	id ruleID,
	name *string,
	match string,
	captures []*captureRule,
) *matchRule {
	return &matchRule{
		baseRule: newBaseRule(location, id, name, nil),
		match:    newRegExpSource(match, id),
		captures: captures,
	}
}

func (r *matchRule) dispose() {
	if r.cachedCompiledPatterns != nil {
		r.cachedCompiledPatterns.dispose()
		r.cachedCompiledPatterns = nil
	}
}

func (r *matchRule) collectPatterns(_ ruleRegistry, out *regexpSourceList) {
	out.push(r.match)
}

func (r *matchRule) compile(grammar ruleRegistry, _ *string) *compiledRule {
	return r.getCachedCompiledPatterns(grammar).compile()
}

func (r *matchRule) compileAG(grammar ruleRegistry, _ *string, allowA, allowG bool) *compiledRule {
	return r.getCachedCompiledPatterns(grammar).compileAG(allowA, allowG)
}

func (r *matchRule) getCachedCompiledPatterns(grammar ruleRegistry) *regexpSourceList {
	if r.cachedCompiledPatterns == nil {
		r.cachedCompiledPatterns = newRuleRegExpSourceList(grammar)
		r.collectPatterns(grammar, r.cachedCompiledPatterns)
	}
	return r.cachedCompiledPatterns
}

type compilePatternsResult struct {
	patterns           []ruleID
	hasMissingPatterns bool
}

type includeOnlyRule struct {
	baseRule
	patterns               []ruleID
	hasMissingPatterns     bool
	cachedCompiledPatterns *regexpSourceList
}

func newIncludeOnlyRule(
	location *Location,
	id ruleID,
	name, contentName *string,
	patterns compilePatternsResult,
) *includeOnlyRule {
	return &includeOnlyRule{
		baseRule:           newBaseRule(location, id, name, contentName),
		patterns:           patterns.patterns,
		hasMissingPatterns: patterns.hasMissingPatterns,
	}
}

func (r *includeOnlyRule) dispose() {
	if r.cachedCompiledPatterns != nil {
		r.cachedCompiledPatterns.dispose()
		r.cachedCompiledPatterns = nil
	}
}

func (r *includeOnlyRule) collectPatterns(grammar ruleRegistry, out *regexpSourceList) {
	for _, pattern := range r.patterns {
		if nested := grammar.getRule(pattern); nested != nil {
			nested.collectPatterns(grammar, out)
		}
	}
}

func (r *includeOnlyRule) compile(grammar ruleRegistry, _ *string) *compiledRule {
	return r.getCachedCompiledPatterns(grammar).compile()
}

func (r *includeOnlyRule) compileAG(grammar ruleRegistry, _ *string, allowA, allowG bool) *compiledRule {
	return r.getCachedCompiledPatterns(grammar).compileAG(allowA, allowG)
}

func (r *includeOnlyRule) getCachedCompiledPatterns(grammar ruleRegistry) *regexpSourceList {
	if r.cachedCompiledPatterns == nil {
		r.cachedCompiledPatterns = newRuleRegExpSourceList(grammar)
		r.collectPatterns(grammar, r.cachedCompiledPatterns)
	}
	return r.cachedCompiledPatterns
}

type beginEndRule struct {
	baseRule
	begin                  *regexpSource
	beginCaptures          []*captureRule
	end                    *regexpSource
	endHasBackReferences   bool
	endCaptures            []*captureRule
	applyEndPatternLast    bool
	patterns               []ruleID
	hasMissingPatterns     bool
	cachedCompiledPatterns *regexpSourceList
}

func newBeginEndRule(
	location *Location,
	id ruleID,
	name, contentName *string,
	begin string,
	beginCaptures []*captureRule,
	end string,
	endCaptures []*captureRule,
	applyEndPatternLast bool,
	patterns compilePatternsResult,
) *beginEndRule {
	endSource := newRegExpSource(end, endRuleID)
	return &beginEndRule{
		baseRule:             newBaseRule(location, id, name, contentName),
		begin:                newRegExpSource(begin, id),
		beginCaptures:        beginCaptures,
		end:                  endSource,
		endHasBackReferences: endSource.hasBackReferences,
		endCaptures:          endCaptures,
		applyEndPatternLast:  applyEndPatternLast,
		patterns:             patterns.patterns,
		hasMissingPatterns:   patterns.hasMissingPatterns,
	}
}

func (r *beginEndRule) dispose() {
	if r.cachedCompiledPatterns != nil {
		r.cachedCompiledPatterns.dispose()
		r.cachedCompiledPatterns = nil
	}
}

func (r *beginEndRule) getEndWithResolvedBackReferences(
	lineText string,
	captureIndices []oniguruma.Capture,
) string {
	return r.end.resolveBackReferences(lineText, captureIndices)
}

func (r *beginEndRule) collectPatterns(_ ruleRegistry, out *regexpSourceList) {
	out.push(r.begin)
}

func (r *beginEndRule) compile(grammar ruleRegistry, endRegexSource *string) *compiledRule {
	return r.getCachedCompiledPatterns(grammar, endRegexSource).compile()
}

func (r *beginEndRule) compileAG(
	grammar ruleRegistry,
	endRegexSource *string,
	allowA, allowG bool,
) *compiledRule {
	return r.getCachedCompiledPatterns(grammar, endRegexSource).compileAG(allowA, allowG)
}

func (r *beginEndRule) getCachedCompiledPatterns(
	grammar ruleRegistry,
	endRegexSource *string,
) *regexpSourceList {
	if r.cachedCompiledPatterns == nil {
		r.cachedCompiledPatterns = newRuleRegExpSourceList(grammar)
		for _, pattern := range r.patterns {
			if nested := grammar.getRule(pattern); nested != nil {
				nested.collectPatterns(grammar, r.cachedCompiledPatterns)
			}
		}

		end := r.end
		if r.endHasBackReferences {
			end = end.clone()
		}
		if r.applyEndPatternLast {
			r.cachedCompiledPatterns.push(end)
		} else {
			r.cachedCompiledPatterns.unshift(end)
		}
	}

	if r.endHasBackReferences {
		resolved := ""
		if endRegexSource != nil {
			resolved = *endRegexSource
		}
		if r.applyEndPatternLast {
			r.cachedCompiledPatterns.setSource(r.cachedCompiledPatterns.length()-1, resolved)
		} else {
			r.cachedCompiledPatterns.setSource(0, resolved)
		}
	}
	return r.cachedCompiledPatterns
}

type beginWhileRule struct {
	baseRule
	begin                       *regexpSource
	beginCaptures               []*captureRule
	whileCaptures               []*captureRule
	while                       *regexpSource
	whileHasBackReferences      bool
	patterns                    []ruleID
	hasMissingPatterns          bool
	cachedCompiledPatterns      *regexpSourceList
	cachedCompiledWhilePatterns *regexpSourceList
}

func newBeginWhileRule(
	location *Location,
	id ruleID,
	name, contentName *string,
	begin string,
	beginCaptures []*captureRule,
	while string,
	whileCaptures []*captureRule,
	patterns compilePatternsResult,
) *beginWhileRule {
	whileSource := newRegExpSource(while, whileRuleID)
	return &beginWhileRule{
		baseRule:               newBaseRule(location, id, name, contentName),
		begin:                  newRegExpSource(begin, id),
		beginCaptures:          beginCaptures,
		whileCaptures:          whileCaptures,
		while:                  whileSource,
		whileHasBackReferences: whileSource.hasBackReferences,
		patterns:               patterns.patterns,
		hasMissingPatterns:     patterns.hasMissingPatterns,
	}
}

func (r *beginWhileRule) dispose() {
	if r.cachedCompiledPatterns != nil {
		r.cachedCompiledPatterns.dispose()
		r.cachedCompiledPatterns = nil
	}
	if r.cachedCompiledWhilePatterns != nil {
		r.cachedCompiledWhilePatterns.dispose()
		r.cachedCompiledWhilePatterns = nil
	}
}

func (r *beginWhileRule) getWhileWithResolvedBackReferences(
	lineText string,
	captureIndices []oniguruma.Capture,
) string {
	return r.while.resolveBackReferences(lineText, captureIndices)
}

func (r *beginWhileRule) collectPatterns(_ ruleRegistry, out *regexpSourceList) {
	out.push(r.begin)
}

func (r *beginWhileRule) compile(grammar ruleRegistry, _ *string) *compiledRule {
	return r.getCachedCompiledPatterns(grammar).compile()
}

func (r *beginWhileRule) compileAG(
	grammar ruleRegistry,
	_ *string,
	allowA, allowG bool,
) *compiledRule {
	return r.getCachedCompiledPatterns(grammar).compileAG(allowA, allowG)
}

func (r *beginWhileRule) getCachedCompiledPatterns(grammar ruleRegistry) *regexpSourceList {
	if r.cachedCompiledPatterns == nil {
		r.cachedCompiledPatterns = newRuleRegExpSourceList(grammar)
		for _, pattern := range r.patterns {
			if nested := grammar.getRule(pattern); nested != nil {
				nested.collectPatterns(grammar, r.cachedCompiledPatterns)
			}
		}
	}
	return r.cachedCompiledPatterns
}

func (r *beginWhileRule) compileWhile(
	grammar ruleRegistry,
	endRegexSource *string,
) *compiledRule {
	return r.getCachedCompiledWhilePatterns(grammar, endRegexSource).compile()
}

func (r *beginWhileRule) compileWhileAG(
	grammar ruleRegistry,
	endRegexSource *string,
	allowA, allowG bool,
) *compiledRule {
	return r.getCachedCompiledWhilePatterns(grammar, endRegexSource).compileAG(allowA, allowG)
}

func (r *beginWhileRule) getCachedCompiledWhilePatterns(
	grammar ruleRegistry,
	endRegexSource *string,
) *regexpSourceList {
	if r.cachedCompiledWhilePatterns == nil {
		r.cachedCompiledWhilePatterns = newRuleRegExpSourceList(grammar)
		item := r.while
		if r.whileHasBackReferences {
			item = item.clone()
		}
		r.cachedCompiledWhilePatterns.push(item)
	}
	if r.whileHasBackReferences {
		resolved := ""
		if endRegexSource != nil {
			resolved = *endRegexSource
		}
		r.cachedCompiledWhilePatterns.setSource(0, resolved)
	}
	return r.cachedCompiledWhilePatterns
}

// RuleFactory turns raw grammar nodes into registered rules. A factory belongs
// to one compiled Grammar; keeping raw node identity here avoids mutating the
// public JSON model with compilation state.
type RuleFactory struct {
	compiledRuleIDs map[*RawRule]ruleID
}

func newRuleFactory() *RuleFactory {
	return &RuleFactory{compiledRuleIDs: make(map[*RawRule]ruleID)}
}

func (f *RuleFactory) createCaptureRule(
	helper ruleFactoryHelper,
	location *Location,
	name, contentName *string,
	retokenizeCapturedWithRuleID ruleID,
) *captureRule {
	return helper.registerRule(func(id ruleID) rule {
		return newCaptureRule(location, id, name, contentName, retokenizeCapturedWithRuleID)
	}).(*captureRule)
}

func (f *RuleFactory) getCompiledRuleID(
	desc *RawRule,
	helper ruleFactoryHelper,
	repository RawRepository,
) ruleID {
	if desc == nil {
		return 0
	}
	if id, ok := f.compiledRuleIDs[desc]; ok {
		return id
	}

	registered := helper.registerRule(func(id ruleID) rule {
		// Reserve the ID before compiling captures or child patterns. Includes
		// can therefore point back at this raw node without allocating forever.
		f.compiledRuleIDs[desc] = id

		if desc.Match != nil && *desc.Match != "" {
			return newMatchRule(
				desc.Location,
				id,
				desc.Name,
				*desc.Match,
				f.compileCaptures(desc.Captures, helper, repository),
			)
		}

		if desc.Begin == nil {
			if desc.Repository != nil {
				repository = mergeRawRepositories(repository, desc.Repository)
			}
			patterns := desc.Patterns
			if patterns == nil && desc.Include != nil && *desc.Include != "" {
				include := *desc.Include
				patterns = []*RawRule{{Include: &include}}
			}
			return newIncludeOnlyRule(
				desc.Location,
				id,
				desc.Name,
				desc.ContentName,
				f.compilePatterns(patterns, helper, repository),
			)
		}

		beginCaptures := desc.BeginCaptures
		if beginCaptures == nil {
			beginCaptures = desc.Captures
		}
		if desc.While != nil && *desc.While != "" {
			whileCaptures := desc.WhileCaptures
			if whileCaptures == nil {
				whileCaptures = desc.Captures
			}
			return newBeginWhileRule(
				desc.Location,
				id,
				desc.Name,
				desc.ContentName,
				*desc.Begin,
				f.compileCaptures(beginCaptures, helper, repository),
				*desc.While,
				f.compileCaptures(whileCaptures, helper, repository),
				f.compilePatterns(desc.Patterns, helper, repository),
			)
		}

		endCaptures := desc.EndCaptures
		if endCaptures == nil {
			endCaptures = desc.Captures
		}
		end := ""
		if desc.End != nil {
			end = *desc.End
		}
		return newBeginEndRule(
			desc.Location,
			id,
			desc.Name,
			desc.ContentName,
			*desc.Begin,
			f.compileCaptures(beginCaptures, helper, repository),
			end,
			f.compileCaptures(endCaptures, helper, repository),
			desc.ApplyEndPatternLast,
			f.compilePatterns(desc.Patterns, helper, repository),
		)
	})
	return registered.getID()
}

func (f *RuleFactory) compileCaptures(
	captures RawCaptures,
	helper ruleFactoryHelper,
	repository RawRepository,
) []*captureRule {
	if captures == nil {
		return nil
	}

	maximumCaptureID := 0
	parsedIDs := make(map[string]int, len(captures))
	for captureID := range captures {
		numericCaptureID, ok := parseDecimalPrefix(captureID)
		if !ok || numericCaptureID < 0 {
			continue
		}
		parsedIDs[captureID] = numericCaptureID
		if numericCaptureID > maximumCaptureID {
			maximumCaptureID = numericCaptureID
		}
	}

	result := make([]*captureRule, maximumCaptureID+1)
	for captureID, capture := range captures {
		if capture == nil {
			continue
		}
		numericCaptureID, ok := parsedIDs[captureID]
		if !ok || numericCaptureID >= len(result) {
			continue
		}

		var retokenizeCapturedWithRuleID ruleID
		if capture.Patterns != nil {
			retokenizeCapturedWithRuleID = f.getCompiledRuleID(capture, helper, repository)
		}
		result[numericCaptureID] = f.createCaptureRule(
			helper,
			capture.Location,
			capture.Name,
			capture.ContentName,
			retokenizeCapturedWithRuleID,
		)
	}
	return result
}

func (f *RuleFactory) compilePatterns(
	patterns []*RawRule,
	helper ruleFactoryHelper,
	repository RawRepository,
) compilePatternsResult {
	result := make([]ruleID, 0, len(patterns))
	for _, pattern := range patterns {
		if pattern == nil {
			continue
		}

		var id ruleID
		found := false
		if pattern.Include != nil && *pattern.Include != "" {
			reference := parseInclude(*pattern.Include)
			switch reference := reference.(type) {
			case baseReference, selfReference:
				if included := repository[*pattern.Include]; included != nil {
					id = f.getCompiledRuleID(included, helper, repository)
					found = true
				}
			case relativeReference:
				if included := repository[reference.ruleName]; included != nil {
					id = f.getCompiledRuleID(included, helper, repository)
					found = true
				}
			case topLevelReference:
				if external := helper.getExternalGrammar(reference.scopeName, repository); external != nil {
					if included := external.Repository["$self"]; included != nil {
						id = f.getCompiledRuleID(included, helper, external.Repository)
						found = true
					}
				}
			case topLevelRepositoryReference:
				if external := helper.getExternalGrammar(reference.scopeName, repository); external != nil {
					if included := external.Repository[reference.ruleName]; included != nil {
						id = f.getCompiledRuleID(included, helper, external.Repository)
						found = true
					}
				}
			}
		} else {
			id = f.getCompiledRuleID(pattern, helper, repository)
			found = id != 0
		}

		if !found {
			continue
		}
		compiled := helper.getRule(id)
		if compiled != nil && ruleShouldBePruned(compiled) {
			continue
		}
		result = append(result, id)
	}

	return compilePatternsResult{
		patterns:           result,
		hasMissingPatterns: len(patterns) != len(result),
	}
}

func ruleShouldBePruned(compiled rule) bool {
	switch typed := compiled.(type) {
	case *includeOnlyRule:
		return typed.hasMissingPatterns && len(typed.patterns) == 0
	case *beginEndRule:
		return typed.hasMissingPatterns && len(typed.patterns) == 0
	case *beginWhileRule:
		return typed.hasMissingPatterns && len(typed.patterns) == 0
	default:
		return false
	}
}

// parseDecimalPrefix matches JavaScript parseInt(value, 10), which is used by
// vscode-textmate for capture-map keys. It intentionally accepts "1-extra".
func parseDecimalPrefix(value string) (int, bool) {
	value = strings.TrimLeft(value, " \t\n\r\f\v")
	if value == "" {
		return 0, false
	}
	sign := 1
	if value[0] == '+' || value[0] == '-' {
		if value[0] == '-' {
			sign = -1
		}
		value = value[1:]
	}
	end := 0
	for end < len(value) && value[end] >= '0' && value[end] <= '9' {
		end++
	}
	if end == 0 {
		return 0, false
	}
	numeric, err := strconv.Atoi(value[:end])
	if err != nil {
		return 0, false
	}
	return sign * numeric, true
}
