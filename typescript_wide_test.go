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

// A function-typed field puts '=' and '=>' on the line, so required-content
// bounds no longer reject the field-declaration lookahead from earlier fields.
func TestTypeScriptWideDeclarationWithFunctionField(t *testing.T) {
	registry := textmate.NewRegistry(textmate.RegistryOptions{LoadGrammar: grammars.Load})
	t.Cleanup(registry.Dispose)
	grammar, err := registry.LoadGrammar("source.ts")
	if err != nil {
		t.Fatal(err)
	}
	budget := time.Second
	if raceEnabled {
		budget *= 10
	}
	for _, fields := range []int{60, 240} {
		wide := strings.TrimSuffix(wideTypeScriptFields(fields), "}; };") + "cb: (x: number) => void; }; };"
		doc := textmate.NewDocument(grammar, textmate.DocumentOptions{
			TokenizeOptions: textmate.TokenizeOptions{TimeLimit: budget},
		})
		doc.SetLines([]string{"export type introspection_types = {", wide, "};"})
		result, ok := doc.Line(1)
		if !ok || result.Stopped {
			t.Fatalf("%d fields: tokenization incomplete: %v", fields, result.StoppedReason)
		}
		line := []rune(wide)
		functionField := false
		for _, token := range result.Tokens {
			if string(line[token.Start:token.End]) != "cb" {
				continue
			}
			for _, scope := range token.Scopes {
				if strings.HasPrefix(scope, "entity.name.function") {
					functionField = true
				}
			}
		}
		if !functionField {
			t.Fatalf("%d fields: cb was not highlighted as a function", fields)
		}
	}
}
