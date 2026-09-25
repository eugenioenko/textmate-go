package textmate

import (
	"fmt"
	"sort"
)

type grammarRepository interface {
	lookup(scopeName string) *RawGrammar
	injections(scopeName string) []string
}

type absoluteRuleReference interface {
	dependencyScopeName() string
	toKey() string
	isAbsoluteRuleReference()
}

type topLevelRuleReference struct {
	scopeName string
}

func (r topLevelRuleReference) dependencyScopeName() string { return r.scopeName }
func (r topLevelRuleReference) toKey() string               { return r.scopeName }
func (topLevelRuleReference) isAbsoluteRuleReference()      {}

type topLevelRepositoryRuleReference struct {
	scopeName string
	ruleName  string
}

func (r topLevelRepositoryRuleReference) dependencyScopeName() string { return r.scopeName }
func (r topLevelRepositoryRuleReference) toKey() string {
	return r.scopeName + "#" + r.ruleName
}
func (topLevelRepositoryRuleReference) isAbsoluteRuleReference() {}

type externalReferenceCollector struct {
	references        []absoluteRuleReference
	seenReferenceKeys map[string]struct{}
	visitedRules      map[*RawRule]struct{}
}

func newExternalReferenceCollector() *externalReferenceCollector {
	return &externalReferenceCollector{
		seenReferenceKeys: make(map[string]struct{}),
		visitedRules:      make(map[*RawRule]struct{}),
	}
}

func (c *externalReferenceCollector) add(reference absoluteRuleReference) {
	key := reference.toKey()
	if _, ok := c.seenReferenceKeys[key]; ok {
		return
	}
	c.seenReferenceKeys[key] = struct{}{}
	c.references = append(c.references, reference)
}

type scopeDependencyProcessor struct {
	repo                     grammarRepository
	initialScopeName         string
	seenFullScopeRequests    map[string]struct{}
	seenPartialScopeRequests map[string]struct{}
	queue                    []absoluteRuleReference
}

func newScopeDependencyProcessor(repo grammarRepository, initialScopeName string) *scopeDependencyProcessor {
	return &scopeDependencyProcessor{
		repo:                     repo,
		initialScopeName:         initialScopeName,
		seenFullScopeRequests:    map[string]struct{}{initialScopeName: {}},
		seenPartialScopeRequests: make(map[string]struct{}),
		queue:                    []absoluteRuleReference{topLevelRuleReference{scopeName: initialScopeName}},
	}
}

func (p *scopeDependencyProcessor) processQueue() error {
	queue := p.queue
	p.queue = nil

	dependencies := newExternalReferenceCollector()
	for _, dependency := range queue {
		if err := collectReferencesOfReference(
			dependency,
			p.initialScopeName,
			p.repo,
			dependencies,
		); err != nil {
			return err
		}
	}

	for _, dependency := range dependencies.references {
		scopeName := dependency.dependencyScopeName()
		switch dependency.(type) {
		case topLevelRuleReference:
			if _, ok := p.seenFullScopeRequests[scopeName]; ok {
				continue
			}
			p.seenFullScopeRequests[scopeName] = struct{}{}
			p.queue = append(p.queue, dependency)
		case topLevelRepositoryRuleReference:
			if _, ok := p.seenFullScopeRequests[scopeName]; ok {
				continue
			}
			if _, ok := p.seenPartialScopeRequests[dependency.toKey()]; ok {
				continue
			}
			p.seenPartialScopeRequests[dependency.toKey()] = struct{}{}
			p.queue = append(p.queue, dependency)
		}
	}
	return nil
}

func collectReferencesOfReference(
	reference absoluteRuleReference,
	baseGrammarScopeName string,
	repo grammarRepository,
	result *externalReferenceCollector,
) error {
	selfGrammar := repo.lookup(reference.dependencyScopeName())
	if selfGrammar == nil {
		if reference.dependencyScopeName() == baseGrammarScopeName {
			return fmt.Errorf("textmate: no grammar provided for <%s>", baseGrammarScopeName)
		}
		return nil
	}

	baseGrammar := repo.lookup(baseGrammarScopeName)
	if baseGrammar == nil {
		return fmt.Errorf("textmate: no grammar provided for <%s>", baseGrammarScopeName)
	}

	context := referenceContext{baseGrammar: baseGrammar, selfGrammar: selfGrammar}
	switch reference := reference.(type) {
	case topLevelRuleReference:
		collectExternalReferencesInTopLevelRule(context, result)
	case topLevelRepositoryRuleReference:
		collectExternalReferencesInTopLevelRepositoryRule(
			reference.ruleName,
			referenceContextWithRepository{
				referenceContext: context,
				repository:       selfGrammar.Repository,
			},
			result,
		)
	}

	for _, injection := range repo.injections(reference.dependencyScopeName()) {
		result.add(topLevelRuleReference{scopeName: injection})
	}
	return nil
}

type referenceContext struct {
	baseGrammar *RawGrammar
	selfGrammar *RawGrammar
}

