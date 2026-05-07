// Package history provides a persistent, append-only log of envpatch
// pipeline operations.
//
// Each time a pipeline runs — merging, patching, validating, or exporting
// environment files — callers can record an Entry that captures the
// operation name, a human-readable message, and an optional map of key/value
// changes involved.
//
// Logs are serialised as newline-delimited JSON so they can be stored
// alongside snapshots or shipped to an audit backend.
//
// Basic usage:
//
//	l := history.New()
//	l.Record("merge", "merged .env.dev into .env.staging", changes)
//
//	var buf bytes.Buffer
//	if err := l.Save(&buf); err != nil { ... }
//
//	loaded, err := history.Load(&buf)
package history
