package textmate

import (
	"strings"
	"time"

	"github.com/eugenioenko/textmate-go/oniguruma"
)

type injection struct {
	debugSelector string
	matcher       scopeMatcher
	priority      int
	ruleID        ruleID
}

// tokenizerGrammar is the deliberately small portion of Grammar needed by
// the hot tokenization loop. Keeping this boundary narrow makes the loop
// independently testable and keeps rule compilation lazy.
type tokenizerGrammar interface {
	ruleRegistry
	getInjections() []injection
	scopeAttributes(scopeName string) uint32
}

// tokenHandler consumes the half-open token range ending at endIndex. Scope
// stacks are immutable, and all indexes are rune offsets.
type tokenHandler interface {
	handle(scopes *attributedScopeStack, endIndex int)
}

type tokenizeStringResult struct {
	stack        *StateStack
	stoppedEarly bool
	stoppedAt    int
}

type tokenizationBudget struct {
	deadline time.Time
	limited  bool
	stop     func() bool
}

func newTokenizationBudget(limit time.Duration) tokenizationBudget {
	if limit <= 0 {
		return tokenizationBudget{}
	}
	return tokenizationBudget{deadline: time.Now().Add(limit), limited: true}
}

func (b tokenizationBudget) exceeded() bool {
	if b.stop != nil {
		return b.stop()
	}
	return b.limited && !time.Now().Before(b.deadline)
}

// tokenizeString is the direct Go counterpart of vscode-textmate's
// _tokenizeString. lineText is expected to include the synthetic trailing
// newline added by Grammar.TokenizeLine.
func tokenizeString(
	grammar tokenizerGrammar,
	lineText *oniguruma.String,
	isFirstLine bool,
	linePos int,
	stack *StateStack,
	handler tokenHandler,
	checkWhileConditions bool,
	timeLimit time.Duration,
) tokenizeStringResult {
	return tokenizeStringWithBudget(
		grammar,
		lineText,
		isFirstLine,
		linePos,
		stack,
		handler,
		checkWhileConditions,
		newTokenizationBudget(timeLimit),
	)
}

