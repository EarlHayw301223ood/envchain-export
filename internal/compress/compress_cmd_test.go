package compress_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/your-org/envchain-export/internal/compress"
)

func TestDecompress_InvalidData_ReturnsError(t *testing.T) {
	_, err := compress.Decompress([]byte("not compressed data"))
	if err == nil {
		t.Fatal("expected error for invalid compressed data, got nil")
	}
}

func TestDecompress_EmptyInput_ReturnsError(t *testing.T) {
	_, err := compress.Decompress([]byte{})
	if err == nil {
		t.Fatal("expected error for empty input, got nil")
	}
}

func TestCompress_EmptyInput_RoundTrips(t *testing.T) {
	original := []byte{}
	compressed, err := compress.Compress(original)
	if err != nil {
		t.Fatalf("Compress empty: %v", err)
	}
	restored, err := compress.Decompress(compressed)
	if err != nil {
		t.Fatalf("Decompress empty: %v", err)
	}
	if !bytes.Equal(original, restored) {
		t.Errorf("round-trip mismatch: got %q, want %q", restored, original)
	}
}

func TestCompress_LargeInput_RoundTrips(t *testing.T) {
	original := []byte(strings.Repeat("envchain-export ", 1000))
	compressed, err := compress.Compress(original)
	if err != nil {
		t.Fatalf("Compress large: %v", err)
	}
	if len(compressed) >= len(original) {
		t.Errorf("expected compression to reduce size: compressed=%d original=%d", len(compressed), len(original))
	}
	restored, err := compress.Decompress(compressed)
	if err != nil {
		t.Fatalf("Decompress large: %v", err)
	}
	if !bytes.Equal(original, restored) {
		t.Error("large input round-trip mismatch")
	}
}

func TestIsCompressed_AfterCompress(t *testing.T) {
	data := []byte("some environment variable data")
	compressed, err := compress.Compress(data)
	if err != nil {
		t.Fatalf("Compress: %v", err)
	}
	if !compress.IsCompressed(compressed) {
		t.Error("expected IsCompressed to return true for compressed data")
	}
}

func TestIsCompressed_PlainText_ReturnsFalse(t *testing.T) {
	data := []byte("KEY=value\nFOO=bar\n")
	if compress.IsCompressed(data) {
		t.Error("expected IsCompressed to return false for plain text")
	}
}
