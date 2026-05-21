// Package rename provides functionality to rename an existing scope.
package rename

import (
	"errors"
	"fmt"

	"github.com/envchain-export/internal/store"
)

// ErrSameName is returned when the source and destination scope names are identical.
var ErrSameName = errors.New("rename: source and destination scope names must differ")

// ErrDestExists is returned when the destination scope already exists.
var ErrDestExists = errors.New("rename: destination scope already exists")

// Rename loads the chain stored under oldScope (decrypted with passphrase),
// saves it under newScope (encrypted with the same passphrase), then deletes
// the old scope file.
func Rename(st *store.Store, oldScope, newScope, passphrase string) error {
	if oldScope == newScope {
		return ErrSameName
	}

	if st.Exists(newScope) {
		return fmt.Errorf("%w: %q", ErrDestExists, newScope)
	}

	chain, err := st.Load(oldScope, passphrase)
	if err != nil {
		return fmt.Errorf("rename: load %q: %w", oldScope, err)
	}

	if err := st.Save(newScope, passphrase, chain); err != nil {
		return fmt.Errorf("rename: save %q: %w", newScope, err)
	}

	if err := st.Delete(oldScope); err != nil {
		// Best-effort rollback: remove the newly created scope.
		_ = st.Delete(newScope)
		return fmt.Errorf("rename: delete old scope %q: %w", oldScope, err)
	}

	return nil
}
