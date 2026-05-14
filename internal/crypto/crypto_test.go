package crypto_test

import (
	"bytes"
	"testing"

	"github.com/user/envchain-export/internal/crypto"
)

func TestEncryptDecrypt_RoundTrip(t *testing.T) {
	passphrase := "supersecret"
	plaintext := []byte("MY_VAR=hello world")

	ciphertext, err := crypto.Encrypt(passphrase, plaintext)
	if err != nil {
		t.Fatalf("Encrypt error: %v", err)
	}

	if bytes.Equal(ciphertext, plaintext) {
		t.Fatal("ciphertext should not equal plaintext")
	}

	decrypted, err := crypto.Decrypt(passphrase, ciphertext)
	if err != nil {
		t.Fatalf("Decrypt error: %v", err)
	}

	if !bytes.Equal(decrypted, plaintext) {
		t.Fatalf("expected %q, got %q", plaintext, decrypted)
	}
}

func TestEncrypt_ProducesUniqueOutput(t *testing.T) {
	passphrase := "supersecret"
	plaintext := []byte("MY_VAR=hello")

	c1, _ := crypto.Encrypt(passphrase, plaintext)
	c2, _ := crypto.Encrypt(passphrase, plaintext)

	if bytes.Equal(c1, c2) {
		t.Fatal("two encryptions of the same plaintext should differ (random nonce)")
	}
}

func TestDecrypt_WrongPassphrase(t *testing.T) {
	ciphertext, err := crypto.Encrypt("correct-passphrase", []byte("secret"))
	if err != nil {
		t.Fatalf("Encrypt error: %v", err)
	}

	_, err = crypto.Decrypt("wrong-passphrase", ciphertext)
	if err == nil {
		t.Fatal("expected error when decrypting with wrong passphrase")
	}
}

func TestDecrypt_TooShortCiphertext(t *testing.T) {
	_, err := crypto.Decrypt("passphrase", []byte("short"))
	if err == nil {
		t.Fatal("expected error for too-short ciphertext")
	}
}
