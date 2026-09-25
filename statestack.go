package textmate

import (
	"strconv"
	"strings"
	"weak"
)

type scopeStack struct {
	parent    *scopeStack
	scopeName string
}

func pushScopeStack(path *scopeStack, scopeNames []string) *scopeStack {
	for _, name := range scopeNames {
		path = &scopeStack{parent: path, scopeName: name}
	}
	return path
}

func (s *scopeStack) push(scopeName string) *scopeStack {
	return &scopeStack{parent: s, scopeName: scopeName}
}

func (s *scopeStack) segments() []string {
	length := 0
	for item := s; item != nil; item = item.parent {
		length++
	}

	result := make([]string, length)
	for item, i := s, length-1; item != nil; item, i = item.parent, i-1 {
		result[i] = item.scopeName
	}
	return result
}

func (s *scopeStack) extends(other *scopeStack) bool {
	for item := s; item != nil; item = item.parent {
		if item == other {
			return true
		}
	}
	return other == nil
}

func (s *scopeStack) extensionFrom(base *scopeStack) ([]string, bool) {
	var reversed []string
	item := s
	for item != nil && item != base {
		reversed = append(reversed, item.scopeName)
		item = item.parent
	}
	if item != base {
		return nil, false
	}

	result := make([]string, len(reversed))
	for i := range reversed {
		result[len(reversed)-1-i] = reversed[i]
	}
	return result, true
}

type attributedScopeStack struct {
	parent          *attributedScopeStack
	scopePath       *scopeStack
	tokenAttributes uint32
	equalityHash    uint64
}

type attributedScopeStackFrame struct {
	tokenAttributes uint32
	scopeNames      []string
}

func newAttributedScopeStack(
	parent *attributedScopeStack,
	scopePath *scopeStack,
	tokenAttributes uint32,
) *attributedScopeStack {
	if scopePath == nil {
		return nil
	}

	parentHash := hashOffset
	if parent != nil {
		parentHash = parent.equalityHash
	}
	return &attributedScopeStack{
		parent:          parent,
		scopePath:       scopePath,
		tokenAttributes: tokenAttributes,
		equalityHash:    hashString(hashUint64(parentHash, uint64(tokenAttributes)), scopePath.scopeName),
	}
}

func newAttributedScopeRoot(scopeName string, tokenAttributes uint32) *attributedScopeStack {
	return newAttributedScopeStack(nil, &scopeStack{scopeName: scopeName}, tokenAttributes)
}

func attributedScopeStackFromExtension(
	base *attributedScopeStack,
	extension []attributedScopeStackFrame,
) *attributedScopeStack {
	current := base
	var scopePath *scopeStack
	if base != nil {
		scopePath = base.scopePath
	}
	for _, frame := range extension {
		scopePath = pushScopeStack(scopePath, frame.scopeNames)
		if scopePath == nil {
			continue
		}
		current = newAttributedScopeStack(current, scopePath, frame.tokenAttributes)
	}
	return current
}

func (s *attributedScopeStack) pushAttributed(scopePath string, tokenAttributes uint32) *attributedScopeStack {
	if s == nil || scopePath == "" {
		return s
	}

	result := s
	for _, scopeName := range strings.Split(scopePath, " ") {
		result = newAttributedScopeStack(result, result.scopePath.push(scopeName), tokenAttributes)
	}
	return result
}

func (s *attributedScopeStack) scopeNames() []string {
	if s == nil {
		return nil
	}
	return s.scopePath.segments()
}

