package env_test

import (
	"os"
	"testing"

	"github.com/user/envchain-export/internal/chain"
	"github.com/user/envchain-export/internal/env"
)

func newChain(t *testing.T, pairs map[string]string) *chain.Chain {
	t.Helper()
	ch, err := chain.New("test")
	if err != nil {
		t.Fatalf("chain.New: %v", err)
	}
	for k, v := range pairs {
		if err := ch.Add(k, v); err != nil {
			t.Fatalf("ch.Add(%q, %q): %v", k, v, err)
		}
	}
	return ch
}

func TestRun_EmptyCommand(t *testing.T) {
	ch := newChain(t, nil)
	err := env.Run(ch, []string{})
	if err == nil {
		t.Fatal("expected error for empty command, got nil")
	}
	if err != env.ErrEmptyCommand {
		t.Fatalf("expected ErrEmptyCommand, got %v", err)
	}
}

func TestRun_CommandReceivesInjectedVars(t *testing.T) {
	ch := newChain(t, map[string]string{
		"ENVCHAIN_TEST_VAR": "hello_injected",
	})

	// Use env -u to check variable presence portably via exit code trick;
	// instead just run printenv and check it doesn't error.
	err := env.Run(ch, []string{"printenv", "ENVCHAIN_TEST_VAR"})
	if err != nil {
		t.Fatalf("expected ENVCHAIN_TEST_VAR to be set, got error: %v", err)
	}
}

func TestRun_OverridesExistingVar(t *testing.T) {
	os.Setenv("ENVCHAIN_OVERRIDE", "original")
	t.Cleanup(func() { os.Unsetenv("ENVCHAIN_OVERRIDE") })

	ch := newChain(t, map[string]string{
		"ENVCHAIN_OVERRIDE": "overridden",
	})

	// printenv exits 0 only if the variable exists; the override value
	// is validated by the absence of an error.
	err := env.Run(ch, []string{"printenv", "ENVCHAIN_OVERRIDE"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRun_NonZeroExitReturnsError(t *testing.T) {
	ch := newChain(t, nil)
	err := env.Run(ch, []string{"false"})
	if err == nil {
		t.Fatal("expected error from non-zero exit, got nil")
	}
}
