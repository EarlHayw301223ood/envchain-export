// main.go is the entry point for the envchain CLI tool.
// It wires together all subcommands and global flags.
package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

const (
	// defaultStoreDir is the default directory where encrypted chain files are stored.
	defaultStoreDir = ".envchain"
)

// globalFlags holds values shared across subcommands.
type globalFlags struct {
	storeDir string
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	flags := &globalFlags{}

	rootCmd := &cobra.Command{
		Use:   "envchain",
		Short: "Safely export and import scoped environment variable sets",
		Long: `envchain stores named sets (scopes) of environment variables,
encrypted at rest, and can export them in POSIX or dotenv format.`,
		SilenceUsage: true,
	}

	// Resolve default store directory relative to user's home.
	homeDir, err := os.UserHomeDir()
	if err != nil {
		homeDir = "."
	}
	defaultDir := filepath.Join(homeDir, defaultStoreDir)

	rootCmd.PersistentFlags().StringVar(
		&flags.storeDir,
		"store",
		defaultDir,
		"directory used to persist encrypted chain files",
	)

	// Register subcommand groups.
	rootCmd.AddCommand(newSetCmd(flags))
	rootCmd.AddCommand(newGetCmd(flags))
	rootCmd.AddCommand(newDeleteCmd(flags))
	rootCmd.AddCommand(newExportCmd(flags))
	rootCmd.AddCommand(newImportCmd(flags))
	rootCmd.AddCommand(newScopeCmd(flags))
	rootCmd.AddCommand(newPassphraseCmd(flags))

	return rootCmd.Execute()
}

// newScopeCmd returns the "scope" subcommand with its children.
func newScopeCmd(flags *globalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "scope",
		Short: "Manage scopes (list, delete, check existence)",
	}

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List all available scopes",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runScopeList(flags.storeDir, cmd.OutOrStdout())
		},
	}

	var forceDelete bool
	delCmd := &cobra.Command{
		Use:   "delete <scope>",
		Short: "Delete a scope and its stored variables",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runScopeDelete(flags.storeDir, args[0], forceDelete, cmd.OutOrStdout())
		},
	}
	delCmd.Flags().BoolVarP(&forceDelete, "force", "f", false, "skip confirmation prompt")

	existsCmd := &cobra.Command{
		Use:   "exists <scope>",
		Short: "Exit 0 if a scope exists, 1 otherwise",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runScopeExists(flags.storeDir, args[0], cmd.OutOrStdout())
		},
	}

	cmd.AddCommand(listCmd, delCmd, existsCmd)
	return cmd
}

// newPassphraseCmd returns the "passphrase" subcommand.
func newPassphraseCmd(flags *globalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "passphrase <scope>",
		Short: "Change the passphrase for a scope",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runPassphraseChange(flags.storeDir, args[0])
		},
	}
	return cmd
}
