package textmate

import (
	"errors"
	"reflect"
	"sort"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func registryIncludeRule(value string) *RawRule {
	return &RawRule{Include: &value}
}

func TestRegistryLoadsTransitiveFullPartialAndInjectionDependencies(t *testing.T) {
	grammars := map[string]*RawGrammar{
		"source.root": {
			ScopeName: "source.root",
			Patterns: []*RawRule{
				registryIncludeRule("source.full"),
				registryIncludeRule("source.partial#entry"),
			},
		},
		"source.full": {
			ScopeName: "source.full",
			Patterns:  []*RawRule{registryIncludeRule("source.leaf")},
		},
		"source.partial": {
			ScopeName: "source.partial",
			Repository: RawRepository{
				"entry":  registryIncludeRule("source.partial-leaf"),
				"unused": registryIncludeRule("source.must-not-load"),
			},
		},
		"source.leaf":         {ScopeName: "source.leaf"},
		"source.partial-leaf": {ScopeName: "source.partial-leaf"},
		"source.inject-one": {
			ScopeName:         "source.inject-one",
			InjectionSelector: "L:source.root",
		},
		"source.inject-two": {
			ScopeName:         "source.inject-two",
			InjectionSelector: "R:source.root",
		},
	}
	var loaded []string
	registry := NewRegistry(RegistryOptions{
		LoadGrammar: func(scopeName string) (*RawGrammar, error) {
			loaded = append(loaded, scopeName)
			return grammars[scopeName], nil
		},
		GetInjections: func(scopeName string) []string {
			if scopeName == "source.root" {
				return []string{"source.inject-two", "source.inject-one"}
			}
			return nil
		},
	})

	grammar, err := registry.LoadGrammar("source.root")
	if err != nil {
		t.Fatal(err)
	}
	if grammar == nil || grammar.rootScopeName != "source.root" {
		t.Fatalf("grammar = %#v", grammar)
	}
	if got, want := registry.injections("source.root"), []string{"source.inject-two", "source.inject-one"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("injections = %v, want %v", got, want)
	}
	sort.Strings(loaded)
	wantLoaded := []string{
		"source.full", "source.inject-one", "source.inject-two", "source.leaf",
		"source.partial", "source.partial-leaf", "source.root",
	}
	if !reflect.DeepEqual(loaded, wantLoaded) {
		t.Fatalf("loaded scopes = %v, want %v", loaded, wantLoaded)
	}
	if registry.lookup("source.must-not-load") != nil {
		t.Fatal("unused part of a partial dependency was traversed")
	}

	again, err := registry.LoadGrammar("source.root")
	if err != nil || again != grammar {
		t.Fatalf("cached load = (%p, %v), want (%p, nil)", again, err, grammar)
	}
}

func TestRegistryIgnoresMissingOptionalDependency(t *testing.T) {
	registry := NewRegistry(RegistryOptions{LoadGrammar: func(scopeName string) (*RawGrammar, error) {
		if scopeName == "source.root" {
			return &RawGrammar{
				ScopeName: "source.root",
				Patterns:  []*RawRule{registryIncludeRule("source.optional")},
			}, nil
		}
		return nil, nil
	}})
	if grammar, err := registry.LoadGrammar("source.root"); err != nil || grammar == nil {
		t.Fatalf("LoadGrammar = (%v, %v), want grammar and nil", grammar, err)
	}
}

func TestRegistryReportsAndCachesMissingRoot(t *testing.T) {
	var calls atomic.Int32
	registry := NewRegistry(RegistryOptions{LoadGrammar: func(string) (*RawGrammar, error) {
		calls.Add(1)
		return nil, nil
	}})
	for range 2 {
		grammar, err := registry.LoadGrammar("source.missing")
		if grammar != nil || err == nil || err.Error() != "textmate: no grammar provided for <source.missing>" {
			t.Fatalf("LoadGrammar = (%v, %v)", grammar, err)
		}
	}
	if got := calls.Load(); got != 1 {
		t.Fatalf("loader calls = %d, want 1", got)
	}
}

func TestRegistryPropagatesAndCachesLoaderError(t *testing.T) {
	wantErr := errors.New("loader failed")
	var calls atomic.Int32
	registry := NewRegistry(RegistryOptions{LoadGrammar: func(string) (*RawGrammar, error) {
		calls.Add(1)
		return nil, wantErr
	}})
	for range 2 {
		grammar, err := registry.LoadGrammar("source.test")
		if grammar != nil || !errors.Is(err, wantErr) {
			t.Fatalf("LoadGrammar = (%v, %v), want nil and %v", grammar, err, wantErr)
		}
	}
	if got := calls.Load(); got != 1 {
		t.Fatalf("loader calls = %d, want 1", got)
	}
}

func TestRegistryRejectsMismatchedLoaderScope(t *testing.T) {
	registry := NewRegistry(RegistryOptions{LoadGrammar: func(string) (*RawGrammar, error) {
		return &RawGrammar{ScopeName: "source.other"}, nil
	}})
	grammar, err := registry.LoadGrammar("source.requested")
	if grammar != nil || err == nil || err.Error() != "textmate: grammar loader for <source.requested> returned grammar for <source.other>" {
		t.Fatalf("LoadGrammar = (%v, %v)", grammar, err)
	}
	if registry.lookup("source.other") != nil {
		t.Fatal("mismatched grammar was cached")
	}
}

func TestRegistryCoalescesConcurrentLoads(t *testing.T) {
	var rootCalls atomic.Int32
	var dependencyCalls atomic.Int32
	registry := NewRegistry(RegistryOptions{LoadGrammar: func(scopeName string) (*RawGrammar, error) {
		time.Sleep(time.Millisecond)
		switch scopeName {
		case "source.root":
			rootCalls.Add(1)
			return &RawGrammar{
				ScopeName: "source.root",
				Patterns:  []*RawRule{registryIncludeRule("source.dependency")},
			}, nil
		case "source.dependency":
			dependencyCalls.Add(1)
			return &RawGrammar{ScopeName: "source.dependency"}, nil
		default:
			return nil, nil
		}
	}})

	const workers = 32
	results := make([]*Grammar, workers)
	errors := make([]error, workers)
	var wg sync.WaitGroup
	start := make(chan struct{})
	for index := range workers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			results[index], errors[index] = registry.LoadGrammar("source.root")
		}()
	}
	close(start)
	wg.Wait()
	for index := range workers {
		if errors[index] != nil || results[index] == nil || results[index] != results[0] {
			t.Fatalf("worker %d = (%p, %v), first=%p", index, results[index], errors[index], results[0])
		}
	}
	if rootCalls.Load() != 1 || dependencyCalls.Load() != 1 {
		t.Fatalf("loader calls root=%d dependency=%d, want 1 each", rootCalls.Load(), dependencyCalls.Load())
	}
}

