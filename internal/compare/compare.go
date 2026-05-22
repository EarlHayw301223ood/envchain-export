// Package compare provides functionality to compare two scopes side-by-side,
// reporting keys that are present in one but not the other, and keys whose
// values differ between the two scopes.
package compare

import (
	"errors"
	"fmt"
	"io"
	"sort"

	"github.com/user/envchain-export/internal/chain"
	"github.com/user/envchain-export/internal/store"
)

// ErrSameName is returned when both scope names are identical.
var ErrSameName = errors.New("compare: source and target scope names are identical")

// Result holds the outcome of comparing two scopes.
type Result struct {
	OnlyInA  []string          // keys present only in scope A
	OnlyInB  []string          // keys present only in scope B
	Changed  []string          // keys present in both but with different values
	Identical []string         // keys present in both with the same value
}

// Compare loads two scopes from the store and compares their contents.
func Compare(st *store.Store, scopeA, scopeB, passA, passB string) (*Result, error) {
	if scopeA == scopeB {
		return nil, ErrSameName
	}

	cA, err := st.Load(scopeA, passA)
	if err != nil {
		return nil, fmt.Errorf("compare: load %q: %w", scopeA, err)
	}

	cB, err := st.Load(scopeB, passB)
	if err != nil {
		return nil, fmt.Errorf("compare: load %q: %w", scopeB, err)
	}

	return compareChains(cA, cB), nil
}

func compareChains(a, b *chain.Chain) *Result {
	mapA := toMap(a)
	mapB := toMap(b)

	r := &Result{}

	for k, va := range mapA {
		if vb, ok := mapB[k]; ok {
			if va == vb {
				r.Identical = append(r.Identical, k)
			} else {
				r.Changed = append(r.Changed, k)
			}
		} else {
			r.OnlyInA = append(r.OnlyInA, k)
		}
	}

	for k := range mapB {
		if _, ok := mapA[k]; !ok {
			r.OnlyInB = append(r.OnlyInB, k)
		}
	}

	sort.Strings(r.OnlyInA)
	sort.Strings(r.OnlyInB)
	sort.Strings(r.Changed)
	sort.Strings(r.Identical)

	return r
}

func toMap(c *chain.Chain) map[string]string {
	m := make(map[string]string)
	for _, k := range c.Keys() {
		v, _ := c.Get(k)
		m[k] = v
	}
	return m
}

// Write renders a human-readable comparison report to w.
func Write(w io.Writer, scopeA, scopeB string, r *Result) {
	fmt.Fprintf(w, "Comparing %q → %q\n", scopeA, scopeB)
	for _, k := range r.OnlyInA {
		fmt.Fprintf(w, "  < %s\n", k)
	}
	for _, k := range r.OnlyInB {
		fmt.Fprintf(w, "  > %s\n", k)
	}
	for _, k := range r.Changed {
		fmt.Fprintf(w, "  ~ %s\n", k)
	}
	for _, k := range r.Identical {
		fmt.Fprintf(w, "  = %s\n", k)
	}
}
