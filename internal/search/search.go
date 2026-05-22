// Package search provides functionality to search for environment variable
// keys or values across all scopes in the store.
package search

import (
	"fmt"
	"strings"

	"github.com/your-org/envchain-export/internal/chain"
	"github.com/your-org/envchain-export/internal/scope"
	"github.com/your-org/envchain-export/internal/store"
)

// Match represents a single search result.
type Match struct {
	Scope string
	Key   string
	Value string
}

// Options controls search behaviour.
type Options struct {
	// SearchValues includes values in the search in addition to keys.
	SearchValues bool
	// CaseInsensitive performs a case-insensitive match.
	CaseInsensitive bool
}

// Search scans every scope accessible with passphrase and returns all entries
// whose key (and optionally value) contain the given query string.
func Search(storeDir, passphrase, query string, opts Options) ([]Match, error) {
	scopes, err := scope.List(storeDir)
	if err != nil {
		return nil, fmt.Errorf("search: list scopes: %w", err)
	}

	cmp := query
	if opts.CaseInsensitive {
		cmp = strings.ToLower(query)
	}

	var matches []Match
	for _, s := range scopes {
		st := store.New(storeDir, s)
		ch, err := st.Load(passphrase)
		if err != nil {
			// Skip scopes that cannot be decrypted with this passphrase.
			continue
		}
		for _, k := range ch.Keys() {
			v, _ := ch.Get(k)
			if matchField(k, cmp, opts.CaseInsensitive) ||
				(opts.SearchValues && matchField(v, cmp, opts.CaseInsensitive)) {
				matches = append(matches, Match{Scope: s, Key: k, Value: v})
			}
		}
	}
	return matches, nil
}

func matchField(field, cmp string, insensitive bool) bool {
	if insensitive {
		return strings.Contains(strings.ToLower(field), cmp)
	}
	return strings.Contains(field, cmp)
}

// FilterByScope returns only those matches belonging to the given scope.
func FilterByScope(matches []Match, scopeName string) []Match {
	var out []Match
	for _, m := range matches {
		if m.Scope == scopeName {
			out = append(out, m)
		}
	}
	return out
}

// Keys returns a deduplicated, sorted list of unique keys found across matches.
func Keys(matches []Match) []string {
	seen := make(map[string]struct{})
	var keys []string
	for _, m := range matches {
		if _, ok := seen[m.Key]; !ok {
			seen[m.Key] = struct{}{}
			keys = append(keys, m.Key)
		}
	}
	return keys
}

// ensure chain import is used via store; keep compiler happy in isolation.
var _ = (*chain.Chain)(nil)
