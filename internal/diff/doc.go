// Package diff compares two parsed .env file maps and produces a structured
// result describing which keys were added, removed, or changed.
//
// Basic usage:
//
//	base, _ := envfile.Parse(baseReader)
//	target, _ := envfile.Parse(targetReader)
//
//	result := diff.Diff(base, target)
//	if result.HasChanges() {
//		diff.Format(os.Stdout, result, diff.FormatOptions{
//			Color:       true,
//			MaskSecrets: true,
//			SecretKeys:  []string{"DB_PASSWORD", "API_KEY"},
//		})
//	}
//
// Change types:
//   - Added:   key exists in target but not in base
//   - Removed: key exists in base but not in target
//   - Changed: key exists in both but values differ
//
// The Format function renders the diff to any io.Writer with optional
// ANSI color highlighting and secret masking.
package diff
