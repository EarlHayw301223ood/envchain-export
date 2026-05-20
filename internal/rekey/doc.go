// Package rekey provides bulk passphrase rotation across all scopes
// stored in a given store directory.
//
// Unlike rotate.Rotate, which operates on a single scope, rekey.Rekey
// iterates every scope and attempts to re-encrypt each one from the old
// passphrase to the new passphrase, collecting per-scope results so the
// caller can report partial failures without aborting the entire operation.
package rekey
