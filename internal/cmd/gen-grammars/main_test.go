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

func TestDisplayNameIsRequiredBeforeOutputMutation(t *testing.T) {
	checkout, revision := makeSourceCheckout(t, []fixtureGrammar{
		{id: "sample", scope: "source.sample", license: "MIT", omitDisplayName: true},
	})
	out := filepath.Join(t.TempDir(), "not-created")
	err := generate(generatorConfig{
		source:    checkout,
		revision:  revision,
		selection: writeSelection(t, "sample\n"),
		out:       out,
	})
	if err == nil || !strings.Contains(err.Error(), "metadata has no displayName") {
		t.Fatalf("generate error = %v, want missing displayName", err)
	}
	if _, statErr := os.Stat(out); !os.IsNotExist(statErr) {
		t.Fatalf("invalid generation created output directory: %v", statErr)
	}
}

func TestGeneratedMetadataIsDeterministicAndKeepsLookupsSeparate(t *testing.T) {
	grammars := []grammarMetadata{
		{
			ID: "zeta", DisplayName: "Zeta Script", ScopeName: "source.zeta",
			Aliases: []string{"z", "shared"}, FileTypes: []string{"zeta"},
			Asset: "data/zeta.json.gz",
		},
		{
			ID: "alpha", DisplayName: "Alpha Script", ScopeName: "source.alpha",
			Aliases: []string{"a", "shared"}, FileTypes: []string{"alpha", "alpha.test"},
			Asset: "data/alpha.json.gz",
		},
		{
			ID: "shared", DisplayName: "Shared", ScopeName: "source.shared",
			Asset: "data/shared.json.gz",
		},
	}
	forward := filepath.Join(t.TempDir(), "forward.go")
	reverse := filepath.Join(t.TempDir(), "reverse.go")
	if err := writeGeneratedGo(forward, "revision", grammars); err != nil {
		t.Fatal(err)
	}
	for left, right := 0, len(grammars)-1; left < right; left, right = left+1, right-1 {
		grammars[left], grammars[right] = grammars[right], grammars[left]
	}
	if err := writeGeneratedGo(reverse, "revision", grammars); err != nil {
		t.Fatal(err)
	}
	forwardData, err := os.ReadFile(forward)
	if err != nil {
		t.Fatal(err)
	}
	reverseData, err := os.ReadFile(reverse)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(forwardData, reverseData) {
		t.Fatal("generated metadata depends on input order")
	}
	contents := string(forwardData)
	for _, want := range []string{
		`ID:          "alpha"`,
		`DisplayName: "Alpha Script"`,
		`ScopeName:   "source.alpha"`,
		`Aliases:     []string{"a", "shared"}`,
		`FileTypes:   []string{"alpha", "alpha.test"}`,
		`"shared": "alpha"`, // Alias collision resolves to lexicographically first ID.
		`"shared": 1`,       // The canonical-ID lookup remains a separate map.
	} {
		if !strings.Contains(contents, want) {
			t.Errorf("generated metadata does not contain %q\n%s", want, contents)
		}
	}
}

type fixtureGrammar struct {
	id              string
	displayName     string
	omitDisplayName bool
	scope           string
	metadataScope   string
	license         string
	aliases         []string
	fileTypes       []string
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
		if !grammar.omitDisplayName {
			displayName := grammar.displayName
			if displayName == "" {
				displayName = "Display " + grammar.id
			}
			fmt.Fprintf(&metadata, "    displayName: '%s',\n", displayName)
		}
		if len(grammar.aliases) != 0 {
			metadata.WriteString("    aliases: [\n")
			for _, alias := range grammar.aliases {
				fmt.Fprintf(&metadata, "      '%s',\n", alias)
			}
			metadata.WriteString("    ],\n")
		}
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

func TestLicenseReviewFillsMissingLicense(t *testing.T) {
	checkout, revision := makeSourceCheckout(t, []fixtureGrammar{
		{id: "missing", scope: "source.missing"},
		{id: "unasserted", scope: "source.unasserted", license: "NOASSERTION"},
	})
	selection := writeSelection(t, "missing\nunasserted\n")
	reviews := filepath.Join(t.TempDir(), "reviews.json")
	mustWrite(t, reviews, []byte(`[
		{"id": "missing", "license": "TextMate-Bundle", "evidence": "https://example.invalid/missing", "text": "Permission granted."},
		{"id": "unasserted", "license": "MIT", "evidence": "https://example.invalid/unasserted", "text": "MIT text."}
	]`))
	out := t.TempDir()
	if err := generate(generatorConfig{
		source: checkout, revision: revision, selection: selection, licenseReviews: reviews, out: out,
	}); err != nil {
		t.Fatal(err)
	}
	notice, err := os.ReadFile(filepath.Join(out, "NOTICE"))
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"fixture notice", "## missing (TextMate-Bundle)", "Permission granted.", "## unasserted (MIT)"} {
		if !strings.Contains(string(notice), want) {
			t.Errorf("NOTICE lacks %q:\n%s", want, notice)
		}
	}
	manifest, err := os.ReadFile(filepath.Join(out, "MANIFEST.md"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(manifest), "| TextMate-Bundle |") {
		t.Errorf("MANIFEST.md lacks reviewed license:\n%s", manifest)
	}
}

func TestLicenseReviewRejectsConflictsAndStaleEntries(t *testing.T) {
	for _, test := range []struct {
		name    string
		reviews string
		want    string
	}{
		{
			name:    "metadata already licensed",
			reviews: `[{"id": "good", "license": "MIT", "evidence": "e", "text": "t"}]`,
			want:    `has a license review but metadata states "MIT"`,
		},
		{
			name:    "unselected grammar",
			reviews: `[{"id": "other", "license": "MIT", "evidence": "e", "text": "t"}]`,
			want:    `license review for "other" matches no selected grammar`,
		},
		{
			name:    "incomplete",
			reviews: `[{"id": "good", "license": "MIT"}]`,
			want:    "must set id, license, evidence, and text",
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			checkout, revision := makeSourceCheckout(t, []fixtureGrammar{
				{id: "good", scope: "source.good", license: "MIT"},
			})
			reviews := filepath.Join(t.TempDir(), "reviews.json")
			mustWrite(t, reviews, []byte(test.reviews))
			err := generate(generatorConfig{
				source: checkout, revision: revision, selection: writeSelection(t, "good\n"),
				licenseReviews: reviews, out: t.TempDir(),
			})
			if err == nil || !strings.Contains(err.Error(), test.want) {
				t.Fatalf("generate error = %v, want substring %q", err, test.want)
			}
		})
	}
}
