package client

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

// Client -
type Client struct {
	HostURL    string
	HTTPClient *http.Client
}

// NewClient initializes a new API client with the given host
func NewClient(host string) *Client {
	return &Client{
		HTTPClient: &http.Client{Timeout: 10 * time.Second},
		HostURL:    host,
	}
}

func (c *Client) doRequest(req *http.Request) ([]byte, error) {
	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if (res.StatusCode != http.StatusOK) && (res.StatusCode != http.StatusCreated) {
		return nil, fmt.Errorf("status: %d, body: %s", res.StatusCode, body)
	}

	return body, err
}
