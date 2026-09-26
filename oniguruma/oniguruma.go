// Package oniguruma adapts the Oniguruma subset used by TextMate grammars to
// regexp2. Offsets are rune indices, deliberately differing from the UTF-16
// offsets exposed by vscode-textmate's JavaScript API.
package oniguruma

import (
	"fmt"
	"strings"
	"sync"
	"sync/atomic"
	"time"
	"unicode"

	"github.com/eugenioenko/regexp2/v2"
)

const (
	// DefaultMatchTimeout bounds each individual pattern evaluation. regexp2's
	// timeout clock is approximate, so this is a safety limit rather than a
	// real-time scheduling guarantee.
	DefaultMatchTimeout = 20 * time.Millisecond
	neverMatch          = `(?!)`
)

// String is the immutable input wrapper corresponding to OnigString.
type String struct {
	content string
	runes   []rune
	id      uint64
}

var nextStringID atomic.Uint64

// NewString prepares text for repeated scanner searches.
func NewString(content string) *String {
	return &String{content: content, runes: []rune(content), id: nextStringID.Add(1)}
}

// Content returns the original UTF-8 string.
func (s *String) Content() string {
	if s == nil {
		return ""
	}
	return s.content
}

// Len returns the number of runes in the string.
func (s *String) Len() int {
	if s == nil {
		return 0
	}
	return len(s.runes)
}

// FindOption mirrors vscode-textmate's Oniguruma search flags.
type FindOption uint8

const (
	FindOptionNone FindOption = 0
	// FindOptionNotBeginString makes \A fail for this search.
	FindOptionNotBeginString FindOption = 1 << 0
	// FindOptionNotEndString makes \z and \Z fail for this search.
	FindOptionNotEndString FindOption = 1 << 1
	// FindOptionNotBeginPosition makes \G fail for this search.
	FindOptionNotBeginPosition FindOption = 1 << 2
	// FindOptionDebugCall is accepted for API parity and has no effect.
	FindOptionDebugCall FindOption = 1 << 3
)

// Capture is a numbered capture span in rune offsets. An unmatched capture has
// Start and End set to -1.
type Capture struct {
	Start int
	End   int
}

// Match is the winning pattern and all of its numbered captures. Captures[0]
// is the whole match.
type Match struct {
	Index    int
	Captures []Capture
}

// Scanner implements ordered OnigScanner searches.
type Scanner interface {
	FindNextMatch(s *String, start int, opts FindOption) *Match
	Diagnostics() []Diagnostic
}

// DiagnosticKind classifies patterns that have been degraded to never match.
type DiagnosticKind string

const (
	DiagnosticCompileError      DiagnosticKind = "compile_error"
	DiagnosticUnsupportedSyntax DiagnosticKind = "unsupported_syntax"
	DiagnosticMatchTimeout      DiagnosticKind = "match_timeout"
	DiagnosticMatchError        DiagnosticKind = "match_error"
	DiagnosticPanic             DiagnosticKind = "regex_panic"
)

// Diagnostic describes a pattern that could not safely participate in a scan.
type Diagnostic struct {
	Kind         DiagnosticKind
	PatternIndex int
	Pattern      string
	Translated   string
	Message      string
}

// partialTranslationError reports syntax that was neutralized locally while
// leaving the rest of the pattern usable. This is intentionally distinct from
// a fatal translation error: a TextMate pattern can contain supported
// alternatives alongside an unsupported construct, and those alternatives
// should continue to participate in scanning.
type partialTranslationError struct {
	message string
}

func (e *partialTranslationError) Error() string { return e.message }

type regexpPanicError struct {
	operation string
	value     any
}

func (e regexpPanicError) Error() string {
	return fmt.Sprintf("regexp2 panicked while %s: %v", e.operation, e.value)
}

// compileHook and matchHook let package tests inject panics at the regexp2
// boundaries. They are nil in production.
var (
	compileHook func(string)
	matchHook   func(*regexp2.Regexp)
)

// ScannerOption configures an OnigScanner.
type ScannerOption func(*scannerConfig)

type scannerConfig struct {
	matchTimeout      time.Duration
	diagnosticHandler func(Diagnostic)
	patternCache      *PatternCache
}

// WithMatchTimeout sets the limit for each pattern evaluation. Non-positive
// values retain DefaultMatchTimeout so scanners always remain bounded.
func WithMatchTimeout(timeout time.Duration) ScannerOption {
	return func(c *scannerConfig) {
		if timeout > 0 {
			c.matchTimeout = timeout
		}
	}
}

// WithDiagnosticHandler registers a callback for each new deduplicated
// diagnostic. The callback runs synchronously after the scanner releases its
// internal diagnostic lock.
func WithDiagnosticHandler(handler func(Diagnostic)) ScannerOption {
	return func(c *scannerConfig) {
		c.diagnosticHandler = handler
	}
}

// OnigScanner is the regexp2-backed Scanner implementation. Searches on one
// scanner are serialized; distinct scanners may search concurrently.
type OnigScanner struct {
	patterns []*compiledPattern

	mu       sync.Mutex
	searches [][variantCount]cachedSearch
	scratch  []regexp2.CaptureIndex
	captures []Capture
	winner   []Capture
	disabled []bool

	diagnosticsMu sync.Mutex
	diagnostics   []Diagnostic
	diagnosticSet map[diagnosticKey]struct{}
	diagnosticFn  func(Diagnostic)
}

type diagnosticKey struct {
	kind  DiagnosticKind
	index int
}

// compiledPattern is immutable after construction so a PatternCache can share
// it between scanners.
type compiledPattern struct {
	source     string
	translated string
	timeout    time.Duration
	hasA       bool
	hasG       bool
	hasZ       bool
	broken     bool
	variants   [variantCount]*regexp2.Regexp
	// compileDiagnostics are replayed with the pattern's index in each scanner
	// that uses it.
	compileDiagnostics []Diagnostic
}

type cachedSearch struct {
	inputID  uint64
	start    int
	captures []Capture
	failed   bool
	valid    bool
}

const (
	variantAllowA uint8 = 1 << iota
	variantAllowG
	variantAllowZ
	variantCount = 8
)

type patternCacheKey struct {
	source  string
	timeout time.Duration
}

// PatternCache shares translated and compiled patterns between scanners.
// TextMate rules with back-referenced end patterns rebuild their scanner
// whenever the resolved end changes; with a cache only the new end pattern is
// compiled. A PatternCache is safe for concurrent use.
type PatternCache struct {
	mu       sync.Mutex
	patterns map[patternCacheKey]*compiledPattern
	limit    int
}

// NewPatternCache returns a cache that holds at most limit patterns, dropping
// all entries when full. A non-positive limit selects 8192.
func NewPatternCache(limit int) *PatternCache {
	if limit <= 0 {
		limit = 8192
	}
	return &PatternCache{patterns: make(map[patternCacheKey]*compiledPattern), limit: limit}
}

