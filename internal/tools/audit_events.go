package tools

import (
	"context"
	"net/url"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/stefanoamorelli/lever-mcp/internal/client"
)

func listAuditEventsTool() *mcp.Tool {
	return &mcp.Tool{
		Annotations: readOnlyAnnotations(),
		Name:        "list_audit_events",
		Description: "List all audit events in your Lever account. Supports pagination.",
		InputSchema: objectSchema(paginationProperties(), nil),
	}
}

func listAuditEventsHandler(c client.LeverClient) mcp.ToolHandler {
	return func(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args, err := extractArgs(req)
		if err != nil {
			return toolError(err.Error()), nil
		}

		params := url.Values{}
		setPagination(params, args)

		data, err := c.Get(ctx, "/audit_events", params)
		if err != nil {
			return toolErrorf("Failed to list audit events: %v", err), nil
		}
		return toolText(string(data)), nil
	}
}
