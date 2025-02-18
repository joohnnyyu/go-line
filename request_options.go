package line

import (
	"context"
	"net/http"
)

type RequestOptionFunc func(*http.Request) error

// WithContext runs the request with the provided context
func WithContext(ctx context.Context) RequestOptionFunc {
	return func(req *http.Request) error {
		*req = *req.WithContext(ctx)
		return nil
	}
}

func WithRequestHeader(name, value string) RequestOptionFunc {
	return func(req *http.Request) error {
		req.Header.Set(name, value)
		return nil
	}
}

func WithRequestHeaders(headers map[string]string) RequestOptionFunc {
	return func(req *http.Request) error {
		for k, v := range headers {
			req.Header.Set(k, v)
		}
		return nil
	}
}

func WithRequestHeaderFunc(fn func(http.Header)) RequestOptionFunc {
	return func(req *http.Request) error {
		fn(req.Header)
		return nil
	}
}
