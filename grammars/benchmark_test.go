package grammars

import (
	"runtime"
	"sync"
	"testing"

	textmate "github.com/eugenioenko/textmate-go"
)

func BenchmarkFirstUseTSX(b *testing.B) {
	for range b.N {
		grammarCache = sync.Map{}
		registry := textmate.NewRegistry(textmate.RegistryOptions{LoadGrammar: Load})
		grammar, err := registry.LoadGrammar("source.tsx")
		if err != nil {
			b.Fatal(err)
		}
		registry.Dispose()
		runtime.KeepAlive(grammar)
	}
}
