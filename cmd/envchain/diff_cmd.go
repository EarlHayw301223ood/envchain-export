package main

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
	"github.com/user/envchain-export/internal/diff"
	"github.com/user/envchain-export/internal/passphrase"
	"github.com/user/envchain-export/internal/store"
)

func newDiffCmd(storeDir string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "diff <scope-a> <scope-b>",
		Short: "Show differences between two scopes",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runDiff(os.Stdout, storeDir, args[0], args[1])
		},
	}
	return cmd
}

func runDiff(w io.Writer, storeDir, scopeA, scopeB string) error {
	passA, err := passphrase.Prompt(fmt.Sprintf("Passphrase for scope %q: ", scopeA))
	if err != nil {
		return fmt.Errorf("passphrase for %q: %w", scopeA, err)
	}

	passB, err := passphrase.Prompt(fmt.Sprintf("Passphrase for scope %q: ", scopeB))
	if err != nil {
		return fmt.Errorf("passphrase for %q: %w", scopeB, err)
	}

	stA := store.New(storeDir)
	chainA, err := stA.Load(scopeA, passA)
	if err != nil {
		return fmt.Errorf("load scope %q: %w", scopeA, err)
	}

	stB := store.New(storeDir)
	chainB, err := stB.Load(scopeB, passB)
	if err != nil {
		return fmt.Errorf("load scope %q: %w", scopeB, err)
	}

	result := diff.Diff(chainA, chainB)
	if !result.HasChanges() {
		fmt.Fprintln(w, "No differences found.")
		return nil
	}

	fmt.Fprintf(w, "Diff: %s → %s\n", scopeA, scopeB)
	for _, e := range result.Entries {
		switch e.Kind {
		case diff.Added:
			fmt.Fprintf(w, "  + %s=%s\n", e.Key, e.NewValue)
		case diff.Removed:
			fmt.Fprintf(w, "  - %s=%s\n", e.Key, e.OldValue)
		case diff.Changed:
			fmt.Fprintf(w, "  ~ %s: %s → %s\n", e.Key, e.OldValue, e.NewValue)
		}
	}
	return nil
}
