// Package pin provides functionality to mark scopes as pinned,
// preventing accidental deletion or overwrite.
package pin

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

var ErrAlreadyPinned = errors.New("scope is already pinned")
var ErrNotPinned = errors.New("scope is not pinned")

const indexFile = "pins.json"

func indexPath(storeDir string) string {
	return filepath.Join(storeDir, indexFile)
}

func loadIndex(storeDir string) (map[string]bool, error) {
	path := indexPath(storeDir)
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return map[string]bool{}, nil
	}
	if err != nil {
		return nil, err
	}
	var index map[string]bool
	if err := json.Unmarshal(data, &index); err != nil {
		return nil, err
	}
	return index, nil
}

func saveIndex(storeDir string, index map[string]bool) error {
	data, err := json.MarshalIndent(index, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(indexPath(storeDir), data, 0600)
}

// Add marks the given scope as pinned.
func Add(storeDir, scope string) error {
	index, err := loadIndex(storeDir)
	if err != nil {
		return err
	}
	if index[scope] {
		return ErrAlreadyPinned
	}
	index[scope] = true
	return saveIndex(storeDir, index)
}

// Remove unmarks the given scope as pinned.
func Remove(storeDir, scope string) error {
	index, err := loadIndex(storeDir)
	if err != nil {
		return err
	}
	if !index[scope] {
		return ErrNotPinned
	}
	delete(index, scope)
	return saveIndex(storeDir, index)
}

// IsPinned reports whether the given scope is pinned.
func IsPinned(storeDir, scope string) (bool, error) {
	index, err := loadIndex(storeDir)
	if err != nil {
		return false, err
	}
	return index[scope], nil
}

// List returns all currently pinned scope names.
func List(storeDir string) ([]string, error) {
	index, err := loadIndex(storeDir)
	if err != nil {
		return nil, err
	}
	var scopes []string
	for scope := range index {
		scopes = append(scopes, scope)
	}
	return scopes, nil
}
