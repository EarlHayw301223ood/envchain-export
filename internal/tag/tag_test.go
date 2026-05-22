package tag_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/envchain/envchain-export/internal/tag"
)

func makeStoreDir(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	return filepath.Join(dir, "store")
}

func TestAdd_AndList(t *testing.T) {
	dir := makeStoreDir(t)
	if err := os.MkdirAll(dir, 0700); err != nil {
		t.Fatal(err)
	}

	if err := tag.Add(dir, "myapp", "production"); err != nil {
		t.Fatalf("Add: %v", err)
	}
	if err := tag.Add(dir, "myapp", "aws"); err != nil {
		t.Fatalf("Add second: %v", err)
	}

	tags, err := tag.List(dir, "myapp")
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(tags) != 2 {
		t.Fatalf("expected 2 tags, got %d", len(tags))
	}
	if tags[0] != "aws" || tags[1] != "production" {
		t.Errorf("unexpected tags: %v", tags)
	}
}

func TestAdd_DuplicateIsIdempotent(t *testing.T) {
	dir := makeStoreDir(t)
	os.MkdirAll(dir, 0700)

	tag.Add(dir, "myapp", "staging")
	tag.Add(dir, "myapp", "staging")

	tags, _ := tag.List(dir, "myapp")
	if len(tags) != 1 {
		t.Errorf("expected 1 tag, got %d", len(tags))
	}
}

func TestRemove_Existing(t *testing.T) {
	dir := makeStoreDir(t)
	os.MkdirAll(dir, 0700)

	tag.Add(dir, "myapp", "production")
	tag.Add(dir, "myapp", "aws")

	if err := tag.Remove(dir, "myapp", "aws"); err != nil {
		t.Fatalf("Remove: %v", err)
	}
	tags, _ := tag.List(dir, "myapp")
	if len(tags) != 1 || tags[0] != "production" {
		t.Errorf("unexpected tags after remove: %v", tags)
	}
}

func TestRemove_NotFound(t *testing.T) {
	dir := makeStoreDir(t)
	os.MkdirAll(dir, 0700)

	tag.Add(dir, "myapp", "production")
	err := tag.Remove(dir, "myapp", "nonexistent")
	if err != tag.ErrTagNotFound {
		t.Errorf("expected ErrTagNotFound, got %v", err)
	}
}

func TestList_NoTags(t *testing.T) {
	dir := makeStoreDir(t)
	os.MkdirAll(dir, 0700)

	_, err := tag.List(dir, "myapp")
	if err != tag.ErrNoTags {
		t.Errorf("expected ErrNoTags, got %v", err)
	}
}

func TestScopesByTag(t *testing.T) {
	dir := makeStoreDir(t)
	os.MkdirAll(dir, 0700)

	tag.Add(dir, "app1", "production")
	tag.Add(dir, "app2", "production")
	tag.Add(dir, "app3", "staging")

	scopes, err := tag.ScopesByTag(dir, "production")
	if err != nil {
		t.Fatalf("ScopesByTag: %v", err)
	}
	if len(scopes) != 2 {
		t.Fatalf("expected 2 scopes, got %d", len(scopes))
	}
	if scopes[0] != "app1" || scopes[1] != "app2" {
		t.Errorf("unexpected scopes: %v", scopes)
	}
}
