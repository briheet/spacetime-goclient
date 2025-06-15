package spacetimedb

import (
	"fmt"
	"time"

	httpClient "github.com/briheet/spacetime-goclient/transport/http"
	websocketsClient "github.com/briheet/spacetime-goclient/transport/websockets"
)

type DBClient interface {
	CallReducer(name string, args map[string]any) error
	Subscribe(query string, handler func(snapshot []any, diff any)) error
	Disconnect() error
	Ping() error

	// Identity Methods
	Identity

	// Database Methods
	Database
}

var _ DBClient = (*Client)(nil)

type Client struct {
	BaseURL string
	DBName  string
	// Http
	HTTPClient *httpClient.Client
	// Websockets for sub
	WebsocketClient *websocketsClient.Conn

	// Identity and Token
	Identity string
	Token    string
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
