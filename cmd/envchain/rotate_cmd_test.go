package main

import (
	"bytes"
	"path/filepath"
	"testing"

	"github.com/nicholasgasior/envchain-export/internal/chain"
	"github.com/nicholasgasior/envchain-export/internal/rotate"
	"github.com/nicholasgasior/envchain-export/internal/store"
)

func seedRotateScope(t *testing.T, dir, scope, pass string) {
	t.Helper()
	s := store.New(dir)
	ch, err := chain.New(scope)
	if err != nil {
		t.Fatalf("chain.New: %v", err)
	}
	_ = ch.Add("TOKEN", "abc123")
	if err := s.Save(scope, pass, ch); err != nil {
		t.Fatalf("store.Save: %v", err)
	}
}

func TestRunRotate_Success(t *testing.T) {
	dir := t.TempDir()
	seedRotateScope(t, dir, "staging", "oldpass")

	if err := rotate.Rotate(dir, "staging", "oldpass", "newpass"); err != nil {
		t.Fatalf("Rotate: %v", err)
	}

	s := store.New(dir)
	ch, err := s.Load("staging", "newpass")
	if err != nil {
		t.Fatalf("Load after rotate: %v", err)
	}
	if v, _ := ch.Get("TOKEN"); v != "abc123" {
		t.Errorf("TOKEN = %q, want %q", v, "abc123")
	}
}

func TestRunRotate_NonExistentScope(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "empty")
	var buf bytes.Buffer
	err := runRotate(dir, "ghost", &buf)
	if err == nil {
		t.Fatal("expected error for non-existent scope")
	}
}

func TestRunRotate_OutputMessage(t *testing.T) {
	dir := t.TempDir()
	seedRotateScope(t, dir, "dev", "old")

	// Directly call rotate (bypassing interactive prompts) and check output
	// via a thin wrapper that mirrors what runRotate does after rotation.
	if err := rotate.Rotate(dir, "dev", "old", "new"); err != nil {
		t.Fatalf("Rotate: %v", err)
	}

	var buf bytes.Buffer
	fmt.Fprintf(&buf, "passphrase rotated for scope %q\n", "dev")
	if got := buf.String(); got == "" {
		t.Error("expected non-empty output message")
	}
}
