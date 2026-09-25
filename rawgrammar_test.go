package textmate

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestParseRawGrammar(t *testing.T) {
	data := []byte(`{
		"scopeName": "source.example",
		"name": "Example",
		"fileTypes": ["example", "ex"],
		"firstLineMatch": "^#!.*example",
		"injectionSelector": "L:source.example",
		"$vscodeTextmateLocation": {"filename":"example.json","line":1,"char":2},
		"patterns": [
			{"include":"#dynamic-key"},
			{
				"begin":"(a)",
				"beginCaptures":{"1":{"name":"punctuation.begin.example"}},
				"end":"(z)",
				"endCaptures":{"1":{"name":"punctuation.end.example"}},
				"while":"continued",
				"whileCaptures":{"0":{"name":"meta.continued.example"}},
				"contentName":"meta.content.example",
				"applyEndPatternLast":1,
				"repository":{"nested":{"match":"nested"}},
				"patterns":[{"match":"body","captures":{"0":{"name":"body.example"}}}],
				"$vscodeTextmateLocation":{"filename":"example.json","line":8,"char":4}
			}
		],
		"repository": {
			"dynamic-key": {"match":"word","name":"keyword.example"},
			"key.with punctuation": {"include":"$self"}
		},
		"injections": {"L:comment": {"match":"TODO", "name":"keyword.todo.example"}}
	}`)

	grammar, err := ParseRawGrammar(data)
	if err != nil {
		t.Fatalf("ParseRawGrammar() error = %v", err)
	}
	if grammar.ScopeName != "source.example" || grammar.Name != "Example" {
		t.Fatalf("grammar identity = (%q, %q)", grammar.ScopeName, grammar.Name)
	}
	if len(grammar.FileTypes) != 2 || grammar.FirstLineMatch != "^#!.*example" {
		t.Fatalf("grammar metadata not decoded: %#v", grammar)
	}
	if grammar.Location == nil || grammar.Location.Filename != "example.json" || grammar.Location.Char != 2 {
		t.Fatalf("grammar location = %#v", grammar.Location)
	}
	if got := grammar.Repository["dynamic-key"]; got == nil || stringValue(got.Match) != "word" {
		t.Fatalf("dynamic repository rule = %#v", got)
	}
	if got := grammar.Repository["key.with punctuation"]; got == nil || stringValue(got.Include) != "$self" {
		t.Fatalf("punctuated repository rule = %#v", got)
	}
	if got := grammar.Injections["L:comment"]; got == nil || stringValue(got.Name) != "keyword.todo.example" {
		t.Fatalf("injection = %#v", got)
	}

	rule := grammar.Patterns[1]
	if !rule.ApplyEndPatternLast {
		t.Fatal("numeric applyEndPatternLast was not decoded as true")
	}
	if stringValue(rule.BeginCaptures["1"].Name) != "punctuation.begin.example" ||
		stringValue(rule.EndCaptures["1"].Name) != "punctuation.end.example" ||
		stringValue(rule.WhileCaptures["0"].Name) != "meta.continued.example" ||
		stringValue(rule.Patterns[0].Captures["0"].Name) != "body.example" {
		t.Fatalf("capture maps not decoded: %#v", rule)
	}
	if got := rule.Repository["nested"]; got == nil || stringValue(got.Match) != "nested" {
		t.Fatalf("nested repository = %#v", got)
	}
	if rule.Location == nil || rule.Location.Line != 8 || rule.Location.Char != 4 {
		t.Fatalf("rule location = %#v", rule.Location)
	}
}

func TestRawRuleApplyEndPatternLastForms(t *testing.T) {
	tests := []struct {
		json string
		want bool
	}{
		{json: `{"applyEndPatternLast":true}`, want: true},
		{json: `{"applyEndPatternLast":false}`, want: false},
		{json: `{"applyEndPatternLast":1}`, want: true},
		{json: `{"applyEndPatternLast":0}`, want: false},
		{json: `{}`, want: false},
	}

	for _, test := range tests {
		var rule RawRule
		if err := json.Unmarshal([]byte(test.json), &rule); err != nil {
			t.Fatalf("json.Unmarshal(%s) error = %v", test.json, err)
		}
		if got := rule.ApplyEndPatternLast; got != test.want {
			t.Errorf("json.Unmarshal(%s) applyEndPatternLast = %v, want %v", test.json, got, test.want)
		}
	}
}

