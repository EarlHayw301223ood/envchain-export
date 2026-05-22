package search_test

import (
	"testing"

	"github.com/your-org/envchain-export/internal/chain"
	"github.com/your-org/envchain-export/internal/search"
	"github.com/your-org/envchain-export/internal/store"
)

const testPass = "hunter2"

func seedScope(t *testing.T, dir, scopeName string, pairs map[string]string) {
	t.Helper()
	ch, err := chain.New(scopeName)
	if err != nil {
		t.Fatalf("chain.New: %v", err)
	}
	for k, v := range pairs {
		if err := ch.Add(k, v); err != nil {
			t.Fatalf("ch.Add(%q): %v", k, err)
		}
	}
	st := store.New(dir, scopeName)
	if err := st.Save(ch, testPass); err != nil {
		t.Fatalf("st.Save: %v", err)
	}
}

func TestSearch_FindsByKey(t *testing.T) {
	dir := t.TempDir()
	seedScope(t, dir, "prod", map[string]string{"DATABASE_URL": "postgres://", "API_KEY": "abc"})
	seedScope(t, dir, "dev", map[string]string{"DATABASE_URL": "sqlite://", "DEBUG": "true"})

	matches, err := search.Search(dir, testPass, "DATABASE", search.Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(matches) != 2 {
		t.Fatalf("want 2 matches, got %d", len(matches))
	}
}

func TestSearch_CaseInsensitive(t *testing.T) {
	dir := t.TempDir()
	seedScope(t, dir, "prod", map[string]string{"Api_Token": "secret"})

	matches, err := search.Search(dir, testPass, "api_token", search.Options{CaseInsensitive: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(matches) != 1 {
		t.Fatalf("want 1 match, got %d", len(matches))
	}
}

func TestSearch_SearchValues(t *testing.T) {
	dir := t.TempDir()
	seedScope(t, dir, "prod", map[string]string{"SOME_KEY": "postgres://localhost/db"})

	matches, err := search.Search(dir, testPass, "postgres", search.Options{SearchValues: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(matches) != 1 {
		t.Fatalf("want 1 match, got %d", len(matches))
	}
	if matches[0].Key != "SOME_KEY" {
		t.Errorf("unexpected key: %s", matches[0].Key)
	}
}

func TestSearch_NoMatch(t *testing.T) {
	dir := t.TempDir()
	seedScope(t, dir, "prod", map[string]string{"FOO": "bar"})

	matches, err := search.Search(dir, testPass, "NONEXISTENT", search.Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(matches) != 0 {
		t.Fatalf("want 0 matches, got %d", len(matches))
	}
}

func TestFilterByScope(t *testing.T) {
	matches := []search.Match{
		{Scope: "prod", Key: "A"},
		{Scope: "dev", Key: "A"},
		{Scope: "prod", Key: "B"},
	}
	got := search.FilterByScope(matches, "prod")
	if len(got) != 2 {
		t.Fatalf("want 2, got %d", len(got))
	}
}

func TestKeys_Deduplicated(t *testing.T) {
	matches := []search.Match{
		{Scope: "prod", Key: "DB"},
		{Scope: "dev", Key: "DB"},
		{Scope: "prod", Key: "API"},
	}
	keys := search.Keys(matches)
	if len(keys) != 2 {
		t.Fatalf("want 2 unique keys, got %d", len(keys))
	}
}
