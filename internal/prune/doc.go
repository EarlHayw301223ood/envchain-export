// Package prune provides utilities for bulk-removing environment variable
// entries from a scope based on key prefix or glob pattern matching.
//
// ByPrefix removes all keys that begin with a given string prefix, which is
// useful for cleaning up namespaced variables (e.g. all keys starting with
// "AWS_" or "OLD_").
//
// ByGlob removes all keys whose names match a shell-style glob pattern using
// filepath.Match semantics (e.g. "DB_*_HOST" or "LEGACY_??_KEY").
//
// Both functions return the list of keys that were actually removed so callers
// can report or audit the changes.
package prune
