package textmate

import (
	"bytes"
	"compress/gzip"
	_ "embed"
	"io"
	"testing"
)

//go:embed grammars/data/swift.json.gz
var swiftGrammarGZIP []byte

func TestSwiftRealGrammarClassifiesTrailingCustomOperatorAsPostfix(t *testing.T) {
	reader, err := gzip.NewReader(bytes.NewReader(swiftGrammarGZIP))
	if err != nil {
		t.Fatal(err)
	}
	data, err := io.ReadAll(reader)
	if err != nil {
		t.Fatal(err)
	}
	if err := reader.Close(); err != nil {
		t.Fatal(err)
	}
	raw, err := ParseRawGrammar(data)
	if err != nil {
		t.Fatal(err)
	}
	grammar := newGrammar(raw.ScopeName, raw, nil)
	result := grammar.TokenizeLine("x=U>", nil)
	for _, token := range result.Tokens {
		if token.Start == 3 && token.End == 4 {
			if got := token.Scopes[len(token.Scopes)-1]; got != "keyword.operator.custom.postfix.swift" {
				t.Fatalf("operator scope = %q, want postfix; tokens=%+v", got, result.Tokens)
			}
			return
		}
	}
	t.Fatalf("operator token not found: %+v", result.Tokens)
}
