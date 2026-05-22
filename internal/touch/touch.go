// Package touch provides functionality to update the passphrase of a scope
// without changing any of its stored values. This is useful for rotating
// credentials on a schedule without altering the underlying data.
package touch

import (
	"fmt"

	"github.com/envchain-export/internal/store"
)

// ErrSamePassphrase is returned when the new passphrase is identical to the old one.
var ErrSamePassphrase = fmt.Errorf("new passphrase must differ from the old passphrase")

// Touch re-encrypts the given scope with newPassphrase after verifying oldPassphrase.
// The stored key-value pairs remain unchanged.
func Touch(st *store.Store, scope, oldPassphrase, newPassphrase string) error {
	if oldPassphrase == newPassphrase {
		return ErrSamePassphrase
	}

	chain, err := st.Load(scope, oldPassphrase)
	if err != nil {
		return fmt.Errorf("touch: load scope %q: %w", scope, err)
	}

	if err := st.Save(scope, chain, newPassphrase); err != nil {
		return fmt.Errorf("touch: save scope %q: %w", scope, err)
	}

	return nil
}
