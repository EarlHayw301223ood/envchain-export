package lock_test

import (
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/nicholasgasior/envchain-export/internal/lock"
)

func makeStoreDir(t *testing.T) string {
	t.Helper()
	dir, err := os.MkdirTemp("", "lock-test-*")
	if err != nil {
		t.Fatalf("MkdirTemp: %v", err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })
	return dir
}

func TestAcquire_CreatesLockFile(t *testing.T) {
	dir := makeStoreDir(t)
	if err := lock.Acquire(dir, "myapp"); err != nil {
		t.Fatalf("Acquire: %v", err)
	}
	path := filepath.Join(dir, "myapp.lock")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("lock file not created: %v", err)
	}
	pid, err := strconv.Atoi(string(data))
	if err != nil {
		t.Fatalf("lock file content not a PID: %q", data)
	}
	if pid != os.Getpid() {
		t.Errorf("expected PID %d, got %d", os.Getpid(), pid)
	}
}

func TestAcquire_StaleLock_Succeeds(t *testing.T) {
	dir := makeStoreDir(t)
	// Write a lock file with PID 1 (init) — not owned by us but always alive
	// on Linux; use PID 0 which is always invalid.
	path := filepath.Join(dir, "stale.lock")
	if err := os.WriteFile(path, []byte("0"), 0o600); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}
	if err := lock.Acquire(dir, "stale"); err != nil {
		t.Fatalf("expected stale lock to be overwritten, got: %v", err)
	}
}

func TestAcquire_AlreadyLocked_ReturnsErrLocked(t *testing.T) {
	dir := makeStoreDir(t)
	if err := lock.Acquire(dir, "prod"); err != nil {
		t.Fatalf("first Acquire: %v", err)
	}
	// Second acquire from same process — our own PID is alive, so it is locked.
	err := lock.Acquire(dir, "prod")
	if err == nil {
		t.Fatal("expected ErrLocked, got nil")
	}
	if !isErrLocked(err) {
		t.Errorf("expected ErrLocked, got: %v", err)
	}
}

func TestRelease_RemovesLockFile(t *testing.T) {
	dir := makeStoreDir(t)
	if err := lock.Acquire(dir, "dev"); err != nil {
		t.Fatalf("Acquire: %v", err)
	}
	if err := lock.Release(dir, "dev"); err != nil {
		t.Fatalf("Release: %v", err)
	}
	path := filepath.Join(dir, "dev.lock")
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Error("expected lock file to be removed")
	}
}

func TestRelease_NotLocked_ReturnsErrNotLocked(t *testing.T) {
	dir := makeStoreDir(t)
	if err := lock.Release(dir, "ghost"); err != lock.ErrNotLocked {
		t.Errorf("expected ErrNotLocked, got: %v", err)
	}
}

func TestIsLocked_TrueAfterAcquire(t *testing.T) {
	dir := makeStoreDir(t)
	if lock.IsLocked(dir, "svc") {
		t.Fatal("expected not locked before Acquire")
	}
	if err := lock.Acquire(dir, "svc"); err != nil {
		t.Fatalf("Acquire: %v", err)
	}
	if !lock.IsLocked(dir, "svc") {
		t.Error("expected IsLocked to return true after Acquire")
	}
}

func isErrLocked(err error) bool {
	return err != nil && (err == lock.ErrLocked ||
		len(err.Error()) > 0 && containsStr(err.Error(), "locked"))
}

func containsStr(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub ||
		func() bool {
			for i := 0; i <= len(s)-len(sub); i++ {
				if s[i:i+len(sub)] == sub {
					return true
				}
			}
			return false
		}())
}
