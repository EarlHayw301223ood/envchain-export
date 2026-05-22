package main

import (
	"fmt"
	"io"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/user/envchain-export/internal/audit"
)

func newAuditCmd(storeDir string) *cobra.Command {
	var scopeFilter string

	cmd := &cobra.Command{
		Use:   "audit",
		Short: "Show the audit log of scope operations",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runAudit(os.Stdout, storeDir, scopeFilter)
		},
	}

	cmd.Flags().StringVarP(&scopeFilter, "scope", "s", "", "Filter entries by scope name")
	return cmd
}

func runAudit(w io.Writer, storeDir, scopeFilter string) error {
	entries, err := audit.Read(storeDir)
	if err != nil {
		return fmt.Errorf("audit: %w", err)
	}

	if len(entries) == 0 {
		fmt.Fprintln(w, "No audit entries found.")
		return nil
	}

	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "TIMESTAMP\tSCOPE\tACTION\tDETAIL")

	count := 0
	for _, e := range entries {
		if scopeFilter != "" && e.Scope != scopeFilter {
			continue
		}
		detail := e.Detail
		if detail == "" {
			detail = "-"
		}
		fmt.Fprintf(tw, "%s\t%s\t%s\t%s\n",
			e.Timestamp.Format("2006-01-02T15:04:05Z"),
			e.Scope,
			e.Action,
			detail,
		)
		count++
	}

	if count == 0 {
		fmt.Fprintf(w, "No audit entries found for scope %q.\n", scopeFilter)
		return nil
	}

	return tw.Flush()
}
