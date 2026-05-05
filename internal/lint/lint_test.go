package lint

import (
	"strings"
	"testing"
)

func TestCheck_EmptyMapReturnsNil(t *testing.T) {
	got := Check(nil)
	if got != nil {
		t.Errorf("expected nil findings for nil map, got %v", got)
	}

	got = Check(map[string]string{})
	if got != nil {
		t.Errorf("expected nil findings for empty map, got %v", got)
	}
}

func TestCheck_LowercaseKeyWarning(t *testing.T) {
	env := map[string]string{"db_host": "localhost"}
	findings := Check(env)
	if !hasMessage(findings, "lowercase") {
		t.Errorf("expected lowercase warning, got %v", findings)
	}
}

func TestCheck_UppercaseKeyClean(t *testing.T) {
	env := map[string]string{"DB_HOST": "localhost"}
	findings := Check(env)
	if hasMessage(findings, "lowercase") {
		t.Errorf("unexpected lowercase warning for all-caps key")
	}
}

func TestCheck_UnresolvedBraceRef(t *testing.T) {
	env := map[string]string{"API_URL": "https://${HOST}/api"}
	findings := Check(env)
	if !hasMessage(findings, "unresolved variable") {
		t.Errorf("expected unresolved variable warning, got %v", findings)
	}
}

func TestCheck_UnresolvedBareRef(t *testing.T) {
	env := map[string]string{"DSN": "postgres://$USER:$PASS@host/db"}
	findings := Check(env)
	if !hasMessage(findings, "unresolved variable") {
		t.Errorf("expected unresolved variable warning for bare $VAR, got %v", findings)
	}
}

func TestCheck_LongValueWarning(t *testing.T) {
	env := map[string]string{"CERT": strings.Repeat("x", 501)}
	findings := Check(env)
	if !hasMessage(findings, "unusually long") {
		t.Errorf("expected long-value warning, got %v", findings)
	}
}

func TestCheck_AcceptableValueLength(t *testing.T) {
	env := map[string]string{"TOKEN": strings.Repeat("a", 100)}
	findings := Check(env)
	if hasMessage(findings, "unusually long") {
		t.Errorf("unexpected long-value warning for short value")
	}
}

func TestCheck_MultipleIssuesSameKey(t *testing.T) {
	env := map[string]string{"my_key": "value with $REF"}
	findings := Check(env)
	// Expect both lowercase and unresolved-ref warnings.
	if !hasMessage(findings, "lowercase") || !hasMessage(findings, "unresolved variable") {
		t.Errorf("expected both warnings, got %v", findings)
	}
}

func TestFinding_String(t *testing.T) {
	f := Finding{Key: "FOO", Message: "some issue", Severity: Warning}
	s := f.String()
	if !strings.Contains(s, "warning") || !strings.Contains(s, "FOO") {
		t.Errorf("unexpected String() output: %s", s)
	}
}

// hasMessage returns true if any finding's message contains substr.
func hasMessage(findings []Finding, substr string) bool {
	for _, f := range findings {
		if strings.Contains(f.Message, substr) {
			return true
		}
	}
	return false
}
