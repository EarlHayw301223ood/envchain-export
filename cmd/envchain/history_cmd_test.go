package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/nicholasgasior/envchain-export/internal/history"
	"github.com/nicholasgasior/envchain-export/internal/store"
)

func seedHistoryScope(t *testing.T, storeDir string) {
	t.Helper()
	st := store.New(storeDir)
	passphrase := "test-pass"

	ch, err := chainNew("history-scope")
	if err != nil {
		t.Fatalf("newChain: %v", err)
	}
	_ = ch.Add("API_KEY", "abc123")
	_ = ch.Add("DB_URL", "postgres://localhost/db")

	if err := st.Save(ch, passphrase); err != nil {
		t.Fatalf("Save: %v", err)
	}

	// Record a few history entries
	for _, action := range []string{"set API_KEY", "set DB_URL", "remove OLD_KEY"} {
		if err := history.Record(storeDir, "history-scope", action); err != nil {
			t.Fatalf("Record: %v", err)
		}
	}
}

func setupHistoryStore(t *testing.T) (storeDir string, cleanup func()) {
	t.Helper()
	dir, err := os.MkdirTemp("", "envchain-history-cmd-*")
	if err != nil {
		t.Fatalf("MkdirTemp: %v", err)
	}
	return dir, func() { os.RemoveAll(dir) }
}

func chainNew(scope string) (interface{ Add(string, string) error }, error) {
	// stub — real code uses internal/chain
	return nil, nil
}

func TestRunHistoryList_PrintsEntries(t *testing.T) {
	storeDir, cleanup := setupHistoryStore(t)
	defer cleanup()
	seedHistoryScope(t, storeDir)

	var buf bytes.Buffer
	err := runHistoryList(&buf, storeDir, "history-scope")
	if err != nil {
		t.Fatalf("runHistoryList: %v", err)
	}

	out := buf.String()
	for _, action := range []string{"set API_KEY", "set DB_URL", "remove OLD_KEY"} {
		if !strings.Contains(out, action) {
			t.Errorf("expected output to contain %q, got:\n%s", action, out)
		}
	}
}

func TestRunHistoryList_NoHistory_ReturnsError(t *testing.T) {
	storeDir, cleanup := setupHistoryStore(t)
	defer cleanup()

	var buf bytes.Buffer
	err := runHistoryList(&buf, storeDir, "nonexistent-scope")
	if err == nil {
		t.Fatal("expected error for scope with no history, got nil")
	}
}

func TestRunHistoryClear_RemovesFile(t *testing.T) {
	storeDir, cleanup := setupHistoryStore(t)
	defer cleanup()
	seedHistoryScope(t, storeDir)

	var buf bytes.Buffer
	err := runHistoryClear(&buf, storeDir, "history-scope")
	if err != nil {
		t.Fatalf("runHistoryClear: %v", err)
	}

	// Verify history file is gone
	hPath := filepath.Join(storeDir, "history-scope.history.json")
	if _, statErr := os.Stat(hPath); !os.IsNotExist(statErr) {
		t.Errorf("expected history file to be removed, but it still exists")
	}

	if !strings.Contains(buf.String(), "history-scope") {
		t.Errorf("expected output to mention scope name, got: %s", buf.String())
	}
}

func TestRunHistoryClear_OutputsMessage(t *testing.T) {
	storeDir, cleanup := setupHistoryStore(t)
	defer cleanup()
	seedHistoryScope(t, storeDir)

	var buf bytes.Buffer
	_ = runHistoryClear(&buf, storeDir, "history-scope")

	out := buf.String()
	if !strings.Contains(out, "cleared") && !strings.Contains(out, "history-scope") {
		t.Errorf("expected confirmation message, got: %s", out)
	}
}
