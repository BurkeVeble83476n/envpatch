package main

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/yourorg/envpatch/internal/envfile"
	"github.com/yourorg/envpatch/internal/profile"
)

// profileStoreFile is the default location for the profile store.
const profileStoreFile = ".envpatch_profiles.json"

// loadOrNewStore opens profileStoreFile if it exists, otherwise returns a new
// empty Store.
func loadOrNewStore() (*profile.Store, error) {
	f, err := os.Open(profileStoreFile)
	if errors.Is(err, os.ErrNotExist) {
		return profile.NewStore(), nil
	}
	if err != nil {
		return nil, fmt.Errorf("open profile store: %w", err)
	}
	defer f.Close()
	return profile.LoadStore(f)
}

// saveStore persists s to profileStoreFile.
func saveStore(s *profile.Store) error {
	f, err := os.Create(profileStoreFile)
	if err != nil {
		return fmt.Errorf("create profile store: %w", err)
	}
	defer f.Close()
	return s.Save(f)
}

// runProfileSet reads envFile, stores it as name in the profile store.
func runProfileSet(name, envFile string) error {
	if name == "" {
		return errors.New("profile name is required")
	}
	data, err := os.ReadFile(envFile)
	if err != nil {
		return fmt.Errorf("read env file: %w", err)
	}
	env, err := envfile.Parse(strings.NewReader(string(data)))
	if err != nil {
		return fmt.Errorf("parse env file: %w", err)
	}
	store, err := loadOrNewStore()
	if err != nil {
		return err
	}
	if err := store.Set(name, env); err != nil {
		return err
	}
	return saveStore(store)
}

// runProfileGet prints the env vars for name to stdout.
func runProfileGet(name string) error {
	store, err := loadOrNewStore()
	if err != nil {
		return err
	}
	p, err := store.Get(name)
	if err != nil {
		return err
	}
	for _, k := range sortedKeys(p.Env) {
		fmt.Printf("%s=%s\n", k, p.Env[k])
	}
	return nil
}

// runProfileList prints all profile names to stdout.
func runProfileList() error {
	store, err := loadOrNewStore()
	if err != nil {
		return err
	}
	for _, n := range store.List() {
		fmt.Println(n)
	}
	return nil
}

// sortedKeys returns map keys in sorted order (local helper).
func sortedKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	// reuse sort from standard library via import in profile.go
	importSort := func(s []string) {
		for i := 1; i < len(s); i++ {
			for j := i; j > 0 && s[j] < s[j-1]; j-- {
				s[j], s[j-1] = s[j-1], s[j]
			}
		}
	}
	importSort(keys)
	return keys
}
