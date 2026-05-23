// Package alias provides support for creating and resolving short aliases
// that point to named scopes within the envchain store.
package alias

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

var (
	ErrAliasNotFound  = errors.New("alias not found")
	ErrAliasExists    = errors.New("alias already exists")
	ErrCircularAlias  = errors.New("alias cannot point to another alias")
)

func indexPath(storeDir string) string {
	return filepath.Join(storeDir, "aliases.json")
}

func loadIndex(storeDir string) (map[string]string, error) {
	path := indexPath(storeDir)
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return map[string]string{}, nil
	}
	if err != nil {
		return nil, err
	}
	var idx map[string]string
	if err := json.Unmarshal(data, &idx); err != nil {
		return nil, err
	}
	return idx, nil
}

func saveIndex(storeDir string, idx map[string]string) error {
	data, err := json.MarshalIndent(idx, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(indexPath(storeDir), data, 0o600)
}

// Add creates a new alias pointing to target scope.
func Add(storeDir, alias, target string) error {
	idx, err := loadIndex(storeDir)
	if err != nil {
		return err
	}
	if _, exists := idx[alias]; exists {
		return fmt.Errorf("%w: %s", ErrAliasExists, alias)
	}
	// Prevent chaining: target must not itself be an alias.
	if _, targetIsAlias := idx[target]; targetIsAlias {
		return ErrCircularAlias
	}
	idx[alias] = target
	return saveIndex(storeDir, idx)
}

// Remove deletes an existing alias.
func Remove(storeDir, alias string) error {
	idx, err := loadIndex(storeDir)
	if err != nil {
		return err
	}
	if _, exists := idx[alias]; !exists {
		return fmt.Errorf("%w: %s", ErrAliasNotFound, alias)
	}
	delete(idx, alias)
	return saveIndex(storeDir, idx)
}

// Resolve returns the scope name that the alias points to.
func Resolve(storeDir, alias string) (string, error) {
	idx, err := loadIndex(storeDir)
	if err != nil {
		return "", err
	}
	target, exists := idx[alias]
	if !exists {
		return "", fmt.Errorf("%w: %s", ErrAliasNotFound, alias)
	}
	return target, nil
}

// List returns all aliases sorted alphabetically.
func List(storeDir string) (map[string]string, error) {
	idx, err := loadIndex(storeDir)
	if err != nil {
		return nil, err
	}
	if len(idx) == 0 {
		return idx, nil
	}
	// Return a stable copy.
	out := make(map[string]string, len(idx))
	keys := make([]string, 0, len(idx))
	for k := range idx {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		out[k] = idx[k]
	}
	return out, nil
}