func (c *PatternCache) get(source string, timeout time.Duration) *compiledPattern {
	if c == nil {
		return newCompiledPattern(source, timeout)
	}
	key := patternCacheKey{source: source, timeout: timeout}
	c.mu.Lock()
	pattern, ok := c.patterns[key]
	c.mu.Unlock()
	if ok {
		return pattern
	}
	pattern = newCompiledPattern(source, timeout)
	c.mu.Lock()
	if len(c.patterns) >= c.limit {
		clear(c.patterns)
	}
	c.patterns[key] = pattern
	c.mu.Unlock()
	return pattern
}

// WithPatternCache shares compiled patterns through cache.
func WithPatternCache(cache *PatternCache) ScannerOption {
	return func(c *scannerConfig) {
		c.patternCache = cache
	}
}

// NewScanner translates and compiles patterns eagerly. A bad pattern is kept
// at its original index as a never-match entry and exposed through Diagnostics.
func NewScanner(sources []string, options ...ScannerOption) *OnigScanner {
	config := scannerConfig{matchTimeout: DefaultMatchTimeout}
	for _, option := range options {
		if option != nil {
			option(&config)
		}
	}

	scanner := &OnigScanner{
		patterns:      make([]*compiledPattern, len(sources)),
		searches:      make([][variantCount]cachedSearch, len(sources)),
		disabled:      make([]bool, len(sources)),
		diagnosticSet: make(map[diagnosticKey]struct{}),
		diagnosticFn:  config.diagnosticHandler,
	}
	for i, source := range sources {
		pattern := config.patternCache.get(source, config.matchTimeout)
		scanner.patterns[i] = pattern
		for _, diagnostic := range pattern.compileDiagnostics {
			diagnostic.PatternIndex = i
			scanner.addDiagnostic(diagnostic)
		}
	}
	return scanner
}

func newCompiledPattern(source string, timeout time.Duration) *compiledPattern {
	translated, err := translatePattern(translateInlineOptions(source))
	_, partiallyTranslated := err.(*partialTranslationError)
	checkSource := stripSyntaxComments(source)
	if duplicate := findDuplicateNamedGroup(checkSource); duplicate != "" {
		err = fmt.Errorf("duplicate named capture %q is not supported", duplicate)
		partiallyTranslated = false
	}
	pattern := &compiledPattern{
		source:     source,
		translated: translated,
		timeout:    timeout,
	}
	if err != nil {
		pattern.compileDiagnostics = append(pattern.compileDiagnostics, Diagnostic{
			Kind: DiagnosticUnsupportedSyntax, Pattern: source,
			Translated: translated, Message: err.Error(),
		})
		if !partiallyTranslated {
			pattern.broken = true
			return pattern
		}
	}
	pattern.hasA = containsAnchor(translated, 'A')
	pattern.hasG = containsAnchor(translated, 'G')
	pattern.hasZ = containsAnchor(translated, 'z') || containsAnchor(translated, 'Z')
	compileErr := func() (err error) {
		defer func() {
			if value := recover(); value != nil {
				err = regexpPanicError{operation: "compiling", value: value}
			}
		}()
		return pattern.compileVariants()
	}()
	if compileErr != nil {
		pattern.broken = true
		pattern.variants = [variantCount]*regexp2.Regexp{}
		kind := DiagnosticCompileError
		if _, panicked := compileErr.(regexpPanicError); panicked {
			kind = DiagnosticPanic
		}
		pattern.compileDiagnostics = append(pattern.compileDiagnostics, Diagnostic{
			Kind: kind, Pattern: source,
			Translated: translated, Message: compileErr.Error(),
		})
	}
	return pattern
}

// Len returns the number of source patterns, including degraded patterns.
func (s *OnigScanner) Len() int {
	if s == nil {
		return 0
	}
	return len(s.patterns)
}

// Diagnostics returns a snapshot of compile, compatibility, and timeout
// diagnostics. A timeout is recorded at most once per pattern.
func (s *OnigScanner) Diagnostics() []Diagnostic {
	if s == nil {
		return nil
	}
	s.diagnosticsMu.Lock()
	defer s.diagnosticsMu.Unlock()
	return append([]Diagnostic(nil), s.diagnostics...)
}

func (s *OnigScanner) addDiagnostic(d Diagnostic) {
	key := diagnosticKey{kind: d.Kind, index: d.PatternIndex}
	s.diagnosticsMu.Lock()
	if _, exists := s.diagnosticSet[key]; exists {
		s.diagnosticsMu.Unlock()
		return
	}
	s.diagnosticSet[key] = struct{}{}
	s.diagnostics = append(s.diagnostics, d)
	handler := s.diagnosticFn
	s.diagnosticsMu.Unlock()
	if handler != nil {
		handler(d)
	}
}

// FindNextMatch returns the earliest match at or after start, breaking ties by
// the lowest source-pattern index.
func (s *OnigScanner) FindNextMatch(input *String, start int, opts FindOption) *Match {
	match, ok := s.FindNextMatchInto(input, start, opts, nil)
	if !ok {
		return nil
	}
	return &match
}

// FindNextMatchInto is the allocation-reusing form of FindNextMatch. Captures
// are written into dst when its capacity is sufficient; callers may reuse the
// returned capture slice as dst on their next call.
func (s *OnigScanner) FindNextMatchInto(input *String, start int, opts FindOption, dst []Capture) (Match, bool) {
	if s == nil || input == nil || start > len(input.runes) {
		return Match{}, false
	}
	if start < 0 {
		start = 0
	}
	allowA := opts&FindOptionNotBeginString == 0
	allowG := opts&FindOptionNotBeginPosition == 0
	allowZ := opts&FindOptionNotEndString == 0

	s.mu.Lock()
	locked := true
	defer func() {
		if locked {
			s.mu.Unlock()
		}
	}()
	var diagnostics []Diagnostic
	bestIndex := -1
	s.winner = s.winner[:0]
	for index, pattern := range s.patterns {
		if pattern.broken || s.disabled[index] {
			continue
		}
		variant := pattern.variantKey(allowA, allowG, allowZ)
		maxStartExclusive := -1
		if len(s.winner) != 0 {
			// Patterns are visited in tie-breaking order. Once an earlier pattern
			// has matched at q, a later one can only win by starting before q.
			maxStartExclusive = s.winner[0].Start
		}
		// A failed unanchored search from p cannot succeed from a later
		// position on the same input. This is not true for \G: its meaning is
		// the current search start, so the same pattern must be retried when
		// that start advances.
		var captures []Capture
		var err error
		if !pattern.hasG {
			cache := &s.searches[index][variant]
			if cached, failed, ok := cache.lookup(input.id, start); ok {
				if failed {
					continue
				}
				captures = cached
			} else {
				captures, err = s.match(pattern.variants[variant], input.runes, start, maxStartExclusive, cache.captures[:0])
				// A bounded miss only proves that the pattern cannot beat this
				// call's winner. It may still match later on the same input.
				if err == nil && (len(captures) != 0 || maxStartExclusive < 0) {
					cache.inputID = input.id
					cache.start = start
					cache.captures = captures
					cache.failed = len(captures) == 0
					cache.valid = true
				}
			}
		} else {
			captures, err = s.match(pattern.variants[variant], input.runes, start, maxStartExclusive, s.captures[:0])
			s.captures = captures
		}
		if err != nil {
			kind := DiagnosticMatchError
			message := "regexp evaluation failed"
			if panicErr, panicked := err.(regexpPanicError); panicked {
				s.disabled[index] = true
				kind = DiagnosticPanic
				message = panicErr.Error()
			} else if strings.Contains(strings.ToLower(err.Error()), "timeout") {
				kind = DiagnosticMatchTimeout
				message = fmt.Sprintf("match exceeded the configured %s per-pattern timeout", pattern.timeout)
			}
			diagnostics = append(diagnostics, Diagnostic{
				Kind: kind, PatternIndex: index,
				Pattern: pattern.source, Translated: pattern.translated, Message: message,
			})
			continue
		}
		if len(captures) == 0 {
			continue
		}
		if len(s.winner) == 0 || captures[0].Start < s.winner[0].Start {
			bestIndex = index
			s.winner = append(s.winner[:0], captures...)
			if captures[0].Start == start {
				break
			}
		}
	}
	found := len(s.winner) != 0
	var result Match
	if found {
		dst = append(dst[:0], s.winner...)
		result = Match{Index: bestIndex, Captures: dst}
	}
	s.mu.Unlock()
	locked = false
	for _, diagnostic := range diagnostics {
		s.addDiagnostic(diagnostic)
	}
	if !found {
		return Match{}, false
	}
	return result, true
}

