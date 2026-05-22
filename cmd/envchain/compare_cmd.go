package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/user/envchain-export/internal/compare"
	"github.com/user/envchain-export/internal/store"
)

func newCompareCmd(st *store.Store) *cobra.Command {
	var passA, passB string

	cmd := &cobra.Command{
		Use:   "compare <scope-a> <scope-b>",
		Short: "Compare two scopes and report key differences",
		Long: `Compare two scopes side-by-side.

Symbols:
  <  key only in scope-a
  >  key only in scope-b
  ~  key present in both but values differ
  =  key present in both with identical value`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runCompare(st, args[0], args[1], passA, passB)
		},
	}

	cmd.Flags().StringVar(&passA, "pass-a", "", "passphrase for scope-a (prompted if omitted)")
	cmd.Flags().StringVar(&passB, "pass-b", "", "passphrase for scope-b (prompted if omitted)")

	return cmd
}

func runCompare(st *store.Store, scopeA, scopeB, passA, passB string) error {
	var err error
	if passA == "" {
		passA, err = promptPassphrase(fmt.Sprintf("Passphrase for %q", scopeA))
		if err != nil {
			return err
		}
	}
	if passB == "" {
		passB, err = promptPassphrase(fmt.Sprintf("Passphrase for %q", scopeB))
		if err != nil {
			return err
		}
	}

	r, err := compare.Compare(st, scopeA, scopeB, passA, passB)
	if err != nil {
		if errors.Is(err, compare.ErrSameName) {
			return fmt.Errorf("source and target scope names must differ")
		}
		return err
	}

	compare.Write(os.Stdout, scopeA, scopeB, r)

	total := len(r.OnlyInA) + len(r.OnlyInB) + len(r.Changed)
	if total == 0 {
		fmt.Fprintln(os.Stdout, "No differences found.")
	}

	return nil
}
