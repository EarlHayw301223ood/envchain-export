package main

import (
	"bytes"
	"path/filepath"
	"strings"
	"testing"

	"github.com/nicholasgasior/envchain-export/internal/chain"
	"github.com/nicholasgasior/envchain-export/internal/store"
)

func seedSnapshotScope(t *testing.T, st *store.Store, scope, passphrase string) {
	t.Helper()
	ch, err := chain.New(scope)
	if err != nil {
		t.Fatalf("chain.New: %v", err)
	}
	_ = ch.Add("KEY", "value")
	if err := st.Save(scope, ch, passphrase); err != nil {
		t.Fatalf("st.Save: %v", err)
	}
}

func setupSnapshotStore(t *testing.T) *store.Store {
	t.Helper()
	st, err := store.New(filepath.Join(t.TempDir(), "store"))
	if err != nil {
		t.Fatalf("store.New: %v", err)
	}
	return st
}

func TestRunSnapshotList_NoSnapshots_ReturnsError(t *testing.T) {
	st := setupSnapshotStore(t)
	seedSnapshotScope(t, st, "myapp", "pass")

	var buf bytes.Buffer
	err := runSnapshotList(st, &buf, "myapp")
	if err == nil {
		t.Fatal("expected error when no snapshots exist")
	}
}

func TestRunSnapshotTake_OutputsLabel(t *testing.T) {
	st := setupSnapshotStore(t)
	seedSnapshotScope(t, st, "myapp", "pass")

	var buf bytes.Buffer
	// Bypass passphrase prompt by calling the inner logic directly.
	import_snapshot := func() (string, error) {
		return snapshotTakeDirect(st, "myapp", "pass")
	}
	label, err := import_snapshot()
	if err != nil {
		t.Fatalf("Take: %v", err)
	}

	_ = runSnapshotListDirect(st, &buf, "myapp")
	if !strings.Contains(buf.String(), label) {
		t.Fatalf("expected label %q in output %q", label, buf.String())
	}
}

func TestRunSnapshotRestore_RestoredValues(t *testing.T) {
	st := setupSnapshotStore(t)
	seedSnapshotScope(t, st, "myapp", "pass")

	label, err := snapshotTakeDirect(st, "myapp", "pass")
	if err != nil {
		t.Fatalf("Take: %v", err)
	}

	// Overwrite the scope.
	seedSnapshotScope(t, st, "myapp", "pass")

	var buf bytes.Buffer
	if err := snapshotRestoreDirect(st, &buf, "myapp", label, "pass"); err != nil {
		t.Fatalf("Restore: %v", err)
	}
	if !strings.Contains(buf.String(), "restored") {
		t.Fatalf("expected 'restored' in output, got %q", buf.String())
	}
}
