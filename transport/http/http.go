package httpClient

import (
	"fmt"
	"io"
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

func (c *Client) Do(method, pathURL string, headers map[string]string, body io.Reader) (*http.Response, error) {

	// Build the url
	fullURL := c.BaseURL + pathURL

	// Create a request
	req, err := http.NewRequest(method, fullURL, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set respective header
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	// Do a req
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}

	return resp, nil
}
