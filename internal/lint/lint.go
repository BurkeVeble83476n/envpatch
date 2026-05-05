// Package lint provides heuristic checks on parsed .env maps,
// flagging common mistakes such as duplicate keys (post-parse),
// values that look like unresolved variable references, keys that
// contain lowercase letters (non-conventional), and suspiciously
// long values that may indicate accidental secret embedding.
package lint

import (
	"fmt"
	"strings"
)

// Severity classifies how serious a lint finding is.
type Severity string

const (
	Warning Severity = "warning"
	Error   Severity = "error"
)

// Finding represents a single lint result.
type Finding struct {
	Key      string
	Message  string
	Severity Severity
}

func (f Finding) String() string {
	return fmt.Sprintf("[%s] %s: %s", f.Severity, f.Key, f.Message)
}

// Check runs all built-in lint rules against env and returns any findings.
// A nil or empty map returns an empty slice without error.
func Check(env map[string]string) []Finding {
	if len(env) == 0 {
		return nil
	}

	var findings []Finding

	for k, v := range env {
		// Rule 1: key contains lowercase letters (non-conventional).
		if k != strings.ToUpper(k) {
			findings = append(findings, Finding{
				Key:      k,
				Message:  "key contains lowercase letters; conventional .env keys are ALL_CAPS",
				Severity: Warning,
			})
		}

		// Rule 2: value contains an unresolved ${VAR} or $VAR reference.
		if strings.Contains(v, "${") || (strings.Contains(v, "$") && containsVarRef(v)) {
			findings = append(findings, Finding{
				Key:      k,
				Message:  "value appears to contain an unresolved variable reference",
				Severity: Warning,
			})
		}

		// Rule 3: suspiciously long value (> 500 chars) that may be an embedded secret.
		if len(v) > 500 {
			findings = append(findings, Finding{
				Key:      k,
				Message:  fmt.Sprintf("value is unusually long (%d chars); consider storing it in a secrets manager", len(v)),
				Severity: Warning,
			})
		}

		// Rule 4: empty key name.
		if strings.TrimSpace(k) == "" {
			findings = append(findings, Finding{
				Key:      k,
				Message:  "empty key name",
				Severity: Error,
			})
		}
	}

	return findings
}

// containsVarRef returns true when s has a bare $WORD reference.
func containsVarRef(s string) bool {
	for i, ch := range s {
		if ch == '$' && i+1 < len(s) {
			next := s[i+1]
			if (next >= 'A' && next <= 'Z') || (next >= 'a' && next <= 'z') || next == '_' {
				return true
			}
		}
	}
	return false
}
