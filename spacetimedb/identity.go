package spacetimedb

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type Identity interface {
	CreateIdentity() (string, string, error)
	CreateIdentityWebsocketToken() (string, error)
	GetPublicKey() (string, error)
	RegisterIdentityWithEmail(string) (string, string, error)
	GetDatabasesByIdentity(string) ([]string, error)
	VerifyIdentityToken(string, string) error
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

// FIXME: Endpoint issue
func (c *Client) RegisterIdentityWithEmail(email string) (string, string, error) {
	identity, token, err := c.CreateIdentity()
	if err != nil {
		return "", "", fmt.Errorf("failed to create identity: %w", err)
	}

	// Base path
	basePath := fmt.Sprintf("/v1/identity/%s/set-email", identity)

	// Attach email as query param
	endpoint := fmt.Sprintf("%s?email=%s", basePath, email)

	fmt.Println(endpoint)

	headers := map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", token),
	}

	resp, err := c.HTTPClient.Do("POST", endpoint, headers, nil)
	if err != nil {
		return "", "", fmt.Errorf("failed to set email: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return "", "", fmt.Errorf("set-email failed: status code %d", resp.StatusCode)
	}

	return identity, token, nil
}

func (c *Client) GetDatabasesByIdentity(identity string) ([]string, error) {

	endpoint := fmt.Sprintf("/v1/identity/%s/databases", url.PathEscape(identity))

	resp, err := c.HTTPClient.Do("GET", endpoint, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get databases: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("get databases failed: status code %d", resp.StatusCode)
	}

	// Parse the JSON response
	var result struct {
		Addresses []string `json:"addresses"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return result.Addresses, nil

}

func (c *Client) VerifyIdentityToken(identity, token string) error {
	// Construct the endpoint URL
	endpoint := fmt.Sprintf("/v1/identity/%s/verify", url.PathEscape(identity))

	// Build the Authorization header using Bearer token
	headers := map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", token),
	}

	// Perform the GET request
	resp, err := c.HTTPClient.Do("GET", endpoint, headers, nil)
	if err != nil {
		return fmt.Errorf("verify request failed: %w", err)
	}
	defer resp.Body.Close()

	// Handle response status codes
	switch resp.StatusCode {
	case http.StatusNoContent:
		return nil
	case http.StatusBadRequest:
		return fmt.Errorf("token valid but does not match identity (400)")
	case http.StatusUnauthorized:
		return fmt.Errorf("invalid or missing token (401)")
	default:
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
}
