package rename_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/envchain-export/internal/chain"
	"github.com/envchain-export/internal/rename"
	"github.com/envchain-export/internal/store"
)

const testPass = "hunter2"

func makeStore(t *testing.T) *store.Store {
	t.Helper()
	dir := filepath.Join(t.TempDir(), "envchain")
	st, err := store.New(dir)
	if err != nil {
		t.Fatalf("store.New: %v", err)
	}
	return st
}

func seedScope(t *testing.T, st *store.Store, scope string) {
	t.Helper()
	ch, _ := chain.New(scope)
	ch.Add("KEY", "value")
	if err := st.Save(scope, testPass, ch); err != nil {
		t.Fatalf("seed %q: %v", scope, err)
	}
}

func TestRename_Success(t *testing.T) {
	st := makeStore(t)
	seedScope(t, st, "old-scope")

	if err := rename.Rename(st, "old-scope", "new-scope", testPass); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if st.Exists("old-scope") {
		t.Error("old scope still exists after rename")
	}
	if !st.Exists("new-scope") {
		t.Error("new scope does not exist after rename")
	}

	ch, err := st.Load("new-scope", testPass)
	if err != nil {
		t.Fatalf("load new scope: %v", err)
	}
	if v, _ := ch.Get("KEY"); v != "value" {
		t.Errorf("expected KEY=value, got %q", v)
	}
}

func TestRename_SameName_ReturnsError(t *testing.T) {
	st := makeStore(t)
	seedScope(t, st, "scope")

	err := rename.Rename(st, "scope", "scope", testPass)
	if err == nil {
		t.Fatal("expected error for same-name rename, got nil")
	}
}

func TestRename_DestExists_ReturnsError(t *testing.T) {
	st := makeStore(t)
	seedScope(t, st, "src")
	seedScope(t, st, "dst")

	err := rename.Rename(st, "src", "dst", testPass)
	if err == nil {
		t.Fatal("expected error when destination exists, got nil")
	}
}

func TestRename_WrongPassphrase_ReturnsError(t *testing.T) {
	st := makeStore(t)
	seedScope(t, st, "src")

	err := rename.Rename(st, "src", "dst", "wrong-passphrase")
	if err == nil {
		t.Fatal("expected error for wrong passphrase, got nil")
	}
	// Destination must NOT have been created on failure.
	if st.Exists("dst") {
		t.Error("destination scope should not exist after failed rename")
	}
}

func TestRename_MissingSource_ReturnsError(t *testing.T) {
	st := makeStore(t)

	err := rename.Rename(st, "ghost", "new", testPass)
	if err == nil {
		t.Fatal("expected error for missing source scope, got nil")
	}
}

func init() {
	// Ensure test binary can locate store files created in temp dirs.
	_ = os.Getenv
	_ = filepath.Join
}
