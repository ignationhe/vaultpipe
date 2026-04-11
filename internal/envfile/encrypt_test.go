package envfile

import (
	"testing"
)

func TestEncrypt_RoundTrip(t *testing.T) {
	secrets := map[string]string{
		"DB_PASSWORD": "s3cr3t",
		"API_KEY":     "abc123",
		"EMPTY_VAL":   "",
	}
	passphrase := "my-strong-passphrase"

	encrypted, err := Encrypt(secrets, passphrase)
	if err != nil {
		t.Fatalf("Encrypt: unexpected error: %v", err)
	}

	for k, v := range secrets {
		if encrypted[k] == v {
			t.Errorf("key %q: encrypted value should differ from plaintext", k)
		}
	}

	decrypted, err := Decrypt(encrypted, passphrase)
	if err != nil {
		t.Fatalf("Decrypt: unexpected error: %v", err)
	}

	for k, want := range secrets {
		if got := decrypted[k]; got != want {
			t.Errorf("key %q: got %q, want %q", k, got, want)
		}
	}
}

func TestEncrypt_DifferentNonceEachCall(t *testing.T) {
	secrets := map[string]string{"KEY": "value"}
	passphrase := "passphrase"

	a, _ := Encrypt(secrets, passphrase)
	b, _ := Encrypt(secrets, passphrase)

	if a["KEY"] == b["KEY"] {
		t.Error("expected different ciphertext on each call due to random nonce")
	}
}

func TestEncrypt_EmptyPassphraseReturnsError(t *testing.T) {
	_, err := Encrypt(map[string]string{"K": "v"}, "")
	if err == nil {
		t.Fatal("expected error for empty passphrase, got nil")
	}
}

func TestDecrypt_EmptyPassphraseReturnsError(t *testing.T) {
	_, err := Decrypt(map[string]string{"K": "v"}, "")
	if err == nil {
		t.Fatal("expected error for empty passphrase, got nil")
	}
}

func TestDecrypt_WrongPassphraseReturnsError(t *testing.T) {
	secrets := map[string]string{"TOKEN": "secret-token"}

	encrypted, err := Encrypt(secrets, "correct-pass")
	if err != nil {
		t.Fatalf("Encrypt: %v", err)
	}

	_, err = Decrypt(encrypted, "wrong-pass")
	if err == nil {
		t.Fatal("expected error when decrypting with wrong passphrase")
	}
}

func TestDecrypt_InvalidBase64ReturnsError(t *testing.T) {
	_, err := Decrypt(map[string]string{"K": "!!!not-base64!!!"}, "pass")
	if err == nil {
		t.Fatal("expected error for invalid base64, got nil")
	}
}

func TestDecrypt_TruncatedCiphertextReturnsError(t *testing.T) {
	// A valid base64 string that is too short to contain a nonce.
	_, err := Decrypt(map[string]string{"K": "dG9vc2hvcnQ="}, "pass")
	if err == nil {
		t.Fatal("expected error for truncated ciphertext, got nil")
	}
}
