package main

import (
	"fmt"
	"io"
	"os"
	"text/tabwriter"

	"github.com/nicholasgasior/envchain-export/internal/history"
	"github.com/spf13/cobra"
)

func newHistoryCmd(storeDir string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "history",
		Short: "Manage scope change history",
	}
	cmd.AddCommand(newHistoryListCmd(storeDir))
	cmd.AddCommand(newHistoryClearCmd(storeDir))
	return cmd
}

func newHistoryListCmd(storeDir string) *cobra.Command {
	return &cobra.Command{
		Use:   "list <scope>",
		Short: "List recorded history entries for a scope",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runHistoryList(storeDir, args[0], cmd.OutOrStdout())
		},
	}
}

func newHistoryClearCmd(storeDir string) *cobra.Command {
	return &cobra.Command{
		Use:   "clear <scope>",
		Short: "Clear all history entries for a scope",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runHistoryClear(storeDir, args[0], cmd.OutOrStdout())
		},
	}
}

func runHistoryList(storeDir, scope string, w io.Writer) error {
	entries, err := history.Read(storeDir, scope)
	if err != nil {
		return fmt.Errorf("read history: %w", err)
	}
	if len(entries) == 0 {
		fmt.Fprintf(w, "No history for scope %q.\n", scope)
		return nil
	}
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "TIMESTAMP\tACTION\tDETAIL")
	for _, e := range entries {
		fmt.Fprintf(tw, "%s\t%s\t%s\n",
			e.Timestamp.Format("2006-01-02 15:04:05 UTC"),
			e.Action,
			e.Detail,
		)
	}
	return tw.Flush()
}

func runHistoryClear(storeDir, scope string, w io.Writer) error {
	if err := history.Clear(storeDir, scope); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("no history found for scope %q", scope)
		}
		return fmt.Errorf("clear history: %w", err)
	}
	fmt.Fprintf(w, "History cleared for scope %q.\n", scope)
	return nil
}
