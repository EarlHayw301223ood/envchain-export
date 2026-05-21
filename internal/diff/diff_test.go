package diff_test

import (
	"testing"

	"github.com/user/envchain-export/internal/chain"
	"github.com/user/envchain-export/internal/diff"
)

func newChain(t *testing.T, scope string, pairs map[string]string) *chain.Chain {
	t.Helper()
	c, err := chain.New(scope)
	if err != nil {
		t.Fatalf("chain.New: %v", err)
	}
	for k, v := range pairs {
		if err := c.Add(k, v); err != nil {
			t.Fatalf("chain.Add(%q): %v", k, err)
		}
	}
	return c
}

func TestDiff_NoChanges(t *testing.T) {
	base := newChain(t, "scope", map[string]string{"A": "1", "B": "2"})
	target := newChain(t, "scope", map[string]string{"A": "1", "B": "2"})

	r := diff.Diff(base, target)
	if r.HasChanges() {
		t.Errorf("expected no changes, got %d entries", len(r.Entries))
	}
}

func TestDiff_Added(t *testing.T) {
	base := newChain(t, "scope", map[string]string{"A": "1"})
	target := newChain(t, "scope", map[string]string{"A": "1", "B": "2"})

	r := diff.Diff(base, target)
	if len(r.Entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(r.Entries))
	}
	e := r.Entries[0]
	if e.Key != "B" || e.Kind != diff.Added || e.NewValue != "2" {
		t.Errorf("unexpected entry: %+v", e)
	}
}

func TestDiff_Removed(t *testing.T) {
	base := newChain(t, "scope", map[string]string{"A": "1", "B": "2"})
	target := newChain(t, "scope", map[string]string{"A": "1"})

	r := diff.Diff(base, target)
	if len(r.Entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(r.Entries))
	}
	e := r.Entries[0]
	if e.Key != "B" || e.Kind != diff.Removed || e.OldValue != "2" {
		t.Errorf("unexpected entry: %+v", e)
	}
}

func TestDiff_Changed(t *testing.T) {
	base := newChain(t, "scope", map[string]string{"A": "old"})
	target := newChain(t, "scope", map[string]string{"A": "new"})

	r := diff.Diff(base, target)
	if len(r.Entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(r.Entries))
	}
	e := r.Entries[0]
	if e.Key != "A" || e.Kind != diff.Changed || e.OldValue != "old" || e.NewValue != "new" {
		t.Errorf("unexpected entry: %+v", e)
	}
}

func TestDiff_SortedByKey(t *testing.T) {
	base := newChain(t, "scope", map[string]string{})
	target := newChain(t, "scope", map[string]string{"Z": "1", "A": "2", "M": "3"})

	r := diff.Diff(base, target)
	if len(r.Entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(r.Entries))
	}
	keys := []string{r.Entries[0].Key, r.Entries[1].Key, r.Entries[2].Key}
	if keys[0] != "A" || keys[1] != "M" || keys[2] != "Z" {
		t.Errorf("entries not sorted: %v", keys)
	}
}
