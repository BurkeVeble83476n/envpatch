package validate_test

import (
	"testing"

	"github.com/yourorg/envpatch/internal/validate"
)

func TestValidate_AllKeysPresent(t *testing.T) {
	schema := map[string]string{"APP_ENV": "production", "PORT": "8080"}
	target := map[string]string{"APP_ENV": "staging", "PORT": "3000"}

	result, err := validate.Validate(schema, target)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Issues) != 0 {
		t.Errorf("expected no issues, got %d", len(result.Issues))
	}
}

func TestValidate_MissingRequiredKey(t *testing.T) {
	schema := map[string]string{"APP_ENV": "", "PORT": "", "DB_URL": ""}
	target := map[string]string{"APP_ENV": "staging"}

	result, err := validate.Validate(schema, target)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.HasErrors() {
		t.Error("expected errors but got none")
	}

	errorKeys := map[string]bool{}
	for _, issue := range result.Issues {
		if issue.Severity == validate.SeverityError {
			errorKeys[issue.Key] = true
		}
	}
	for _, k := range []string{"PORT", "DB_URL"} {
		if !errorKeys[k] {
			t.Errorf("expected error for key %q", k)
		}
	}
}

func TestValidate_ExtraKeyIsWarning(t *testing.T) {
	schema := map[string]string{"APP_ENV": ""}
	target := map[string]string{"APP_ENV": "prod", "EXTRA_KEY": "value"}

	result, err := validate.Validate(schema, target)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.HasErrors() {
		t.Error("expected no errors, only warnings")
	}
	if len(result.Issues) != 1 {
		t.Fatalf("expected 1 warning, got %d", len(result.Issues))
	}
	if result.Issues[0].Severity != validate.SeverityWarning {
		t.Errorf("expected warning severity, got %s", result.Issues[0].Severity)
	}
	if result.Issues[0].Key != "EXTRA_KEY" {
		t.Errorf("expected EXTRA_KEY warning, got %s", result.Issues[0].Key)
	}
}

func TestValidate_NilSchemaReturnsError(t *testing.T) {
	_, err := validate.Validate(nil, map[string]string{})
	if err == nil {
		t.Error("expected error for nil schema")
	}
}

func TestValidate_NilTargetReturnsError(t *testing.T) {
	_, err := validate.Validate(map[string]string{}, nil)
	if err == nil {
		t.Error("expected error for nil target")
	}
}
