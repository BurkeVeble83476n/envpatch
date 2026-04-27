// Package validate provides functionality for validating .env files
// against a schema or reference environment definition.
package validate

import (
	"fmt"
	"sort"
)

// Severity indicates how serious a validation issue is.
type Severity string

const (
	SeverityError   Severity = "error"
	SeverityWarning Severity = "warning"
)

// Issue represents a single validation finding for a key.
type Issue struct {
	Key      string
	Message  string
	Severity Severity
}

func (i Issue) Error() string {
	return fmt.Sprintf("%s [%s]: %s", i.Severity, i.Key, i.Message)
}

// Result holds all issues found during validation.
type Result struct {
	Issues []Issue
}

// HasErrors returns true if any issue has error severity.
func (r *Result) HasErrors() bool {
	for _, issue := range r.Issues {
		if issue.Severity == SeverityError {
			return true
		}
	}
	return false
}

// Validate checks the target env map against the schema env map.
// Keys present in schema but missing in target are errors.
// Keys present in target but missing in schema are warnings.
func Validate(schema, target map[string]string) (*Result, error) {
	if schema == nil {
		return nil, fmt.Errorf("validate: schema must not be nil")
	}
	if target == nil {
		return nil, fmt.Errorf("validate: target must not be nil")
	}

	result := &Result{}

	schemaKeys := sortedKeys(schema)
	for _, key := range schemaKeys {
		if _, ok := target[key]; !ok {
			result.Issues = append(result.Issues, Issue{
				Key:      key,
				Message:  "required key is missing from target",
				Severity: SeverityError,
			})
		}
	}

	targetKeys := sortedKeys(target)
	for _, key := range targetKeys {
		if _, ok := schema[key]; !ok {
			result.Issues = append(result.Issues, Issue{
				Key:      key,
				Message:  "key is not defined in schema",
				Severity: SeverityWarning,
			})
		}
	}

	return result, nil
}

func sortedKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
