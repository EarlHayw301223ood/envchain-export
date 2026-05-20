package rekey_test

import (
	"path/filepath"
	"testing"

	"github.com/user/envchain-export/internal/chain"
	"github.com/user/envchain-export/internal/rekey"
	"github.com/user/envchain-export/internal/store"
)

func seedScope(t *testing.T, dir, name, pass string) {
	t.Helper()
	s, err := store.New(dir, name)
	if err != nil {
		t.Fatalf("store.New: %v", err)
	}
	c, err := chain.New(name)
	if err != nil {
		t.Fatalf("chain.New: %v", err)
	}
	if err := c.Add("KEY", "value"); err != nil {
		t.Fatalf("chain.Add: %v", err)
	}
	if err := s.Save(c, pass); err != nil {
		t.Fatalf("store.Save: %v", err)
	}
}

func TestRekey_Success(t *testing.T) {
	dir := t.TempDir()
	seedScope(t, dir, "alpha", "old")
	seedScope(t, dir, "beta", "old")

	results, err := rekey.Rekey(dir, "old", "new")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	for _, r := range results {
		if !r.Success {
			t.Errorf("scope %q failed: %v", r.Scope, r.Err)
		}
	}

	// Verify new passphrase works
	s, _ := store.New(dir, "alpha")
	if _, err := s.Load("new"); err != nil {
		t.Errorf("load with new passphrase failed: %v", err)
	}
}

func TestRekey_WrongOldPassphrase(t *testing.T) {
	dir := t.TempDir()
	seedScope(t, dir, "gamma", "correct")

	results, err := rekey.Rekey(dir, "wrong", "new")
	if err != nil {
		t.Fatalf("unexpected top-level error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Success {
		t.Error("expected failure for wrong passphrase")
	}
}

func TestRekey_EmptyStore(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "nostore")
	_, err := rekey.Rekey(dir, "old", "new")
	if err == nil {
		t.Error("expected error for missing store directory")
	}
}

func TestRekey_PartialFailure(t *testing.T) {
	dir := t.TempDir()
	seedScope(t, dir, "ok", "pass")
	seedScope(t, dir, "bad", "different")

	results, err := rekey.Rekey(dir, "pass", "newpass")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var successes, failures int
	for _, r := range results {
		if r.Success {
			successes++
		} else {
			failures++
		}
	}
	if successes != 1 || failures != 1 {
		t.Errorf("expected 1 success and 1 failure, got %d/%d", successes, failures)
	}
}
