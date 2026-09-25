package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"unicode/utf16"

	textmate "github.com/eugenioenko/textmate-go"
)

type request struct {
	ID                json.RawMessage     `json:"id"`
	Op                string              `json:"op"`
	Grammars          []json.RawMessage   `json:"grammars"`
	Injections        map[string][]string `json:"injections"`
	Registry          string              `json:"registry"`
	ScopeName         string              `json:"scopeName"`
	EmbeddedLanguages map[string]int      `json:"embeddedLanguages"`
	Grammar           string              `json:"grammar"`
	Line              *string             `json:"line"`
	State             *string             `json:"state"`
}

type response struct {
	ID          json.RawMessage   `json:"id"`
	Registry    string            `json:"registry,omitempty"`
	Grammar     string            `json:"grammar,omitempty"`
	Tokens      []token           `json:"tokens,omitempty"`
	State       string            `json:"state,omitempty"`
	Stopped     bool              `json:"stopped,omitempty"`
	Disposed    bool              `json:"disposed,omitempty"`
	Diagnostics []regexDiagnostic `json:"diagnostics,omitempty"`
	Error       string            `json:"error,omitempty"`
	Code        string            `json:"code,omitempty"`
}

type token struct {
	Start  int      `json:"start"`
	End    int      `json:"end"`
	Scopes []string `json:"scopes"`
}

type registryEntry struct {
	registry       *textmate.Registry
	grammars       map[string]*textmate.RawGrammar
	grammarHandles map[string]struct{}
	stateHandles   map[string]struct{}
}

type grammarEntry struct {
	registry            string
	grammar             *textmate.Grammar
	reportedDiagnostics map[string]struct{}
}

type regexDiagnostic struct {
	Kind       string `json:"kind"`
	Pattern    string `json:"pattern"`
	Translated string `json:"translated,omitempty"`
	Message    string `json:"message"`
}

type stateEntry struct {
	registry string
	grammar  string
	state    *textmate.StateStack
}

type protocolError struct {
	code    string
	message string
}

func (e *protocolError) Error() string { return e.message }

type server struct {
	registries map[string]*registryEntry
	grammars   map[string]grammarEntry
	states     map[string]stateEntry
	next       uint64
}

func newServer() *server {
	return &server{
		registries: make(map[string]*registryEntry),
		grammars:   make(map[string]grammarEntry),
		states:     make(map[string]stateEntry),
	}
}

func (s *server) handle(data []byte) response {
	var req request
	if err := json.Unmarshal(data, &req); err != nil {
		return errorResponse(nil, "invalid_json", fmt.Sprintf("invalid JSON request: %v", err))
	}
	if len(req.ID) == 0 {
		req.ID = json.RawMessage("null")
	}

	var resp response
	var err error
	switch req.Op {
	case "newRegistry":
		resp, err = s.newRegistry(req)
	case "loadGrammar":
		resp, err = s.loadGrammar(req)
	case "tokenizeLine":
		resp, err = s.tokenizeLine(req)
	case "dispose":
		resp, err = s.dispose(req)
	default:
		err = &protocolError{code: "unknown_operation", message: fmt.Sprintf("unknown operation %q", req.Op)}
	}
	if err != nil {
		var protocolErr *protocolError
		if errors.As(err, &protocolErr) {
			return errorResponse(req.ID, protocolErr.code, protocolErr.message)
		}
		return errorResponse(req.ID, "internal_error", err.Error())
	}
	resp.ID = req.ID
	return resp
}

