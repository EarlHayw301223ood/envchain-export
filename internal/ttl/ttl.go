// Package ttl provides expiry tracking for scopes.
// A scope can be assigned an expiry time; after that time
// it is considered expired and callers may act accordingly.
package ttl

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// ErrNoExpiry is returned when a scope has no expiry set.
var ErrNoExpiry = errors.New("ttl: no expiry set for scope")

// ErrExpired is returned when a scope's expiry has passed.
var ErrExpired = errors.New("ttl: scope has expired")

type entry struct {
	ExpiresAt time.Time `json:"expires_at"`
}

func indexPath(storeDir, scope string) string {
	return filepath.Join(storeDir, scope+".ttl.json")
}

// Set assigns an expiry duration to a scope, measured from now.
func Set(storeDir, scope string, d time.Duration) error {
	if err := os.MkdirAll(storeDir, 0700); err != nil {
		return fmt.Errorf("ttl: mkdir: %w", err)
	}
	e := entry{ExpiresAt: time.Now().UTC().Add(d)}
	data, err := json.Marshal(e)
	if err != nil {
		return fmt.Errorf("ttl: marshal: %w", err)
	}
	if err := os.WriteFile(indexPath(storeDir, scope), data, 0600); err != nil {
		return fmt.Errorf("ttl: write: %w", err)
	}
	return nil
}

// Get returns the expiry time for a scope.
// Returns ErrNoExpiry if no expiry has been set.
func Get(storeDir, scope string) (time.Time, error) {
	data, err := os.ReadFile(indexPath(storeDir, scope))
	if errors.Is(err, os.ErrNotExist) {
		return time.Time{}, ErrNoExpiry
	}
	if err != nil {
		return time.Time{}, fmt.Errorf("ttl: read: %w", err)
	}
	var e entry
	if err := json.Unmarshal(data, &e); err != nil {
		return time.Time{}, fmt.Errorf("ttl: unmarshal: %w", err)
	}
	return e.ExpiresAt, nil
}

// Check returns ErrNoExpiry if no expiry is set, ErrExpired if the
// scope has expired, or nil if the scope is still valid.
func Check(storeDir, scope string) error {
	t, err := Get(storeDir, scope)
	if err != nil {
		return err
	}
	if time.Now().UTC().After(t) {
		return ErrExpired
	}
	return nil
}

// Clear removes the expiry record for a scope.
// Returns nil if no record existed.
func Clear(storeDir, scope string) error {
	err := os.Remove(indexPath(storeDir, scope))
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("ttl: remove: %w", err)
	}
	return nil
}
