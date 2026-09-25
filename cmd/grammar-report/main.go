// Command grammar-report scans a tm-grammars checkout for parser, dependency,
// and regular-expression compatibility failures.
package main

import (
	"errors"
	"flag"
	"fmt"
	"html"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	textmate "github.com/eugenioenko/textmate-go"
	"github.com/eugenioenko/textmate-go/oniguruma"
)

const pinnedGrammarCount = 260

type options struct {
	grammarDir string
	revision   string
	wantCount  int
}

type grammarFile struct {
	name    string
	grammar *textmate.RawGrammar
}

type loadFailure struct {
	file    string
	scope   string
	message string
}

type regexDiagnostic struct {
	file            string
	scope           string
	path            string
	field           string
	pattern         string
	compiledPattern string
	diagnostic      oniguruma.Diagnostic
}

type report struct {
	revision         string
	grammarDir       string
	files            int
	parsed           int
	loaded           int
	regexFields      int
	patternsAffected int
	loadFailures     []loadFailure
	diagnostics      []regexDiagnostic
}

func main() {
	var opts options
	flag.StringVar(&opts.grammarDir, "grammars", "../tm-grammars/packages/tm-grammars/grammars", "directory containing tm-grammars JSON files")
	flag.StringVar(&opts.revision, "revision", "", "required tm-grammars git revision; the checkout must match it")
	flag.IntVar(&opts.wantCount, "want-grammars", pinnedGrammarCount, "required number of grammar JSON files; zero disables the check")
	flag.Parse()

	if opts.revision == "" {
		fatal(errors.New("-revision is required"))
	}
	result, err := scan(opts)
	if err != nil {
		fatal(err)
	}
	if err := writeReport(os.Stdout, result); err != nil {
		fatal(err)
	}
}

func fatal(err error) {
	fmt.Fprintf(os.Stderr, "grammar-report: %v\n", err)
	os.Exit(1)
}

func scan(opts options) (*report, error) {
	if err := verifyGrammarCheckout(opts.grammarDir, opts.revision); err != nil {
		return nil, err
	}

	paths, err := filepath.Glob(filepath.Join(opts.grammarDir, "*.json"))
	if err != nil {
		return nil, fmt.Errorf("list grammars: %w", err)
	}
	sort.Strings(paths)
	if opts.wantCount > 0 && len(paths) != opts.wantCount {
		return nil, fmt.Errorf("found %d grammar JSON files in %s, want %d", len(paths), opts.grammarDir, opts.wantCount)
	}

	result := &report{
		revision:   opts.revision,
		grammarDir: filepath.ToSlash(filepath.Clean(opts.grammarDir)),
		files:      len(paths),
	}
	files := make([]grammarFile, 0, len(paths))
	byScope := make(map[string]*textmate.RawGrammar, len(paths))
	fileByScope := make(map[string]string, len(paths))

	for _, path := range paths {
		name := filepath.Base(path)
		data, readErr := os.ReadFile(path)
		if readErr != nil {
			result.loadFailures = append(result.loadFailures, loadFailure{file: name, message: "read: " + readErr.Error()})
			continue
		}
		grammar, parseErr := textmate.ParseRawGrammar(data)
		if parseErr != nil {
			result.loadFailures = append(result.loadFailures, loadFailure{file: name, message: parseErr.Error()})
			continue
		}
		result.parsed++
		if previous, exists := fileByScope[grammar.ScopeName]; exists {
			result.loadFailures = append(result.loadFailures, loadFailure{
				file: name, scope: grammar.ScopeName,
				message: fmt.Sprintf("duplicate scope name (also in %s)", previous),
			})
			continue
		}
		fileByScope[grammar.ScopeName] = name
		byScope[grammar.ScopeName] = grammar
		files = append(files, grammarFile{name: name, grammar: grammar})
	}

	registry := textmate.NewRegistry(textmate.RegistryOptions{
		LoadGrammar: func(scopeName string) (*textmate.RawGrammar, error) {
			return byScope[scopeName], nil
		},
	})
	defer registry.Dispose()
	for _, file := range files {
		if _, loadErr := registry.LoadGrammar(file.grammar.ScopeName); loadErr != nil {
			result.loadFailures = append(result.loadFailures, loadFailure{
				file: file.name, scope: file.grammar.ScopeName, message: loadErr.Error(),
			})
			continue
		}
		result.loaded++
	}

	for _, file := range files {
		scanGrammar(result, file)
	}
	sort.Slice(result.loadFailures, func(i, j int) bool {
		if result.loadFailures[i].file != result.loadFailures[j].file {
			return result.loadFailures[i].file < result.loadFailures[j].file
		}
		return result.loadFailures[i].message < result.loadFailures[j].message
	})
	sort.SliceStable(result.diagnostics, func(i, j int) bool {
		a, b := result.diagnostics[i], result.diagnostics[j]
		if a.file != b.file {
			return a.file < b.file
		}
		if a.path != b.path {
			return a.path < b.path
		}
		return a.diagnostic.Kind < b.diagnostic.Kind
	})
	return result, nil
}