func (s *server) newRegistry(req request) (response, error) {
	grammars := make(map[string]*textmate.RawGrammar, len(req.Grammars))
	for i, data := range req.Grammars {
		grammar, err := textmate.ParseRawGrammar(data)
		if err != nil {
			return response{}, &protocolError{
				code:    "invalid_grammar",
				message: fmt.Sprintf("grammar %d is invalid: %v", i, err),
			}
		}
		if _, exists := grammars[grammar.ScopeName]; exists {
			return response{}, &protocolError{
				code:    "duplicate_scope",
				message: fmt.Sprintf("duplicate grammar scopeName %q", grammar.ScopeName),
			}
		}
		grammars[grammar.ScopeName] = grammar
	}
	injections := make(map[string][]string, len(req.Injections))
	for scopeName, values := range req.Injections {
		injections[scopeName] = append([]string(nil), values...)
	}
	registry := textmate.NewRegistry(textmate.RegistryOptions{
		LoadGrammar: func(scopeName string) (*textmate.RawGrammar, error) {
			return grammars[scopeName], nil
		},
		GetInjections: func(scopeName string) []string {
			return append([]string(nil), injections[scopeName]...)
		},
	})

	handle := s.newHandle("r")
	s.registries[handle] = &registryEntry{
		registry:       registry,
		grammars:       grammars,
		grammarHandles: make(map[string]struct{}),
		stateHandles:   make(map[string]struct{}),
	}
	return response{Registry: handle}, nil
}

func (s *server) loadGrammar(req request) (response, error) {
	registry, ok := s.registries[req.Registry]
	if !ok {
		return response{}, unknownHandle("registry", req.Registry)
	}
	if req.ScopeName == "" {
		return response{}, &protocolError{code: "invalid_scope", message: "scopeName must not be empty"}
	}
	if _, ok := registry.grammars[req.ScopeName]; !ok {
		return response{}, &protocolError{
			code:    "grammar_not_found",
			message: fmt.Sprintf("grammar %q is not registered", req.ScopeName),
		}
	}
	grammar, err := registry.registry.LoadGrammar(req.ScopeName)
	if err != nil {
		return response{}, fmt.Errorf("load grammar %q: %w", req.ScopeName, err)
	}

	handle := s.newHandle("g")
	s.grammars[handle] = grammarEntry{
		registry:            req.Registry,
		grammar:             grammar,
		reportedDiagnostics: make(map[string]struct{}),
	}
	registry.grammarHandles[handle] = struct{}{}
	return response{Grammar: handle}, nil
}

func (s *server) tokenizeLine(req request) (response, error) {
	grammar, ok := s.grammars[req.Grammar]
	if !ok {
		return response{}, unknownHandle("grammar", req.Grammar)
	}
	if req.Line == nil {
		return response{}, &protocolError{code: "invalid_line", message: "line is required"}
	}

	var previousState *textmate.StateStack
	if req.State != nil {
		state, ok := s.states[*req.State]
		if !ok {
			return response{}, unknownHandle("state", *req.State)
		}
		if state.grammar != req.Grammar {
			return response{}, &protocolError{
				code:    "state_mismatch",
				message: fmt.Sprintf("state %q does not belong to grammar %q", *req.State, req.Grammar),
			}
		}
		previousState = state.state
	}

	result := grammar.grammar.TokenizeLine(*req.Line, previousState)
	diagnostics := newDiagnostics(grammar)
	stateHandle := s.newHandle("s")
	s.states[stateHandle] = stateEntry{
		registry: grammar.registry,
		grammar:  req.Grammar,
		state:    result.RuleStack,
	}
	s.registries[grammar.registry].stateHandles[stateHandle] = struct{}{}
	tokens := make([]token, len(result.Tokens))
	offsets := utf16Offsets(*req.Line)
	for i, value := range result.Tokens {
		tokens[i] = token{
			Start:  utf16Offset(offsets, value.Start),
			End:    utf16Offset(offsets, value.End),
			Scopes: append([]string(nil), value.Scopes...),
		}
	}

	return response{
		Tokens:      tokens,
		State:       stateHandle,
		Stopped:     result.Stopped,
		Diagnostics: diagnostics,
	}, nil
}

