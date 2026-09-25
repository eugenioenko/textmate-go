package main

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func TestProtocolLifecycle(t *testing.T) {
	input := strings.Join([]string{
		`{"id":1,"op":"newRegistry","grammars":[{"scopeName":"source.test","patterns":[{"match":"😀","name":"constant.emoji"}]}]}`,
		`{"id":2,"op":"loadGrammar","registry":"r1","scopeName":"source.test"}`,
		`{"id":3,"op":"tokenizeLine","grammar":"g2","line":"a😀b","state":null}`,
		`{"id":4,"op":"tokenizeLine","grammar":"g2","line":"next","state":"s3"}`,
		`{"id":5,"op":"dispose","registry":"r1"}`,
		`{"id":6,"op":"tokenizeLine","grammar":"g2","line":"gone","state":"s3"}`,
	}, "\n")

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	if err := run(strings.NewReader(input), &stdout, &stderr); err != nil {
		t.Fatalf("run: %v", err)
	}

	var got []response
	decoder := json.NewDecoder(&stdout)
	for decoder.More() {
		var resp response
		if err := decoder.Decode(&resp); err != nil {
			t.Fatalf("decode response: %v", err)
		}
		got = append(got, resp)
	}
	if len(got) != 6 {
		t.Fatalf("got %d responses, want 6", len(got))
	}
	if got[0].Registry != "r1" || got[1].Grammar != "g2" {
		t.Fatalf("unexpected handles: registry=%q grammar=%q", got[0].Registry, got[1].Grammar)
	}
	if len(got[2].Tokens) != 3 {
		t.Fatalf("got %d tokens, want 3: %#v", len(got[2].Tokens), got[2].Tokens)
	}
	token := got[2].Tokens[1]
	if token.Start != 1 || token.End != 3 {
		t.Fatalf("emoji token range = [%d,%d), want [1,3) UTF-16 units", token.Start, token.End)
	}
	wantScopes := []string{"source.test", "constant.emoji"}
	if len(token.Scopes) != len(wantScopes) || token.Scopes[0] != wantScopes[0] || token.Scopes[1] != wantScopes[1] {
		t.Fatalf("token scopes = %v, want %v", token.Scopes, wantScopes)
	}
	if got[2].State != "s3" || got[3].State != "s4" {
		t.Fatalf("immutable state handles = (%q, %q), want (s3, s4)", got[2].State, got[3].State)
	}
	if !got[4].Disposed {
		t.Fatal("dispose response did not report success")
	}
	if got[5].Code != "unknown_grammar" {
		t.Fatalf("post-dispose code = %q, want unknown_grammar", got[5].Code)
	}
	if stderr.Len() == 0 {
		t.Fatal("expected protocol error to be logged on stderr")
	}
}

func TestMalformedRequestDoesNotCorruptStdout(t *testing.T) {
	input := "not json\n" +
		`{"id":"after","op":"newRegistry","grammars":[]}` + "\n"

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	if err := run(strings.NewReader(input), &stdout, &stderr); err != nil {
		t.Fatalf("run: %v", err)
	}

	lines := strings.Split(strings.TrimSpace(stdout.String()), "\n")
	if len(lines) != 2 {
		t.Fatalf("stdout contains %d lines, want 2: %q", len(lines), stdout.String())
	}
	for i, line := range lines {
		if !json.Valid([]byte(line)) {
			t.Fatalf("stdout line %d is not JSON: %q", i+1, line)
		}
	}
	if !strings.Contains(lines[0], `"code":"invalid_json"`) {
		t.Fatalf("malformed-request response = %s", lines[0])
	}
	if !strings.Contains(lines[1], `"registry":"r1"`) {
		t.Fatalf("following response = %s", lines[1])
	}
}

func TestRegexDiagnosticsAreStructuredLoggedAndDeduplicated(t *testing.T) {
	input := strings.Join([]string{
		`{"id":1,"op":"newRegistry","grammars":[{"scopeName":"source.test","patterns":[{"match":"(","name":"invalid"}]}]}`,
		`{"id":2,"op":"loadGrammar","registry":"r1","scopeName":"source.test"}`,
		`{"id":3,"op":"tokenizeLine","grammar":"g2","line":"first","state":null}`,
		`{"id":4,"op":"tokenizeLine","grammar":"g2","line":"second","state":"s3"}`,
	}, "\n")

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	if err := run(strings.NewReader(input), &stdout, &stderr); err != nil {
		t.Fatalf("run: %v", err)
	}

	var got []response
	decoder := json.NewDecoder(&stdout)
	for decoder.More() {
		var resp response
		if err := decoder.Decode(&resp); err != nil {
			t.Fatalf("decode response: %v", err)
		}
		got = append(got, resp)
	}
	if len(got) != 4 {
		t.Fatalf("got %d responses, want 4", len(got))
	}
	if len(got[2].Diagnostics) != 1 || got[2].Diagnostics[0].Kind != "compile_error" || got[2].Diagnostics[0].Pattern != "(" {
		t.Fatalf("first diagnostics = %+v", got[2].Diagnostics)
	}
	if len(got[3].Diagnostics) != 0 {
		t.Fatalf("duplicate diagnostics = %+v, want none", got[3].Diagnostics)
	}
	if count := strings.Count(stderr.String(), "textmate-cli: compile_error:"); count != 1 {
		t.Fatalf("compile diagnostics logged %d times, stderr=%q", count, stderr.String())
	}
}

