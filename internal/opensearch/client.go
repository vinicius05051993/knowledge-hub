package opensearch

import (
	"fmt"
	"net/http"

	"indexer/internal/config"
)

type Client struct {
	baseURL string

	httpClient *http.Client
}

func NewClient(
	cfg *config.Config,
) *Client {

	return &Client{
		baseURL: fmt.Sprintf(
			"http://%s:%s",
			cfg.OpenSearchHost,
			cfg.OpenSearchPort,
		),

		httpClient: &http.Client{},
	}
}

func (c *Client) Do(
	req *http.Request,
) (*http.Response, error) {

	req.Header.Set(
		"Content-Type",
		"application/json",
	)

	return c.httpClient.Do(req)
}

func (c *Client) URL(
	path string,
) string {

	return fmt.Sprintf(
		"%s/%s",
		c.baseURL,
		path,
	)
}