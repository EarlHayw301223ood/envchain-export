// Package expire provides functionality to check and purge scopes
// whose TTL has elapsed, removing them from the store entirely.
package expire

import (
	"fmt"
	"time"

	"github.com/user/envchain-export/internal/scope"
	"github.com/user/envchain-export/internal/ttl"
)

// Result holds the outcome of an expiry check for a single scope.
type Result struct {
	Scope   string
	Expired bool
	Expiry  time.Time
}

// ErrNoScopes is returned when the store contains no scopes.
var ErrNoScopes = scope.ErrNoScopes

// Check inspects all scopes in storeDir and returns a Result for each one
// that has a TTL configured. Scopes without a TTL are omitted.
func Check(storeDir string) ([]Result, error) {
	scopes, err := scope.List(storeDir)
	if err != nil {
		return nil, err
	}

	var results []Result
	for _, s := range scopes {
		expiry, err := ttl.Get(storeDir, s)
		if err == ttl.ErrNoExpiry {
			continue
		}
		if err != nil {
			return nil, fmt.Errorf("expire: check ttl for %q: %w", s, err)
		}
		results = append(results, Result{
			Scope:   s,
			Expired: time.Now().UTC().After(expiry),
			Expiry:  expiry,
		})
	}
	return results, nil
}

// Purge deletes all scopes in storeDir whose TTL has elapsed.
// It returns the names of the scopes that were removed.
func Purge(storeDir string) ([]string, error) {
	results, err := Check(storeDir)
	if err != nil {
		return nil, err
	}

	var removed []string
	for _, r := range results {
		if !r.Expired {
			continue
		}
		if err := scope.Delete(storeDir, r.Scope); err != nil {
			return removed, fmt.Errorf("expire: purge %q: %w", r.Scope, err)
		}
		removed = append(removed, r.Scope)
	}
	return removed, nil
}
