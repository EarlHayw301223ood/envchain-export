// Package merge provides utilities for combining multiple envchain
// Chain instances into a single unified set of variables.
package merge

import (
	"fmt"

	"github.com/user/envchain-export/internal/chain"
)

// ConflictPolicy determines how key conflicts are handled during a merge.
type ConflictPolicy int

const (
	// PolicyError returns an error when a duplicate key is encountered.
	PolicyError ConflictPolicy = iota
	// PolicyOverwrite allows later chains to overwrite earlier ones.
	PolicyOverwrite
	// PolicySkip keeps the first value and ignores subsequent duplicates.
	PolicySkip
)

// Merge combines one or more Chain instances into dst according to the
// provided ConflictPolicy. The destination chain is mutated in place.
// Chains are processed in the order they are provided.
func Merge(dst *chain.Chain, policy ConflictPolicy, sources ...*chain.Chain) error {
	for _, src := range sources {
		for _, key := range src.Keys() {
			val, _ := src.Get(key)
			_, exists := dst.Get(key)

			switch {
			case !exists:
				if err := dst.Add(key, val); err != nil {
					return fmt.Errorf("merge: add %q: %w", key, err)
				}
			case policy == PolicyError:
				return fmt.Errorf("merge: duplicate key %q", key)
			case policy == PolicyOverwrite:
				if err := dst.Add(key, val); err != nil {
					return fmt.Errorf("merge: overwrite %q: %w", key, err)
				}
			case policy == PolicySkip:
				// intentionally keep existing value
			}
		}
	}
	return nil
}
