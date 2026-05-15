package main

import (
	"fmt"
	"io"

	"github.com/nicholasgasior/envchain-export/internal/rotate"
	"github.com/spf13/cobra"
)

// newRotateCmd builds the `passphrase rotate` sub-command.
func newRotateCmd(storeDir string, out io.Writer) *cobra.Command {
	return &cobra.Command{
		Use:   "rotate <scope>",
		Short: "Re-encrypt a scope with a new passphrase",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runRotate(storeDir, args[0], out)
		},
	}
}

func runRotate(storeDir, scope string, out io.Writer) error {
	oldPass, err := promptPassphrase("Current passphrase: ")
	if err != nil {
		return err
	}

	newPass, err := promptNewPassphrase()
	if err != nil {
		return err
	}

	if err := rotate.Rotate(storeDir, scope, oldPass, newPass); err != nil {
		return err
	}

	fmt.Fprintf(out, "passphrase rotated for scope %q\n", scope)
	return nil
}
