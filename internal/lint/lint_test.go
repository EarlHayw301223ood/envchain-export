package lint_test

import (
	"testing"

	"github.com/user/envchain-export/internal/chain"
	"github.com/user/envchain-export/internal/lint"
)

func newChain(t *testing.T, pairs map[string]string) *chain.Chain {
	t.Helper()
	c, err := chain.New("test")
	if err != nil {
		t.Fatalf("chain.New: %v", err)
	}
	for k, v := range pairs {
		if err := c.Add(k, v); err != nil {
			t.Fatalf("c.Add(%q): %v", k, err)
		}
	}
	return c
}

func TestLint_Clean(t *testing.T) {
	c := newChain(t, map[string]string{
		"DATABASE_URL": "postgres://localhost/mydb",
		"APP_ENV":      "production",
	})
	findings := lint.Lint(c)
	if len(findings) != 0 {
		t.Fatalf("expected no findings, got %v", findings)
	}
}

func TestLint_ShortSecret(t *testing.T) {
	c := newChain(t, map[string]string{
		"API_TOKEN": "abc",
	})
	findings := lint.Lint(c)
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding, got %d: %v", len(findings), findings)
	}
	if findings[0].Severity != lint.Warn {
		t.Errorf("expected WARN, got %s", findings[0].Severity)
	}
}

func TestLint_NewlineInValue(t *testing.T) {
	c := newChain(t, map[string]string{
		"MULTI_LINE": "line1\nline2",
	})
	findings := lint.Lint(c)
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding, got %d: %v", len(findings), findings)
	}
	if findings[0].Severity != lint.Error {
		t.Errorf("expected ERROR, got %s", findings[0].Severity)
	}
}

func TestLint_DollarSign(t *testing.T) {
	c := newChain(t, map[string]string{
		"GREETING": "hello $USER",
	})
	findings := lint.Lint(c)
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding, got %d: %v", len(findings), findings)
	}
	if findings[0].Severity != lint.Warn {
		t.Errorf("expected WARN, got %s", findings[0].Severity)
	}
}

func TestLint_MultipleIssues(t *testing.T) {
	c := newChain(t, map[string]string{
		"SECRET_KEY": "x",          // short secret → WARN
		"BAD_VAR":    "val\nnext",  // newline → ERROR
	})
	findings := lint.Lint(c)
	if len(findings) != 2 {
		t.Fatalf("expected 2 findings, got %d: %v", len(findings), findings)
	}
}

func TestFinding_String(t *testing.T) {
	f := lint.Finding{Key: "FOO", Severity: lint.Warn, Message: "test message"}
	got := f.String()
	expected := "[WARN] FOO: test message"
	if got != expected {
		t.Errorf("String() = %q, want %q", got, expected)
	}
}
