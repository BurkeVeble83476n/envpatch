package audit_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/envpatch/internal/audit"
	"github.com/user/envpatch/internal/diff"
)

func TestNew_NilWriterReturnsError(t *testing.T) {
	_, err := audit.New(nil)
	if err == nil {
		t.Fatal("expected error for nil writer, got nil")
	}
}

func TestNew_ValidWriter(t *testing.T) {
	var buf bytes.Buffer
	l, err := audit.New(&buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if l == nil {
		t.Fatal("expected non-nil logger")
	}
}

func TestRecord_WritesKindAndMessage(t *testing.T) {
	var buf bytes.Buffer
	l, _ := audit.New(&buf)

	err := l.Record(audit.KindMerge, "merged staging into base", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "merge") {
		t.Errorf("expected 'merge' in output, got: %s", out)
	}
	if !strings.Contains(out, "merged staging into base") {
		t.Errorf("expected message in output, got: %s", out)
	}
	if !strings.Contains(out, "0 change(s)") {
		t.Errorf("expected change count in output, got: %s", out)
	}
}

func TestRecord_WritesChanges(t *testing.T) {
	var buf bytes.Buffer
	l, _ := audit.New(&buf)

	changes := []diff.Change{
		{Key: "DB_HOST", Op: diff.OpAdded},
		{Key: "API_KEY", Op: diff.OpChanged},
	}

	err := l.Record(audit.KindPatch, "applied patch", changes)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "2 change(s)") {
		t.Errorf("expected 2 changes in output, got: %s", out)
	}
	if !strings.Contains(out, "key=DB_HOST") {
		t.Errorf("expected DB_HOST in output, got: %s", out)
	}
	if !strings.Contains(out, "key=API_KEY") {
		t.Errorf("expected API_KEY in output, got: %s", out)
	}
}

func TestRecord_AllKinds(t *testing.T) {
	kinds := []audit.EntryKind{
		audit.KindMerge,
		audit.KindPatch,
		audit.KindValidate,
		audit.KindSnapshot,
	}
	for _, k := range kinds {
		var buf bytes.Buffer
		l, _ := audit.New(&buf)
		if err := l.Record(k, "test", nil); err != nil {
			t.Errorf("kind %s: unexpected error: %v", k, err)
		}
		if !strings.Contains(buf.String(), string(k)) {
			t.Errorf("kind %s not found in output: %s", k, buf.String())
		}
	}
}
