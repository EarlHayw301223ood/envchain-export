package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/user/envchain-export/internal/env"
	"github.com/user/envchain-export/internal/store"
)

func newEnvCmd(storeDir string) *cobra.Command {
	var passphrase string

	cmd := &cobra.Command{
		Use:   "env <scope> -- <command> [args...]",
		Short: "Run a command with a scope's variables injected into its environment",
		Args:  cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			scope := args[0]
			command := args[1:]
			return runEnv(storeDir, scope, passphrase, command, cmd)
		},
	}

	cmd.Flags().StringVarP(&passphrase, "passphrase", "p", "", "passphrase for the scope (prompted if omitted)")
	return cmd
}

func runEnv(storeDir, scope, passphrase string, command []string, cmd *cobra.Command) error {
	if passphrase == "" {
		var err error
		passphrase, err = promptPassphrase(fmt.Sprintf("Passphrase for scope %q: ", scope), cmd)
		if err != nil {
			return err
		}
	}

	s := store.New(storeDir)
	ch, err := s.Load(scope, passphrase)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return err
	}

	if err := env.Run(ch, command); err != nil {
		return err
	}
	return nil
}
