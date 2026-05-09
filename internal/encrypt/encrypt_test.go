package encrypt_test

import (
	"strings"
	"testing"

	"github.com/yourorg/envpatch/internal/encrypt"
)

const testPass = "hunter2-super-secret"

func TestEncryptDecryptValue_RoundTrip(t *testing.T) {
	plain := "my-db-password"
	enc, err := encrypt.EncryptValue(testPass, plain)
	if err != nil {
		t.Fatalf("EncryptValue: %v", err)
	}
	if enc == plain {
		t.Fatal("ciphertext must differ from plaintext")
	}
	dec, err := encrypt.DecryptValue(testPass, enc)
	if err != nil {
		t.Fatalf("DecryptValue: %v", err)
	}
	if dec != plain {
		t.Fatalf("got %q, want %q", dec, plain)
	}
}

func TestEncryptValue_DifferentCiphertextsEachCall(t *testing.T) {
	a, _ := encrypt.EncryptValue(testPass, "value")
	b, _ := encrypt.EncryptValue(testPass, "value")
	if a == b {
		t.Fatal("expected different ciphertexts due to random nonce")
	}
}

func TestDecryptValue_WrongPassphrase(t *testing.T) {
	enc, _ := encrypt.EncryptValue(testPass, "secret")
	_, err := encrypt.DecryptValue("wrong-pass", enc)
	if err != encrypt.ErrInvalidCiphertext {
		t.Fatalf("expected ErrInvalidCiphertext, got %v", err)
	}
}

func TestDecryptValue_TamperedCiphertext(t *testing.T) {
	enc, _ := encrypt.EncryptValue(testPass, "secret")
	tampered := enc[:len(enc)-4] + "AAAA"
	_, err := encrypt.DecryptValue(testPass, tampered)
	if err != encrypt.ErrInvalidCiphertext {
		t.Fatalf("expected ErrInvalidCiphertext, got %v", err)
	}
}

func TestDecryptValue_InvalidBase64(t *testing.T) {
	_, err := encrypt.DecryptValue(testPass, "!!!not-base64!!!")
	if err != encrypt.ErrInvalidCiphertext {
		t.Fatalf("expected ErrInvalidCiphertext, got %v", err)
	}
}

func TestEncryptMap_EncryptsSelectedKeys(t *testing.T) {
	src := map[string]string{
		"DB_PASSWORD": "secret",
		"APP_NAME":    "envpatch",
	}
	out, err := encrypt.EncryptMap(testPass, src, []string{"DB_PASSWORD"})
	if err != nil {
		t.Fatalf("EncryptMap: %v", err)
	}
	if out["APP_NAME"] != "envpatch" {
		t.Errorf("APP_NAME should be unchanged")
	}
	if out["DB_PASSWORD"] == "secret" {
		t.Errorf("DB_PASSWORD should be encrypted")
	}
	if !strings.Contains(out["DB_PASSWORD"], "=") && len(out["DB_PASSWORD"]) < 20 {
		t.Errorf("DB_PASSWORD looks too short to be valid ciphertext")
	}
}

func TestDecryptMap_RoundTrip(t *testing.T) {
	src := map[string]string{
		"DB_PASSWORD": "secret",
		"API_KEY":     "key123",
		"APP_NAME":    "envpatch",
	}
	keys := []string{"DB_PASSWORD", "API_KEY"}
	enc, err := encrypt.EncryptMap(testPass, src, keys)
	if err != nil {
		t.Fatalf("EncryptMap: %v", err)
	}
	dec, err := encrypt.DecryptMap(testPass, enc, keys)
	if err != nil {
		t.Fatalf("DecryptMap: %v", err)
	}
	for k, want := range src {
		if got := dec[k]; got != want {
			t.Errorf("key %q: got %q, want %q", k, got, want)
		}
	}
}
