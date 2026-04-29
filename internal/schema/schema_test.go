package schema_test

import (
	"strings"
	"testing"

	"github.com/yourorg/envpatch/internal/schema"
)

func TestLoad_RequiredAndOptionalKeys(t *testing.T) {
	input := `
DB_HOST required # database hostname
DB_PORT optional # database port
APP_SECRET required
`
	s, err := schema.Load(strings.NewReader(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(s) != 3 {
		t.Fatalf("expected 3 keys, got %d", len(s))
	}

	if !s["DB_HOST"].Required {
		t.Error("DB_HOST should be required")
	}
	if s["DB_PORT"].Required {
		t.Error("DB_PORT should be optional")
	}
	if s["DB_HOST"].Description != "database hostname" {
		t.Errorf("unexpected description: %q", s["DB_HOST"].Description)
	}
}

func TestLoad_DefaultsToRequired(t *testing.T) {
	input := "SOME_KEY\n"
	s, err := schema.Load(strings.NewReader(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !s["SOME_KEY"].Required {
		t.Error("key with no qualifier should default to required")
	}
}

func TestLoad_SkipsBlankAndCommentLines(t *testing.T) {
	input := "# this is a comment\n\nVALID_KEY required\n"
	s, err := schema.Load(strings.NewReader(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(s) != 1 {
		t.Fatalf("expected 1 key, got %d", len(s))
	}
}

func TestLoad_UnknownQualifierReturnsError(t *testing.T) {
	input := "MY_KEY bogus\n"
	_, err := schema.Load(strings.NewReader(input))
	if err == nil {
		t.Fatal("expected error for unknown qualifier, got nil")
	}
}

func TestLoad_NilReaderReturnsError(t *testing.T) {
	_, err := schema.Load(nil)
	if err == nil {
		t.Fatal("expected error for nil reader, got nil")
	}
}

func TestLoad_KeySpecFieldsPopulated(t *testing.T) {
	input := "API_KEY required # secret api key\n"
	s, err := schema.Load(strings.NewReader(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	spec := s["API_KEY"]
	if spec.Key != "API_KEY" {
		t.Errorf("expected Key=API_KEY, got %q", spec.Key)
	}
	if spec.Description != "secret api key" {
		t.Errorf("expected description 'secret api key', got %q", spec.Description)
	}
}
