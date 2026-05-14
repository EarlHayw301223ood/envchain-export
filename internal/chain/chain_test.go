package chain

import (
	"testing"
)

func TestNew_ValidScope(t *testing.T) {
	c, err := New("my-app")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if c.Scope != "my-app" {
		t.Errorf("expected scope 'my-app', got %q", c.Scope)
	}
	if c.Len() != 0 {
		t.Errorf("expected empty chain, got len %d", c.Len())
	}
}

func TestNew_InvalidScope(t *testing.T) {
	invalid := []string{"", " spaces", "!bang", "-leading-dash"}
	for _, s := range invalid {
		_, err := New(s)
		if err == nil {
			t.Errorf("expected error for scope %q, got nil", s)
		}
	}
}

func TestAdd_AndGet(t *testing.T) {
	c, _ := New("test")
	c.Add("FOO", "bar")

	val, ok := c.Get("FOO")
	if !ok {
		t.Fatal("expected key FOO to exist")
	}
	if val != "bar" {
		t.Errorf("expected 'bar', got %q", val)
	}
}

func TestAdd_UpdatesExisting(t *testing.T) {
	c, _ := New("test")
	c.Add("FOO", "bar")
	c.Add("FOO", "baz")

	if c.Len() != 1 {
		t.Errorf("expected 1 var after update, got %d", c.Len())
	}
	val, _ := c.Get("FOO")
	if val != "baz" {
		t.Errorf("expected updated value 'baz', got %q", val)
	}
}

func TestRemove_Existing(t *testing.T) {
	c, _ := New("test")
	c.Add("FOO", "bar")
	c.Add("BAZ", "qux")

	removed := c.Remove("FOO")
	if !removed {
		t.Error("expected Remove to return true")
	}
	if c.Len() != 1 {
		t.Errorf("expected 1 var remaining, got %d", c.Len())
	}
	_, ok := c.Get("FOO")
	if ok {
		t.Error("expected FOO to be removed")
	}
}

func TestRemove_NonExistent(t *testing.T) {
	c, _ := New("test")
	removed := c.Remove("MISSING")
	if removed {
		t.Error("expected Remove to return false for missing key")
	}
}

func TestGet_Missing(t *testing.T) {
	c, _ := New("test")
	_, ok := c.Get("NOPE")
	if ok {
		t.Error("expected Get to return false for missing key")
	}
}
