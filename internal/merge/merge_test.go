package merge_test

import (
	"testing"

	"github.com/user/envchain-export/internal/chain"
	"github.com/user/envchain-export/internal/merge"
)

func newChain(t *testing.T, scope string, pairs ...string) *chain.Chain {
	t.Helper()
	c, err := chain.New(scope)
	if err != nil {
		t.Fatalf("chain.New(%q): %v", scope, err)
	}
	for i := 0; i+1 < len(pairs); i += 2 {
		if err := c.Add(pairs[i], pairs[i+1]); err != nil {
			t.Fatalf("Add(%q): %v", pairs[i], err)
		}
	}
	return c
}

func TestMerge_NoConflict(t *testing.T) {
	dst := newChain(t, "dst")
	a := newChain(t, "a", "FOO", "1")
	b := newChain(t, "b", "BAR", "2")

	if err := merge.Merge(dst, merge.PolicyError, a, b); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	for key, want := range map[string]string{"FOO": "1", "BAR": "2"} {
		got, ok := dst.Get(key)
		if !ok || got != want {
			t.Errorf("Get(%q) = %q, %v; want %q, true", key, got, ok, want)
		}
	}
}

func TestMerge_PolicyError_OnDuplicate(t *testing.T) {
	dst := newChain(t, "dst", "FOO", "original")
	src := newChain(t, "src", "FOO", "new")

	if err := merge.Merge(dst, merge.PolicyError, src); err == nil {
		t.Fatal("expected error for duplicate key, got nil")
	}
}

func TestMerge_PolicyOverwrite(t *testing.T) {
	dst := newChain(t, "dst", "FOO", "original")
	src := newChain(t, "src", "FOO", "new")

	if err := merge.Merge(dst, merge.PolicyOverwrite, src); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got, _ := dst.Get("FOO")
	if got != "new" {
		t.Errorf("Get(FOO) = %q; want %q", got, "new")
	}
}

func TestMerge_PolicySkip(t *testing.T) {
	dst := newChain(t, "dst", "FOO", "original")
	src := newChain(t, "src", "FOO", "new")

	if err := merge.Merge(dst, merge.PolicySkip, src); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got, _ := dst.Get("FOO")
	if got != "original" {
		t.Errorf("Get(FOO) = %q; want %q", got, "original")
	}
}

func TestMerge_MultipleSources(t *testing.T) {
	dst := newChain(t, "dst")
	a := newChain(t, "a", "A", "1")
	b := newChain(t, "b", "B", "2")
	c := newChain(t, "c", "C", "3")

	if err := merge.Merge(dst, merge.PolicyError, a, b, c); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(dst.Keys()) != 3 {
		t.Errorf("expected 3 keys, got %d", len(dst.Keys()))
	}
}