func tokenizeStringWithBudget(
	grammar tokenizerGrammar,
	lineText *oniguruma.String,
	isFirstLine bool,
	linePos int,
	stack *StateStack,
	handler tokenHandler,
	checkWhileConditions bool,
	budget tokenizationBudget,
) tokenizeStringResult {
	lineLength := lineText.Len()
	anchorPosition := -1

	if budget.exceeded() {
		return tokenizeStringResult{stack: stack, stoppedEarly: true, stoppedAt: linePos}
	}

	if checkWhileConditions {
		checked := checkWhileWithBudget(grammar, lineText, isFirstLine, linePos, stack, handler, budget)
		stack = checked.stack
		linePos = checked.linePos
		isFirstLine = checked.isFirstLine
		anchorPosition = checked.anchorPosition
		if checked.stoppedEarly {
			return tokenizeStringResult{stack: stack, stoppedEarly: true, stoppedAt: linePos}
		}
	}

	for {
		if budget.exceeded() {
			return tokenizeStringResult{stack: stack, stoppedEarly: true, stoppedAt: linePos}
		}

		matched := matchRuleOrInjections(grammar, lineText, isFirstLine, linePos, stack, anchorPosition)
		if budget.exceeded() {
			return tokenizeStringResult{stack: stack, stoppedEarly: true, stoppedAt: linePos}
		}
		if matched == nil || len(matched.captureIndices) == 0 {
			handler.handle(stack.contentNameScopesList, lineLength)
			return tokenizeStringResult{stack: stack}
		}

		captures := matched.captureIndices
		wholeMatch := captures[0]
		hasAdvanced := wholeMatch.End > linePos

		if matched.matchedRuleID == endRuleID {
			poppedRule, ok := grammar.getRule(stack.ruleID).(*beginEndRule)
			if !ok {
				// A compiled end sentinel can only belong to a begin/end rule.
				// Treat an inconsistent registry as no further match instead of
				// panicking in an editor process.
				handler.handle(stack.contentNameScopesList, lineLength)
				return tokenizeStringResult{stack: stack}
			}

			handler.handle(stack.contentNameScopesList, wholeMatch.Start)
			stack = stack.withContentNameScopesList(stack.nameScopesList)
			handleCapturesWithBudget(grammar, lineText, isFirstLine, stack, handler, poppedRule.endCaptures, captures, budget)
			handler.handle(stack.contentNameScopesList, wholeMatch.End)

			popped := stack
			stack = stack.parent
			anchorPosition = popped.getAnchorPos()

			if !hasAdvanced && popped.getEnterPos() == linePos {
				// Guard 1: a rule was pushed and popped without consuming text.
				stack = popped
				handler.handle(stack.contentNameScopesList, lineLength)
				return tokenizeStringResult{stack: stack}
			}
		} else {
			matchedRule := grammar.getRule(matched.matchedRuleID)
			if matchedRule == nil {
				handler.handle(stack.contentNameScopesList, lineLength)
				return tokenizeStringResult{stack: stack}
			}

			handler.handle(stack.contentNameScopesList, wholeMatch.Start)

			beforePush := stack
			nameScopes := pushAttributedScope(
				stack.contentNameScopesList,
				matchedRule.getName(lineText.Content(), captures),
				grammar,
			)
			stack = stack.push(
				matched.matchedRuleID,
				linePos,
				anchorPosition,
				wholeMatch.End == lineLength,
				nil,
				nameScopes,
				nameScopes,
			)

			switch typedRule := matchedRule.(type) {
			case *beginEndRule:
				handleCapturesWithBudget(grammar, lineText, isFirstLine, stack, handler, typedRule.beginCaptures, captures, budget)
				handler.handle(stack.contentNameScopesList, wholeMatch.End)
				anchorPosition = wholeMatch.End

				contentScopes := pushAttributedScope(
					nameScopes,
					typedRule.getContentName(lineText.Content(), captures),
					grammar,
				)
				stack = stack.withContentNameScopesList(contentScopes)
				if typedRule.endHasBackReferences {
					stack = stack.withEndRule(typedRule.getEndWithResolvedBackReferences(lineText.Content(), captures))
				}

				if !hasAdvanced && beforePush.hasSameRuleAs(stack) {
					// Guard 2: the same begin/end rule was pushed without
					// consuming text.
					stack = stack.pop()
					handler.handle(stack.contentNameScopesList, lineLength)
					return tokenizeStringResult{stack: stack}
				}

			case *beginWhileRule:
				handleCapturesWithBudget(grammar, lineText, isFirstLine, stack, handler, typedRule.beginCaptures, captures, budget)
				handler.handle(stack.contentNameScopesList, wholeMatch.End)
				anchorPosition = wholeMatch.End

				contentScopes := pushAttributedScope(
					nameScopes,
					typedRule.getContentName(lineText.Content(), captures),
					grammar,
				)
				stack = stack.withContentNameScopesList(contentScopes)
				if typedRule.whileHasBackReferences {
					stack = stack.withEndRule(typedRule.getWhileWithResolvedBackReferences(lineText.Content(), captures))
				}

				if !hasAdvanced && beforePush.hasSameRuleAs(stack) {
					// Guard 3: the same begin/while rule was pushed without
					// consuming text.
					stack = stack.pop()
					handler.handle(stack.contentNameScopesList, lineLength)
					return tokenizeStringResult{stack: stack}
				}

			case *matchRule:
				handleCapturesWithBudget(grammar, lineText, isFirstLine, stack, handler, typedRule.captures, captures, budget)
				handler.handle(stack.contentNameScopesList, wholeMatch.End)
				stack = stack.pop()

				if !hasAdvanced {
					// Guard 4: a match rule neither consumed input nor left a
					// new rule on the stack.
					stack = stack.safePop()
					handler.handle(stack.contentNameScopesList, lineLength)
					return tokenizeStringResult{stack: stack}
				}

			default:
				// Include-only rules are flattened during compilation and
				// therefore cannot be returned as direct matches.
				stack = stack.pop()
				handler.handle(stack.contentNameScopesList, lineLength)
				return tokenizeStringResult{stack: stack}
			}
		}

		if wholeMatch.End > linePos {
			linePos = wholeMatch.End
			isFirstLine = false
		}
	}
}

type whileCheckResult struct {
	stack          *StateStack
	linePos        int
	anchorPosition int
	isFirstLine    bool
	stoppedEarly   bool
}

// checkWhile walks active begin/while rules from the bottom of the stack to
// the top. A failed condition removes that rule and every rule above it.
func checkWhile(
	grammar tokenizerGrammar,
	lineText *oniguruma.String,
	isFirstLine bool,
	linePos int,
	stack *StateStack,
	handler tokenHandler,
) whileCheckResult {
	return checkWhileWithBudget(
		grammar,
		lineText,
		isFirstLine,
		linePos,
		stack,
		handler,
		tokenizationBudget{},
	)
}

