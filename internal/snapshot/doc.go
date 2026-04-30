// Package snapshot provides point-in-time capture and drift detection for
// .env configurations.
//
// # Overview
//
// A Snapshot records the full key-value state of an environment at a specific
// moment, serialised as JSON. Snapshots can be saved to any io.Writer (a file,
// a buffer, etc.) and later reloaded for comparison.
//
// # Drift Detection
//
// Compare returns a slice of Drift values describing every key that was added,
// removed, or changed between the snapshot and a live env map. This makes it
// straightforward to alert when production configuration drifts from a known
// good baseline.
//
// # Usage
//
//	env, _ := envfile.Parse(r)
//	s, _ := snapshot.New(env, "prod-baseline")
//	snapshot.Save(s, file)
//
//	// later …
//	s2, _ := snapshot.Load(file)
//	drifts, _ := snapshot.Compare(s2, currentEnv)
package snapshot
