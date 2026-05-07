package interpolate_test

import (
	"testing"

	"github.com/yourorg/envpatch/internal/interpolate"
)

func TestResolve_SubstitutesBaseToken(t *testing.T) {
	base := map[string]string{"HOST": "db.internal"}
	overlay := map[string]string{"DATABASE_URL": "postgres://${BASE:HOST}/mydb"}

	res, err := interpolate.Resolve(base, overlay)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := res.Env["DATABASE_URL"]
	want := "postgres://db.internal/mydb"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
	if len(res.Warnings) != 0 {
		t.Errorf("expected no warnings, got %v", res.Warnings)
	}
}

func TestResolve_UnresolvedTokenProducesWarning(t *testing.T) {
	base := map[string]string{}
	overlay := map[string]string{"REDIS_URL": "redis://${BASE:REDIS_HOST}:6379"}

	res, err := interpolate.Resolve(base, overlay)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Env["REDIS_URL"] != "redis://${BASE:REDIS_HOST}:6379" {
		t.Errorf("expected token to remain intact, got %q", res.Env["REDIS_URL"])
	}
	if len(res.Warnings) != 1 {
		t.Fatalf("expected 1 warning, got %d", len(res.Warnings))
	}
}

func TestResolve_NoTokensPassThrough(t *testing.T) {
	base := map[string]string{"X": "1"}
	overlay := map[string]string{"PLAIN": "hello"}

	res, err := interpolate.Resolve(base, overlay)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Env["PLAIN"] != "hello" {
		t.Errorf("expected 'hello', got %q", res.Env["PLAIN"])
	}
	if len(res.Warnings) != 0 {
		t.Errorf("expected no warnings")
	}
}

func TestResolve_MultipleTokensInSingleValue(t *testing.T) {
	base := map[string]string{"USER": "admin", "PASS": "s3cr3t"}
	overlay := map[string]string{"DSN": "${BASE:USER}:${BASE:PASS}@host"}

	res, err := interpolate.Resolve(base, overlay)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := "admin:s3cr3t@host"
	if res.Env["DSN"] != want {
		t.Errorf("got %q, want %q", res.Env["DSN"], want)
	}
}

func TestResolve_NilBaseReturnsError(t *testing.T) {
	_, err := interpolate.Resolve(nil, map[string]string{})
	if err == nil {
		t.Error("expected error for nil base")
	}
}

func TestResolve_NilOverlayReturnsError(t *testing.T) {
	_, err := interpolate.Resolve(map[string]string{}, nil)
	if err == nil {
		t.Error("expected error for nil overlay")
	}
}

func TestResolve_DoesNotMutateInputs(t *testing.T) {
	base := map[string]string{"KEY": "val"}
	overlay := map[string]string{"A": "${BASE:KEY}"}

	_, _ = interpolate.Resolve(base, overlay)

	if overlay["A"] != "${BASE:KEY}" {
		t.Error("overlay was mutated")
	}
}