func checkWhileWithBudget(
	grammar tokenizerGrammar,
	lineText *oniguruma.String,
	isFirstLine bool,
	linePos int,
	stack *StateStack,
	handler tokenHandler,
	budget tokenizationBudget,
) whileCheckResult {
	anchorPosition := -1
	if stack != nil && stack.beginRuleCapturedEOL {
		anchorPosition = 0
	}

	type whileStack struct {
		stack *StateStack
		rule  *beginWhileRule
	}
	var whileRules []whileStack
	for node := stack; node != nil; node = node.pop() {
		if candidate, ok := grammar.getRule(node.ruleID).(*beginWhileRule); ok {
			whileRules = append(whileRules, whileStack{stack: node, rule: candidate})
		}
	}

	for index := len(whileRules) - 1; index >= 0; index-- {
		if budget.exceeded() {
			return whileCheckResult{
				stack: stack, linePos: linePos, anchorPosition: anchorPosition,
				isFirstLine: isFirstLine, stoppedEarly: true,
			}
		}
		active := whileRules[index]
		scanner := active.rule.compileWhileAG(
			grammar,
			active.stack.endRule,
			isFirstLine,
			linePos == anchorPosition,
		)
		matched := scanner.findNextMatch(lineText, linePos, oniguruma.FindOptionNone)
		if budget.exceeded() {
			return whileCheckResult{
				stack: stack, linePos: linePos, anchorPosition: anchorPosition,
				isFirstLine: isFirstLine, stoppedEarly: true,
			}
		}
		if matched == nil || matched.ruleID != whileRuleID {
			stack = active.stack.pop()
			break
		}

		if len(matched.captureIndices) == 0 {
			continue
		}
		wholeMatch := matched.captureIndices[0]
		handler.handle(active.stack.contentNameScopesList, wholeMatch.Start)
		handleCapturesWithBudget(
			grammar,
			lineText,
			isFirstLine,
			active.stack,
			handler,
			active.rule.whileCaptures,
			matched.captureIndices,
			budget,
		)
		handler.handle(active.stack.contentNameScopesList, wholeMatch.End)
		anchorPosition = wholeMatch.End
		if wholeMatch.End > linePos {
			linePos = wholeMatch.End
			isFirstLine = false
		}
	}

	return whileCheckResult{
		stack:          stack,
		linePos:        linePos,
		anchorPosition: anchorPosition,
		isFirstLine:    isFirstLine,
	}
}

type matchResult struct {
	captureIndices []oniguruma.Capture
	matchedRuleID  ruleID
	priorityMatch  bool
}

func matchRuleOrInjections(
	grammar tokenizerGrammar,
	lineText *oniguruma.String,
	isFirstLine bool,
	linePos int,
	stack *StateStack,
	anchorPosition int,
) *matchResult {
	regular := matchRuleAt(grammar, lineText, isFirstLine, linePos, stack, anchorPosition)
	injections := grammar.getInjections()
	if len(injections) == 0 {
		return regular
	}

	injected := matchInjections(injections, grammar, lineText, isFirstLine, linePos, stack, anchorPosition)
	if injected == nil {
		return regular
	}
	if regular == nil {
		return injected
	}

	regularStart := regular.captureIndices[0].Start
	injectedStart := injected.captureIndices[0].Start
	if injectedStart < regularStart || (injected.priorityMatch && injectedStart == regularStart) {
		return injected
	}
	return regular
}

func matchRuleAt(
	grammar tokenizerGrammar,
	lineText *oniguruma.String,
	isFirstLine bool,
	linePos int,
	stack *StateStack,
	anchorPosition int,
) *matchResult {
	if stack == nil {
		return nil
	}
	activeRule := grammar.getRule(stack.ruleID)
	if activeRule == nil {
		return nil
	}
	scanner := activeRule.compileAG(grammar, stack.endRule, isFirstLine, linePos == anchorPosition)
	matched := scanner.findNextMatch(lineText, linePos, oniguruma.FindOptionNone)
	if matched == nil || len(matched.captureIndices) == 0 {
		return nil
	}
	return &matchResult{captureIndices: matched.captureIndices, matchedRuleID: matched.ruleID}
}

