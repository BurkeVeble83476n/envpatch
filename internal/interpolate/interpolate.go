// Package interpolate resolves cross-file variable references between two
// env maps. A value in the overlay may reference a key from the base using
// the ${BASE:KEY} syntax, which is substituted at resolution time.
package interpolate

import (
	"fmt"
	"regexp"
	"strings"
)

// baseSyntax matches ${BASE:SOME_KEY} placeholders.
var baseSyntax = regexp.MustCompile(`\$\{BASE:([A-Z0-9_]+)\}`)

// Result holds the resolved env map and any warnings produced during
// interpolation (e.g. unresolved BASE references that were left in place).
type Result struct {
	Env      map[string]string
	Warnings []string
}

// Resolve walks every value in overlay and substitutes ${BASE:KEY} tokens
// with the corresponding value from base. Keys that are not found in base
// are left unchanged and a warning is appended to Result.Warnings.
//
// Neither base nor overlay is mutated; a fresh map is returned inside Result.
func Resolve(base, overlay map[string]string) (*Result, error) {
	if base == nil {
		return nil, fmt.Errorf("interpolate: base map must not be nil")
	}
	if overlay == nil {
		return nil, fmt.Errorf("interpolate: overlay map must not be nil")
	}

	out := make(map[string]string, len(overlay))
	var warnings []string

	for k, v := range overlay {
		resolved, w := resolveValue(v, base)
		out[k] = resolved
		warnings = append(warnings, w...)
	}

	return &Result{Env: out, Warnings: warnings}, nil
}

// resolveValue replaces all ${BASE:KEY} tokens in s using the provided lookup
// map. It returns the substituted string and any warnings for missing keys.
func resolveValue(s string, lookup map[string]string) (string, []string) {
	var warnings []string

	result := baseSyntax.ReplaceAllStringFunc(s, func(match string) string {
		key := baseSyntax.FindStringSubmatch(match)[1]
		if val, ok := lookup[key]; ok {
			return val
		}
		warnings = append(warnings, fmt.Sprintf("unresolved BASE reference: %s", key))
		return match // leave token intact
	})

	_ = strings.Contains // satisfy import if needed
	return result, warnings
}
