package export_test

import (
	"strings"
	"testing"

	"github.com/yourorg/envpatch/internal/export"
)

func TestMarshal_SortedKeys(t *testing.T) {
	env := map[string]string{"ZEBRA": "1", "ALPHA": "2", "MIDDLE": "3"}
	out, err := export.Marshal(env, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	s := string(out)
	alphaIdx := strings.Index(s, "ALPHA")
	middleIdx := strings.Index(s, "MIDDLE")
	zebraIdx := strings.Index(s, "ZEBRA")
	if !(alphaIdx < middleIdx && middleIdx < zebraIdx) {
		t.Errorf("keys not in sorted order:\n%s", s)
	}
}

func TestMarshal_QuotesValueWithSpaces(t *testing.T) {
	env := map[string]string{"MSG": "hello world"}
	out, err := export.Marshal(env, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(string(out), `MSG="hello world"`) {
		t.Errorf("expected quoted value, got: %s", out)
	}
}

func TestMarshal_QuotesEmptyValue(t *testing.T) {
	env := map[string]string{"EMPTY": ""}
	out, err := export.Marshal(env, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(string(out), `EMPTY=""`) {
		t.Errorf("expected empty quoted value, got: %s", out)
	}
}

func TestMarshal_RedactsSensitiveKeys(t *testing.T) {
	env := map[string]string{"DB_PASSWORD": "s3cr3t", "APP_NAME": "myapp"}
	out, err := export.Marshal(env, &export.Options{RedactSensitive: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	s := string(out)
	if strings.Contains(s, "s3cr3t") {
		t.Errorf("sensitive value was not redacted:\n%s", s)
	}
	if !strings.Contains(s, "APP_NAME=myapp") {
		t.Errorf("safe value should not be redacted:\n%s", s)
	}
}

func TestMarshal_HeaderPrependedAsComments(t *testing.T) {
	env := map[string]string{"KEY": "val"}
	out, err := export.Marshal(env, &export.Options{Header: "Generated file\nDo not edit"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	s := string(out)
	if !strings.HasPrefix(s, "# Generated file") {
		t.Errorf("expected header comment at start, got:\n%s", s)
	}
	if !strings.Contains(s, "# Do not edit") {
		t.Errorf("expected second header line, got:\n%s", s)
	}
}

func TestMarshal_NilEnvReturnsError(t *testing.T) {
	_, err := export.Marshal(nil, nil)
	if err == nil {
		t.Fatal("expected error for nil env map, got nil")
	}
}
