// Package inspect provides functionality to read and display the contents
// of an encrypted envchain scope without exporting it to a shell format.
package inspect

import (
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/user/envchain-export/internal/chain"
	"github.com/user/envchain-export/internal/store"
)

// Options controls how the inspection output is rendered.
type Options struct {
	// MaskValues replaces variable values with asterisks.
	MaskValues bool
	// SortKeys sorts the output alphabetically by key.
	SortKeys bool
}

// Inspect loads the named scope from the store and writes a human-readable
// summary of its key/value pairs to w.
func Inspect(st *store.Store, scope, passphrase string, opts Options, w io.Writer) error {
	ch, err := st.Load(scope, passphrase)
	if err != nil {
		return fmt.Errorf("inspect: %w", err)
	}

	keys := ch.Keys()
	if opts.SortKeys {
		sort.Strings(keys)
	}

	if len(keys) == 0 {
		fmt.Fprintf(w, "scope %q is empty\n", scope)
		return nil
	}

	fmt.Fprintf(w, "scope: %s (%d variable(s))\n", scope, len(keys))
	fmt.Fprintln(w, strings.Repeat("-", 40))

	for _, k := range keys {
		v, _ := ch.Get(k)
		if opts.MaskValues {
			v = maskValue(v)
		}
		fmt.Fprintf(w, "  %-24s = %s\n", k, v)
	}

	return nil
}

// maskValue replaces all but the first two characters of v with asterisks.
// If v is two characters or fewer the entire value is masked.
func maskValue(v string) string {
	if len(v) <= 2 {
		return strings.Repeat("*", len(v))
	}
	return v[:2] + strings.Repeat("*", len(v)-2)
}

// Keys returns the keys stored in ch in insertion order.
// It is a thin helper so Inspect does not depend on chain internals.
func keys(ch *chain.Chain) []string { return ch.Keys() }
