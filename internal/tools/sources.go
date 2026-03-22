package tools

import (
	"context"
	"net/url"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/stefanoamorelli/lever-mcp/internal/client"
)

func listSourcesTool() *mcp.Tool {
	return &mcp.Tool{
		Annotations: readOnlyAnnotations(),
		Name:        "list_sources",
		Description: "List all sources in your Lever account.",
		InputSchema: objectSchema(paginationProperties(), nil),
	}
}

func listSourcesHandler(c client.LeverClient) mcp.ToolHandler {
	return func(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args, err := extractArgs(req)
		if err != nil {
			return toolError(err.Error()), nil
		}

		params := url.Values{}
		setPagination(params, args)

		data, err := c.Get(ctx, "/sources", params)
		if err != nil {
			return toolErrorf("Failed to list sources: %v", err), nil
		}
		return toolText(string(data)), nil
	}
}
