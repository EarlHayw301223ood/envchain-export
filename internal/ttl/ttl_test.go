package ttl_test

import (
	"os"
	"testing"
	"time"

	"github.com/user/envchain-export/internal/ttl"
)

func makeStoreDir(t *testing.T) string {
	t.Helper()
	dir, err := os.MkdirTemp("", "ttl-test-*")
	if err != nil {
		t.Fatalf("mkdirtemp: %v", err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })
	return dir
}

func TestSet_AndGet_RoundTrip(t *testing.T) {
	dir := makeStoreDir(t)
	d := 10 * time.Minute
	before := time.Now().UTC()
	if err := ttl.Set(dir, "myscope", d); err != nil {
		t.Fatalf("Set: %v", err)
	}
	got, err := ttl.Get(dir, "myscope")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if got.Before(before.Add(d - time.Second)) {
		t.Errorf("expiry too early: %v", got)
	}
}

func TestGet_NoExpiry_ReturnsErrNoExpiry(t *testing.T) {
	dir := makeStoreDir(t)
	_, err := ttl.Get(dir, "ghost")
	if err != ttl.ErrNoExpiry {
		t.Fatalf("expected ErrNoExpiry, got %v", err)
	}
}

func TestCheck_Valid(t *testing.T) {
	dir := makeStoreDir(t)
	if err := ttl.Set(dir, "s", time.Hour); err != nil {
		t.Fatal(err)
	}
	if err := ttl.Check(dir, "s"); err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestCheck_Expired(t *testing.T) {
	dir := makeStoreDir(t)
	if err := ttl.Set(dir, "s", -time.Second); err != nil {
		t.Fatal(err)
	}
	if err := ttl.Check(dir, "s"); err != ttl.ErrExpired {
		t.Fatalf("expected ErrExpired, got %v", err)
	}
}

func TestCheck_NoExpiry(t *testing.T) {
	dir := makeStoreDir(t)
	if err := ttl.Check(dir, "none"); err != ttl.ErrNoExpiry {
		t.Fatalf("expected ErrNoExpiry, got %v", err)
	}
}

func TestClear_RemovesExpiry(t *testing.T) {
	dir := makeStoreDir(t)
	if err := ttl.Set(dir, "s", time.Hour); err != nil {
		t.Fatal(err)
	}
	if err := ttl.Clear(dir, "s"); err != nil {
		t.Fatalf("Clear: %v", err)
	}
	if err := ttl.Check(dir, "s"); err != ttl.ErrNoExpiry {
		t.Fatalf("expected ErrNoExpiry after clear, got %v", err)
	}
}

func TestClear_NonExistent_IsNoop(t *testing.T) {
	dir := makeStoreDir(t)
	if err := ttl.Clear(dir, "ghost"); err != nil {
		t.Fatalf("Clear on missing scope should be nil, got %v", err)
	}
}
