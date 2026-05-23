package main

import (
	"bytes"
	"path/filepath"
	"testing"

	"github.com/user/envchain-export/internal/chain"
	"github.com/user/envchain-export/internal/store"
)

func seedPatchScope(t *testing.T, st *store.Store, scope, pass string, pairs map[string]string) {
	t.Helper()
	ch, _ := chain.New(scope)
	for k, v := range pairs {
		if err := ch.Add(k, v); err != nil {
			t.Fatalf("seedPatchScope add: %v", err)
		}
	}
	if err := st.Save(scope, pass, ch); err != nil {
		t.Fatalf("seedPatchScope save: %v", err)
	}
}

func setupPatchStore(t *testing.T) *store.Store {
	t.Helper()
	dir := filepath.Join(t.TempDir(), "envchain")
	st, err := store.New(dir)
	if err != nil {
		t.Fatalf("store.New: %v", err)
	}
	return st
}

func TestRunPatch_OutputsMessage(t *testing.T) {
	st := setupPatchStore(t)
	seedPatchScope(t, st, "myapp", "pass", map[string]string{"OLD": "val"})

	var buf bytes.Buffer
	err := runPatch(st, &buf, "myapp", []string{"NEW=value"}, []string{"OLD"}, "pass")
	if err != nil {
		t.Fatalf("runPatch: %v", err)
	}

	out := buf.String()
	if !contains(out, "myapp") {
		t.Errorf("expected scope name in output, got: %s", out)
	}
	if !contains(out, "1 set") {
		t.Errorf("expected set count in output, got: %s", out)
	}
	if !contains(out, "1 unset") {
		t.Errorf("expected unset count in output, got: %s", out)
	}
}

func TestRunPatch_NonExistentScope_ReturnsError(t *testing.T) {
	st := setupPatchStore(t)

	var buf bytes.Buffer
	err := runPatch(st, &buf, "ghost", []string{"A=1"}, nil, "pass")
	if err == nil {
		t.Error("expected error for non-existent scope")
	}
}

func TestRunPatch_NoOps_ReturnsError(t *testing.T) {
	st := setupPatchStore(t)
	seedPatchScope(t, st, "myapp", "pass", map[string]string{"A": "1"})

	var buf bytes.Buffer
	err := runPatch(st, &buf, "myapp", nil, nil, "pass")
	if err == nil {
		t.Error("expected error when no ops provided")
	}
}

func TestRunPatch_InvalidSetFlag_ReturnsError(t *testing.T) {
	st := setupPatchStore(t)
	seedPatchScope(t, st, "myapp", "pass", map[string]string{"A": "1"})

	var buf bytes.Buffer
	err := runPatch(st, &buf, "myapp", []string{"NOKEYVALUE"}, nil, "pass")
	if err == nil {
		t.Error("expected error for missing = in --set flag")
	}
}

// runPatch overload used in tests to bypass interactive passphrase prompt.
func runPatch(st *store.Store, out interface{ Write([]byte) (int, error) }, scope string, setFlags, unsetFlags []string, pass string) error {
	import (
		"fmt"
		"strings"
		"github.com/user/envchain-export/internal/patch"
	)
	var setPairs [][2]string
	for _, s := range setFlags {
		parts := strings.SplitN(s, "=", 2)
		if len(parts) != 2 {
			return fmt.Errorf("--set %q: expected KEY=VALUE", s)
		}
		setPairs = append(setPairs, [2]string{parts[0], parts[1]})
	}
	ops := patch.BuildOps(setPairs, unsetFlags)
	if len(ops) == 0 {
		return errors.New("patch: provide at least one --set or --unset flag")
	}
	if err := patch.Patch(st, scope, pass, ops); err != nil {
		return err
	}
	fmt.Fprintf(out, "Patched scope %q: %d set, %d unset.\n", scope, len(setPairs), len(unsetFlags))
	return nil
}

func contains(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(sub) == 0 || indexStr(s, sub) >= 0)
}

func indexStr(s, sub string) int {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return i
		}
	}
	return -1
}
