package spacetimedb

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"github.com/gorilla/websocket"
)

type Database interface {
	PublishDatabase(wasmFile string, token string) (string, string, error)
	PublishNamedDatabase(nameOrIdentity string, wasmFile string, token string, clear bool) (string, string, *string, error)
	GetDatabaseInfo(nameOrIdentity string) (string, string, string, string, error)
	DeleteDatabase(nameOrIdentity, token string) error
	GetDatabaseNames(nameOrIdentity string) ([]string, error)
	AddDatabaseName(nameOrIdentity, newName, token string) error
	GetDatabaseIdentity(nameOrIdentity string) (string, error)
	WebsocketSubscribe(dbNameOrIden, token, protocol string) (*websocket.Conn, error)
	GetDatabaseLogs(dbNameOrIden string, token string, numLines int, follow bool) (io.ReadCloser, error)
	RunSQLQuery(query, token, dbName string) ([]SQLResult, error)
}

func (c *Client) PublishDatabase(wasmFile string, token string) (string, string, error) {
	// Need to read the file
	data, err := os.ReadFile(wasmFile)
	if err != nil {
		return "", "", fmt.Errorf("reading wasm: %w", err)
	}

	url := fmt.Sprintf("%s/v1/database", c.HttpBaseURL)

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
		DatabaseIdentity struct {
			Identity string `json:"__identity__"`
		} `json:"database_identity"`
		OwnerIdentity struct {
			Identity string `json:"__identity__"`
		} `json:"owner_identity"`
		HostType       map[string]interface{} `json:"host_type"`
		InitialProgram string                 `json:"initial_program"`
	}

	bodyBytes, _ := io.ReadAll(resp.Body)

	if err := json.Unmarshal(bodyBytes, &result); err != nil {
		return "", "", "", "", fmt.Errorf("failed to decode response: %w", err)
	}

	hostType := ""
	for k := range result.HostType {
		hostType = k
		break
	}
	return result.DatabaseIdentity.Identity, result.OwnerIdentity.Identity, hostType, result.InitialProgram, nil
}

func (c *Client) DeleteDatabase(nameOrIdentity string, token string) error {

	endpoint := fmt.Sprintf("/v1/database/%s", url.PathEscape(nameOrIdentity))

	// Build up the header
	headers := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": fmt.Sprintf("Bearer %s", token),
	}

	resp, err := c.HTTPClient.Do("DELETE", endpoint, headers, nil)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	// Read response body (always)
	bodyBytes, _ := io.ReadAll(resp.Body)
	bodyText := string(bodyBytes)

	if resp.StatusCode != 200 {
		return fmt.Errorf("unexpected status code %d: %s", resp.StatusCode, bodyText)
	}

	return nil
}

func (c *Client) GetDatabaseNames(nameOrIdentity string) ([]string, error) {
	endpoint := fmt.Sprintf("/v1/database/%s/names", url.PathEscape(nameOrIdentity))

	resp, err := c.HTTPClient.Do("GET", endpoint, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status code %d: %s", resp.StatusCode, string(body))
	}

	var result struct {
		Names []string `json:"names"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return result.Names, nil
}

func (c *Client) AddDatabaseName(nameOrIdentity, newName, token string) error {
	endpoint := fmt.Sprintf("/v1/database/%s/names", url.PathEscape(nameOrIdentity))

	headers := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": fmt.Sprintf("Bearer %s", token),
	}

	body, err := json.Marshal(newName)
	if err != nil {
		return fmt.Errorf("failed to marshal name: %w", err)
	}

	resp, err := c.HTTPClient.Do("POST", endpoint, headers, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != 200 {
		return fmt.Errorf("failed to set name, status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	// Check for success or permission denied
	var result map[string]json.RawMessage
	if err := json.Unmarshal(bodyBytes, &result); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	if _, ok := result["Success"]; ok {
		return nil
	}
	if val, ok := result["PermissionDenied"]; ok {
		var denial struct {
			Domain string `json:"domain"`
		}
		json.Unmarshal(val, &denial)
		return fmt.Errorf("permission denied for domain: %s", denial.Domain)
	}

	return fmt.Errorf("unexpected response: %s", string(bodyBytes))
}

func (c *Client) GetDatabaseIdentity(nameOrIdentity string) (string, error) {
	endpoint := fmt.Sprintf("/v1/database/%s/identity", url.PathEscape(nameOrIdentity))

	resp, err := c.HTTPClient.Do("GET", endpoint, nil, nil)
	if err != nil {
		return "", fmt.Errorf("failed to get database identity: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("unexpected status code %d: %s", resp.StatusCode, string(body))
	}

	// Read raw string (not JSON)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	return string(body), nil
}

func (c *Client) WebsocketSubscribe(dbNameOrIden, token, protocol string) (*websocket.Conn, error) {

	if protocol == "" {
		protocol = "v1.json.spacetimedb" // default to JSON protocol
	}

	// Prepare headers
	headers := http.Header{}
	headers.Set("Sec-WebSocket-Protocol", protocol)
	// headers.Set("Sec-WebSocket-Version", "13")

	if token != "" {
		headers.Set("Authorization", "Bearer "+token)
	}

	endpoint := fmt.Sprintf("/v1/database/%s/subscribe", dbNameOrIden)
	fullURL := c.WssBaseURL + endpoint

	conn, resp, err := c.WebsocketClient.Connect(fullURL, headers)
	if err != nil {
		if resp != nil {
			return nil, fmt.Errorf("WebSocket handshake failed: status %s", resp.Status)
		}
		return nil, fmt.Errorf("WebSocket connection error: %w", err)
	}

	return conn, nil

}

func (c *Client) GetDatabaseLogs(dbNameOrIden string, token string, numLines int, follow bool) (io.ReadCloser, error) {
	endpoint := fmt.Sprintf("%s/v1/database/%s/logs", c.HttpBaseURL, dbNameOrIden)

	// Build query params
	params := url.Values{}
	if numLines > 0 {
		params.Set("num_lines", strconv.Itoa(numLines))
	}
	if follow {
		params.Set("follow", "true")
	}
	if len(params) > 0 {
		endpoint += "?" + params.Encode()
	}

	headers := map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", token),
	}

	resp, err := c.HTTPClient.Do("GET", endpoint, headers, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to request logs: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		defer resp.Body.Close()
		return nil, fmt.Errorf("unexpected status code: %s", resp.Status)
	}

	// Caller is responsible for closing resp.Body
	return resp.Body, nil
}

type SQLResult struct {
	Schema any           `json:"schema"`
	Rows   []interface{} `json:"rows"`
}

func (c *Client) RunSQLQuery(query, token, dbName string) ([]SQLResult, error) {
	if query == "" {
		return nil, fmt.Errorf("query cannot be empty")
	}

	endpoint := fmt.Sprintf("/v1/database/%s/sql", dbName)
	body := bytes.NewBufferString(query)

	headers := map[string]string{
		"Content-Type":  "text/plain", // Spacetime expects raw SQL
		"Authorization": fmt.Sprintf("Bearer %s", token),
	}

	resp, err := c.HTTPClient.Do("POST", endpoint, headers, body)
	if err != nil {
		return nil, fmt.Errorf("failed to run SQL query: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		msg, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("SQL query failed: %s — %s", resp.Status, msg)
	}

	var results []SQLResult
	if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
		return nil, fmt.Errorf("failed to decode SQL response: %w", err)
	}

	return results, nil
}
