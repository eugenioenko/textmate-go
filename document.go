package textmate

import (
	"errors"
	"sync"
)

// DefaultDocumentCacheCapacity is the number of distinct line/state results
// retained by a Document when DocumentOptions.CacheCapacity is zero.
const DefaultDocumentCacheCapacity = 20_000

// ErrInvalidLineRange is returned when ReplaceLines receives a range outside
// the current document or one whose start is after its end.
var ErrInvalidLineRange = errors.New("textmate: invalid document line range")

// DocumentOptions configures a Document. Options are fixed at construction.
type DocumentOptions struct {
	// TokenizeOptions is applied to every line tokenized by the document.
	TokenizeOptions TokenizeOptions

	// CacheCapacity bounds cached results keyed by line text and canonical
	// start state. Zero selects DefaultDocumentCacheCapacity. A negative value
	// disables result caching.
	CacheCapacity int
}

// Document incrementally tokenizes an editable sequence of lines. It owns a
// shallow copy of the supplied line slice, carries start states between lines,
// and reuses an unchanged suffix after an edit as soon as tokenizer state
// converges with the previous suffix.
//
// All methods are safe to call concurrently. Results returned by Line are
// immutable, library-owned views and must not be modified; this permits cache
// hits without allocating. A result remains valid after later calls and cache
// eviction. Copy the Tokens slice before changing Token fields. ScopeStack and
// Token.Scopes remain library-owned and read-only, as they are for
// Grammar.TokenizeLine.
type Document struct {
	mu      sync.Mutex
	grammar *Grammar
	options TokenizeOptions
	lines   []string

	// states[i] is the state in which line i starts. stateSrc[i] is the text
	// of line i that produced states[i+1]. The table is populated lazily.
	states   []*StateStack
	stateSrc []string
	tail     *documentTail

	cacheCapacity int
	cache         map[documentCacheKey]LineResult
	cacheOrder    []documentCacheKey
	cacheNext     int
}

type documentTail struct {
	states   []*StateStack
	stateSrc []string
	oldStart int
	newStart int
}

type documentCacheKey struct {
	line  string
	state *StateStack
}

// NewDocument creates an empty document backed by grammar. A nil grammar is
// valid: lines remain addressable, but Line returns empty tokenization results.
func NewDocument(grammar *Grammar, options DocumentOptions) *Document {
	capacity := options.CacheCapacity
	if capacity == 0 {
		capacity = DefaultDocumentCacheCapacity
	}
	document := &Document{
		grammar:       grammar,
		options:       options.TokenizeOptions,
		states:        []*StateStack{nil},
		cacheCapacity: capacity,
	}
	return document
}

// Len returns the current number of lines.
func (d *Document) Len() int {
	if d == nil {
		return 0
	}
	d.mu.Lock()
	defer d.mu.Unlock()
	return len(d.lines)
}

// SetLines replaces the complete document with a shallow copy of lines.
// Tokenizer states above the common prefix are retained. States in the common
// suffix become reusable once recomputation reaches an equal canonical state.
func (d *Document) SetLines(lines []string) {
	if d == nil {
		return
	}
	owned := append([]string(nil), lines...)
	d.mu.Lock()
	defer d.mu.Unlock()
	d.setLinesLocked(owned)
}

// ReplaceLines replaces the half-open line range [start, end) with a shallow
// copy of replacement. It returns ErrInvalidLineRange without changing the
// document when the range is invalid.
func (d *Document) ReplaceLines(start, end int, replacement []string) error {
	if d == nil {
		return ErrInvalidLineRange
	}
	d.mu.Lock()
	defer d.mu.Unlock()
	if start < 0 || start > end || end > len(d.lines) {
		return ErrInvalidLineRange
	}

	lines := make([]string, 0, len(d.lines)-(end-start)+len(replacement))
	lines = append(lines, d.lines[:start]...)
	lines = append(lines, replacement...)
	lines = append(lines, d.lines[end:]...)
	d.setLinesLocked(lines)
	return nil
}

// Line returns the tokenization result for line index. The boolean is false
// when index is outside the current document. Requesting a line computes any
// missing start states above it and records the line's resulting state.
//
// The returned result and its Tokens slice are immutable, library-owned views.
// They remain valid after later calls and cache eviction. See Document's
// ownership documentation.
func (d *Document) Line(index int) (LineResult, bool) {
	if d == nil {
		return LineResult{}, false
	}
	d.mu.Lock()
	defer d.mu.Unlock()
	if index < 0 || index >= len(d.lines) {
		return LineResult{}, false
	}

	state := d.stateAtLocked(index)
	result := d.tokenizeLocked(d.lines[index], state)
	if len(d.states) == index+1 {
		d.states = append(d.states, result.RuleStack)
		d.stateSrc = append(d.stateSrc, d.lines[index])
		d.spliceTailLocked(index + 1)
	} else if !documentStatesEqual(d.states[index+1], result.RuleStack) {
		d.replaceMaterializedResultLocked(index, result.RuleStack)
	}
	return result, true
}

// InvalidateFrom discards materialized states for line index and all following
// lines, and clears the result cache. Index may equal Len. The next Line or
// StateAt call recomputes the discarded range.
//
// Time-limit stops are never cached, so calling Line again retries that line
// directly. InvalidateFrom is useful when a stopped line has already supplied
// state to later lines and the caller wants the complete tail recomputed.
func (d *Document) InvalidateFrom(index int) error {
	if d == nil {
		return ErrInvalidLineRange
	}
	d.mu.Lock()
	defer d.mu.Unlock()
	if index < 0 || index > len(d.lines) {
		return ErrInvalidLineRange
	}

	keep := min(index, len(d.stateSrc))
	d.states = d.states[: keep+1 : keep+1]
	d.stateSrc = d.stateSrc[:keep:keep]
	d.tail = nil
	clear(d.cache)
	d.cacheOrder = d.cacheOrder[:0]
	d.cacheNext = 0
	return nil
}