type referenceContextWithRepository struct {
	referenceContext
	repository RawRepository
}

func collectExternalReferencesInTopLevelRepositoryRule(
	ruleName string,
	context referenceContextWithRepository,
	result *externalReferenceCollector,
) {
	if rule := context.repository[ruleName]; rule != nil {
		collectExternalReferencesInRules([]*RawRule{rule}, context, result)
	}
}

func collectExternalReferencesInTopLevelRule(
	context referenceContext,
	result *externalReferenceCollector,
) {
	childContext := referenceContextWithRepository{
		referenceContext: context,
		repository:       context.selfGrammar.Repository,
	}
	collectExternalReferencesInRules(context.selfGrammar.Patterns, childContext, result)

	// JavaScript preserves object insertion order. RawInjections is a Go map,
	// so sort selectors to make dependency discovery stable across runs.
	selectors := make([]string, 0, len(context.selfGrammar.Injections))
	for selector := range context.selfGrammar.Injections {
		selectors = append(selectors, selector)
	}
	sort.Strings(selectors)
	for _, selector := range selectors {
		collectExternalReferencesInRules(
			[]*RawRule{context.selfGrammar.Injections[selector]},
			childContext,
			result,
		)
	}
}

func collectExternalReferencesInRules(
	rules []*RawRule,
	context referenceContextWithRepository,
	result *externalReferenceCollector,
) {
	for _, rule := range rules {
		if rule == nil {
			continue
		}
		if _, ok := result.visitedRules[rule]; ok {
			continue
		}
		result.visitedRules[rule] = struct{}{}

		patternRepository := context.repository
		if rule.Repository != nil {
			patternRepository = mergeRawRepositories(context.repository, rule.Repository)
		}
		childContext := context
		childContext.repository = patternRepository

		collectExternalReferencesInRules(rule.Patterns, childContext, result)

		if rule.Include == nil || *rule.Include == "" {
			continue
		}

		reference := parseInclude(*rule.Include)
		switch reference := reference.(type) {
		case baseReference:
			baseContext := context.referenceContext
			baseContext.selfGrammar = context.baseGrammar
			collectExternalReferencesInTopLevelRule(baseContext, result)
		case selfReference:
			collectExternalReferencesInTopLevelRule(context.referenceContext, result)
		case relativeReference:
			collectExternalReferencesInTopLevelRepositoryRule(reference.ruleName, childContext, result)
		case topLevelReference:
			collectExternalReferencesForTopLevelInclude(
				reference.scopeName,
				"",
				false,
				childContext,
				result,
			)
		case topLevelRepositoryReference:
			collectExternalReferencesForTopLevelInclude(
				reference.scopeName,
				reference.ruleName,
				true,
				childContext,
				result,
			)
		}
	}
}

func collectExternalReferencesForTopLevelInclude(
	scopeName string,
	ruleName string,
	isRepositoryReference bool,
	context referenceContextWithRepository,
	result *externalReferenceCollector,
) {
	var selfGrammar *RawGrammar
	switch scopeName {
	case context.selfGrammar.ScopeName:
		selfGrammar = context.selfGrammar
	case context.baseGrammar.ScopeName:
		selfGrammar = context.baseGrammar
	}

	if selfGrammar == nil {
		if isRepositoryReference {
			result.add(topLevelRepositoryRuleReference{scopeName: scopeName, ruleName: ruleName})
		} else {
			result.add(topLevelRuleReference{scopeName: scopeName})
		}
		return
	}

	localContext := context
	localContext.selfGrammar = selfGrammar
	if isRepositoryReference {
		collectExternalReferencesInTopLevelRepositoryRule(ruleName, localContext, result)
	} else {
		collectExternalReferencesInTopLevelRule(localContext.referenceContext, result)
	}
}

func mergeRawRepositories(base RawRepository, overlay RawRepository) RawRepository {
	merged := make(RawRepository, len(base)+len(overlay))
	for key, rule := range base {
		merged[key] = rule
	}
	for key, rule := range overlay {
		merged[key] = rule
	}
	return merged
}

type includeReference interface {
	isIncludeReference()
}

type baseReference struct{}

func (baseReference) isIncludeReference() {}

type selfReference struct{}

func (selfReference) isIncludeReference() {}

type relativeReference struct {
	ruleName string
}

func (relativeReference) isIncludeReference() {}

type topLevelReference struct {
	scopeName string
}

func (topLevelReference) isIncludeReference() {}

type topLevelRepositoryReference struct {
	scopeName string
	ruleName  string
}

func (topLevelRepositoryReference) isIncludeReference() {}

func parseInclude(include string) includeReference {
	switch include {
	case "$base":
		return baseReference{}
	case "$self":
		return selfReference{}
	}

	for index, char := range include {
		if char != '#' {
			continue
		}
		if index == 0 {
			return relativeReference{ruleName: include[1:]}
		}
		return topLevelRepositoryReference{
			scopeName: include[:index],
			ruleName:  include[index+1:],
		}
	}
	return topLevelReference{scopeName: include}
}