func matchInjections(
	injections []injection,
	grammar tokenizerGrammar,
	lineText *oniguruma.String,
	isFirstLine bool,
	linePos int,
	stack *StateStack,
	anchorPosition int,
) *matchResult {
	if stack == nil || stack.contentNameScopesList == nil {
		return nil
	}

	bestStart := int(^uint(0) >> 1)
	var best *matchResult
	scopes := stack.contentNameScopesList.scopeNames()
	for _, candidate := range injections {
		if candidate.matcher == nil || !candidate.matcher(scopes) {
			continue
		}
		candidateRule := grammar.getRule(candidate.ruleID)
		if candidateRule == nil {
			continue
		}
		scanner := candidateRule.compileAG(grammar, nil, isFirstLine, linePos == anchorPosition)
		matched := scanner.findNextMatch(lineText, linePos, oniguruma.FindOptionNone)
		if matched == nil || len(matched.captureIndices) == 0 {
			continue
		}

		matchStart := matched.captureIndices[0].Start
		if matchStart >= bestStart {
			// Injections are already sorted by priority, so an equal match
			// keeps the earlier candidate.
			continue
		}
		bestStart = matchStart
		best = &matchResult{
			captureIndices: matched.captureIndices,
			matchedRuleID:  matched.ruleID,
			priorityMatch:  candidate.priority == -1,
		}
		if bestStart == linePos {
			break
		}
	}
	return best
}

func handleCaptures(
	grammar tokenizerGrammar,
	lineText *oniguruma.String,
	isFirstLine bool,
	stack *StateStack,
	handler tokenHandler,
	captureRules []*captureRule,
	captureIndices []oniguruma.Capture,
) {
	handleCapturesWithBudget(
		grammar,
		lineText,
		isFirstLine,
		stack,
		handler,
		captureRules,
		captureIndices,
		tokenizationBudget{},
	)
}

func handleCapturesWithBudget(
	grammar tokenizerGrammar,
	lineText *oniguruma.String,
	isFirstLine bool,
	stack *StateStack,
	handler tokenHandler,
	captureRules []*captureRule,
	captureIndices []oniguruma.Capture,
	budget tokenizationBudget,
) {
	if len(captureRules) == 0 || len(captureIndices) == 0 {
		return
	}

	type localStackElement struct {
		scopes *attributedScopeStack
		endPos int
	}
	localStack := make([]localStackElement, 0, len(captureRules))
	maxEnd := captureIndices[0].End
	length := min(len(captureRules), len(captureIndices))

	for index := 0; index < length; index++ {
		captureRule := captureRules[index]
		if captureRule == nil {
			continue
		}
		capture := captureIndices[index]
		if capture.End <= capture.Start {
			continue
		}
		if capture.Start > maxEnd {
			break
		}

		for len(localStack) > 0 && localStack[len(localStack)-1].endPos <= capture.Start {
			last := localStack[len(localStack)-1]
			handler.handle(last.scopes, last.endPos)
			localStack = localStack[:len(localStack)-1]
		}

		if len(localStack) > 0 {
			handler.handle(localStack[len(localStack)-1].scopes, capture.Start)
		} else {
			handler.handle(stack.contentNameScopesList, capture.Start)
		}

		if captureRule.retokenizeCapturedWithRuleID != 0 {
			nameScopes := pushAttributedScope(
				stack.contentNameScopesList,
				captureRule.getName(lineText.Content(), captureIndices),
				grammar,
			)
			contentScopes := pushAttributedScope(
				nameScopes,
				captureRule.getContentName(lineText.Content(), captureIndices),
				grammar,
			)
			captureStack := stack.push(
				captureRule.retokenizeCapturedWithRuleID,
				capture.Start,
				-1,
				false,
				nil,
				nameScopes,
				contentScopes,
			)

			lineRunes := []rune(lineText.Content())
			end := capture.End
			if end > len(lineRunes) {
				end = len(lineRunes)
			}
			if end >= 0 {
				captureText := oniguruma.NewString(string(lineRunes[:end]))
				tokenizeStringWithBudget(
					grammar,
					captureText,
					isFirstLine && capture.Start == 0,
					capture.Start,
					captureStack,
					handler,
					false,
					budget,
				)
			}
			continue
		}

		captureScopeName := captureRule.getName(lineText.Content(), captureIndices)
		if captureScopeName != "" {
			base := stack.contentNameScopesList
			if len(localStack) > 0 {
				base = localStack[len(localStack)-1].scopes
			}
			localStack = append(localStack, localStackElement{
				scopes: pushAttributedScope(base, captureScopeName, grammar),
				endPos: capture.End,
			})
		}
	}

	for len(localStack) > 0 {
		last := localStack[len(localStack)-1]
		handler.handle(last.scopes, last.endPos)
		localStack = localStack[:len(localStack)-1]
	}
}

func pushAttributedScope(
	stack *attributedScopeStack,
	scopePath string,
	grammar tokenizerGrammar,
) *attributedScopeStack {
	for _, scopeName := range strings.Fields(scopePath) {
		stack = stack.pushAttributed(scopeName, grammar.scopeAttributes(scopeName))
	}
	return stack
}
