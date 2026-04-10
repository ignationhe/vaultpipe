package vault

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func makeServer(t *testing.T, status int, payload interface{}) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(status)
		if payload != nil {
			_ = json.NewEncoder(w).Encode(payload)
		}
	}))
}

func kvPayload(data map[string]interface{}) map[string]interface{} {
	return map[string]interface{}{
		"data": map[string]interface{}{
			"data": data,
		},
	}
}

func TestGetSecrets_ReturnsKeyValues(t *testing.T) {
	srv := makeServer(t, http.StatusOK, kvPayload(map[string]interface{}{
		"DB_HOST": "localhost",
		"DB_PORT": "5432",
	}))
	defer srv.Close()

	client := NewClient(srv.URL, "test-token")
	secrets, err := client.GetSecrets("secret/data/myapp")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if secrets["DB_HOST"] != "localhost" {
		t.Errorf("expected DB_HOST=localhost, got %q", secrets["DB_HOST"])
	}
	if secrets["DB_PORT"] != "5432" {
		t.Errorf("expected DB_PORT=5432, got %q", secrets["DB_PORT"])
	}
}

func TestGetSecrets_ReturnsErrorOn403(t *testing.T) {
	srv := makeServer(t, http.StatusForbidden, nil)
	defer srv.Close()

	client := NewClient(srv.URL, "bad-token")
	_, err := client.GetSecrets("secret/data/myapp")
	if err == nil {
		t.Fatal("expected error for 403, got nil")
	}
}

func TestGetSecrets_ReturnsErrorOn404(t *testing.T) {
	srv := makeServer(t, http.StatusNotFound, nil)
	defer srv.Close()

	client := NewClient(srv.URL, "test-token")
	_, err := client.GetSecrets("secret/data/missing")
	if err == nil {
		t.Fatal("expected error for 404, got nil")
	}
}

func TestParseResponse_HandlesEmptyData(t *testing.T) {
	body, _ := json.Marshal(kvPayload(map[string]interface{}{}))
	result, err := parseResponse(body)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 0 {
		t.Errorf("expected empty map, got %v", result)
	}
}

func TestParseResponse_InvalidJSON(t *testing.T) {
	_, err := parseResponse([]byte("not-json"))
	if err == nil {
		t.Fatal("expected error for invalid JSON, got nil")
	}
}
