package grammars

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"unicode/utf8"

	textmate "github.com/eugenioenko/textmate-go"
)

const sourceRevision = "37edd1b26f18838050661d912334aba0ca7f4931"

type corpusManifest struct {
	Files []corpusFile `json:"files"`
}

type corpusFile struct {
	File    string `json:"file"`
	Grammar string `json:"grammar"`
	Source  string `json:"source"`
}

func TestEmbeddedCorpus(t *testing.T) {
	if os.Getenv("TEXTMATE_GO_REQUIRE_EMBEDDED_CORPUS") != "1" {
		t.Skip("set TEXTMATE_GO_REQUIRE_EMBEDDED_CORPUS=1 to run the embedded production corpus")
	}

	repositoryRoot, err := filepath.Abs("..")
	if err != nil {
		t.Fatal(err)
	}
	packageRoot := os.Getenv("TM_GRAMMARS_DIR")
	if packageRoot == "" {
		packageRoot = filepath.Join(repositoryRoot, "..", "tm-grammars", "packages", "tm-grammars")
	} else if !filepath.IsAbs(packageRoot) {
		packageRoot = filepath.Join(repositoryRoot, packageRoot)
	}
	packageRoot, err = filepath.Abs(packageRoot)
	if err != nil {
		t.Fatal(err)
	}
	checkoutRoot := filepath.Clean(filepath.Join(packageRoot, "..", ".."))
	verifySourceCheckout(t, checkoutRoot)

	manifestData, err := os.ReadFile(filepath.Join(repositoryRoot, "conformance", "corpus", "manifest.json"))
	if err != nil {
		t.Fatal(err)
	}
	var manifest corpusManifest
	if err := json.Unmarshal(manifestData, &manifest); err != nil {
		t.Fatal(err)
	}
	if len(manifest.Files) != 166 {
		t.Fatalf("corpus entries = %d, want 166", len(manifest.Files))
	}
	sourceGrammars := loadSourceGrammars(t, filepath.Join(packageRoot, "grammars"))

	for _, fixture := range manifest.Files {
		fixture := fixture
		t.Run(strings.ReplaceAll(fixture.File, "/", "_"), func(t *testing.T) {
			scopeName := scopeForGrammarID(fixture.Grammar)
			if scopeName == "" {
				t.Fatalf("corpus grammar %q is unexpectedly not embedded", fixture.Grammar)
			}
			filename := filepath.Join(repositoryRoot, "conformance", "corpus", fixture.File)
			if fixture.Source == "tm-grammars" {
				filename = filepath.Join(checkoutRoot, fixture.File)
			} else if fixture.Source != "" {
				t.Fatalf("unknown corpus source %q", fixture.Source)
			}
			contents, err := os.ReadFile(filename)
			if err != nil {
				t.Fatal(err)
			}

			registry := textmate.NewRegistry(textmate.RegistryOptions{LoadGrammar: Load})
			defer registry.Dispose()
			grammar, err := registry.LoadGrammar(scopeName)
			if err != nil {
				t.Fatal(err)
			}
			sourceRegistry := textmate.NewRegistry(textmate.RegistryOptions{
				LoadGrammar: func(scopeName string) (*textmate.RawGrammar, error) {
					return sourceGrammars[scopeName], nil
				},
			})
			defer sourceRegistry.Dispose()
			sourceGrammar, err := sourceRegistry.LoadGrammar(scopeName)
			if err != nil {
				t.Fatal(err)
			}

			var state, sourceState *textmate.StateStack
			for lineIndex, line := range splitLines(string(contents)) {
				result := grammar.TokenizeLine(line, state)
				sourceResult := sourceGrammar.TokenizeLine(line, sourceState)
				assertContiguousTokens(t, fixture.File, lineIndex+1, line, result.Tokens)
				if !sameTokens(result.Tokens, sourceResult.Tokens) {
					t.Fatalf(
						"%s:%d embedded tokens differ from the complete source repository\nembedded: %#v\nsource:   %#v",
						fixture.File,
						lineIndex+1,
						result.Tokens,
						sourceResult.Tokens,
					)
				}
				if result.RuleStack == nil {
					t.Fatalf("%s:%d returned nil state", fixture.File, lineIndex+1)
				}
				state = result.RuleStack
				sourceState = sourceResult.RuleStack
			}
		})
	}
}

func loadSourceGrammars(t *testing.T, directory string) map[string]*textmate.RawGrammar {
	t.Helper()
	paths, err := filepath.Glob(filepath.Join(directory, "*.json"))
	if err != nil {
		t.Fatal(err)
	}
	result := make(map[string]*textmate.RawGrammar, len(paths))
	for _, path := range paths {
		data, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		grammar, err := textmate.ParseRawGrammar(data)
		if err != nil {
			t.Fatalf("parse %s: %v", path, err)
		}
		result[grammar.ScopeName] = grammar
	}
	if len(result) != 260 {
		t.Fatalf("source grammars = %d, want 260", len(result))
	}
	return result
}

