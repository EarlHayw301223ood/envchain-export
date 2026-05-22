package main

import (
	"fmt"
	"io"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/your-org/envchain-export/internal/search"
)

func newSearchCmd(storeDir *string) *cobra.Command {
	var (
		searchValues    bool
		caseInsensitive bool
		scopeFilter     string
	)

	cmd := &cobra.Command{
		Use:   "search <query>",
		Short: "Search for keys (and optionally values) across all scopes",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			pass, err := promptPassphrase(cmd)
			if err != nil {
				return err
			}
			return runSearch(os.Stdout, *storeDir, pass, args[0], scopeFilter, search.Options{
				SearchValues:    searchValues,
				CaseInsensitive: caseInsensitive,
			})
		},
	}

	cmd.Flags().BoolVarP(&searchValues, "values", "v", false, "also search inside values")
	cmd.Flags().BoolVarP(&caseInsensitive, "ignore-case", "i", false, "case-insensitive matching")
	cmd.Flags().StringVarP(&scopeFilter, "scope", "s", "", "restrict search to a single scope")
	return cmd
}

func runSearch(w io.Writer, storeDir, passphrase, query, scopeFilter string, opts search.Options) error {
	matches, err := search.Search(storeDir, passphrase, query, opts)
	if err != nil {
		return fmt.Errorf("search: %w", err)
	}

	if scopeFilter != "" {
		matches = search.FilterByScope(matches, scopeFilter)
	}

	if len(matches) == 0 {
		return fmt.Errorf("search: no matches found for %q", query)
	}

	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "SCOPE\tKEY\tVALUE")
	for _, m := range matches {
		display := m.Value
		if len(display) > 40 {
			display = display[:37] + "..."
		}
		fmt.Fprintf(tw, "%s\t%s\t%s\n", m.Scope, m.Key, display)
	}
	return tw.Flush()
}
