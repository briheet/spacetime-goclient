package websocketsClient

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
)

type Client struct {
	BaseURL string
	Dialer  *websocket.Dialer
}

type Option func(*Client) error

func WithBaseURL(base string) Option {
	return func(c *Client) error {
		if base == "" {
			return fmt.Errorf("BaseURL cannot be empty")
		}
		c.BaseURL = base
		return nil
	}
}

func WithDialTimeout(timeout time.Duration) Option {
	return func(c *Client) error {
		c.Dialer.HandshakeTimeout = timeout
		return nil
	}
}

func WithCustomDialer(d *websocket.Dialer) Option {
	return func(c *Client) error {
		if d == nil {
			return fmt.Errorf("custom dialer cannot be nil")
		}
		c.Dialer = d
		return nil
	}
}

func NewClient(opts ...Option) (*Client, error) {
	client := &Client{
		Dialer: &websocket.Dialer{
			Proxy:             http.ProxyFromEnvironment,
			HandshakeTimeout:  10 * time.Second,
			EnableCompression: true,
		},
	}

	for _, opt := range opts {
		if err := opt(client); err != nil {
			return nil, err
		}
	}

	if client.BaseURL == "" {
		return nil, fmt.Errorf("BaseURL is required (use WithBaseURL)")
	}

	return client, nil
}

// Connect opens a websocket connection to the path with *dynamic headers*
func (c *Client) Connect(path string, headers http.Header) (*websocket.Conn, *http.Response, error) {
	wsURL, err := url.JoinPath(c.BaseURL, path)
	if err != nil {
		return nil, nil, fmt.Errorf("invalid URL: %w", err)
	}

	conn, resp, err := c.Dialer.Dial(wsURL, headers)
	if err != nil {
		return nil, resp, fmt.Errorf("failed to connect: %w", err)
	}

	return conn, resp, nil
}
