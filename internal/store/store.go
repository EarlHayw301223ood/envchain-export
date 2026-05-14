package store

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/user/envchain-export/internal/chain"
	"github.com/user/envchain-export/internal/crypto"
)

const defaultDir = ".envchain"

// Store handles persisting and loading encrypted chain files.
type Store struct {
	dir string
}

// New creates a Store rooted at dir. If dir is empty, defaultDir is used.
func New(dir string) *Store {
	if dir == "" {
		dir = defaultDir
	}
	return &Store{dir: dir}
}

// Save encrypts and writes the chain to disk.
func (s *Store) Save(c *chain.Chain, passphrase string) error {
	if err := os.MkdirAll(s.dir, 0700); err != nil {
		return fmt.Errorf("store: mkdir: %w", err)
	}

	data, err := json.Marshal(c)
	if err != nil {
		return fmt.Errorf("store: marshal: %w", err)
	}

	encrypted, err := crypto.Encrypt(passphrase, data)
	if err != nil {
		return fmt.Errorf("store: encrypt: %w", err)
	}

	path := s.filePath(c.Scope())
	if err := os.WriteFile(path, encrypted, 0600); err != nil {
		return fmt.Errorf("store: write: %w", err)
	}
	return nil
}

// Load decrypts and reads the chain for the given scope.
func (s *Store) Load(scope, passphrase string) (*chain.Chain, error) {
	path := s.filePath(scope)

	encrypted, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("store: read: %w", err)
	}

	data, err := crypto.Decrypt(passphrase, encrypted)
	if err != nil {
		return nil, fmt.Errorf("store: decrypt: %w", err)
	}

	var c chain.Chain
	if err := json.Unmarshal(data, &c); err != nil {
		return nil, fmt.Errorf("store: unmarshal: %w", err)
	}
	return &c, nil
}

// filePath returns the file path for the given scope.
func (s *Store) filePath(scope string) string {
	return filepath.Join(s.dir, scope+".enc")
}
