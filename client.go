package line

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strings"

	"github.com/google/go-querystring/query"
)

const (
	defaultBaseURL = "https://api.line.me/"
	apiVersionPath = "v2/"
	userAgent      = "go-line"

	contentType = "application/json"
)

type AuthType int

const (
	BasicAuth AuthType = iota
)

var ErrNotFound = errors.New("404 Not Found")

type Client struct {
	client                *http.Client
	authType              AuthType
	baseURL               *url.URL
	apiVersionPath        string
	defaultRequestOptions []RequestOptionFunc
	token                 string
	UserAgent             string
	ContentType           string
	Bot                   *BotService
	Message               *MessageService
}

type Response struct {
	*http.Response
}

func NewClient(token string, options ...ClientOptionFunc) (*Client, error) {
	client, err := newClient(append(
		options,
		WithUserAgent(userAgent),
		WithApiVersionPath(apiVersionPath),
		WithToken(token))...,
	)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func newClient(options ...ClientOptionFunc) (*Client, error) {
	c := &Client{
		client:      NewRetryableHTTPClient(),
		ContentType: contentType,
	}

	err := c.setBaseURL(defaultBaseURL)
	if err != nil {
		return nil, err
	}

	for _, fn := range options {
		if fn == nil {
			continue
		}
		if err := fn(c); err != nil {
			return nil, err
		}
	}

	c.Bot = &BotService{client: c}
	c.Message = &MessageService{client: c}
	return c, nil
}

func (c *Client) setBaseURL(urlStr string) error {
	// Make sure the given URL end with a slash
	if !strings.HasSuffix(urlStr, "/") {
		urlStr += "/"
	}

	baseURL, err := url.Parse(urlStr)
	if err != nil {
		return err
	}

	// Update the base URL of the client.
	c.baseURL = baseURL

	return nil
}

func (c *Client) NewRequest(ctx context.Context, method, path string, opt interface{}, options []RequestOptionFunc) (*http.Request, error) {
	u := *c.baseURL
	unescaped, err := url.PathUnescape(path)
	if err != nil {
		return nil, err
	}

	// Set the encoded path data
	baseURL := c.baseURL.Path
	if !strings.HasSuffix(baseURL, c.apiVersionPath) {
		baseURL += c.apiVersionPath
	}
	u.RawPath = baseURL + path
	u.Path = baseURL + unescaped

	// Create a request specific headers map.
	reqHeaders := make(http.Header)
	reqHeaders.Set("Content-Type", c.ContentType)
	reqHeaders.Set("Accept", "application/json")
	reqHeaders.Set("Authorization", c.token)

	if c.UserAgent != "" {
		reqHeaders.Set("User-Agent", c.UserAgent)
	}

	// Initialize body as an io.Reader
	var body io.Reader
	switch {
	case method == http.MethodPatch || method == http.MethodPost || method == http.MethodPut:
		if opt != nil {
			jsonData, err := json.Marshal(opt)
			if err != nil {
				return nil, err
			}
			body = bytes.NewReader(jsonData)
		}
	case opt != nil:
		q, err := query.Values(opt)
		if err != nil {
			return nil, err
		}
		u.RawQuery = q.Encode()
	}

	req, err := http.NewRequestWithContext(ctx, method, u.String(), body)
	if err != nil {
		return nil, err
	}

	for _, opt := range append(c.defaultRequestOptions, options...) {
		if opt == nil {
			continue
		}
		if err := opt(req); err != nil {
			return nil, err
		}
	}

	// Set the request specific headers.
	for k, v := range reqHeaders {
		req.Header[k] = v
	}

	return req, nil
}

func (c *Client) Do(req *http.Request, v interface{}) (*Response, error) {
	// Set the correct authentication header. If using basic auth, then check
	// if we already have a token and if not first authenticate and get one.
	var basicAuthToken string
	switch c.authType {
	case BasicAuth:
		basicAuthToken = c.token
		req.Header.Set("Authorization", "Bearer "+basicAuthToken)
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	defer func() {
		_, _ = io.Copy(io.Discard, resp.Body)
	}()

	response := newResponse(resp)

	err = CheckResponse(resp)
	if err != nil {
		// Even though there was an error, we still return the response
		// in case the caller wants to inspect it further.
		return response, err
	}

	if v != nil {
		if w, ok := v.(io.Writer); ok {
			_, err = io.Copy(w, resp.Body)
		} else {
			err = json.NewDecoder(resp.Body).Decode(v)
		}
	}

	return response, err
}

func newResponse(r *http.Response) *Response {
	response := &Response{Response: r}
	// response.populatePageValues()
	// response.populateLinkValues()
	return response
}

type ErrorResponse struct {
	Body     []byte
	Response *http.Response
	Message  string
}

func (e *ErrorResponse) Error() string {
	path, _ := url.QueryUnescape(e.Response.Request.URL.Path)
	errorURL := fmt.Sprintf("%s://%s%s", e.Response.Request.URL.Scheme, e.Response.Request.URL.Host, path)

	if e.Message == "" {
		return fmt.Sprintf("%s %s: %d", e.Response.Request.Method, errorURL, e.Response.StatusCode)
	}
	return fmt.Sprintf("%s %s: %d %s", e.Response.Request.Method, errorURL, e.Response.StatusCode, e.Message)
}

func CheckResponse(r *http.Response) error {
	switch r.StatusCode {
	case 200, 201, 202, 204, 304:
		return nil
	case 404:
		return ErrNotFound
	}

	errorResponse := &ErrorResponse{Response: r}

	data, err := io.ReadAll(r.Body)
	if err == nil && strings.TrimSpace(string(data)) != "" {
		errorResponse.Body = data

		var raw interface{}
		if err := json.Unmarshal(data, &raw); err != nil {
			errorResponse.Message = fmt.Sprintf("failed to parse unknown error format: %s", data)
		} else {
			errorResponse.Message = parseError(raw)
		}
	}

	return errorResponse
}

func parseError(raw interface{}) string {
	switch raw := raw.(type) {
	case string:
		return raw

	case []interface{}:
		var errs []string
		for _, v := range raw {
			errs = append(errs, parseError(v))
		}
		return fmt.Sprintf("[%s]", strings.Join(errs, ", "))

	case map[string]interface{}:
		var errs []string
		for k, v := range raw {
			errs = append(errs, fmt.Sprintf("{%s: %s}", k, parseError(v)))
		}
		sort.Strings(errs)
		return strings.Join(errs, ", ")

	default:
		return fmt.Sprintf("failed to parse unexpected error type: %T", raw)
	}
}
