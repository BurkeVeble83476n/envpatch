// Package export provides functionality for serialising env maps back
// to .env file format, with optional redaction of sensitive values.
package export

import (
	"fmt"
	"sort"
	"strings"

	"github.com/yourorg/envpatch/internal/redact"
)

// Options controls how the exported output is rendered.
type Options struct {
	// RedactSensitive replaces sensitive values with a masked placeholder.
	RedactSensitive bool
	// CustomPatterns extends the default sensitive-key detection patterns.
	CustomPatterns []string
	// Header is an optional comment block prepended to the output.
	Header string
}

// Marshal converts an env map into a .env-formatted byte slice.
// Keys are written in sorted order for deterministic output.
// If opts is nil, default options (no redaction, no header) are used.
func Marshal(env map[string]string, opts *Options) ([]byte, error) {
	if env == nil {
		return nil, fmt.Errorf("export: env map must not be nil")
	}

	if opts == nil {
		opts = &Options{}
	}

	var sb strings.Builder

	if opts.Header != "" {
		for _, line := range strings.Split(opts.Header, "\n") {
			sb.WriteString("# ")
			sb.WriteString(line)
			sb.WriteByte('\n')
		}
		sb.WriteByte('\n')
	}

	keys := make([]string, 0, len(env))
	for k := range env {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		v := env[k]
		if opts.RedactSensitive {
			v = redact.Value(k, v, opts.CustomPatterns...)
		}
		v = quoteIfNeeded(v)
		fmt.Fprintf(&sb, "%s=%s\n", k, v)
	}

	return []byte(sb.String()), nil
}

// quoteIfNeeded wraps the value in double quotes when it contains
// whitespace, a hash character, or is the empty string.
func quoteIfNeeded(v string) string {
	if v == "" {
		return `""`
	}
	if strings.ContainsAny(v, " \t\n#") {
		escaped := strings.ReplaceAll(v, `"`, `\"`)
		return `"` + escaped + `"`
	}
	return v
}