func verifyGrammarCheckout(grammarDir, wantRevision string) error {
	root, err := gitOutput(grammarDir, "rev-parse", "--show-toplevel")
	if err != nil {
		return fmt.Errorf("locate tm-grammars checkout: %w", err)
	}

	revision, err := gitOutput(root, "rev-parse", "--verify", "HEAD^{commit}")
	if err != nil {
		return fmt.Errorf("read tm-grammars checkout revision: %w", err)
	}
	if revision != wantRevision {
		return fmt.Errorf("tm-grammars checkout revision is %s, want %s", revision, wantRevision)
	}

	status, err := gitOutput(root, "status", "--porcelain=v1", "--untracked-files=all")
	if err != nil {
		return fmt.Errorf("inspect tm-grammars checkout cleanliness: %w", err)
	}
	if status != "" {
		return fmt.Errorf("tm-grammars checkout is not clean:\n%s", status)
	}
	return nil
}

func gitOutput(dir string, args ...string) (string, error) {
	commandArgs := append([]string{"-C", dir}, args...)
	output, err := exec.Command("git", commandArgs...).CombinedOutput()
	trimmed := strings.TrimRight(string(output), "\r\n")
	if err != nil {
		if trimmed == "" {
			return "", err
		}
		return "", fmt.Errorf("git %s: %s: %w", strings.Join(args, " "), trimmed, err)
	}
	return trimmed, nil
}

