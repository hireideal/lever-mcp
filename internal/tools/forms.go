package tools

import (
	"context"
	"fmt"
	"net/url"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/stefanoamorelli/lever-mcp/internal/client"
)

func listFormsTool() *mcp.Tool {
	return &mcp.Tool{
		Annotations: readOnlyAnnotations(),
		Name:        "list_opportunity_forms",
		Description: "List all profile forms for a specific opportunity in Lever.",
		InputSchema: objectSchema(mergeProperties(paginationProperties(), map[string]any{
			"opportunity_id": prop("string", "The opportunity ID"),
		}), []string{"opportunity_id"}),
	}
}

func listFormsHandler(c client.LeverClient) mcp.ToolHandler {
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
		setPagination(params, args)

		data, err := c.Get(ctx, fmt.Sprintf("/opportunities/%s/forms", id), params)
		if err != nil {
			return toolErrorf("Failed to list forms: %v", err), nil
		}
		return toolText(string(data)), nil
	}
}
