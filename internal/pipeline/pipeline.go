package pipeline

import (
	"fmt"

	"github.com/yourorg/envpatch/internal/diff"
	"github.com/yourorg/envpatch/internal/envfile"
	"github.com/yourorg/envpatch/internal/export"
	"github.com/yourorg/envpatch/internal/merge"
	"github.com/yourorg/envpatch/internal/patch"
	"github.com/yourorg/envpatch/internal/validate"
)

// Result holds the outcome of a pipeline run.
type Result struct {
	Env    map[string]string
	Diff   []diff.Change
	Issues []validate.Issue
}

// Options configures a pipeline execution.
type Options struct {
	BaseSource    string
	OverlaySource string
	PatchOps      []patch.Op
	Schema        map[string]validate.Rule
	MergeStrategy merge.Strategy
	ExportHeader  string
	RedactOutput  bool
}

// Run executes the full envpatch pipeline: parse → merge → patch → validate → diff.
func Run(base, overlay map[string]string, opts Options) (*Result, error) {
	if base == nil {
		return nil, fmt.Errorf("pipeline: base env must not be nil")
	}

	// Merge overlay into base.
	merged, err := merge.Merge(base, overlay, opts.MergeStrategy)
	if err != nil {
		return nil, fmt.Errorf("pipeline: merge failed: %w", err)
	}

	// Apply patch operations.
	patched, err := patch.Apply(merged, opts.PatchOps)
	if err != nil {
		return nil, fmt.Errorf("pipeline: patch failed: %w", err)
	}

	// Validate against schema if provided.
	var issues []validate.Issue
	if opts.Schema != nil {
		issues, err = validate.Validate(opts.Schema, patched)
		if err != nil {
			return nil, fmt.Errorf("pipeline: validate failed: %w", err)
		}
	}

	// Compute diff between original base and final result.
	changes := diff.Diff(base, patched)

	return &Result{
		Env:    patched,
		Diff:   changes,
		Issues: issues,
	}, nil
}

// RunFromStrings parses raw .env content and executes the pipeline.
func RunFromStrings(baseRaw, overlayRaw string, opts Options) (*Result, error) {
	base, err := envfile.Parse(baseRaw)
	if err != nil {
		return nil, fmt.Errorf("pipeline: parse base: %w", err)
	}

	var overlay map[string]string
	if overlayRaw != "" {
		overlay, err = envfile.Parse(overlayRaw)
		if err != nil {
			return nil, fmt.Errorf("pipeline: parse overlay: %w", err)
		}
	}

	return Run(base, overlay, opts)
}

// Marshal serialises a Result's Env map to .env format.
func Marshal(r *Result, header string, redact bool) (string, error) {
	if r == nil {
		return "", fmt.Errorf("pipeline: result must not be nil")
	}
	return export.Marshal(r.Env, header, redact)
}
