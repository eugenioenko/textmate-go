package textmate

import (
	"reflect"
	"strings"
	"testing"
)

type dependencyTestRepository struct {
	grammars        map[string]*RawGrammar
	injectionScopes map[string][]string
}

func (r *dependencyTestRepository) lookup(scopeName string) *RawGrammar {
	return r.grammars[scopeName]
}

func (r *dependencyTestRepository) injections(scopeName string) []string {
	return r.injectionScopes[scopeName]
}

func TestParseInclude(t *testing.T) {
	tests := []struct {
		include string
		want    includeReference
	}{
		{include: "$base", want: baseReference{}},
		{include: "$self", want: selfReference{}},
		{include: "#rule", want: relativeReference{ruleName: "rule"}},
		{include: "#", want: relativeReference{ruleName: ""}},
		{include: "source.js", want: topLevelReference{scopeName: "source.js"}},
		{include: "", want: topLevelReference{scopeName: ""}},
		{
			include: "source.js#expression",
			want:    topLevelRepositoryReference{scopeName: "source.js", ruleName: "expression"},
		},
		{
			include: "source.js#rule#suffix",
			want:    topLevelRepositoryReference{scopeName: "source.js", ruleName: "rule#suffix"},
		},
	}

	for _, test := range tests {
		t.Run(test.include, func(t *testing.T) {
			if got := parseInclude(test.include); !reflect.DeepEqual(got, test.want) {
				t.Fatalf("parseInclude(%q) = %#v, want %#v", test.include, got, test.want)
			}
		})
	}
}

func TestScopeDependencyProcessorFullAndPartialRequests(t *testing.T) {
	root := grammarWithPatterns("source.root",
		includeRule("source.full"),
		includeRule("source.full#already-covered"),
		includeRule("source.partial#entry"),
		includeRule("source.partial#entry"),
	)
	root.Injections = RawInjections{
		"L:source.root": includeRule("source.from-raw-injection#item"),
	}
	full := grammarWithPatterns("source.full", includeRule("source.transitive"))
	partial := grammarWithPatterns("source.partial", includeRule("source.must-not-load"))
	partial.Repository = RawRepository{
		"entry": includeRule("source.from-partial"),
	}

	repo := &dependencyTestRepository{
		grammars: map[string]*RawGrammar{
			root.ScopeName:    root,
			full.ScopeName:    full,
			partial.ScopeName: partial,
		},
		injectionScopes: map[string][]string{
			"source.root":    {"source.registry-injection"},
			"source.partial": {"source.partial-injection"},
		},
	}
	processor := newScopeDependencyProcessor(repo, root.ScopeName)

	if err := processor.processQueue(); err != nil {
		t.Fatalf("first processQueue() error = %v", err)
	}
	if got, want := dependencyKeys(processor.queue), []string{
		"source.full",
		"source.partial#entry",
		"source.from-raw-injection#item",
		"source.registry-injection",
	}; !reflect.DeepEqual(got, want) {
		t.Fatalf("first queue = %#v, want %#v", got, want)
	}
	if _, ok := processor.seenPartialScopeRequests["source.full#already-covered"]; ok {
		t.Fatal("partial request was retained after a full request for the same scope")
	}

	if err := processor.processQueue(); err != nil {
		t.Fatalf("second processQueue() error = %v", err)
	}
	got := dependencyKeys(processor.queue)
	if !containsString(got, "source.transitive") {
		t.Fatalf("full grammar patterns were not traversed: queue = %#v", got)
	}
	if !containsString(got, "source.from-partial") {
		t.Fatalf("partial repository rule was not traversed: queue = %#v", got)
	}
	if !containsString(got, "source.partial-injection") {
		t.Fatalf("injections for a partial grammar request were not collected: queue = %#v", got)
	}
	if containsString(got, "source.must-not-load") {
		t.Fatalf("partial request traversed the grammar's top-level patterns: queue = %#v", got)
	}
}

func TestDependencyTraversalResolvesBaseSelfRelativeAndNestedRepositories(t *testing.T) {
	baseRepositoryRule := includeRule("source.from-base-repository")
	base := grammarWithPatterns("source.base", includeRule("#base-rule"))
	base.Repository = RawRepository{
		"base-rule": baseRepositoryRule,
		"target":    includeRule("source.from-outer-repository"),
	}

	overlayRule := includeRule("source.from-inner-repository")
	nested := &RawRule{
		Repository: RawRepository{"target": overlayRule},
		Patterns: []*RawRule{
			includeRule("#target"),
			includeRule("$self"),
			includeRule("$base"),
		},
	}
	external := grammarWithPatterns("source.external", nested)
	repo := &dependencyTestRepository{grammars: map[string]*RawGrammar{
		base.ScopeName:     base,
		external.ScopeName: external,
	}}

	collector := newExternalReferenceCollector()
	err := collectReferencesOfReference(
		topLevelRuleReference{scopeName: external.ScopeName},
		base.ScopeName,
		repo,
		collector,
	)
	if err != nil {
		t.Fatalf("collectReferencesOfReference() error = %v", err)
	}
	got := referenceKeys(collector.references)
	for _, key := range []string{
		"source.from-inner-repository",
		"source.from-base-repository",
	} {
		if !containsString(got, key) {
			t.Errorf("dependency %q missing from %#v", key, got)
		}
	}
	if containsString(got, "source.from-outer-repository") {
		t.Fatalf("nested repository did not override outer repository: %#v", got)
	}
}

