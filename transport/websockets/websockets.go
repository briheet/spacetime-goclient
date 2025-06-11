package websocketsClient

import (
	"net/url"
	"fmt"
	"github.com/gorilla/websocket"
)

type Conn struct {
	WS *websocket.Conn
}

func NewConn(baseURL string) (*Conn, error) {
	u := url.URL{
		Scheme: "ws",
		Host:   baseURL,
		Path:   "/v1/ws", // TODO: Need to change this accordingly
	}

	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("websocket dial failed: %w", err)
	}

	return &Conn{WS: conn}, nil
}

func (c *Conn) Close() error {
	if c.WS != nil {
		return c.WS.Close()
	}
	return nil
}