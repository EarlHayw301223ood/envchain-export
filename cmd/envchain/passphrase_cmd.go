package main

import (
	"fmt"
	"os"

	"github.com/yourusername/envchain-export/internal/passphrase"
	"github.com/yourusername/envchain-export/internal/store"
)

// promptPassphrase reads a passphrase for an existing scope (single prompt).
func promptPassphrase(scope string) (string, error) {
	return passphrase.Prompt(fmt.Sprintf("Passphrase for scope %q: ", scope))
}

// promptNewPassphrase reads and confirms a passphrase for a new scope.
func promptNewPassphrase(scope string) (string, error) {
	return passphrase.PromptConfirm(
		fmt.Sprintf("New passphrase for scope %q: ", scope),
		"Confirm passphrase: ",
	)
}

// runPassphraseChange re-encrypts a scope under a new passphrase.
func runPassphraseChange(storeDir, scope string) error {
	currentPass, err := passphrase.Prompt("Current passphrase: ")
	if err != nil {
		return err
	}

	s, err := store.New(storeDir)
	if err != nil {
		return fmt.Errorf("opening store: %w", err)
	}

	chain, err := s.Load(scope, currentPass)
	if err != nil {
		return fmt.Errorf("loading scope: %w", err)
	}

	newPass, err := passphrase.PromptConfirm("New passphrase: ", "Confirm new passphrase: ")
	if err != nil {
		return err
	}

	if err := s.Save(scope, chain, newPass); err != nil {
		return fmt.Errorf("saving scope with new passphrase: %w", err)
	}

	fmt.Fprintf(os.Stdout, "Passphrase for scope %q updated.\n", scope)
	return nil
}
