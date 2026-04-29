// Package pipeline provides a high-level orchestration layer for the envpatch
// workflow. It composes the parse, merge, patch, validate, and diff stages into
// a single, easy-to-use API.
//
// # Typical usage
//
//	result, err := pipeline.Run(base, overlay, pipeline.Options{
//		MergeStrategy: merge.Overlay,
//		PatchOps:      ops,
//		Schema:        schema,
//	})
//
// The returned Result contains the final environment map, a list of diff
// changes relative to the original base, and any validation issues raised
// by the schema check.
//
// RunFromStrings is a convenience wrapper that parses raw .env content before
// delegating to Run, making it straightforward to drive the pipeline directly
// from file bytes or test strings.
//
// Marshal serialises a Result back to .env format, optionally redacting
// sensitive values before output.
package pipeline
