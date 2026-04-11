package envfile

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRender_BasicSubstitution(t *testing.T) {
	vars := map[string]string{"HOST": "localhost", "PORT": "5432"}
	got, err := Render("DB=${HOST}:${PORT}", vars, DefaultTemplateOptions())
	if err != nil {
		t.Fatal(err)
	}
	if got != "DB=localhost:5432" {
		t.Errorf("got %q", got)
	}
}

func TestRender_MissingKeyError(t *testing.T) {
	_, err := Render("X=${MISSING}", map[string]string{}, DefaultTemplateOptions())
	if err == nil {
		t.Fatal("expected error for missing key")
	}
}

func TestRender_MissingKeyKeep(t *testing.T) {
	opts := TemplateOptions{MissingKey: "keep"}
	got, err := Render("X=${MISSING}", map[string]string{}, opts)
	if err != nil {
		t.Fatal(err)
	}
	if got != "X=${MISSING}" {
		t.Errorf("got %q", got)
	}
}

func TestRender_MissingKeyEmpty(t *testing.T) {
	opts := TemplateOptions{MissingKey: "empty"}
	got, err := Render("X=${MISSING}", map[string]string{}, opts)
	if err != nil {
		t.Fatal(err)
	}
	if got != "X=" {
		t.Errorf("got %q", got)
	}
}

func TestRenderMap_SubstitutesValues(t *testing.T) {
	vars := map[string]string{"ENV": "production"}
	m := map[string]string{"APP_ENV": "${ENV}", "PLAIN": "hello"}
	out, err := RenderMap(m, vars, DefaultTemplateOptions())
	if err != nil {
		t.Fatal(err)
	}
	if out["APP_ENV"] != "production" {
		t.Errorf("APP_ENV: got %q", out["APP_ENV"])
	}
	if out["PLAIN"] != "hello" {
		t.Errorf("PLAIN: got %q", out["PLAIN"])
	}
}

func TestRenderFile_CreatesOutputFile(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "tpl.env")
	dst := filepath.Join(dir, "out.env")

	_ = os.WriteFile(src, []byte("DB_URL=${SCHEME}://${HOST}\n"), 0o600)

	vars := map[string]string{"SCHEME": "postgres", "HOST": "db.local"}
	if err := RenderFile(src, dst, vars, DefaultTemplateOptions()); err != nil {
		t.Fatal(err)
	}

	data, _ := os.ReadFile(dst)
	if string(data) != "DB_URL=postgres://db.local\n" {
		t.Errorf("unexpected output: %q", string(data))
	}
}

func TestRenderFile_MissingSrcReturnsError(t *testing.T) {
	err := RenderFile("/nonexistent/tpl.env", "/tmp/out.env", nil, DefaultTemplateOptions())
	if err == nil {
		t.Fatal("expected error for missing source file")
	}
}
