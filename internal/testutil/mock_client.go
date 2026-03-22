package testutil

import (
	"context"
	"encoding/json"
	"net/url"
)

type MockLeverClient struct {
	GetFunc    func(ctx context.Context, path string, params url.Values) (json.RawMessage, error)
	PostFunc   func(ctx context.Context, path string, params url.Values, body json.RawMessage) (json.RawMessage, error)
	PutFunc    func(ctx context.Context, path string, params url.Values, body json.RawMessage) (json.RawMessage, error)
	DeleteFunc func(ctx context.Context, path string, params url.Values) (json.RawMessage, error)
}

func (m *MockLeverClient) Get(ctx context.Context, path string, params url.Values) (json.RawMessage, error) {
	if m.GetFunc != nil {
		return m.GetFunc(ctx, path, params)
	}
	return json.RawMessage(`{"data":[]}`), nil
}

func (m *MockLeverClient) Post(ctx context.Context, path string, params url.Values, body json.RawMessage) (json.RawMessage, error) {
	if m.PostFunc != nil {
		return m.PostFunc(ctx, path, params, body)
	}
	return json.RawMessage(`{"data":{}}`), nil
}

func (m *MockLeverClient) Put(ctx context.Context, path string, params url.Values, body json.RawMessage) (json.RawMessage, error) {
	if m.PutFunc != nil {
		return m.PutFunc(ctx, path, params, body)
	}
	return json.RawMessage(`{"data":{}}`), nil
}

func (m *MockLeverClient) Delete(ctx context.Context, path string, params url.Values) (json.RawMessage, error) {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, path, params)
	}
	return json.RawMessage(`{}`), nil
}
