package line

import (
	"context"
	"net/http"

	"github.com/hashicorp/go-retryablehttp"
)

type (
	RetryableHTTPClientOption func(client *retryablehttp.Client)
	RequestOptionFunc         func(*http.Request) error
)

// WithContext runs the request with the provided context
func WithContext(ctx context.Context) RequestOptionFunc {
	return func(req *http.Request) error {
		*req = *req.WithContext(ctx)
		return nil
	}
}

// WithHeader takes a header name and value and appends it to the request headers.
func WithHeader(name, value string) RequestOptionFunc {
	return func(req *http.Request) error {
		req.Header.Set(name, value)
		return nil
	}
}

// WithHeaders takes a map of header name/value pairs and appends them to the
func WithHeaders(headers map[string]string) RequestOptionFunc {
	return func(req *http.Request) error {
		for k, v := range headers {
			req.Header.Set(k, v)
		}
		return nil
	}
}

func NewRetryableHTTPClient(opts ...RetryableHTTPClientOption) *http.Client {
	retryClient := retryablehttp.NewClient()
	retryClient.Logger = nil
	for _, opt := range opts {
		opt(retryClient)
	}
	return retryClient.StandardClient()
}
