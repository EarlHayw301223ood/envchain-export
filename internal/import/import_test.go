package imp_test

import (
	"strings"
	"testing"

	"github.com/user/envchain-export/internal/chain"
	imp "github.com/user/envchain-export/internal/import"
)

func newChain(t *testing.T) *chain.Chain {
	t.Helper()
	c, err := chain.New("test")
	if err != nil {
		t.Fatalf("chain.New: %v", err)
	}
	return c
}

func TestRead_DotenvFormat(t *testing.T) {
	input := "FOO=bar\nBAZ=qux\n"
	c := newChain(t)
	if err := imp.Read(c, strings.NewReader(input), imp.FormatDotenv); err != nil {
		t.Fatalf("Read: %v", err)
	}
	if v, _ := c.Get("FOO"); v != "bar" {
		t.Errorf("FOO = %q, want %q", v, "bar")
	}
	if v, _ := c.Get("BAZ"); v != "qux" {
		t.Errorf("BAZ = %q, want %q", v, "qux")
	}
}

func TestRead_PosixFormat(t *testing.T) {
	input := "export FOO='hello world'\nexport BAR='it'\''s fine'\n"
	c := newChain(t)
	if err := imp.Read(c, strings.NewReader(input), imp.FormatPosix); err != nil {
		t.Fatalf("Read: %v", err)
	}
	if v, _ := c.Get("FOO"); v != "hello world" {
		t.Errorf("FOO = %q, want %q", v, "hello world")
	}
}

func TestRead_SkipsBlankAndComments(t *testing.T) {
	input := "# comment\n\nKEY=value\n"
	c := newChain(t)
	if err := imp.Read(c, strings.NewReader(input), imp.FormatDotenv); err != nil {
		t.Fatalf("Read: %v", err)
	}
	if v, _ := c.Get("KEY"); v != "value" {
		t.Errorf("KEY = %q, want %q", v, "value")
	}
}

func TestRead_QuotedValues(t *testing.T) {
	input := `GREETING="hello world"` + "\n"
	c := newChain(t)
	if err := imp.Read(c, strings.NewReader(input), imp.FormatDotenv); err != nil {
		t.Fatalf("Read: %v", err)
	}
	if v, _ := c.Get("GREETING"); v != "hello world" {
		t.Errorf("GREETING = %q, want %q", v, "hello world")
	}
}

func TestRead_UnknownFormat(t *testing.T) {
	c := newChain(t)
	err := imp.Read(c, strings.NewReader("KEY=val\n"), "json")
	if err == nil {
		t.Fatal("expected error for unknown format, got nil")
	}
}

func TestRead_MissingEquals(t *testing.T) {
	c := newChain(t)
	err := imp.Read(c, strings.NewReader("NODEFINITION\n"), imp.FormatDotenv)
	if err == nil {
		t.Fatal("expected error for line without '=', got nil")
	}
}

func TestRead_ExportPrefix(t *testing.T) {
	input := "export MY_VAR=hello\n"
	c := newChain(t)
	if err := imp.Read(c, strings.NewReader(input), imp.FormatDotenv); err != nil {
		t.Fatalf("Read: %v", err)
	}
	if v, _ := c.Get("MY_VAR"); v != "hello" {
		t.Errorf("MY_VAR = %q, want %q", v, "hello")
	}
}
