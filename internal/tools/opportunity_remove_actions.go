package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/stefanoamorelli/lever-mcp/internal/client"
)

func removeOpportunityTagsTool() *mcp.Tool {
	return &mcp.Tool{
		Annotations: updateAnnotations(),
		Name:        "remove_opportunity_tags",
		Description: "Remove tags from an opportunity in Lever.",
		InputSchema: objectSchema(map[string]any{
			"opportunity_id": prop("string", "The opportunity ID"),
			"perform_as":     prop("string", "User ID to perform this action as"),
			"tags":           map[string]any{"type": "array", "items": map[string]any{"type": "string"}, "description": "Tags to remove"},
		}, []string{"opportunity_id", "perform_as", "tags"}),
	}
}

func removeOpportunityTagsHandler(c client.LeverClient) mcp.ToolHandler {
	return func(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args, err := extractArgs(req)
		if err != nil {
			return toolError(err.Error()), nil
		}

		id := getString(args, "opportunity_id", "")
		performAs := getString(args, "perform_as", "")
		if id == "" || performAs == "" {
			return toolError("opportunity_id and perform_as are required"), nil
		}

		params := url.Values{}
		params.Set("perform_as", performAs)

		tags := getStringSlice(args, "tags")
		body := map[string]any{"tags": tags}
		bodyJSON, _ := json.Marshal(body)

		data, err := c.Post(ctx, fmt.Sprintf("/opportunities/%s/removeTags", id), params, bodyJSON)
		if err != nil {
			return toolErrorf("Failed to remove tags: %v", err), nil
		}
		return toolText(string(data)), nil
	}
}

func removeOpportunitySourcesTool() *mcp.Tool {
	return &mcp.Tool{
		Annotations: updateAnnotations(),
		Name:        "remove_opportunity_sources",
		Description: "Remove sources from an opportunity in Lever.",
		InputSchema: objectSchema(map[string]any{
			"opportunity_id": prop("string", "The opportunity ID"),
			"perform_as":     prop("string", "User ID to perform this action as"),
			"sources":        map[string]any{"type": "array", "items": map[string]any{"type": "string"}, "description": "Sources to remove"},
		}, []string{"opportunity_id", "perform_as", "sources"}),
	}
}

func removeOpportunitySourcesHandler(c client.LeverClient) mcp.ToolHandler {
	return func(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args, err := extractArgs(req)
		if err != nil {
			return toolError(err.Error()), nil
		}

		id := getString(args, "opportunity_id", "")
		performAs := getString(args, "perform_as", "")
		if id == "" || performAs == "" {
			return toolError("opportunity_id and perform_as are required"), nil
		}

		params := url.Values{}
		params.Set("perform_as", performAs)

		sources := getStringSlice(args, "sources")
		body := map[string]any{"sources": sources}
		bodyJSON, _ := json.Marshal(body)

		data, err := c.Post(ctx, fmt.Sprintf("/opportunities/%s/removeSources", id), params, bodyJSON)
		if err != nil {
			return toolErrorf("Failed to remove sources: %v", err), nil
		}
		return toolText(string(data)), nil
	}
}

func removeOpportunityLinksTool() *mcp.Tool {
	return &mcp.Tool{
		Annotations: updateAnnotations(),
		Name:        "remove_opportunity_links",
		Description: "Remove links from the contact associated with an opportunity in Lever.",
		InputSchema: objectSchema(map[string]any{
			"opportunity_id": prop("string", "The opportunity ID"),
			"perform_as":     prop("string", "User ID to perform this action as"),
			"links":          map[string]any{"type": "array", "items": map[string]any{"type": "string"}, "description": "Links to remove"},
		}, []string{"opportunity_id", "perform_as", "links"}),
	}
}

func removeOpportunityLinksHandler(c client.LeverClient) mcp.ToolHandler {
	return func(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args, err := extractArgs(req)
		if err != nil {
			return toolError(err.Error()), nil
		}

		id := getString(args, "opportunity_id", "")
		performAs := getString(args, "perform_as", "")
		if id == "" || performAs == "" {
			return toolError("opportunity_id and perform_as are required"), nil
		}

		params := url.Values{}
		params.Set("perform_as", performAs)

		links := getStringSlice(args, "links")
		body := map[string]any{"links": links}
		bodyJSON, _ := json.Marshal(body)

		data, err := c.Post(ctx, fmt.Sprintf("/opportunities/%s/removeLinks", id), params, bodyJSON)
		if err != nil {
			return toolErrorf("Failed to remove links: %v", err), nil
		}
		return toolText(string(data)), nil
	}
}
