package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/user/envchain-export/internal/rename"
	"github.com/user/envchain-export/internal/store"
)

func newRenameCmd(storeDir string) *cobra.Command {
	var forceFlag bool

	cmd := &cobra.Command{
		Use:   "rename <old-scope> <new-scope>",
		Short: "Rename an existing scope",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runRename(storeDir, args[0], args[1], forceFlag)
		},
	}

	cmd.Flags().BoolVarP(&forceFlag, "force", "f", false, "overwrite destination scope if it exists")
	return cmd
}

func runRename(storeDir, oldScope, newScope string, force bool) error {
	passphrase, err := promptPassphrase("Passphrase: ")
	if err != nil {
		return err
	}

	s := store.New(storeDir)

	err = rename.Rename(s, oldScope, newScope, passphrase, force)
	if err != nil {
		if errors.Is(err, rename.ErrDestExists) {
			return fmt.Errorf("scope %q already exists; use --force to overwrite", newScope)
		}
		if errors.Is(err, rename.ErrSameName) {
			return fmt.Errorf("old and new scope names are identical: %q", oldScope)
		}
		return err
	}

	fmt.Fprintf(os.Stdout, "Renamed scope %q → %q\n", oldScope, newScope)
	return nil
}
