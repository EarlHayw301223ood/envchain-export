package main

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
	"github.com/user/envchain-export/internal/rekey"
)

func newRekeyCmd(storeDir string, out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rekey",
		Short: "Re-encrypt all scopes with a new passphrase",
		Long: `Rekey decrypts every scope using the current passphrase and
re-encrypts each one with the new passphrase in a single pass.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runRekey(storeDir, out)
		},
	}
	return cmd
}

func runRekey(storeDir string, out io.Writer) error {
	oldPass, err := promptPassphrase("Current passphrase: ")
	if err != nil {
		return err
	}
	newPass, err := promptNewPassphrase()
	if err != nil {
		return err
	}

	results, err := rekey.Rekey(storeDir, oldPass, newPass)
	if err != nil {
		return fmt.Errorf("rekey: %w", err)
	}

	var failed int
	for _, r := range results {
		if r.Success {
			fmt.Fprintf(out, "  ✓ %s\n", r.Scope)
		} else {
			fmt.Fprintf(os.Stderr, "  ✗ %s: %v\n", r.Scope, r.Err)
			failed++
		}
	}

	if failed > 0 {
		return fmt.Errorf("%d scope(s) could not be rekeyed", failed)
	}
	fmt.Fprintf(out, "All %d scope(s) rekeyed successfully.\n", len(results))
	return nil
}
