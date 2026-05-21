package snapshot_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/nicholasgasior/envchain-export/internal/chain"
	"github.com/nicholasgasior/envchain-export/internal/snapshot"
	"github.com/nicholasgasior/envchain-export/internal/store"
)

func makeStore(t *testing.T) *store.Store {
	t.Helper()
	dir := t.TempDir()
	st, err := store.New(filepath.Join(dir, "store"))
	if err != nil {
		t.Fatalf("store.New: %v", err)
	}
	return st
}

func seedScope(t *testing.T, st *store.Store, scope, passphrase string, pairs map[string]string) {
	t.Helper()
	ch, err := chain.New(scope)
	if err != nil {
		t.Fatalf("chain.New: %v", err)
	}
	for k, v := range pairs {
		if err := ch.Add(k, v); err != nil {
			t.Fatalf("ch.Add(%q): %v", k, err)
		}
	}
	if err := st.Save(scope, ch, passphrase); err != nil {
		t.Fatalf("st.Save: %v", err)
	}
}

func TestTake_CreatesSnapshot(t *testing.T) {
	st := makeStore(t)
	seedScope(t, st, "myapp", "secret", map[string]string{"FOO": "bar", "BAZ": "qux"})

	label, err := snapshot.Take(st, "myapp", "secret")
	if err != nil {
		t.Fatalf("Take: %v", err)
	}
	if label == "" {
		t.Fatal("expected non-empty label")
	}

	labels, err := snapshot.List(st, "myapp")
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(labels) != 1 || labels[0] != label {
		t.Fatalf("expected [%q], got %v", label, labels)
	}
}

func TestRestore_RoundTrip(t *testing.T) {
	st := makeStore(t)
	seedScope(t, st, "myapp", "secret", map[string]string{"FOO": "original"})

	label, err := snapshot.Take(st, "myapp", "secret")
	if err != nil {
		t.Fatalf("Take: %v", err)
	}

	// Mutate the live scope.
	seedScope(t, st, "myapp", "secret", map[string]string{"FOO": "mutated"})

	if err := snapshot.Restore(st, "myapp", label, "secret"); err != nil {
		t.Fatalf("Restore: %v", err)
	}

	ch, err := st.Load("myapp", "secret")
	if err != nil {
		t.Fatalf("Load after restore: %v", err)
	}
	v, ok := ch.Get("FOO")
	if !ok || v != "original" {
		t.Fatalf("expected FOO=original, got %q (ok=%v)", v, ok)
	}
}

func TestList_NoSnapshots_ReturnsError(t *testing.T) {
	st := makeStore(t)
	seedScope(t, st, "myapp", "secret", map[string]string{"X": "1"})

	_, err := snapshot.List(st, "myapp")
	if err != snapshot.ErrNoSnapshots {
		t.Fatalf("expected ErrNoSnapshots, got %v", err)
	}
}

func TestTake_NonExistentScope_ReturnsError(t *testing.T) {
	st := makeStore(t)
	// Ensure store dir exists but scope does not.
	if err := os.MkdirAll(st.Dir(), 0700); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	_, err := snapshot.Take(st, "ghost", "secret")
	if err == nil {
		t.Fatal("expected error for non-existent scope")
	}
}