// lookup reuses the earliest result of a previous search when the new start
// has not advanced past it. A regexp search from p that found its first match
// at q has the same answer from every p' in [p,q]. A failed search is reusable
// for every later start. Callers exclude \G patterns because their anchor
// moves with the search start.
func (c *cachedSearch) lookup(inputID uint64, start int) ([]Capture, bool, bool) {
	if !c.valid || c.inputID != inputID || start < c.start {
		return nil, false, false
	}
	if c.failed {
		return nil, true, true
	}
	if len(c.captures) != 0 && start <= c.captures[0].Start {
		return c.captures, false, true
	}
	return nil, false, false
}

func (p *compiledPattern) variantKey(allowA, allowG, allowZ bool) uint8 {
	var key uint8
	if !p.hasA || allowA {
		key |= variantAllowA
	}
	if !p.hasG || allowG {
		key |= variantAllowG
	}
	if !p.hasZ || allowZ {
		key |= variantAllowZ
	}
	return key
}

// compileVariants fills every slot of variants: a slot whose anchor the
// pattern lacks aliases the slot where that anchor is allowed.
func (p *compiledPattern) compileVariants() error {
	for key := range uint8(variantCount) {
		allowA := key&variantAllowA != 0
		allowG := key&variantAllowG != 0
		allowZ := key&variantAllowZ != 0
		canonical := p.variantKey(allowA, allowG, allowZ)
		if canonical != key {
			continue
		}
		source := p.translated
		if p.hasA && !allowA {
			source = neutralizeAnchors(source, "A")
		}
		if p.hasG && !allowG {
			source = neutralizeAnchors(source, "G")
		}
		if p.hasZ && !allowZ {
			source = neutralizeAnchors(source, "zZ")
		}
		if compileHook != nil {
			compileHook(source)
		}
		regex, err := regexp2.Compile(source, regexp2.Multiline, regexp2.OptionMaintainCaptureOrder())
		if err != nil {
			return err
		}
		regex.MatchTimeout = p.timeout
		p.variants[key] = regex
	}
	for key := range uint8(variantCount) {
		allowA := key&variantAllowA != 0
		allowG := key&variantAllowG != 0
		allowZ := key&variantAllowZ != 0
		p.variants[key] = p.variants[p.variantKey(allowA, allowG, allowZ)]
	}
	return nil
}

// match must be called with s.mu held because it reuses s.scratch.
func (s *OnigScanner) match(
	regex *regexp2.Regexp,
	text []rune,
	start, maxStartExclusive int,
	dst []Capture,
) (captures []Capture, err error) {
	captures = dst[:0]
	defer func() {
		if value := recover(); value != nil {
			captures = dst[:0]
			err = regexpPanicError{operation: "matching", value: value}
		}
	}()
	if matchHook != nil {
		matchHook(regex)
	}
	var indices []regexp2.CaptureIndex
	if maxStartExclusive < 0 {
		indices, err = regex.FindRunesCaptureIndicesStartingAt(text, start, s.scratch[:0])
	} else {
		indices, err = regex.FindRunesCaptureIndicesStartingAtBefore(text, start, maxStartExclusive, s.scratch[:0])
	}
	if cap(indices) > cap(s.scratch) {
		s.scratch = indices[:0]
	}
	if err != nil || len(indices) == 0 {
		return captures, err
	}
	for _, capture := range indices {
		if capture.RuneIndex < 0 {
			captures = append(captures, Capture{Start: -1, End: -1})
			continue
		}
		captures = append(captures, Capture{Start: capture.RuneIndex, End: capture.RuneIndex + capture.RuneLength})
	}
	return captures, nil
}

func containsAnchor(source string, anchor rune) bool {
	runes := []rune(source)
	inClass := false
	for index := 0; index+1 < len(runes); index++ {
		if !inClass && index+2 < len(runes) && runes[index] == '(' && runes[index+1] == '?' && runes[index+2] == '#' {
			index = commentEnd(runes, index+3) - 1
			continue
		}
		if runes[index] == '[' {
			inClass = true
			continue
		}
		if runes[index] == ']' {
			inClass = false
			continue
		}
		if runes[index] != '\\' {
			continue
		}
		if !inClass && runes[index+1] == anchor {
			return true
		}
		index++
	}
	return false
}

func neutralizeAnchors(source, anchors string) string {
	var result strings.Builder
	runes := []rune(source)
	inClass := false
	for index := 0; index < len(runes); index++ {
		if !inClass && index+2 < len(runes) && runes[index] == '(' && runes[index+1] == '?' && runes[index+2] == '#' {
			end := commentEnd(runes, index+3)
			result.WriteString(string(runes[index:end]))
			index = end - 1
			continue
		}
		if runes[index] == '\\' && index+1 < len(runes) {
			next := runes[index+1]
			if !inClass && strings.ContainsRune(anchors, next) {
				result.WriteString(neverMatch)
				index++
				continue
			}
			result.WriteRune(runes[index])
			result.WriteRune(next)
			index++
			continue
		}
		switch runes[index] {
		case '[':
			inClass = true
		case ']':
			inClass = false
		}
		result.WriteRune(runes[index])
	}
	return result.String()
}