func newDiagnostics(grammar grammarEntry) []regexDiagnostic {
	var result []regexDiagnostic
	for _, diagnostic := range grammar.grammar.Diagnostics() {
		key := fmt.Sprintf(
			"%s\x00%s\x00%s\x00%s",
			diagnostic.Kind,
			diagnostic.Pattern,
			diagnostic.Translated,
			diagnostic.Message,
		)
		if _, exists := grammar.reportedDiagnostics[key]; exists {
			continue
		}
		grammar.reportedDiagnostics[key] = struct{}{}
		result = append(result, regexDiagnostic{
			Kind:       string(diagnostic.Kind),
			Pattern:    diagnostic.Pattern,
			Translated: diagnostic.Translated,
			Message:    diagnostic.Message,
		})
	}
	return result
}

// utf16Offsets maps each rune boundary in line to its UTF-16 code-unit offset.
// Grammar token offsets are rune-based; the JSON protocol mirrors JavaScript.
func utf16Offsets(line string) []int {
	runes := []rune(line)
	offsets := make([]int, len(runes)+1)
	for i, value := range runes {
		offsets[i+1] = offsets[i] + len(utf16.Encode([]rune{value}))
	}
	return offsets
}

func utf16Offset(offsets []int, runeOffset int) int {
	if runeOffset <= 0 {
		return 0
	}
	if runeOffset >= len(offsets) {
		return offsets[len(offsets)-1]
	}
	return offsets[runeOffset]
}

func (s *server) dispose(req request) (response, error) {
	registry, ok := s.registries[req.Registry]
	if !ok {
		return response{}, unknownHandle("registry", req.Registry)
	}
	registry.registry.Dispose()
	for handle := range registry.grammarHandles {
		delete(s.grammars, handle)
	}
	for handle := range registry.stateHandles {
		delete(s.states, handle)
	}
	delete(s.registries, req.Registry)
	return response{Disposed: true}, nil
}

func (s *server) newHandle(prefix string) string {
	s.next++
	return fmt.Sprintf("%s%d", prefix, s.next)
}

func unknownHandle(kind, handle string) error {
	return &protocolError{
		code:    "unknown_" + kind,
		message: fmt.Sprintf("unknown %s handle %q", kind, handle),
	}
}

func errorResponse(id json.RawMessage, code, message string) response {
	if len(id) == 0 {
		id = json.RawMessage("null")
	}
	return response{ID: id, Error: message, Code: code}
}

func run(in io.Reader, out io.Writer, errOut io.Writer) error {
	server := newServer()
	reader := bufio.NewReader(in)
	encoder := json.NewEncoder(out)
	for {
		line, readErr := reader.ReadBytes('\n')
		line = bytes.TrimSpace(line)
		if len(line) != 0 {
			resp := server.handle(line)
			if resp.Error != "" {
				if _, err := fmt.Fprintf(errOut, "textmate-cli: %s: %s\n", resp.Code, resp.Error); err != nil {
					return fmt.Errorf("write diagnostic: %w", err)
				}
			}
			if len(resp.Diagnostics) != 0 {
				diagnostic := resp.Diagnostics[0]
				if _, err := fmt.Fprintf(
					errOut,
					"textmate-cli: %s: %s: %s (%d new diagnostic(s))\n",
					diagnostic.Kind,
					abbreviateDiagnostic(diagnostic.Pattern),
					abbreviateDiagnostic(diagnostic.Message),
					len(resp.Diagnostics),
				); err != nil {
					return fmt.Errorf("write regex diagnostic: %w", err)
				}
			}
			if err := encoder.Encode(resp); err != nil {
				return fmt.Errorf("write response: %w", err)
			}
		}
		if readErr != nil {
			if errors.Is(readErr, io.EOF) {
				return nil
			}
			return fmt.Errorf("read request: %w", readErr)
		}
	}
}

func abbreviateDiagnostic(value string) string {
	const limit = 160
	value = strings.Join(strings.Fields(value), " ")
	runes := []rune(value)
	if len(runes) <= limit {
		return value
	}
	return string(runes[:limit]) + "…"
}

func main() {
	if err := run(os.Stdin, os.Stdout, os.Stderr); err != nil {
		fmt.Fprintln(os.Stderr, "textmate-cli:", err)
		os.Exit(1)
	}
}
