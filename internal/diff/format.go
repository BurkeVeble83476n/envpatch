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
	Color       bool
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

	colorize := func(color, line string) string {
		if opts.Color {
			return color + line + colorReset
		}
		return line
	}

	for _, c := range result.Changes {
		switch c.Type {
		case Added:
			line := colorize(colorGreen, fmt.Sprintf("+ %s=%s", c.Key, maskIfNeeded(c.Key, c.NewValue)))
			fmt.Fprintln(w, line)
		case Removed:
			line := colorize(colorRed, fmt.Sprintf("- %s=%s", c.Key, maskIfNeeded(c.Key, c.OldValue)))
			fmt.Fprintln(w, line)
		case Changed:
			old := maskIfNeeded(c.Key, c.OldValue)
			new := maskIfNeeded(c.Key, c.NewValue)
			line := colorize(colorYellow, fmt.Sprintf("~ %s: %s → %s", c.Key, old, new))
			fmt.Fprintln(w, line)
		}
	}
}
