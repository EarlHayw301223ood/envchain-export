// Package snapshot provides point-in-time capture and restore of scoped
// environment variable sets.
//
// Snapshots are stored as specially-named scopes within the same store,
// using the naming convention __snapshot__<scope>__<label>. The label is
// a UTC timestamp in the format 20060102T150405Z.
//
// Example usage:
//
//	label, err := snapshot.Take(st, "myapp", passphrase)
//	// label == "20240115T103045Z"
//
//	err = snapshot.Restore(st, "myapp", label, passphrase)
package snapshot
