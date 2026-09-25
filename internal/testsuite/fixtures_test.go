package testsuite

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	textmate "github.com/eugenioenko/textmate-go"
)

type fixtureCase struct {
	Description       string        `json:"desc"`
	Grammars          []string      `json:"grammars"`
	GrammarPath       string        `json:"grammarPath"`
	GrammarScopeName  string        `json:"grammarScopeName"`
	GrammarInjections []string      `json:"grammarInjections"`
	Lines             []fixtureLine `json:"lines"`
}

type fixtureLine struct {
	Line   string         `json:"line"`
	Tokens []fixtureToken `json:"tokens"`
}

type fixtureToken struct {
	Value  string   `json:"value"`
	Scopes []string `json:"scopes"`
}

type suiteResult struct {
	passed, failed, skipped int
	failures                []string
}

func TestTokenizationFixtures(t *testing.T) {
	checkout := vscodeTextmateCheckout(t)
	for _, suite := range []struct{ name, manifest string }{
		{"first-mate", "test-cases/first-mate/tests.json"},
		{"suite1", "test-cases/suite1/tests.json"},
		{"suite1-while", "test-cases/suite1/whileTests.json"},
	} {
		suite := suite
		t.Run(suite.name, func(t *testing.T) {
			result, err := runFixtureSuite(filepath.Join(checkout, suite.manifest))
			if err != nil {
				t.Fatal(err)
			}
			t.Logf("%s: %d passed, %d failed, %d skipped (plist)", suite.manifest, result.passed, result.failed, result.skipped)
			if result.failed == 0 {
				return
			}
			message := strings.Join(result.failures, "\n")
			t.Fatalf("%d fixture cases failed:\n%s", result.failed, message)
		})
	}
}

func vscodeTextmateCheckout(t *testing.T) string {
	t.Helper()
	if configured := os.Getenv("VSCODE_TEXTMATE_DIR"); configured != "" {
		if path, err := filepath.Abs(configured); err == nil {
			return path
		}
		return configured
	}
	dir, err := os.Getwd()
	if err != nil {
		t.Skipf("cannot locate vscode-textmate checkout: %v", err)
	}
	for {
		candidate := filepath.Join(filepath.Dir(dir), "vscode-textmate")
		if _, err := os.Stat(filepath.Join(candidate, "test-cases", "first-mate", "tests.json")); err == nil {
			return candidate
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	t.Skip("set VSCODE_TEXTMATE_DIR to run vscode-textmate fixtures")
	return ""
}

func runFixtureSuite(manifestPath string) (suiteResult, error) {
	data, err := os.ReadFile(manifestPath)
	if err != nil {
		return suiteResult{}, fmt.Errorf("read fixture manifest: %w", err)
	}
	var cases []fixtureCase
	if err := json.Unmarshal(data, &cases); err != nil {
		return suiteResult{}, fmt.Errorf("parse fixture manifest: %w", err)
	}

	var result suiteResult
	for index, test := range cases {
		if hasPlistGrammar(test.Grammars) {
			result.skipped++
			continue
		}
		if err := runFixtureCase(filepath.Dir(manifestPath), test); err != nil {
			result.failed++
			result.failures = append(result.failures, fmt.Sprintf("[%d] %s: %v", index+1, test.Description, err))
			continue
		}
		result.passed++
	}
	return result, nil
}

func hasPlistGrammar(paths []string) bool {
	for _, path := range paths {
		if !strings.HasSuffix(strings.ToLower(path), ".json") {
			return true
		}
	}
	return false
}

func runFixtureCase(baseDir string, test fixtureCase) error {
	grammars := make(map[string]*textmate.RawGrammar, len(test.Grammars))
	targetScope := test.GrammarScopeName
	for _, relativePath := range test.Grammars {
		data, err := os.ReadFile(filepath.Join(baseDir, relativePath))
		if err != nil {
			return fmt.Errorf("read grammar %q: %w", relativePath, err)
		}
		raw, err := textmate.ParseRawGrammar(data)
		if err != nil {
			return fmt.Errorf("parse grammar %q: %w", relativePath, err)
		}
		grammars[raw.ScopeName] = raw
		if targetScope == "" && relativePath == test.GrammarPath {
			targetScope = raw.ScopeName
		}
	}
	if targetScope == "" {
		return fmt.Errorf("fixture does not identify its root grammar")
	}

	registry := textmate.NewRegistry(textmate.RegistryOptions{
		LoadGrammar: func(scopeName string) (*textmate.RawGrammar, error) { return grammars[scopeName], nil },
		GetInjections: func(scopeName string) []string {
			if scopeName == targetScope {
				return test.GrammarInjections
			}
			return nil
		},
	})
	grammar, err := registry.LoadGrammar(targetScope)
	if err != nil {
		return fmt.Errorf("load grammar %q: %w", targetScope, err)
	}

	var state *textmate.StateStack
	for lineIndex, line := range test.Lines {
		actual := grammar.TokenizeLine(line.Line, state)
		state = actual.RuleStack
		runes := []rune(line.Line)
		tokens := make([]fixtureToken, 0, len(actual.Tokens))
		for _, token := range actual.Tokens {
			if token.Start < 0 || token.End < token.Start || token.End > len(runes) {
				return fmt.Errorf("line %d returned invalid rune range [%d,%d) for length %d", lineIndex+1, token.Start, token.End, len(runes))
			}
			tokens = append(tokens, fixtureToken{Value: string(runes[token.Start:token.End]), Scopes: token.Scopes})
		}
		expected := line.Tokens
		if len(runes) != 0 {
			expected = removeEmptyTokens(expected)
		}
		if !reflect.DeepEqual(tokens, expected) {
			return fmt.Errorf("line %d %q: %s", lineIndex+1, line.Line, firstDifference(expected, tokens))
		}
	}
	return nil
}

func removeEmptyTokens(tokens []fixtureToken) []fixtureToken {
	result := make([]fixtureToken, 0, len(tokens))
	for _, token := range tokens {
		if token.Value != "" {
			result = append(result, token)
		}
	}
	return result
}

func firstDifference(want, got []fixtureToken) string {
	limit := min(len(want), len(got))
	for index := 0; index < limit; index++ {
		if !reflect.DeepEqual(want[index], got[index]) {
			return fmt.Sprintf("token %d = %#v, want %#v", index+1, got[index], want[index])
		}
	}
	return fmt.Sprintf("got %d tokens, want %d", len(got), len(want))
}
