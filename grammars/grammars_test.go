package grammars

import (
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
