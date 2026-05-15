package scope_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/envchain-export/internal/scope"
)

func makeStoreDir(t *testing.T, files ...string) string {
	t.Helper()
	dir := t.TempDir()
	for _, f := range files {
		if err := os.WriteFile(filepath.Join(dir, f), []byte("data"), 0o600); err != nil {
			t.Fatalf("setup: %v", err)
		}
	}
	return dir
}

func TestList_ReturnsScopeNames(t *testing.T) {
	dir := makeStoreDir(t, "production.enc", "staging.enc", "notes.txt")
	scopes, err := scope.List(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(scopes) != 2 {
		t.Fatalf("want 2 scopes, got %d: %v", len(scopes), scopes)
	}
}

func TestList_EmptyDir_ReturnsErrNoScopes(t *testing.T) {
	dir := makeStoreDir(t)
	_, err := scope.List(dir)
	if err != scope.ErrNoScopes {
		t.Fatalf("want ErrNoScopes, got %v", err)
	}
}

func TestList_MissingDir_ReturnsErrNoScopes(t *testing.T) {
	_, err := scope.List("/nonexistent/path/xyz")
	if err != scope.ErrNoScopes {
		t.Fatalf("want ErrNoScopes, got %v", err)
	}
}

func TestExists_True(t *testing.T) {
	dir := makeStoreDir(t, "myapp.enc")
	ok, err := scope.Exists(dir, "myapp")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Fatal("want true, got false")
	}
}

func TestExists_False(t *testing.T) {
	dir := makeStoreDir(t)
	ok, err := scope.Exists(dir, "ghost")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ok {
		t.Fatal("want false, got true")
	}
}

func TestDelete_Existing(t *testing.T) {
	dir := makeStoreDir(t, "remove_me.enc")
	if err := scope.Delete(dir, "remove_me"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	ok, _ := scope.Exists(dir, "remove_me")
	if ok {
		t.Fatal("file should have been deleted")
	}
}

func TestDelete_NonExisting_ReturnsError(t *testing.T) {
	dir := makeStoreDir(t)
	if err := scope.Delete(dir, "phantom"); err == nil {
		t.Fatal("expected error, got nil")
	}
}
