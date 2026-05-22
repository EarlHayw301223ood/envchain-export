package main

import (
	"fmt"
	"os"

	"github.com/nicholasgasior/envchain-export/internal/store"
	"github.com/nicholasgasior/envchain-export/internal/template"
	"github.com/spf13/cobra"
)

func newTemplateCmd(storeDir string) *cobra.Command {
	var templateFile string

	cmd := &cobra.Command{
		Use:   "template <scope> -f <template-file>",
		Short: "Render a text template using variables from a scope",
		Long: `Reads a Go text/template from the given file and renders it by
substituting values from the named scope. Access variables with
{{ index . "VAR_NAME" }}.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runTemplate(storeDir, args[0], templateFile)
		},
	}

	cmd.Flags().StringVarP(&templateFile, "file", "f", "", "path to the template file (required)")
	_ = cmd.MarkFlagRequired("file")
	return cmd
}

func runTemplate(storeDir, scope, templateFile string) error {
	pass, err := promptPassphrase(scope)
	if err != nil {
		return err
	}

	s := store.New(storeDir)
	c, err := s.Load(scope, pass)
	if err != nil {
		return fmt.Errorf("load scope %q: %w", scope, err)
	}

	src, err := os.ReadFile(templateFile)
	if err != nil {
		return fmt.Errorf("read template file: %w", err)
	}

	if err := template.Render(os.Stdout, c, string(src)); err != nil {
		return err
	}
	return nil
}
