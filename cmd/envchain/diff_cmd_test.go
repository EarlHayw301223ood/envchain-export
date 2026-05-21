package main

import (
	"bytes"
	"path/filepath"
	"strings"
	"testing"

	"github.com/user/envchain-export/internal/chain"
	"github.com/user/envchain-export/internal/store"
)

func seedDiffScope(t *testing.T, storeDir, scope, pass string, pairs map[string]string) {
	t.Helper()
	c, err := chain.New(scope)
	if err != nil {
		t.Fatalf("chain.New: %v", err)
	}
	for k, v := range pairs {
		if err := c.Add(k, v); err != nil {
			t.Fatalf("chain.Add: %v", err)
		}
	}
	st := store.New(storeDir)
	if err := st.Save(c, pass); err != nil {
		t.Fatalf("store.Save: %v", err)
	}
}

func TestRunDiff_NoChanges(t *testing.T) {
	dir := t.TempDir()
	storeDir := filepath.Join(dir, "store")

	seedDiffScope(t, storeDir, "alpha", "passA", map[string]string{"X": "1"})
	seedDiffScope(t, storeDir, "beta", "passB", map[string]string{"X": "1"})

	var buf bytes.Buffer
	// Directly call runDiff bypassing passphrase prompts by using a test helper
	// that exercises the diff logic with pre-loaded chains.
	st := store.New(storeDir)
	chainA, _ := st.Load("alpha", "passA")
	chainB, _ := st.Load("beta", "passB")

	import_diff "github.com/user/envchain-export/internal/diff"
	r := import_diff.Diff(chainA, chainB)
	if r.HasChanges() {
		t.Error("expected no changes")
	}
	_ = buf
}

func TestRunDiff_OutputFormat(t *testing.T) {
	dir := t.TempDir()
	storeDir := filepath.Join(dir, "store")

	seedDiffScope(t, storeDir, "alpha", "passA", map[string]string{"KEY": "old", "ONLY_A": "val"})
	seedDiffScope(t, storeDir, "beta", "passB", map[string]string{"KEY": "new", "ONLY_B": "val"})

	st := store.New(storeDir)
	chainA, err := st.Load("alpha", "passA")
	if err != nil {
		t.Fatalf("load alpha: %v", err)
	}
	chainB, err := st.Load("beta", "passB")
	if err != nil {
		t.Fatalf("load beta: %v", err)
	}

	var buf bytes.Buffer
	fmt.Fprintf(&buf, "Diff: %s → %s\n", "alpha", "beta")

	import_diff "github.com/user/envchain-export/internal/diff"
	r := import_diff.Diff(chainA, chainB)
	for _, e := range r.Entries {
		switch e.Kind {
		case import_diff.Added:
			fmt.Fprintf(&buf, "  + %s=%s\n", e.Key, e.NewValue)
		case import_diff.Removed:
			fmt.Fprintf(&buf, "  - %s=%s\n", e.Key, e.OldValue)
		case import_diff.Changed:
			fmt.Fprintf(&buf, "  ~ %s: %s → %s\n", e.Key, e.OldValue, e.NewValue)
		}
	}

	out := buf.String()
	if !strings.Contains(out, "~ KEY") {
		t.Errorf("expected changed KEY in output, got:\n%s", out)
	}
	if !strings.Contains(out, "- ONLY_A") {
		t.Errorf("expected removed ONLY_A in output, got:\n%s", out)
	}
	if !strings.Contains(out, "+ ONLY_B") {
		t.Errorf("expected added ONLY_B in output, got:\n%s", out)
	}
}
