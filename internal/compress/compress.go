// Package compress provides optional gzip compression and decompression
// for scope payloads before encryption, reducing storage size for large
// environment variable sets.
package compress

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
)

// Compress compresses src using gzip and returns the compressed bytes.
func Compress(src []byte) ([]byte, error) {
	var buf bytes.Buffer
	w := gzip.NewWriter(&buf)
	if _, err := w.Write(src); err != nil {
		return nil, fmt.Errorf("compress: write: %w", err)
	}
	if err := w.Close(); err != nil {
		return nil, fmt.Errorf("compress: close: %w", err)
	}
	return buf.Bytes(), nil
}

// Decompress decompresses gzip-compressed src and returns the original bytes.
func Decompress(src []byte) ([]byte, error) {
	r, err := gzip.NewReader(bytes.NewReader(src))
	if err != nil {
		return nil, fmt.Errorf("decompress: new reader: %w", err)
	}
	defer r.Close()

	out, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("decompress: read: %w", err)
	}
	return out, nil
}

// IsCompressed reports whether data begins with the gzip magic bytes.
func IsCompressed(data []byte) bool {
	return len(data) >= 2 && data[0] == 0x1f && data[1] == 0x8b
}
