package textmate_test

import (
	"fmt"
	"strings"
	"testing"
	"time"

	textmate "github.com/eugenioenko/textmate-go"
	"github.com/eugenioenko/textmate-go/grammars"
)

func wideTypeScriptFields(count int) string {
	var line strings.Builder
	line.WriteString("'T': { kind: 'OBJECT'; name: 'T'; fields: { ")
	for i := range count {
		fmt.Fprintf(&line, "'f%d': { name: 'f%d'; type: { kind: 'NON_NULL'; name: never; ofType: { kind: 'SCALAR'; name: 'String'; ofType: null; } }; }; ", i, i)
	}
	line.WriteString("}; };")
	return line.String()
}

func TestTypeScriptWideNestedDeclaration(t *testing.T) {
	registry := textmate.NewRegistry(textmate.RegistryOptions{LoadGrammar: grammars.Load})
	t.Cleanup(registry.Dispose)
	grammar, err := registry.LoadGrammar("source.ts")
	if err != nil {
		t.Fatal(err)
	}
	doc := textmate.NewDocument(grammar, textmate.DocumentOptions{
		TokenizeOptions: textmate.TokenizeOptions{TimeLimit: time.Second},
	})
	wide := wideTypeScriptFields(60)
	doc.SetLines([]string{
		"export type introspection_types = {",
		wide,
		"};",
		"const handler = (value: string) => value;",
	})
	for pass := range 2 {
		for i := range doc.Len() {
			result, ok := doc.Line(i)
			if !ok || result.Stopped {
				t.Fatalf("pass %d line %d: tokenization incomplete: %v", pass, i, result.StoppedReason)
			}
			if i == 1 {
				foundLastField := false
				line := []rune(wide)
				for _, token := range result.Tokens {
					if strings.Contains(string(line[token.Start:token.End]), "f59") {
						foundLastField = true
					}
				}
				if !foundLastField {
					t.Fatal("last field was not tokenized")
				}
			}
		}
	}
}

func BenchmarkTypeScriptWideNestedDeclaration(b *testing.B) {
	for _, fields := range []int{10, 30, 60} {
		b.Run(fmt.Sprintf("fields_%d", fields), func(b *testing.B) {
			registry := textmate.NewRegistry(textmate.RegistryOptions{LoadGrammar: grammars.Load})
			b.Cleanup(registry.Dispose)
			grammar, err := registry.LoadGrammar("source.ts")
			if err != nil {
				b.Fatal(err)
			}
			state := grammar.TokenizeLine("export type introspection_types = {", nil).RuleStack
			line := wideTypeScriptFields(fields)
			grammar.TokenizeLine(line, state)
			b.ReportAllocs()
			b.SetBytes(int64(len(line)))
			b.ResetTimer()
			for range b.N {
				grammar.TokenizeLine(line, state)
			}
		})
	}
}