func TestDependencyTraversalUsesRawRuleIdentityForCycleSafety(t *testing.T) {
	shared := includeRule("#local")
	grammarA := grammarWithPatterns("source.a", shared)
	grammarA.Repository = RawRepository{"local": includeRule("source.from-a")}
	grammarB := grammarWithPatterns("source.b", shared)
	grammarB.Repository = RawRepository{"local": includeRule("source.from-b")}
	base := grammarWithPatterns("source.base")
	repo := &dependencyTestRepository{grammars: map[string]*RawGrammar{
		base.ScopeName:     base,
		grammarA.ScopeName: grammarA,
		grammarB.ScopeName: grammarB,
	}}
	collector := newExternalReferenceCollector()

	for _, reference := range []absoluteRuleReference{
		topLevelRuleReference{scopeName: grammarA.ScopeName},
		topLevelRuleReference{scopeName: grammarB.ScopeName},
	} {
		if err := collectReferencesOfReference(reference, base.ScopeName, repo, collector); err != nil {
			t.Fatalf("collectReferencesOfReference() error = %v", err)
		}
	}

	if got, want := referenceKeys(collector.references), []string{"source.from-a"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("dependencies from shared raw-rule identity = %#v, want %#v", got, want)
	}

	cycle := includeRule("$self")
	cyclicGrammar := grammarWithPatterns("source.cycle", cycle)
	repo.grammars[cyclicGrammar.ScopeName] = cyclicGrammar
	if err := collectReferencesOfReference(
		topLevelRuleReference{scopeName: cyclicGrammar.ScopeName},
		base.ScopeName,
		repo,
		newExternalReferenceCollector(),
	); err != nil {
		t.Fatalf("self cycle returned error = %v", err)
	}
}

func TestScopeDependencyProcessorMissingGrammars(t *testing.T) {
	repo := &dependencyTestRepository{grammars: make(map[string]*RawGrammar)}
	processor := newScopeDependencyProcessor(repo, "source.missing")
	if err := processor.processQueue(); err == nil || !strings.Contains(err.Error(), "source.missing") {
		t.Fatalf("missing initial grammar error = %v", err)
	}

	base := grammarWithPatterns("source.base")
	repo.grammars[base.ScopeName] = base
	collector := newExternalReferenceCollector()
	if err := collectReferencesOfReference(
		topLevelRuleReference{scopeName: "source.optional"},
		base.ScopeName,
		repo,
		collector,
	); err != nil {
		t.Fatalf("missing external grammar error = %v", err)
	}
	if len(collector.references) != 0 {
		t.Fatalf("missing external grammar added dependencies: %#v", collector.references)
	}
}

func TestExternalReferenceCollectorDeduplicatesByDependencyKey(t *testing.T) {
	collector := newExternalReferenceCollector()
	collector.add(topLevelRuleReference{scopeName: "source.a"})
	collector.add(topLevelRuleReference{scopeName: "source.a"})
	collector.add(topLevelRepositoryRuleReference{scopeName: "source.a", ruleName: "part"})
	collector.add(topLevelRepositoryRuleReference{scopeName: "source.a", ruleName: "part"})
	if got, want := referenceKeys(collector.references), []string{"source.a", "source.a#part"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("references = %#v, want %#v", got, want)
	}
}

func grammarWithPatterns(scopeName string, patterns ...*RawRule) *RawGrammar {
	return &RawGrammar{ScopeName: scopeName, Patterns: patterns}
}

func includeRule(include string) *RawRule {
	return &RawRule{Include: &include}
}

func dependencyKeys(references []absoluteRuleReference) []string {
	return referenceKeys(references)
}

func referenceKeys(references []absoluteRuleReference) []string {
	result := make([]string, len(references))
	for index, reference := range references {
		result[index] = reference.toKey()
	}
	return result
}

func containsString(values []string, want string) bool {
	for _, value := range values {
		if value == want {
			return true
		}
	}
	return false
}
