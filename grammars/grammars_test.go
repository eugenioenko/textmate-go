package grammars

import (
	"reflect"
	"sort"
	"sync"
	"testing"

	textmate "github.com/eugenioenko/textmate-go"
)

func TestLoad(t *testing.T) {
	grammar, err := Load("source.go")
	if err != nil {
		t.Fatal(err)
	}
	if grammar == nil || grammar.ScopeName != "source.go" {
		t.Fatalf("Load(source.go) = %#v", grammar)
	}
	again, err := Load("source.go")
	if err != nil {
		t.Fatal(err)
	}
	if again != grammar {
		t.Fatal("Load did not return its cached result")
	}
	missing, err := Load("source.not-embedded")
	if err != nil || missing != nil {
		t.Fatalf("Load(missing) = (%#v, %v), want (nil, nil)", missing, err)
	}
}

func TestLoadAllConcurrently(t *testing.T) {
	scopes := Scopes()
	if len(scopes) != 41 {
		t.Fatalf("len(Scopes()) = %d, want 41", len(scopes))
	}
	var wait sync.WaitGroup
	for _, scopeName := range scopes {
		scopeName := scopeName
		wait.Add(2)
		for range 2 {
			go func() {
				defer wait.Done()
				grammar, err := Load(scopeName)
				if err != nil {
					t.Errorf("Load(%q): %v", scopeName, err)
				} else if grammar == nil || grammar.ScopeName != scopeName {
					t.Errorf("Load(%q) returned %#v", scopeName, grammar)
				}
			}()
		}
	}
	wait.Wait()
}

func TestForFilename(t *testing.T) {
	tests := map[string]string{
		"main.go":                  "source.go",
		"program.cs":               "source.cs",
		"script.js":                "source.js",
		"script.py":                "source.python",
		"script.sh":                "source.shell",
		"README.md":                "text.html.markdown",
		"SRC/HELLO.CPP":            "source.cpp",
		"component.test.tsx":       "source.tsx",
		"Dockerfile":               "source.dockerfile",
		"build/Makefile":           "source.makefile",
		"config/.bashrc":           "source.shell",
		"settings.jsonc":           "source.json.comments",
		"project/.env":             "source.dotenv",
		"main.m":                   "source.objc",
		"main.mm":                  "source.objcpp",
		"script.pl":                "source.perl",
		"source.cc":                "source.cpp",
		"source.cxx":               "source.cpp",
		"source.c++":               "source.cpp",
		"header.hh":                "source.cpp",
		"header.hpp":               "source.cpp",
		"header.hxx":               "source.cpp",
		"header.h++":               "source.cpp",
		"unknown.extension-absent": "",
	}
	for filename, want := range tests {
		if got := ScopeForFilename(filename); got != want {
			t.Errorf("ScopeForFilename(%q) = %q, want %q", filename, got, want)
		}
		grammar, err := ForFilename(filename)
		if err != nil {
			t.Errorf("ForFilename(%q): %v", filename, err)
			continue
		}
		if want == "" && grammar != nil {
			t.Errorf("ForFilename(%q) = %#v, want nil", filename, grammar)
		} else if want != "" && (grammar == nil || grammar.ScopeName != want) {
			t.Errorf("ForFilename(%q) = %#v, want %q", filename, grammar, want)
		}
	}
}

func TestRegistryIntegration(t *testing.T) {
	registry := textmate.NewRegistry(textmate.RegistryOptions{LoadGrammar: Load})
	defer registry.Dispose()
	grammar, err := registry.LoadGrammar("source.go")
	if err != nil {
		t.Fatal(err)
	}
	result := grammar.TokenizeLine("package main", textmate.InitialState)
	if len(result.Tokens) == 0 {
		t.Fatal("embedded Go grammar returned no tokens")
	}
}

func TestScopesReturnsCopy(t *testing.T) {
	scopes := Scopes()
	scopes[0] = "changed"
	if Scopes()[0] == "changed" {
		t.Fatal("Scopes returned mutable package storage")
	}
}

func TestGrammarInfoLookups(t *testing.T) {
	wantJavaScript := GrammarInfo{
		ID:          "javascript",
		DisplayName: "JavaScript",
		ScopeName:   "source.js",
		Aliases:     []string{"js", "cjs", "mjs"},
	}
	for name, lookup := range map[string]func() (GrammarInfo, bool){
		"ID":       func() (GrammarInfo, bool) { return InfoForID("  JavaScript ") },
		"alias":    func() (GrammarInfo, bool) { return InfoForAlias(" JS ") },
		"scope":    func() (GrammarInfo, bool) { return InfoForScope("source.js") },
		"filename": func() (GrammarInfo, bool) { return InfoForFilename("src/APP.MJS") },
	} {
		t.Run(name, func(t *testing.T) {
			got, ok := lookup()
			if !ok {
				t.Fatal("lookup did not find JavaScript")
			}
			if !reflect.DeepEqual(got, wantJavaScript) {
				t.Fatalf("lookup = %#v, want %#v", got, wantJavaScript)
			}
		})
	}

	graphql, ok := InfoForID("graphql")
	if !ok {
		t.Fatal("InfoForID(graphql) did not find metadata")
	}
	if want := []string{"graphql", "graphqls", "gql", "graphcool"}; !reflect.DeepEqual(graphql.FileTypes, want) {
		t.Fatalf("GraphQL file types = %v, want %v", graphql.FileTypes, want)
	}
}

