// Package export serialises an in-memory env map back to the standard
// .env file format.
//
// # Overview
//
// After diffing, merging, or patching an env map with the other envpatch
// packages, export.Marshal writes the result to a deterministic,
// human-readable .env byte slice that can be written directly to disk.
//
// # Key ordering
//
// Keys are always emitted in lexicographic order so that the output is
// stable across runs and produces clean version-control diffs.
//
// # Quoting
//
// Values that are empty or contain whitespace or hash characters are
// automatically wrapped in double quotes.  Embedded double-quote
// characters are escaped with a backslash.
//
// # Redaction
//
// When Options.RedactSensitive is true, values whose keys match the
// default sensitive patterns (password, secret, token, key, …) are
// replaced with the redact package placeholder before writing.  Custom
// patterns can be supplied via Options.CustomPatterns.
package export