func (s *attributedScopeStack) extensionFrom(base *attributedScopeStack) ([]attributedScopeStackFrame, bool) {
	var reversed []attributedScopeStackFrame
	item := s
	for item != nil && item != base {
		var basePath *scopeStack
		if item.parent != nil {
			basePath = item.parent.scopePath
		}
		scopeNames, ok := item.scopePath.extensionFrom(basePath)
		if !ok {
			return nil, false
		}
		reversed = append(reversed, attributedScopeStackFrame{
			tokenAttributes: item.tokenAttributes,
			scopeNames:      scopeNames,
		})
		item = item.parent
	}
	if item != base {
		return nil, false
	}

	result := make([]attributedScopeStackFrame, len(reversed))
	for i := range reversed {
		result[len(reversed)-1-i] = reversed[i]
	}
	return result, true
}

func attributedScopeStacksEqual(a, b *attributedScopeStack) bool {
	if a == b {
		return true
	}
	if a == nil || b == nil || a.equalityHash != b.equalityHash {
		return false
	}
	for a != nil && b != nil {
		if a == b {
			return true
		}
		if a.scopePath.scopeName != b.scopePath.scopeName || a.tokenAttributes != b.tokenAttributes {
			return false
		}
		a = a.parent
		b = b.parent
	}
	return a == b
}

type stateStackFrame struct {
	ruleID                ruleID
	enterPos              int
	anchorPos             int
	beginRuleCapturedEOL  bool
	endRule               *string
	nameScopesList        []attributedScopeStackFrame
	contentNameScopesList []attributedScopeStackFrame
}

// StateStack is an immutable tokenizer state. A returned state can safely be
// retained and compared after subsequent lines have been tokenized.
//
// While two equal states returned by the same Grammar are both live, they are
// represented by the same *StateStack. This makes pointers suitable as
// short-lived cache keys. The identity is scoped to one Grammar and lasts only
// as long as the state remains reachable; use Equal after discarding all
// previously returned equal states. States are not transferable between
// grammars, and no pointer identity is promised across Grammar instances.
type StateStack struct {
	parent                *StateStack
	ruleID                ruleID
	enterPos              int
	anchorPos             int
	beginRuleCapturedEOL  bool
	endRule               *string
	nameScopesList        *attributedScopeStack
	contentNameScopesList *attributedScopeStack
	depth                 int
	structuralHash        uint64
	equalityHash          uint64
}

// stateStackInterner canonicalizes returned reusable states without making
// their lifetimes equal to the Grammar's lifetime. Grammar.mu serializes all
// access in production; the type itself intentionally has no second lock.
type stateStackInterner struct {
	byHash map[uint64][]weak.Pointer[StateStack]
	calls  uint64
}

const stateStackInternerSweepInterval = 256

func (i *stateStackInterner) intern(stack *StateStack) *StateStack {
	if stack == nil || stack == InitialState {
		return stack
	}
	if i.byHash == nil {
		i.byHash = make(map[uint64][]weak.Pointer[StateStack])
	}

	hash := stack.equalityHash
	bucket := i.byHash[hash]
	live := bucket[:0]
	var canonical *StateStack
	for _, reference := range bucket {
		candidate := reference.Value()
		if candidate == nil {
			continue
		}
		live = append(live, reference)
		if canonical == nil && candidate.Equal(stack) {
			canonical = candidate
		}
	}
	if canonical != nil {
		i.storeBucket(hash, live, bucket)
		i.afterLookup()
		return canonical
	}

	live = append(live, weak.Make(stack))
	i.storeBucket(hash, live, bucket)
	i.afterLookup()
	return stack
}

func (i *stateStackInterner) storeBucket(
	hash uint64,
	live []weak.Pointer[StateStack],
	original []weak.Pointer[StateStack],
) {
	if len(live) == 0 {
		delete(i.byHash, hash)
		return
	}
	// Clear discarded tail entries so a future append cannot expose stale weak
	// references from the reused backing array.
	if len(live) < len(original) {
		clear(original[len(live):])
	}
	i.byHash[hash] = live
}

func (i *stateStackInterner) afterLookup() {
	i.calls++
	if i.calls%stateStackInternerSweepInterval != 0 {
		return
	}
	for hash, bucket := range i.byHash {
		live := bucket[:0]
		for _, reference := range bucket {
			if reference.Value() != nil {
				live = append(live, reference)
			}
		}
		i.storeBucket(hash, live, bucket)
	}
}

