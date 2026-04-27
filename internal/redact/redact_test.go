package redact_test

import (
	"testing"

	"github.com/yourorg/envpatch/internal/redact"
)

func TestIsSensitive_MatchesDefaultPatterns(t *testing.T) {
	cases := []struct {
		key  string
		want bool
	}{
		{"DB_PASSWORD", true},
		{"AWS_SECRET_ACCESS_KEY", true},
		{"API_TOKEN", true},
		{"GITHUB_TOKEN", true},
		{"APP_NAME", false},
		{"PORT", false},
		{"DATABASE_URL", false},
	}
	for _, tc := range cases {
		got := redact.IsSensitive(tc.key, nil)
		if got != tc.want {
			t.Errorf("IsSensitive(%q) = %v, want %v", tc.key, got, tc.want)
		}
	}
}

func TestIsSensitive_CustomPatterns(t *testing.T) {
	patterns := []string{"MAGIC", "WAND"}
	if !redact.IsSensitive("MAGIC_VALUE", patterns) {
		t.Error("expected MAGIC_VALUE to be sensitive with custom patterns")
	}
	if redact.IsSensitive("DB_PASSWORD", patterns) {
		t.Error("expected DB_PASSWORD NOT to be sensitive with custom patterns")
	}
}

func TestValue_MasksSensitive(t *testing.T) {
	got := redact.Value("DB_PASSWORD", "supersecret", nil)
	if got != "***" {
		t.Errorf("expected masked value, got %q", got)
	}
}

func TestValue_PassesThroughSafe(t *testing.T) {
	got := redact.Value("APP_ENV", "production", nil)
	if got != "production" {
		t.Errorf("expected original value, got %q", got)
	}
}

func TestMap_RedactsCorrectKeys(t *testing.T) {
	env := map[string]string{
		"APP_NAME":    "myapp",
		"DB_PASSWORD": "hunter2",
		"API_TOKEN":   "tok_abc123",
		"PORT":        "8080",
	}
	result := redact.Map(env, nil)

	if result["APP_NAME"] != "myapp" {
		t.Errorf("APP_NAME should not be redacted, got %q", result["APP_NAME"])
	}
	if result["PORT"] != "8080" {
		t.Errorf("PORT should not be redacted, got %q", result["PORT"])
	}
	if result["DB_PASSWORD"] != "***" {
		t.Errorf("DB_PASSWORD should be redacted, got %q", result["DB_PASSWORD"])
	}
	if result["API_TOKEN"] != "***" {
		t.Errorf("API_TOKEN should be redacted, got %q", result["API_TOKEN"])
	}
}

func TestMap_DoesNotMutateOriginal(t *testing.T) {
	env := map[string]string{"SECRET_KEY": "real-value"}
	_ = redact.Map(env, nil)
	if env["SECRET_KEY"] != "real-value" {
		t.Error("original map was mutated")
	}
}
