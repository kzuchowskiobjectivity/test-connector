package miro

import (
	"fmt"
	"io"
	"net/http"
)

type Client struct {
	client  *http.Client
	url     string
	ownerId string
	apiKey  string
}

func NewClient(url string, ownerId, apiKey string) *Client {
	return &Client{
		client:  http.DefaultClient,
		url:     url,
		ownerId: ownerId,
		apiKey:  apiKey,
	}
}

func (c *Client) doRequest(req *http.Request) (*http.Response, error) {
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to do http request: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		responseBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read response body: %w", err)
		}
		return nil, fmt.Errorf("failed to perform request: %s", string(responseBytes))
	}

	return resp, nil
}
