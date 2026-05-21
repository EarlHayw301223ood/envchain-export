// Package copy provides functionality to duplicate an existing scope
// into a new scope, optionally under a different passphrase.
package copy

import (
	"errors"
	"fmt"

	"github.com/user/envchain-export/internal/store"
)

// ErrSameName is returned when the source and destination scope names are identical.
var ErrSameName = errors.New("copy: source and destination scope names must differ")

// Copy reads the chain stored at srcScope (decrypted with srcPassphrase) and
// writes it to dstScope encrypted with dstPassphrase. If dstPassphrase is
// empty, srcPassphrase is reused. The destination scope must not already exist
// unless the caller has removed it beforehand.
func Copy(st *store.Store, srcScope, dstScope, srcPassphrase, dstPassphrase string) error {
	if srcScope == dstScope {
		return ErrSameName
	}

	if dstPassphrase == "" {
		dstPassphrase = srcPassphrase
	}

	chain, err := st.Load(srcScope, srcPassphrase)
	if err != nil {
		return fmt.Errorf("copy: load source scope %q: %w", srcScope, err)
	}

	if err := st.Save(dstScope, dstPassphrase, chain); err != nil {
		return fmt.Errorf("copy: save destination scope %q: %w", dstScope, err)
	}

	return nil
}
