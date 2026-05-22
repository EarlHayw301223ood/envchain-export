package main

import (
	"fmt"

	"github.com/envchain-export/internal/touch"
	"github.com/spf13/cobra"
)

func newTouchCmd(cfg *config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "touch <scope>",
		Short: "Re-encrypt a scope with a new passphrase without altering its values",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runTouch(cfg, args[0])
		},
	}
	return cmd
}

func runTouch(cfg *config, scope string) error {
	oldPass, err := promptPassphrase(fmt.Sprintf("Current passphrase for scope %q", scope))
	if err != nil {
		return err
	}

	newPass, err := promptNewPassphrase()
	if err != nil {
		return err
	}

	if err := touch.Touch(cfg.store, scope, oldPass, newPass); err != nil {
		return err
	}

	fmt.Fprintf(cfg.stdout, "scope %q re-encrypted with new passphrase\n", scope)
	return nil
}