func TestGrammarInfoIDAndAliasLookupsAreSeparate(t *testing.T) {
	if _, ok := InfoForID("js"); ok {
		t.Fatal("InfoForID searched aliases")
	}
	if _, ok := InfoForAlias("javascript"); ok {
		t.Fatal("InfoForAlias searched canonical IDs")
	}
	if info, ok := InfoForID("javascript"); !ok || info.ID != "javascript" {
		t.Fatalf("InfoForID(javascript) = (%#v, %v)", info, ok)
	}
	if info, ok := InfoForAlias("js"); !ok || info.ID != "javascript" {
		t.Fatalf("InfoForAlias(js) = (%#v, %v)", info, ok)
	}
}

func TestGrammarInfoMissingLookups(t *testing.T) {
	for name, lookup := range map[string]func() (GrammarInfo, bool){
		"ID":       func() (GrammarInfo, bool) { return InfoForID("not-embedded") },
		"alias":    func() (GrammarInfo, bool) { return InfoForAlias("not-embedded") },
		"scope":    func() (GrammarInfo, bool) { return InfoForScope("source.not-embedded") },
		"filename": func() (GrammarInfo, bool) { return InfoForFilename("file.not-embedded") },
	} {
		t.Run(name, func(t *testing.T) {
			info, ok := lookup()
			if ok || !reflect.DeepEqual(info, GrammarInfo{}) {
				t.Fatalf("missing lookup = (%#v, %v), want zero, false", info, ok)
			}
		})
	}
}

func TestGrammarInfosAreCompleteSortedAndIndependent(t *testing.T) {
	infos := Infos()
	if len(infos) != len(Scopes()) {
		t.Fatalf("len(Infos()) = %d, want %d", len(infos), len(Scopes()))
	}
	ids := make([]string, len(infos))
	seenIDs := make(map[string]bool, len(infos))
	seenScopes := make(map[string]bool, len(infos))
	for index, info := range infos {
		ids[index] = info.ID
		if info.ID == "" || info.DisplayName == "" || info.ScopeName == "" {
			t.Errorf("incomplete grammar info: %#v", info)
		}
		if seenIDs[info.ID] {
			t.Errorf("duplicate grammar ID %q", info.ID)
		}
		seenIDs[info.ID] = true
		if seenScopes[info.ScopeName] {
			t.Errorf("duplicate grammar scope %q", info.ScopeName)
		}
		seenScopes[info.ScopeName] = true
		if byID, ok := InfoForID(info.ID); !ok || !reflect.DeepEqual(byID, info) {
			t.Errorf("InfoForID(%q) = (%#v, %v), want %#v", info.ID, byID, ok, info)
		}
		for _, alias := range info.Aliases {
			if byAlias, ok := InfoForAlias(alias); !ok || byAlias.ID != info.ID {
				t.Errorf("InfoForAlias(%q) = (%#v, %v), want ID %q", alias, byAlias, ok, info.ID)
			}
		}
	}
	if !sort.StringsAreSorted(ids) {
		t.Fatalf("Infos IDs are not sorted: %v", ids)
	}

	javascript, ok := InfoForID("javascript")
	if !ok {
		t.Fatal("JavaScript metadata missing")
	}
	javascript.Aliases[0] = "changed"
	infos[0].ID = "changed"
	if again, _ := InfoForID("javascript"); again.Aliases[0] != "js" {
		t.Fatal("InfoForID returned mutable package alias storage")
	}
	if Infos()[0].ID == "changed" {
		t.Fatal("Infos returned mutable package metadata storage")
	}
}

func TestEveryFilenameMappingHasGrammarInfo(t *testing.T) {
	for filename, scopeName := range generatedFilenames {
		info, ok := InfoForFilename(filename)
		if !ok || info.ScopeName != scopeName {
			t.Errorf("InfoForFilename(%q) = (%#v, %v), want scope %q", filename, info, ok, scopeName)
		}
	}
	for extension, scopeName := range generatedExtensions {
		info, ok := InfoForFilename("file." + extension)
		if !ok || info.ScopeName != scopeName {
			t.Errorf("InfoForFilename(file.%s) = (%#v, %v), want scope %q", extension, info, ok, scopeName)
		}
	}
}
