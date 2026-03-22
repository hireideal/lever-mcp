package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/stefanoamorelli/lever-mcp/internal/client"
)

func listOpportunitiesTool() *mcp.Tool {
	return &mcp.Tool{
		Annotations: readOnlyAnnotations(),
		Name:        "list_opportunities",
		Description: "List all opportunities in your Lever account. Supports filtering by posting, stage, contact, tag, source, origin, archived status, confidentiality, and date ranges. Results are paginated.",
		InputSchema: objectSchema(mergeProperties(paginationProperties(), map[string]any{
			"posting_id":       prop("string", "Filter by posting ID"),
			"stage_id":         prop("string", "Filter by stage ID"),
			"contact_id":       prop("string", "Filter by contact ID"),
			"tag":              prop("string", "Filter by tag"),
			"source":           prop("string", "Filter by source"),
			"origin":           prop("string", "Filter by origin"),
			"confidentiality":  prop("string", "Filter: all, non-confidential, or confidential"),
			"archived":         prop("boolean", "Filter by archived status"),
			"expand":           prop("string", "Comma-separated fields to expand (e.g. applications,stage,owner,followers,sourcedBy,contact)"),
			"created_at_start": prop("integer", "Filter by creation date start (Unix timestamp in ms)"),
			"created_at_end":   prop("integer", "Filter by creation date end (Unix timestamp in ms)"),
			"updated_at_start": prop("integer", "Filter by update date start (Unix timestamp in ms)"),
			"updated_at_end":   prop("integer", "Filter by update date end (Unix timestamp in ms)"),
		}), nil),
	}
}

func listOpportunitiesHandler(c client.LeverClient) mcp.ToolHandler {
	return func(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args, err := extractArgs(req)
		if err != nil {
			return toolError(err.Error()), nil
		}

		params := url.Values{}
		setPagination(params, args)
		setString(params, "posting_id", getString(args, "posting_id", ""))
		setString(params, "stage_id", getString(args, "stage_id", ""))
		setString(params, "contact_id", getString(args, "contact_id", ""))
		setString(params, "tag", getString(args, "tag", ""))
		setString(params, "source", getString(args, "source", ""))
		setString(params, "origin", getString(args, "origin", ""))
		setString(params, "confidentiality", getString(args, "confidentiality", ""))
		setBool(params, "archived", args)
		setString(params, "expand", getString(args, "expand", ""))
		setInt(params, "created_at_start", args)
		setInt(params, "created_at_end", args)
		setInt(params, "updated_at_start", args)
		setInt(params, "updated_at_end", args)

		data, err := c.Get(ctx, "/opportunities", params)
		if err != nil {
			return toolErrorf("Failed to list opportunities: %v", err), nil
		}
		return toolText(string(data)), nil
	}
}

func getOpportunityTool() *mcp.Tool {
	return &mcp.Tool{
		Annotations: readOnlyAnnotations(),
		Name:        "get_opportunity",
		Description: "Retrieve a single opportunity by ID from Lever. Optionally expand nested objects.",
		InputSchema: objectSchema(map[string]any{
			"opportunity_id": prop("string", "The opportunity ID"),
			"expand":         prop("string", "Comma-separated fields to expand"),
		}, []string{"opportunity_id"}),
	}
}

func getOpportunityHandler(c client.LeverClient) mcp.ToolHandler {
	return func(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args, err := extractArgs(req)
		if err != nil {
			return toolError(err.Error()), nil
		}

		id := getString(args, "opportunity_id", "")
		if id == "" {
			return toolError("opportunity_id is required"), nil
		}

		params := url.Values{}
		setString(params, "expand", getString(args, "expand", ""))

		data, err := c.Get(ctx, fmt.Sprintf("/opportunities/%s", id), params)
		if err != nil {
			return toolErrorf("Failed to get opportunity: %v", err), nil
		}
		return toolText(string(data)), nil
	}
}

func createOpportunityTool() *mcp.Tool {
	return &mcp.Tool{
		Annotations: createAnnotations(),
		Name:        "create_opportunity",
		Description: "Create a new opportunity (and optionally a new candidate) in Lever. Deduplicates candidates by email. Provide perform_as (user ID) to attribute the action.",
		InputSchema: objectSchema(map[string]any{
			"perform_as": prop("string", "User ID to perform this action as (required)"),
			"name":       prop("string", "Candidate full name"),
			"headline":   prop("string", "Candidate headline"),
			"location":   prop("string", "Candidate location"),
			"emails":     map[string]any{"type": "array", "items": map[string]any{"type": "string"}, "description": "Candidate email addresses"},
			"phones":     map[string]any{"type": "array", "items": map[string]any{"type": "object"}, "description": "Candidate phone numbers [{type, value}]"},
			"links":      map[string]any{"type": "array", "items": map[string]any{"type": "string"}, "description": "Candidate links (LinkedIn, portfolio, etc.)"},
			"tags":       map[string]any{"type": "array", "items": map[string]any{"type": "string"}, "description": "Tags to apply"},
			"sources":    map[string]any{"type": "array", "items": map[string]any{"type": "string"}, "description": "Sources to apply"},
			"origin":     prop("string", "Origin of the opportunity (e.g. sourced, applied, referred)"),
			"posting_id": prop("string", "Posting ID to associate with"),
			"stage_id":   prop("string", "Stage ID to place the opportunity in"),
		}, []string{"perform_as"}),
	}
}

func createOpportunityHandler(c client.LeverClient) mcp.ToolHandler {
	return func(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args, err := extractArgs(req)
		if err != nil {
			return toolError(err.Error()), nil
		}

		performAs := getString(args, "perform_as", "")
		if performAs == "" {
			return toolError("perform_as is required"), nil
		}

		params := url.Values{}
		params.Set("perform_as", performAs)

		body := make(map[string]any)
		for _, key := range []string{"name", "headline", "location", "origin", "posting_id", "stage_id"} {
			if v := getString(args, key, ""); v != "" {
				body[key] = v
			}
		}
		for _, key := range []string{"emails", "phones", "links", "tags", "sources"} {
			if v, ok := args[key]; ok {
				body[key] = v
			}
		}

		bodyJSON, _ := json.Marshal(body)
		data, err := c.Post(ctx, "/opportunities", params, bodyJSON)
		if err != nil {
			return toolErrorf("Failed to create opportunity: %v", err), nil
		}
		return toolText(string(data)), nil
	}
}

func listDeletedOpportunitiesTool() *mcp.Tool {
	return &mcp.Tool{
		Annotations: readOnlyAnnotations(),
		Name:        "list_deleted_opportunities",
		Description: "List all deleted opportunities in your Lever account.",
		InputSchema: objectSchema(paginationProperties(), nil),
	}
}

func listDeletedOpportunitiesHandler(c client.LeverClient) mcp.ToolHandler {
	return func(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args, err := extractArgs(req)
		if err != nil {
			return toolError(err.Error()), nil
		}

		params := url.Values{}
		setPagination(params, args)

		data, err := c.Get(ctx, "/opportunities/deleted", params)
		if err != nil {
			return toolErrorf("Failed to list deleted opportunities: %v", err), nil
		}
		return toolText(string(data)), nil
	}
}
