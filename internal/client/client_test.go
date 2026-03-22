package client

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNew_DefaultBaseURL(t *testing.T) {
	c := New("test-key").(*httpLeverClient)
	if c.baseURL != defaultBaseURL {
		t.Errorf("expected %s, got %s", defaultBaseURL, c.baseURL)
	}
}

func TestNew_WithBaseURL(t *testing.T) {
	c := New("test-key", WithBaseURL("https://custom.api")).(*httpLeverClient)
	if c.baseURL != "https://custom.api" {
		t.Errorf("expected https://custom.api, got %s", c.baseURL)
	}
}

func TestGet(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}

		user, _, ok := r.BasicAuth()
		if !ok || user != "test-key" {
			t.Errorf("expected basic auth with test-key, got %s", user)
		}

		if r.URL.Path != "/v1/opportunities" {
			t.Errorf("expected /v1/opportunities, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{"data": []any{}})
	}))
	defer server.Close()

	c := New("test-key", WithBaseURL(server.URL+"/v1"))
	_, err := c.Get(context.Background(), "/opportunities", nil)
	if err != nil {
		t.Fatal(err)
	}
}

func TestPost(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.Header.Get("Content-Type") != "application/json" {
			t.Error("expected application/json content type")
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]any{"data": map[string]any{"id": "123"}})
	}))
	defer server.Close()

	c := New("test-key", WithBaseURL(server.URL+"/v1"))
	_, err := c.Post(context.Background(), "/opportunities", nil, json.RawMessage(`{"name":"test"}`))
	if err != nil {
		t.Fatal(err)
	}
}

func TestPut(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		json.NewEncoder(w).Encode(map[string]any{"data": map[string]any{}})
	}))
	defer server.Close()

	c := New("test-key", WithBaseURL(server.URL+"/v1"))
	_, err := c.Put(context.Background(), "/opportunities/123/stage", nil, json.RawMessage(`{"stage":"456"}`))
	if err != nil {
		t.Fatal(err)
	}
}

func TestDelete(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		json.NewEncoder(w).Encode(map[string]any{})
	}))
	defer server.Close()

	c := New("test-key", WithBaseURL(server.URL+"/v1"))
	_, err := c.Delete(context.Background(), "/webhooks/123", nil)
	if err != nil {
		t.Fatal(err)
	}
}

func TestGet_ErrorStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"error":"unauthorized"}`))
	}))
	defer server.Close()

	c := New("bad-key", WithBaseURL(server.URL+"/v1"))
	_, err := c.Get(context.Background(), "/opportunities", nil)
	if err == nil {
		t.Fatal("expected error for 401 response")
	}
}

func TestGet_WithParams(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("limit") != "10" {
			t.Errorf("expected limit=10, got %s", r.URL.Query().Get("limit"))
		}
		json.NewEncoder(w).Encode(map[string]any{"data": []any{}})
	}))
	defer server.Close()

	c := New("test-key", WithBaseURL(server.URL+"/v1"))
	params := map[string][]string{"limit": {"10"}}
	_, err := c.Get(context.Background(), "/opportunities", params)
	if err != nil {
		t.Fatal(err)
	}
}
