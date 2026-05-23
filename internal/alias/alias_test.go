package alias_test

import (
	"os"
	"testing"

	"github.com/nicholasgasior/envchain-export/internal/alias"
)

func makeStoreDir(t *testing.T) string {
	t.Helper()
	dir, err := os.MkdirTemp("", "alias-test-*")
	if err != nil {
		t.Fatalf("MkdirTemp: %v", err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })
	return dir
}

func TestAdd_AndResolve(t *testing.T) {
	dir := makeStoreDir(t)
	if err := alias.Add(dir, "prod", "production"); err != nil {
		t.Fatalf("Add: %v", err)
	}
	target, err := alias.Resolve(dir, "prod")
	if err != nil {
		t.Fatalf("Resolve: %v", err)
	}
	if target != "production" {
		t.Errorf("expected production, got %s", target)
	}
}

func TestAdd_DuplicateReturnsError(t *testing.T) {
	dir := makeStoreDir(t)
	_ = alias.Add(dir, "prod", "production")
	err := alias.Add(dir, "prod", "production-v2")
	if err == nil {
		t.Fatal("expected error for duplicate alias")
	}
}

func TestAdd_CircularAlias_ReturnsError(t *testing.T) {
	dir := makeStoreDir(t)
	_ = alias.Add(dir, "prod", "production")
	// Attempting to alias something that is itself an alias key.
	err := alias.Add(dir, "p", "prod")
	if err == nil {
		t.Fatal("expected ErrCircularAlias")
	}
}

func TestRemove_Existing(t *testing.T) {
	dir := makeStoreDir(t)
	_ = alias.Add(dir, "prod", "production")
	if err := alias.Remove(dir, "prod"); err != nil {
		t.Fatalf("Remove: %v", err)
	}
	_, err := alias.Resolve(dir, "prod")
	if err == nil {
		t.Fatal("expected error after removal")
	}
}

func TestRemove_NotFound_ReturnsError(t *testing.T) {
	dir := makeStoreDir(t)
	err := alias.Remove(dir, "nonexistent")
	if err == nil {
		t.Fatal("expected ErrAliasNotFound")
	}
}

func TestList_ReturnsSortedAliases(t *testing.T) {
	dir := makeStoreDir(t)
	_ = alias.Add(dir, "z-scope", "zz")
	_ = alias.Add(dir, "a-scope", "aa")
	_ = alias.Add(dir, "m-scope", "mm")

	list, err := alias.List(dir)
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(list) != 3 {
		t.Errorf("expected 3 aliases, got %d", len(list))
	}
	if list["a-scope"] != "aa" || list["z-scope"] != "zz" {
		t.Errorf("unexpected alias map: %v", list)
	}
}

func TestList_EmptyStore_ReturnsEmptyMap(t *testing.T) {
	dir := makeStoreDir(t)
	list, err := alias.List(dir)
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(list) != 0 {
		t.Errorf("expected empty map, got %v", list)
	}
}

func TestResolve_NotFound_ReturnsError(t *testing.T) {
	dir := makeStoreDir(t)
	_, err := alias.Resolve(dir, "missing")
	if err == nil {
		t.Fatal("expected ErrAliasNotFound")
	}
}
