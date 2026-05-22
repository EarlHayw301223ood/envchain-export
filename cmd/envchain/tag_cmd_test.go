package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/envchain/envchain-export/internal/tag"
)

func setupTagStore(t *testing.T) string {
	t.Helper()
	dir := filepath.Join(t.TempDir(), "store")
	if err := os.MkdirAll(dir, 0700); err != nil {
		t.Fatal(err)
	}
	return dir
}

func TestRunTagAdd_OutputsMessage(t *testing.T) {
	dir := setupTagStore(t)
	var buf bytes.Buffer
	if err := runTagAdd(dir, &buf, "myapp", "production"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "production") {
		t.Errorf("expected output to mention tag, got: %s", buf.String())
	}
}

func TestRunTagList_PrintsTags(t *testing.T) {
	dir := setupTagStore(t)
	tag.Add(dir, "myapp", "production")
	tag.Add(dir, "myapp", "aws")

	var buf bytes.Buffer
	if err := runTagList(dir, &buf, "myapp"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "production") || !strings.Contains(out, "aws") {
		t.Errorf("expected both tags in output, got: %s", out)
	}
}

func TestRunTagList_NoTags_ReturnsError(t *testing.T) {
	dir := setupTagStore(t)
	var buf bytes.Buffer
	err := runTagList(dir, &buf, "myapp")
	if err == nil {
		t.Fatal("expected error for scope with no tags")
	}
	if !strings.Contains(err.Error(), "no tags") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestRunTagRemove_OutputsMessage(t *testing.T) {
	dir := setupTagStore(t)
	tag.Add(dir, "myapp", "staging")

	var buf bytes.Buffer
	if err := runTagRemove(dir, &buf, "myapp", "staging"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "staging") {
		t.Errorf("expected output to mention tag, got: %s", buf.String())
	}
}

func TestRunTagScopes_PrintsMatchingScopes(t *testing.T) {
	dir := setupTagStore(t)
	tag.Add(dir, "app1", "production")
	tag.Add(dir, "app2", "production")
	tag.Add(dir, "app3", "staging")

	var buf bytes.Buffer
	if err := runTagScopes(dir, &buf, "production"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "app1") || !strings.Contains(out, "app2") {
		t.Errorf("expected app1 and app2 in output, got: %s", out)
	}
	if strings.Contains(out, "app3") {
		t.Errorf("app3 should not appear for tag 'production', got: %s", out)
	}
}
