package spacetimedb

import (
	"encoding/json"
	"fmt"
	"io"
	"time"

	httpClient "github.com/briheet/spacetime-goclient/transport/http"
	websocketsClient "github.com/briheet/spacetime-goclient/transport/websockets"
)

type DBClient interface {
	CallReducer(name string, args map[string]any) error
	Subscribe(query string, handler func(snapshot []any, diff any)) error
	Disconnect() error
	Ping() error
	CreateIdentity() (string, string, error)
	CreateIdentityWebsocketToken() (string, error)
}

var _ DBClient = (*Client)(nil)

type Client struct {
	BaseURL string
	DBName  string
	// Http
	HTTPClient *httpClient.Client
	// Websockets for sub
	WebsocketClient *websocketsClient.Conn
}

func Connect(url string, port string, dbName string) (DBClient, error) {
	base := fmt.Sprintf("%s:%s", url, port)

	httpClient, err := httpClient.NewClient(
		httpClient.WithBaseURL(base),
		httpClient.WithTimeout(3*time.Second),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP client: %w", err)
	}

	// TODO: Needs to be changed
	websocketClient := &websocketsClient.Conn{}

	return &Client{
		BaseURL:         base,
		DBName:          dbName,
		HTTPClient:      httpClient,
		WebsocketClient: websocketClient,
	}, nil
}

func (c *Client) CallReducer(name string, args map[string]any) error {
	// TODO: Implement actual reducer call
	return nil
}

func (c *Client) Subscribe(query string, handler func(snapshot []any, diff any)) error {
	// TODO: Implement actual websocket subscription
	return nil
}

func (c *Client) Disconnect() error {
	// Clean up the Http and Websocket conn

	if c.WebsocketClient != nil {
		if err := c.WebsocketClient.Close(); err != nil {
			return fmt.Errorf("failed to close websocket connection: %w", err)
		}
	}
	return nil
}

func (c *Client) Ping() error {
	resp, err := c.HTTPClient.Do("GET", "/v1/ping", nil, nil)
	if err != nil {
		return fmt.Errorf("ping failed: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("ping failed: status code %d", resp.StatusCode)
	}

	return nil
}

func (c *Client) CreateIdentity() (string, string, error) {
	resp, err := c.HTTPClient.Do("POST", "/v1/identity", nil, nil)
	if err != nil {
		return "", "", fmt.Errorf("failed to create identity: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", "", fmt.Errorf("create identity failed: status code %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", fmt.Errorf("failed to read response body: %w", err)
	}

	var parsed struct {
		Identity string `json:"identity"`
		Token    string `json:"token"`
	}

	if err := json.Unmarshal(body, &parsed); err != nil {
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
