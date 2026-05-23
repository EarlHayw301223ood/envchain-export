package compress_test

import (
	"bytes"
	"testing"

	"github.com/nicholasgasior/envchain-export/internal/compress"
)

func TestCompressDecompress_RoundTrip(t *testing.T) {
	orig := []byte("DATABASE_URL=postgres://localhost/mydb\nSECRET_KEY=supersecret\n")

	compressed, err := compress.Compress(orig)
	if err != nil {
		t.Fatalf("Compress: %v", err)
	}

	if bytes.Equal(compressed, orig) {
		t.Fatal("expected compressed output to differ from input")
	}

	got, err := compress.Decompress(compressed)
	if err != nil {
		t.Fatalf("Decompress: %v", err)
	}

	if !bytes.Equal(got, orig) {
		t.Fatalf("round-trip mismatch: got %q, want %q", got, orig)
	}
}

func TestCompress_ReducesSize(t *testing.T) {
	// Highly repetitive data compresses well.
	orig := bytes.Repeat([]byte("KEY=value\n"), 100)

	compressed, err := compress.Compress(orig)
	if err != nil {
		t.Fatalf("Compress: %v", err)
	}

	if len(compressed) >= len(orig) {
		t.Errorf("expected compressed size (%d) < original size (%d)", len(compressed), len(orig))
	}
}

func TestIsCompressed_True(t *testing.T) {
	data, err := compress.Compress([]byte("hello"))
	if err != nil {
		t.Fatalf("Compress: %v", err)
	}
	if !compress.IsCompressed(data) {
		t.Fatal("expected IsCompressed to return true for gzip data")
	}
}

func TestIsCompressed_False(t *testing.T) {
	if compress.IsCompressed([]byte("plain text")) {
		t.Fatal("expected IsCompressed to return false for plain text")
	}
}

func TestIsCompressed_ShortSlice(t *testing.T) {
	if compress.IsCompressed([]byte{0x1f}) {
		t.Fatal("expected IsCompressed to return false for single-byte slice")
	}
}

func TestDecompress_InvalidData_ReturnsError(t *testing.T) {
	_, err := compress.Decompress([]byte("this is not gzip data"))
	if err == nil {
		t.Fatal("expected error decompressing invalid data")
	}
}