// InitialState is the sentinel passed when tokenizing the first line.
var InitialState = newStateStack(nil, 0, 0, 0, false, nil, nil, nil)

func newStateStack(
	parent *StateStack,
	ruleID ruleID,
	enterPos int,
	anchorPos int,
	beginRuleCapturedEOL bool,
	endRule *string,
	nameScopesList *attributedScopeStack,
	contentNameScopesList *attributedScopeStack,
) *StateStack {
	depth := 1
	parentHash := hashOffset
	if parent != nil {
		depth = parent.depth + 1
		parentHash = parent.structuralHash
	}

	structuralHash := hashUint64(parentHash, uint64(ruleID))
	if endRule == nil {
		structuralHash = hashUint64(structuralHash, 0)
	} else {
		structuralHash = hashString(hashUint64(structuralHash, 1), *endRule)
	}
	var scopesHash uint64
	if contentNameScopesList != nil {
		scopesHash = contentNameScopesList.equalityHash
	}

	return &StateStack{
		parent:                parent,
		ruleID:                ruleID,
		enterPos:              enterPos,
		anchorPos:             anchorPos,
		beginRuleCapturedEOL:  beginRuleCapturedEOL,
		endRule:               cloneString(endRule),
		nameScopesList:        nameScopesList,
		contentNameScopesList: contentNameScopesList,
		depth:                 depth,
		structuralHash:        structuralHash,
		equalityHash:          hashUint64(structuralHash, scopesHash),
	}
}

// Equal reports whether two stacks represent the same reusable tokenizer
// state. Line-local enter and anchor positions are intentionally ignored, as
// they are by vscode-textmate.
func (s *StateStack) Equal(other *StateStack) bool {
	if s == other {
		return true
	}
	if s == nil || other == nil || s.equalityHash != other.equalityHash {
		return false
	}
	if !stateStacksStructurallyEqual(s, other) {
		return false
	}
	return attributedScopeStacksEqual(s.contentNameScopesList, other.contentNameScopesList)
}

func stateStacksStructurallyEqual(a, b *StateStack) bool {
	for a != nil && b != nil {
		if a == b {
			return true
		}
		if a.depth != b.depth || a.ruleID != b.ruleID || !equalStrings(a.endRule, b.endRule) {
			return false
		}
		a = a.parent
		b = b.parent
	}
	return a == b
}

func (s *StateStack) clone() *StateStack {
	return s
}

// reset clears transient line positions without mutating this stack or any of
// its ancestors.
func (s *StateStack) reset() *StateStack {
	if s == nil {
		return nil
	}
	parent := s.parent.reset()
	if parent == s.parent && s.enterPos == -1 && s.anchorPos == -1 {
		return s
	}
	return newStateStack(
		parent,
		s.ruleID,
		-1,
		-1,
		s.beginRuleCapturedEOL,
		s.endRule,
		s.nameScopesList,
		s.contentNameScopesList,
	)
}

func (s *StateStack) pop() *StateStack {
	if s == nil {
		return nil
	}
	return s.parent
}

func (s *StateStack) safePop() *StateStack {
	if s == nil || s.parent == nil {
		return s
	}
	return s.parent
}

func (s *StateStack) push(
	ruleID ruleID,
	enterPos int,
	anchorPos int,
	beginRuleCapturedEOL bool,
	endRule *string,
	nameScopesList *attributedScopeStack,
	contentNameScopesList *attributedScopeStack,
) *StateStack {
	return newStateStack(
		s,
		ruleID,
		enterPos,
		anchorPos,
		beginRuleCapturedEOL,
		endRule,
		nameScopesList,
		contentNameScopesList,
	)
}

func (s *StateStack) getEnterPos() int {
	return s.enterPos
}

