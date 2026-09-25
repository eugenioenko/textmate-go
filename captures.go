package textmate

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/eugenioenko/textmate-go/oniguruma"
)

// captureReferencePattern is deliberately narrower than a general TextMate
// variable parser. It mirrors vscode-textmate's RegexSource substitutions:
// plain numeric references and the two case-conversion commands.
var captureReferencePattern = regexp.MustCompile(`\$(\d+)|\$\{(\d+):/(downcase|upcase)\}`)

// hasCaptureReferences reports whether a rule name needs match-time capture
// substitution.
func hasCaptureReferences(value string) bool {
	return captureReferencePattern.MatchString(value)
}

// replaceCaptureReferences expands capture references in a scope name. The
// line and capture offsets are runes, matching the public Go tokenizer API.
// References to missing capture slots remain unchanged, as in
// vscode-textmate. A present but unmatched or otherwise invalid slot expands
// to the empty string.
func replaceCaptureReferences(value string, line []rune, captures []oniguruma.Capture) string {
	return captureReferencePattern.ReplaceAllStringFunc(value, func(reference string) string {
		parts := captureReferencePattern.FindStringSubmatch(reference)
		indexText := parts[1]
		if indexText == "" {
			indexText = parts[2]
		}

		index, err := strconv.Atoi(indexText)
		if err != nil || index < 0 || index >= len(captures) {
			return reference
		}

		result := string(captureRunes(line, captures[index]))
		// Leading dots would make the generated scope selector invalid.
		result = strings.TrimLeft(result, ".")
		switch parts[3] {
		case "downcase":
			return strings.ToLower(result)
		case "upcase":
			return strings.ToUpper(result)
		default:
			return result
		}
	})
}

// captureRunes returns a defensive slice for a capture span. Malformed and
// unmatched spans produce nil rather than risking an out-of-bounds slice.
func captureRunes(line []rune, capture oniguruma.Capture) []rune {
	if capture.Start < 0 || capture.End < capture.Start || capture.End > len(line) {
		return nil
	}
	result := make([]rune, capture.End-capture.Start)
	copy(result, line[capture.Start:capture.End])
	return result
}

// capturedRunes returns one defensive rune slice per numbered capture while
// preserving capture order and unmatched slots.
func capturedRunes(line []rune, captures []oniguruma.Capture) [][]rune {
	result := make([][]rune, len(captures))
	for index, capture := range captures {
		result[index] = captureRunes(line, capture)
	}
	return result
}

// orderedRawCaptures turns a sparse numeric capture map into the indexed form
// consumed by compiled rules. Missing indexes remain nil.
func orderedRawCaptures(captures RawCaptures) []*RawRule {
	maximum := -1
	indexed := make(map[int]*RawRule, len(captures))
	for captureID, rule := range captures {
		index, err := strconv.Atoi(captureID)
		if err != nil || index < 0 {
			continue
		}
		indexed[index] = rule
		if index > maximum {
			maximum = index
		}
	}
	if maximum < 0 {
		return nil
	}

	result := make([]*RawRule, maximum+1)
	for index, rule := range indexed {
		result[index] = rule
	}
	return result
}