func TestRegistryCallbacksRunWithoutRegistryLock(t *testing.T) {
	var registry *Registry
	registry = NewRegistry(RegistryOptions{
		LoadGrammar: func(scopeName string) (*RawGrammar, error) {
			_ = registry.lookup(scopeName)
			return &RawGrammar{ScopeName: scopeName}, nil
		},
		GetInjections: func(scopeName string) []string {
			_ = registry.injections(scopeName)
			return nil
		},
	})
	if _, err := registry.LoadGrammar("source.test"); err != nil {
		t.Fatal(err)
	}
}

func TestRegistryDisposeIsIdempotentAndPreventsLoads(t *testing.T) {
	registry := NewRegistry(RegistryOptions{LoadGrammar: func(scopeName string) (*RawGrammar, error) {
		return &RawGrammar{ScopeName: scopeName}, nil
	}})
	if _, err := registry.LoadGrammar("source.test"); err != nil {
		t.Fatal(err)
	}
	registry.Dispose()
	registry.Dispose()
	if registry.options.LoadGrammar != nil || registry.options.GetInjections != nil {
		t.Fatal("Dispose retained registry callback closures")
	}
	if registry.rawGrammars != nil || registry.grammars != nil || registry.injectionMap != nil || registry.grammarLoads != nil {
		t.Fatal("Dispose retained registry maps")
	}
	if grammar, err := registry.LoadGrammar("source.test"); grammar != nil || !errors.Is(err, errRegistryDisposed) {
		t.Fatalf("load after dispose = (%v, %v)", grammar, err)
	}
}

func TestRegistryDisposeDuringInflightLoadReleasesWaiters(t *testing.T) {
	started := make(chan struct{})
	release := make(chan struct{})
	registry := NewRegistry(RegistryOptions{LoadGrammar: func(scopeName string) (*RawGrammar, error) {
		close(started)
		<-release
		return &RawGrammar{ScopeName: scopeName}, nil
	}})

	done := make(chan error, 1)
	go func() {
		_, err := registry.LoadGrammar("source.test")
		done <- err
	}()
	<-started
	registry.Dispose()
	close(release)

	select {
	case err := <-done:
		if !errors.Is(err, errRegistryDisposed) {
			t.Fatalf("in-flight load error = %v, want disposed", err)
		}
	case <-time.After(time.Second):
		t.Fatal("in-flight load did not finish after Dispose")
	}
}
