// Package template provides variable substitution for .env file values,
// replacing {{VAR}} or ${VAR} placeholders with values from a context map.
package template

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	// bracePattern matches {{VAR_NAME}} style references.
	bracePattern = regexp.MustCompile(`\{\{([A-Z0-9_]+)\}\}`)
	// shellPattern matches ${VAR_NAME} style references.
	shellPattern = regexp.MustCompile(`\$\{([A-Z0-9_]+)\}`)
)

// ErrUnresolved is returned when a placeholder cannot be resolved from the context.
type ErrUnresolved struct {
	Key string
	Placeholder string
}

func (e *ErrUnresolved) Error() string {
	return fmt.Sprintf("template: unresolved placeholder %q for key %q", e.Placeholder, e.Key)
}

// Expand replaces all {{VAR}} and ${VAR} placeholders in env with values
// looked up from ctx. Missing keys are collected and returned as errors;
// resolved values are written to a new map.
func Expand(env map[string]string, ctx map[string]string) (map[string]string, []error) {
	if env == nil {
		return nil, []error{fmt.Errorf("template: env map must not be nil")}
	}
	if ctx == nil {
		ctx = map[string]string{}
	}

	out := make(map[string]string, len(env))
	var errs []error

	for k, v := range env {
		resolved, keyErrs := expandValue(k, v, ctx)
		out[k] = resolved
		errs = append(errs, keyErrs...)
	}
	return out, errs
}

// expandValue substitutes all placeholders in a single value string.
func expandValue(key, value string, ctx map[string]string) (string, []error) {
	var errs []error

	replace := func(match, varName string) string {
		if val, ok := ctx[varName]; ok {
			return val
		}
		errs = append(errs, &ErrUnresolved{Key: key, Placeholder: match})
		return match
	}

	result := bracePattern.ReplaceAllStringFunc(value, func(m string) string {
		sub := bracePattern.FindStringSubmatch(m)
		return replace(m, sub[1])
	})
	result = shellPattern.ReplaceAllStringFunc(result, func(m string) string {
		sub := shellPattern.FindStringSubmatch(m)
		return replace(m, sub[1])
	})

	_ = strings.TrimSpace // imported for future use
	return result, errs
}
