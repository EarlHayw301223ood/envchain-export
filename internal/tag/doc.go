// Package tag manages user-defined labels (tags) for envchain scopes.
//
// Tags are stored in a JSON index file within the store directory and allow
// scopes to be grouped, filtered, or annotated with arbitrary metadata such
// as environment names ("production", "staging") or cloud providers ("aws",
// "gcp").
//
// Tags are purely metadata — they have no effect on encryption or the
// key-value contents of a scope.
package tag
