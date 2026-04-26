package envfile

import (
	"strings"
	"testing"
)

func TestParse_BasicKeyValue(t *testing.T) {
	input := `APP_ENV=production
DATABASE_URL=postgres://localhost/mydb
`
	ef, err := Parse(strings.NewReader(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(ef.Index) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(ef.Index))
	}
	assertEntry(t, ef, "APP_ENV", "production")
	assertEntry(t, ef, "DATABASE_URL", "postgres://localhost/mydb")
}

func TestParse_QuotedValues(t *testing.T) {
	input := `SECRET="hello world"
TOKEN='abc123'
`
	ef, err := Parse(strings.NewReader(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	assertEntry(t, ef, "SECRET", "hello world")
	assertEntry(t, ef, "TOKEN", "abc123")
}

func TestParse_InlineComment(t *testing.T) {
	input := `PORT=8080 # default port
`
	ef, err := Parse(strings.NewReader(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	assertEntry(t, ef, "PORT", "8080")
	idx := ef.Index["PORT"]
	if ef.Entries[idx].Comment != "# default port" {
		t.Errorf("expected inline comment, got %q", ef.Entries[idx].Comment)
	}
}

func TestParse_SkipsBlankAndCommentLines(t *testing.T) {
	input := `# this is a comment

FOO=bar
`
	ef, err := Parse(strings.NewReader(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(ef.Index) != 1 {
		t.Fatalf("expected 1 key, got %d", len(ef.Index))
	}
}

func TestParse_MissingEquals(t *testing.T) {
	input := `BADLINE
`
	_, err := Parse(strings.NewReader(input))
	if err == nil {
		t.Fatal("expected error for missing '=', got nil")
	}
}

func assertEntry(t *testing.T, ef *EnvFile, key, wantValue string) {
	t.Helper()
	idx, ok := ef.Index[key]
	if !ok {
		t.Errorf("key %q not found", key)
		return
	}
	got := ef.Entries[idx].Value
	if got != wantValue {
		t.Errorf("key %q: expected value %q, got %q", key, wantValue, got)
	}
}
