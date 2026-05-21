package copy_test

import (
	"os"
	"testing"

	"github.com/user/envchain-export/internal/chain"
	"github.com/user/envchain-export/internal/copy"
	"github.com/user/envchain-export/internal/store"
)

func makeStore(t *testing.T) *store.Store {
	t.Helper()
	dir := t.TempDir()
	st, err := store.New(dir)
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
			t.Fatalf("chain.Add(%q): %v", k, err)
		}
	}
	if err := st.Save(scope, passphrase, ch); err != nil {
		t.Fatalf("st.Save: %v", err)
	}
}

func TestCopy_Success(t *testing.T) {
	st := makeStore(t)
	seedScope(t, st, "src", "pass", map[string]string{"FOO": "bar", "BAZ": "qux"})

	if err := copy.Copy(st, "src", "dst", "pass", ""); err != nil {
		t.Fatalf("Copy: %v", err)
	}

	ch, err := st.Load("dst", "pass")
	if err != nil {
		t.Fatalf("Load dst: %v", err)
	}
	if v, _ := ch.Get("FOO"); v != "bar" {
		t.Errorf("expected FOO=bar, got %q", v)
	}
	if v, _ := ch.Get("BAZ"); v != "qux" {
		t.Errorf("expected BAZ=qux, got %q", v)
	}
}

func TestCopy_DifferentPassphrase(t *testing.T) {
	st := makeStore(t)
	seedScope(t, st, "src", "oldpass", map[string]string{"KEY": "val"})

	if err := copy.Copy(st, "src", "dst", "oldpass", "newpass"); err != nil {
		t.Fatalf("Copy: %v", err)
	}

	if _, err := st.Load("dst", "oldpass"); err == nil {
		t.Error("expected error loading dst with old passphrase")
	}
	ch, err := st.Load("dst", "newpass")
	if err != nil {
		t.Fatalf("Load dst with new pass: %v", err)
	}
	if v, _ := ch.Get("KEY"); v != "val" {
		t.Errorf("expected KEY=val, got %q", v)
	}
}

func TestCopy_SameName_ReturnsError(t *testing.T) {
	st := makeStore(t)
	seedScope(t, st, "src", "pass", map[string]string{"A": "1"})

	if err := copy.Copy(st, "src", "src", "pass", ""); !os.IsNotExist(err) && err != copy.ErrSameName {
		if err != copy.ErrSameName {
			t.Errorf("expected ErrSameName, got %v", err)
		}
	}
}

func TestCopy_WrongPassphrase_ReturnsError(t *testing.T) {
	st := makeStore(t)
	seedScope(t, st, "src", "correct", map[string]string{"X": "y"})

	if err := copy.Copy(st, "src", "dst", "wrong", ""); err == nil {
		t.Error("expected error for wrong passphrase, got nil")
	}
}

func TestCopy_MissingSource_ReturnsError(t *testing.T) {
	st := makeStore(t)

	if err := copy.Copy(st, "nonexistent", "dst", "pass", ""); err == nil {
		t.Error("expected error for missing source scope, got nil")
	}
}
