package textmatebench

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	textmate "github.com/eugenioenko/textmate-go"
	"github.com/eugenioenko/textmate-go/grammars"
)

type corpusCase struct {
	Name      string `json:"name"`
	Group     string `json:"group"`
	Filename  string `json:"filename"`
	Fixture   string `json:"fixture"`
	Lines     int    `json:"lines"`
	ScopeName string `json:"scopeName"`
}

var benchmarkTokenCount int

const (
	repoRoot     = "../.."
	warmupPasses = 5
)

func loadCorpus(tb testing.TB) []corpusCase {
	tb.Helper()
	data, err := os.ReadFile(filepath.Join(repoRoot, "benchmarks", "corpus.json"))
	if err != nil {
		tb.Fatal(err)
	}
	var manifest struct {
		Cases []corpusCase `json:"cases"`
	}
	if err := json.Unmarshal(data, &manifest); err != nil {
		tb.Fatal(err)
	}
	return manifest.Cases
}

func TestCorpusManifest(t *testing.T) {
	for _, test := range loadCorpus(t) {
		if test.Group != "core" && test.Group != "extended" {
			t.Errorf("%s: unknown group %q", test.Name, test.Group)
		}
		if got := grammars.ScopeForFilename(test.Filename); got != test.ScopeName {
			t.Errorf("%s: scope %q, grammars resolves %q", test.Name, test.ScopeName, got)
		}
		if _, err := os.Stat(filepath.Join(repoRoot, filepath.FromSlash(test.Fixture))); err != nil {
			t.Errorf("%s: %v", test.Name, err)
		}
		if test.Lines <= 0 {
			t.Errorf("%s: lines must be positive", test.Name)
		}
	}
}

func BenchmarkWarmLineByLine(b *testing.B) {
	for _, test := range loadCorpus(b) {
		b.Run(test.Name, func(b *testing.B) {
			lines := loadRepeatedLines(b, test)
			registry := textmate.NewRegistry(textmate.RegistryOptions{LoadGrammar: grammars.Load})
			b.Cleanup(registry.Dispose)
			grammar, err := registry.LoadGrammar(test.ScopeName)
			if err != nil {
				b.Fatal(err)
			}
			for range warmupPasses {
				benchmarkTokenCount = consume(b, grammar, lines)
			}

			b.ReportAllocs()
			b.ResetTimer()
			for b.Loop() {
				benchmarkTokenCount = consume(b, grammar, lines)
			}
			b.StopTimer()
			b.ReportMetric(float64(b.Elapsed().Nanoseconds())/float64(b.N*len(lines)), "ns/line")
		})
	}
}

func loadRepeatedLines(tb testing.TB, test corpusCase) []string {
	tb.Helper()
	data, err := os.ReadFile(filepath.Join(repoRoot, filepath.FromSlash(test.Fixture)))
	if err != nil {
		tb.Fatal(err)
	}
	seed := strings.Split(strings.TrimSuffix(string(data), "\n"), "\n")
	lines := make([]string, test.Lines)
	for index := range lines {
		lines[index] = seed[index%len(seed)]
	}
	return lines
}

func consume(tb testing.TB, grammar *textmate.Grammar, lines []string) int {
	tb.Helper()
	var state *textmate.StateStack
	tokenCount := 0
	for _, line := range lines {
		result := grammar.TokenizeLine(line, state)
		if result.Stopped {
			tb.Fatal("textmate-go stopped while tokenizing benchmark corpus")
		}
		state = result.RuleStack
		tokenCount += len(result.Tokens)
	}
	runtime.KeepAlive(state)
	return tokenCount
}
