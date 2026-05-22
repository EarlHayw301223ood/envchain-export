// Package clone provides functionality to deep-copy a scope into a new scope,
// optionally re-encrypting with a different passphrase.
package clone

import (
	"fmt"

	"github.com/nicholasgasior/envchain-export/internal/store"
)

// Clone copies all key-value pairs from srcScope into dstScope.
// srcPass is used to decrypt the source; dstPass is used to encrypt the
// destination. They may be identical or different.
// Returns an error if srcScope does not exist, dstScope already exists,
// or the source passphrase is wrong.
func Clone(st *store.Store, srcScope, dstScope, srcPass, dstPass string) error {
	if srcScope == dstScope {
		return fmt.Errorf("clone: source and destination scope must differ: %q", srcScope)
	}

	exists, err := st.Exists(dstScope)
	if err != nil {
		return fmt.Errorf("clone: checking destination scope: %w", err)
	}
	if exists {
		return fmt.Errorf("clone: destination scope already exists: %q", dstScope)
	}

	chain, err := st.Load(srcScope, srcPass)
	if err != nil {
		return fmt.Errorf("clone: loading source scope %q: %w", srcScope, err)
	}

	if err := st.Save(dstScope, dstPass, chain); err != nil {
		return fmt.Errorf("clone: saving destination scope %q: %w", dstScope, err)
	}

	return nil
}
