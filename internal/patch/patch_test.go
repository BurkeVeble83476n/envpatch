package patch_test

import (
	"testing"

	"github.com/yourorg/envpatch/internal/patch"
)

func TestApply_SetNewKey(t *testing.T) {
	base := map[string]string{"FOO": "bar"}
	p := patch.Patch{Entries: []patch.Entry{
		{Op: patch.OpSet, Key: "BAZ", Value: "qux"},
	}}
	out, err := patch.Apply(base, p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["BAZ"] != "qux" {
		t.Errorf("expected BAZ=qux, got %q", out["BAZ"])
	}
	if out["FOO"] != "bar" {
		t.Errorf("expected FOO=bar, got %q", out["FOO"])
	}
}

func TestApply_UpdateExistingKey(t *testing.T) {
	base := map[string]string{"FOO": "old"}
	p := patch.Patch{Entries: []patch.Entry{
		{Op: patch.OpSet, Key: "FOO", Value: "new"},
	}}
	out, err := patch.Apply(base, p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["FOO"] != "new" {
		t.Errorf("expected FOO=new, got %q", out["FOO"])
	}
}

func TestApply_RemoveKey(t *testing.T) {
	base := map[string]string{"FOO": "bar", "DEL": "me"}
	p := patch.Patch{Entries: []patch.Entry{
		{Op: patch.OpRemove, Key: "DEL"},
	}}
	out, err := patch.Apply(base, p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := out["DEL"]; ok {
		t.Error("expected DEL to be removed")
	}
}

func TestApply_RemoveMissingKeyIsNoop(t *testing.T) {
	base := map[string]string{"FOO": "bar"}
	p := patch.Patch{Entries: []patch.Entry{
		{Op: patch.OpRemove, Key: "MISSING"},
	}}
	out, err := patch.Apply(base, p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 1 {
		t.Errorf("expected 1 key, got %d", len(out))
	}
}

func TestApply_NilBaseReturnsError(t *testing.T) {
	_, err := patch.Apply(nil, patch.Patch{})
	if err == nil {
		t.Error("expected error for nil base")
	}
}

func TestApply_EmptyKeyReturnsError(t *testing.T) {
	base := map[string]string{}
	p := patch.Patch{Entries: []patch.Entry{
		{Op: patch.OpSet, Key: "", Value: "v"},
	}}
	_, err := patch.Apply(base, p)
	if err == nil {
		t.Error("expected error for empty key")
	}
}

func TestApply_UnknownOpReturnsError(t *testing.T) {
	base := map[string]string{}
	p := patch.Patch{Entries: []patch.Entry{
		{Op: patch.Op("invalid"), Key: "FOO"},
	}}
	_, err := patch.Apply(base, p)
	if err == nil {
		t.Error("expected error for unknown op")
	}
}

func TestApply_DoesNotMutateBase(t *testing.T) {
	base := map[string]string{"FOO": "original"}
	p := patch.Patch{Entries: []patch.Entry{
		{Op: patch.OpSet, Key: "FOO", Value: "mutated"},
	}}
	_, _ = patch.Apply(base, p)
	if base["FOO"] != "original" {
		t.Error("Apply must not mutate the base map")
	}
}
