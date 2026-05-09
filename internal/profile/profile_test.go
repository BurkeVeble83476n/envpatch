package profile_test

import (
	"bytes"
	"testing"

	"github.com/yourorg/envpatch/internal/profile"
)

func TestSet_AddsProfile(t *testing.T) {
	s := profile.NewStore()
	env := map[string]string{"FOO": "bar"}
	if err := s.Set("dev", env); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	p, err := s.Get("dev")
	if err != nil {
		t.Fatalf("expected profile: %v", err)
	}
	if p.Env["FOO"] != "bar" {
		t.Errorf("expected FOO=bar, got %q", p.Env["FOO"])
	}
}

func TestSet_EmptyNameReturnsError(t *testing.T) {
	s := profile.NewStore()
	if err := s.Set("", map[string]string{}); err == nil {
		t.Fatal("expected error for empty name")
	}
}

func TestSet_NilEnvReturnsError(t *testing.T) {
	s := profile.NewStore()
	if err := s.Set("dev", nil); err == nil {
		t.Fatal("expected error for nil env")
	}
}

func TestGet_UnknownProfileReturnsError(t *testing.T) {
	s := profile.NewStore()
	if _, err := s.Get("missing"); err == nil {
		t.Fatal("expected error for missing profile")
	}
}

func TestDelete_RemovesProfile(t *testing.T) {
	s := profile.NewStore()
	_ = s.Set("prod", map[string]string{"X": "1"})
	if err := s.Delete("prod"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := s.Get("prod"); err == nil {
		t.Fatal("expected profile to be deleted")
	}
}

func TestDelete_MissingProfileReturnsError(t *testing.T) {
	s := profile.NewStore()
	if err := s.Delete("ghost"); err == nil {
		t.Fatal("expected error")
	}
}

func TestList_SortedNames(t *testing.T) {
	s := profile.NewStore()
	for _, n := range []string{"staging", "dev", "prod"} {
		_ = s.Set(n, map[string]string{})
	}
	names := s.List()
	want := []string{"dev", "prod", "staging"}
	for i, n := range names {
		if n != want[i] {
			t.Errorf("position %d: got %q want %q", i, n, want[i])
		}
	}
}

func TestSaveAndLoad_RoundTrip(t *testing.T) {
	s := profile.NewStore()
	_ = s.Set("ci", map[string]string{"CI": "true", "PORT": "8080"})
	var buf bytes.Buffer
	if err := s.Save(&buf); err != nil {
		t.Fatalf("save: %v", err)
	}
	s2, err := profile.LoadStore(&buf)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	p, err := s2.Get("ci")
	if err != nil {
		t.Fatalf("get after load: %v", err)
	}
	if p.Env["PORT"] != "8080" {
		t.Errorf("expected PORT=8080, got %q", p.Env["PORT"])
	}
}

func TestLoadStore_NilReaderReturnsError(t *testing.T) {
	if _, err := profile.LoadStore(nil); err == nil {
		t.Fatal("expected error for nil reader")
	}
}
