// Package prune removes environment variable entries from a scope
// whose keys match a given pattern or predicate.
package prune

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/user/envchain-export/internal/chain"
	"github.com/user/envchain-export/internal/store"
)

// Result holds the outcome of a prune operation.
type Result struct {
	Scope   string
	Removed []string
}

// ByPrefix removes all keys that start with the given prefix from the
// named scope. It decrypts with passphrase, mutates the chain, and
// re-encrypts before returning the Result.
func ByPrefix(storeDir, scope, prefix, passphrase string) (Result, error) {
	if prefix == "" {
		return Result{}, fmt.Errorf("prune: prefix must not be empty")
	}

	s := store.New(storeDir)
	ch, err := s.Load(scope, passphrase)
	if err != nil {
		return Result{}, fmt.Errorf("prune: load %q: %w", scope, err)
	}

	var removed []string
	for _, k := range ch.Keys() {
		if strings.HasPrefix(k, prefix) {
			ch.Remove(k)
			removed = append(removed, k)
		}
	}

	if len(removed) == 0 {
		return Result{Scope: scope, Removed: nil}, nil
	}

	if err := s.Save(scope, ch, passphrase); err != nil {
		return Result{}, fmt.Errorf("prune: save %q: %w", scope, err)
	}

	return Result{Scope: scope, Removed: removed}, nil
}

// ByGlob removes all keys whose names match the given shell glob pattern.
func ByGlob(storeDir, scope, pattern, passphrase string) (Result, error) {
	if pattern == "" {
		return Result{}, fmt.Errorf("prune: pattern must not be empty")
	}

	// Validate pattern early so we surface syntax errors before loading.
	if _, err := filepath.Match(pattern, ""); err != nil {
		return Result{}, fmt.Errorf("prune: invalid pattern %q: %w", pattern, err)
	}

	s := store.New(storeDir)
	ch, err := s.Load(scope, passphrase)
	if err != nil {
		return Result{}, fmt.Errorf("prune: load %q: %w", scope, err)
	}

	var removed []string
	for _, k := range ch.Keys() {
		matched, _ := filepath.Match(pattern, k)
		if matched {
			ch.Remove(k)
			removed = append(removed, k)
		}
	}

	if len(removed) == 0 {
		return Result{Scope: scope, Removed: nil}, nil
	}

	if err := s.Save(scope, ch, passphrase); err != nil {
		return Result{}, fmt.Errorf("prune: save %q: %w", scope, err)
	}

	return Result{Scope: scope, Removed: removed}, nil
}

// helper used by tests — returns a bare *chain.Chain via the store interface.
func loadChain(storeDir, scope, passphrase string) (*chain.Chain, error) {
	return store.New(storeDir).Load(scope, passphrase)
}
