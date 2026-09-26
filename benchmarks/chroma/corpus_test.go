package chromabench

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/eugenioenko/textmate-go/grammars"
)

// The JavaScript harness cannot consult the Go filename tables, so corpus.json
// records each scope explicitly; keep it in agreement with them.
func TestCorpusManifest(t *testing.T) {
	for _, test := range loadCorpus(t) {
		if got := grammars.ScopeForFilename(test.Filename); got != test.ScopeName {
			t.Errorf("%s: scope %q, grammars resolves %q", test.Name, test.ScopeName, got)
		}
		if lexers.Match(test.Filename) == nil {
			t.Errorf("%s: no Chroma lexer for %s", test.Name, test.Filename)
		}
		if _, err := os.Stat(filepath.Join(repoRoot, filepath.FromSlash(test.Fixture))); err != nil {
			t.Errorf("%s: %v", test.Name, err)
		}
		if test.Lines <= 0 {
			t.Errorf("%s: lines must be positive", test.Name)
		}
	}
}
