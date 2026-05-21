// Package diff computes the difference between two chains, reporting
// keys that were added, removed, or changed between a base and a target.
package diff

import (
	"sort"

	"github.com/user/envchain-export/internal/chain"
)

// ChangeKind describes the type of change for a key.
type ChangeKind string

const (
	Added   ChangeKind = "added"
	Removed ChangeKind = "removed"
	Changed ChangeKind = "changed"
)

// Entry represents a single diff entry.
type Entry struct {
	Key      string
	Kind     ChangeKind
	OldValue string // empty when Kind == Added
	NewValue string // empty when Kind == Removed
}

// Result holds all diff entries between two chains.
type Result struct {
	Entries []Entry
}

// HasChanges returns true if there is at least one diff entry.
func (r *Result) HasChanges() bool {
	return len(r.Entries) > 0
}

// Diff computes the difference between base and target chains.
// Keys present only in target are Added; keys missing from target are Removed;
// keys in both with different values are Changed.
func Diff(base, target *chain.Chain) *Result {
	result := &Result{}

	baseKeys := base.Keys()
	targetKeys := target.Keys()

	baseMap := make(map[string]string, len(baseKeys))
	for _, k := range baseKeys {
		v, _ := base.Get(k)
		baseMap[k] = v
	}

	targetMap := make(map[string]string, len(targetKeys))
	for _, k := range targetKeys {
		v, _ := target.Get(k)
		targetMap[k] = v
	}

	// Detect added and changed
	for _, k := range targetKeys {
		oldVal, exists := baseMap[k]
		newVal := targetMap[k]
		if !exists {
			result.Entries = append(result.Entries, Entry{Key: k, Kind: Added, NewValue: newVal})
		} else if oldVal != newVal {
			result.Entries = append(result.Entries, Entry{Key: k, Kind: Changed, OldValue: oldVal, NewValue: newVal})
		}
	}

	// Detect removed
	for _, k := range baseKeys {
		if _, exists := targetMap[k]; !exists {
			v := baseMap[k]
			result.Entries = append(result.Entries, Entry{Key: k, Kind: Removed, OldValue: v})
		}
	}

	sort.Slice(result.Entries, func(i, j int) bool {
		return result.Entries[i].Key < result.Entries[j].Key
	})

	return result
}
