package template_test

import (
	"strings"
	"testing"

	"github.com/nicholasgasior/envchain-export/internal/chain"
	tmpl "github.com/nicholasgasior/envchain-export/internal/template"
)

func newChain(t *testing.T, pairs map[string]string) *chain.Chain {
	t.Helper()
	c, err := chain.New("test")
	if err != nil {
		t.Fatalf("chain.New: %v", err)
	}
	for k, v := range pairs {
		if err := c.Add(k, v); err != nil {
			t.Fatalf("chain.Add(%q): %v", k, err)
		}
	}
	return c
}

func TestRender_BasicSubstitution(t *testing.T) {
	c := newChain(t, map[string]string{"APP_HOST": "localhost", "APP_PORT": "8080"})
	src := `host={{ index . "APP_HOST" }} port={{ index . "APP_PORT" }}`
	var sb strings.Builder
	if err := tmpl.Render(&sb, c, src); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := "host=localhost port=8080"
	if sb.String() != want {
		t.Errorf("got %q, want %q", sb.String(), want)
	}
}

func TestRender_MissingKey_ReturnsError(t *testing.T) {
	c := newChain(t, map[string]string{"PRESENT": "yes"})
	src := `{{ index . "MISSING" }}`
	var sb strings.Builder
	err := tmpl.Render(&sb, c, src)
	if err == nil {
		t.Fatal("expected error for missing key, got nil")
	}
	if !strings.Contains(err.Error(), "render error") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestRender_InvalidTemplate_ReturnsError(t *testing.T) {
	c := newChain(t, map[string]string{})
	src := `{{ .Unclosed`
	var sb strings.Builder
	if err := tmpl.Render(&sb, c, src); err == nil {
		t.Fatal("expected parse error, got nil")
	}
}

func TestRender_MultilineTemplate(t *testing.T) {
	c := newChain(t, map[string]string{"DB_USER": "admin", "DB_PASS": "secret"})
	src := "user={{ index . \"DB_USER\" }}\npass={{ index . \"DB_PASS\" }}"
	var sb strings.Builder
	if err := tmpl.Render(&sb, c, src); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(sb.String(), "user=admin") || !strings.Contains(sb.String(), "pass=secret") {
		t.Errorf("unexpected output: %q", sb.String())
	}
}

func TestRender_EmptyTemplate(t *testing.T) {
	c := newChain(t, map[string]string{"FOO": "bar"})
	var sb strings.Builder
	if err := tmpl.Render(&sb, c, ""); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sb.String() != "" {
		t.Errorf("expected empty output, got %q", sb.String())
	}
}
