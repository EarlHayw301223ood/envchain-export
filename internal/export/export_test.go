package export_test

import (
	"strings"
	"testing"

	"github.com/user/envchain-export/internal/chain"
	"github.com/user/envchain-export/internal/export"
)

func newChain(t *testing.T) *chain.Chain {
	t.Helper()
	c, err := chain.New("test")
	if err != nil {
		t.Fatalf("chain.New: %v", err)
	}
	return c
}

func TestWrite_PosixFormat(t *testing.T) {
	c := newChain(t)
	c.Add("FOO", "bar")
	c.Add("BAZ", "qux")

	var sb strings.Builder
	if err := export.Write(&sb, c, export.FormatPosix); err != nil {
		t.Fatalf("Write: %v", err)
	}
	out := sb.String()
	if !strings.Contains(out, "export FOO='bar'") {
		t.Errorf("expected posix export for FOO, got:\n%s", out)
	}
	if !strings.Contains(out, "export BAZ='qux'") {
		t.Errorf("expected posix export for BAZ, got:\n%s", out)
	}
}

func TestWrite_DotenvFormat(t *testing.T) {
	c := newChain(t)
	c.Add("KEY", "value with spaces")

	var sb strings.Builder
	if err := export.Write(&sb, c, export.FormatDotenv); err != nil {
		t.Fatalf("Write: %v", err)
	}
	out := sb.String()
	if !strings.Contains(out, "KEY='value with spaces'") {
		t.Errorf("expected dotenv line, got:\n%s", out)
	}
}

func TestWrite_QuotesSingleQuotes(t *testing.T) {
	c := newChain(t)
	c.Add("MSG", "it's alive")

	var sb strings.Builder
	if err := export.Write(&sb, c, export.FormatPosix); err != nil {
		t.Fatalf("Write: %v", err)
	}
	out := sb.String()
	// single quote inside value must be escaped
	if !strings.Contains(out, `'it'\''s alive'`) {
		t.Errorf("expected escaped single quote, got:\n%s", out)
	}
}

func TestWrite_UnknownFormat(t *testing.T) {
	c := newChain(t)
	c.Add("X", "y")

	var sb strings.Builder
	err := export.Write(&sb, c, export.Format("xml"))
	if err == nil {
		t.Fatal("expected error for unknown format, got nil")
	}
}

func TestWrite_EmptyChain(t *testing.T) {
	c := newChain(t)

	var sb strings.Builder
	if err := export.Write(&sb, c, export.FormatDotenv); err != nil {
		t.Fatalf("Write: %v", err)
	}
	if sb.Len() != 0 {
		t.Errorf("expected empty output for empty chain, got %q", sb.String())
	}
}
