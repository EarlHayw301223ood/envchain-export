package validate_test

import (
	"testing"

	"github.com/yourorg/envchain-export/internal/validate"
)

func TestKey_Valid(t *testing.T) {
	valid := []string{
		"FOO",
		"_BAR",
		"foo_bar",
		"FOO123",
		"_",
		"a",
	}
	for _, k := range valid {
		if err := validate.Key(k); err != nil {
			t.Errorf("Key(%q) unexpected error: %v", k, err)
		}
	}
}

func TestKey_Invalid(t *testing.T) {
	cases := []struct {
		input string
		wantErr error
	}{
		{"", validate.ErrEmptyKey},
		{"1FOO", validate.ErrInvalidKey},
		{"FOO-BAR", validate.ErrInvalidKey},
		{"FOO BAR", validate.ErrInvalidKey},
		{"FOO=BAR", validate.ErrInvalidKey},
		{"$FOO", validate.ErrInvalidKey},
	}
	for _, tc := range cases {
		err := validate.Key(tc.input)
		if err != tc.wantErr {
			t.Errorf("Key(%q) = %v, want %v", tc.input, err, tc.wantErr)
		}
	}
}

func TestValue_Valid(t *testing.T) {
	valid := []string{
		"",
		"hello world",
		"it's a value",
		"line1\nline2",
		"unicode: 日本語",
	}
	for _, v := range valid {
		if err := validate.Value(v); err != nil {
			t.Errorf("Value(%q) unexpected error: %v", v, err)
		}
	}
}

func TestValue_NullByte(t *testing.T) {
	err := validate.Value("bad\x00value")
	if err != validate.ErrNullByteInValue {
		t.Errorf("Value with null byte = %v, want ErrNullByteInValue", err)
	}
}

func TestPair_Valid(t *testing.T) {
	if err := validate.Pair("MY_VAR", "some value"); err != nil {
		t.Errorf("Pair() unexpected error: %v", err)
	}
}

func TestPair_InvalidKey(t *testing.T) {
	err := validate.Pair("123BAD", "value")
	if err != validate.ErrInvalidKey {
		t.Errorf("Pair() = %v, want ErrInvalidKey", err)
	}
}

func TestPair_InvalidValue(t *testing.T) {
	err := validate.Pair("GOOD_KEY", "bad\x00val")
	if err != validate.ErrNullByteInValue {
		t.Errorf("Pair() = %v, want ErrNullByteInValue", err)
	}
}
