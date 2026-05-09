// Package profile provides a named-profile store for envpatch.
//
// A Profile is a labelled snapshot of an environment variable map
// (e.g. "development", "staging", "production"). Profiles are grouped
// inside a Store which can be serialised to / deserialised from JSON,
// making it straightforward to persist them alongside the project.
//
// Typical usage:
//
//	store := profile.NewStore()
//	_ = store.Set("staging", envMap)
//
//	// later …
//	p, err := store.Get("staging")
//	if err != nil { … }
//	fmt.Println(p.Env)
//
// Profiles are independent of secrets: callers should redact sensitive
// values (see internal/redact) before persisting a Store to disk.
package profile
