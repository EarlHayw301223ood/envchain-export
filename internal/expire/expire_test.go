package expire_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/user/envchain-export/internal/chain"
	"github.com/user/envchain-export/internal/expire"
	"github.com/user/envchain-export/internal/store"
	"github.com/user/envchain-export/internal/ttl"
)

const testPass = "s3cret"

func makeStoreDir(t *testing.T) string {
	t.Helper()
	dir, err := os.MkdirTemp("", "expire-test-*")
	if err != nil {
		t.Fatalf("MkdirTemp: %v", err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })
	return filepath.Join(dir, "store")
}

func seedScope(t *testing.T, storeDir, name string) {
	t.Helper()
	ch, err := chain.New(name)
	if err != nil {
		t.Fatalf("chain.New: %v", err)
	}
	ch.Add("KEY", "value")
	st := store.New(storeDir)
	if err := st.Save(name, ch, testPass); err != nil {
		t.Fatalf("store.Save: %v", err)
	}
}

func TestCheck_ReturnsExpiredScope(t *testing.T) {
	storeDir := makeStoreDir(t)
	seedScope(t, storeDir, "alpha")
	past := time.Now().UTC().Add(-1 * time.Hour)
	if err := ttl.Set(storeDir, "alpha", past); err != nil {
		t.Fatalf("ttl.Set: %v", err)
	}

	results, err := expire.Check(storeDir)
	if err != nil {
		t.Fatalf("Check: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if !results[0].Expired {
		t.Errorf("expected scope to be expired")
	}
	if results[0].Scope != "alpha" {
		t.Errorf("expected scope %q, got %q", "alpha", results[0].Scope)
	}
}

func TestCheck_SkopesWithoutTTL(t *testing.T) {
	storeDir := makeStoreDir(t)
	seedScope(t, storeDir, "beta")

	results, err := expire.Check(storeDir)
	if err != nil {
		t.Fatalf("Check: %v", err)
	}
	if len(results) != 0 {
		t.Errorf("expected 0 results for scope without TTL, got %d", len(results))
	}
}

func TestPurge_RemovesExpiredScope(t *testing.T) {
	storeDir := makeStoreDir(t)
	seedScope(t, storeDir, "gamma")
	seedScope(t, storeDir, "delta")

	// only gamma is expired
	past := time.Now().UTC().Add(-2 * time.Hour)
	if err := ttl.Set(storeDir, "gamma", past); err != nil {
		t.Fatalf("ttl.Set: %v", err)
	}
	future := time.Now().UTC().Add(2 * time.Hour)
	if err := ttl.Set(storeDir, "delta", future); err != nil {
		t.Fatalf("ttl.Set: %v", err)
	}

	removed, err := expire.Purge(storeDir)
	if err != nil {
		t.Fatalf("Purge: %v", err)
	}
	if len(removed) != 1 || removed[0] != "gamma" {
		t.Errorf("expected [gamma] removed, got %v", removed)
	}
}

func TestPurge_EmptyStore_ReturnsError(t *testing.T) {
	storeDir := makeStoreDir(t)
	_, err := expire.Purge(storeDir)
	if err == nil {
		t.Fatal("expected error for empty store, got nil")
	}
}
