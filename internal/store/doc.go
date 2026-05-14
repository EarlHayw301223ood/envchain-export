// Package store provides encrypted persistence for envchain Chain objects.
//
// Each Chain is serialized to JSON, encrypted with AES-GCM using a
// caller-supplied passphrase, and written to a file named after the
// chain's scope (e.g. "myapp.enc") inside a configurable directory.
//
// Example usage:
//
//	s := store.New(".envchain")
//
//	// Save
//	if err := s.Save(myChain, "my-passphrase"); err != nil {
//		log.Fatal(err)
//	}
//
//	// Load
//	loaded, err := s.Load("myapp", "my-passphrase")
//	if err != nil {
//		log.Fatal(err)
//	}
package store
