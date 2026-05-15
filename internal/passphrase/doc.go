// Package passphrase provides terminal-based passphrase prompting
// utilities for the envchain-export CLI.
//
// It supports single-prompt and confirmation-prompt flows, both of
// which suppress terminal echo via the golang.org/x/term package.
//
// Typical usage:
//
//	// Prompt once (e.g. for decryption / read operations)
//	pass, err := passphrase.Prompt("Passphrase: ")
//
//	// Prompt twice (e.g. for encryption / write operations)
//	pass, err := passphrase.PromptConfirm("New passphrase: ", "Confirm passphrase: ")
package passphrase
