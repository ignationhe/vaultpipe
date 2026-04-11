package envfile_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourorg/vaultpipe/internal/envfile"
)

// TestRenderMap_ThenWrite verifies the full render→write pipeline.
func TestRenderMap_ThenWrite(t *testing.T) {
	dir := t.TempDir()
	dst := filepath.Join(dir, "out.env")

	vars := map[string]string{"HOST": "db.prod", "PORT": "5432"}
	raw := map[string]string{
		"DATABASE_URL": "postgres://${HOST}:${PORT}/app",
		"REDIS_URL":    "redis://${HOST}:6379",
	}

	rendered, err := envfile.RenderMap(raw, vars, envfile.DefaultTemplateOptions())
	if err != nil {
		t.Fatal(err)
	}

	if err := envfile.Write(dst, rendered); err != nil {
		t.Fatal(err)
	}

	parsed, err := envfile.Parse(dst)
	if err != nil {
		t.Fatal(err)
	}

	if parsed["DATABASE_URL"] != "postgres://db.prod:5432/app" {
		t.Errorf("DATABASE_URL: got %q", parsed["DATABASE_URL"])
	}
	if parsed["REDIS_URL"] != "redis://db.prod:6379" {
		t.Errorf("REDIS_URL: got %q", parsed["REDIS_URL"])
	}
}

// TestRenderFile_RoundTrip writes a template, renders it, then parses the result.
func TestRenderFile_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "tpl.env")
	dst := filepath.Join(dir, "rendered.env")

	tpl := "APP_SECRET=${SECRET}\nAPP_ENV=${ENV}\n"
	_ = os.WriteFile(src, []byte(tpl), 0o600)

	vars := map[string]string{"SECRET": "s3cr3t", "ENV": "staging"}
	if err := envfile.RenderFile(src, dst, vars, envfile.DefaultTemplateOptions()); err != nil {
		t.Fatal(err)
	}

	parsed, err := envfile.Parse(dst)
	if err != nil {
		t.Fatal(err)
	}

	if parsed["APP_SECRET"] != "s3cr3t" {
		t.Errorf("APP_SECRET: got %q", parsed["APP_SECRET"])
	}
	if parsed["APP_ENV"] != "staging" {
		t.Errorf("APP_ENV: got %q", parsed["APP_ENV"])
	}
}

// TestRenderMap_MissingKeyPropagatesError ensures errors surface from nested maps.
func TestRenderMap_MissingKeyPropagatesError(t *testing.T) {
	m := map[string]string{"KEY": "${UNDEFINED}"}
	_, err := envfile.RenderMap(m, map[string]string{}, envfile.DefaultTemplateOptions())
	if err == nil {
		t.Fatal("expected error")
	}
}
