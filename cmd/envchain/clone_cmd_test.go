package main

import (
	"bytes"
	"path/filepath"
	"strings"
	"testing"

	"github.com/nicholasgasior/envchain-export/internal/chain"
	"github.com/nicholasgasior/envchain-export/internal/store"
)

func seedCloneScope(t *testing.T, st *store.Store, scope, pass string) {
	t.Helper()
	ch, err := chain.New(scope)
	if err != nil {
		t.Fatalf("chain.New: %v", err)
	}
	_ = ch.Add("ALPHA", "one")
	_ = ch.Add("BETA", "two")
	if err := st.Save(scope, pass, ch); err != nil {
		t.Fatalf("st.Save: %v", err)
	}
}

func setupCloneStore(t *testing.T) *store.Store {
	t.Helper()
	dir := filepath.Join(t.TempDir(), ".envchain")
	st, err := store.New(dir)
	if err != nil {
		t.Fatalf("store.New: %v", err)
	}
	return st
}

func TestRunClone_OutputsMessage(t *testing.T) {
	st := setupCloneStore(t)
	seedCloneScope(t, st, "prod", "secret")

	buf := &bytes.Buffer{}
	cmd := newCloneCmd(st)
	cmd.SetOut(buf)

	// Inject passphrase via env to avoid interactive prompt in tests.
	t.Setenv("ENVCHAIN_PASSPHRASE", "secret")

	cmd.SetArgs([]string{"prod", "staging"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "prod") || !strings.Contains(out, "staging") {
		t.Errorf("expected output to mention both scopes, got: %q", out)
	}
}

func TestRunClone_NonExistentScope_ReturnsError(t *testing.T) {
	st := setupCloneStore(t)

	t.Setenv("ENVCHAIN_PASSPHRASE", "secret")

	cmd := newCloneCmd(st)
	cmd.SetArgs([]string{"ghost", "copy"})
	if err := cmd.Execute(); err == nil {
		t.Error("expected error for non-existent source scope, got nil")
	}
}

func TestRunClone_SameName_ReturnsError(t *testing.T) {
	st := setupCloneStore(t)
	seedCloneScope(t, st, "prod", "secret")

	t.Setenv("ENVCHAIN_PASSPHRASE", "secret")

	cmd := newCloneCmd(st)
	cmd.SetArgs([]string{"prod", "prod"})
	if err := cmd.Execute(); err == nil {
		t.Error("expected error when src == dst, got nil")
	}
}
