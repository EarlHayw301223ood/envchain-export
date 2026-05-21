package main

import (
	"errors"
	"fmt"
	"io"

	"github.com/spf13/cobra"

	envcopy "github.com/nicholasgasior/envchain-export/internal/copy"
	"github.com/nicholasgasior/envchain-export/internal/passphrase"
	"github.com/nicholasgasior/envchain-export/internal/store"
)

func newCopyCmd(storeDir string, out io.Writer) *cobra.Command {
	var newPassphrase string

	cmd := &cobra.Command{
		Use:   "copy <source-scope> <dest-scope>",
		Short: "Copy a scope to a new name",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runCopy(storeDir, args[0], args[1], newPassphrase, out)
		},
	}

	cmd.Flags().StringVar(&newPassphrase, "new-passphrase", "", "Passphrase for the destination scope (defaults to source passphrase)")
	return cmd
}

func runCopy(storeDir, src, dst, newPass string, out io.Writer) error {
	s := store.New(storeDir)

	oldPass, err := passphrase.Prompt("Source passphrase: ")
	if err != nil {
		return fmt.Errorf("reading source passphrase: %w", err)
	}

	if newPass == "" {
		newPass, err = promptNewPassphrase("Destination passphrase: ")
		if err != nil {
			if errors.Is(err, passphrase.ErrMismatch) {
				return fmt.Errorf("passphrases do not match")
			}
			return fmt.Errorf("reading destination passphrase: %w", err)
		}
	}

	if err := envcopy.Copy(s, src, dst, oldPass, newPass); err != nil {
		return fmt.Errorf("copy: %w", err)
	}

	fmt.Fprintf(out, "Scope %q copied to %q successfully.\n", src, dst)
	return nil
}