func translateInlineOptions(source string) string {
	runes := []rune(source)
	var result strings.Builder
	extended := false
	extendedStack := make([]bool, 0, 4)
	for index := 0; index < len(runes); {
		if runes[index] == '\\' {
			result.WriteRune(runes[index])
			index++
			if index < len(runes) {
				result.WriteRune(runes[index])
				index++
			}
			continue
		}
		if runes[index] == '[' {
			end := classEnd(runes, index)
			if end == -1 {
				end = len(runes)
			}
			result.WriteString(string(runes[index:end]))
			index = end
			continue
		}
		if extended && runes[index] == '#' {
			for index < len(runes) && runes[index] != '\n' {
				result.WriteRune(runes[index])
				index++
			}
			continue
		}
		if index+2 < len(runes) && runes[index] == '(' && runes[index+1] == '?' && runes[index+2] == '#' {
			end := commentEnd(runes, index+3)
			result.WriteString(string(runes[index:end]))
			index = end
			continue
		}
		if runes[index] == ')' {
			result.WriteRune(')')
			if len(extendedStack) != 0 {
				extended = extendedStack[len(extendedStack)-1]
				extendedStack = extendedStack[:len(extendedStack)-1]
			}
			index++
			continue
		}
		if index+2 >= len(runes) || runes[index] != '(' || runes[index+1] != '?' {
			if runes[index] == '(' {
				extendedStack = append(extendedStack, extended)
			}
			result.WriteRune(runes[index])
			index++
			continue
		}
		end := index + 2
		for end < len(runes) && strings.ContainsRune("imsx-", runes[end]) {
			end++
		}
		if end == index+2 || end >= len(runes) || runes[end] != ')' && runes[end] != ':' {
			extendedStack = append(extendedStack, extended)
			result.WriteRune(runes[index])
			index++
			continue
		}
		result.WriteString("(?")
		newExtended := applyExtendedOption(extended, runes[index+2:end])
		for _, option := range runes[index+2 : end] {
			if option == 'm' {
				option = 's'
			}
			result.WriteRune(option)
		}
		result.WriteRune(runes[end])
		if runes[end] == ':' {
			extendedStack = append(extendedStack, extended)
		}
		extended = newExtended
		index = end + 1
	}
	return result.String()
}

func applyExtendedOption(current bool, options []rune) bool {
	enable := true
	for _, option := range options {
		if option == '-' {
			enable = false
			continue
		}
		if option == 'x' {
			current = enable
		}
	}
	return current
}

func translatePattern(source string) (string, error) {
	source = expandResolvableSubroutineCalls(source)
	runes := []rune(source)
	out := make([]rune, 0, len(runes)+16)
	groupStarts := make([]int, 0, 4)
	extended := false
	extendedStack := make([]bool, 0, 4)
	lastAtom := -1
	firstSubroutine := ""
	for index := 0; index < len(runes); {
		if extended && runes[index] == '#' {
			for index < len(runes) && runes[index] != '\n' {
				out = append(out, runes[index])
				index++
			}
			continue
		}
		switch runes[index] {
		case '\\':
			if subroutine, next, ok := subroutineCallAt(runes, index); ok {
				if firstSubroutine == "" {
					firstSubroutine = subroutine
				}
				// regexp2 has no subexpression-call operator. Neutralize only
				// the unsupported atom so other alternatives in the same
				// TextMate pattern remain usable (for example, the primitive
				// type branch in the Thrift field grammar).
				lastAtom = len(out)
				out = append(out, []rune(neverMatch)...)
				index = next
				continue
			}
			translated, next, err := translateEscape(runes, index)
			if err != nil {
				return string(out), err
			}
			lastAtom = len(out)
			out = append(out, []rune(translated)...)
			index = next
		case '[':
			end := classEnd(runes, index)
			if end == -1 {
				lastAtom = len(out)
				out = append(out, runes[index])
				index++
				continue
			}
			translated, err := translateClassWithIntersection(runes[index:end])
			if err != nil {
				return string(out), err
			}
			lastAtom = len(out)
			out = append(out, []rune(translated)...)
			index = end
		case '(':
			if index+2 < len(runes) && runes[index+1] == '?' && runes[index+2] == '#' {
				end := commentEnd(runes, index+3)
				out = append(out, runes[index:end]...)
				index = end
				continue
			}
			if end, scoped, newExtended, ok := inlineOptionGroup(runes, index, extended); ok {
				if scoped {
					groupStarts = append(groupStarts, len(out))
					extendedStack = append(extendedStack, extended)
					lastAtom = -1
				}
				out = append(out, runes[index:end]...)
				extended = newExtended
				index = end
				continue
			}
			groupStarts = append(groupStarts, len(out))
			extendedStack = append(extendedStack, extended)
			out = append(out, '(')
			lastAtom = -1
			index++
		case ')':
			out = append(out, ')')
			if len(groupStarts) != 0 {
				lastAtom = groupStarts[len(groupStarts)-1]
				groupStarts = groupStarts[:len(groupStarts)-1]
			}
			if len(extendedStack) != 0 {
				extended = extendedStack[len(extendedStack)-1]
				extendedStack = extendedStack[:len(extendedStack)-1]
			}
			index++
		case '|':
			out = append(out, '|')
			lastAtom = -1
			index++
		case '*', '+', '?':
			if lastAtom < 0 || isGroupPrefix(out) {
				out = append(out, runes[index])
				index++
				continue
			}
			out = append(out, runes[index])
			index++
			if index < len(runes) && runes[index] == '+' {
				out = atomicWrap(out, lastAtom)
				index++
			}
		case '{':
			interval, minimum, maximum, end, ok := parseInterval(runes, index)
			if !ok || lastAtom < 0 {
				lastAtom = len(out)
				out = append(out, '{')
				index++
				continue
			}
			if len(out) != 0 && strings.ContainsRune("*+?", out[len(out)-1]) {
				// Oniguruma permits a repeat operator to be quantified by an
				// interval (for example, `x+{0,1}` in the current Go and C++
				// grammars). regexp2 rejects adjacent quantifiers. Group the
				// already-quantified atom so the interval applies to that whole
				// repetition, preserving Oniguruma's nested-repeat behavior.
				out = groupWrap(out, lastAtom)
			}
			if maximum != "" && decimalGreater(minimum, maximum) {
				interval = "{" + maximum + "," + minimum + "}"
				out = append(out, []rune(interval)...)
				out = atomicWrap(out, lastAtom)
				index = end
				continue
			}
			out = append(out, []rune(interval)...)
			index = end
			if index < len(runes) && runes[index] == '+' {
				out = repeatedWrap(out, lastAtom)
				index++
			}
		case '^':
			// Oniguruma does not treat the position immediately after a final
			// newline as the beginning of another line. regexp2's Multiline ^
			// does, which makes a pattern such as ^$ spuriously match at EOF in
			// the synthetic "line\n" strings used by TextMate tokenizers. Keep
			// the true beginning-of-string case, and require a following rune
			// for all other multiline beginnings. The lookbehind deliberately
			// avoids \A so FindOptionNotBeginString continues to affect only
			// an explicit Oniguruma \A anchor.
			lastAtom = -1
			out = append(out, []rune(`(?:(?<![\s\S])|^(?=[\s\S]))`)...)
			index++
		case '$':
			lastAtom = -1
			out = append(out, runes[index])
			index++
		default:
			if len(out) == 0 || out[len(out)-1] != '(' || runes[index] != '?' {
				lastAtom = len(out)
			}
			out = append(out, runes[index])
			index++
		}
	}
	translated := string(out)
	if firstSubroutine != "" {
		return translated, &partialTranslationError{
			message: fmt.Sprintf("oniguruma subroutine call %q was neutralized", firstSubroutine),
		}
	}
	return translated, nil
}

