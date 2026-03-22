package tools

import (
	"context"
	"fmt"
	"net/url"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/stefanoamorelli/lever-mcp/internal/client"
)

func listArchiveReasonsTool() *mcp.Tool {
	return &mcp.Tool{
		Annotations: readOnlyAnnotations(),
		Name:        "list_archive_reasons",
		Description: "List all archive reasons in your Lever account. Optionally filter by type (hired or non-hired).",
		InputSchema: objectSchema(mergeProperties(paginationProperties(), map[string]any{
			"type": prop("string", "Filter by type: hired or non-hired"),
		}), nil),
	}
}

func listArchiveReasonsHandler(c client.LeverClient) mcp.ToolHandler {
	return func(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args, err := extractArgs(req)
		if err != nil {
			return toolError(err.Error()), nil
		}

		params := url.Values{}
		setPagination(params, args)
		setString(params, "type", getString(args, "type", ""))

		data, err := c.Get(ctx, "/archive_reasons", params)
		if err != nil {
			return toolErrorf("Failed to list archive reasons: %v", err), nil
		}
		return toolText(string(data)), nil
	}
}

func getArchiveReasonTool() *mcp.Tool {
	return &mcp.Tool{
		Annotations: readOnlyAnnotations(),
		Name:        "get_archive_reason",
		Description: "Retrieve a single archive reason by ID from Lever.",
		InputSchema: objectSchema(map[string]any{
			"archive_reason_id": prop("string", "The archive reason ID"),
		}, []string{"archive_reason_id"}),
	}
}

func getArchiveReasonHandler(c client.LeverClient) mcp.ToolHandler {
	return func(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args, err := extractArgs(req)
		if err != nil {
			return toolError(err.Error()), nil
		}

		id := getString(args, "archive_reason_id", "")
		if id == "" {
			return toolError("archive_reason_id is required"), nil
		}

		data, err := c.Get(ctx, fmt.Sprintf("/archive_reasons/%s", id), nil)
		if err != nil {
			return toolErrorf("Failed to get archive reason: %v", err), nil
		}
		return toolText(string(data)), nil
	}
}
