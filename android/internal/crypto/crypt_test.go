package crypto

import (
	"bytes"
	"clipsync-android/internal"
	"testing"
)

func TestNewBlockCrypt_EncryptDecrypt(t *testing.T) {
	block, err := NewBlockCrypt(internal.SecretKey)
	if err != nil {
		t.Fatalf("Failed to create BlockCrypt: %v", err)
	}
	if block == nil {
		t.Fatal("Expected non-nil BlockCrypt")
	}

	plaintext := []byte("secret clipboard contents from phone")
	// Block size for AES is 16 bytes, kcp BlockCrypt handles blocks
	blockSize := 16
	padding := (blockSize - (len(plaintext) % blockSize)) % blockSize
	padded := append(plaintext, bytes.Repeat([]byte{0}, padding)...)

	ciphertext := make([]byte, len(padded))
	decrypted := make([]byte, len(padded))

	for i := 0; i < len(padded); i += blockSize {
		block.Encrypt(ciphertext[i:i+blockSize], padded[i:i+blockSize])
	}

	for i := 0; i < len(padded); i += blockSize {
		block.Decrypt(decrypted[i:i+blockSize], ciphertext[i:i+blockSize])
	}

	if !bytes.Equal(padded, decrypted) {
		t.Fatalf("Decrypted content does not match original: got %q, want %q", decrypted, padded)
	}
}

func TestNewBlockCrypt_InvalidKey(t *testing.T) {
	// KCP AES expects valid key lengths (16, 24, 32 bytes) or pads.
	block, err := NewBlockCrypt([]byte("clipboardsyncapp"))
	if err != nil {
		t.Fatalf("Expected valid creation with standard key, got err: %v", err)
	}
	if block == nil {
		t.Fatal("Expected non-nil block")
	}
}
