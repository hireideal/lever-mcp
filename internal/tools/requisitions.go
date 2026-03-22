package tools

import (
	"context"
	"fmt"
	"net/url"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/stefanoamorelli/lever-mcp/internal/client"
)

func listRequisitionsTool() *mcp.Tool {
	return &mcp.Tool{
		Annotations: readOnlyAnnotations(),
		Name:        "list_requisitions",
		Description: "List all requisitions in your Lever account. Supports filtering by status, requisition code, and date ranges.",
		InputSchema: objectSchema(mergeProperties(paginationProperties(), map[string]any{
			"status":           prop("string", "Filter by status: open, onHold, closed, draft"),
			"requisition_code": prop("string", "Filter by requisition code"),
			"created_at_start": prop("integer", "Filter by creation date start (Unix timestamp in ms)"),
			"created_at_end":   prop("integer", "Filter by creation date end (Unix timestamp in ms)"),
			"confidentiality":  prop("string", "Filter: all, non-confidential, or confidential"),
		}), nil),
	}
}

func listRequisitionsHandler(c client.LeverClient) mcp.ToolHandler {
	return func(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args, err := extractArgs(req)
		if err != nil {
			return toolError(err.Error()), nil
		}

		params := url.Values{}
		setPagination(params, args)
		setString(params, "status", getString(args, "status", ""))
		setString(params, "requisition_code", getString(args, "requisition_code", ""))
		setInt(params, "created_at_start", args)
		setInt(params, "created_at_end", args)
		setString(params, "confidentiality", getString(args, "confidentiality", ""))

		data, err := c.Get(ctx, "/requisitions", params)
		if err != nil {
			return toolErrorf("Failed to list requisitions: %v", err), nil
		}
		return toolText(string(data)), nil
	}
}

func getRequisitionTool() *mcp.Tool {
	return &mcp.Tool{
		Annotations: readOnlyAnnotations(),
		Name:        "get_requisition",
		Description: "Retrieve a single requisition by ID from Lever.",
		InputSchema: objectSchema(map[string]any{
			"requisition_id": prop("string", "The requisition ID"),
		}, []string{"requisition_id"}),
	}
}

func getRequisitionHandler(c client.LeverClient) mcp.ToolHandler {
	return func(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args, err := extractArgs(req)
		if err != nil {
			return toolError(err.Error()), nil
		}

		id := getString(args, "requisition_id", "")
		if id == "" {
			return toolError("requisition_id is required"), nil
		}

		data, err := c.Get(ctx, fmt.Sprintf("/requisitions/%s", id), nil)
		if err != nil {
			return toolErrorf("Failed to get requisition: %v", err), nil
		}
		return toolText(string(data)), nil
	}
}
