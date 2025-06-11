package spacetimedb

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Identity interface {
	CreateIdentity() (string, string, error)
	CreateIdentityWebsocketToken() (string, error)
	GetPublicKey() (string, error)
	RegisterIdentityWithEmail(string) error
}

func (c *Client) CreateIdentity() (string, string, error) {
	// Make a request
	resp, err := c.HTTPClient.Do("POST", "/v1/identity", nil, nil)
	if err != nil {
		return "", "", fmt.Errorf("failed to create identity: %w", err)
	}

	defer resp.Body.Close()

	// Check for status code
	if resp.StatusCode != 200 {
		return "", "", fmt.Errorf("create identity failed: status code %d", resp.StatusCode)
	}

	var parsed struct {
		Identity string `json:"identity"`
		Token    string `json:"token"`
	}

	// Parse that shii
	if err := json.NewDecoder(resp.Body).Decode(&parsed); err != nil {
		return "", "", fmt.Errorf("failed to parse identity response: %w", err)
	}

	return parsed.Identity, parsed.Token, nil
}

func (c *Client) CreateIdentityWebsocketToken() (string, error) {
	// First get a token which is needed for websocketToken call
	_, token, err := c.CreateIdentity()
	if err != nil {
		return "", fmt.Errorf("unable to generate token %w", err)
	}

	// Construct headers which would be added in the request
	headers := map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", token),
	}

	// Make the request
	resp, err := c.HTTPClient.Do("POST", "/v1/identity/websocket-token", headers, nil)
	if err != nil {
		return "", fmt.Errorf("websocket-token request failed: %w", err)
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("websocket-token request failed: status code %d", resp.StatusCode)
	}

	var result struct {
		Token string `json:"token"`
	}

	// Parse that shii
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to decode websocket token: %w", err)
	}

	return result.Token, nil
}

func (c *Client) GetPublicKey() (string, error) {
	// Make a request
	resp, err := c.HTTPClient.Do("GET", "/v1/identity/public-key", nil, nil)
	if err != nil {
		return "", fmt.Errorf("public key request failed: %w", err)
	}

	defer resp.Body.Close()

	// Check for non-200 status code
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("public key request failed: status code %d", resp.StatusCode)
	}

	// Read the PEM-encoded public key as plain text
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read public key response: %w", err)
	}

	return string(body), nil
}

func (c *Client) RegisterIdentityWithEmail(string) error {
	return nil
}