func TestStateCannotCrossGrammars(t *testing.T) {
	s := newServer()
	registry := s.handle([]byte(`{"id":1,"op":"newRegistry","grammars":[{"scopeName":"source.a"},{"scopeName":"source.b"}]}`))
	grammarA := s.handle([]byte(`{"id":2,"op":"loadGrammar","registry":"r1","scopeName":"source.a"}`))
	grammarB := s.handle([]byte(`{"id":3,"op":"loadGrammar","registry":"r1","scopeName":"source.b"}`))
	first := s.handle([]byte(`{"id":4,"op":"tokenizeLine","grammar":"g2","line":"x","state":null}`))
	crossed := s.handle([]byte(`{"id":5,"op":"tokenizeLine","grammar":"g3","line":"y","state":"s4"}`))

	if registry.Error != "" || grammarA.Error != "" || grammarB.Error != "" || first.Error != "" {
		t.Fatalf("setup failed: registry=%q grammarA=%q grammarB=%q first=%q", registry.Error, grammarA.Error, grammarB.Error, first.Error)
	}
	if crossed.Code != "state_mismatch" {
		t.Fatalf("cross-grammar state code = %q, want state_mismatch", crossed.Code)
	}
}

func TestStateHandlesRemainReusable(t *testing.T) {
	s := newServer()
	registry := s.handle([]byte(`{"id":1,"op":"newRegistry","grammars":[{"scopeName":"source.test","patterns":[]}]}`))
	grammar := s.handle([]byte(`{"id":2,"op":"loadGrammar","registry":"r1","scopeName":"source.test"}`))
	first := s.handle([]byte(`{"id":3,"op":"tokenizeLine","grammar":"g2","line":"first","state":null}`))
	left := s.handle([]byte(`{"id":4,"op":"tokenizeLine","grammar":"g2","line":"left","state":"s3"}`))
	right := s.handle([]byte(`{"id":5,"op":"tokenizeLine","grammar":"g2","line":"right","state":"s3"}`))

	if registry.Error != "" || grammar.Error != "" || first.Error != "" || left.Error != "" || right.Error != "" {
		t.Fatalf(
			"branching setup failed: registry=%q grammar=%q first=%q left=%q right=%q",
			registry.Error,
			grammar.Error,
			first.Error,
			left.Error,
			right.Error,
		)
	}
	if first.State != "s3" || left.State != "s4" || right.State != "s5" {
		t.Fatalf("state handles = (%q, %q, %q), want (s3, s4, s5)", first.State, left.State, right.State)
	}
	if _, ok := s.states[first.State]; !ok {
		t.Fatalf("original state %q was discarded after reuse", first.State)
	}
}

func TestRegistryInjectionMapIsApplied(t *testing.T) {
	s := newServer()
	registry := s.handle([]byte(`{
		"id":1,
		"op":"newRegistry",
		"grammars":[
			{"scopeName":"source.root","patterns":[]},
			{
				"scopeName":"source.injected",
				"injectionSelector":"L:source.root",
				"patterns":[{"match":"x","name":"injected.match"}]
			}
		],
		"injections":{"source.root":["source.injected"]}
	}`))
	grammar := s.handle([]byte(`{"id":2,"op":"loadGrammar","registry":"r1","scopeName":"source.root"}`))
	line := s.handle([]byte(`{"id":3,"op":"tokenizeLine","grammar":"g2","line":"x","state":null}`))

	if registry.Error != "" || grammar.Error != "" || line.Error != "" {
		t.Fatalf("injection setup failed: registry=%q grammar=%q line=%q", registry.Error, grammar.Error, line.Error)
	}
	if len(line.Tokens) != 1 {
		t.Fatalf("tokens = %#v, want one injected token", line.Tokens)
	}
	wantScopes := []string{"source.root", "injected.match"}
	if len(line.Tokens[0].Scopes) != len(wantScopes) || line.Tokens[0].Scopes[0] != wantScopes[0] || line.Tokens[0].Scopes[1] != wantScopes[1] {
		t.Fatalf("token scopes = %v, want %v", line.Tokens[0].Scopes, wantScopes)
	}
}
