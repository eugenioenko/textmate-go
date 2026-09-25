package main

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/eugenioenko/textmate-go/oniguruma"
)

func TestScanReportsLoadsAndRegexDiagnostics(t *testing.T) {
	dir := newGrammarCheckout(t)
	writeGrammar(t, dir, "a.json", `{
  "scopeName": "source.a",
  "patterns": [{"match": "("}]
}`)
	writeGrammar(t, dir, "b.json", `{
  "scopeName": "source.b",
  "patterns": [{
    "begin": "(x)",
    "end": "\\1",
    "patterns": [{"include": "source.a"}]
  }]
}`)
	revision := commitCheckout(t, dir)

	got, err := scan(options{grammarDir: dir, revision: revision, wantCount: 2})
	if err != nil {
		t.Fatal(err)
	}
	if got.files != 2 || got.parsed != 2 || got.loaded != 2 {
		t.Fatalf("grammar counts = files %d, parsed %d, loaded %d", got.files, got.parsed, got.loaded)
	}
	if got.regexFields != 3 {
		t.Fatalf("regex fields = %d, want 3", got.regexFields)
	}
	if len(got.loadFailures) != 0 {
		t.Fatalf("load failures = %+v", got.loadFailures)
	}
	if len(got.diagnostics) != 1 {
		t.Fatalf("diagnostics = %+v, want one", got.diagnostics)
	}
	if got.diagnostics[0].diagnostic.Kind != oniguruma.DiagnosticCompileError {
		t.Fatalf("diagnostic kind = %q", got.diagnostics[0].diagnostic.Kind)
	}
	if got.diagnostics[0].file != "a.json" || got.diagnostics[0].path != "$.patterns[0].match" {
		t.Fatalf("diagnostic location = %+v", got.diagnostics[0])
	}

	var output bytes.Buffer
	if err := writeReport(&output, got); err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{
		"revision `" + revision + "`",
		"| Grammar JSON files | 2 |",
		"| Pattern compile failures | 1 |",
		"### 1. `a.json` — `$.patterns[0].match`",
	} {
		if !strings.Contains(output.String(), want) {
			t.Errorf("report missing %q", want)
		}
	}
}

func TestScanRejectsIncompleteCorpus(t *testing.T) {
	dir := newGrammarCheckout(t)
	writeGrammar(t, dir, "only.json", `{"scopeName":"source.only"}`)
	revision := commitCheckout(t, dir)
	_, err := scan(options{grammarDir: dir, revision: revision, wantCount: 2})
	if err == nil || !strings.Contains(err.Error(), "found 1 grammar JSON files") {
		t.Fatalf("scan error = %v", err)
	}
}

func TestScanRejectsCheckoutAtWrongRevision(t *testing.T) {
	dir := newGrammarCheckout(t)
	writeGrammar(t, dir, "only.json", `{"scopeName":"source.only"}`)
	revision := commitCheckout(t, dir)

	_, err := scan(options{grammarDir: dir, revision: strings.Repeat("0", len(revision)), wantCount: 1})
	if err == nil || !strings.Contains(err.Error(), "checkout revision is "+revision) {
		t.Fatalf("scan error = %v", err)
	}
}

func TestScanRejectsDirtyCheckout(t *testing.T) {
	tests := []struct {
		name  string
		dirty func(*testing.T, string)
		want  string
	}{
		{
			name: "tracked modification",
			dirty: func(t *testing.T, dir string) {
				writeGrammar(t, dir, "only.json", `{"scopeName":"source.changed"}`)
			},
			want: " M packages/tm-grammars/grammars/only.json",
		},
		{
			name: "untracked file outside grammar directory",
			dirty: func(t *testing.T, dir string) {
				if err := os.WriteFile(filepath.Join(checkoutRoot(dir), "scratch.txt"), []byte("dirty"), 0o600); err != nil {
					t.Fatal(err)
				}
			},
			want: "?? scratch.txt",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			dir := newGrammarCheckout(t)
			writeGrammar(t, dir, "only.json", `{"scopeName":"source.only"}`)
			revision := commitCheckout(t, dir)
			test.dirty(t, dir)

			_, err := scan(options{grammarDir: dir, revision: revision, wantCount: 1})
			if err == nil || !strings.Contains(err.Error(), "checkout is not clean") || !strings.Contains(err.Error(), test.want) {
				t.Fatalf("scan error = %v", err)
			}
		})
	}
}

func TestReplaceNumericBackReferences(t *testing.T) {
	got := replaceNumericBackReferences(`^(\1)-(\12)-\\3$`, "x")
	if want := `^(x)-(x)-\x$`; got != want {
		t.Fatalf("replaceNumericBackReferences = %q, want %q", got, want)
	}
}

func writeGrammar(t *testing.T, dir, name, contents string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(dir, name), []byte(contents), 0o600); err != nil {
		t.Fatal(err)
	}
}

func newGrammarCheckout(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	runGit(t, root, "init", "--quiet")
	dir := filepath.Join(root, "packages", "tm-grammars", "grammars")
	if err := os.MkdirAll(dir, 0o700); err != nil {
		t.Fatal(err)
	}
	return dir
}

func commitCheckout(t *testing.T, grammarDir string) string {
	t.Helper()
	root := checkoutRoot(grammarDir)
	runGit(t, root, "add", ".")
	runGit(t, root,
		"-c", "user.name=grammar-report tests",
		"-c", "user.email=grammar-report@example.invalid",
		"commit", "--quiet", "-m", "fixture",
	)
	return runGit(t, root, "rev-parse", "HEAD")
}

func checkoutRoot(grammarDir string) string {
	return filepath.Clean(filepath.Join(grammarDir, "..", "..", ".."))
}

func runGit(t *testing.T, dir string, args ...string) string {
	t.Helper()
	commandArgs := append([]string{"-C", dir}, args...)
	output, err := exec.Command("git", commandArgs...).CombinedOutput()
	if err != nil {
		t.Fatalf("git %s: %s: %v", strings.Join(args, " "), output, err)
	}
	return strings.TrimSpace(string(output))
}
