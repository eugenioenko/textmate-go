package main

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"testing"
)

func TestRejectedLicenseDoesNotMutateOutput(t *testing.T) {
	for _, test := range []struct {
		name    string
		license string
		want    string
	}{
		{name: "missing", want: "metadata has no license"},
		{name: "unreviewed", license: "GPL-3.0", want: `unreviewed license "GPL-3.0"`},
	} {
		t.Run(test.name, func(t *testing.T) {
			checkout, revision := makeSourceCheckout(t, []fixtureGrammar{
				{id: "good", scope: "source.good", license: "MIT"},
				{id: "sample", scope: "source.sample", license: test.license},
			})
			selection := writeSelection(t, "good\nsample\n")
			out := t.TempDir()
			dataDir := filepath.Join(out, "data")
			if err := os.Mkdir(dataDir, 0o755); err != nil {
				t.Fatal(err)
			}
			mustWrite(t, filepath.Join(dataDir, "existing.json.gz"), []byte("existing asset"))
			mustWrite(t, filepath.Join(out, "generated.go"), []byte("existing generated file"))
			before := snapshotTree(t, out)

			err := generate(generatorConfig{
				source: checkout, revision: revision, selection: selection, out: out,
			})
			if err == nil || !strings.Contains(err.Error(), test.want) {
				t.Fatalf("generate error = %v, want substring %q", err, test.want)
			}
			if after := snapshotTree(t, out); !equalSnapshots(before, after) {
				t.Fatalf("rejected generation mutated output\nbefore: %v\nafter:  %v", before, after)
			}
			for _, asset := range []string{"good.json.gz", "sample.json.gz"} {
				if _, err := os.Stat(filepath.Join(dataDir, asset)); !os.IsNotExist(err) {
					t.Fatalf("asset %s exists after rejected generation: %v", asset, err)
				}
			}
		})
	}
}

func TestGenerationIsDeterministic(t *testing.T) {
	checkout, revision := makeSourceCheckout(t, []fixtureGrammar{
		{id: "alpha", scope: "source.alpha", license: "MIT", fileTypes: []string{"a"}},
		{id: "beta", scope: "source.beta", license: "Apache-2.0", fileTypes: []string{"b"}},
	})
	selection := writeSelection(t, "alpha\nbeta\n")
	outA := filepath.Join(t.TempDir(), "generated")
	outB := filepath.Join(t.TempDir(), "generated")
	for _, out := range []string{outA, outB} {
		if err := generate(generatorConfig{
			source: checkout, revision: revision, selection: selection, out: out,
		}); err != nil {
			t.Fatalf("generate %s: %v", out, err)
		}
	}
	if a, b := snapshotTree(t, outA), snapshotTree(t, outB); !equalSnapshots(a, b) {
		t.Fatalf("generated trees differ\nA: %v\nB: %v", a, b)
	}
}

func TestScopeMetadataMustMatchBeforeOutputMutation(t *testing.T) {
	checkout, revision := makeSourceCheckout(t, []fixtureGrammar{
		{id: "good", scope: "source.good", license: "MIT"},
		{id: "sample", scope: "source.sample", metadataScope: "source.other", license: "MIT"},
	})
	out := filepath.Join(t.TempDir(), "not-created")
	err := generate(generatorConfig{
		source:    checkout,
		revision:  revision,
		selection: writeSelection(t, "good\nsample\n"),
		out:       out,
	})
	if err == nil || !strings.Contains(err.Error(), "does not match metadata") {
		t.Fatalf("generate error = %v, want scope mismatch", err)
	}
	if _, statErr := os.Stat(out); !os.IsNotExist(statErr) {
		t.Fatalf("invalid generation created output directory: %v", statErr)
	}
}

type fixtureGrammar struct {
	id            string
	scope         string
	metadataScope string
	license       string
	fileTypes     []string
}

func makeSourceCheckout(t *testing.T, grammars []fixtureGrammar) (string, string) {
	t.Helper()
	root := t.TempDir()
	packageDir := filepath.Join(root, "packages", "tm-grammars")
	grammarDir := filepath.Join(packageDir, "grammars")
	if err := os.MkdirAll(grammarDir, 0o755); err != nil {
		t.Fatal(err)
	}
	var metadata strings.Builder
	metadata.WriteString("export const grammars = [\n")
	for _, grammar := range grammars {
		metadataScope := grammar.metadataScope
		if metadataScope == "" {
			metadataScope = grammar.scope
		}
		metadata.WriteString("  {\n")
		if grammar.license != "" {
			fmt.Fprintf(&metadata, "    license: '%s',\n", grammar.license)
		}
		fmt.Fprintf(&metadata, "    name: '%s',\n", grammar.id)
		fmt.Fprintf(&metadata, "    scopeName: '%s',\n", metadataScope)
		fmt.Fprintf(&metadata, "    sha: '%040d',\n", len(grammar.id))
		fmt.Fprintf(&metadata, "    source: 'https://example.invalid/%s',\n", grammar.id)
		metadata.WriteString("  },\n")

		fileTypes := ""
		for index, fileType := range grammar.fileTypes {
			if index != 0 {
				fileTypes += ","
			}
			fileTypes += fmt.Sprintf("%q", fileType)
		}
		contents := fmt.Sprintf(
			"{\"scopeName\":%q,\"fileTypes\":[%s],\"patterns\":[]}",
			grammar.scope, fileTypes,
		)
		mustWrite(t, filepath.Join(grammarDir, grammar.id+".json"), []byte(contents))
	}
	metadata.WriteString("]\n")
	mustWrite(t, filepath.Join(packageDir, "index.js"), []byte(metadata.String()))
	mustWrite(t, filepath.Join(packageDir, "NOTICE"), []byte("fixture notice\n"))
	mustWrite(t, filepath.Join(packageDir, "LICENSE"), []byte("fixture license\n"))

	runGit(t, root, "init", "-q")
	runGit(t, root, "config", "user.email", "generator-test@example.invalid")
	runGit(t, root, "config", "user.name", "Generator Test")
	runGit(t, root, "add", ".")
	runGit(t, root, "commit", "-qm", "fixture")
	revision := strings.TrimSpace(runGit(t, root, "rev-parse", "HEAD"))
	return root, revision
}

func writeSelection(t *testing.T, contents string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "selection.txt")
	mustWrite(t, path, []byte(contents))
	return path
}

func mustWrite(t *testing.T, path string, contents []byte) {
	t.Helper()
	if err := os.WriteFile(path, contents, 0o644); err != nil {
		t.Fatal(err)
	}
}

func runGit(t *testing.T, dir string, arguments ...string) string {
	t.Helper()
	command := exec.Command("git", arguments...)
	command.Dir = dir
	output, err := command.CombinedOutput()
	if err != nil {
		t.Fatalf("git %s: %v\n%s", strings.Join(arguments, " "), err, output)
	}
	return string(output)
}

func snapshotTree(t *testing.T, root string) map[string][]byte {
	t.Helper()
	result := make(map[string][]byte)
	err := filepath.WalkDir(root, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			return nil
		}
		relative, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		result[filepath.ToSlash(relative)], err = os.ReadFile(path)
		return err
	})
	if err != nil {
		t.Fatal(err)
	}
	return result
}

func equalSnapshots(left, right map[string][]byte) bool {
	if len(left) != len(right) {
		return false
	}
	keys := make([]string, 0, len(left))
	for key := range left {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		if !bytes.Equal(left[key], right[key]) {
			return false
		}
	}
	return true
}
