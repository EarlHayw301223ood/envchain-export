// Package passphrase provides utilities for securely prompting
// and confirming passphrases from the terminal.
package passphrase

import (
	"errors"
	"fmt"
	"os"

	"golang.org/x/term"
)

// ErrMismatch is returned when two passphrase entries do not match.
var ErrMismatch = errors.New("passphrases do not match")

// ErrEmpty is returned when an empty passphrase is provided.
var ErrEmpty = errors.New("passphrase must not be empty")

// Prompt reads a passphrase from the terminal without echo.
// The prompt string is written to stderr before reading.
func Prompt(prompt string) (string, error) {
	fmt.Fprint(os.Stderr, prompt)
	bytes, err := term.ReadPassword(int(os.Stdin.Fd()))
	fmt.Fprintln(os.Stderr)
	if err != nil {
		return "", fmt.Errorf("reading passphrase: %w", err)
	}
	if len(bytes) == 0 {
		return "", ErrEmpty
	}
	return string(bytes), nil
}

// PromptConfirm prompts for a passphrase twice and returns it only
// if both entries match. Returns ErrMismatch if they differ.
func PromptConfirm(prompt, confirmPrompt string) (string, error) {
	first, err := Prompt(prompt)
	if err != nil {
		return "", err
	}
	second, err := Prompt(confirmPrompt)
	if err != nil {
		return "", err
	}
	if first != second {
		return "", ErrMismatch
	}
	return first, nil
}
