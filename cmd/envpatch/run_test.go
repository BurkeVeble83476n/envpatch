package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourorg/envpatch/internal/merge"
)

func writeTmp(t *testing.T, name, content string) string {
	t.Helper()
	p := filepath.Join(t.TempDir(), name)
	if err := os.WriteFile(p, []byte(content), 0o600); err != nil {
		t.Fatalf("writeTmp: %v", err)
	}
	return p
}

func TestRunPipeline_BasicMerge(t *testing.T) {
	base := writeTmp(t, ".env.base", "APP_ENV=development\nPORT=3000\n")
	overlay := writeTmp(t, ".env.overlay", "APP_ENV=staging\n")
	out := filepath.Join(t.TempDir(), ".env.out")

	err := runPipeline(runConfig{
		BaseFile:    base,
		OverlayFile: overlay,
		Strategy:    merge.Overlay,
		OutputFile:  out,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, _ := os.ReadFile(out)
	content := string(data)
	if content == "" {
		t.Error("expected non-empty output file")
	}
}

func TestRunPipeline_MissingBaseFile(t *testing.T) {
	err := runPipeline(runConfig{
		BaseFile: "/nonexistent/.env",
		Strategy: merge.Overlay,
	})
	if err == nil {
		t.Error("expected error for missing base file")
	}
}

func TestRunPipeline_WithSchema(t *testing.T) {
	base := writeTmp(t, ".env", "APP_ENV=production\nDB_URL=postgres://localhost/db\n")
	schemaFile := writeTmp(t, ".env.schema", "APP_ENV required\nDB_URL required\nLOG_LEVEL optional\n")

	err := runPipeline(runConfig{
		BaseFile:   base,
		SchemaFile: schemaFile,
		Strategy:   merge.Overlay,
		OutputFile: "-",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRunPipeline_PrintDiff(t *testing.T) {
	base := writeTmp(t, ".env.base", "APP_ENV=development\n")
	overlay := writeTmp(t, ".env.overlay", "APP_ENV=production\n")

	err := runPipeline(runConfig{
		BaseFile:    base,
		OverlayFile: overlay,
		Strategy:    merge.Overlay,
		PrintDiff:   true,
		OutputFile:  "-",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
