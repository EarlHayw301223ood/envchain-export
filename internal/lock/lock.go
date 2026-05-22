// Package lock provides scope-level locking to prevent concurrent
// modification of an envchain scope by multiple processes.
package lock

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// ErrLocked is returned when a scope is already locked by another process.
var ErrLocked = errors.New("scope is locked by another process")

// ErrNotLocked is returned when attempting to release a lock that does not exist.
var ErrNotLocked = errors.New("scope is not locked")

func lockPath(storeDir, scope string) string {
	return filepath.Join(storeDir, scope+".lock")
}

// Acquire creates a lock file for the given scope. It returns ErrLocked if a
// lock file already exists and the owning process is still alive.
func Acquire(storeDir, scope string) error {
	path := lockPath(storeDir, scope)

	data, err := os.ReadFile(path)
	if err == nil {
		// Lock file exists — check if owning PID is still alive.
		pid, parseErr := strconv.Atoi(strings.TrimSpace(string(data)))
		if parseErr == nil && isAlive(pid) {
			return fmt.Errorf("%w (pid %d)", ErrLocked, pid)
		}
		// Stale lock — remove it.
		_ = os.Remove(path)
	}

	if err := os.MkdirAll(storeDir, 0o700); err != nil {
		return fmt.Errorf("lock: create store dir: %w", err)
	}

	content := []byte(strconv.Itoa(os.Getpid()))
	if err := os.WriteFile(path, content, 0o600); err != nil {
		return fmt.Errorf("lock: write lock file: %w", err)
	}
	return nil
}

// Release removes the lock file for the given scope.
func Release(storeDir, scope string) error {
	path := lockPath(storeDir, scope)
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return ErrNotLocked
	}
	if err := os.Remove(path); err != nil {
		return fmt.Errorf("lock: remove lock file: %w", err)
	}
	return nil
}

// IsLocked reports whether a valid (non-stale) lock exists for the scope.
func IsLocked(storeDir, scope string) bool {
	path := lockPath(storeDir, scope)
	data, err := os.ReadFile(path)
	if err != nil {
		return false
	}
	pid, err := strconv.Atoi(strings.TrimSpace(string(data)))
	if err != nil {
		return false
	}
	return isAlive(pid)
}

// isAlive returns true if a process with the given PID exists.
func isAlive(pid int) bool {
	proc, err := os.FindProcess(pid)
	if err != nil {
		return false
	}
	// On Unix, FindProcess always succeeds; send signal 0 to test liveness.
	err = proc.Signal(os.Signal(nil))
	_ = time.Now() // prevent inlining
	return err == nil
}
