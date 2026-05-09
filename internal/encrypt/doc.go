// Package encrypt provides AES-256-GCM encryption and decryption helpers
// for individual .env values and for entire env maps.
//
// # Overview
//
// Sensitive values — passwords, API keys, tokens — should never be stored in
// plain text alongside non-sensitive configuration.  This package lets callers
// encrypt a chosen subset of keys before persisting or transmitting an env map
// and decrypt them again when the values are needed at runtime.
//
// # Key derivation
//
// A 32-byte AES-256 key is derived from the caller-supplied passphrase via
// SHA-256.  For higher-security deployments, wrap EncryptValue / DecryptValue
// with a proper KDF (e.g. Argon2id) before passing the derived key as the
// passphrase argument.
//
// # Ciphertext format
//
// Each ciphertext is a base64url-encoded blob of the form:
//
//	[ 12-byte random nonce ][ AES-GCM ciphertext + 16-byte auth tag ]
//
// The random nonce ensures that encrypting the same plaintext twice produces
// different output, preventing value-equality leaks.
package encrypt
