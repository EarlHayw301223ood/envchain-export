package audit_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/user/envchain-export/internal/audit"
)

func makeStoreDir(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	return dir
}

func TestLog_CreatesEntryAndRead(t *testing.T) {
	dir := makeStoreDir(t)

	if err := audit.Log(dir, "prod", "set", "API_KEY"); err != nil {
		t.Fatalf("Log: %v", err)
	}

	entries, err := audit.Read(dir)
	if err != nil {
		t.Fatalf("Read: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	e := entries[0]
	if e.Scope != "prod" {
		t.Errorf("scope: got %q, want %q", e.Scope, "prod")
	}
	if e.Action != "set" {
		t.Errorf("action: got %q, want %q", e.Action, "set")
	}
	if e.Detail != "API_KEY" {
		t.Errorf("detail: got %q, want %q", e.Detail, "API_KEY")
	}
	if e.Timestamp.IsZero() {
		t.Error("timestamp should not be zero")
	}
}

func TestLog_AppendsMultipleEntries(t *testing.T) {
	dir := makeStoreDir(t)

	actions := []string{"set", "delete", "rotate"}
	for _, a := range actions {
		if err := audit.Log(dir, "staging", a, ""); err != nil {
			t.Fatalf("Log(%s): %v", a, err)
		}
	}

	entries, err := audit.Read(dir)
	if err != nil {
		t.Fatalf("Read: %v", err)
	}
	if len(entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(entries))
	}
	for i, a := range actions {
		if entries[i].Action != a {
			t.Errorf("entry %d action: got %q, want %q", i, entries[i].Action, a)
		}
	}
}

func TestRead_NoLogFile_ReturnsEmpty(t *testing.T) {
	dir := makeStoreDir(t)

	entries, err := audit.Read(dir)
	if err != nil {
		t.Fatalf("Read: %v", err)
	}
	if len(entries) != 0 {
		t.Errorf("expected 0 entries, got %d", len(entries))
	}
}

func TestLog_TimestampIsUTC(t *testing.T) {
	dir := makeStoreDir(t)
	before := time.Now().UTC().Add(-time.Second)

	if err := audit.Log(dir, "dev", "inspect", ""); err != nil {
		t.Fatalf("Log: %v", err)
	}

	after := time.Now().UTC().Add(time.Second)
	entries, _ := audit.Read(dir)
	ts := entries[0].Timestamp
	if ts.Before(before) || ts.After(after) {
		t.Errorf("timestamp %v not in expected range [%v, %v]", ts, before, after)
	}
}

func TestLog_FilePermissions(t *testing.T) {
	dir := makeStoreDir(t)
	if err := audit.Log(dir, "prod", "set", ""); err != nil {
		t.Fatalf("Log: %v", err)
	}

	info, err := os.Stat(filepath.Join(dir, "audit.log"))
	if err != nil {
		t.Fatalf("Stat: %v", err)
	}
	if perm := info.Mode().Perm(); perm != 0600 {
		t.Errorf("expected perm 0600, got %o", perm)
	}
}
