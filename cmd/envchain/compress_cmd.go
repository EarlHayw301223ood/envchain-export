package main

import (
	"fmt"
	"os"

	"github.com/nicholasgasior/envchain-export/internal/compress"
	"github.com/spf13/cobra"
)

func newCompressCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "compress",
		Short: "Compress or inspect compression of a raw payload file",
	}
	cmd.AddCommand(newCompressCheckCmd())
	return cmd
}

func newCompressCheckCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "check <file>",
		Short: "Report whether a file is gzip-compressed",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runCompressCheck(cmd, args)
		},
	}
}

func runCompressCheck(cmd *cobra.Command, args []string) error {
	path := args[0]

	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("compress check: read file: %w", err)
	}

	if compress.IsCompressed(data) {
		fmt.Fprintf(cmd.OutOrStdout(), "%s: compressed (gzip)\n", path)
	} else {
		fmt.Fprintf(cmd.OutOrStdout(), "%s: not compressed\n", path)
	}
	return nil
}
