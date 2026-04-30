package snapshot_test

import {
	"bytes"
	"testing"

	"github.com/yourorg/envpatch/internal/snapshot"
)

func TestNew_CopiesMap(t *testing.T) {
	env := map[string]string{"FOO": "bar", "BAZ": "qux"}
	s, err := snapshot.New(env, "test")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	env["FOO"] = "mutated"
	if s.Values["FOO"] != "bar" {
		t.Errorf("snapshot should be independent of source map")
	}
}

func TestNew_NilMapReturnsError(t *testing.T) {
	_, err := snapshot.New(nil, "label")
	if err == nil {
		t.Fatal("expected error for nil map")
	}
}

func TestSaveAndLoad_RoundTrip(t *testing.T) {
	env := map[string]string{"KEY": "value", "OTHER": "123"}
	s, _ := snapshot.New(env, "prod-2024")

	var buf bytes.Buffer
	if err := snapshot.Save(s, &buf); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	loaded, err := snapshot.Load(&buf)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	if loaded.Label != "prod-2024" {
		t.Errorf("expected label prod-2024, got %s", loaded.Label)
	}
	if loaded.Values["KEY"] != "value" {
		t.Errorf("expected KEY=value, got %s", loaded.Values["KEY"])
	}
}

func TestLoad_NilReaderReturnsError(t *testing.T) {
	_, err := snapshot.Load(nil)
	if err == nil {
		t.Fatal("expected error for nil reader")
	}
}

func TestCompare_DetectsDrift(t *testing.T) {
	base := map[string]string{"A": "1", "B": "2", "C": "3"}
	s, _ := snapshot.New(base, "snap")

	current := map[string]string{
		"A": "1",   // unchanged
		"B": "999", // changed
		// C removed
		"D": "new", // added
	}

	drifts, err := snapshot.Compare(s, current)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	kinds := map[string]string{}
	for _, d := range drifts {
		kinds[d.Key] = d.Kind
	}
	if kinds["B"] != "changed" {
		t.Errorf("expected B to be changed")
	}
	if kinds["C"] != "removed" {
		t.Errorf("expected C to be removed")
	}
	if kinds["D"] != "added" {
		t.Errorf("expected D to be added")
	}
	if _, ok := kinds["A"]; ok {
		t.Errorf("A should not appear in drifts")
	}
}

func TestCompare_NilSnapshotReturnsError(t *testing.T) {
	_, err := snapshot.Compare(nil, map[string]string{})
	if err == nil {
		t.Fatal("expected error for nil snapshot")
	}
}

func TestCompare_NilCurrentReturnsError(t *testing.T) {
	s, _ := snapshot.New(map[string]string{}, "s")
	_, err := snapshot.Compare(s, nil)
	if err == nil {
		t.Fatal("expected error for nil current map")
	}
}
