package vault

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Client is a minimal HashiCorp Vault HTTP client.
type Client struct {
	Address string
	Token   string
	http    *http.Client
}

// NewClient creates a new Vault client with the given address and token.
func NewClient(address, token string) *Client {
	return &Client{
		Address: address,
		Token:   token,
		http: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// GetSecrets reads a KV v2 secret at the given path and returns
// the key/value pairs stored under the "data" field.
func (c *Client) GetSecrets(path string) (map[string]string, error) {
	url := fmt.Sprintf("%s/v1/%s", c.Address, path)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("vault: building request: %w", err)
	}
	req.Header.Set("X-Vault-Token", c.Token)

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("vault: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusForbidden || resp.StatusCode == http.StatusUnauthorized {
		return nil, fmt.Errorf("vault: authentication failed (status %d)", resp.StatusCode)
	}
	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("vault: secret not found at path %q", path)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("vault: unexpected status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("vault: reading response: %w", err)
	}

	return parseResponse(body)
}

type kvResponse struct {
	Data struct {
		Data map[string]interface{} `json:"data"`
	} `json:"data"`
}

func parseResponse(body []byte) (map[string]string, error) {
	var resp kvResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("vault: parsing response: %w", err)
	}

	result := make(map[string]string, len(resp.Data.Data))
	for k, v := range resp.Data.Data {
		result[k] = fmt.Sprintf("%v", v)
	}
	return result, nil
}
