package pin_test

import (
	"os"
	"testing"

	"github.com/nicholasgasior/envchain-export/internal/pin"
)

func makeStoreDir(t *testing.T) string {
	t.Helper()
	dir, err := os.MkdirTemp("", "pin-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })
	return dir
}

func TestAdd_AndIsPinned(t *testing.T) {
	dir := makeStoreDir(t)
	if err := pin.Add(dir, "myscope"); err != nil {
		t.Fatalf("Add: %v", err)
	}
	ok, err := pin.IsPinned(dir, "myscope")
	if err != nil {
		t.Fatalf("IsPinned: %v", err)
	}
	if !ok {
		t.Error("expected scope to be pinned")
	}
}

func TestAdd_DuplicateReturnsError(t *testing.T) {
	dir := makeStoreDir(t)
	_ = pin.Add(dir, "myscope")
	err := pin.Add(dir, "myscope")
	if err != pin.ErrAlreadyPinned {
		t.Errorf("expected ErrAlreadyPinned, got %v", err)
	}
}

func TestRemove_Existing(t *testing.T) {
	dir := makeStoreDir(t)
	_ = pin.Add(dir, "myscope")
	if err := pin.Remove(dir, "myscope"); err != nil {
		t.Fatalf("Remove: %v", err)
	}
	ok, _ := pin.IsPinned(dir, "myscope")
	if ok {
		t.Error("expected scope to be unpinned after Remove")
	}
}

func TestRemove_NotFound(t *testing.T) {
	dir := makeStoreDir(t)
	err := pin.Remove(dir, "ghost")
	if err != pin.ErrNotPinned {
		t.Errorf("expected ErrNotPinned, got %v", err)
	}
}

func TestList_ReturnsPinnedScopes(t *testing.T) {
	dir := makeStoreDir(t)
	_ = pin.Add(dir, "alpha")
	_ = pin.Add(dir, "beta")
	scopes, err := pin.List(dir)
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(scopes) != 2 {
		t.Errorf("expected 2 scopes, got %d", len(scopes))
	}
}

func TestList_EmptyReturnsNil(t *testing.T) {
	dir := makeStoreDir(t)
	scopes, err := pin.List(dir)
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(scopes) != 0 {
		t.Errorf("expected 0 scopes, got %d", len(scopes))
	}
}

func TestIsPinned_UnknownScope_ReturnsFalse(t *testing.T) {
	dir := makeStoreDir(t)
	ok, err := pin.IsPinned(dir, "unknown")
	if err != nil {
		t.Fatalf("IsPinned: %v", err)
	}
	if ok {
		t.Error("expected false for unknown scope")
	}
}
