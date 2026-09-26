package chromabench

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/chroma/v2/lexers"
	textmate "github.com/eugenioenko/textmate-go"
	"github.com/eugenioenko/textmate-go/grammars"
)

type corpusCase struct {
	Name      string `json:"name"`
	Filename  string `json:"filename"`
	Fixture   string `json:"fixture"`
	Lines     int    `json:"lines"`
	ScopeName string `json:"scopeName"`
}

// repoRoot is relative to this package; corpus.json paths are relative to it.
var repoRoot = filepath.Join("..", "..")

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

var benchmarkTokenCount int

const warmupPasses = 5

// BenchmarkWarmWholeDocument compares each engine in its correctness-preserving
// mode. TextMate carries state across physical lines. Chroma receives the whole
// document because released Chroma cannot return state for a later call.
func BenchmarkWarmWholeDocument(b *testing.B) {
	for _, test := range loadCorpus(b) {
		lines := loadRepeatedLines(b, test)
		document := strings.Join(lines, "\n") + "\n"
		b.SetBytes(int64(len(document)))

		b.Run(test.Name, func(b *testing.B) {
			b.Run("textmate-go", func(b *testing.B) {
				registry := textmate.NewRegistry(textmate.RegistryOptions{LoadGrammar: grammars.Load})
				b.Cleanup(registry.Dispose)
				grammar, err := registry.LoadGrammar(test.ScopeName)
				if err != nil {
					b.Fatal(err)
				}
				warmTextMate(b, grammar, lines)

				b.ReportAllocs()
				b.ResetTimer()
				for b.Loop() {
					benchmarkTokenCount = consumeTextMate(b, grammar, lines)
				}
				b.StopTimer()
				b.ReportMetric(float64(b.Elapsed().Nanoseconds())/float64(b.N*len(lines)), "ns/line")
			})

			b.Run("chroma-v2.24.1", func(b *testing.B) {
				lexer := chroma.Coalesce(lexers.Match(test.Filename))
				if lexer == nil {
					b.Fatalf("no Chroma lexer for %s", test.Filename)
				}
				warmChromaDocument(b, lexer, document)

				b.ReportAllocs()
				b.ResetTimer()
				for b.Loop() {
					benchmarkTokenCount = consumeChromaDocument(b, lexer, document)
				}
				b.StopTimer()
				b.ReportMetric(float64(b.Elapsed().Nanoseconds())/float64(b.N*len(lines)), "ns/line")
			})
		})
	}
}

// BenchmarkWarmLineByLine compares the tokenization primitive used by ttt's
// Chroma highlighter with textmate-go on identical lines. TextMate carries its
// immutable state between lines; released Chroma has no equivalent state API,
// so ttt starts Chroma at its root state for each line and separately patches
// a small set of multiline regions.
func BenchmarkWarmLineByLine(b *testing.B) {
	for _, test := range loadCorpus(b) {
		lines := loadRepeatedLines(b, test)
		b.Run(test.Name, func(b *testing.B) {
			b.Run("textmate-go", func(b *testing.B) {
				registry := textmate.NewRegistry(textmate.RegistryOptions{LoadGrammar: grammars.Load})
				b.Cleanup(registry.Dispose)
				grammar, err := registry.LoadGrammar(test.ScopeName)
				if err != nil {
					b.Fatal(err)
				}
				warmTextMate(b, grammar, lines)

				b.ReportAllocs()
				b.ResetTimer()
				for b.Loop() {
					benchmarkTokenCount = consumeTextMate(b, grammar, lines)
				}
				b.StopTimer()
				b.ReportMetric(float64(b.Elapsed().Nanoseconds())/float64(b.N*len(lines)), "ns/line")
			})

			b.Run("chroma-v2.24.1", func(b *testing.B) {
				lexer := chroma.Coalesce(lexers.Match(test.Filename))
				if lexer == nil {
					b.Fatalf("no Chroma lexer for %s", test.Filename)
				}
				warmChromaLines(b, lexer, lines)

				b.ReportAllocs()
				b.ResetTimer()
				for b.Loop() {
					benchmarkTokenCount = consumeChroma(b, lexer, lines)
				}
				b.StopTimer()
				b.ReportMetric(float64(b.Elapsed().Nanoseconds())/float64(b.N*len(lines)), "ns/line")
			})
		})
	}
}

func warmTextMate(b *testing.B, grammar *textmate.Grammar, lines []string) {
	b.Helper()
	for range warmupPasses {
		benchmarkTokenCount = consumeTextMate(b, grammar, lines)
	}
}

func warmChromaLines(b *testing.B, lexer chroma.Lexer, lines []string) {
	b.Helper()
	for range warmupPasses {
		benchmarkTokenCount = consumeChroma(b, lexer, lines)
	}
}

func warmChromaDocument(b *testing.B, lexer chroma.Lexer, document string) {
	b.Helper()
	for range warmupPasses {
		benchmarkTokenCount = consumeChromaDocument(b, lexer, document)
	}
}

func consumeTextMate(b *testing.B, grammar *textmate.Grammar, lines []string) int {
	b.Helper()
	var state *textmate.StateStack
	tokenCount := 0
	for _, line := range lines {
		result := grammar.TokenizeLine(line, state)
		if result.Stopped {
			b.Fatal("textmate-go stopped while tokenizing benchmark corpus")
		}
		state = result.RuleStack
		tokenCount += len(result.Tokens)
	}
	runtime.KeepAlive(state)
	return tokenCount
}

func consumeChroma(b *testing.B, lexer chroma.Lexer, lines []string) int {
	b.Helper()
	tokenCount := 0
	for _, line := range lines {
		iterator, err := lexer.Tokenise(nil, line+"\n")
		if err != nil {
			b.Fatal(err)
		}
		tokenCount += len(iterator.Tokens())
	}
	return tokenCount
}

func consumeChromaDocument(b *testing.B, lexer chroma.Lexer, document string) int {
	b.Helper()
	iterator, err := lexer.Tokenise(nil, document)
	if err != nil {
		b.Fatal(err)
	}
	return len(iterator.Tokens())
}

func loadRepeatedLines(b *testing.B, test corpusCase) []string {
	b.Helper()
	data, err := os.ReadFile(filepath.Join(repoRoot, filepath.FromSlash(test.Fixture)))
	if err != nil {
		b.Fatal(err)
	}
	seed := strings.Split(strings.TrimSuffix(string(data), "\n"), "\n")
	lines := make([]string, test.Lines)
	for index := range lines {
		lines[index] = seed[index%len(seed)]
	}
	return lines
}
