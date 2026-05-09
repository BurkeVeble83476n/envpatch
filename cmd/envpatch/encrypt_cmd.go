package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/yourorg/envpatch/internal/encrypt"
	"github.com/yourorg/envpatch/internal/envfile"
	"github.com/yourorg/envpatch/internal/export"
	"github.com/yourorg/envpatch/internal/redact"
)

// encryptCmd holds the parsed flags for the encrypt / decrypt sub-commands.
type encryptCmd struct {
	envFile    string
	passphrase string
	decrypt    bool
}

// run executes the encrypt (or decrypt) command and writes the result to
// stdout so callers can redirect it to a file.
func (c *encryptCmd) run() error {
	if c.passphrase == "" {
		return errors.New("encrypt: --passphrase must not be empty")
	}
	f, err := os.Open(c.envFile)
	if err != nil {
		return fmt.Errorf("encrypt: open %q: %w", c.envFile, err)
	}
	defer f.Close()

	env, err := envfile.Parse(f)
	if err != nil {
		return fmt.Errorf("encrypt: parse: %w", err)
	}

	// Determine which keys are sensitive using the shared redact heuristic.
	var sensitiveKeys []string
	for k := range env {
		if redact.IsSensitive(k, nil) {
			sensitiveKeys = append(sensitiveKeys, k)
		}
	}

	var result map[string]string
	if c.decrypt {
		result, err = encrypt.DecryptMap(c.passphrase, env, sensitiveKeys)
	} else {
		result, err = encrypt.EncryptMap(c.passphrase, env, sensitiveKeys)
	}
	if err != nil {
		return fmt.Errorf("encrypt: transform: %w", err)
	}

	out, err := export.Marshal(result, export.Options{Redact: false})
	if err != nil {
		return fmt.Errorf("encrypt: marshal: %w", err)
	}
	_, err = fmt.Fprint(os.Stdout, out)
	return err
}