type namedSubroutineDefinition struct {
	name      string
	body      string
	start     int
	end       int
	zeroCount bool
}

// expandResolvableSubroutineCalls converts Oniguruma's named subexpression
// calls to non-capturing copies of their definitions. regexp2 does not expose
// a subexpression-call operator, but expansion is exact for acyclic
// definitions and for the common right-recursive `item self?` form used by
// generated TextMate grammars. Calls involved in any other cycle are left for
// translatePattern to neutralize with an unsupported-syntax diagnostic.
func expandResolvableSubroutineCalls(source string) string {
	definitions := collectNamedSubroutineDefinitions(source)
	if len(definitions) == 0 {
		return source
	}

	memo := make(map[string]string, len(definitions))
	visiting := make(map[string]bool, len(definitions))
	var expandDefinition func(string) (string, bool)
	expandDefinition = func(name string) (string, bool) {
		if expanded, ok := memo[name]; ok {
			return expanded, true
		}
		definition, ok := definitions[name]
		if !ok || visiting[name] {
			return "", false
		}
		visiting[name] = true
		defer delete(visiting, name)

		body := definition.body
		repeat := false
		if prefix, ok := trimRightRecursiveTail(body, name); ok {
			body = prefix
			repeat = true
		}
		expanded, ok := rewriteSubroutineCalls(body, expandDefinition)
		if !ok {
			return "", false
		}
		if repeat {
			expanded = "(?:" + expanded + ")+"
		}
		memo[name] = expanded
		return expanded, true
	}

	runes := []rune(source)
	zeroDefinitions := make(map[int]namedSubroutineDefinition)
	for _, definition := range definitions {
		if definition.zeroCount {
			zeroDefinitions[definition.start] = definition
		}
	}

	var result strings.Builder
	for index := 0; index < len(runes); {
		if definition, ok := zeroDefinitions[index]; ok {
			// Preserve the capture slot while removing dead definition
			// bodies. Their {0} quantifier means they never participate in a
			// match, but retaining a plain group keeps later capture numbers
			// stable and also accepts Oniguruma names containing hyphens.
			result.WriteString("(" + neverMatch + "){0}")
			index = definition.end
			continue
		}
		if call, end, ok := subroutineCallAt(runes, index); ok {
			name := call[3 : len(call)-1]
			if expanded, ok := expandDefinition(name); ok {
				result.WriteString("(?:" + expanded + ")")
				index = end
				continue
			}
		}
		result.WriteRune(runes[index])
		index++
	}
	return result.String()
}

func collectNamedSubroutineDefinitions(source string) map[string]namedSubroutineDefinition {
	runes := []rune(source)
	result := make(map[string]namedSubroutineDefinition)
	for index := 0; index+4 < len(runes); index++ {
		if runes[index] == '\\' {
			index++
			continue
		}
		if runes[index] == '[' {
			if end := classEnd(runes, index); end != -1 {
				index = end - 1
			}
			continue
		}
		if runes[index] != '(' || runes[index+1] != '?' || runes[index+2] != '<' {
			continue
		}
		if runes[index+3] == '=' || runes[index+3] == '!' {
			continue
		}
		nameEnd := index + 3
		for nameEnd < len(runes) && runes[nameEnd] != '>' {
			nameEnd++
		}
		if nameEnd == len(runes) || nameEnd == index+3 {
			continue
		}
		groupEnd := namedGroupEnd(runes, index)
		if groupEnd == -1 {
			continue
		}
		end := groupEnd
		zeroCount := end+3 <= len(runes) && string(runes[end:end+3]) == "{0}"
		if zeroCount {
			end += 3
		}
		name := string(runes[index+3 : nameEnd])
		if _, exists := result[name]; !exists {
			result[name] = namedSubroutineDefinition{
				name:      name,
				body:      string(runes[nameEnd+1 : groupEnd-1]),
				start:     index,
				end:       end,
				zeroCount: zeroCount,
			}
		}
	}
	return result
}

func namedGroupEnd(runes []rune, start int) int {
	depth := 0
	for index := start; index < len(runes); index++ {
		if runes[index] == '\\' {
			index++
			continue
		}
		if runes[index] == '[' {
			if end := classEnd(runes, index); end != -1 {
				index = end - 1
			}
			continue
		}
		if index+2 < len(runes) && runes[index] == '(' && runes[index+1] == '?' && runes[index+2] == '#' {
			index = commentEnd(runes, index+3) - 1
			continue
		}
		switch runes[index] {
		case '(':
			depth++
		case ')':
			depth--
			if depth == 0 {
				return index + 1
			}
		}
	}
	return -1
}

func trimRightRecursiveTail(body, name string) (string, bool) {
	trimmed := strings.TrimRightFunc(body, unicode.IsSpace)
	for _, call := range []string{`\g<` + name + `>?`, `\g'` + name + `'?`} {
		if strings.HasSuffix(trimmed, call) {
			return strings.TrimSuffix(trimmed, call), true
		}
	}
	return body, false
}

func rewriteSubroutineCalls(
	source string,
	expand func(string) (string, bool),
) (string, bool) {
	runes := []rune(source)
	var result strings.Builder
	for index := 0; index < len(runes); {
		if runes[index] == '[' {
			if end := classEnd(runes, index); end != -1 {
				result.WriteString(string(runes[index:end]))
				index = end
				continue
			}
		}
		if call, end, ok := subroutineCallAt(runes, index); ok {
			name := call[3 : len(call)-1]
			expanded, ok := expand(name)
			if !ok {
				return "", false
			}
			result.WriteString("(?:" + expanded + ")")
			index = end
			continue
		}
		result.WriteRune(runes[index])
		index++
	}
	return result.String(), true
}

func subroutineCallAt(runes []rune, start int) (text string, end int, ok bool) {
	if start+3 >= len(runes) || runes[start] != '\\' || runes[start+1] != 'g' {
		return "", start, false
	}
	close := rune('>')
	switch runes[start+2] {
	case '<':
	case '\'':
		close = '\''
	default:
		return "", start, false
	}
	end = start + 3
	for end < len(runes) && runes[end] != close {
		end++
	}
	if end == len(runes) {
		return "", start, false
	}
	end++
	return string(runes[start:end]), end, true
}

