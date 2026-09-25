package chromabench

import (
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
	name     string
	filename string
	fixture  string
	lines    int
}

var corpusCases = []corpusCase{
	{name: "TSX_2000", filename: "component.tsx", fixture: "conformance/corpus/web/component.tsx", lines: 2000},
	{name: "HTML_2000", filename: "multiline.html", fixture: "conformance/corpus/web/multiline.html", lines: 2000},
	{name: "Go_2000", filename: "main.go", fixture: "conformance/corpus/general/main.go", lines: 2000},
	{name: "Markdown_1000", filename: "article.md", fixture: "conformance/corpus/markup/article.md", lines: 1000},
}

var benchmarkTokenCount int

const warmupPasses = 5

// BenchmarkWarmWholeDocument compares each engine in its correctness-preserving
// mode. TextMate carries state across physical lines. Chroma receives the whole
// document because released Chroma cannot return state for a later call.
func BenchmarkWarmWholeDocument(b *testing.B) {
	for _, test := range corpusCases {
		lines := loadRepeatedLines(b, test)
		document := strings.Join(lines, "\n") + "\n"
		b.SetBytes(int64(len(document)))

		b.Run(test.name, func(b *testing.B) {
			b.Run("textmate-go", func(b *testing.B) {
				registry := textmate.NewRegistry(textmate.RegistryOptions{LoadGrammar: grammars.Load})
				b.Cleanup(registry.Dispose)
				grammar, err := registry.LoadGrammar(grammars.ScopeForFilename(test.filename))
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
				lexer := chroma.Coalesce(lexers.Match(test.filename))
				if lexer == nil {
					b.Fatalf("no Chroma lexer for %s", test.filename)
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
	for _, test := range corpusCases {
		lines := loadRepeatedLines(b, test)
		b.Run(test.name, func(b *testing.B) {
			b.Run("textmate-go", func(b *testing.B) {
				registry := textmate.NewRegistry(textmate.RegistryOptions{LoadGrammar: grammars.Load})
				b.Cleanup(registry.Dispose)
				grammar, err := registry.LoadGrammar(grammars.ScopeForFilename(test.filename))
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
				lexer := chroma.Coalesce(lexers.Match(test.filename))
				if lexer == nil {
					b.Fatalf("no Chroma lexer for %s", test.filename)
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
	root := filepath.Clean(filepath.Join("..", ".."))
	data, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(test.fixture)))
	if err != nil {
		b.Fatal(err)
	}
	seed := strings.Split(strings.TrimSuffix(string(data), "\n"), "\n")
	lines := make([]string, test.lines)
	for index := range lines {
		lines[index] = seed[index%len(seed)]
	}
	return lines
}
