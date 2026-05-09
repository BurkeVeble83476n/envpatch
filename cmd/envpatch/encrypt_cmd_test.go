package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeEncTmp(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0o600); err != nil {
		t.Fatalf("writeEncTmp: %v", err)
	}
	return p
}

func TestEncryptCmd_EmptyPassphraseReturnsError(t *testing.T) {
	cmd := &encryptCmd{
		envFile:    writeEncTmp(t, "APP_NAME=envpatch\n"),
		passphrase: "",
	}
	if err := cmd.run(); err == nil {
		t.Fatal("expected error for empty passphrase")
	}
}

func TestEncryptCmd_MissingFileReturnsError(t *testing.T) {
	cmd := &encryptCmd{
		envFile:    "/nonexistent/.env",
		passphrase: "pass",
	}
	if err := cmd.run(); err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestEncryptCmd_NonSensitiveKeyPassesThrough(t *testing.T) {
	path := writeEncTmp(t, "APP_NAME=envpatch\n")
	// Capture stdout by redirecting os.Stdout.
	r, w, _ := os.Pipe()
	old := os.Stdout
	os.Stdout = w

	cmd := &encryptCmd{envFile: path, passphrase: "secret"}
	if err := cmd.run(); err != nil {
		w.Close()
		os.Stdout = old
		t.Fatalf("run: %v", err)
	}
	w.Close()
	os.Stdout = old

	buf := make([]byte, 512)
	n, _ := r.Read(buf)
	output := string(buf[:n])

	if !strings.Contains(output, "APP_NAME=envpatch") {
		t.Errorf("expected APP_NAME to be unchanged, got: %q", output)
	}
}

func TestEncryptCmd_SensitiveKeyIsTransformed(t *testing.T) {
	path := writeEncTmp(t, "DB_PASSWORD=hunter2\nAPP_NAME=envpatch\n")
	r, w, _ := os.Pipe()
	old := os.Stdout
	os.Stdout = w

	cmd := &encryptCmd{envFile: path, passphrase: "my-pass"}
	if err := cmd.run(); err != nil {
		w.Close()
		os.Stdout = old
		t.Fatalf("run: %v", err)
	}
	w.Close()
	os.Stdout = old

	buf := make([]byte, 1024)
	n, _ := r.Read(buf)
	output := string(buf[:n])

	if strings.Contains(output, "hunter2") {
		t.Errorf("plaintext password must not appear in output: %q", output)
	}
	if !strings.Contains(output, "APP_NAME=envpatch") {
		t.Errorf("non-sensitive key should be unchanged: %q", output)
	}
}
