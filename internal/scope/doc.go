// Package scope provides helpers for managing named environment-variable
// scopes stored by envchain-export.
//
// A scope is a named collection of encrypted key/value pairs persisted as a
// single file ("<scope>.enc") inside the configured store directory.
//
// The package intentionally contains no I/O beyond the store directory so
// that it remains easy to test without touching production state.
package scope
