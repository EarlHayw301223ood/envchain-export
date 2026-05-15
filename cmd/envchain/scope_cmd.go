package main

import (
	"fmt"
	"os"

	"github.com/user/envchain-export/internal/scope"
)

// runScopeList prints all scopes found in the configured store directory.
func runScopeList(storeDir string) error {
	scopes, err := scope.List(storeDir)
	if err != nil {
		return fmt.Errorf("list: %w", err)
	}
	for _, s := range scopes {
		fmt.Println(s)
	}
	return nil
}

// runScopeDelete removes a named scope from the store directory after
// prompting for confirmation unless --force is set.
func runScopeDelete(storeDir, name string, force bool) error {
	if !force {
		fmt.Fprintf(os.Stderr, "Delete scope %q? [y/N] ", name)
		var answer string
		if _, err := fmt.Scanln(&answer); err != nil || (answer != "y" && answer != "Y") {
			fmt.Fprintln(os.Stderr, "Aborted.")
			return nil
		}
	}
	if err := scope.Delete(storeDir, name); err != nil {
		return fmt.Errorf("delete: %w", err)
	}
	fmt.Fprintf(os.Stderr, "Scope %q deleted.\n", name)
	return nil
}

// runScopeExists exits with code 0 if the scope exists, 1 otherwise.
func runScopeExists(storeDir, name string) error {
	ok, err := scope.Exists(storeDir, name)
	if err != nil {
		return fmt.Errorf("exists: %w", err)
	}
	if !ok {
		os.Exit(1)
	}
	return nil
}