func inlineOptionGroup(runes []rune, start int, extended bool) (end int, scoped, newExtended, ok bool) {
	if start+2 >= len(runes) || runes[start] != '(' || runes[start+1] != '?' {
		return 0, false, extended, false
	}
	index := start + 2
	for index < len(runes) && strings.ContainsRune("imsx-", runes[index]) {
		index++
	}
	if index == start+2 || index >= len(runes) || runes[index] != ')' && runes[index] != ':' {
		return 0, false, extended, false
	}
	return index + 1, runes[index] == ':', applyExtendedOption(extended, runes[start+2:index]), true
}

func isGroupPrefix(out []rune) bool {
	return len(out) != 0 && out[len(out)-1] == '('
}

func atomicWrap(out []rune, atomStart int) []rune {
	result := make([]rune, 0, len(out)+4)
	result = append(result, out[:atomStart]...)
	result = append(result, '(', '?', '>')
	result = append(result, out[atomStart:]...)
	result = append(result, ')')
	return result
}

func repeatedWrap(out []rune, atomStart int) []rune {
	result := make([]rune, 0, len(out)+5)
	result = append(result, out[:atomStart]...)
	result = append(result, '(', '?', ':')
	result = append(result, out[atomStart:]...)
	result = append(result, ')', '+')
	return result
}

func groupWrap(out []rune, atomStart int) []rune {
	result := make([]rune, 0, len(out)+4)
	result = append(result, out[:atomStart]...)
	result = append(result, '(', '?', ':')
	result = append(result, out[atomStart:]...)
	result = append(result, ')')
	return result
}

func parseInterval(runes []rune, start int) (text, minimum, maximum string, end int, ok bool) {
	index := start + 1
	minimumStart := index
	for index < len(runes) && runes[index] >= '0' && runes[index] <= '9' {
		index++
	}
	minimum = string(runes[minimumStart:index])
	if index < len(runes) && runes[index] == '}' && minimum != "" {
		return "{" + minimum + "}", minimum, minimum, index + 1, true
	}
	if index >= len(runes) || runes[index] != ',' {
		return "", "", "", start, false
	}
	index++
	maximumStart := index
	for index < len(runes) && runes[index] >= '0' && runes[index] <= '9' {
		index++
	}
	maximum = string(runes[maximumStart:index])
	if index >= len(runes) || runes[index] != '}' || minimum == "" && maximum == "" {
		return "", "", "", start, false
	}
	if minimum == "" {
		minimum = "0"
	}
	return "{" + minimum + "," + maximum + "}", minimum, maximum, index + 1, true
}

func decimalGreater(left, right string) bool {
	left = strings.TrimLeft(left, "0")
	right = strings.TrimLeft(right, "0")
	if len(left) != len(right) {
		return len(left) > len(right)
	}
	return left > right
}

