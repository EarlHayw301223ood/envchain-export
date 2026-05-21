package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/nicholasgasior/envchain-export/internal/snapshot"
	"github.com/nicholasgasior/envchain-export/internal/store"
	"github.com/spf13/cobra"
)

func newSnapshotCmd(st *store.Store, out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "snapshot",
		Short: "Manage point-in-time snapshots of a scope",
	}

	cmd.AddCommand(newSnapshotTakeCmd(st, out))
	cmd.AddCommand(newSnapshotRestoreCmd(st, out))
	cmd.AddCommand(newSnapshotListCmd(st, out))

	return cmd
}

func newSnapshotTakeCmd(st *store.Store, out io.Writer) *cobra.Command {
	return &cobra.Command{
		Use:   "take <scope>",
		Short: "Capture a snapshot of a scope",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSnapshotTake(st, out, args[0])
		},
	}
}

func newSnapshotRestoreCmd(st *store.Store, out io.Writer) *cobra.Command {
	return &cobra.Command{
		Use:   "restore <scope> <label>",
		Short: "Restore a scope from a snapshot",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSnapshotRestore(st, out, args[0], args[1])
		},
	}
}

func newSnapshotListCmd(st *store.Store, out io.Writer) *cobra.Command {
	return &cobra.Command{
		Use:   "list <scope>",
		Short: "List snapshots for a scope",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSnapshotList(st, out, args[0])
		},
	}
}

func runSnapshotTake(st *store.Store, out io.Writer, scope string) error {
	passphrase, err := promptPassphrase(scope)
	if err != nil {
		return err
	}
	label, err := snapshot.Take(st, scope, passphrase)
	if err != nil {
		return fmt.Errorf("snapshot take: %w", err)
	}
	fmt.Fprintf(out, "snapshot taken: %s\n", label)
	return nil
}

func runSnapshotRestore(st *store.Store, out io.Writer, scope, label string) error {
	passphrase, err := promptPassphrase(scope)
	if err != nil {
		return err
	}
	if err := snapshot.Restore(st, scope, label, passphrase); err != nil {
		return fmt.Errorf("snapshot restore: %w", err)
	}
	fmt.Fprintf(out, "scope %q restored from snapshot %s\n", scope, label)
	return nil
}

func runSnapshotList(st *store.Store, out io.Writer, scope string) error {
	labels, err := snapshot.List(st, scope)
	if err != nil {
		if errors.Is(err, snapshot.ErrNoSnapshots) {
			fmt.Fprintf(os.Stderr, "no snapshots for scope %q\n", scope)
			return err
		}
		return fmt.Errorf("snapshot list: %w", err)
	}
	fmt.Fprintln(out, strings.Join(labels, "\n"))
	return nil
}