func (s *StateStack) getAnchorPos() int {
	return s.anchorPos
}

func (s *StateStack) getRuleID() ruleID {
	return s.ruleID
}

func (s *StateStack) withContentNameScopesList(contentNameScopesList *attributedScopeStack) *StateStack {
	if s == nil || s.contentNameScopesList == contentNameScopesList {
		return s
	}
	return newStateStack(
		s.parent,
		s.ruleID,
		s.enterPos,
		s.anchorPos,
		s.beginRuleCapturedEOL,
		s.endRule,
		s.nameScopesList,
		contentNameScopesList,
	)
}

func (s *StateStack) withEndRule(endRule string) *StateStack {
	if s == nil || (s.endRule != nil && *s.endRule == endRule) {
		return s
	}
	return newStateStack(
		s.parent,
		s.ruleID,
		s.enterPos,
		s.anchorPos,
		s.beginRuleCapturedEOL,
		&endRule,
		s.nameScopesList,
		s.contentNameScopesList,
	)
}

func (s *StateStack) hasSameRuleAs(other *StateStack) bool {
	if s == nil || other == nil {
		return false
	}
	for item := s; item != nil && item.enterPos == other.enterPos; item = item.parent {
		if item.ruleID == other.ruleID {
			return true
		}
	}
	return false
}

func (s *StateStack) toStateStackFrame() stateStackFrame {
	var parentNames *attributedScopeStack
	if s.parent != nil {
		parentNames = s.parent.nameScopesList
	}
	nameScopes, _ := s.nameScopesList.extensionFrom(parentNames)
	contentScopes, _ := s.contentNameScopesList.extensionFrom(s.nameScopesList)
	return stateStackFrame{
		ruleID:                s.ruleID,
		enterPos:              s.enterPos,
		anchorPos:             s.anchorPos,
		beginRuleCapturedEOL:  s.beginRuleCapturedEOL,
		endRule:               cloneString(s.endRule),
		nameScopesList:        nameScopes,
		contentNameScopesList: contentScopes,
	}
}

func pushStateStackFrame(stack *StateStack, frame stateStackFrame) *StateStack {
	var parentNames *attributedScopeStack
	if stack != nil {
		parentNames = stack.nameScopesList
	}
	nameScopes := attributedScopeStackFromExtension(parentNames, frame.nameScopesList)
	return newStateStack(
		stack,
		frame.ruleID,
		frame.enterPos,
		frame.anchorPos,
		frame.beginRuleCapturedEOL,
		frame.endRule,
		nameScopes,
		attributedScopeStackFromExtension(nameScopes, frame.contentNameScopesList),
	)
}

func (s *StateStack) String() string {
	if s == nil {
		return "[]"
	}
	frames := make([]string, s.depth)
	for item, i := s, s.depth-1; item != nil; item, i = item.parent, i-1 {
		frames[i] = "(" + strconv.Itoa(int(item.ruleID)) + ", " +
			strings.Join(item.nameScopesList.scopeNames(), " ") + ", " +
			strings.Join(item.contentNameScopesList.scopeNames(), " ") + ")"
	}
	return "[" + strings.Join(frames, ",") + "]"
}

func cloneString(value *string) *string {
	if value == nil {
		return nil
	}
	result := *value
	return &result
}

func equalStrings(a, b *string) bool {
	if a == nil || b == nil {
		return a == b
	}
	return *a == *b
}

const hashOffset uint64 = 14695981039346656037

func hashUint64(hash, value uint64) uint64 {
	const prime uint64 = 1099511628211
	for range 8 {
		hash ^= value & 0xff
		hash *= prime
		value >>= 8
	}
	return hash
}

func hashString(hash uint64, value string) uint64 {
	const prime uint64 = 1099511628211
	for i := 0; i < len(value); i++ {
		hash ^= uint64(value[i])
		hash *= prime
	}
	return hashUint64(hash, uint64(len(value)))
}
