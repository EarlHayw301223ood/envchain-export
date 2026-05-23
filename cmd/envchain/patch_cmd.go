package main

import (
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/spf13/cobra"
	"github.com/user/envchain-export/internal/patch"
	"github.com/user/envchain-export/internal/store"
)

func newPatchCmd(st *store.Store, out io.Writer) *cobra.Command {
	var (
		setFlags   []string
		unsetFlags []string
	)

	cmd := &cobra.Command{
		Use:   "patch <scope>",
		Short: "Apply key mutations to a scope without a full rewrite",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runPatch(st, out, args[0], setFlags, unsetFlags)
		},
	}

	cmd.Flags().StringArrayVarP(&setFlags, "set", "s", nil, "KEY=VALUE pair to set (repeatable)")
	cmd.Flags().StringArrayVarP(&unsetFlags, "unset", "u", nil, "KEY to remove (repeatable)")
	return cmd
}

func runPatch(st *store.Store, out io.Writer, scope string, setFlags, unsetFlags []string) error {
	var setPairs [][2]string
	for _, s := range setFlags {
		parts := strings.SplitN(s, "=", 2)
		if len(parts) != 2 {
			return fmt.Errorf("--set %q: expected KEY=VALUE", s)
		}
		setPairs = append(setPairs, [2]string{parts[0], parts[1]})
	}

	ops := patch.BuildOps(setPairs, unsetFlags)
	if len(ops) == 0 {
		return errors.New("patch: provide at least one --set or --unset flag")
	}

	pass, err := promptPassphrase(fmt.Sprintf("Passphrase for scope %q: ", scope))
	if err != nil {
		return err
	}

	if err := patch.Patch(st, scope, pass, ops); err != nil {
		return err
	}

	setCount := len(setPairs)
	unsetCount := len(unsetFlags)
	fmt.Fprintf(out, "Patched scope %q: %d set, %d unset.\n", scope, setCount, unsetCount)
	return nil
}
