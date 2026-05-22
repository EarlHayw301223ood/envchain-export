package main

import (
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/envchain/envchain-export/internal/tag"
	"github.com/spf13/cobra"
)

func newTagCmd(storeDir string, out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tag",
		Short: "Manage tags on scopes",
	}
	cmd.AddCommand(newTagAddCmd(storeDir, out))
	cmd.AddCommand(newTagRemoveCmd(storeDir, out))
	cmd.AddCommand(newTagListCmd(storeDir, out))
	cmd.AddCommand(newTagScopesCmd(storeDir, out))
	return cmd
}

func newTagAddCmd(storeDir string, out io.Writer) *cobra.Command {
	return &cobra.Command{
		Use:   "add <scope> <tag>",
		Short: "Add a tag to a scope",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runTagAdd(storeDir, out, args[0], args[1])
		},
	}
}

func newTagRemoveCmd(storeDir string, out io.Writer) *cobra.Command {
	return &cobra.Command{
		Use:   "remove <scope> <tag>",
		Short: "Remove a tag from a scope",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runTagRemove(storeDir, out, args[0], args[1])
		},
	}
}

func newTagListCmd(storeDir string, out io.Writer) *cobra.Command {
	return &cobra.Command{
		Use:   "list <scope>",
		Short: "List tags on a scope",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runTagList(storeDir, out, args[0])
		},
	}
}

func newTagScopesCmd(storeDir string, out io.Writer) *cobra.Command {
	return &cobra.Command{
		Use:   "scopes <tag>",
		Short: "List scopes that carry a given tag",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runTagScopes(storeDir, out, args[0])
		},
	}
}

func runTagAdd(storeDir string, out io.Writer, scope, t string) error {
	if err := tag.Add(storeDir, scope, t); err != nil {
		return err
	}
	fmt.Fprintf(out, "Tagged scope %q with %q.\n", scope, t)
	return nil
}

func runTagRemove(storeDir string, out io.Writer, scope, t string) error {
	if err := tag.Remove(storeDir, scope, t); err != nil {
		if errors.Is(err, tag.ErrTagNotFound) {
			return fmt.Errorf("scope %q does not have tag %q", scope, t)
		}
		return err
	}
	fmt.Fprintf(out, "Removed tag %q from scope %q.\n", t, scope)
	return nil
}

func runTagList(storeDir string, out io.Writer, scope string) error {
	tags, err := tag.List(storeDir, scope)
	if err != nil {
		if errors.Is(err, tag.ErrNoTags) {
			return fmt.Errorf("scope %q has no tags", scope)
		}
		return err
	}
	fmt.Fprintln(out, strings.Join(tags, "\n"))
	return nil
}

func runTagScopes(storeDir string, out io.Writer, t string) error {
	scopes, err := tag.ScopesByTag(storeDir, t)
	if err != nil {
		return err
	}
	if len(scopes) == 0 {
		return fmt.Errorf("no scopes found with tag %q", t)
	}
	fmt.Fprintln(out, strings.Join(scopes, "\n"))
	return nil
}
