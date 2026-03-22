package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

const defaultBaseURL = "https://api.lever.co/v1"

type LeverClient interface {
	Get(ctx context.Context, path string, params url.Values) (json.RawMessage, error)
	Post(ctx context.Context, path string, params url.Values, body json.RawMessage) (json.RawMessage, error)
	Put(ctx context.Context, path string, params url.Values, body json.RawMessage) (json.RawMessage, error)
	Delete(ctx context.Context, path string, params url.Values) (json.RawMessage, error)
}

type Option func(*httpLeverClient)

func WithBaseURL(baseURL string) Option {
	return func(c *httpLeverClient) { c.baseURL = baseURL }
}

func WithHTTPClient(hc *http.Client) Option {
	return func(c *httpLeverClient) { c.httpClient = hc }
}

func New(apiKey string, opts ...Option) LeverClient {
	c := &httpLeverClient{
		apiKey:     apiKey,
		baseURL:    defaultBaseURL,
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

type httpLeverClient struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
}

func (c *httpLeverClient) do(ctx context.Context, method, path string, params url.Values, body io.Reader) (json.RawMessage, error) {
	reqURL := c.baseURL + path
	if len(params) > 0 {
		reqURL += "?" + params.Encode()
	}

	req, err := http.NewRequestWithContext(ctx, method, reqURL, body)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	req.SetBasicAuth(c.apiKey, "")
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(respBody))
	}

	return json.RawMessage(respBody), nil
}

func (c *httpLeverClient) Get(ctx context.Context, path string, params url.Values) (json.RawMessage, error) {
	return c.do(ctx, http.MethodGet, path, params, nil)
}

func (c *httpLeverClient) Post(ctx context.Context, path string, params url.Values, body json.RawMessage) (json.RawMessage, error) {
	var r io.Reader
	if body != nil {
		r = bytes.NewReader(body)
	}
	return c.do(ctx, http.MethodPost, path, params, r)
}

func (c *httpLeverClient) Put(ctx context.Context, path string, params url.Values, body json.RawMessage) (json.RawMessage, error) {
	var r io.Reader
	if body != nil {
		r = bytes.NewReader(body)
	}
	return c.do(ctx, http.MethodPut, path, params, r)
}

func (c *httpLeverClient) Delete(ctx context.Context, path string, params url.Values) (json.RawMessage, error) {
	return c.do(ctx, http.MethodDelete, path, params, nil)
}
