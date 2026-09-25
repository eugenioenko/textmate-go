// Command gen-grammars deterministically generates the embedded grammar set.
package main

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"encoding/json"
	"flag"
	"fmt"
	"go/format"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

type rawMetadata struct {
	ScopeName string   `json:"scopeName"`
	FileTypes []string `json:"fileTypes"`
}

type grammarMetadata struct {
	ID          string
	DisplayName string
	ScopeName   string
	FileTypes   []string
	Aliases     []string
	Asset       string
	License     string
	Source      string
	SHA         string
	RawBytes    int64
	GzipBytes   int64
}

type sourceMetadata struct {
	Name        string
	DisplayName string
	ScopeName   string
	License     string
	Source      string
	SHA         string
	Aliases     []string
}

type preparedGrammar struct {
	metadata grammarMetadata
	contents []byte
}

type generationPlan struct {
	revision  string
	selection string
	grammars  []preparedGrammar
	notice    []byte
	license   []byte
}

type generatorConfig struct {
	source    string
	revision  string
	selection string
	out       string
}

type candidate struct {
	ScopeName string
	Priority  int
}

func main() {
	source := flag.String("source", "", "path to the textmate-grammars-themes checkout")
	revision := flag.String("revision", "", "required source Git revision")
	selection := flag.String("selection", "", "newline-delimited grammar IDs, or 'all'")
	out := flag.String("out", ".", "output package directory")
	flag.Parse()
	if *source == "" || *revision == "" || *selection == "" {
		fatalf("-source, -revision, and -selection are required")
	}
	if err := generate(generatorConfig{
		source: *source, revision: *revision, selection: *selection, out: *out,
	}); err != nil {
		fatalf("%v", err)
	}
}

func generate(config generatorConfig) error {
	root, err := filepath.Abs(config.source)
	if err != nil {
		return err
	}
	if err := verifyRevision(root, config.revision); err != nil {
		return err
	}
	plan, err := prepareGeneration(root, config.revision, config.selection)
	if err != nil {
		return err
	}
	outDir, err := filepath.Abs(config.out)
	if err != nil {
		return err
	}
	return stageAndInstall(outDir, plan)
}

