package spacetimedb

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"os"
)

type Database interface {
	PublishDatabase(wasmFile string, token string) (string, string, error)
	PublishNamedDatabase(nameOrIdentity string, wasmFile string, token string, clear bool) (string, string, *string, error)
	GetDatabaseInfo(nameOrIdentity string) (string, string, string, string, error)
}

func (c *Client) PublishDatabase(wasmFile string, token string) (string, string, error) {
	// Need to read the file
	data, err := os.ReadFile(wasmFile)
	if err != nil {
		return "", "", fmt.Errorf("reading wasm: %w", err)
	}

	url := fmt.Sprintf("%s/v1/database", c.BaseURL)

	headers := map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", token),
	}

	resp, err := c.HTTPClient.Do("POST", url, headers, bytes.NewReader(data))
	if err != nil {
		return "", "", fmt.Errorf("failed to publish data: %w", err)
	}

	var result struct {
		Success struct {
			DatabaseIdentity string `json:"database_identity"`
			Op               string `json:"op"` // "created" or "updated"
		} `json:"Success"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", "", fmt.Errorf("failed to parse publish response: %w", err)
	}

	return result.Success.DatabaseIdentity, result.Success.Op, nil
}

func (c *Client) PublishNamedDatabase(nameOrIdentity string, wasmFile string, token string, clear bool) (string, string, *string, error) {
	// Read the Wasm binary
	data, err := os.ReadFile(wasmFile)
	if err != nil {
		return "", "", nil, fmt.Errorf("reading wasm file: %w", err)
	}

	query := url.Values{}
	if clear {
		query.Set("clear", "true")
	}

	fullPath := fmt.Sprintf("/v1/database/%s", nameOrIdentity)
	if q := query.Encode(); q != "" {
		fullPath += "?" + q
	}

	// Set headers
	headers := map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", token),
	}

	// Make the request
	resp, err := c.HTTPClient.Do("POST", fullPath, headers, bytes.NewReader(data))
	if err != nil {
		return "", "", nil, fmt.Errorf("failed to publish named database: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	// Handle permission denied
	if resp.StatusCode == 401 {
		var denied struct {
			PermissionDenied struct {
				Name string `json:"name"`
			} `json:"PermissionDenied"`
		}
		if err := json.Unmarshal(body, &denied); err != nil {
			return "", "", nil, fmt.Errorf("unauthorized, and failed to parse PermissionDenied: %w", err)
		}
		return "", "", nil, fmt.Errorf("permission denied to publish database: %s", denied.PermissionDenied.Name)
	}

	// Handle success
	if resp.StatusCode != 200 {
		return "", "", nil, fmt.Errorf("unexpected status code %d: %s", resp.StatusCode, string(body))
	}

	var result struct {
		Success struct {
			Domain           *string `json:"domain"`
			DatabaseIdentity string  `json:"database_identity"`
			Op               string  `json:"op"`
		} `json:"Success"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", "", nil, fmt.Errorf("failed to parse publish response: %w", err)
	}

	return result.Success.DatabaseIdentity, result.Success.Op, result.Success.Domain, nil
}

func (c *Client) GetDatabaseInfo(nameOrIdentity string) (string, string, string, string, error) {
	endpoint := fmt.Sprintf("/v1/database/%s", url.PathEscape(nameOrIdentity))

	resp, err := c.HTTPClient.Do("GET", endpoint, nil, nil)
	if err != nil {
		return "", "", "", "", fmt.Errorf("failed to get database info: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return "", "", "", "", fmt.Errorf("unexpected status code %d: %s", resp.StatusCode, string(body))
	}

	var result struct {
		DatabaseIdentity string `json:"database_identity"`
		OwnerIdentity    string `json:"owner_identity"`
		HostType         string `json:"host_type"`
		InitialProgram   string `json:"initial_program"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", "", "", "", fmt.Errorf("failed to decode response: %w", err)
	}

	return result.DatabaseIdentity, result.OwnerIdentity, result.HostType, result.InitialProgram, nil
}
