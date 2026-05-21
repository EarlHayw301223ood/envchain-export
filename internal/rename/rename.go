// Package rename provides scope rename functionality for envchain-export.
package rename

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

// ErrSameName is returned when the old and new scope names are identical.
var ErrSameName = errors.New("rename: old and new scope names are identical")

// ErrDestExists is returned when the destination scope already exists and
// force is false.
var ErrDestExists = errors.New("rename: destination scope already exists")

// Storer is the minimal interface required by Rename.
type Storer interface {
	Load(scope, passphrase string) (interface{ Add(k, v string) error }, error)
	Save(scope string, c interface{ Add(k, v string) error }, passphrase string) error
	Dir() string
}

// concreteStore matches *store.Store without importing it directly.
type concreteStore interface {
	Dir() string
}

// Rename renames oldScope to newScope within s using passphrase.
// If force is false and newScope already exists, ErrDestExists is returned.
func Rename(s interface {
	Load(scope, passphrase string) (chainVal interface{ Add(k, v string) error }, err error)
	Save(scope string, c interface{ Add(k, v string) error }, passphrase string) error
	Dir() string
}, oldScope, newScope, passphrase string, force bool) error {
	if oldScope == newScope {
		return ErrSameName
	}

	destPath := filepath.Join(s.Dir(), newScope+".enc")
	if !force {
		if _, err := os.Stat(destPath); err == nil {
			return ErrDestExists
		}
	}

	c, err := s.Load(oldScope, passphrase)
	if err != nil {
		return fmt.Errorf("rename: load %q: %w", oldScope, err)
	}

	if err := s.Save(newScope, c, passphrase); err != nil {
		return fmt.Errorf("rename: save %q: %w", newScope, err)
	}

	srcPath := filepath.Join(s.Dir(), oldScope+".enc")
	if err := os.Remove(srcPath); err != nil {
		return fmt.Errorf("rename: remove old scope %q: %w", oldScope, err)
	}

	return nil
}