func scanGrammar(result *report, file grammarFile) {
	affected := make(map[string]struct{})
	scanPattern := func(path, field, pattern string, substituteBackReferences bool) {
		result.regexFields++
		compiledPattern := pattern
		if substituteBackReferences {
			compiledPattern = replaceNumericBackReferences(pattern, "x")
		}
		scanner := oniguruma.NewScanner([]string{compiledPattern})
		for _, diagnostic := range scanner.Diagnostics() {
			result.diagnostics = append(result.diagnostics, regexDiagnostic{
				file: file.name, scope: file.grammar.ScopeName, path: path, field: field,
				pattern: pattern, compiledPattern: compiledPattern, diagnostic: diagnostic,
			})
			affected[file.name+"\x00"+path] = struct{}{}
		}
	}

	if file.grammar.FirstLineMatch != "" {
		scanPattern("$.firstLineMatch", "firstLineMatch", file.grammar.FirstLineMatch, false)
	}
	seen := make(map[*textmate.RawRule]struct{})
	var walkRule func(string, *textmate.RawRule)
	var walkRules func(string, []*textmate.RawRule)
	var walkMap func(string, map[string]*textmate.RawRule)
	walkRule = func(path string, rule *textmate.RawRule) {
		if rule == nil {
			return
		}
		if _, exists := seen[rule]; exists {
			return
		}
		seen[rule] = struct{}{}
		if rule.Match != nil {
			scanPattern(path+".match", "match", *rule.Match, false)
		}
		if rule.Begin != nil {
			scanPattern(path+".begin", "begin", *rule.Begin, false)
		}
		if rule.End != nil {
			scanPattern(path+".end", "end", *rule.End, true)
		}
		if rule.While != nil {
			scanPattern(path+".while", "while", *rule.While, true)
		}
		walkRules(path+".patterns", rule.Patterns)
		walkMap(path+".captures", map[string]*textmate.RawRule(rule.Captures))
		walkMap(path+".beginCaptures", map[string]*textmate.RawRule(rule.BeginCaptures))
		walkMap(path+".endCaptures", map[string]*textmate.RawRule(rule.EndCaptures))
		walkMap(path+".whileCaptures", map[string]*textmate.RawRule(rule.WhileCaptures))
		walkMap(path+".repository", map[string]*textmate.RawRule(rule.Repository))
	}
	walkRules = func(path string, rules []*textmate.RawRule) {
		for index, rule := range rules {
			walkRule(path+"["+strconv.Itoa(index)+"]", rule)
		}
	}
	walkMap = func(path string, rules map[string]*textmate.RawRule) {
		keys := make([]string, 0, len(rules))
		for key := range rules {
			keys = append(keys, key)
		}
		sort.Strings(keys)
		for _, key := range keys {
			walkRule(path+"["+strconv.Quote(key)+"]", rules[key])
		}
	}

	walkRules("$.patterns", file.grammar.Patterns)
	walkMap("$.repository", map[string]*textmate.RawRule(file.grammar.Repository))
	walkMap("$.injections", map[string]*textmate.RawRule(file.grammar.Injections))
	result.patternsAffected += len(affected)
}

func replaceNumericBackReferences(source, replacement string) string {
	var result strings.Builder
	last := 0
	for pos := 0; pos+1 < len(source); pos++ {
		if source[pos] != '\\' || source[pos+1] < '0' || source[pos+1] > '9' {
			continue
		}
		end := pos + 2
		for end < len(source) && source[end] >= '0' && source[end] <= '9' {
			end++
		}
		result.WriteString(source[last:pos])
		result.WriteString(replacement)
		last = end
		pos = end - 1
	}
	if last == 0 {
		return source
	}
	result.WriteString(source[last:])
	return result.String()
}

func writeReport(w io.Writer, result *report) error {
	unsupported, compileErrors, other := partitionDiagnostics(result.diagnostics)
	var buffer strings.Builder
	output := &reportWriter{writer: &buffer}
	write := output.printf

	write("# Grammar compatibility report\n\n")
	write("This report scans the pinned `textmate-grammars-themes` grammar corpus at revision `%s`. ", result.revision)
	write("It parses and registry-loads every grammar, then sends every regex field reachable through the parsed grammar model through the same `oniguruma.NewScanner` translation and compilation path used by the tokenizer.\n\n")
	write("Reproduce from the `textmate-go` repository root:\n\n")
	write("```sh\n")
	write("report_file=$(mktemp)\n")
	write("go run ./cmd/grammar-report -grammars %s -revision %s > \"$report_file\"\n", shellQuote(result.grammarDir), shellQuote(result.revision))
	write("mv \"$report_file\" docs/grammar-report.md\n")
	write("```\n\n")
	write("Numeric capture references in `end` and `while` fields are replaced with a safe literal before compilation. At runtime those references are replaced with escaped text captured by the corresponding `begin`; compiling them as standalone backreferences would report false failures. The report always shows the original grammar pattern.\n\n")

	write("## Summary\n\n")
	write("| Check | Count |\n|---|---:|\n")
	write("| Grammar JSON files | %d |\n", result.files)
	write("| Parsed grammars | %d |\n", result.parsed)
	write("| Registry load successes | %d |\n", result.loaded)
	write("| Grammar load failures | %d |\n", len(result.loadFailures))
	write("| Regex fields scanned | %d |\n", result.regexFields)
	write("| Patterns with diagnostics | %d |\n", result.patternsAffected)
	write("| Expected unsupported-syntax diagnostics | %d |\n", len(unsupported))
	write("| Pattern compile failures | %d |\n", len(compileErrors))
	write("| Other diagnostics | %d |\n\n", len(other))

	write("`unsupported_syntax` entries are known Oniguruma constructs that regexp2 cannot faithfully execute; the adapter degrades the affected construct or pattern and records it explicitly. `compile_error` entries are separate: translation completed, but regexp2 rejected the result. A static compile scan cannot produce match-time timeout or match-error diagnostics.\n\n")

	write("## Grammar load failures\n\n")
	if len(result.loadFailures) == 0 {
		write("None.\n\n")
	} else {
		for _, failure := range result.loadFailures {
			write("- `%s`", markdownCode(failure.file))
			if failure.scope != "" {
				write(" (`%s`)", markdownCode(failure.scope))
			}
			write(": %s\n", html.EscapeString(failure.message))
		}
		write("\n")
	}

	writeDiagnosticSection(output, "Expected unsupported syntax", unsupported)
	writeDiagnosticSection(output, "Pattern compile failures", compileErrors)
	writeDiagnosticSection(output, "Other regex diagnostics", other)
	if output.err != nil {
		return output.err
	}
	_, err := io.WriteString(w, strings.TrimRight(buffer.String(), "\n")+"\n")
	return err
}

