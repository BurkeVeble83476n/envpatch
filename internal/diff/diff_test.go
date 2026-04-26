package diff_test

import (
	"testing"

	"github.com/yourorg/envpatch/internal/diff"
)

func TestDiff_Added(t *testing.T) {
	base := map[string]string{"FOO": "bar"}
	target := map[string]string{"FOO": "bar", "NEW_KEY": "value"}

	result := diff.Diff(base, target)
	if len(result.Changes) != 1 {
		t.Fatalf("expected 1 change, got %d", len(result.Changes))
	}
	if result.Changes[0].Type != diff.Added || result.Changes[0].Key != "NEW_KEY" {
		t.Errorf("unexpected change: %+v", result.Changes[0])
	}
}

func TestDiff_Removed(t *testing.T) {
	base := map[string]string{"FOO": "bar", "OLD_KEY": "gone"}
	target := map[string]string{"FOO": "bar"}

	result := diff.Diff(base, target)
	if len(result.Changes) != 1 {
		t.Fatalf("expected 1 change, got %d", len(result.Changes))
	}
	if result.Changes[0].Type != diff.Removed || result.Changes[0].Key != "OLD_KEY" {
		t.Errorf("unexpected change: %+v", result.Changes[0])
	}
}

func TestDiff_Changed(t *testing.T) {
	base := map[string]string{"FOO": "old"}
	target := map[string]string{"FOO": "new"}

	result := diff.Diff(base, target)
	if len(result.Changes) != 1 {
		t.Fatalf("expected 1 change, got %d", len(result.Changes))
	}
	c := result.Changes[0]
	if c.Type != diff.Changed || c.OldValue != "old" || c.NewValue != "new" {
		t.Errorf("unexpected change: %+v", c)
	}
}

func TestDiff_NoChanges(t *testing.T) {
	env := map[string]string{"A": "1", "B": "2"}
	result := diff.Diff(env, env)
	if result.HasChanges() {
		t.Errorf("expected no changes, got %d", len(result.Changes))
	}
}

func TestDiff_SortedOutput(t *testing.T) {
	base := map[string]string{}
	target := map[string]string{"Z_KEY": "z", "A_KEY": "a", "M_KEY": "m"}

	result := diff.Diff(base, target)
	if len(result.Changes) != 3 {
		t.Fatalf("expected 3 changes, got %d", len(result.Changes))
	}
	if result.Changes[0].Key != "A_KEY" || result.Changes[1].Key != "M_KEY" || result.Changes[2].Key != "Z_KEY" {
		t.Errorf("changes not sorted: %+v", result.Changes)
	}
}