// StateAt returns the state in which line index starts. Index may equal Len,
// in which case StateAt returns the state immediately after the final line.
// The boolean is false for all other out-of-range indices.
func (d *Document) StateAt(index int) (*StateStack, bool) {
	if d == nil {
		return nil, false
	}
	d.mu.Lock()
	defer d.mu.Unlock()
	if index < 0 || index > len(d.lines) {
		return nil, false
	}
	return d.stateAtLocked(index), true
}

func (d *Document) stateAtLocked(index int) *StateStack {
	for len(d.states) <= index {
		lineIndex := len(d.states) - 1
		result := d.tokenizeLocked(d.lines[lineIndex], d.states[lineIndex])
		d.states = append(d.states, result.RuleStack)
		d.stateSrc = append(d.stateSrc, d.lines[lineIndex])
		d.spliceTailLocked(lineIndex + 1)
	}
	return d.states[index]
}

func (d *Document) tokenizeLocked(line string, state *StateStack) LineResult {
	key := documentCacheKey{line: line, state: state}
	if result, ok := d.cache[key]; ok {
		return result
	}

	var result LineResult
	if d.grammar != nil {
		result = d.grammar.TokenizeLineWithOptions(line, state, d.options)
	}
	// A time stop can be caused by transient lock or CPU pressure. Retaining it
	// would turn a soft per-call budget into a permanent partial result. A line
	// limit is deterministic for this document's fixed options and is safe to
	// cache.
	if d.cacheCapacity > 0 && result.StoppedReason != StopReasonTimeLimit {
		d.cacheResultLocked(key, result)
	}
	return result
}

func (d *Document) replaceMaterializedResultLocked(index int, state *StateStack) {
	oldStates := d.states
	oldStateSrc := d.stateSrc
	d.states = oldStates[: index+1 : index+1]
	d.stateSrc = oldStateSrc[:index:index]
	d.tail = &documentTail{
		states:   oldStates,
		stateSrc: oldStateSrc,
		oldStart: index + 1,
		newStart: index + 1,
	}
	d.states = append(d.states, state)
	d.stateSrc = append(d.stateSrc, d.lines[index])
	d.spliceTailLocked(index + 1)
}

func (d *Document) cacheResultLocked(key documentCacheKey, result LineResult) {
	if d.cache == nil {
		d.cache = make(map[documentCacheKey]LineResult)
	}
	if len(d.cacheOrder) < d.cacheCapacity {
		d.cacheOrder = append(d.cacheOrder, key)
	} else {
		oldest := d.cacheOrder[d.cacheNext]
		delete(d.cache, oldest)
		d.cacheOrder[d.cacheNext] = key
		d.cacheNext = (d.cacheNext + 1) % d.cacheCapacity
	}
	d.cache[key] = result
}

func (d *Document) setLinesLocked(lines []string) {
	if equalDocumentLines(d.lines, lines) {
		return
	}

	oldLines := d.lines
	oldStates := d.states
	oldStateSrc := d.stateSrc

	prefix := commonLinePrefix(oldLines, lines)
	suffix := commonLineSuffix(oldLines, lines, prefix)
	keep := min(prefix, len(oldStateSrc))

	d.lines = lines
	d.states = oldStates[: keep+1 : keep+1]
	d.stateSrc = oldStateSrc[:keep:keep]
	d.tail = nil

	oldTailStart := len(oldLines) - suffix
	if suffix > 0 && oldTailStart <= len(oldStateSrc) {
		d.tail = &documentTail{
			states:   oldStates,
			stateSrc: oldStateSrc,
			oldStart: oldTailStart,
			newStart: len(lines) - suffix,
		}
		// A deletion can expose an unchanged suffix at a state that has
		// already converged, so no new tokenization need be performed first.
		d.spliceTailLocked(keep)
	}
}

func (d *Document) spliceTailLocked(newIndex int) {
	tail := d.tail
	if tail == nil || newIndex < tail.newStart {
		return
	}
	oldIndex := tail.oldStart + newIndex - tail.newStart
	if oldIndex < 0 || oldIndex >= len(tail.states) {
		d.tail = nil
		return
	}

	current, previous := d.states[newIndex], tail.states[oldIndex]
	if current != previous && !current.Equal(previous) {
		return
	}

	d.states = append(d.states, tail.states[oldIndex+1:]...)
	d.stateSrc = append(d.stateSrc, tail.stateSrc[oldIndex:]...)
	d.tail = nil
}

func equalDocumentLines(left, right []string) bool {
	if len(left) != len(right) {
		return false
	}
	for index := range left {
		if left[index] != right[index] {
			return false
		}
	}
	return true
}

func commonLinePrefix(left, right []string) int {
	limit := min(len(left), len(right))
	index := 0
	for index < limit && left[index] == right[index] {
		index++
	}
	return index
}

func commonLineSuffix(left, right []string, prefix int) int {
	limit := min(len(left), len(right)) - prefix
	count := 0
	for count < limit && left[len(left)-1-count] == right[len(right)-1-count] {
		count++
	}
	return count
}

func documentStatesEqual(left, right *StateStack) bool {
	return left == right || left.Equal(right)
}
