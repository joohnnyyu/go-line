package line

import (
	"net/http"
	"time"

	"github.com/hashicorp/go-retryablehttp"
)

// RetryableHTTPClientOption defines a function type that takes a pointer
// to a retryablehttp.Client and modifies its configuration.
type RetryableHTTPClientOption func(client *retryablehttp.Client)

// WithRetryableHTTPClientLogger allows setting a custom logger for the
func WithRetryableHTTPClientLogger(logger retryablehttp.Logger) RetryableHTTPClientOption {
	return func(client *retryablehttp.Client) {
		client.Logger = logger
	}
}

// WithRetryableHTTPClientRetryWaitMin sets the minimum wait duration
func WithRetryableHTTPClientRetryWaitMin(waitMin time.Duration) RetryableHTTPClientOption {
	return func(client *retryablehttp.Client) {
		client.RetryWaitMin = waitMin
	}
}

// WithRetryableHTTPClientRetryWaitMax sets the maximum wait duration
func WithRetryableHTTPClientRetryWaitMax(waitMax time.Duration) RetryableHTTPClientOption {
	return func(client *retryablehttp.Client) {
		client.RetryWaitMax = waitMax
	}
}

// WithRetryableHTTPClientRetryMax sets the maximum number of retry attempts
func WithRetryableHTTPClientRetryMax(retryMax int) RetryableHTTPClientOption {
	return func(client *retryablehttp.Client) {
		client.RetryMax = retryMax
	}
}

// WithRetryableHTTPClientCheckRetry allows setting a custom check retry
func WithRetryableHTTPClientCheckRetry(checkRetry retryablehttp.CheckRetry) RetryableHTTPClientOption {
	return func(client *retryablehttp.Client) {
		client.CheckRetry = checkRetry
	}
}

// WithRetryableHTTPClientBackoff allows setting a custom backoff strategy
func WithRetryableHTTPClientBackoff(backoff retryablehttp.Backoff) RetryableHTTPClientOption {
	return func(client *retryablehttp.Client) {
		client.Backoff = backoff
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
