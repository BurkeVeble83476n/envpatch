package envfile

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

// Entry represents a single key-value pair in a .env file.
type Entry struct {
	Key     string
	Value   string
	Comment string // inline comment, if any
	Raw     string // original line as-is
}

// EnvFile holds all parsed entries from a .env file.
type EnvFile struct {
	Entries []Entry
	Index   map[string]int // key -> index in Entries
}

// Parse reads an .env file from r and returns a structured EnvFile.
func Parse(r io.Reader) (*EnvFile, error) {
	ef := &EnvFile{
		Index: make(map[string]int),
	}

	scanner := bufio.NewScanner(r)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		raw := scanner.Text()
		trimmed := strings.TrimSpace(raw)

		// Skip blank lines and full-line comments
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			ef.Entries = append(ef.Entries, Entry{Raw: raw})
			continue
		}

		eqIdx := strings.IndexByte(trimmed, '=')
		if eqIdx < 0 {
			return nil, fmt.Errorf("line %d: missing '=' in %q", lineNum, trimmed)
		}

		key := strings.TrimSpace(trimmed[:eqIdx])
		rest := trimmed[eqIdx+1:]

		value, comment := splitValueComment(rest)

		entry := Entry{
			Key:     key,
			Value:   unquote(value),
			Comment: comment,
			Raw:     raw,
		}

		ef.Index[key] = len(ef.Entries)
		ef.Entries = append(ef.Entries, entry)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scanner error: %w", err)
	}

	return ef, nil
}

// splitValueComment separates the value from an optional inline comment.
func splitValueComment(s string) (value, comment string) {
	if len(s) > 0 && (s[0] == '"' || s[0] == '\'') {
		quote := s[0]
		end := strings.IndexByte(s[1:], quote)
		if end >= 0 {
			return s[:end+2], strings.TrimSpace(s[end+2:])
		}
	}
	parts := strings.SplitN(s, " #", 2)
	if len(parts) == 2 {
		return strings.TrimSpace(parts[0]), "#" + parts[1]
	}
	return strings.TrimSpace(s), ""
}

// unquote strips surrounding single or double quotes from a value.
func unquote(s string) string {
	if len(s) >= 2 {
		if (s[0] == '"' && s[len(s)-1] == '"') || (s[0] == '\'' && s[len(s)-1] == '\'') {
			return s[1 : len(s)-1]
		}
	}
	return s
}
