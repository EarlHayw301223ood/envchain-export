// Package compress provides transparent gzip compression and decompression
// helpers used when persisting scope payloads to disk.
//
// Compression is applied before encryption so that the plaintext size is
// reduced prior to being sealed. Decompression is applied after decryption
// when loading a scope. The IsCompressed helper allows callers to detect
// legacy (uncompressed) payloads and handle them gracefully.
package compress
