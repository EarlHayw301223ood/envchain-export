package rotate_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/nicholasgasior/envchain-export/internal/chain"
	"github.com/nicholasgasior/envchain-export/internal/rotate"
	"github.com/nicholasgasior/envchain-export/internal/store"
)

func makeStore(t *testing.T) (string, *store.Store) {
	t.Helper()
	dir := t.TempDir()
	return dir, store.New(dir)
}

func seedScope(t *testing.T, s *store.Store, scope, pass string, pairs map[string]string) {
	t.Helper()
	ch, err := chain.New(scope)
	if err != nil {
		t.Fatalf("chain.New: %v", err)
	}
	for k, v := range pairs {
		if err := ch.Add(k, v); err != nil {
			t.Fatalf("ch.Add: %v", err)
		}
	}
	if err := s.Save(scope, pass, ch); err != nil {
		t.Fatalf("store.Save: %v", err)
	}
}

func TestRotate_Success(t *testing.T) {
	dir, s := makeStore(t)
	seedScope(t, s, "prod", "old-secret", map[string]string{"FOO": "bar", "BAZ": "qux"})

	if err := rotate.Rotate(dir, "prod", "old-secret", "new-secret"); err != nil {
		t.Fatalf("Rotate: %v", err)
	}

	// Must be loadable with new passphrase.
	ch, err := s.Load("prod", "new-secret")
	if err != nil {
		t.Fatalf("Load after rotate: %v", err)
	}
	if v, _ := ch.Get("FOO"); v != "bar" {
		t.Errorf("FOO = %q, want %q", v, "bar")
	}
}

func TestRotate_WrongOldPassphrase(t *testing.T) {
	dir, s := makeStore(t)
	seedScope(t, s, "prod", "correct", map[string]string{"X": "1"})

	if err := rotate.Rotate(dir, "prod", "wrong", "new-secret"); err == nil {
		t.Fatal("expected error for wrong old passphrase, got nil")
	}

	// Original file must still be intact.
	ch, err := s.Load("prod", "correct")
	if err != nil {
		t.Fatalf("Load with original pass after failed rotate: %v", err)
	}
	if v, _ := ch.Get("X"); v != "1" {
		t.Errorf("X = %q, want %q", v, "1")
	}
}

func TestRotate_SamePassphrase(t *testing.T) {
	dir, _ := makeStore(t)
	if err := rotate.Rotate(dir, "prod", "same", "same"); err == nil {
		t.Fatal("expected ErrSamePassphrase, got nil")
	}
}

func TestRotate_MissingScope(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "noexist")
	_ = os.MkdirAll(dir, 0o700)
	if err := rotate.Rotate(dir, "ghost", "old", "new"); err == nil {
		t.Fatal("expected error for missing scope, got nil")
	}
}
