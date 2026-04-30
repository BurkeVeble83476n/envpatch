// Package audit provides structured audit logging for envpatch operations.
//
// An audit Logger records timestamped entries for significant lifecycle events
// such as merges, patches, validations, and snapshot comparisons. Each entry
// captures the event kind, a human-readable message, and the list of key-level
// changes that occurred.
//
// Usage:
//
//	var buf bytes.Buffer
//	logger, err := audit.New(&buf)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	changes := diff.Diff(base, overlay)
//	logger.Record(audit.KindMerge, "merged production overlay", changes)
//
// Output is written in a line-oriented human-readable format suitable for
// appending to a log file or streaming to stdout. Each entry is followed by
// one line per changed key.
package audit
