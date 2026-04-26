package diff

import (
	"fmt"
	"io"
	"strings"
)

const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
)

// FormatOptions controls how diff output is rendered.
type FormatOptions struct {
	Color   bool
	MaskSecrets bool
	SecretKeys  []string
}

// Format writes a human-readable diff to w.
func Format(w io.Writer, result *Result, opts FormatOptions) {
	secretSet := make(map[string]struct{}, len(opts.SecretKeys))
	for _, k := range opts.SecretKeys {
		secretSet[strings.ToUpper(k)] = struct{}{}
	}

	maskIfNeeded := func(key, val string) string {
		if opts.MaskSecrets {
			if _, secret := secretSet[strings.ToUpper(key)]; secret {
				return "***"
			}
		}
		return val
	}

	for _, c := range result.Changes {
		switch c.Type {
		case Added:
			line := fmt.Sprintf("+ %s=%s", c.Key, maskIfNeeded(c.Key, c.NewValue))
			if opts.Color {
				line = colorGreen + line + colorReset
			}
			fmt.Fprintln(w, line)
		case Removed:
			line := fmt.Sprintf("- %s=%s", c.Key, maskIfNeeded(c.Key, c.OldValue))
			if opts.Color {
				line = colorRed + line + colorReset
			}
			fmt.Fprintln(w, line)
		case Changed:
			old := maskIfNeeded(c.Key, c.OldValue)
			new := maskIfNeeded(c.Key, c.NewValue)
			line := fmt.Sprintf("~ %s: %s → %s", c.Key, old, new)
			if opts.Color {
				line = colorYellow + line + colorReset
			}
			fmt.Fprintln(w, line)
		}
	}
}
