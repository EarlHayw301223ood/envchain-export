package main

import (
	"fmt"

	"github.com/nicholasgasior/envchain-export/internal/clone"
	"github.com/nicholasgasior/envchain-export/internal/store"
	"github.com/spf13/cobra"
)

func newCloneCmd(st *store.Store) *cobra.Command {
	var diffPass bool

	cmd := &cobra.Command{
		Use:   "clone <src-scope> <dst-scope>",
		Short: "Clone a scope into a new scope, optionally with a new passphrase",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runClone(cmd, st, args[0], args[1], diffPass)
		},
	}

	cmd.Flags().BoolVar(&diffPass, "new-passphrase", false,
		"Encrypt the cloned scope with a different passphrase")

	return cmd
}

func runClone(cmd *cobra.Command, st *store.Store, src, dst string, diffPass bool) error {
	srcPass, err := promptPassphrase(cmd, fmt.Sprintf("Passphrase for source scope %q", src))
	if err != nil {
		return err
	}

	dstPass := srcPass
	if diffPass {
		dstPass, err = promptNewPassphrase(cmd, fmt.Sprintf("New passphrase for cloned scope %q", dst))
		if err != nil {
			return err
		}
	}

	if err := clone.Clone(st, src, dst, srcPass, dstPass); err != nil {
		return err
	}

	fmt.Fprintf(cmd.OutOrStdout(), "Cloned scope %q → %q\n", src, dst)
	return nil
}
