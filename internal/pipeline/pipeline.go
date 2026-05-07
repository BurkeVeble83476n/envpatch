// Package pipeline wires together parse, merge, patch, validate, and export
// into a single high-level operation.
package pipeline

import (
	"fmt"
	"io"
	"strings"

	"github.com/yourorg/envpatch/internal/diff"
	"github.com/yourorg/envpatch/internal/envfile"
	"github.com/yourorg/envpatch/internal/export"
	"github.com/yourorg/envpatch/internal/merge"
	"github.com/yourorg/envpatch/internal/patch"
	"github.com/yourorg/envpatch/internal/template"
	"github.com/yourorg/envpatch/internal/validate"
)

// Options controls pipeline behaviour.
type Options struct {
	Header       string
	RedactKeys   []string
	PatchOps     []patch.Op
	Schema       map[string]validate.KeySpec
	// TemplateCtx, when non-nil, triggers variable expansion after merge.
	TemplateCtx  map[string]string
}

// Result holds the outputs produced by Run.
type Result struct {
	Env              map[string]string
	Diff             []diff.Change
	ValidationIssues []validate.Issue
	TemplateErrors   []error
}

// Run merges base and overlay, applies patches, optionally expands templates,
// validates against a schema, and returns a Result.
func Run(base, overlay map[string]string, opts Options) (*Result, error) {
	if base == nil {
		return nil, fmt.Errorf("pipeline: base map must not be nil")
	}

	merged, err := merge.Merge(base, overlay, merge.OverlayStrategy)
	if err != nil {
		return nil, fmt.Errorf("pipeline: merge failed: %w", err)
	}

	if len(opts.PatchOps) > 0 {
		merged, err = patch.Apply(merged, opts.PatchOps)
		if err != nil {
			return nil, fmt.Errorf("pipeline: patch failed: %w", err)
		}
	}

	var tmplErrs []error
	if opts.TemplateCtx != nil {
		merged, tmplErrs = template.Expand(merged, opts.TemplateCtx)
	}

	changes := diff.Diff(base, merged)

	var issues []validate.Issue
	if opts.Schema != nil {
		issues, err = validate.Validate(opts.Schema, merged)
		if err != nil {
			return nil, fmt.Errorf("pipeline: validate failed: %w", err)
		}
	}

	return &Result{
		Env:              merged,
		Diff:             changes,
		ValidationIssues: issues,
		TemplateErrors:   tmplErrs,
	}, nil
}

// RunFromStrings parses base and overlay from raw .env strings before running
// the pipeline.
func RunFromStrings(baseRaw, overlayRaw string, opts Options) (*Result, error) {
	base, err := envfile.Parse(strings.NewReader(baseRaw))
	if err != nil {
		return nil, fmt.Errorf("pipeline: parse base: %w", err)
	}
	var overlay map[string]string
	if overlayRaw != "" {
		overlay, err = envfile.Parse(strings.NewReader(overlayRaw))
		if err != nil {
			return nil, fmt.Errorf("pipeline: parse overlay: %w", err)
		}
	}
	return Run(base, overlay, opts)
}

// Marshal serialises a Result's Env map to a .env string.
func Marshal(r *Result, opts Options) (string, error) {
	if r == nil {
		return "", fmt.Errorf("pipeline: result must not be nil")
	}
	return export.Marshal(r.Env, export.Options{
		Header:     opts.Header,
		RedactKeys: opts.RedactKeys,
	})
}

// FormatDiff returns a human-readable diff string for the result's changes.
func FormatDiff(r *Result, maskSecrets bool) string {
	if r == nil {
		return ""
	}
	var sb strings.Builder
	for _, line := range diff.Format(r.Diff, maskSecrets) {
		sb.WriteString(line)
		sb.WriteByte('\n')
	}
	_ = io.Discard // suppress unused import
	return sb.String()
}
