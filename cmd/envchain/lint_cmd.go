package main

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
	"github.com/user/envchain-export/internal/lint"
	"github.com/user/envchain-export/internal/store"
)

func newLintCmd(storeDir string) *cobra.Command {
	var passphrase string

	cmd := &cobra.Command{
		Use:   "lint <scope>",
		Short: "Check a scope for suspicious or unsafe variable entries",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if passphrase == "" {
				var err error
				passphrase, err = promptPassphrase()
				if err != nil {
					return err
				}
			}
			return runLint(cmd.OutOrStdout(), storeDir, args[0], passphrase)
		},
	}

	cmd.Flags().StringVarP(&passphrase, "passphrase", "p", "", "passphrase for the scope")
	return cmd
}

func runLint(w io.Writer, storeDir, scope, passphrase string) error {
	s := store.New(storeDir)
	c, err := s.Load(scope, passphrase)
	if err != nil {
		return fmt.Errorf("load scope %q: %w", scope, err)
	}

	findings := lint.Lint(c)
	if len(findings) == 0 {
		fmt.Fprintf(w, "scope %q: no issues found\n", scope)
		return nil
	}

	fmt.Fprintf(w, "scope %q: %d issue(s) found\n", scope, len(findings))
	for _, f := range findings {
		fmt.Fprintln(w, " ", f.String())
	}

	// Exit with a non-zero status when errors (not just warnings) are present.
	for _, f := range findings {
		if f.Severity == lint.Error {
			os.Exit(1)
		}
	}
	return nil
}
