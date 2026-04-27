// Package patch provides primitives for applying a declarative set of changes
// to an environment variable map.
//
// # Overview
//
// A [Patch] is an ordered list of [Entry] values. Each entry carries an [Op]
// that is either [OpSet] (add or overwrite a key) or [OpRemove] (delete a
// key). Entries are applied in order so that later entries can overwrite
// earlier ones.
//
// # Usage
//
//	out, err := patch.Apply(base, patch.Patch{
//		Entries: []patch.Entry{
//			{Op: patch.OpSet,    Key: "LOG_LEVEL", Value: "debug"},
//			{Op: patch.OpRemove, Key: "LEGACY_FLAG"},
//		},
//	})
//
// Apply never mutates the base map; it always returns a new map.
package patch
