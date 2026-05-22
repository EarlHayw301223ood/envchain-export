package main

import (
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/nicholasgasior/envchain-export/internal/pin"
	"github.com/spf13/cobra"
)

func newPinCmd(storeDir string, out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pin",
		Short: "Manage pinned scopes",
	}
	cmd.AddCommand(newPinAddCmd(storeDir, out))
	cmd.AddCommand(newPinRemoveCmd(storeDir, out))
	cmd.AddCommand(newPinListCmd(storeDir, out))
	return cmd
}

func newPinAddCmd(storeDir string, out io.Writer) *cobra.Command {
	return &cobra.Command{
		Use:   "add <scope>",
		Short: "Pin a scope to protect it from deletion",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runPinAdd(storeDir, args[0], out)
		},
	}
}

func newPinRemoveCmd(storeDir string, out io.Writer) *cobra.Command {
	return &cobra.Command{
		Use:   "remove <scope>",
		Short: "Unpin a scope",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runPinRemove(storeDir, args[0], out)
		},
	}
}

func newPinListCmd(storeDir string, out io.Writer) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all pinned scopes",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runPinList(storeDir, out)
		},
	}
}

func runPinAdd(storeDir, scope string, out io.Writer) error {
	if err := pin.Add(storeDir, scope); err != nil {
		if errors.Is(err, pin.ErrAlreadyPinned) {
			return fmt.Errorf("scope %q is already pinned", scope)
		}
		return err
	}
	fmt.Fprintf(out, "Pinned scope %q.\n", scope)
	return nil
}

func runPinRemove(storeDir, scope string, out io.Writer) error {
	if err := pin.Remove(storeDir, scope); err != nil {
		if errors.Is(err, pin.ErrNotPinned) {
			return fmt.Errorf("scope %q is not pinned", scope)
		}
		return err
	}
	fmt.Fprintf(out, "Unpinned scope %q.\n", scope)
	return nil
}

func runPinList(storeDir string, out io.Writer) error {
	scopes, err := pin.List(storeDir)
	if err != nil {
		return err
	}
	if len(scopes) == 0 {
		return fmt.Errorf("no pinned scopes")
	}
	fmt.Fprintln(out, strings.Join(scopes, "\n"))
	return nil
}