func TestRawCapturesAcceptsConvertedGrammarForms(t *testing.T) {
	tests := []struct {
		name string
		json string
		want map[string]string
	}{
		{
			name: "object drops malformed entries",
			json: `{
				"0":{"name":"valid.zero"},
				"1":"invalid shorthand",
				"2":[{"name":"invalid.array"}],
				"name":"misplaced rule name"
			}`,
			want: map[string]string{"0": "valid.zero"},
		},
		{
			name: "array uses numeric indexes",
			json: `[
				{"name":"valid.zero"},
				{"name":"valid.one"},
				"invalid shorthand"
			]`,
			want: map[string]string{"0": "valid.zero", "1": "valid.one"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var captures RawCaptures
			if err := json.Unmarshal([]byte(test.json), &captures); err != nil {
				t.Fatalf("json.Unmarshal() error = %v", err)
			}
			if len(captures) != len(test.want) {
				t.Fatalf("len(captures) = %d, want %d: %#v", len(captures), len(test.want), captures)
			}
			for captureID, wantName := range test.want {
				if got := captures[captureID]; got == nil || stringValue(got.Name) != wantName {
					t.Errorf("captures[%q] = %#v, want name %q", captureID, got, wantName)
				}
			}
		})
	}
}

func TestRawRulePreservesAbsentAndEmptyStrings(t *testing.T) {
	var rules []*RawRule
	if err := json.Unmarshal([]byte(`[
		{},
		{"include":"","name":"","contentName":"","match":"","begin":"","end":"","while":""}
	]`), &rules); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}
	if rules[0].Include != nil || rules[0].Name != nil || rules[0].ContentName != nil ||
		rules[0].Match != nil || rules[0].Begin != nil || rules[0].End != nil || rules[0].While != nil {
		t.Fatalf("absent optional strings decoded as present: %#v", rules[0])
	}
	if rules[1].Include == nil || rules[1].Name == nil || rules[1].ContentName == nil ||
		rules[1].Match == nil || rules[1].Begin == nil || rules[1].End == nil || rules[1].While == nil {
		t.Fatalf("present empty optional strings decoded as absent: %#v", rules[1])
	}
}

func TestRawRepositoryDropsMalformedEntries(t *testing.T) {
	var repository RawRepository
	if err := json.Unmarshal([]byte(`{
		"valid":{"match":"word"},
		"array":[{"match":"ignored"}],
		"string":"ignored"
	}`), &repository); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}
	if len(repository) != 1 || repository["valid"] == nil || stringValue(repository["valid"].Match) != "word" {
		t.Fatalf("repository = %#v", repository)
	}
}

func TestParseRawGrammarErrors(t *testing.T) {
	tests := []struct {
		name string
		data string
		want string
	}{
		{name: "invalid JSON", data: `{`, want: "parse raw grammar"},
		{name: "missing scope", data: `{"patterns":[]}`, want: "missing scopeName"},
		{name: "invalid flexible boolean", data: `{"scopeName":"source.bad","patterns":[{"applyEndPatternLast":"yes"}]}`, want: "applyEndPatternLast"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := ParseRawGrammar([]byte(test.data))
			if err == nil || !strings.Contains(err.Error(), test.want) {
				t.Fatalf("ParseRawGrammar() error = %v, want containing %q", err, test.want)
			}
		})
	}
}

func stringValue[T ~string](value *T) string {
	if value == nil {
		return ""
	}
	return string(*value)
}

func TestParseRawGrammarReferenceCorpora(t *testing.T) {
	roots := []string{
		"../tm-grammars/packages/tm-grammars/grammars",
		"../vscode-textmate/test-cases",
	}

	parsed := 0
	for _, root := range roots {
		if _, err := os.Stat(root); err != nil {
			if os.IsNotExist(err) {
				continue
			}
			t.Fatalf("os.Stat(%q) error = %v", root, err)
		}
		err := filepath.WalkDir(root, func(path string, entry os.DirEntry, walkErr error) error {
			if walkErr != nil {
				return walkErr
			}
			if entry.IsDir() || filepath.Ext(path) != ".json" {
				return nil
			}
			// tests.json and similar files are fixture manifests, not grammars.
			data, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			var identity struct {
				ScopeName string `json:"scopeName"`
			}
			if json.Unmarshal(data, &identity) != nil || identity.ScopeName == "" {
				return nil
			}
			if _, err := ParseRawGrammar(data); err != nil {
				t.Errorf("ParseRawGrammar(%q) error = %v", path, err)
			}
			parsed++
			return nil
		})
		if err != nil {
			t.Fatalf("filepath.WalkDir(%q) error = %v", root, err)
		}
	}
	if parsed == 0 {
		t.Skip("reference grammar repositories are not cloned next to this repository")
	}
	t.Logf("parsed %d JSON TextMate grammars from reference repositories", parsed)
}
