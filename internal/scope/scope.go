// Package scope provides utilities for listing, validating, and managing
// named scopes stored on disk by the envchain-export store.
package scope

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/user/envchain-export/internal/validate"
)

// ErrNoScopes is returned when the store directory contains no scopes.
var ErrNoScopes = errors.New("scope: no scopes found")

// List returns all scope names persisted under storeDir.
// Each scope corresponds to a file named "<scope>.enc" in the directory.
func List(storeDir string) ([]string, error) {
	entries, err := os.ReadDir(storeDir)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, ErrNoScopes
		}
		return nil, err
	}

	var scopes []string
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		if name := e.Name(); strings.HasSuffix(name, ".enc") {
			scopes = append(scopes, strings.TrimSuffix(name, ".enc"))
		}
	}

	if len(scopes) == 0 {
		return nil, ErrNoScopes
	}
	return scopes, nil
}

// Exists reports whether the given scope name has a corresponding file in
// storeDir.
func Exists(storeDir, name string) (bool, error) {
	if err := validate.Key(name); err != nil {
		return false, err
	}
	path := filepath.Join(storeDir, name+".enc")
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	return false, err
}

// Delete removes the encrypted file for the given scope from storeDir.
func Delete(storeDir, name string) error {
	if err := validate.Key(name); err != nil {
		return err
	}
	path := filepath.Join(storeDir, name+".enc")
	if err := os.Remove(path); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return errors.New("scope: " + name + " does not exist")
		}
		return err
	}
	return nil
}