var grammarIDPattern = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._+-]*$`)

func prepareGeneration(root, revision, selection string) (*generationPlan, error) {
	grammarDir := filepath.Join(root, "packages", "tm-grammars", "grammars")
	ids, err := selectedIDs(grammarDir, selection)
	if err != nil {
		return nil, err
	}
	if len(ids) == 0 {
		return nil, fmt.Errorf("selection is empty")
	}
	sources, err := readSourceMetadata(filepath.Join(root, "packages", "tm-grammars", "index.js"))
	if err != nil {
		return nil, err
	}

	plan := &generationPlan{revision: revision, selection: selection}
	seenScopes := make(map[string]string, len(ids))
	for _, id := range ids {
		if !grammarIDPattern.MatchString(id) || filepath.Base(id) != id {
			return nil, fmt.Errorf("invalid grammar ID %q", id)
		}
		source, ok := sources[id]
		if !ok {
			return nil, fmt.Errorf("grammar %q is absent from package metadata", id)
		}
		if source.ScopeName == "" {
			return nil, fmt.Errorf("grammar %q metadata has no scopeName", id)
		}
		if source.DisplayName == "" {
			return nil, fmt.Errorf("grammar %q metadata has no displayName", id)
		}
		if source.License == "" {
			return nil, fmt.Errorf("grammar %q metadata has no license", id)
		}
		if !permissiveLicense(source.License) {
			return nil, fmt.Errorf("grammar %q has unreviewed license %q", id, source.License)
		}
		if source.Source == "" || source.SHA == "" {
			return nil, fmt.Errorf("grammar %q metadata has incomplete source provenance", id)
		}

		input := filepath.Join(grammarDir, id+".json")
		data, err := os.ReadFile(input)
		if err != nil {
			return nil, fmt.Errorf("read grammar %q: %w", id, err)
		}
		var raw rawMetadata
		if err := json.Unmarshal(data, &raw); err != nil {
			return nil, fmt.Errorf("parse grammar %q: %w", id, err)
		}
		if raw.ScopeName == "" {
			return nil, fmt.Errorf("grammar %q has no scopeName", id)
		}
		if raw.ScopeName != source.ScopeName {
			return nil, fmt.Errorf(
				"grammar %q scopeName %q does not match metadata %q",
				id, raw.ScopeName, source.ScopeName,
			)
		}
		if previous, exists := seenScopes[raw.ScopeName]; exists {
			return nil, fmt.Errorf("grammars %q and %q share scopeName %q", previous, id, raw.ScopeName)
		}
		seenScopes[raw.ScopeName] = id

		var compact bytes.Buffer
		if err := json.Compact(&compact, data); err != nil {
			return nil, fmt.Errorf("compact grammar %q: %w", id, err)
		}
		compact.WriteByte('\n')
		asset := filepath.ToSlash(filepath.Join("data", id+".json.gz"))
		plan.grammars = append(plan.grammars, preparedGrammar{
			contents: compact.Bytes(),
			metadata: grammarMetadata{
				ID: id, DisplayName: source.DisplayName,
				ScopeName: raw.ScopeName, FileTypes: raw.FileTypes,
				Asset: asset, Aliases: source.Aliases, License: source.License,
				Source: source.Source, SHA: source.SHA, RawBytes: int64(compact.Len()),
			},
		})
	}

	plan.notice, err = os.ReadFile(filepath.Join(root, "packages", "tm-grammars", "NOTICE"))
	if err != nil {
		return nil, fmt.Errorf("read source NOTICE: %w", err)
	}
	plan.license, err = os.ReadFile(filepath.Join(root, "packages", "tm-grammars", "LICENSE"))
	if err != nil {
		return nil, fmt.Errorf("read source LICENSE: %w", err)
	}
	return plan, nil
}

func readSourceMetadata(path string) (map[string]sourceMetadata, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open package metadata: %w", err)
	}
	defer func() { _ = file.Close() }()
	result := make(map[string]sourceMetadata)
	var current sourceMetadata
	inEntry := false
	inAliases := false
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		switch line {
		case "  {":
			inEntry = true
			inAliases = false
			current = sourceMetadata{}
		case "  },":
			if inEntry && current.Name != "" {
				if _, exists := result[current.Name]; exists {
					return nil, fmt.Errorf("duplicate package metadata for grammar %q", current.Name)
				}
				result[current.Name] = current
			}
			inEntry = false
			inAliases = false
		default:
			if !inEntry {
				continue
			}
			trimmed := strings.TrimSpace(line)
			if trimmed == "aliases: [" {
				inAliases = true
				continue
			}
			if inAliases {
				if trimmed == "]," {
					inAliases = false
					continue
				}
				if strings.HasPrefix(trimmed, "'") && strings.HasSuffix(trimmed, "',") {
					current.Aliases = append(current.Aliases, strings.TrimSuffix(strings.TrimPrefix(trimmed, "'"), "',"))
				}
				continue
			}
			key, value, ok := generatedStringField(line)
			if !ok {
				continue
			}
			switch key {
			case "name":
				current.Name = value
			case "displayName":
				current.DisplayName = value
			case "scopeName":
				current.ScopeName = value
			case "license":
				current.License = value
			case "source":
				current.Source = value
			case "sha":
				current.SHA = value
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("read package metadata: %w", err)
	}
	return result, nil
}

func generatedStringField(line string) (string, string, bool) {
	line = strings.TrimSpace(line)
	colon := strings.Index(line, ": '")
	if colon < 1 || !strings.HasSuffix(line, "',") {
		return "", "", false
	}
	key := line[:colon]
	value := strings.TrimSuffix(line[colon+3:], "',")
	value = strings.ReplaceAll(value, "\\'", "'")
	return key, value, true
}

func permissiveLicense(license string) bool {
	switch license {
	case "MIT", "Apache-2.0", "BSD-3-Clause", "ISC":
		return true
	default:
		return false
	}
}

func verifyRevision(root, want string) error {
	command := exec.Command("git", "-C", root, "rev-parse", "HEAD")
	output, err := command.Output()
	if err != nil {
		return fmt.Errorf("read source revision: %w", err)
	}
	got := strings.TrimSpace(string(output))
	if got != want {
		return fmt.Errorf("source revision is %s, want %s", got, want)
	}
	command = exec.Command("git", "-C", root, "status", "--porcelain")
	output, err = command.Output()
	if err != nil {
		return fmt.Errorf("read source status: %w", err)
	}
	if len(output) != 0 {
		return fmt.Errorf("source checkout is dirty")
	}
	return nil
}

func selectedIDs(grammarDir, selection string) ([]string, error) {
	if selection == "all" {
		matches, err := filepath.Glob(filepath.Join(grammarDir, "*.json"))
		if err != nil {
			return nil, err
		}
		ids := make([]string, 0, len(matches))
		for _, match := range matches {
			ids = append(ids, strings.TrimSuffix(filepath.Base(match), ".json"))
		}
		sort.Strings(ids)
		return ids, nil
	}
	file, err := os.Open(selection)
	if err != nil {
		return nil, fmt.Errorf("open selection: %w", err)
	}
	defer func() { _ = file.Close() }()
	seen := make(map[string]bool)
	var ids []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		id := strings.TrimSpace(strings.SplitN(scanner.Text(), "#", 2)[0])
		if id == "" {
			continue
		}
		if seen[id] {
			return nil, fmt.Errorf("duplicate grammar %q in selection", id)
		}
		seen[id] = true
		ids = append(ids, id)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	if !sort.StringsAreSorted(ids) {
		return nil, fmt.Errorf("selection must be sorted")
	}
	return ids, nil
}

func stageAndInstall(outDir string, plan *generationPlan) error {
	if err := os.MkdirAll(filepath.Dir(outDir), 0o755); err != nil {
		return fmt.Errorf("create output parent: %w", err)
	}
	stageDir, err := os.MkdirTemp(filepath.Dir(outDir), ".gen-grammars-stage-")
	if err != nil {
		return fmt.Errorf("create staging directory: %w", err)
	}
	defer func() { _ = os.RemoveAll(stageDir) }()
	stageDataDir := filepath.Join(stageDir, "data")
	if err := os.Mkdir(stageDataDir, 0o755); err != nil {
		return fmt.Errorf("create staged data directory: %w", err)
	}

	grammars := make([]grammarMetadata, len(plan.grammars))
	for index, grammar := range plan.grammars {
		metadata := grammar.metadata
		metadata.GzipBytes, err = writeGzip(
			filepath.Join(stageDir, filepath.FromSlash(metadata.Asset)),
			grammar.contents,
		)
		if err != nil {
			return fmt.Errorf("stage grammar %q: %w", metadata.ID, err)
		}
		grammars[index] = metadata
	}
	if err := writeGeneratedGo(filepath.Join(stageDir, "generated.go"), plan.revision, grammars); err != nil {
		return fmt.Errorf("stage generated.go: %w", err)
	}
	if err := os.WriteFile(filepath.Join(stageDir, "NOTICE"), plan.notice, 0o644); err != nil {
		return fmt.Errorf("stage NOTICE: %w", err)
	}
	if err := os.WriteFile(filepath.Join(stageDir, "LICENSE"), plan.license, 0o644); err != nil {
		return fmt.Errorf("stage LICENSE: %w", err)
	}
	if err := writeSource(filepath.Join(stageDir, "SOURCE"), plan.revision, plan.selection, grammars); err != nil {
		return fmt.Errorf("stage SOURCE: %w", err)
	}
	if err := writeManifest(filepath.Join(stageDir, "MANIFEST.md"), grammars); err != nil {
		return fmt.Errorf("stage MANIFEST.md: %w", err)
	}

	dataDir := filepath.Join(outDir, "data")
	if err := os.MkdirAll(dataDir, 0o755); err != nil {
		return fmt.Errorf("create output data directory: %w", err)
	}
	expectedAssets := make(map[string]struct{}, len(grammars))
	for _, grammar := range grammars {
		expectedAssets[filepath.Base(grammar.Asset)] = struct{}{}
		if err := installFile(
			filepath.Join(stageDir, filepath.FromSlash(grammar.Asset)),
			filepath.Join(outDir, filepath.FromSlash(grammar.Asset)),
		); err != nil {
			return fmt.Errorf("install grammar %q: %w", grammar.ID, err)
		}
	}
	for _, name := range []string{"generated.go", "NOTICE", "LICENSE", "SOURCE", "MANIFEST.md"} {
		if err := installFile(filepath.Join(stageDir, name), filepath.Join(outDir, name)); err != nil {
			return fmt.Errorf("install %s: %w", name, err)
		}
	}
	if err := removeStaleGeneratedAssets(dataDir, expectedAssets); err != nil {
		return err
	}
	return nil
}

func installFile(source, destination string) error {
	input, err := os.Open(source)
	if err != nil {
		return err
	}
	defer func() { _ = input.Close() }()
	temporary, err := os.CreateTemp(filepath.Dir(destination), "."+filepath.Base(destination)+".tmp-")
	if err != nil {
		return err
	}
	temporaryName := temporary.Name()
	removeTemporary := true
	defer func() {
		if removeTemporary {
			_ = os.Remove(temporaryName)
		}
	}()
	if _, err := io.Copy(temporary, input); err != nil {
		_ = temporary.Close()
		return err
	}
	if err := temporary.Chmod(0o644); err != nil {
		_ = temporary.Close()
		return err
	}
	if err := temporary.Close(); err != nil {
		return err
	}
	if err := os.Rename(temporaryName, destination); err != nil {
		return err
	}
	removeTemporary = false
	return nil
}

func removeStaleGeneratedAssets(dataDir string, expected map[string]struct{}) error {
	entries, err := os.ReadDir(dataDir)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".json.gz") {
			continue
		}
		if _, keep := expected[entry.Name()]; keep {
			continue
		}
		if err := os.Remove(filepath.Join(dataDir, entry.Name())); err != nil {
			return err
		}
	}
	return nil
}

func writeGzip(path string, data []byte) (int64, error) {
	file, err := os.Create(path)
	if err != nil {
		return 0, err
	}
	writer, err := gzip.NewWriterLevel(file, gzip.BestCompression)
	if err != nil {
		_ = file.Close()
		return 0, err
	}
	writer.OS = 255
	_, writeErr := writer.Write(data)
	closeErr := writer.Close()
	fileCloseErr := file.Close()
	if writeErr != nil {
		return 0, writeErr
	}
	if closeErr != nil {
		return 0, closeErr
	}
	if fileCloseErr != nil {
		return 0, fileCloseErr
	}
	info, err := os.Stat(path)
	if err != nil {
		return 0, err
	}
	return info.Size(), nil
}

func writeGeneratedGo(path, revision string, grammars []grammarMetadata) error {
	grammars = append([]grammarMetadata(nil), grammars...)
	sort.Slice(grammars, func(i, j int) bool { return grammars[i].ID < grammars[j].ID })

	extensions := make(map[string]candidate)
	filenames := map[string]candidate{
		"dockerfile":    {ScopeName: "source.dockerfile", Priority: 100},
		"containerfile": {ScopeName: "source.dockerfile", Priority: 100},
		"makefile":      {ScopeName: "source.makefile", Priority: 100},
		"gnumakefile":   {ScopeName: "source.makefile", Priority: 100},
		".bashrc":       {ScopeName: "source.shell", Priority: 100},
		".zshrc":        {ScopeName: "source.shell", Priority: 100},
		".profile":      {ScopeName: "source.shell", Priority: 100},
	}
	reviewedExtensions := map[string]candidate{
		"m":   {ScopeName: "source.objc", Priority: 100},
		"mm":  {ScopeName: "source.objcpp", Priority: 100},
		"pl":  {ScopeName: "source.perl", Priority: 100},
		"cc":  {ScopeName: "source.cpp", Priority: 100},
		"cxx": {ScopeName: "source.cpp", Priority: 100},
		"c++": {ScopeName: "source.cpp", Priority: 100},
		"hh":  {ScopeName: "source.cpp", Priority: 100},
		"hpp": {ScopeName: "source.cpp", Priority: 100},
		"hxx": {ScopeName: "source.cpp", Priority: 100},
		"h++": {ScopeName: "source.cpp", Priority: 100},
	}
	filenames[".env"] = candidate{ScopeName: "source.dotenv", Priority: 100}
	available := make(map[string]bool, len(grammars))
	for _, grammar := range grammars {
		available[grammar.ScopeName] = true
		addCandidate(extensions, grammar.ID, grammar.ScopeName, 20)
		for _, alias := range grammar.Aliases {
			addCandidate(extensions, strings.ToLower(alias), grammar.ScopeName, 25)
		}
		for _, fileType := range grammar.FileTypes {
			key := strings.ToLower(strings.TrimPrefix(strings.TrimSpace(fileType), "."))
			if key == "" {
				continue
			}
			priority := 10
			if key == grammar.ID {
				priority = 30
			}
			addCandidate(extensions, key, grammar.ScopeName, priority)
			if strings.Contains(key, ".") || fileType != strings.ToLower(fileType) {
				addCandidate(filenames, strings.ToLower(fileType), grammar.ScopeName, priority)
			}
		}
	}
	for key, value := range filenames {
		if !available[value.ScopeName] {
			delete(filenames, key)
		}
	}
	for key, value := range reviewedExtensions {
		if available[value.ScopeName] {
			addCandidate(extensions, key, value.ScopeName, value.Priority)
		}
	}
	var buffer bytes.Buffer
	fmt.Fprintln(&buffer, "// Code generated by gen-grammars; DO NOT EDIT.")
	fmt.Fprintln(&buffer, "package grammars")
	fmt.Fprintln(&buffer)
	fmt.Fprintf(&buffer, "const SourceRevision = %q\n\n", revision)
	fmt.Fprintln(&buffer, "var generatedGrammarInfos = []GrammarInfo{")
	for _, grammar := range grammars {
		fmt.Fprintln(&buffer, "\t{")
		fmt.Fprintf(&buffer, "\t\tID: %q,\n", grammar.ID)
		fmt.Fprintf(&buffer, "\t\tDisplayName: %q,\n", grammar.DisplayName)
		fmt.Fprintf(&buffer, "\t\tScopeName: %q,\n", grammar.ScopeName)
		writeStringSliceField(&buffer, "Aliases", grammar.Aliases)
		writeStringSliceField(&buffer, "FileTypes", grammar.FileTypes)
		fmt.Fprintln(&buffer, "\t},")
	}
	fmt.Fprintln(&buffer, "}")
	infoByID := make(map[string]int, len(grammars))
	infoByScope := make(map[string]int, len(grammars))
	aliases := make(map[string]string)
	for index, grammar := range grammars {
		infoByID[strings.ToLower(grammar.ID)] = index
		infoByScope[grammar.ScopeName] = index
		for _, alias := range grammar.Aliases {
			key := strings.ToLower(strings.TrimSpace(alias))
			if key == "" {
				continue
			}
			// Keep future catalog collisions reproducible. Canonical IDs live in a
			// separate lookup, so an alias can never displace an ID.
			if current, ok := aliases[key]; !ok || grammar.ID < current {
				aliases[key] = grammar.ID
			}
		}
	}
	writeIndexMap(&buffer, "generatedInfoByID", infoByID)
	writeIndexMap(&buffer, "generatedInfoByScope", infoByScope)
	writeStringMap(&buffer, "generatedAliases", aliases)
	fmt.Fprintln(&buffer, "var generatedAssets = map[string]string{")
	for _, grammar := range grammars {
		fmt.Fprintf(&buffer, "\t%q: %q,\n", grammar.ScopeName, grammar.Asset)
	}
	fmt.Fprintln(&buffer, "}")
	fmt.Fprintln(&buffer, "var generatedScopeNames = []string{")
	scopes := make([]string, 0, len(grammars))
	for _, grammar := range grammars {
		scopes = append(scopes, grammar.ScopeName)
	}
	sort.Strings(scopes)
	for _, scopeName := range scopes {
		fmt.Fprintf(&buffer, "\t%q,\n", scopeName)
	}
	fmt.Fprintln(&buffer, "}")
	writeCandidateMap(&buffer, "generatedExtensions", extensions)
	writeCandidateMap(&buffer, "generatedFilenames", filenames)
	formatted, err := format.Source(buffer.Bytes())
	if err != nil {
		return err
	}
	return os.WriteFile(path, formatted, 0o644)
}

func writeStringSliceField(buffer *bytes.Buffer, name string, values []string) {
	if len(values) == 0 {
		return
	}
	fmt.Fprintf(buffer, "\t\t%s: []string{", name)
	for index, value := range values {
		if index != 0 {
			fmt.Fprint(buffer, ", ")
		}
		fmt.Fprintf(buffer, "%q", value)
	}
	fmt.Fprintln(buffer, "},")
}

func writeIndexMap(buffer *bytes.Buffer, name string, values map[string]int) {
	fmt.Fprintf(buffer, "var %s = map[string]int{\n", name)
	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		fmt.Fprintf(buffer, "\t%q: %d,\n", key, values[key])
	}
	fmt.Fprintln(buffer, "}")
}

func writeStringMap(buffer *bytes.Buffer, name string, values map[string]string) {
	fmt.Fprintf(buffer, "var %s = map[string]string{\n", name)
	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		fmt.Fprintf(buffer, "\t%q: %q,\n", key, values[key])
	}
	fmt.Fprintln(buffer, "}")
}

func addCandidate(values map[string]candidate, key, scopeName string, priority int) {
	current, ok := values[key]
	if !ok || priority > current.Priority || (priority == current.Priority && scopeName < current.ScopeName) {
		values[key] = candidate{ScopeName: scopeName, Priority: priority}
	}
}

func writeCandidateMap(buffer *bytes.Buffer, name string, values map[string]candidate) {
	fmt.Fprintf(buffer, "var %s = map[string]string{\n", name)
	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		fmt.Fprintf(buffer, "\t%q: %q,\n", key, values[key].ScopeName)
	}
	fmt.Fprintln(buffer, "}")
}

func writeSource(path, revision, selection string, grammars []grammarMetadata) error {
	var rawBytes, gzipBytes int64
	for _, grammar := range grammars {
		rawBytes += grammar.RawBytes
		gzipBytes += grammar.GzipBytes
	}
	contents := fmt.Sprintf("Source: https://github.com/shikijs/textmate-grammars-themes\nRevision: %s\nSelection: %s\nGrammars: %d\nCompacted JSON bytes: %d\nIndividually gzip-compressed bytes: %d\n", revision, filepath.Base(selection), len(grammars), rawBytes, gzipBytes)
	return os.WriteFile(path, []byte(contents), 0o644)
}

func writeManifest(path string, grammars []grammarMetadata) error {
	var buffer bytes.Buffer
	fmt.Fprintln(&buffer, "# Embedded grammar manifest")
	fmt.Fprintln(&buffer)
	fmt.Fprintln(&buffer, "| ID | Scope | License | Upstream source |")
	fmt.Fprintln(&buffer, "| --- | --- | --- | --- |")
	for _, grammar := range grammars {
		source := grammar.Source
		if grammar.SHA != "" {
			source = fmt.Sprintf("[%s](%s) (`%s`)", grammar.ID, grammar.Source, grammar.SHA)
		}
		fmt.Fprintf(&buffer, "| `%s` | `%s` | %s | %s |\n", grammar.ID, grammar.ScopeName, grammar.License, source)
	}
	return os.WriteFile(path, buffer.Bytes(), 0o644)
}

func fatalf(format string, arguments ...any) {
	fmt.Fprintf(os.Stderr, "gen-grammars: "+format+"\n", arguments...)
	os.Exit(1)
}
