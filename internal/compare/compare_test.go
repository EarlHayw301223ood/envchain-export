package compare_test

import (
	"bytes"
	"os"
	"testing"

	"github.com/user/envchain-export/internal/chain"
	"github.com/user/envchain-export/internal/compare"
	"github.com/user/envchain-export/internal/store"
)

const pass = "testpass"

func makeStore(t *testing.T) *store.Store {
	t.Helper()
	dir := t.TempDir()
	st, err := store.New(dir)
	if err != nil {
		t.Fatalf("store.New: %v", err)
	}
	return st
}

func seedScope(t *testing.T, st *store.Store, scope string, pairs map[string]string) {
	t.Helper()
	c, err := chain.New(scope)
	if err != nil {
		t.Fatalf("chain.New: %v", err)
	}
	for k, v := range pairs {
		if err := c.Add(k, v); err != nil {
			t.Fatalf("chain.Add: %v", err)
		}
	}
	if err := st.Save(c, pass); err != nil {
		t.Fatalf("store.Save: %v", err)
	}
}

func TestCompare_NoChanges(t *testing.T) {
	st := makeStore(t)
	seedScope(t, st, "alpha", map[string]string{"FOO": "bar", "BAZ": "qux"})
	seedScope(t, st, "beta", map[string]string{"FOO": "bar", "BAZ": "qux"})

	r, err := compare.Compare(st, "alpha", "beta", pass, pass)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(r.OnlyInA) != 0 || len(r.OnlyInB) != 0 || len(r.Changed) != 0 {
		t.Errorf("expected no differences, got %+v", r)
	}
	if len(r.Identical) != 2 {
		t.Errorf("expected 2 identical keys, got %d", len(r.Identical))
	}
}

func TestCompare_Differences(t *testing.T) {
	st := makeStore(t)
	seedScope(t, st, "alpha", map[string]string{"FOO": "bar", "ONLY_A": "1"})
	seedScope(t, st, "beta", map[string]string{"FOO": "changed", "ONLY_B": "2"})

	r, err := compare.Compare(st, "alpha", "beta", pass, pass)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(r.OnlyInA) != 1 || r.OnlyInA[0] != "ONLY_A" {
		t.Errorf("OnlyInA: got %v", r.OnlyInA)
	}
	if len(r.OnlyInB) != 1 || r.OnlyInB[0] != "ONLY_B" {
		t.Errorf("OnlyInB: got %v", r.OnlyInB)
	}
	if len(r.Changed) != 1 || r.Changed[0] != "FOO" {
		t.Errorf("Changed: got %v", r.Changed)
	}
}

func TestCompare_SameName_ReturnsError(t *testing.T) {
	st := makeStore(t)
	_, err := compare.Compare(st, "alpha", "alpha", pass, pass)
	if err == nil {
		t.Fatal("expected error for same scope name")
	}
}

func TestCompare_MissingScope_ReturnsError(t *testing.T) {
	st := makeStore(t)
	seedScope(t, st, "alpha", map[string]string{"FOO": "bar"})
	_, err := compare.Compare(st, "alpha", "missing", pass, pass)
	if err == nil {
		t.Fatal("expected error for missing scope")
	}
}

func TestWrite_OutputFormat(t *testing.T) {
	st := makeStore(t)
	seedScope(t, st, "alpha", map[string]string{"FOO": "same", "ONLY_A": "1"})
	seedScope(t, st, "beta", map[string]string{"FOO": "same", "ONLY_B": "2"})

	r, err := compare.Compare(st, "alpha", "beta", pass, pass)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var buf bytes.Buffer
	compare.Write(&buf, "alpha", "beta", r)
	out := buf.String()

	for _, want := range []string{"alpha", "beta", "ONLY_A", "ONLY_B", "FOO"} {
		if !bytes.Contains([]byte(out), []byte(want)) {
			t.Errorf("output missing %q:\n%s", want, out)
		}
	}
	_ = os.DevNull
}