type reportWriter struct {
	writer io.Writer
	err    error
}

func (w *reportWriter) printf(format string, args ...any) {
	if w.err != nil {
		return
	}
	_, w.err = fmt.Fprintf(w.writer, format, args...)
}

func partitionDiagnostics(all []regexDiagnostic) (unsupported, compileErrors, other []regexDiagnostic) {
	for _, diagnostic := range all {
		switch diagnostic.diagnostic.Kind {
		case oniguruma.DiagnosticUnsupportedSyntax:
			unsupported = append(unsupported, diagnostic)
		case oniguruma.DiagnosticCompileError:
			compileErrors = append(compileErrors, diagnostic)
		default:
			other = append(other, diagnostic)
		}
	}
	return unsupported, compileErrors, other
}

func writeDiagnosticSection(w *reportWriter, title string, diagnostics []regexDiagnostic) {
	w.printf("## %s\n\n", title)
	if len(diagnostics) == 0 {
		w.printf("None.\n\n")
		return
	}
	for index, item := range diagnostics {
		w.printf("### %d. `%s` — `%s`\n\n", index+1, markdownCode(item.file), markdownCode(item.path))
		w.printf("- Scope: `%s`\n", markdownCode(item.scope))
		w.printf("- Field: `%s`\n", markdownCode(item.field))
		w.printf("- Kind: `%s`\n", item.diagnostic.Kind)
		w.printf("- Reason: %s\n", html.EscapeString(item.diagnostic.Message))
		w.printf("- Original pattern:\n\n")
		writeCodeBlock(w, item.pattern)
		translated := item.diagnostic.Translated
		if translated != "" && translated != item.compiledPattern {
			w.printf("- Adapter translation:\n\n")
			writeCodeBlock(w, translated)
		}
	}
}

func writeCodeBlock(w *reportWriter, value string) {
	longest := 2
	current := 0
	for _, ch := range value {
		if ch == '`' {
			current++
			if current > longest {
				longest = current
			}
		} else {
			current = 0
		}
	}
	fence := strings.Repeat("`", longest+1)
	w.printf("%stext\n%s\n%s\n\n", fence, value, fence)
}

func markdownCode(value string) string {
	return strings.ReplaceAll(value, "`", "&#96;")
}

func shellQuote(value string) string {
	if value != "" && strings.IndexFunc(value, func(r rune) bool {
		safe := r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z' || r >= '0' && r <= '9' || strings.ContainsRune("_@%+=:,./-", r)
		return !safe
	}) == -1 {
		return value
	}
	return "'" + strings.ReplaceAll(value, "'", "'\\''") + "'"
}
