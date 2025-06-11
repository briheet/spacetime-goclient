package httpClient

import (
	"fmt"
	"net/http"
	"time"
)

type Client struct {
	BaseURL    string
	HTTPClient *http.Client
}

type Option func(*Client) error

func WithBaseURL(url string) Option {
	return func(c *Client) error {
		if url == "" {
			return fmt.Errorf("base URL cannot be empty")
		}
		c.BaseURL = url
		return nil
	}
}

func WithTimeout(timeout time.Duration) Option {
	return func(c *Client) error {
		c.HTTPClient.Timeout = timeout
		return nil
	}
}

// TODO: Custom checks needs to be done
// WithCustomHTTPClient allows injecting a fully configured *http.Client.
func WithCustomHTTPClient(hc *http.Client) Option {
	return func(c *Client) error {
		if hc == nil {
			return fmt.Errorf("custom http client cannot be nil")
		}
		c.HTTPClient = hc
		return nil
	}
}

func NewClient(opts ...Option) (*Client, error) {
	client := &Client{
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}

	for _, opt := range opts {
		if err := opt(client); err != nil {
			return nil, err
		}
	}

	// Recheck
	if client.BaseURL == "" {
		return nil, fmt.Errorf("BaseURL is required (use WithBaseURL)")
	}

	return client, nil
}
