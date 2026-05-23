package main

import (
	"errors"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/nicholasgasior/envchain-export/internal/alias"
	"github.com/spf13/cobra"
)

func newAliasCmd(storeDir string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "alias",
		Short: "Manage scope aliases",
	}
	cmd.AddCommand(newAliasAddCmd(storeDir))
	cmd.AddCommand(newAliasRemoveCmd(storeDir))
	cmd.AddCommand(newAliasListCmd(storeDir))
	cmd.AddCommand(newAliasResolveCmd(storeDir))
	return cmd
}

func newAliasAddCmd(storeDir string) *cobra.Command {
	return &cobra.Command{
		Use:   "add <alias> <scope>",
		Short: "Create an alias pointing to a scope",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runAliasAdd(storeDir, args[0], args[1], cmd)
		},
	}
}

func runAliasAdd(storeDir, name, target string, cmd *cobra.Command) error {
	if err := alias.Add(storeDir, name, target); err != nil {
		if errors.Is(err, alias.ErrAliasExists) {
			return fmt.Errorf("alias %q already exists", name)
		}
		return err
	}
	fmt.Fprintf(cmd.OutOrStdout(), "alias %q → %q created\n", name, target)
	return nil
}

func newAliasRemoveCmd(storeDir string) *cobra.Command {
	return &cobra.Command{
		Use:   "remove <alias>",
		Short: "Remove an existing alias",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := alias.Remove(storeDir, args[0]); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "alias %q removed\n", args[0])
			return nil
		},
	}
}

func newAliasListCmd(storeDir string) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all aliases",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			list, err := alias.List(storeDir)
			if err != nil {
				return err
			}
			if len(list) == 0 {
				return fmt.Errorf("no aliases defined")
			}
			w := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 0, 2, ' ', 0)
			fmt.Fprintln(w, "ALIAS\tSCOPE")
			for k, v := range list {
				fmt.Fprintf(w, "%s\t%s\n", k, v)
			}
			return w.Flush()
		},
	}
}

func newAliasResolveCmd(storeDir string) *cobra.Command {
	return &cobra.Command{
		Use:   "resolve <alias>",
		Short: "Print the scope an alias points to",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			target, err := alias.Resolve(storeDir, args[0])
			if err != nil {
				return err
			}
			fmt.Fprintln(cmd.OutOrStdout(), target)
			return nil
		},
	}
}

var _ = os.Stderr // suppress unused import
