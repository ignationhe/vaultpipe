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

// newTestSyncer creates a Syncer backed by a test Vault server returning the
// given secrets. It registers cleanup of the server via t.Cleanup.
func newTestSyncer(t *testing.T, secrets map[string]interface{}) *Syncer {
	t.Helper()
	srv := makeVaultServer(t, secrets)
	t.Cleanup(srv.Close)
	client, err := vaultclient.NewClient(srv.URL, "test-token")
	if err != nil {
		t.Fatalf("failed to create vault client: %v", err)
	}
	return New(client)
}

func TestSync_WritesSecretsToEnvFile(t *testing.T) {
	syncer := newTestSyncer(t, map[string]interface{}{"KEY": "value"})

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
	syncer := newTestSyncer(t, map[string]interface{}{"KEY": "value"})

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
	syncer := newTestSyncer(t, map[string]interface{}{"KEY": "value"})

	tmp := filepath.Join(t.TempDir(), ".env")
	_ = envfile.Write(tmp, map[string]string{"KEY": "value"})

	res, err := syncer.Sync(Options{VaultPath: "secret/data/app", EnvFile: tmp})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Written {
		t.Error("expected Written=false when secrets are unchanged")
	}
}

func TestSync_FilterApplied(t *testing.T) {
	syncer := newTestSyncer(t, map[string]interface{}{"APP_KEY": "v1", "OTHER": "v2"})

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
