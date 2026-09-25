package textmate

import "time"

// TokenizeOptions bounds work performed by TokenizeLineWithOptions. Zero or
// negative values disable the corresponding limit.
//
// TimeLimit is a soft budget for the complete call, including waiting for the
// grammar lock, lazy root compilation, and begin/while condition setup. The
// tokenizer checks the budget between regexp searches; an in-flight regexp
// search cannot be interrupted and can therefore overrun the budget.
//
// MaxLineBytes and MaxLineRunes are independent. A line is skipped when it
// exceeds either enabled limit. Limits are checked before the newline-appended
// Oniguruma input and its rune buffer are allocated.
type TokenizeOptions struct {
	TimeLimit    time.Duration
	MaxLineBytes int
	MaxLineRunes int
}

// StopReason explains why a LineResult stopped before fully parsing its line.
type StopReason uint8

const (
	// StopReasonNone means the line was fully tokenized.
	StopReasonNone StopReason = iota
	// StopReasonTimeLimit means the soft time budget was exhausted.
	StopReasonTimeLimit
	// StopReasonLineLimit means a byte or rune line-length cap was exceeded.
	StopReasonLineLimit
)

func (r StopReason) String() string {
	switch r {
	case StopReasonNone:
		return "none"
	case StopReasonTimeLimit:
		return "time_limit"
	case StopReasonLineLimit:
		return "line_limit"
	default:
		return "unknown"
	}
}
