// Package env provides functionality to inject scoped environment variables
// into a subprocess, launching it with the decrypted key-value pairs merged
// into the current process environment.
package env

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/user/envchain-export/internal/chain"
)

// ErrEmptyCommand is returned when no command is provided to Run.
var ErrEmptyCommand = fmt.Errorf("env: command must not be empty")

// Run executes the given command with the variables from ch merged into
// the current process environment. The child process inherits all existing
// environment variables, with the chain's variables taking precedence.
func Run(ch *chain.Chain, command []string) error {
	if len(command) == 0 {
		return ErrEmptyCommand
	}

	env := buildEnv(ch)

	cmd := exec.Command(command[0], command[1:]...)
	cmd.Env = env
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("env: command failed: %w", err)
	}
	return nil
}

// buildEnv merges the current OS environment with the variables in ch.
// Variables from ch override any existing OS environment variables with
// the same key.
func buildEnv(ch *chain.Chain) []string {
	overrides := make(map[string]string)
	for _, k := range ch.Keys() {
		v, _ := ch.Get(k)
		overrides[k] = v
	}

	base := os.Environ()
	result := make([]string, 0, len(base)+len(overrides))

	seen := make(map[string]bool)
	for _, entry := range base {
		for k, v := range overrides {
			if len(entry) > len(k)+1 && entry[:len(k)+1] == k+"=" {
				result = append(result, k+"="+v)
				seen[k] = true
				goto next
			}
		}
		result = append(result, entry)
	next:
	}

	for k, v := range overrides {
		if !seen[k] {
			result = append(result, k+"="+v)
		}
	}

	return result
}
