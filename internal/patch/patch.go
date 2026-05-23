// Package patch applies a set of key mutations (set, unset) to an existing
// scope without requiring a full rewrite of the chain from scratch.
package patch

import (
	"errors"
	"fmt"

	"github.com/user/envchain-export/internal/chain"
	"github.com/user/envchain-export/internal/store"
	"github.com/user/envchain-export/internal/validate"
)

// ErrNoOps is returned when Patch is called with no operations.
var ErrNoOps = errors.New("patch: no operations provided")

// Op represents a single patch operation.
type Op struct {
	Key   string
	Value string // empty string + Unset==true means delete
	Unset bool
}

// Patch applies ops to the named scope, loading and saving via st.
// All keys are validated before any mutation is applied.
func Patch(st *store.Store, scope, passphrase string, ops []Op) error {
	if len(ops) == 0 {
		return ErrNoOps
	}

	// Validate all keys (and values for set ops) up front.
	for _, op := range ops {
		if err := validate.Key(op.Key); err != nil {
			return fmt.Errorf("patch: invalid key %q: %w", op.Key, err)
		}
		if !op.Unset {
			if err := validate.Value(op.Value); err != nil {
				return fmt.Errorf("patch: invalid value for key %q: %w", op.Key, err)
			}
		}
	}

	ch, err := st.Load(scope, passphrase)
	if err != nil {
		return fmt.Errorf("patch: load %q: %w", scope, err)
	}

	for _, op := range ops {
		if op.Unset {
			ch.Remove(op.Key)
		} else {
			if err := ch.Add(op.Key, op.Value); err != nil {
				return fmt.Errorf("patch: add key %q: %w", op.Key, err)
			}
		}
	}

	if err := st.Save(scope, passphrase, ch); err != nil {
		return fmt.Errorf("patch: save %q: %w", scope, err)
	}
	return nil
}

// BuildOps is a convenience helper that converts parallel slices of set-pairs
// and unset-keys into an Op slice.
func BuildOps(setPairs [][2]string, unsetKeys []string) []Op {
	ops := make([]Op, 0, len(setPairs)+len(unsetKeys))
	for _, kv := range setPairs {
		ops = append(ops, Op{Key: kv[0], Value: kv[1]})
	}
	for _, k := range unsetKeys {
		ops = append(ops, Op{Key: k, Unset: true})
	}
	return ops
}