func TestEmbeddedExternalDependenciesAreDocumented(t *testing.T) {
	var missing []string
	for _, scopeName := range Scopes() {
		grammar, err := Load(scopeName)
		if err != nil {
			t.Fatal(err)
		}
		walkRawRules(grammar.Patterns, grammar.Repository, func(include string) {
			external := strings.SplitN(include, "#", 2)[0]
			if external == "" || external == "$self" || external == "$base" {
				return
			}
			if embedded, _ := Load(external); embedded == nil {
				missing = append(missing, scopeName+" -> "+external)
			}
		})
	}
	slices.Sort(missing)
	missing = slices.Compact(missing)
	data, err := os.ReadFile("OPTIONAL_DEPENDENCIES")
	if err != nil {
		t.Fatal(err)
	}
	var documented []string
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(strings.SplitN(line, "#", 2)[0])
		if line != "" {
			documented = append(documented, line)
		}
	}
	if !slices.Equal(missing, documented) {
		t.Fatalf(
			"external dependency audit changed\ngot:\n%s\nwant OPTIONAL_DEPENDENCIES:\n%s",
			strings.Join(missing, "\n"),
			strings.Join(documented, "\n"),
		)
	}
}

func scopeForGrammarID(id string) string {
	for scopeName, asset := range generatedAssets {
		if strings.TrimSuffix(filepath.Base(asset), ".json.gz") == id {
			return scopeName
		}
	}
	return ""
}

func splitLines(contents string) []string {
	contents = strings.ReplaceAll(contents, "\r\n", "\n")
	contents = strings.ReplaceAll(contents, "\r", "\n")
	return strings.Split(contents, "\n")
}

func assertContiguousTokens(t *testing.T, filename string, lineNumber int, line string, tokens []textmate.Token) {
	t.Helper()
	wantEnd := utf8.RuneCountInString(line)
	position := 0
	for index, token := range tokens {
		if token.Start != position || token.End < token.Start || token.End > wantEnd {
			t.Fatalf("%s:%d token %d invalid after %d for rune length %d: %#v", filename, lineNumber, index, position, wantEnd, tokens)
		}
		if len(token.Scopes) == 0 {
			t.Fatalf("%s:%d token %d has no scopes", filename, lineNumber, index)
		}
		position = token.End
	}
	if position != wantEnd {
		t.Fatalf("%s:%d tokens end at %d, want %d: %#v", filename, lineNumber, position, wantEnd, tokens)
	}
}

func walkRawRules(patterns []*textmate.RawRule, repository textmate.RawRepository, visit func(string)) {
	seen := make(map[*textmate.RawRule]bool)
	var walk func(*textmate.RawRule)
	walk = func(rule *textmate.RawRule) {
		if rule == nil || seen[rule] {
			return
		}
		seen[rule] = true
		if rule.Include != nil {
			visit(*rule.Include)
		}
		for _, child := range rule.Patterns {
			walk(child)
		}
		for _, child := range rule.Repository {
			walk(child)
		}
		for _, captures := range []textmate.RawCaptures{
			rule.Captures,
			rule.BeginCaptures,
			rule.EndCaptures,
			rule.WhileCaptures,
		} {
			for _, child := range captures {
				walk(child)
			}
		}
	}
	for _, rule := range patterns {
		walk(rule)
	}
	for _, rule := range repository {
		walk(rule)
	}
}

func verifySourceCheckout(t *testing.T, root string) {
	t.Helper()
	command := exec.Command("git", "-C", root, "rev-parse", "HEAD")
	output, err := command.Output()
	if err != nil {
		t.Fatalf("read tm-grammars revision: %v", err)
	}
	if got := strings.TrimSpace(string(output)); got != sourceRevision {
		t.Fatalf("tm-grammars revision = %s, want %s", got, sourceRevision)
	}
	command = exec.Command("git", "-C", root, "status", "--porcelain")
	output, err = command.Output()
	if err != nil {
		t.Fatalf("read tm-grammars status: %v", err)
	}
	if len(output) != 0 {
		t.Fatal("tm-grammars checkout is dirty")
	}
}

// sameTokens ignores ScopeStack, whose interned pointers differ between the
// two registries being compared.
func sameTokens(a, b []textmate.Token) bool {
	return slices.EqualFunc(a, b, func(x, y textmate.Token) bool {
		return x.Start == y.Start && x.End == y.End && slices.Equal(x.Scopes, y.Scopes)
	})
}
