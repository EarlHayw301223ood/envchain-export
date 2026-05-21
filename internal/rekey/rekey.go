// Package rekey provides functionality to re-encrypt all scopes
// in the store from an old passphrase to a new passphrase in a single
// atomic-style operation.
package rekey

import (
	"fmt"

	"github.com/user/envchain-export/internal/rotate"
	"github.com/user/envchain-export/internal/scope"
	"github.com/user/envchain-export/internal/store"
)

// ErrNoScopes is returned when the store contains no scopes to rekey.
var ErrNoScopes = scope.ErrNoScopes

// Result holds the outcome of rekeying a single scope.
type Result struct {
	Scope   string
	Success bool
	Err     error
}

// Rekey re-encrypts every scope in storeDir from oldPass to newPass.
// It returns a slice of Result, one per scope, and a non-nil error only
// when the scope list itself cannot be retrieved.
func Rekey(storeDir, oldPass, newPass string) ([]Result, error) {
	scopes, err := scope.List(storeDir)
	if err != nil {
		return nil, fmt.Errorf("rekey: list scopes: %w", err)
	}

	results := make([]Result, 0, len(scopes))
	for _, name := range scopes {
		results = append(results, rekeyScope(storeDir, name, oldPass, newPass))
	}
	return results, nil
}

// rekeyScope attempts to rekey a single scope and returns the corresponding Result.
func rekeyScope(storeDir, name, oldPass, newPass string) Result {
	s, err := store.New(storeDir, name)
	if err != nil {
		return Result{Scope: name, Success: false, Err: fmt.Errorf("open store: %w", err)}
	}
	if err := rotate.Rotate(s, oldPass, newPass); err != nil {
		return Result{Scope: name, Success: false, Err: fmt.Errorf("rotate: %w", err)}
	}
	return Result{Scope: name, Success: true}
}
