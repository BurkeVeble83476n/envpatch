package history_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/yourorg/envpatch/internal/history"
)

func TestNew_EmptyLog(t *testing.T) {
	l := history.New()
	if l.Len() != 0 {
		t.Fatalf("expected 0 entries, got %d", l.Len())
	}
}

func TestRecord_AppendsEntry(t *testing.T) {
	l := history.New()
	l.Record("merge", "merged dev into staging", map[string]string{"DB_HOST": "localhost"})
	if l.Len() != 1 {
		t.Fatalf("expected 1 entry, got %d", l.Len())
	}
	entry := l.Entries[0]
	if entry.Operation != "merge" {
		t.Errorf("expected operation 'merge', got %q", entry.Operation)
	}
	if entry.Message != "merged dev into staging" {
		t.Errorf("unexpected message: %q", entry.Message)
	}
	if entry.Changes["DB_HOST"] != "localhost" {
		t.Errorf("expected change DB_HOST=localhost")
	}
	if entry.Timestamp.IsZero() {
		t.Error("timestamp should not be zero")
	}
}

func TestSaveAndLoad_RoundTrip(t *testing.T) {
	l := history.New()
	l.Record("patch", "applied hotfix", nil)
	l.Record("validate", "schema check passed", nil)

	var buf bytes.Buffer
	if err := l.Save(&buf); err != nil {
		t.Fatalf("Save: %v", err)
	}

	loaded, err := history.Load(&buf)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if loaded.Len() != 2 {
		t.Fatalf("expected 2 entries after round-trip, got %d", loaded.Len())
	}
	if loaded.Entries[0].Operation != "patch" {
		t.Errorf("expected first op 'patch', got %q", loaded.Entries[0].Operation)
	}
}

func TestSave_NilWriterReturnsError(t *testing.T) {
	l := history.New()
	if err := l.Save(nil); err == nil {
		t.Fatal("expected error for nil writer")
	}
}

func TestLoad_NilReaderReturnsError(t *testing.T) {
	if _, err := history.Load(nil); err == nil {
		t.Fatal("expected error for nil reader")
	}
}

func TestLoad_InvalidJSONReturnsError(t *testing.T) {
	r := strings.NewReader("{not valid json}")
	if _, err := history.Load(r); err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestRecord_MultipleEntries(t *testing.T) {
	l := history.New()
	ops := []string{"merge", "patch", "validate", "export"}
	for _, op := range ops {
		l.Record(op, "", nil)
	}
	if l.Len() != len(ops) {
		t.Fatalf("expected %d entries, got %d", len(ops), l.Len())
	}
}
