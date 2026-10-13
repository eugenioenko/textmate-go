package textmate

import (
	"errors"
	"slices"
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

	// partialState is the earliest index i > 0 whose states[i] came from a
	// line stopped by the time limit, or 0. A time stop can be caused by
	// transient lock or CPU pressure, so that state, and every state after it,
	// is recomputed on the next call instead of being kept until an edit.
	partialState int
	// partialRetries bounds those recomputations for a line that keeps
	// exceeding its budget; retryLine is the line last retried.
	partialRetries int
	retryLine      int
}

// maxTimeLimitRetries is how many times a state from a time-stopped line is
// recomputed before it is kept. A line whose own work exceeds the budget would
// otherwise cost a full budget on every call.
const maxTimeLimitRetries = 3

type documentTail struct {
	states   []*StateStack
	stateSrc []string
	oldStart int
	newStart int
	// inPlace marks a tail that ReplaceLines shifted into the arrays the
	// materialized prefix grows into: states and stateSrc already have the new
	// layout, so convergence adopts them instead of copying the tail.
	inPlace bool
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
//
// Lines and states are spliced in place, so an edit costs the replaced lines
// plus shifting the ones after them, rather than copying and comparing the
// whole document as SetLines does.
func (d *Document) ReplaceLines(start, end int, replacement []string) error {
	if d == nil {
		return ErrInvalidLineRange
	}
	d.mu.Lock()
	defer d.mu.Unlock()
	if start < 0 || start > end || end > len(d.lines) {
		return ErrInvalidLineRange
	}
	if len(replacement) == end-start && slices.Equal(d.lines[start:end], replacement) {
		return nil
	}
	if d.tail != nil {
		// An earlier edit has not converged; SetLines semantics settle both.
		lines := make([]string, 0, len(d.lines)-(end-start)+len(replacement))
		lines = append(lines, d.lines[:start]...)
		lines = append(lines, replacement...)
		lines = append(lines, d.lines[end:]...)
		d.setLinesLocked(lines)
		return nil
	}

	d.lines = slices.Replace(d.lines, start, end, replacement...)
	d.partialRetries = 0
	materialized := len(d.stateSrc)
	keep := min(start, materialized)
	reusable := materialized >= end
	if d.partialState > keep {
		// The old suffix may hold states derived from a time-stopped line.
		d.partialState = 0
		reusable = false
	}
	if !reusable {
		d.states = d.states[:keep+1]
		d.stateSrc = d.stateSrc[:keep]
		return nil
	}

	n := len(replacement)
	if n == 0 {
		// The line after a deletion starts where the deleted range started.
		converged := documentStatesEqual(d.states[start], d.states[end])
		d.states = slices.Delete(d.states, start+1, end+1)
		d.stateSrc = slices.Delete(d.stateSrc, start, end)
		if converged || start+1 >= len(d.states) {
			if !converged {
				d.states = d.states[:start+1]
				d.stateSrc = d.stateSrc[:start]
			}
			return nil
		}
	} else {
		// states[start+n] is where the first unchanged line now starts; keep its
		// old state there as the convergence candidate.
		candidate := d.states[end]
		d.states = slices.Replace(d.states, start+1, end+1, make([]*StateStack, n)...)
		d.states[start+n] = candidate
		d.stateSrc = slices.Replace(d.stateSrc, start, end, make([]string, n)...)
	}
	newStart := start + max(n, 1)
	d.tail = &documentTail{
		states:   d.states,
		stateSrc: d.stateSrc,
		oldStart: newStart,
		newStart: newStart,
		inPlace:  true,
	}
	d.states = d.states[:start+1]
	d.stateSrc = d.stateSrc[:start]
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

	d.retryPartialStateLocked()
	state := d.stateAtLocked(index)
	result := d.tokenizeLocked(d.lines[index], state)
	if len(d.states) == index+1 {
		d.appendStateLocked(index, result)
	} else if !documentStatesEqual(d.states[index+1], result.RuleStack) {
		d.replaceMaterializedResultLocked(index, result.RuleStack)
		d.noteStateLocked(index, result)
	}
	return result, true
}

// MaterializedLines returns how many leading lines have a known start state.
// Line(index) for an index below it tokenizes only that line; a larger index
// first tokenizes every line from MaterializedLines up to it. An edit lowers it
// to the first changed line until recomputation converges with the old tail.
func (d *Document) MaterializedLines() int {
	if d == nil {
		return 0
	}
	d.mu.Lock()
	defer d.mu.Unlock()
	d.retryPartialStateLocked()
	return min(len(d.states), len(d.lines))
}

// InvalidateFrom discards materialized states for line index and all following
// lines, and clears the result cache. Index may equal Len. The next Line or
// StateAt call recomputes the discarded range.
//
// Time-limit stops are never cached, and a state that came from one is
// recomputed on the next Line, StateAt or MaterializedLines call, a bounded
// number of times. InvalidateFrom forces the complete tail to be recomputed
// regardless, for example after a grammar or theme change on the caller side.
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
	if d.partialState > keep {
		d.partialState = 0
	}
	d.partialRetries = 0
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
	d.retryPartialStateLocked()
	return d.stateAtLocked(index), true
}

