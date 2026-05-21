// Package rename provides functionality to rename an existing
// envchain scope to a new name within the same store.
//
// The rename operation loads the scope under the old name using the
// provided passphrase, saves it under the new name, and removes the
// old scope file — all atomically from the caller's perspective.
package rename
