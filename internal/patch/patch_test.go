package patch_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/envchain-export/internal/chain"
	"github.com/user/envchain-export/internal/patch"
	"github.com/user/envchain-export/internal/store"
)

func makeStore(t *testing.T) *store.Store {
	t.Helper()
	dir := filepath.Join(t.TempDir(), "envchain")
	st, err := store.New(dir)
	if err != nil {
		t.Fatalf("store.New: %v", err)
	}
	return st
}

func seedScope(t *testing.T, st *store.Store, scope, pass string, pairs map[string]string) {
	t.Helper()
	ch, _ := chain.New(scope)
	for k, v := range pairs {
		if err := ch.Add(k, v); err != nil {
			t.Fatalf("seed add: %v", err)
		}
	}
	if err := st.Save(scope, pass, ch); err != nil {
		t.Fatalf("seed save: %v", err)
	}
}

func TestPatch_SetNewKey(t *testing.T) {
	st := makeStore(t)
	seedScope(t, st, "app", "secret", map[string]string{"FOO": "bar"})

	ops := patch.BuildOps([][2]string{{"NEW_KEY", "new_val"}}, nil)
	if err := patch.Patch(st, "app", "secret", ops); err != nil {
		t.Fatalf("Patch: %v", err)
	}

	ch, _ := st.Load("app", "secret")
	if v, _ := ch.Get("NEW_KEY"); v != "new_val" {
		t.Errorf("expected new_val, got %q", v)
	}
	if v, _ := ch.Get("FOO"); v != "bar" {
		t.Errorf("FOO should be unchanged, got %q", v)
	}
}

func TestPatch_UnsetKey(t *testing.T) {
	st := makeStore(t)
	seedScope(t, st, "app", "secret", map[string]string{"FOO": "bar", "BAZ": "qux"})

	ops := patch.BuildOps(nil, []string{"FOO"})
	if err := patch.Patch(st, "app", "secret", ops); err != nil {
		t.Fatalf("Patch: %v", err)
	}

	ch, _ := st.Load("app", "secret")
	if _, err := ch.Get("FOO"); err == nil {
		t.Error("expected FOO to be removed")
	}
	if v, _ := ch.Get("BAZ"); v != "qux" {
		t.Errorf("BAZ should be unchanged, got %q", v)
	}
}

func TestPatch_NoOps_ReturnsError(t *testing.T) {
	st := makeStore(t)
	seedScope(t, st, "app", "secret", map[string]string{"A": "1"})

	if err := patch.Patch(st, "app", "secret", nil); err == nil {
		t.Error("expected ErrNoOps")
	}
}

func TestPatch_InvalidKey_ReturnsError(t *testing.T) {
	st := makeStore(t)
	seedScope(t, st, "app", "secret", map[string]string{"A": "1"})

	ops := patch.BuildOps([][2]string{{"123INVALID", "v"}}, nil)
	if err := patch.Patch(st, "app", "secret", ops); err == nil {
		t.Error("expected error for invalid key")
	}
}

func TestPatch_WrongPassphrase_ReturnsError(t *testing.T) {
	st := makeStore(t)
	seedScope(t, st, "app", "correct", map[string]string{"A": "1"})

	ops := patch.BuildOps([][2]string{{"B", "2"}}, nil)
	if err := patch.Patch(st, "app", "wrong", ops); err == nil {
		t.Error("expected error for wrong passphrase")
	}
}

func TestPatch_MissingScope_ReturnsError(t *testing.T) {
	st := makeStore(t)
	_ = os.MkdirAll(filepath.Join(t.TempDir(), "envchain"), 0o700)

	ops := patch.BuildOps([][2]string{{"A", "1"}}, nil)
	if err := patch.Patch(st, "ghost", "pass", ops); err == nil {
		t.Error("expected error for missing scope")
	}
}
