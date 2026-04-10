package sync

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/yourusername/vaultpipe/internal/envfile"
	vaultclient "github.com/yourusername/vaultpipe/internal/vault"
)

func makeVaultServer(t *testing.T, data map[string]interface{}) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"data": map[string]interface{}{"data": data},
		})
	}))
}

func TestSync_WritesSecretsToEnvFile(t *testing.T) {
	srv := makeVaultServer(t, map[string]interface{}{"KEY": "value"})
	defer srv.Close()

	client, _ := vaultclient.NewClient(srv.URL, "test-token")
	syncer := New(client)

	tmp := filepath.Join(t.TempDir(), ".env")
	res, err := syncer.Sync(Options{VaultPath: "secret/data/app", EnvFile: tmp})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.Written {
		t.Error("expected Written=true")
	}
	parsed, _ := envfile.Parse(tmp)
	if parsed["KEY"] != "value" {
		t.Errorf("expected KEY=value, got %s", parsed["KEY"])
	}
}

func TestSync_DryRunDoesNotWrite(t *testing.T) {
	srv := makeVaultServer(t, map[string]interface{}{"KEY": "value"})
	defer srv.Close()

	client, _ := vaultclient.NewClient(srv.URL, "test-token")
	syncer := New(client)

	tmp := filepath.Join(t.TempDir(), ".env")
	res, err := syncer.Sync(Options{VaultPath: "secret/data/app", EnvFile: tmp, DryRun: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Written {
		t.Error("expected Written=false for dry run")
	}
	if _, err := os.Stat(tmp); !os.IsNotExist(err) {
		t.Error("env file should not exist after dry run")
	}
}

func TestSync_NoWriteWhenUnchanged(t *testing.T) {
	srv := makeVaultServer(t, map[string]interface{}{"KEY": "value"})
	defer srv.Close()

	tmp := filepath.Join(t.TempDir(), ".env")
	_ = envfile.Write(tmp, map[string]string{"KEY": "value"})

	client, _ := vaultclient.NewClient(srv.URL, "test-token")
	syncer := New(client)

	res, err := syncer.Sync(Options{VaultPath: "secret/data/app", EnvFile: tmp})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Written {
		t.Error("expected Written=false when secrets are unchanged")
	}
}

func TestSync_FilterApplied(t *testing.T) {
	srv := makeVaultServer(t, map[string]interface{}{"APP_KEY": "v1", "OTHER": "v2"})
	defer srv.Close()

	client, _ := vaultclient.NewClient(srv.URL, "test-token")
	syncer := New(client)

	tmp := filepath.Join(t.TempDir(), ".env")
	_, err := syncer.Sync(Options{
		VaultPath: "secret/data/app",
		EnvFile:   tmp,
		Filter:    envfile.FilterOptions{Prefix: "APP_"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	parsed, _ := envfile.Parse(tmp)
	if _, ok := parsed["OTHER"]; ok {
		t.Error("OTHER should have been filtered out")
	}
	if parsed["APP_KEY"] != "v1" {
		t.Errorf("expected APP_KEY=v1, got %s", parsed["APP_KEY"])
	}
}
