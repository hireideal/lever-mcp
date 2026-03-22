package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/stefanoamorelli/lever-mcp/internal/client"
)

func archiveOpportunityTool() *mcp.Tool {
	return &mcp.Tool{
		Annotations: updateAnnotations(),
		Name:        "archive_opportunity",
		Description: "Archive or unarchive an opportunity in Lever. Provide an archive reason ID to archive, or set archived to false to unarchive.",
		InputSchema: objectSchema(map[string]any{
			"opportunity_id":   prop("string", "The opportunity ID"),
			"perform_as":       prop("string", "User ID to perform this action as"),
			"reason_id":        prop("string", "Archive reason ID (required when archiving)"),
			"clean_up_posting": prop("boolean", "Close posting if no remaining active opportunities"),
		}, []string{"opportunity_id", "perform_as"}),
	}
}

func archiveOpportunityHandler(c client.LeverClient) mcp.ToolHandler {
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

		body := make(map[string]any)
		if reasonID := getString(args, "reason_id", ""); reasonID != "" {
			body["reason"] = reasonID
		}
		if v, ok := getBool(args, "clean_up_posting"); ok {
			body["cleanUpPosting"] = v
		}

		bodyJSON, _ := json.Marshal(body)
		data, err := c.Put(ctx, fmt.Sprintf("/opportunities/%s/archived", id), params, bodyJSON)
		if err != nil {
			return toolErrorf("Failed to archive opportunity: %v", err), nil
		}
		return toolText(string(data)), nil
	}
}

func changeOpportunityStageTool() *mcp.Tool {
	return &mcp.Tool{
		Annotations: updateAnnotations(),
		Name:        "change_opportunity_stage",
		Description: "Move an opportunity to a different pipeline stage in Lever.",
		InputSchema: objectSchema(map[string]any{
			"opportunity_id": prop("string", "The opportunity ID"),
			"perform_as":     prop("string", "User ID to perform this action as"),
			"stage_id":       prop("string", "Target stage ID"),
		}, []string{"opportunity_id", "perform_as", "stage_id"}),
	}
}

func changeOpportunityStageHandler(c client.LeverClient) mcp.ToolHandler {
	return func(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args, err := extractArgs(req)
		if err != nil {
			return toolError(err.Error()), nil
		}

		id := getString(args, "opportunity_id", "")
		performAs := getString(args, "perform_as", "")
		stageID := getString(args, "stage_id", "")
		if id == "" || performAs == "" || stageID == "" {
			return toolError("opportunity_id, perform_as, and stage_id are required"), nil
		}

		params := url.Values{}
		params.Set("perform_as", performAs)

		body := map[string]any{"stage": stageID}
		bodyJSON, _ := json.Marshal(body)

		data, err := c.Put(ctx, fmt.Sprintf("/opportunities/%s/stage", id), params, bodyJSON)
		if err != nil {
			return toolErrorf("Failed to change stage: %v", err), nil
		}
		return toolText(string(data)), nil
	}
}

func addOpportunityTagsTool() *mcp.Tool {
	return &mcp.Tool{
		Annotations: createAnnotations(),
		Name:        "add_opportunity_tags",
		Description: "Add tags to an opportunity in Lever.",
		InputSchema: objectSchema(map[string]any{
			"opportunity_id": prop("string", "The opportunity ID"),
			"perform_as":     prop("string", "User ID to perform this action as"),
			"tags":           map[string]any{"type": "array", "items": map[string]any{"type": "string"}, "description": "Tags to add"},
		}, []string{"opportunity_id", "perform_as", "tags"}),
	}
}

func addOpportunityTagsHandler(c client.LeverClient) mcp.ToolHandler {
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

		data, err := c.Post(ctx, fmt.Sprintf("/opportunities/%s/addTags", id), params, bodyJSON)
		if err != nil {
			return toolErrorf("Failed to add tags: %v", err), nil
		}
		return toolText(string(data)), nil
	}
}

func addOpportunitySourcesTool() *mcp.Tool {
	return &mcp.Tool{
		Annotations: createAnnotations(),
		Name:        "add_opportunity_sources",
		Description: "Add sources to an opportunity in Lever.",
		InputSchema: objectSchema(map[string]any{
			"opportunity_id": prop("string", "The opportunity ID"),
			"perform_as":     prop("string", "User ID to perform this action as"),
			"sources":        map[string]any{"type": "array", "items": map[string]any{"type": "string"}, "description": "Sources to add"},
		}, []string{"opportunity_id", "perform_as", "sources"}),
	}
}

func addOpportunitySourcesHandler(c client.LeverClient) mcp.ToolHandler {
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

		data, err := c.Post(ctx, fmt.Sprintf("/opportunities/%s/addSources", id), params, bodyJSON)
		if err != nil {
			return toolErrorf("Failed to add sources: %v", err), nil
		}
		return toolText(string(data)), nil
	}
}

func addOpportunityLinksTool() *mcp.Tool {
	return &mcp.Tool{
		Annotations: createAnnotations(),
		Name:        "add_opportunity_links",
		Description: "Add links (LinkedIn, portfolio, etc.) to the contact associated with an opportunity in Lever.",
		InputSchema: objectSchema(map[string]any{
			"opportunity_id": prop("string", "The opportunity ID"),
			"perform_as":     prop("string", "User ID to perform this action as"),
			"links":          map[string]any{"type": "array", "items": map[string]any{"type": "string"}, "description": "Links to add"},
		}, []string{"opportunity_id", "perform_as", "links"}),
	}
}

func addOpportunityLinksHandler(c client.LeverClient) mcp.ToolHandler {
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

		data, err := c.Post(ctx, fmt.Sprintf("/opportunities/%s/addLinks", id), params, bodyJSON)
		if err != nil {
			return toolErrorf("Failed to add links: %v", err), nil
		}
		return toolText(string(data)), nil
	}
}
