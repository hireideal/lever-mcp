package tools

import (
	"context"
	"fmt"
	"net/url"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/stefanoamorelli/lever-mcp/internal/client"
)

func listApplicationsTool() *mcp.Tool {
	return &mcp.Tool{
		Annotations: readOnlyAnnotations(),
		Name:        "list_opportunity_applications",
		Description: "List all applications for a specific opportunity in Lever.",
		InputSchema: objectSchema(mergeProperties(paginationProperties(), map[string]any{
			"opportunity_id": prop("string", "The opportunity ID"),
			"expand":         prop("string", "Comma-separated fields to expand"),
		}), []string{"opportunity_id"}),
	}
}

func listApplicationsHandler(c client.LeverClient) mcp.ToolHandler {
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
		setString(params, "expand", getString(args, "expand", ""))

		data, err := c.Get(ctx, fmt.Sprintf("/opportunities/%s/applications", id), params)
		if err != nil {
			return toolErrorf("Failed to list applications: %v", err), nil
		}
		return toolText(string(data)), nil
	}
}

func getApplicationTool() *mcp.Tool {
	return &mcp.Tool{
		Annotations: readOnlyAnnotations(),
		Name:        "get_opportunity_application",
		Description: "Retrieve a single application for an opportunity in Lever.",
		InputSchema: objectSchema(map[string]any{
			"opportunity_id": prop("string", "The opportunity ID"),
			"application_id": prop("string", "The application ID"),
		}, []string{"opportunity_id", "application_id"}),
	}
}

func getApplicationHandler(c client.LeverClient) mcp.ToolHandler {
	return func(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args, err := extractArgs(req)
		if err != nil {
			return toolError(err.Error()), nil
		}

		id := getString(args, "opportunity_id", "")
		appID := getString(args, "application_id", "")
		if id == "" || appID == "" {
			return toolError("opportunity_id and application_id are required"), nil
		}

		data, err := c.Get(ctx, fmt.Sprintf("/opportunities/%s/applications/%s", id, appID), nil)
		if err != nil {
			return toolErrorf("Failed to get application: %v", err), nil
		}
		return toolText(string(data)), nil
	}
}
