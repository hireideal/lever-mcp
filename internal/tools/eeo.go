package tools

import (
	"context"
	"net/url"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/stefanoamorelli/lever-mcp/internal/client"
)

func listEEOResponsesTool() *mcp.Tool {
	return &mcp.Tool{
		Annotations: readOnlyAnnotations(),
		Name:        "list_eeo_responses",
		Description: "List anonymous Equal Employment Opportunity (EEO) responses from Lever. Supports date range filtering.",
		InputSchema: objectSchema(mergeProperties(paginationProperties(), map[string]any{
			"from_date": prop("string", "Start date filter (ISO 8601)"),
			"to_date":   prop("string", "End date filter (ISO 8601)"),
		}), nil),
	}
}

func listEEOResponsesHandler(c client.LeverClient) mcp.ToolHandler {
	return func(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args, err := extractArgs(req)
		if err != nil {
			return toolError(err.Error()), nil
		}

		params := url.Values{}
		setPagination(params, args)
		setString(params, "fromDate", getString(args, "from_date", ""))
		setString(params, "toDate", getString(args, "to_date", ""))

		data, err := c.Get(ctx, "/eeo/responses", params)
		if err != nil {
			return toolErrorf("Failed to list EEO responses: %v", err), nil
		}
		return toolText(string(data)), nil
	}
}
