// Package tag provides functionality for tagging scopes with arbitrary
// labels, enabling grouping and filtering of envchain scopes.
package tag

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

// ErrNoTags is returned when a scope has no tags.
var ErrNoTags = errors.New("tag: no tags found")

// ErrTagNotFound is returned when a specific tag does not exist on a scope.
var ErrTagNotFound = errors.New("tag: tag not found")

const tagFile = "tags.json"

// tagIndex maps scope names to their list of tags.
type tagIndex map[string][]string

func indexPath(storeDir string) string {
	return filepath.Join(storeDir, tagFile)
}

func loadIndex(storeDir string) (tagIndex, error) {
	path := indexPath(storeDir)
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return make(tagIndex), nil
	}
	if err != nil {
		return nil, fmt.Errorf("tag: read index: %w", err)
	}
	var idx tagIndex
	if err := json.Unmarshal(data, &idx); err != nil {
		return nil, fmt.Errorf("tag: parse index: %w", err)
	}
	return idx, nil
}

func saveIndex(storeDir string, idx tagIndex) error {
	data, err := json.MarshalIndent(idx, "", "  ")
	if err != nil {
		return fmt.Errorf("tag: marshal index: %w", err)
	}
	if err := os.WriteFile(indexPath(storeDir), data, 0600); err != nil {
		return fmt.Errorf("tag: write index: %w", err)
	}
	return nil
}

// Add adds a tag to the given scope. Duplicate tags are ignored.
func Add(storeDir, scope, tag string) error {
	idx, err := loadIndex(storeDir)
	if err != nil {
		return err
	}
	for _, t := range idx[scope] {
		if t == tag {
			return nil
		}
	}
	idx[scope] = append(idx[scope], tag)
	sort.Strings(idx[scope])
	return saveIndex(storeDir, idx)
}

// Remove removes a tag from the given scope.
func Remove(storeDir, scope, tag string) error {
	idx, err := loadIndex(storeDir)
	if err != nil {
		return err
	}
	tags := idx[scope]
	filtered := tags[:0]
	found := false
	for _, t := range tags {
		if t == tag {
			found = true
			continue
		}
		filtered = append(filtered, t)
	}
	if !found {
		return ErrTagNotFound
	}
	idx[scope] = filtered
	return saveIndex(storeDir, idx)
}

// List returns all tags for the given scope.
func List(storeDir, scope string) ([]string, error) {
	idx, err := loadIndex(storeDir)
	if err != nil {
		return nil, err
	}
	tags, ok := idx[scope]
	if !ok || len(tags) == 0 {
		return nil, ErrNoTags
	}
	return tags, nil
}

// ScopesByTag returns all scopes that carry the given tag.
func ScopesByTag(storeDir, tag string) ([]string, error) {
	idx, err := loadIndex(storeDir)
	if err != nil {
		return nil, err
	}
	var scopes []string
	for scope, tags := range idx {
		for _, t := range tags {
			if t == tag {
				scopes = append(scopes, scope)
				break
			}
		}
	}
	sort.Strings(scopes)
	return scopes, nil
}
