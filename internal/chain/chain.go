package chain

import (
	"errors"
	"regexp"
)

// ErrInvalidScope is returned when a scope name is invalid.
var ErrInvalidScope = errors.New("invalid scope name: must be alphanumeric with hyphens or underscores")

// ErrEmptyChain is returned when attempting to operate on an empty chain.
var ErrEmptyChain = errors.New("chain contains no environment variables")

var scopePattern = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9_-]*$`)

// EnvVar represents a single environment variable key-value pair.
type EnvVar struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// Chain holds a named, scoped set of environment variables.
type Chain struct {
	Scope string   `json:"scope"`
	Vars  []EnvVar `json:"vars"`
}

// New creates a new Chain with the given scope name.
func New(scope string) (*Chain, error) {
	if !scopePattern.MatchString(scope) {
		return nil, ErrInvalidScope
	}
	return &Chain{
		Scope: scope,
		Vars:  []EnvVar{},
	}, nil
}

// Add inserts or updates an environment variable in the chain.
func (c *Chain) Add(key, value string) {
	for i, v := range c.Vars {
		if v.Key == key {
			c.Vars[i].Value = value
			return
		}
	}
	c.Vars = append(c.Vars, EnvVar{Key: key, Value: value})
}

// Remove deletes an environment variable by key. Returns false if not found.
func (c *Chain) Remove(key string) bool {
	for i, v := range c.Vars {
		if v.Key == key {
			c.Vars = append(c.Vars[:i], c.Vars[i+1:]...)
			return true
		}
	}
	return false
}

// Get retrieves the value of an environment variable by key.
func (c *Chain) Get(key string) (string, bool) {
	for _, v := range c.Vars {
		if v.Key == key {
			return v.Value, true
		}
	}
	return "", false
}

// Len returns the number of variables stored in the chain.
func (c *Chain) Len() int {
	return len(c.Vars)
}
