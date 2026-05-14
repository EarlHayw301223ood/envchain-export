// Package validate provides validation helpers for environment variable
// names and values used within envchain-export.
package validate

import (
	"errors"
	"regexp"
	"strings"
)

// validKeyRe matches POSIX-compliant environment variable names:
// must start with a letter or underscore, followed by letters, digits, or underscores.
var validKeyRe = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*$`)

// ErrEmptyKey is returned when an empty key is provided.
var ErrEmptyKey = errors.New("validate: key must not be empty")

// ErrInvalidKey is returned when a key contains invalid characters or
// starts with a digit.
var ErrInvalidKey = errors.New("validate: key must match [A-Za-z_][A-Za-z0-9_]*")

// ErrNullByteInValue is returned when a value contains a null byte.
var ErrNullByteInValue = errors.New("validate: value must not contain null bytes")

// Key validates that s is a legal environment variable name.
func Key(s string) error {
	if s == "" {
		return ErrEmptyKey
	}
	if !validKeyRe.MatchString(s) {
		return ErrInvalidKey
	}
	return nil
}

// Value validates that v is safe to use as an environment variable value.
// Currently the only hard restriction is the absence of null bytes, which
// cannot be represented in a C-string environment.
func Value(v string) error {
	if strings.ContainsRune(v, '\x00') {
		return ErrNullByteInValue
	}
	return nil
}

// Pair validates both the key and the value together.
func Pair(key, value string) error {
	if err := Key(key); err != nil {
		return err
	}
	return Value(value)
}
