package main

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/nicholasgasior/envchain-export/internal/lock"
	"github.com/spf13/cobra"
)

func newLockCmd(storeDir string, out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "lock",
		Short: "Manage advisory locks on scopes",
	}
	cmd.AddCommand(newLockAcquireCmd(storeDir, out))
	cmd.AddCommand(newLockReleaseCmd(storeDir, out))
	cmd.AddCommand(newLockStatusCmd(storeDir, out))
	return cmd
}

func newLockAcquireCmd(storeDir string, out io.Writer) *cobra.Command {
	return &cobra.Command{
		Use:   "acquire <scope>",
		Short: "Acquire an advisory lock on a scope",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runLockAcquire(storeDir, args[0], out)
		},
	}
}

func newLockReleaseCmd(storeDir string, out io.Writer) *cobra.Command {
	return &cobra.Command{
		Use:   "release <scope>",
		Short: "Release an advisory lock on a scope",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runLockRelease(storeDir, args[0], out)
		},
	}
}

func newLockStatusCmd(storeDir string, out io.Writer) *cobra.Command {
	return &cobra.Command{
		Use:   "status <scope>",
		Short: "Report whether a scope is currently locked",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runLockStatus(storeDir, args[0], out)
		},
	}
}

func runLockAcquire(storeDir, scope string, out io.Writer) error {
	if err := lock.Acquire(storeDir, scope); err != nil {
		if errors.Is(err, lock.ErrLocked) {
			return fmt.Errorf("scope %q is already locked", scope)
		}
		return err
	}
	fmt.Fprintf(out, "Lock acquired on scope %q (pid %d)\n", scope, os.Getpid())
	return nil
}

func runLockRelease(storeDir, scope string, out io.Writer) error {
	if err := lock.Release(storeDir, scope); err != nil {
		if errors.Is(err, lock.ErrNotLocked) {
			return fmt.Errorf("scope %q is not locked", scope)
		}
		return err
	}
	fmt.Fprintf(out, "Lock released on scope %q\n", scope)
	return nil
}

func runLockStatus(storeDir, scope string, out io.Writer) error {
	if lock.IsLocked(storeDir, scope) {
		fmt.Fprintf(out, "scope %q is LOCKED\n", scope)
	} else {
		fmt.Fprintf(out, "scope %q is unlocked\n", scope)
	}
	return nil
}
