package pipeline_test

import (
	"testing"

	"github.com/yourorg/envpatch/internal/merge"
	"github.com/yourorg/envpatch/internal/patch"
	"github.com/yourorg/envpatch/internal/pipeline"
	"github.com/yourorg/envpatch/internal/validate"
)

func TestRun_MergeAndPatch(t *testing.T) {
	base := map[string]string{"APP_ENV": "development", "PORT": "3000"}
	overlay := map[string]string{"APP_ENV": "staging"}
	ops := []patch.Op{{Action: patch.Set, Key: "LOG_LEVEL", Value: "debug"}}

	result, err := pipeline.Run(base, overlay, pipeline.Options{
		MergeStrategy: merge.Overlay,
		PatchOps:      ops,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Env["APP_ENV"] != "staging" {
		t.Errorf("expected APP_ENV=staging, got %q", result.Env["APP_ENV"])
	}
	if result.Env["LOG_LEVEL"] != "debug" {
		t.Errorf("expected LOG_LEVEL=debug, got %q", result.Env["LOG_LEVEL"])
	}
	if len(result.Diff) == 0 {
		t.Error("expected non-empty diff")
	}
}

func TestRun_ValidationIssues(t *testing.T) {
	base := map[string]string{"APP_ENV": "production"}
	schema := map[string]validate.Rule{
		"APP_ENV":    {Required: true},
		"DB_URL":     {Required: true},
		"LOG_LEVEL":  {Required: false},
	}

	result, err := pipeline.Run(base, nil, pipeline.Options{
		MergeStrategy: merge.Overlay,
		Schema:        schema,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Issues) == 0 {
		t.Error("expected validation issues for missing DB_URL")
	}
}

func TestRun_NilBaseReturnsError(t *testing.T) {
	_, err := pipeline.Run(nil, nil, pipeline.Options{})
	if err == nil {
		t.Error("expected error for nil base")
	}
}

func TestRunFromStrings_Basic(t *testing.T) {
	baseRaw := "APP_ENV=development\nPORT=3000\n"
	overlayRaw := "APP_ENV=production\n"

	result, err := pipeline.RunFromStrings(baseRaw, overlayRaw, pipeline.Options{
		MergeStrategy: merge.Overlay,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Env["APP_ENV"] != "production" {
		t.Errorf("expected APP_ENV=production, got %q", result.Env["APP_ENV"])
	}
}

func TestRunFromStrings_EmptyOverlay(t *testing.T) {
	baseRaw := "APP_ENV=development\n"

	result, err := pipeline.RunFromStrings(baseRaw, "", pipeline.Options{
		MergeStrategy: merge.Overlay,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Env["APP_ENV"] != "development" {
		t.Errorf("expected APP_ENV=development, got %q", result.Env["APP_ENV"])
	}
}

func TestMarshal_NilResultReturnsError(t *testing.T) {
	_, err := pipeline.Marshal(nil, "", false)
	if err == nil {
		t.Error("expected error for nil result")
	}
}
