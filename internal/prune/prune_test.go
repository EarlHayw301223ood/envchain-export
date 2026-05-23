package prune_test

import (
	"os"
	"testing"

	"github.com/user/envchain-export/internal/chain"
	"github.com/user/envchain-export/internal/prune"
	"github.com/user/envchain-export/internal/store"
)

const testPass = "hunter2"

func makeStore(t *testing.T) string {
	t.Helper()
	dir, err := os.MkdirTemp("", "prune-test-*")
	if err != nil {
		t.Fatalf("mkdirtemp: %v", err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })
	return dir
}

func seedScope(t *testing.T, dir, scope string, pairs map[string]string) {
	t.Helper()
	ch, err := chain.New(scope)
	if err != nil {
		t.Fatalf("chain.New: %v", err)
	}
	for k, v := range pairs {
		if err := ch.Add(k, v); err != nil {
			t.Fatalf("chain.Add(%q): %v", k, err)
		}
	}
	if err := store.New(dir).Save(scope, ch, testPass); err != nil {
		t.Fatalf("store.Save: %v", err)
	}
}

func TestByPrefix_RemovesMatchingKeys(t *testing.T) {
	dir := makeStore(t)
	seedScope(t, dir, "myapp", map[string]string{
		"AWS_ACCESS_KEY": "abc",
		"AWS_SECRET_KEY": "xyz",
		"DATABASE_URL":   "postgres://",
	})

	res, err := prune.ByPrefix(dir, "myapp", "AWS_", testPass)
	if err != nil {
		t.Fatalf("ByPrefix: %v", err)
	}
	if len(res.Removed) != 2 {
		t.Fatalf("expected 2 removed, got %d: %v", len(res.Removed), res.Removed)
	}

	ch, err := store.New(dir).Load("myapp", testPass)
	if err != nil {
		t.Fatalf("load after prune: %v", err)
	}
	if _, ok := ch.Get("DATABASE_URL"); !ok {
		t.Error("DATABASE_URL should still exist")
	}
	if _, ok := ch.Get("AWS_ACCESS_KEY"); ok {
		t.Error("AWS_ACCESS_KEY should have been removed")
	}
}

func TestByPrefix_NoMatch_ReturnsEmptyRemoved(t *testing.T) {
	dir := makeStore(t)
	seedScope(t, dir, "myapp", map[string]string{"FOO": "bar"})

	res, err := prune.ByPrefix(dir, "myapp", "NOMATCH_", testPass)
	if err != nil {
		t.Fatalf("ByPrefix: %v", err)
	}
	if len(res.Removed) != 0 {
		t.Errorf("expected no removals, got %v", res.Removed)
	}
}

func TestByPrefix_EmptyPrefix_ReturnsError(t *testing.T) {
	dir := makeStore(t)
	seedScope(t, dir, "myapp", map[string]string{"FOO": "bar"})

	_, err := prune.ByPrefix(dir, "myapp", "", testPass)
	if err == nil {
		t.Fatal("expected error for empty prefix")
	}
}

func TestByGlob_RemovesMatchingKeys(t *testing.T) {
	dir := makeStore(t)
	seedScope(t, dir, "svc", map[string]string{
		"CACHE_HOST": "localhost",
		"CACHE_PORT": "6379",
		"APP_PORT":   "8080",
	})

	res, err := prune.ByGlob(dir, "svc", "CACHE_*", testPass)
	if err != nil {
		t.Fatalf("ByGlob: %v", err)
	}
	if len(res.Removed) != 2 {
		t.Fatalf("expected 2 removed, got %d", len(res.Removed))
	}
}

func TestByGlob_InvalidPattern_ReturnsError(t *testing.T) {
	dir := makeStore(t)
	seedScope(t, dir, "svc", map[string]string{"FOO": "bar"})

	_, err := prune.ByGlob(dir, "svc", "[", testPass)
	if err == nil {
		t.Fatal("expected error for invalid glob pattern")
	}
}

func TestByPrefix_WrongPassphrase_ReturnsError(t *testing.T) {
	dir := makeStore(t)
	seedScope(t, dir, "sec", map[string]string{"KEY": "val"})

	_, err := prune.ByPrefix(dir, "sec", "KEY", "wrongpass")
	if err == nil {
		t.Fatal("expected error with wrong passphrase")
	}
}
