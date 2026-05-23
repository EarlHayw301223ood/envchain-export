package history_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/nicholasgasior/envchain-export/internal/history"
)

func makeStoreDir(t *testing.T) string {
	t.Helper()
	dir, err := os.MkdirTemp("", "history-test-*")
	if err != nil {
		t.Fatalf("MkdirTemp: %v", err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })
	return dir
}

func TestRecord_AndRead_RoundTrip(t *testing.T) {
	dir := makeStoreDir(t)
	scope := "myapp"

	if err := history.Record(dir, scope, "set", "Added KEY1"); err != nil {
		t.Fatalf("Record: %v", err)
	}

	entries, err := history.Read(dir, scope)
	if err != nil {
		t.Fatalf("Read: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].Action != "set" {
		t.Errorf("Action = %q, want %q", entries[0].Action, "set")
	}
	if entries[0].Detail != "Added KEY1" {
		t.Errorf("Detail = %q, want %q", entries[0].Detail, "Added KEY1")
	}
}

func TestRecord_AppendsMultiple(t *testing.T) {
	dir := makeStoreDir(t)
	scope := "myapp"

	for i, action := range []string{"set", "remove", "rotate"} {
		if err := history.Record(dir, scope, action, ""); err != nil {
			t.Fatalf("Record[%d]: %v", i, err)
		}
	}

	entries, err := history.Read(dir, scope)
	if err != nil {
		t.Fatalf("Read: %v", err)
	}
	if len(entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(entries))
	}
}

func TestRead_NoHistory_ReturnsEmpty(t *testing.T) {
	dir := makeStoreDir(t)
	entries, err := history.Read(dir, "ghost")
	if err != nil {
		t.Fatalf("Read: %v", err)
	}
	if len(entries) != 0 {
		t.Errorf("expected empty slice, got %d entries", len(entries))
	}
}

func TestClear_RemovesFile(t *testing.T) {
	dir := makeStoreDir(t)
	scope := "myapp"

	_ = history.Record(dir, scope, "set", "")

	if err := history.Clear(dir, scope); err != nil {
		t.Fatalf("Clear: %v", err)
	}

	path := filepath.Join(dir, scope+".history.json")
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Errorf("expected history file to be removed")
	}
}

func TestRecord_TimestampIsUTC(t *testing.T) {
	dir := makeStoreDir(t)
	before := time.Now().UTC().Add(-time.Second)

	_ = history.Record(dir, "myapp", "set", "")

	entries, _ := history.Read(dir, "myapp")
	if len(entries) == 0 {
		t.Fatal("no entries")
	}
	ts := entries[0].Timestamp
	if ts.Location() != time.UTC {
		t.Errorf("timestamp not UTC: %v", ts.Location())
	}
	if ts.Before(before) {
		t.Errorf("timestamp %v is before test start %v", ts, before)
	}
}
