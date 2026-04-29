// Package schema provides loading and parsing of .env schema files,
// which define required and optional keys along with their descriptions.
package schema

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

// KeySpec describes a single key in the schema.
type KeySpec struct {
	Key         string
	Required    bool
	Description string
}

// Schema is a map of key name to KeySpec.
type Schema map[string]KeySpec

// Load reads a schema definition from r.
//
// Schema file format (one entry per line):
//
//	KEY_NAME [required|optional] # optional description
//
// Lines starting with '#' or blank lines are ignored.
func Load(r io.Reader) (Schema, error) {
	if r == nil {
		return nil, fmt.Errorf("schema: reader must not be nil")
	}

	schema := make(Schema)
	scanner := bufio.NewScanner(r)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())

		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Strip inline comment
		desc := ""
		if idx := strings.Index(line, " #"); idx != -1 {
			desc = strings.TrimSpace(line[idx+2:])
			line = strings.TrimSpace(line[:idx])
		}

		parts := strings.Fields(line)
		if len(parts) == 0 {
			continue
		}

		key := parts[0]
		required := true // default to required

		if len(parts) >= 2 {
			switch strings.ToLower(parts[1]) {
			case "required":
				required = true
			case "optional":
				required = false
			default:
				return nil, fmt.Errorf("schema: line %d: unknown qualifier %q (expected required|optional)", lineNum, parts[1])
			}
		}

		schema[key] = KeySpec{
			Key:         key,
			Required:    required,
			Description: desc,
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("schema: scanning failed: %w", err)
	}

	return schema, nil
}
