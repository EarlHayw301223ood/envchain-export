// Package rotate provides functionality to re-encrypt an existing scope
// with a new passphrase without losing any stored variables.
package rotate

import (
	"fmt"

	"github.com/nicholasgasior/envchain-export/internal/store"
)

// ErrSamePassphrase is returned when the new passphrase matches the old one.
var ErrSamePassphrase = fmt.Errorf("rotate: new passphrase must differ from the current passphrase")

// Rotate loads a scope with oldPass, then saves it again with newPass,
// effectively re-encrypting all stored variables under the new passphrase.
//
// The operation is atomic from the caller's perspective: if loading fails
// (e.g. wrong passphrase) the stored file is never touched.
func Rotate(storeDir, scope, oldPass, newPass string) error {
	if oldPass == newPass {
		return ErrSamePassphrase
	}

	s := store.New(storeDir)

	ch, err := s.Load(scope, oldPass)
	if err != nil {
		return fmt.Errorf("rotate: load scope %q: %w", scope, err)
	}

	if err := s.Save(scope, newPass, ch); err != nil {
		return fmt.Errorf("rotate: save scope %q: %w", scope, err)
	}

	return nil
}
