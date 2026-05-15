package main

import (
	"os"
	"path/filepath"
	"testing"
)

func setupStore(t *testing.T, scopes ...string) string {
	t.Helper()
	dir := t.TempDir()
	for _, s := range scopes {
		path := filepath.Join(dir, s+".enc")
		if err := os.WriteFile(path, []byte("enc"), 0o600); err != nil {
			t.Fatalf("setup: %v", err)
		}
	}
	return dir
}

func TestRunScopeList_PrintsScopes(t *testing.T) {
	dir := setupStore(t, "alpha", "beta")
	if err := runScopeList(dir); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRunScopeList_EmptyStore_ReturnsError(t *testing.T) {
	dir := setupStore(t)
	if err := runScopeList(dir); err == nil {
		t.Fatal("expected error for empty store")
	}
}

func TestRunScopeDelete_Force_Removes(t *testing.T) {
	dir := setupStore(t, "todelete")
	if err := runScopeDelete(dir, "todelete", true); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := os.Stat(filepath.Join(dir, "todelete.enc")); !os.IsNotExist(err) {
		t.Fatal("expected file to be removed")
	}
}

func TestRunScopeDelete_NonExistent_ReturnsError(t *testing.T) {
	dir := setupStore(t)
	if err := runScopeDelete(dir, "ghost", true); err == nil {
		t.Fatal("expected error for missing scope")
	}
}

func TestRunScopeExists_Present(t *testing.T) {
	dir := setupStore(t, "present")
	if err := runScopeExists(dir, "present"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