func translateEscape(runes []rune, start int) (string, int, error) {
	if start+1 >= len(runes) {
		return `\`, start + 1, nil
	}
	next := runes[start+1]
	switch next {
	case 'N':
		return `[^\n]`, start + 2, nil
	case 'h':
		return `[0-9A-Fa-f]`, start + 2, nil
	case 'H':
		return `[^0-9A-Fa-f]`, start + 2, nil
	case 'x':
		if start+2 < len(runes) && runes[start+2] == '{' {
			end := start + 3
			for end < len(runes) && runes[end] != '}' {
				end++
			}
			if end == len(runes) {
				return "", end, fmt.Errorf("unterminated hexadecimal escape")
			}
			digits := string(runes[start+3 : end])
			value, ok := parseHex(digits)
			if !ok || value > unicode.MaxRune {
				return "", end + 1, fmt.Errorf("invalid hexadecimal escape \\x{%s}", digits)
			}
			if value <= 0xffff {
				return fmt.Sprintf(`\u%04X`, value), end + 1, nil
			}
			return fmt.Sprintf(`\x{%X}`, value), end + 1, nil
		}
	case 'p', 'P':
		if start+2 < len(runes) && runes[start+2] == '{' {
			end := start + 3
			for end < len(runes) && runes[end] != '}' {
				end++
			}
			if end == len(runes) {
				return "", end, fmt.Errorf("unterminated Unicode property")
			}
			return translateProperty(string(runes[start+3:end]), next == 'P'), end + 1, nil
		}
	}
	return string(runes[start : start+2]), start + 2, nil
}

func parseHex(digits string) (rune, bool) {
	if digits == "" || len(digits) > 8 {
		return 0, false
	}
	var value rune
	for _, digit := range digits {
		value *= 16
		switch {
		case digit >= '0' && digit <= '9':
			value += digit - '0'
		case digit >= 'a' && digit <= 'f':
			value += digit - 'a' + 10
		case digit >= 'A' && digit <= 'F':
			value += digit - 'A' + 10
		default:
			return 0, false
		}
	}
	return value, true
}

var onigurumaProperties = map[string]string{
	"alnum":           `\p{L}\p{Nl}\p{Other_Alphabetic}\p{Nd}`,
	"alpha":           `\p{L}\p{Nl}\p{Other_Alphabetic}`,
	"alphabetic":      `\p{L}\p{Nl}\p{Other_Alphabetic}`,
	"any":             `\s\S`,
	"ascii":           `\x00-\x7F`,
	"assigned":        `\P{Cn}`,
	"blank":           `\t\p{Zs}`,
	"cntrl":           `\p{Cc}`,
	"control":         `\p{Cc}`,
	"decimalnumber":   `\p{Nd}`,
	"digit":           `\p{Nd}`,
	"graph":           `\p{L}\p{M}\p{N}\p{P}\p{S}`,
	"letter":          `\p{L}`,
	"lower":           `\p{Ll}\p{Other_Lowercase}`,
	"lowercaseletter": `\p{Ll}`,
	"mark":            `\p{M}`,
	"number":          `\p{N}`,
	"print":           `\p{L}\p{M}\p{N}\p{P}\p{S}\p{Zs}`,
	"punct":           `\p{P}`,
	"punctuation":     `\p{P}`,
	"separator":       `\p{Z}`,
	"space":           `\s`,
	"titlecaseletter": `\p{Lt}`,
	"upper":           `\p{Lu}\p{Other_Uppercase}`,
	"uppercaseletter": `\p{Lu}`,
	"word":            `\p{L}\p{M}\p{N}\p{Pc}`,
	"xdigit":          `0-9A-Fa-f`,
	"xids":            `\p{L}\p{Nl}\p{Other_ID_Start}`,
	"xidc":            `\p{L}\p{Nl}\p{Mn}\p{Mc}\p{Nd}\p{Pc}\p{Other_ID_Start}\p{Other_ID_Continue}`,
}

var canonicalUnicodeProperties = func() map[string]string {
	result := make(map[string]string)
	for name := range unicode.Categories {
		result[normalizePropertyName(name)] = name
	}
	for name := range unicode.Scripts {
		result[normalizePropertyName(name)] = name
	}
	for name := range unicode.Properties {
		result[normalizePropertyName(name)] = name
	}
	return result
}()

func translateProperty(name string, negate bool) string {
	normalized := propertyKey(name)
	body, exists := onigurumaProperties[normalized]
	if !exists {
		if canonical, ok := canonicalUnicodeProperties[normalized]; ok {
			if negate {
				return `\P{` + canonical + `}`
			}
			return `\p{` + canonical + `}`
		}
		if negate {
			return `\P{` + name + `}`
		}
		return `\p{` + name + `}`
	}
	if strings.HasPrefix(body, `\p{`) && strings.Count(body, `\p{`) == 1 {
		if negate {
			return `\P{` + body[3:]
		}
		return body
	}
	if body == `\P{Cn}` {
		if negate {
			return `\p{Cn}`
		}
		return body
	}
	if normalized == "any" {
		if negate {
			return neverMatch
		}
		return `[` + body + `]`
	}
	if negate {
		return `[^` + body + `]`
	}
	return `[` + body + `]`
}

func normalizePropertyName(name string) string {
	return strings.Map(func(r rune) rune {
		switch r {
		case '_', '-', ' ', '\t':
			return -1
		default:
			return unicode.ToLower(r)
		}
	}, name)
}

func propertyKey(name string) string {
	normalized := normalizePropertyName(name)
	if _, exists := onigurumaProperties[normalized]; exists {
		return normalized
	}
	if _, exists := canonicalUnicodeProperties[normalized]; exists {
		return normalized
	}
	return strings.TrimPrefix(normalized, "is")
}

func classEnd(runes []rune, start int) int {
	index := start + 1
	if index < len(runes) && runes[index] == '^' {
		index++
	}
	if index < len(runes) && runes[index] == ']' {
		index++
	}
	for index < len(runes) {
		if runes[index] == '\\' {
			index += 2
			continue
		}
		if runes[index] == '[' && index+1 < len(runes) && runes[index+1] == ':' {
			index += 2
			for index+1 < len(runes) && (runes[index] != ':' || runes[index+1] != ']') {
				index++
			}
			if index+1 < len(runes) {
				index += 2
			}
			continue
		}
		if runes[index] == '[' {
			// Each nested class has its own leading-] rule. Recursing matters
			// for forms such as `[^]...]`, where that first ] is a literal and
			// must not be mistaken for the nested class terminator. If no
			// complete nested class follows, `[` is simply literal in the outer
			// class (the conventional `[[]` spelling).
			if end := classEnd(runes, index); end != -1 {
				index = end
				continue
			}
		}
		if runes[index] == ']' {
			return index + 1
		}
		index++
	}
	return -1
}

func translateClass(source string) (string, error) {
	inner := []rune(source[1 : len(source)-1])
	outerNegated := len(inner) != 0 && inner[0] == '^'
	if outerNegated {
		inner = inner[1:]
	}

	var positive strings.Builder
	negative := make([]string, 0, 1)
	nestedClasses := make([]string, 0, 1)
	for index := 0; index < len(inner); {
		if inner[index] == '&' && index+1 < len(inner) && inner[index+1] == '&' {
			return "", fmt.Errorf("oniguruma character-class intersection is not supported")
		}
		if inner[index] == '[' && (index+1 >= len(inner) || inner[index+1] != ':') {
			end := classEnd(inner, index)
			if end == -1 {
				return "", fmt.Errorf("unterminated nested character class")
			}
			nestedPattern, err := translateClass(string(inner[index:end]))
			if err != nil {
				return "", err
			}
			nestedClasses = append(nestedClasses, nestedPattern)
			index = end
			continue
		}
		if inner[index] == '[' && index+1 < len(inner) && inner[index+1] == ':' {
			end := index + 2
			for end+1 < len(inner) && (inner[end] != ':' || inner[end+1] != ']') {
				end++
			}
			if end+1 >= len(inner) {
				positive.WriteRune(inner[index])
				index++
				continue
			}
			name := string(inner[index+2 : end])
			posixNegated := strings.HasPrefix(name, "^")
			name = strings.TrimPrefix(name, "^")
			mapped, exists := onigurumaProperties[normalizePropertyName(name)]
			if !exists {
				return "", fmt.Errorf("unknown POSIX character class %q", name)
			}
			if posixNegated {
				negative = append(negative, mapped)
			} else {
				positive.WriteString(mapped)
			}
			index = end + 2
			continue
		}
		if inner[index] == '\\' && index+1 < len(inner) {
			if inner[index+1] == '-' && index+3 < len(inner) && inner[index+2] == '-' {
				// Oniguruma allows an escaped hyphen to be a range endpoint:
				// `[\--9]` spans '-' through '9'. regexp2 interprets `\-` as
				// an unconditionally literal member, so spell the endpoint by
				// code point and leave the following range separator intact.
				positive.WriteString(`\x2D`)
				index += 2
				continue
			}
			if isClassIdentityEscape(inner[index+1]) {
				// Oniguruma accepts a backslash before punctuation that is
				// already literal in a character class. regexp2 rejects some
				// of those identity escapes (notably \_), so emit the literal
				// rune while retaining escapes for class metacharacters such
				// as '-', ']', '^', and '\\'.
				positive.WriteRune(inner[index+1])
				index += 2
				continue
			}
			switch inner[index+1] {
			case 'A', 'G', 'Z', 'z':
				positive.WriteRune(inner[index+1])
				index += 2
				continue
			case 'h':
				positive.WriteString("0-9A-Fa-f")
				index += 2
				continue
			case 'H':
				negative = append(negative, "0-9A-Fa-f")
				index += 2
				continue
			case 'p', 'P':
				if index+2 < len(inner) && inner[index+2] == '{' {
					end := index + 3
					for end < len(inner) && inner[end] != '}' {
						end++
					}
					if end == len(inner) {
						return "", fmt.Errorf("unterminated Unicode property")
					}
					body := propertyClassBody(string(inner[index+3 : end]))
					if inner[index+1] == 'P' {
						negative = append(negative, body)
					} else {
						positive.WriteString(body)
					}
					index = end + 1
					continue
				}
			}
			translated, next, err := translateEscape(inner, index)
			if err != nil {
				return "", err
			}
			positive.WriteString(translated)
			index = next
			continue
		}
		positive.WriteRune(inner[index])
		index++
	}

	positiveBody := positive.String()
	if len(nestedClasses) == 0 {
		if len(negative) == 0 {
			if outerNegated {
				return "[^" + positiveBody + "]", nil
			}
			return "[" + positiveBody + "]", nil
		}
		if outerNegated && len(negative) > 1 {
			return "", fmt.Errorf("multiple negated members in an outer-negated character class are not supported")
		}
		if outerNegated {
			if positiveBody == "" {
				return "[" + negative[0] + "]", nil
			}
			return "[" + negative[0] + "-[" + positiveBody + "]]", nil
		}
		parts := make([]string, 0, len(negative)+1)
		for _, body := range negative {
			parts = append(parts, "[^"+body+"]")
		}
		if positiveBody != "" {
			parts = append(parts, "["+positiveBody+"]")
		}
		if len(parts) == 1 {
			return parts[0], nil
		}
		return "(?:" + strings.Join(parts, "|") + ")", nil
	}

	parts := make([]string, 0, len(negative)+len(nestedClasses)+1)
	for _, body := range negative {
		parts = append(parts, "[^"+body+"]")
	}
	if positiveBody != "" {
		parts = append(parts, "["+positiveBody+"]")
	}
	parts = append(parts, nestedClasses...)

	var union string
	switch len(parts) {
	case 0:
		union = neverMatch
	case 1:
		union = parts[0]
	default:
		union = "(?:" + strings.Join(parts, "|") + ")"
	}
	if outerNegated {
		// Oniguruma negates the union represented by every direct and nested
		// class member. Express that set complement as a one-rune negative
		// lookahead so nested complements remain exact as well.
		return "(?:(?!(?:" + union + "))[\\s\\S])", nil
	}
	return union, nil
}

func isClassIdentityEscape(value rune) bool {
	switch value {
	case '!', '"', '#', '$', '%', '&', '\'', '(', ')', '*', '+', ',', '.', '/',
		':', ';', '<', '=', '>', '?', '@', '_', '`', '{', '|', '}', '~':
		return true
	default:
		return false
	}
}

// stripSyntaxComments removes text that cannot contain regex constructs. It
// tracks scoped and unscoped extended-mode options so compatibility checks do
// not mistake comment examples for active character classes or named groups.
func stripSyntaxComments(source string) string {
	runes := []rune(source)
	var result strings.Builder
	extended := false
	extendedStack := make([]bool, 0, 4)
	for index := 0; index < len(runes); {
		if runes[index] == '\\' {
			result.WriteRune(runes[index])
			index++
			if index < len(runes) {
				result.WriteRune(runes[index])
				index++
			}
			continue
		}
		if runes[index] == '[' {
			end := classEnd(runes, index)
			if end == -1 {
				end = len(runes)
			}
			result.WriteString(string(runes[index:end]))
			index = end
			continue
		}
		if extended && runes[index] == '#' {
			for index < len(runes) && runes[index] != '\n' {
				index++
			}
			continue
		}
		if index+2 < len(runes) && runes[index] == '(' && runes[index+1] == '?' && runes[index+2] == '#' {
			index = commentEnd(runes, index+3)
			continue
		}
		if runes[index] == ')' {
			result.WriteRune(runes[index])
			if len(extendedStack) != 0 {
				extended = extendedStack[len(extendedStack)-1]
				extendedStack = extendedStack[:len(extendedStack)-1]
			}
			index++
			continue
		}
		if end, scoped, newExtended, ok := inlineOptionGroup(runes, index, extended); ok {
			result.WriteString(string(runes[index:end]))
			if scoped {
				extendedStack = append(extendedStack, extended)
			}
			extended = newExtended
			index = end
			continue
		}
		if runes[index] == '(' {
			extendedStack = append(extendedStack, extended)
		}
		result.WriteRune(runes[index])
		index++
	}
	return result.String()
}

// translateClassWithIntersection rewrites Oniguruma's `[A&&B]`, which regexp2
// lacks. A class consumes exactly one rune, so the intersection is that rune
// matching A while lookaheads check the other operands: `(?:(?=B)A)`. In a
// lookbehind the class is matched first and the lookahead then inspects the
// same rune, so the rewrite holds in both directions.
func translateClassWithIntersection(class []rune) (string, error) {
	body := class[1 : len(class)-1]
	negated := len(body) > 0 && body[0] == '^'
	if negated {
		body = body[1:]
	}
	operands := splitClassIntersection(body)
	if len(operands) == 1 {
		return translateClass(string(class))
	}
	var translated []string
	for _, operand := range operands {
		if len(operand) == 0 {
			return "", fmt.Errorf("empty operand in oniguruma character-class intersection")
		}
		source := "[" + string(operand) + "]"
		if operand[0] == '[' && classEnd(operand, 0) == len(operand) {
			source = string(operand)
		}
		class, err := translateClassWithIntersection([]rune(source))
		if err != nil {
			return "", err
		}
		translated = append(translated, class)
	}
	var match strings.Builder
	for _, operand := range translated[1:] {
		match.WriteString("(?=" + operand + ")")
	}
	match.WriteString(translated[0])
	if negated {
		return `(?:(?!` + match.String() + `)[\s\S])`, nil
	}
	return "(?:" + match.String() + ")", nil
}

// splitClassIntersection splits a class body at `&&` outside nested classes.
func splitClassIntersection(body []rune) [][]rune {
	var operands [][]rune
	start := 0
	for index := 0; index < len(body); {
		switch {
		case body[index] == '\\':
			index += 2
		case body[index] == '[':
			if end := classEnd(body, index); end != -1 {
				index = end
			} else {
				index++
			}
		case body[index] == '&' && index+1 < len(body) && body[index+1] == '&':
			operands = append(operands, body[start:index])
			index += 2
			start = index
		default:
			index++
		}
	}
	return append(operands, body[start:])
}

func propertyClassBody(name string) string {
	normalized := propertyKey(name)
	if body, exists := onigurumaProperties[normalized]; exists {
		return body
	}
	if canonical, exists := canonicalUnicodeProperties[normalized]; exists {
		return `\p{` + canonical + `}`
	}
	return `\p{` + name + `}`
}

func findDuplicateNamedGroup(source string) string {
	runes := []rune(source)
	names := make(map[string]struct{})
	inClass := false
	for index := 0; index+3 < len(runes); index++ {
		if runes[index] == '\\' {
			index++
			continue
		}
		if runes[index] == '[' {
			inClass = true
			continue
		}
		if runes[index] == ']' {
			inClass = false
			continue
		}
		if inClass || runes[index] != '(' || runes[index+1] != '?' {
			continue
		}
		var close rune
		nameStart := index + 3
		switch runes[index+2] {
		case '<':
			if runes[nameStart] == '=' || runes[nameStart] == '!' {
				continue
			}
			close = '>'
		case '\'':
			close = '\''
		default:
			continue
		}
		nameEnd := nameStart
		for nameEnd < len(runes) && runes[nameEnd] != close {
			nameEnd++
		}
		if nameEnd == len(runes) {
			continue
		}
		name := string(runes[nameStart:nameEnd])
		if _, exists := names[name]; exists {
			return name
		}
		names[name] = struct{}{}
	}
	return ""
}

func commentEnd(runes []rune, start int) int {
	for index := start; index < len(runes); index++ {
		if runes[index] == '\\' {
			index++
			continue
		}
		if runes[index] == ')' {
			return index + 1
		}
	}
	return len(runes)
}
