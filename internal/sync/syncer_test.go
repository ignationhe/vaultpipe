package sync_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/yourusername/vaultpipe/internal/sync"
	"github.com/yourusername/vaultpipe/internal/vault"
)

func makeVaultServer(t *testing.T, data map[string]string) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		payload := map[string]interface{}{
			"data": map[string]interface{}{"data": data},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(payload) //nolint:errcheck
	}))
}

func TestSync_WritesSecretsToEnvFile(t *testing.T) {
	srv := makeVaultServer(t, map[string]string{"FOO": "bar", "BAZ": "qux"})
	defer srv.Close()

	client, err := vault.NewClient(srv.URL, "test-token")
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}

	tmpDir := t.TempDir()
	envPath := filepath.Join(tmpDir, ".env")

	s := sync.New(client)
	result, err := s.Sync("secret/data/app", envPath, sync.Options{})
	if err != nil {
		t.Fatalf("Sync: %v", err)
	}

	if !result.Written {
		t.Error("expected Written=true for new file")
	}

	if _, err := os.Stat(envPath); os.IsNotExist(err) {
		t.Error("expected env file to exist after sync")
	}
}

func TestSync_DryRunDoesNotWrite(t *testing.T) {
	srv := makeVaultServer(t, map[string]string{"KEY": "value"})
	defer srv.Close()

	client, err := vault.NewClient(srv.URL, "test-token")
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}

	tmpDir := t.TempDir()
	envPath := filepath.Join(tmpDir, ".env")

	s := sync.New(client)
	result, err := s.Sync("secret/data/app", envPath, sync.Options{DryRun: true})
	if err != nil {
		t.Fatalf("Sync: %v", err)
	}

	if result.Written {
		t.Error("expected Written=false in dry-run mode")
	}

	if _, err := os.Stat(envPath); !os.IsNotExist(err) {
		t.Error("expected env file NOT to exist in dry-run mode")
	}
}

func TestSync_NoWriteWhenUnchanged(t *testing.T) {
	srv := makeVaultServer(t, map[string]string{"STABLE": "yes"})
	defer srv.Close()

	client, err := vault.NewClient(srv.URL, "test-token")
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}

	tmpDir := t.TempDir()
	envPath := filepath.Join(tmpDir, ".env")

	// Pre-populate the file with the same value Vault will return.
	if err := os.WriteFile(envPath, []byte("STABLE=yes\n"), 0o600); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	s := sync.New(client)
	result, err := s.Sync("secret/data/app", envPath, sync.Options{})
	if err != nil {
		t.Fatalf("Sync: %v", err)
	}

	if result.Written {
		t.Error("expected Written=false when secrets are unchanged")
	}
}
