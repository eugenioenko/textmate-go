package textmate

import (
	"fmt"
	"testing"
)

var (
	benchmarkDocumentResult LineResult
	benchmarkDocumentState  *StateStack
)

func BenchmarkDocumentWarmLine(b *testing.B) {
	document := NewDocument(newDocumentTestGrammar(), DocumentOptions{})
	document.SetLines([]string{`const message = "hello"`})
	benchmarkDocumentResult, _ = document.Line(0)
	b.ReportAllocs()
	b.ResetTimer()
	for range b.N {
		benchmarkDocumentResult, _ = document.Line(0)
	}
}

func BenchmarkDocumentSequential(b *testing.B) {
	lines := make([]string, 1_000)
	for index := range lines {
		lines[index] = fmt.Sprintf("const value%d = %d", index, index)
	}
	b.ReportAllocs()
	b.ResetTimer()
	for range b.N {
		document := NewDocument(newDocumentTestGrammar(), DocumentOptions{})
		document.SetLines(lines)
		benchmarkDocumentState, _ = document.StateAt(len(lines))
	}
}

func BenchmarkDocumentEditTailReuse(b *testing.B) {
	lines := make([]string, 1_000)
	for index := range lines {
		lines[index] = fmt.Sprintf("plain line %d", index)
	}
	document := NewDocument(newDocumentTestGrammar(), DocumentOptions{})
	document.SetLines(lines)
	benchmarkDocumentState, _ = document.StateAt(len(lines))

	b.ReportAllocs()
	b.ResetTimer()
	for iteration := range b.N {
		text := "changed even"
		if iteration%2 != 0 {
			text = "changed odd"
		}
		if err := document.ReplaceLines(0, 1, []string{text}); err != nil {
			b.Fatal(err)
		}
		benchmarkDocumentResult, _ = document.Line(0)
		benchmarkDocumentState, _ = document.StateAt(len(lines))
	}
}
