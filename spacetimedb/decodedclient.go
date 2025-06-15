package spacetimedb

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
)

type DecodedClient interface {
	SendMessageDatabase(string, string, string, string) error
}

func (c *Client) SendMessageDatabase(reducerName string, dbID string, token string, text string) error {

	// Setup the data that needs to be sent
	payload := map[string]interface{}{
		"text": text,
	}
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	// Build up the header
	headers := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": fmt.Sprintf("Bearer %s", token),
	}

	// Buildup the url
	url := fmt.Sprintf("/v1/database/%s/call/%s", dbID, reducerName)
	resp, err := c.HTTPClient.Do("POST", url, headers, bytes.NewReader(jsonData))
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
