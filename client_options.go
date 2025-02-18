package line

import (
	"net/http"
)

type (
	ClientOptionFunc func(client *Client) error
)

// WithBaseURL sets the base URL for API requests to a custom endpoint.
func WithBaseURL(urlStr string) ClientOptionFunc {
	return func(c *Client) error {
		return c.setBaseURL(urlStr)
	}
}

// WithToken sets the token for API requests to a custom endpoint.
func WithToken(token string) ClientOptionFunc {
	return func(c *Client) error {
		c.token = token
		return nil
	}
}

func WithUserAgent(userAgent string) ClientOptionFunc {
	return func(c *Client) error {
		c.UserAgent = userAgent
		return nil
	}
}

func WithClient(httpClient *http.Client) ClientOptionFunc {
	return func(c *Client) error {
		c.client = httpClient
		return nil
	}
}

func WithApiVersionPath(apiVersionPath string) ClientOptionFunc {
	return func(c *Client) error {
		c.apiVersionPath = apiVersionPath
		return nil
	}
}
