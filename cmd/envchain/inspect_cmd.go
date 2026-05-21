package main

import (
	"fmt"
	"io"
	"os"

	"github.com/urfave/cli/v2"

	"github.com/yourorg/envchain-export/internal/inspect"
	"github.com/yourorg/envchain-export/internal/store"
)

// newInspectCmd returns the CLI command for inspecting a scope's contents.
// By default, values are masked; use --reveal to show plaintext values.
func newInspectCmd(storeDir string) *cli.Command {
	return &cli.Command{
		Name:      "inspect",
		Usage:     "Display keys (and optionally values) stored in a scope",
		ArgsUsage: "<scope>",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "reveal",
				Aliases: []string{"r"},
				Usage:   "Show plaintext values instead of masked output",
			},
		},
		Action: func(c *cli.Context) error {
			if c.NArg() < 1 {
				return fmt.Errorf("scope name is required")
			}
			scope := c.Args().First()
			mask := !c.Bool("reveal")

			passphrase, err := promptPassphrase(scope)
			if err != nil {
				return err
			}

			s, err := store.New(storeDir)
			if err != nil {
				return fmt.Errorf("opening store: %w", err)
			}

			return runInspect(c.App.Writer, s, scope, passphrase, mask)
		},
	}
}

// runInspect loads the given scope and writes an inspection table to w.
func runInspect(w io.Writer, s *store.Store, scope, passphrase string, mask bool) error {
	rows, err := inspect.Inspect(s, scope, passphrase, mask)
	if err != nil {
		return fmt.Errorf("inspect %q: %w", scope, err)
	}

	if len(rows) == 0 {
		fmt.Fprintf(w, "scope %q is empty\n", scope)
		return nil
	}

	// Determine column widths for aligned output.
	maxKey := len("KEY")
	for _, row := range rows {
		if len(row.Key) > maxKey {
			maxKey = len(row.Key)
		}
	}

	fmt.Fprintf(w, "Scope: %s\n", scope)
	fmt.Fprintf(w, "%-*s  %s\n", maxKey, "KEY", "VALUE")
	fmt.Fprintf(w, "%s  %s\n", repeatDash(maxKey), repeatDash(20))
	for _, row := range rows {
		fmt.Fprintf(w, "%-*s  %s\n", maxKey, row.Key, row.Value)
	}
	return nil
}

// repeatDash returns a string of n dash characters.
func repeatDash(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = '-'
	}
	return string(b)
}

// init registers the inspect command; called from main via run.
var _ = os.Stderr // ensure os import used via promptPassphrase dependency
