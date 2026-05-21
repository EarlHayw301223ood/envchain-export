package inspect_test

import {
	"bytes"
	"path/filepath"
	"strings"
	"testing"

	"github.com/user/envchain-export/internal/chain"
	"github.com/user/envchain-export/internal/inspect"
	"github.com/user/envchain-export/internal/store"
)

const testPass = "hunter2"

func makeStore(t *testing.T) *store.Store {
	t.Helper()
	st, err := store.New(filepath.Join(t.TempDir(), "store"))
	if err != nil {
		t.Fatalf("store.New: %v", err)
	}
	return st
}

func seedScope(t *testing.T, st *store.Store, scope string, pairs map[string]string) {
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
	if err := st.Save(ch, testPass); err != nil {
		t.Fatalf("st.Save: %v", err)
	}
}

func TestInspect_ShowsKeysAndValues(t *testing.T) {
	st := makeStore(t)
	seedScope(t, st, "myapp", map[string]string{"API_KEY": "secret123", "PORT": "8080"})

	var buf bytes.Buffer
	err := inspect.Inspect(st, "myapp", testPass, inspect.Options{SortKeys: true}, &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "API_KEY") {
		t.Errorf("expected API_KEY in output, got:\n%s", out)
	}
	if !strings.Contains(out, "secret123") {
		t.Errorf("expected value in output, got:\n%s", out)
	}
}

func TestInspect_MaskValues(t *testing.T) {
	st := makeStore(t)
	seedScope(t, st, "prod", map[string]string{"DB_PASS": "supersecret"})

	var buf bytes.Buffer
	err := inspect.Inspect(st, "prod", testPass, inspect.Options{MaskValues: true}, &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if strings.Contains(out, "supersecret") {
		t.Errorf("value should be masked, got:\n%s", out)
	}
	if !strings.Contains(out, "su*********") {
		t.Errorf("expected masked value su********* in output, got:\n%s", out)
	}
}

func TestInspect_WrongPassphrase(t *testing.T) {
	st := makeStore(t)
	seedScope(t, st, "secure", map[string]string{"TOKEN": "abc"})

	var buf bytes.Buffer
	err := inspect.Inspect(st, "secure", "wrongpass", inspect.Options{}, &buf)
	if err == nil {
		t.Fatal("expected error for wrong passphrase, got nil")
	}
}

func TestInspect_EmptyScope(t *testing.T) {
	st := makeStore(t)
	seedScope(t, st, "empty", map[string]string{})

	var buf bytes.Buffer
	err := inspect.Inspect(st, "empty", testPass, inspect.Options{}, &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "empty") {
		t.Errorf("expected empty-scope message, got: %s", buf.String())
	}
}
