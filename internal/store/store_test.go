package store_test

import (
	"os"
	"testing"

	"github.com/user/envchain-export/internal/chain"
	"github.com/user/envchain-export/internal/store"
)

func TestSaveAndLoad_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	s := store.New(dir)
	passphrase := "test-pass"

	c, err := chain.New("myapp")
	if err != nil {
		t.Fatalf("chain.New: %v", err)
	}
	c.Add("DB_URL", "postgres://localhost/db")
	c.Add("API_KEY", "secret123")

	if err := s.Save(c, passphrase); err != nil {
		t.Fatalf("Save: %v", err)
	}

	loaded, err := s.Load("myapp", passphrase)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	if loaded.Scope() != "myapp" {
		t.Errorf("expected scope myapp, got %s", loaded.Scope())
	}
	if v, _ := loaded.Get("DB_URL"); v != "postgres://localhost/db" {
		t.Errorf("unexpected DB_URL value: %s", v)
	}
	if v, _ := loaded.Get("API_KEY"); v != "secret123" {
		t.Errorf("unexpected API_KEY value: %s", v)
	}
}

func TestLoad_WrongPassphrase(t *testing.T) {
	dir := t.TempDir()
	s := store.New(dir)

	c, _ := chain.New("myapp")
	c.Add("KEY", "value")
	s.Save(c, "correct")

	_, err := s.Load("myapp", "wrong")
	if err == nil {
		t.Fatal("expected error with wrong passphrase")
	}
}

func TestLoad_MissingFile(t *testing.T) {
	dir := t.TempDir()
	s := store.New(dir)

	_, err := s.Load("nonexistent", "pass")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestSave_CreatesDirectory(t *testing.T) {
	dir := t.TempDir()
	nestedDir := dir + "/nested/store"
	s := store.New(nestedDir)

	c, _ := chain.New("scope1")
	if err := s.Save(c, "pass"); err != nil {
		t.Fatalf("Save: %v", err)
	}

	if _, err := os.Stat(nestedDir); os.IsNotExist(err) {
		t.Fatal("expected directory to be created")
	}
}
