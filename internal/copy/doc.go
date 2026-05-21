// Package copy provides functionality to duplicate an existing scope
// under a new name within the same store, optionally using a different
// passphrase for the destination scope.
//
// Usage:
//
//	err := copy.Copy(store, "source-scope", "dest-scope", "old-pass", "new-pass")
//
// If source and destination scope names are identical, Copy returns an error.
// The destination scope is created fresh; any existing data at that path
// is overwritten.
package copy
