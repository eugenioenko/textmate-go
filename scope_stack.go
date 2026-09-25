package textmate

import "sync/atomic"

// ScopeStackID identifies one interned scope-name stack. IDs are comparable
// and unique within the current process, so they can be used directly as map
// keys. The zero value does not identify a scope stack.
//
// An ID remains valid for as long as the ScopeStack that returned it is
// reachable. IDs are runtime identities and must not be persisted.
type ScopeStackID struct {
	owner uint64
	index uint64
}

// IsZero reports whether id is the zero value.
func (id ScopeStackID) IsZero() bool {
	return id == ScopeStackID{}
}

// ScopeStack is an immutable, interned sequence of TextMate scope names.
// Equal scope sequences produced by one Grammar have the same *ScopeStack and
// ScopeStackID, including when the tokenizer reconstructed its internal stack.
// A ScopeStack remains valid after its Grammar or Registry is no longer used.
//
// ScopeStack values are owned by the library and must not be modified. Its
// accessors are safe to call concurrently.
type ScopeStack struct {
	id    ScopeStackID
	names []string

	// compatibilityNames backs Token.Scopes. It is deliberately separate from
	// names so mutation through that legacy field cannot corrupt the interner.
	compatibilityNames []string
}

// ID returns the stable runtime identity of s. A nil ScopeStack has the zero
// ID.
func (s *ScopeStack) ID() ScopeStackID {
	if s == nil {
		return ScopeStackID{}
	}
	return s.id
}

// Len returns the number of scope names in s.
func (s *ScopeStack) Len() int {
	if s == nil {
		return 0
	}
	return len(s.names)
}

// At returns the scope name at index and whether index is in range.
func (s *ScopeStack) At(index int) (string, bool) {
	if s == nil || index < 0 || index >= len(s.names) {
		return "", false
	}
	return s.names[index], true
}

// Names returns a fresh copy of the scope names, ordered outermost first.
func (s *ScopeStack) Names() []string {
	if s == nil {
		return nil
	}
	return append([]string(nil), s.names...)
}

// Range calls yield for each scope name, ordered outermost first. Iteration
// stops when yield returns false.
func (s *ScopeStack) Range(yield func(string) bool) {
	if s == nil || yield == nil {
		return
	}
	for _, name := range s.names {
		if !yield(name) {
			return
		}
	}
}

type scopeStackInterner struct {
	owner  uint64
	nextID uint64
	byHash map[uint64][]*ScopeStack

	// recent is a small bounded identity cache for internal paths that survive
	// across adjacent tokens or lines (especially the root path). It avoids
	// hashing and comparing those common hits without retaining every transient
	// path created while matching captures.
	recent      [8]scopeStackCacheEntry
	recentIndex uint8
}

type scopeStackCacheEntry struct {
	path  *scopeStack
	stack *ScopeStack
}

var nextScopeStackOwner atomic.Uint64

func (i *scopeStackInterner) intern(scopes *attributedScopeStack) *ScopeStack {
	var path *scopeStack
	if scopes != nil {
		path = scopes.scopePath
	}
	for index := range i.recent {
		entry := i.recent[index]
		if entry.stack != nil && entry.path == path {
			return entry.stack
		}
	}
	stack := i.internPathWithHash(path, hashScopePath(path))
	i.recent[i.recentIndex%uint8(len(i.recent))] = scopeStackCacheEntry{path: path, stack: stack}
	i.recentIndex++
	return stack
}

// internPathWithHash takes an explicit hash so collision handling can be
// tested without depending on a naturally occurring 64-bit collision.
func (i *scopeStackInterner) internPathWithHash(path *scopeStack, hash uint64) *ScopeStack {
	if i.owner == 0 {
		i.owner = nextScopeStackOwner.Add(1)
	}
	if i.byHash == nil {
		i.byHash = make(map[uint64][]*ScopeStack)
	}
	for _, candidate := range i.byHash[hash] {
		if scopePathEqualsNames(path, candidate.names) {
			return candidate
		}
	}

	names := scopePathNames(path)
	i.nextID++
	stack := &ScopeStack{
		id:                 ScopeStackID{owner: i.owner, index: i.nextID},
		names:              names,
		compatibilityNames: append([]string(nil), names...),
	}
	i.byHash[hash] = append(i.byHash[hash], stack)
	return stack
}

func hashScopePath(path *scopeStack) uint64 {
	hash := hashOffset
	for item := path; item != nil; item = item.parent {
		// Walking innermost-first is sufficient for a bucket key and avoids
		// materializing the outermost-first name slice on an interner hit.
		hash = hashString(hash, item.scopeName)
	}
	return hash
}

func scopePathEqualsNames(path *scopeStack, names []string) bool {
	index := len(names) - 1
	for item := path; item != nil; item = item.parent {
		if index < 0 || names[index] != item.scopeName {
			return false
		}
		index--
	}
	return index == -1
}

func scopePathNames(path *scopeStack) []string {
	length := 0
	for item := path; item != nil; item = item.parent {
		length++
	}
	names := make([]string, length)
	for item, index := path, length-1; item != nil; item, index = item.parent, index-1 {
		names[index] = item.scopeName
	}
	return names
}
