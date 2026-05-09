// Package profile manages named environment profiles, allowing users to
// store and retrieve environment variable sets identified by a profile name
// (e.g. "staging", "production").
package profile

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"sort"
	"time"
)

// Profile holds a named snapshot of environment variables.
type Profile struct {
	Name      string            `json:"name"`
	CreatedAt time.Time         `json:"created_at"`
	UpdatedAt time.Time         `json:"updated_at"`
	Env       map[string]string `json:"env"`
}

// Store is a collection of named profiles persisted as JSON.
type Store struct {
	Profiles map[string]*Profile `json:"profiles"`
}

// NewStore returns an empty Store.
func NewStore() *Store {
	return &Store{Profiles: make(map[string]*Profile)}
}

// Set creates or replaces the profile with the given name.
func (s *Store) Set(name string, env map[string]string) error {
	if name == "" {
		return errors.New("profile name must not be empty")
	}
	if env == nil {
		return errors.New("env map must not be nil")
	}
	now := time.Now().UTC()
	p, exists := s.Profiles[name]
	if !exists {
		p = &Profile{Name: name, CreatedAt: now}
	}
	copy := make(map[string]string, len(env))
	for k, v := range env {
		copy[k] = v
	}
	p.Env = copy
	p.UpdatedAt = now
	s.Profiles[name] = p
	return nil
}

// Get retrieves a profile by name. Returns an error if not found.
func (s *Store) Get(name string) (*Profile, error) {
	p, ok := s.Profiles[name]
	if !ok {
		return nil, fmt.Errorf("profile %q not found", name)
	}
	return p, nil
}

// Delete removes a profile. Returns an error if not found.
func (s *Store) Delete(name string) error {
	if _, ok := s.Profiles[name]; !ok {
		return fmt.Errorf("profile %q not found", name)
	}
	delete(s.Profiles, name)
	return nil
}

// List returns profile names in sorted order.
func (s *Store) List() []string {
	names := make([]string, 0, len(s.Profiles))
	for n := range s.Profiles {
		names = append(names, n)
	}
	sort.Strings(names)
	return names
}

// Save serialises the Store as JSON to w.
func (s *Store) Save(w io.Writer) error {
	if w == nil {
		return errors.New("writer must not be nil")
	}
	return json.NewEncoder(w).Encode(s)
}

// LoadStore deserialises a Store from r.
func LoadStore(r io.Reader) (*Store, error) {
	if r == nil {
		return nil, errors.New("reader must not be nil")
	}
	var s Store
	if err := json.NewDecoder(r).Decode(&s); err != nil {
		return nil, fmt.Errorf("decode profile store: %w", err)
	}
	if s.Profiles == nil {
		s.Profiles = make(map[string]*Profile)
	}
	return &s, nil
}
