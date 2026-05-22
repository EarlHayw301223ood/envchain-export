package clone_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/nicholasgasior/envchain-export/internal/chain"
	"github.com/nicholasgasior/envchain-export/internal/clone"
	"github.com/nicholasgasior/envchain-export/internal/store"
)

func makeStore(t *testing.T) *store.Store {
	t.Helper()
	dir := filepath.Join(t.TempDir(), ".envchain")
	st, err := store.New(dir)
	if err != nil {
		t.Fatalf("store.New: %v", err)
	}
	return st
}

func seedScope(t *testing.T, st *store.Store, scope, pass string) {
	t.Helper()
	ch, err := chain.New(scope)
	if err != nil {
		t.Fatalf("chain.New: %v", err)
	}
	if err := ch.Add("KEY1", "value1"); err != nil {
		t.Fatalf("chain.Add: %v", err)
	}
	if err := ch.Add("KEY2", "value2"); err != nil {
		t.Fatalf("chain.Add: %v", err)
	}
	if err := st.Save(scope, pass, ch); err != nil {
		t.Fatalf("st.Save: %v", err)
	}
}

func TestClone_Success(t *testing.T) {
	st := makeStore(t)
	seedScope(t, st, "src", "pass")

	if err := clone.Clone(st, "src", "dst", "pass", "pass"); err != nil {
		t.Fatalf("Clone: %v", err)
	}

	ch, err := st.Load("dst", "pass")
	if err != nil {
		t.Fatalf("Load dst: %v", err)
	}
	v, _ := ch.Get("KEY1")
	if v != "value1" {
		t.Errorf("KEY1 = %q, want %q", v, "value1")
	}
}

func TestClone_DifferentPassphrase(t *testing.T) {
	st := makeStore(t)
	seedScope(t, st, "src", "oldpass")

	if err := clone.Clone(st, "src", "dst", "oldpass", "newpass"); err != nil {
		t.Fatalf("Clone: %v", err)
	}

	if _, err := st.Load("dst", "newpass"); err != nil {
		t.Errorf("expected dst loadable with newpass, got: %v", err)
	}
}

func TestClone_SameName_ReturnsError(t *testing.T) {
	st := makeStore(t)
	seedScope(t, st, "src", "pass")

	if err := clone.Clone(st, "src", "src", "pass", "pass"); err == nil {
		t.Error("expected error for same src/dst, got nil")
	}
}

func TestClone_DestExists_ReturnsError(t *testing.T) {
	st := makeStore(t)
	seedScope(t, st, "src", "pass")
	seedScope(t, st, "dst", "pass")

	if err := clone.Clone(st, "src", "dst", "pass", "pass"); err == nil {
		t.Error("expected error when dst already exists, got nil")
	}
}

func TestClone_WrongPassphrase_ReturnsError(t *testing.T) {
	st := makeStore(t)
	seedScope(t, st, "src", "pass")

	if err := clone.Clone(st, "src", "dst", "wrong", "pass"); err == nil {
		t.Error("expected error for wrong passphrase, got nil")
	}
}

func TestClone_MissingSource_ReturnsError(t *testing.T) {
	st := makeStore(t)

	if err := clone.Clone(st, "ghost", "dst", "pass", "pass"); err == nil {
		t.Error("expected error for missing source, got nil")
	}
	_ = os.Getenv("HOME") // keep import happy
}
