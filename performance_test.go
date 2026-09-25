package textmate_test

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	textmate "github.com/eugenioenko/textmate-go"
	"github.com/eugenioenko/textmate-go/grammars"
)

func BenchmarkTokenizeWarm(b *testing.B) {
	tests := []struct {
		name     string
		filename string
		fixture  string
		lines    int
	}{
		{name: "TSX_2000", filename: "component.tsx", fixture: "conformance/corpus/web/component.tsx", lines: 2000},
		{name: "HTML_2000", filename: "multiline.html", fixture: "conformance/corpus/web/multiline.html", lines: 2000},
		{name: "Go_2000", filename: "main.go", fixture: "conformance/corpus/general/main.go", lines: 2000},
		{name: "Markdown_1000", filename: "article.md", fixture: "conformance/corpus/markup/article.md", lines: 1000},
	}

	for _, test := range tests {
		b.Run(test.name, func(b *testing.B) {
			seed, err := os.ReadFile(filepath.FromSlash(test.fixture))
			if err != nil {
				b.Fatal(err)
			}
			lines := repeatLines(string(seed), test.lines)
			scopeName := grammars.ScopeForFilename(test.filename)
			if scopeName == "" {
				b.Fatalf("no embedded grammar for %s", test.filename)
			}
			registry := textmate.NewRegistry(textmate.RegistryOptions{LoadGrammar: grammars.Load})
			b.Cleanup(registry.Dispose)
			grammar, err := registry.LoadGrammar(scopeName)
			if err != nil {
				b.Fatal(err)
			}

			// Exercise lazy rule/regex compilation before measuring the warm path.
			var state *textmate.StateStack
			for _, line := range lines {
				state = grammar.TokenizeLine(line, state).RuleStack
			}

			var maxLine time.Duration
			b.ReportAllocs()
			b.ResetTimer()
			for range b.N {
				state = nil
				for _, line := range lines {
					started := time.Now()
					state = grammar.TokenizeLine(line, state).RuleStack
					maxLine = max(maxLine, time.Since(started))
				}
			}
			b.StopTimer()
			b.ReportMetric(float64(b.Elapsed().Nanoseconds())/float64(b.N*len(lines)), "ns/line")
			b.ReportMetric(float64(maxLine.Nanoseconds()), "max-ns/line")
			runtime.KeepAlive(state)
		})
	}
}

func repeatLines(contents string, count int) []string {
	seed := strings.Split(strings.TrimSuffix(contents, "\n"), "\n")
	result := make([]string, count)
	for index := range result {
		result[index] = seed[index%len(seed)]
	}
	return result
}
