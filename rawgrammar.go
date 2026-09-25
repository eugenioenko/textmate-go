package textmate

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
)

// IncludeString is a TextMate include reference such as "#rule", "$self", or
// "source.js#rule".
type IncludeString = string

// RegExpString is a regular expression in the grammar's Oniguruma dialect.
type RegExpString = string

// RawRepository maps grammar-defined repository names to rules.
type RawRepository map[string]*RawRule

// UnmarshalJSON ignores malformed non-object entries found in a small number
// of converted grammars. JavaScript loads these values but they do not compile
// into meaningful rules.
func (repository *RawRepository) UnmarshalJSON(data []byte) error {
	result, err := decodeRawRuleMap(data, false)
	if err != nil {
		return err
	}
	*repository = RawRepository(result)
	return nil
}

// RawCaptures maps numeric capture group names to their rules.
type RawCaptures map[string]*RawRule

// UnmarshalJSON accepts both the specified object form and the array form
// emitted by a few converted TextMate grammars. Invalid capture entries are
// ignored, matching JavaScript's permissive property access on those values.
func (captures *RawCaptures) UnmarshalJSON(data []byte) error {
	result, err := decodeRawRuleMap(data, true)
	if err != nil {
		return err
	}
	*captures = RawCaptures(result)
	return nil
}

// RawInjections maps scope selector expressions to injected rules.
type RawInjections map[string]*RawRule

// Location identifies a value in a source grammar. Plain JSON grammars do not
// normally contain it, but vscode-textmate can attach it while parsing.
type Location struct {
	Filename string `json:"filename"`
	Line     int    `json:"line"`
	Char     int    `json:"char"`
}

// RawGrammar is the JSON representation of a TextMate grammar.
type RawGrammar struct {
	Repository        RawRepository `json:"repository"`
	ScopeName         string        `json:"scopeName"`
	Patterns          []*RawRule    `json:"patterns"`
	Injections        RawInjections `json:"injections,omitempty"`
	InjectionSelector string        `json:"injectionSelector,omitempty"`
	FileTypes         []string      `json:"fileTypes,omitempty"`
	Name              string        `json:"name,omitempty"`
	FirstLineMatch    RegExpString  `json:"firstLineMatch,omitempty"`
	Location          *Location     `json:"$vscodeTextmateLocation,omitempty"`
}

// RawRule is one rule in a TextMate grammar.
type RawRule struct {
	Include *IncludeString `json:"include,omitempty"`

	Name        *string `json:"name,omitempty"`
	ContentName *string `json:"contentName,omitempty"`

	Match         *RegExpString `json:"match,omitempty"`
	Captures      RawCaptures   `json:"captures,omitempty"`
	Begin         *RegExpString `json:"begin,omitempty"`
	BeginCaptures RawCaptures   `json:"beginCaptures,omitempty"`
	End           *RegExpString `json:"end,omitempty"`
	EndCaptures   RawCaptures   `json:"endCaptures,omitempty"`
	While         *RegExpString `json:"while,omitempty"`
	WhileCaptures RawCaptures   `json:"whileCaptures,omitempty"`
	Patterns      []*RawRule    `json:"patterns,omitempty"`

	Repository          RawRepository `json:"repository,omitempty"`
	ApplyEndPatternLast bool          `json:"applyEndPatternLast,omitempty"`
	Location            *Location     `json:"$vscodeTextmateLocation,omitempty"`
}

// ParseRawGrammar decodes a JSON TextMate grammar.
func ParseRawGrammar(data []byte) (*RawGrammar, error) {
	var grammar RawGrammar
	if err := json.Unmarshal(data, &grammar); err != nil {
		return nil, fmt.Errorf("textmate: parse raw grammar: %w", err)
	}
	if grammar.ScopeName == "" {
		return nil, fmt.Errorf("textmate: parse raw grammar: missing scopeName")
	}
	return &grammar, nil
}

func (r *RawRule) UnmarshalJSON(data []byte) error {
	// Some widely used grammars encode applyEndPatternLast as 0 or 1 even
	// though the TextMate interface declares it as a boolean.
	type rawRule RawRule
	decoded := struct {
		*rawRule
		ApplyEndPatternLast json.RawMessage `json:"applyEndPatternLast"`
	}{rawRule: (*rawRule)(r)}
	if err := json.Unmarshal(data, &decoded); err != nil {
		return err
	}
	if len(decoded.ApplyEndPatternLast) == 0 || string(decoded.ApplyEndPatternLast) == "null" {
		return nil
	}
	if err := json.Unmarshal(decoded.ApplyEndPatternLast, &r.ApplyEndPatternLast); err == nil {
		return nil
	}

	var number json.Number
	if err := json.Unmarshal(decoded.ApplyEndPatternLast, &number); err != nil {
		return fmt.Errorf("applyEndPatternLast must be a boolean or number: %w", err)
	}
	value, err := number.Float64()
	if err != nil {
		return fmt.Errorf("applyEndPatternLast must be a boolean or number: %w", err)
	}
	r.ApplyEndPatternLast = value != 0
	return nil
}

func decodeRawRuleMap(data []byte, allowArray bool) (map[string]*RawRule, error) {
	data = bytes.TrimSpace(data)
	if len(data) == 0 || bytes.Equal(data, []byte("null")) {
		return nil, nil
	}

	values := make(map[string]json.RawMessage)
	switch data[0] {
	case '{':
		if err := json.Unmarshal(data, &values); err != nil {
			return nil, err
		}
	case '[':
		if !allowArray {
			return nil, fmt.Errorf("repository must be an object")
		}
		var array []json.RawMessage
		if err := json.Unmarshal(data, &array); err != nil {
			return nil, err
		}
		for index, value := range array {
			values[strconv.Itoa(index)] = value
		}
	default:
		if allowArray {
			return nil, fmt.Errorf("captures must be an object or array")
		}
		return nil, fmt.Errorf("repository must be an object")
	}

	result := make(map[string]*RawRule, len(values))
	for key, value := range values {
		value = bytes.TrimSpace(value)
		if len(value) == 0 || value[0] != '{' {
			continue
		}
		var rule RawRule
		if err := json.Unmarshal(value, &rule); err != nil {
			return nil, fmt.Errorf("rule %q: %w", key, err)
		}
		result[key] = &rule
	}
	return result, nil
}
