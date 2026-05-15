package passphrase_test

import (
	"testing"

	"github.com/yourusername/envchain-export/internal/passphrase"
)

// promptFn is the signature shared by Prompt and the first call inside
// PromptConfirm; we test the exported error sentinels directly since
// terminal interaction cannot be exercised in unit tests without a pty.

func TestErrMismatch_IsDistinct(t *testing.T) {
	if passphrase.ErrMismatch == passphrase.ErrEmpty {
		t.Fatal("ErrMismatch and ErrEmpty must be distinct errors")
	}
}

func TestErrEmpty_Message(t *testing.T) {
	if passphrase.ErrEmpty.Error() == "" {
		t.Fatal("ErrEmpty must have a non-empty message")
	}
}

func TestErrMismatch_Message(t *testing.T) {
	if passphrase.ErrMismatch.Error() == "" {
		t.Fatal("ErrMismatch must have a non-empty message")
	}
}

// mockPrompt exercises the PromptConfirm mismatch path by replacing
// the reader; since Prompt wraps term.ReadPassword we test via the
// exported Confirm helper that accepts an injected reader.
func TestConfirm_Match(t *testing.T) {
	err := passphrase.Confirm("secret", "secret")
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestConfirm_Mismatch(t *testing.T) {
	err := passphrase.Confirm("secret", "other")
	if err != passphrase.ErrMismatch {
		t.Fatalf("expected ErrMismatch, got %v", err)
	}
}

func TestConfirm_Empty(t *testing.T) {
	err := passphrase.Confirm("", "")
	if err != passphrase.ErrEmpty {
		t.Fatalf("expected ErrEmpty, got %v", err)
	}
}
