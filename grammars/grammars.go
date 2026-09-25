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

//go:generate go run ../internal/cmd/gen-grammars -source ../../tm-grammars -revision 37edd1b26f18838050661d912334aba0ca7f4931 -selection curated.txt -out .

//go:embed data/*.json.gz
var assets embed.FS

type cacheEntry struct {
	once    sync.Once
	grammar *textmate.RawGrammar
	err     error
}

var grammarCache sync.Map

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
