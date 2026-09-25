package textmate

import (
	"errors"
	"fmt"
	"sync"
)

var errRegistryDisposed = errors.New("textmate: registry disposed")

// RegistryOptions provides the callbacks used to lazily resolve grammars.
// A loader may return (nil, nil) when a scope is not available. Missing root
// grammars are errors; missing optional dependencies are ignored.
type RegistryOptions struct {
	LoadGrammar   func(scopeName string) (*RawGrammar, error)
	GetInjections func(scopeName string) []string
}

// Registry owns raw and compiled grammars and resolves their transitive
// dependencies. Callback invocations never occur while the registry lock is
// held.
type Registry struct {
	options RegistryOptions

	mu           sync.RWMutex
	rawGrammars  map[string]*RawGrammar
	grammars     map[string]*Grammar
	injectionMap map[string][]string
	grammarLoads map[string]*grammarLoad
	disposed     bool
}

// grammarLoad is retained after completion. Besides coalescing concurrent
// requests, this mirrors vscode-textmate's ensure-grammar cache: an absent
// grammar or loader error is not repeatedly requested.
type grammarLoad struct {
	done chan struct{}
	err  error
}

// NewRegistry constructs an empty lazy-loading grammar registry.
func NewRegistry(options RegistryOptions) *Registry {
	return &Registry{
		options:      options,
		rawGrammars:  make(map[string]*RawGrammar),
		grammars:     make(map[string]*Grammar),
		injectionMap: make(map[string][]string),
		grammarLoads: make(map[string]*grammarLoad),
	}
}

// LoadGrammar loads scopeName and all reachable includes and injections, then
// returns the cached compiled grammar for scopeName.
func (r *Registry) LoadGrammar(scopeName string) (*Grammar, error) {
	if r == nil {
		return nil, fmt.Errorf("textmate: nil registry")
	}
	r.mu.RLock()
	disposed := r.disposed
	r.mu.RUnlock()
	if disposed {
		return nil, errRegistryDisposed
	}
	processor := newScopeDependencyProcessor(r, scopeName)
	for len(processor.queue) != 0 {
		queue := append([]absoluteRuleReference(nil), processor.queue...)
		for _, request := range queue {
			if err := r.loadSingleGrammar(request.dependencyScopeName()); err != nil {
				return nil, err
			}
		}
		if err := processor.processQueue(); err != nil {
			return nil, err
		}
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	if r.disposed {
		return nil, errRegistryDisposed
	}
	if grammar := r.grammars[scopeName]; grammar != nil {
		return grammar, nil
	}
	raw := r.rawGrammars[scopeName]
	if raw == nil {
		return nil, fmt.Errorf("textmate: no grammar provided for <%s>", scopeName)
	}
	grammar := newGrammar(scopeName, raw, r)
	r.grammars[scopeName] = grammar
	return grammar, nil
}

func (r *Registry) loadSingleGrammar(scopeName string) error {
	r.mu.Lock()
	if r.disposed {
		r.mu.Unlock()
		return errRegistryDisposed
	}
	if load := r.grammarLoads[scopeName]; load != nil {
		r.mu.Unlock()
		<-load.done
		return load.err
	}
	load := &grammarLoad{done: make(chan struct{})}
	r.grammarLoads[scopeName] = load
	loader := r.options.LoadGrammar
	getInjections := r.options.GetInjections
	r.mu.Unlock()

	var raw *RawGrammar
	var err error
	if loader != nil {
		raw, err = loader(scopeName)
	}
	if err == nil && raw != nil && raw.ScopeName != scopeName {
		err = fmt.Errorf(
			"textmate: grammar loader for <%s> returned grammar for <%s>",
			scopeName,
			raw.ScopeName,
		)
	}
	var injections []string
	if err == nil && raw != nil && getInjections != nil {
		injections = append([]string(nil), getInjections(scopeName)...)
	}

	r.mu.Lock()
	if r.disposed {
		err = errRegistryDisposed
	} else if err == nil && raw != nil {
		r.rawGrammars[raw.ScopeName] = raw
		r.injectionMap[raw.ScopeName] = injections
	}
	load.err = err
	close(load.done)
	r.mu.Unlock()
	return err
}

// Dispose releases the registry's references to loaded grammars. It is
// idempotent; subsequent loads fail.
func (r *Registry) Dispose() {
	if r == nil {
		return
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.disposed {
		return
	}
	r.disposed = true
	r.options = RegistryOptions{}
	r.rawGrammars = nil
	r.grammars = nil
	r.injectionMap = nil
	r.grammarLoads = nil
}

func (r *Registry) lookup(scopeName string) *RawGrammar {
	if r == nil {
		return nil
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.rawGrammars[scopeName]
}

func (r *Registry) injections(scopeName string) []string {
	if r == nil {
		return nil
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	return append([]string(nil), r.injectionMap[scopeName]...)
}
