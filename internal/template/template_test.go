package template

import (
	"errors"
	"testing"
)

func TestExpand_BracePlaceholder(t *testing.T) {
	env := map[string]string{"GREETING": "Hello, {{NAME}}!"}
	ctx := map[string]string{"NAME": "World"}
	out, errs := Expand(env, ctx)
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if got := out["GREETING"]; got != "Hello, World!" {
		t.Errorf("got %q, want %q", got, "Hello, World!")
	}
}

func TestExpand_ShellPlaceholder(t *testing.T) {
	env := map[string]string{"DSN": "postgres://${DB_HOST}:5432/app"}
	ctx := map[string]string{"DB_HOST": "localhost"}
	out, errs := Expand(env, ctx)
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if got := out["DSN"]; got != "postgres://localhost:5432/app" {
		t.Errorf("got %q, want %q", got, "postgres://localhost:5432/app")
	}
}

func TestExpand_UnresolvedPlaceholderReturnsError(t *testing.T) {
	env := map[string]string{"URL": "https://{{HOST}}/path"}
	_, errs := Expand(env, map[string]string{})
	if len(errs) == 0 {
		t.Fatal("expected an error for unresolved placeholder")
	}
	var unres *ErrUnresolved
	if !errors.As(errs[0], &unres) {
		t.Fatalf("expected *ErrUnresolved, got %T", errs[0])
	}
	if unres.Key != "URL" {
		t.Errorf("expected key URL, got %q", unres.Key)
	}
}

func TestExpand_NoPlaceholdersPassThrough(t *testing.T) {
	env := map[string]string{"PLAIN": "just-a-value"}
	out, errs := Expand(env, nil)
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if got := out["PLAIN"]; got != "just-a-value" {
		t.Errorf("got %q, want %q", got, "just-a-value")
	}
}

func TestExpand_MixedPatterns(t *testing.T) {
	env := map[string]string{"ADDR": "{{SCHEME}}://${HOST}"}
	ctx := map[string]string{"SCHEME": "https", "HOST": "example.com"}
	out, errs := Expand(env, ctx)
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if got := out["ADDR"]; got != "https://example.com" {
		t.Errorf("got %q, want %q", got, "https://example.com")
	}
}

func TestExpand_NilEnvReturnsError(t *testing.T) {
	_, errs := Expand(nil, nil)
	if len(errs) == 0 {
		t.Fatal("expected error for nil env")
	}
}
