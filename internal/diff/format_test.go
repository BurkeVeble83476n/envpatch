package diff_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/yourorg/envpatch/internal/diff"
)

func TestFormat_Added(t *testing.T) {
	result := &diff.Result{
		Changes: []diff.Change{
			{Key: "NEW", Type: diff.Added, NewValue: "val"},
		},
	}
	var buf bytes.Buffer
	diff.Format(&buf, result, diff.FormatOptions{})
	if !strings.Contains(buf.String(), "+ NEW=val") {
		t.Errorf("unexpected output: %q", buf.String())
	}
}

func TestFormat_Removed(t *testing.T) {
	result := &diff.Result{
		Changes: []diff.Change{
			{Key: "OLD", Type: diff.Removed, OldValue: "gone"},
		},
	}
	var buf bytes.Buffer
	diff.Format(&buf, result, diff.FormatOptions{})
	if !strings.Contains(buf.String(), "- OLD=gone") {
		t.Errorf("unexpected output: %q", buf.String())
	}
}

func TestFormat_Changed(t *testing.T) {
	result := &diff.Result{
		Changes: []diff.Change{
			{Key: "FOO", Type: diff.Changed, OldValue: "old", NewValue: "new"},
		},
	}
	var buf bytes.Buffer
	diff.Format(&buf, result, diff.FormatOptions{})
	out := buf.String()
	if !strings.Contains(out, "~ FOO") || !strings.Contains(out, "old") || !strings.Contains(out, "new") {
		t.Errorf("unexpected output: %q", out)
	}
}

func TestFormat_MaskSecrets(t *testing.T) {
	result := &diff.Result{
		Changes: []diff.Change{
			{Key: "DB_PASSWORD", Type: diff.Added, NewValue: "supersecret"},
		},
	}
	var buf bytes.Buffer
	diff.Format(&buf, result, diff.FormatOptions{
		MaskSecrets: true,
		SecretKeys:  []string{"DB_PASSWORD"},
	})
	out := buf.String()
	if strings.Contains(out, "supersecret") {
		t.Errorf("secret value should be masked, got: %q", out)
	}
	if !strings.Contains(out, "***") {
		t.Errorf("expected masked value in output: %q", out)
	}
}
