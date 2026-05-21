// Package snapshot provides functionality to capture and restore
// point-in-time copies of a scope's environment variables.
package snapshot

import (
	"fmt"
	"time"

	"github.com/nicholasgasior/envchain-export/internal/chain"
	"github.com/nicholasgasior/envchain-export/internal/store"
)

// ErrNoSnapshots is returned when no snapshots exist for a scope.
var ErrNoSnapshots = fmt.Errorf("no snapshots found for scope")

// snapshotScopeName returns the internal scope name used to store a snapshot.
func snapshotScopeName(scope, label string) string {
	return fmt.Sprintf("__snapshot__%s__%s", scope, label)
}

// Take captures the current state of the given scope and saves it as a
// snapshot identified by a timestamp label.
func Take(st *store.Store, scope, passphrase string) (string, error) {
	ch, err := st.Load(scope, passphrase)
	if err != nil {
		return "", fmt.Errorf("load scope: %w", err)
	}

	label := time.Now().UTC().Format("20060102T150405Z")
	snapScope := snapshotScopeName(scope, label)

	snap, err := chain.New(snapScope)
	if err != nil {
		return "", fmt.Errorf("create snapshot chain: %w", err)
	}

	for _, k := range ch.Keys() {
		v, _ := ch.Get(k)
		if err := snap.Add(k, v); err != nil {
			return "", fmt.Errorf("copy key %q: %w", k, err)
		}
	}

	if err := st.Save(snapScope, snap, passphrase); err != nil {
		return "", fmt.Errorf("save snapshot: %w", err)
	}

	return label, nil
}

// Restore loads a previously taken snapshot back into the live scope,
// overwriting its current contents.
func Restore(st *store.Store, scope, label, passphrase string) error {
	snapScope := snapshotScopeName(scope, label)

	snap, err := st.Load(snapScope, passphrase)
	if err != nil {
		return fmt.Errorf("load snapshot %q: %w", label, err)
	}

	ch, err := chain.New(scope)
	if err != nil {
		return fmt.Errorf("create chain: %w", err)
	}

	for _, k := range snap.Keys() {
		v, _ := snap.Get(k)
		if err := ch.Add(k, v); err != nil {
			return fmt.Errorf("restore key %q: %w", k, err)
		}
	}

	if err := st.Save(scope, ch, passphrase); err != nil {
		return fmt.Errorf("save restored scope: %w", err)
	}

	return nil
}

// List returns all snapshot labels for the given scope, sorted oldest first.
func List(st *store.Store, scope string) ([]string, error) {
	prefix := fmt.Sprintf("__snapshot__%s__", scope)
	all, err := st.ListScopes()
	if err != nil {
		return nil, err
	}

	var labels []string
	for _, s := range all {
		if len(s) > len(prefix) && s[:len(prefix)] == prefix {
			labels = append(labels, s[len(prefix):])
		}
	}

	if len(labels) == 0 {
		return nil, ErrNoSnapshots
	}

	return labels, nil
}
