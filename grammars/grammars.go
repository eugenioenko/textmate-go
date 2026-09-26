// Package grammars provides a curated, compressed set of TextMate grammars.
//
// Grammar data is decompressed and parsed on first use. Successful results and
// errors are cached, so callers must treat returned grammars as read-only.
package grammars

import (
	"bytes"
	"compress/gzip"
	"embed"
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"sync"

	textmate "github.com/eugenioenko/textmate-go"
)

//go:generate go run ../internal/cmd/gen-grammars -source ../../tm-grammars -revision 37edd1b26f18838050661d912334aba0ca7f4931 -selection curated.txt -license-reviews license-reviews.json -out .

//go:embed data/*.json.gz
var assets embed.FS

type cacheEntry struct {
	once    sync.Once
	grammar *textmate.RawGrammar
	err     error
}

var grammarCache sync.Map

// GrammarInfo describes an embedded grammar and the language names callers can
// use to select it. ID is the canonical language ID from the pinned grammar
// catalog. Aliases are alternate language labels suitable for Markdown code
// fences. FileTypes are the file types declared by the grammar itself.
//
// Lookup functions return independent copies, so callers may modify the slice
// fields without changing package data.
type GrammarInfo struct {
	ID          string
	DisplayName string
	ScopeName   string
	Aliases     []string
	FileTypes   []string
}

// Infos returns metadata for every embedded grammar, sorted by canonical ID.
func Infos() []GrammarInfo {
	infos := make([]GrammarInfo, len(generatedGrammarInfos))
	for index, info := range generatedGrammarInfos {
		infos[index] = cloneGrammarInfo(info)
	}
	return infos
}

// InfoForScope returns metadata for scopeName.
func InfoForScope(scopeName string) (GrammarInfo, bool) {
	index, ok := generatedInfoByScope[scopeName]
	if !ok {
		return GrammarInfo{}, false
	}
	return cloneGrammarInfo(generatedGrammarInfos[index]), true
}

// InfoForID returns metadata for a canonical language ID. Matching is
// case-insensitive and ignores surrounding whitespace. It does not search
// aliases; use InfoForAlias for alternate language labels.
func InfoForID(id string) (GrammarInfo, bool) {
	index, ok := generatedInfoByID[strings.ToLower(strings.TrimSpace(id))]
	if !ok {
		return GrammarInfo{}, false
	}
	return cloneGrammarInfo(generatedGrammarInfos[index]), true
}

// InfoForAlias returns metadata for an alternate language label, such as a
// Markdown code-fence label. Matching is case-insensitive and ignores
// surrounding whitespace. Canonical IDs are intentionally handled separately
// by InfoForID.
func InfoForAlias(alias string) (GrammarInfo, bool) {
	id, ok := generatedAliases[strings.ToLower(strings.TrimSpace(alias))]
	if !ok {
		return GrammarInfo{}, false
	}
	return InfoForID(id)
}

// InfoForFilename returns metadata for the grammar selected for name. It uses
// the same exact-name, compound-suffix, and extension matching as
// ScopeForFilename.
func InfoForFilename(name string) (GrammarInfo, bool) {
	scopeName := ScopeForFilename(name)
	if scopeName == "" {
		return GrammarInfo{}, false
	}
	return InfoForScope(scopeName)
}

func cloneGrammarInfo(info GrammarInfo) GrammarInfo {
	info.Aliases = append([]string(nil), info.Aliases...)
	info.FileTypes = append([]string(nil), info.FileTypes...)
	return info
}

// Load returns the embedded grammar with scopeName. It returns (nil, nil) when
// the scope is not in the curated set, making it directly usable as a
// textmate.RegistryOptions.LoadGrammar callback.
func Load(scopeName string) (*textmate.RawGrammar, error) {
	asset, ok := generatedAssets[scopeName]
	if !ok {
		return nil, nil
	}
	value, _ := grammarCache.LoadOrStore(scopeName, new(cacheEntry))
	entry := value.(*cacheEntry)
	entry.once.Do(func() {
		entry.grammar, entry.err = loadAsset(scopeName, asset)
	})
	return entry.grammar, entry.err
}

func loadAsset(scopeName, asset string) (*textmate.RawGrammar, error) {
	compressed, err := assets.ReadFile(asset)
	if err != nil {
		return nil, fmt.Errorf("grammars: read %s: %w", scopeName, err)
	}
	reader, err := gzip.NewReader(bytes.NewReader(compressed))
	if err != nil {
		return nil, fmt.Errorf("grammars: decompress %s: %w", scopeName, err)
	}
	data, readErr := io.ReadAll(reader)
	closeErr := reader.Close()
	if readErr != nil {
		return nil, fmt.Errorf("grammars: decompress %s: %w", scopeName, readErr)
	}
	if closeErr != nil {
		return nil, fmt.Errorf("grammars: decompress %s: %w", scopeName, closeErr)
	}
	grammar, err := textmate.ParseRawGrammar(data)
	if err != nil {
		return nil, fmt.Errorf("grammars: load %s: %w", scopeName, err)
	}
	if grammar.ScopeName != scopeName {
		return nil, fmt.Errorf("grammars: asset for %s contains %s", scopeName, grammar.ScopeName)
	}
	return grammar, nil
}

// ForFilename returns the embedded grammar associated with name. It recognizes
// exact conventional names such as Dockerfile and Makefile, compound suffixes,
// and ordinary extensions. Matching is case-insensitive. It returns (nil, nil)
// when no curated grammar is associated with the name.
func ForFilename(name string) (*textmate.RawGrammar, error) {
	scopeName := ScopeForFilename(name)
	if scopeName == "" {
		return nil, nil
	}
	return Load(scopeName)
}

// ScopeForFilename returns the scope associated with name, or an empty string
// when the curated set has no match.
func ScopeForFilename(name string) string {
	base := strings.ToLower(filepath.Base(name))
	if scopeName := generatedFilenames[base]; scopeName != "" {
		return scopeName
	}
	trimmed := strings.TrimPrefix(base, ".")
	for {
		if scopeName := generatedExtensions[trimmed]; scopeName != "" {
			return scopeName
		}
		dot := strings.IndexByte(trimmed, '.')
		if dot < 0 {
			return ""
		}
		trimmed = trimmed[dot+1:]
	}
}

// Scopes returns the sorted scope names available to Load.
func Scopes() []string {
	return append([]string(nil), generatedScopeNames...)
}
