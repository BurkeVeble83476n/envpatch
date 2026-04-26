package merge

import (
	"testing"
)

func TestMerge_OverlayStrategy(t *testing.T) {
	base := map[string]string{"HOST": "localhost", "PORT": "5432"}
	overlay := map[string]string{"PORT": "9999", "DEBUG": "true"}

	result, err := Merge(base, overlay, StrategyOverlay)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Merged["PORT"] != "9999" {
		t.Errorf("expected PORT=9999, got %s", result.Merged["PORT"])
	}
	if result.Merged["HOST"] != "localhost" {
		t.Errorf("expected HOST=localhost, got %s", result.Merged["HOST"])
	}
	if result.Merged["DEBUG"] != "true" {
		t.Errorf("expected DEBUG=true, got %s", result.Merged["DEBUG"])
	}
	if len(result.Overridden) != 1 || result.Overridden[0] != "PORT" {
		t.Errorf("expected Overridden=[PORT], got %v", result.Overridden)
	}
	if len(result.Added) != 1 || result.Added[0] != "DEBUG" {
		t.Errorf("expected Added=[DEBUG], got %v", result.Added)
	}
}

func TestMerge_KeepBaseStrategy(t *testing.T) {
	base := map[string]string{"HOST": "localhost", "PORT": "5432"}
	overlay := map[string]string{"PORT": "9999", "DEBUG": "true"}

	result, err := Merge(base, overlay, StrategyKeepBase)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Merged["PORT"] != "5432" {
		t.Errorf("expected PORT=5432 (base kept), got %s", result.Merged["PORT"])
	}
	if len(result.Overridden) != 0 {
		t.Errorf("expected no overrides with KeepBase strategy, got %v", result.Overridden)
	}
	if len(result.Added) != 1 || result.Added[0] != "DEBUG" {
		t.Errorf("expected Added=[DEBUG], got %v", result.Added)
	}
}

func TestMerge_EmptyOverlay(t *testing.T) {
	base := map[string]string{"A": "1", "B": "2"}
	result, err := Merge(base, map[string]string{}, StrategyOverlay)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Merged) != 2 {
		t.Errorf("expected 2 keys, got %d", len(result.Merged))
	}
	if len(result.Added) != 0 || len(result.Overridden) != 0 {
		t.Errorf("expected no changes with empty overlay")
	}
}

func TestMerge_NilBaseReturnsError(t *testing.T) {
	_, err := Merge(nil, map[string]string{}, StrategyOverlay)
	if err == nil {
		t.Error("expected error for nil base, got nil")
	}
}

func TestMerge_NilOverlayReturnsError(t *testing.T) {
	_, err := Merge(map[string]string{}, nil, StrategyOverlay)
	if err == nil {
		t.Error("expected error for nil overlay, got nil")
	}
}
