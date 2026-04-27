// Package redact provides utilities for masking sensitive values
// in .env file entries before they are displayed or logged.
package redact

import "strings"

// DefaultSecretPatterns holds common key substrings that indicate a secret.
var DefaultSecretPatterns = []string{
	"SECRET",
	"PASSWORD",
	"PASSWD",
	"TOKEN",
	"API_KEY",
	"APIKEY",
	"PRIVATE",
	"CREDENTIAL",
	"AUTH",
	"ACCESS_KEY",
}

const masked = "***"

// IsSensitive reports whether the given key name matches any of the
// provided patterns (case-insensitive substring match).
func IsSensitive(key string, patterns []string) bool {
	upper := strings.ToUpper(key)
	for _, p := range patterns {
		if strings.Contains(upper, strings.ToUpper(p)) {
			return true
		}
	}
	return false
}

// Value returns the original value if the key is not sensitive, or the
// masked placeholder if it is. patterns may be nil, in which case
// DefaultSecretPatterns is used.
func Value(key, value string, patterns []string) string {
	if patterns == nil {
		patterns = DefaultSecretPatterns
	}
	if IsSensitive(key, patterns) {
		return masked
	}
	return value
}

// Map returns a new map with sensitive values replaced by the masked
// placeholder. The original map is never modified.
func Map(env map[string]string, patterns []string) map[string]string {
	if patterns == nil {
		patterns = DefaultSecretPatterns
	}
	out := make(map[string]string, len(env))
	for k, v := range env {
		out[k] = Value(k, v, patterns)
	}
	return out
}