func (d *Document) stateAtLocked(index int) *StateStack {
	for len(d.states) <= index {
		lineIndex := len(d.states) - 1
		result := d.tokenizeLocked(d.lines[lineIndex], d.states[lineIndex])
		d.appendStateLocked(lineIndex, result)
	}
	return d.states[index]
}

// appendStateLocked records the state after lineIndex, which must be the last
// materialized line, and reuses the pending tail once states converge.
func (d *Document) appendStateLocked(lineIndex int, result LineResult) {
	next := lineIndex + 1
	tail := d.tail
	var candidate *StateStack
	compare := tail != nil && tail.inPlace && next >= tail.newStart && next < len(tail.states)
	if compare {
		// The append below overwrites this slot of the shared array.
		candidate = tail.states[next]
	}
	d.states = append(d.states, result.RuleStack)
	d.stateSrc = append(d.stateSrc, d.lines[lineIndex])
	d.noteStateLocked(lineIndex, result)
	switch {
	case tail == nil:
	case !tail.inPlace:
		d.spliceTailLocked(next)
	case compare && documentStatesEqual(candidate, result.RuleStack):
		d.states = tail.states
		d.stateSrc = tail.stateSrc
		d.tail = nil
	case next+1 >= len(tail.states):
		d.tail = nil
	}
}

// noteStateLocked records that states[lineIndex+1] came from result.
func (d *Document) noteStateLocked(lineIndex int, result LineResult) {
	if result.StoppedReason == StopReasonTimeLimit {
		if d.partialState == 0 || lineIndex+1 < d.partialState {
			d.partialState = lineIndex + 1
		}
		return
	}
	if lineIndex == d.retryLine {
		d.partialRetries = 0
	}
}

// retryPartialStateLocked discards the states from the earliest time-stopped
// line on, so they are recomputed lazily. States derived from the partial
// state are not reused through a tail: an equal state found later could still
// lead back into states built on another stopped line.
func (d *Document) retryPartialStateLocked() {
	index := d.partialState
	if index == 0 || d.partialRetries >= maxTimeLimitRetries {
		return
	}
	d.partialState = 0
	d.partialRetries++
	d.retryLine = index - 1
	keep := min(index-1, len(d.stateSrc))
	d.states = d.states[: keep+1 : keep+1]
	d.stateSrc = d.stateSrc[:keep:keep]
	d.tail = nil
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
	d.partialRetries = 0
	reusable := true
	if d.partialState > keep {
		// The old suffix may hold states derived from a time-stopped line.
		d.partialState = 0
		reusable = false
	}

	oldTailStart := len(oldLines) - suffix
	if reusable && suffix > 0 && oldTailStart <= len(oldStateSrc) {
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
