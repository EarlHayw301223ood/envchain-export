package main

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/user/envchain-export/internal/chain"
	"github.com/user/envchain-export/internal/store"
)

func seedEnvScope(t *testing.T, dir, scope, pass string, pairs map[string]string) {
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
	s := store.New(dir)
	if err := s.Save(ch, pass); err != nil {
		t.Fatalf("store.Save: %v", err)
	}
}

func TestRunEnv_Success(t *testing.T) {
	dir := t.TempDir()
	seedEnvScope(t, dir, "myscope", "secret", map[string]string{
		"ENVCHAIN_HELLO": "world",
	})

	cmd := newEnvCmd(dir)
	cmd.SetArgs([]string{"myscope", "--", "printenv", "ENVCHAIN_HELLO"})
	cmd.Flags().Set("passphrase", "secret")

	var buf bytes.Buffer
	cmd.SetOut(&buf)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRunEnv_NonExistentScope_ReturnsError(t *testing.T) {
	dir := t.TempDir()

	cmd := newEnvCmd(dir)
	cmd.SetArgs([]string{"missing", "--", "printenv", "X"})
	cmd.Flags().Set("passphrase", "pass")
	cmd.SilenceErrors()
	cmd.SilenceUsage()

	if err := cmd.Execute(); err == nil {
		t.Fatal("expected error for missing scope, got nil")
	}
}

func TestRunEnv_StoreFilePresent(t *testing.T) {
	dir := t.TempDir()
	seedEnvScope(t, dir, "checkscope", "pw", map[string]string{
		"MY_KEY": "MY_VAL",
	})

	expected := filepath.Join(dir, "checkscope.enc")
	if _, err := os.Stat(expected); err != nil {
		t.Fatalf("store file not found: %v", err)
	}
}
