package spacetimedb

import (
	"fmt"
	"time"

	"github.com/briheet/spacetimedb-go-client/transport/http"
	"github.com/briheet/spacetimedb-go-client/transport/websockets"
)

type DBClient interface {
	CallReducer(name string, args map[string]any) error
	Subscribe(query string, handler func(snapshot []any, diff any)) error
	Disconnect() error
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