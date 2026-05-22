package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/user/envchain-export/internal/ttl"
)

func newTTLCmd(storeDir string, out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ttl",
		Short: "Manage expiry (time-to-live) for scopes",
	}
	cmd.AddCommand(newTTLSetCmd(storeDir, out))
	cmd.AddCommand(newTTLGetCmd(storeDir, out))
	cmd.AddCommand(newTTLClearCmd(storeDir, out))
	return cmd
}

func newTTLSetCmd(storeDir string, out io.Writer) *cobra.Command {
	return &cobra.Command{
		Use:   "set <scope> <duration>",
		Short: "Set expiry for a scope (e.g. 24h, 30m)",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runTTLSet(storeDir, args[0], args[1], out)
		},
	}
}

func newTTLGetCmd(storeDir string, out io.Writer) *cobra.Command {
	return &cobra.Command{
		Use:   "get <scope>",
		Short: "Show expiry time for a scope",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runTTLGet(storeDir, args[0], out)
		},
	}
}

func newTTLClearCmd(storeDir string, out io.Writer) *cobra.Command {
	return &cobra.Command{
		Use:   "clear <scope>",
		Short: "Remove expiry for a scope",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runTTLClear(storeDir, args[0], out)
		},
	}
}

func runTTLSet(storeDir, scope, rawDuration string, out io.Writer) error {
	d, err := time.ParseDuration(rawDuration)
	if err != nil {
		return fmt.Errorf("invalid duration %q: %w", rawDuration, err)
	}
	if err := ttl.Set(storeDir, scope, d); err != nil {
		return err
	}
	fmt.Fprintf(out, "expiry set for scope %q: %s from now\n", scope, d)
	return nil
}

func runTTLGet(storeDir, scope string, out io.Writer) error {
	t, err := ttl.Get(storeDir, scope)
	if errors.Is(err, ttl.ErrNoExpiry) {
		fmt.Fprintf(out, "scope %q has no expiry set\n", scope)
		return nil
	}
	if err != nil {
		return err
	}
	remaining := time.Until(t).Round(time.Second)
	status := "valid"
	if remaining <= 0 {
		status = "EXPIRED"
		remaining = 0
	}
	fmt.Fprintf(out, "scope %q expires at %s (%s remaining) [%s]\n",
		scope, t.Format(time.RFC3339), remaining, status)
	return nil
}

func runTTLClear(storeDir, scope string, out io.Writer) error {
	if err := ttl.Clear(storeDir, scope); err != nil {
		return err
	}
	fmt.Fprintf(out, "expiry cleared for scope %q\n", scope)
	return nil
}

var _ = os.Stderr // satisfy import
