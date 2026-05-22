package touch_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/envchain-export/internal/chain"
	"github.com/envchain-export/internal/store"
	"github.com/envchain-export/internal/touch"
)

func makeStore(t *testing.T) *store.Store {
	t.Helper()
	dir := filepath.Join(t.TempDir(), "envchain")
	if err := os.MkdirAll(dir, 0700); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	return store.New(dir)
}

func seedScope(t *testing.T, st *store.Store, scope, passphrase string) {
	t.Helper()
	ch, _ := chain.New(scope)
	ch.Add("KEY", "value")
	if err := st.Save(scope, ch, passphrase); err != nil {
		t.Fatalf("seed: %v", err)
	}
}

func TestTouch_Success(t *testing.T) {
	st := makeStore(t)
	seedScope(t, st, "prod", "old-pass")

	if err := touch.Touch(st, "prod", "old-pass", "new-pass"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Should be loadable with new passphrase.
	ch, err := st.Load("prod", "new-pass")
	if err != nil {
		t.Fatalf("load with new passphrase: %v", err)
	}
	if v, _ := ch.Get("KEY"); v != "value" {
		t.Errorf("expected KEY=value, got %q", v)
	}
}

func TestTouch_WrongOldPassphrase(t *testing.T) {
	st := makeStore(t)
	seedScope(t, st, "prod", "correct")

	err := touch.Touch(st, "prod", "wrong", "new-pass")
	if err == nil {
		t.Fatal("expected error for wrong passphrase")
	}
}

func TestTouch_SamePassphrase(t *testing.T) {
	st := makeStore(t)
	seedScope(t, st, "prod", "same")

	err := touch.Touch(st, "prod", "same", "same")
	if err != touch.ErrSamePassphrase {
		t.Errorf("expected ErrSamePassphrase, got %v", err)
	}
}

func TestTouch_NonExistentScope(t *testing.T) {
	st := makeStore(t)

	err := touch.Touch(st, "ghost", "old", "new")
	if err == nil {
		t.Fatal("expected error for non-existent scope")
	}
}
