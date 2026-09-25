// Package textmate provides pure-Go, line-oriented tokenization using TextMate
// grammars. Token offsets are rune indices. Callers pass each LineResult's
// immutable RuleStack into the next line and may compare saved stacks with
// StateStack.Equal when incrementally re-tokenizing edited documents.
package textmate
