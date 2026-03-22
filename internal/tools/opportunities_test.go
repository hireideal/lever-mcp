package tools

import (
	"context"
	"encoding/json"
	"net/url"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/stefanoamorelli/lever-mcp/internal/testutil"
)

func makeReq(args map[string]any) *mcp.CallToolRequest {
	data, _ := json.Marshal(args)
	return &mcp.CallToolRequest{
		Params: &mcp.CallToolParamsRaw{
			Arguments: json.RawMessage(data),
		},
	}
}

func TestListOpportunitiesHandler(t *testing.T) {
	mock := &testutil.MockLeverClient{
		GetFunc: func(ctx context.Context, path string, params url.Values) (json.RawMessage, error) {
			if path != "/opportunities" {
				t.Errorf("expected path /opportunities, got %s", path)
			}
			return testutil.OpportunityListFixture(), nil
		},
	}

	handler := listOpportunitiesHandler(mock)
	result, err := handler(context.Background(), makeReq(map[string]any{}))
	if err != nil {
		t.Fatal(err)
	}
	if result.IsError {
		t.Fatal("expected success")
	}
}

func TestListOpportunitiesHandler_WithFilters(t *testing.T) {
	mock := &testutil.MockLeverClient{
		GetFunc: func(ctx context.Context, path string, params url.Values) (json.RawMessage, error) {
			if params.Get("posting_id") != "post-123" {
				t.Errorf("expected posting_id=post-123, got %s", params.Get("posting_id"))
			}
			if params.Get("tag") != "engineering" {
				t.Errorf("expected tag=engineering, got %s", params.Get("tag"))
			}
			return testutil.OpportunityListFixture(), nil
		},
	}

	handler := listOpportunitiesHandler(mock)
	result, err := handler(context.Background(), makeReq(map[string]any{
		"posting_id": "post-123",
		"tag":        "engineering",
	}))
	if err != nil {
		t.Fatal(err)
	}
	if result.IsError {
		t.Fatal("expected success")
	}
}

func TestGetOpportunityHandler(t *testing.T) {
	mock := &testutil.MockLeverClient{
		GetFunc: func(ctx context.Context, path string, params url.Values) (json.RawMessage, error) {
			if path != "/opportunities/opp-123" {
				t.Errorf("expected path /opportunities/opp-123, got %s", path)
			}
			return testutil.OpportunityFixture(), nil
		},
	}

	handler := getOpportunityHandler(mock)
	result, err := handler(context.Background(), makeReq(map[string]any{"opportunity_id": "opp-123"}))
	if err != nil {
		t.Fatal(err)
	}
	if result.IsError {
		t.Fatal("expected success")
	}
}

func TestGetOpportunityHandler_MissingID(t *testing.T) {
	mock := &testutil.MockLeverClient{}
	handler := getOpportunityHandler(mock)
	result, err := handler(context.Background(), makeReq(map[string]any{}))
	if err != nil {
		t.Fatal(err)
	}
	if !result.IsError {
		t.Fatal("expected error when ID missing")
	}
}

func TestCreateOpportunityHandler(t *testing.T) {
	mock := &testutil.MockLeverClient{
		PostFunc: func(ctx context.Context, path string, params url.Values, body json.RawMessage) (json.RawMessage, error) {
			if path != "/opportunities" {
				t.Errorf("expected path /opportunities, got %s", path)
			}
			if params.Get("perform_as") != "user-123" {
				t.Errorf("expected perform_as=user-123, got %s", params.Get("perform_as"))
			}
			var b map[string]any
			json.Unmarshal(body, &b)
			if b["name"] != "Jane Doe" {
				t.Errorf("expected name=Jane Doe, got %v", b["name"])
			}
			return testutil.OpportunityFixture(), nil
		},
	}

	handler := createOpportunityHandler(mock)
	result, err := handler(context.Background(), makeReq(map[string]any{
		"perform_as": "user-123",
		"name":       "Jane Doe",
		"emails":     []string{"jane@example.com"},
	}))
	if err != nil {
		t.Fatal(err)
	}
	if result.IsError {
		t.Fatal("expected success")
	}
}

func TestCreateOpportunityHandler_MissingPerformAs(t *testing.T) {
	mock := &testutil.MockLeverClient{}
	handler := createOpportunityHandler(mock)
	result, err := handler(context.Background(), makeReq(map[string]any{"name": "Jane Doe"}))
	if err != nil {
		t.Fatal(err)
	}
	if !result.IsError {
		t.Fatal("expected error when perform_as missing")
	}
}
