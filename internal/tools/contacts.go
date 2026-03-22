package tools

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/stefanoamorelli/lever-mcp/internal/client"
)

func getContactTool() *mcp.Tool {
	return &mcp.Tool{
		Annotations: readOnlyAnnotations(),
		Name:        "get_contact",
		Description: "Retrieve a single contact by ID from Lever. Returns name, headline, location, emails, phones, and links.",
		InputSchema: objectSchema(map[string]any{
			"contact_id": prop("string", "The contact ID"),
		}, []string{"contact_id"}),
	}
}

func getContactHandler(c client.LeverClient) mcp.ToolHandler {
	return func(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args, err := extractArgs(req)
		if err != nil {
			return toolError(err.Error()), nil
		}

		id := getString(args, "contact_id", "")
		if id == "" {
			return toolError("contact_id is required"), nil
		}

		data, err := c.Get(ctx, fmt.Sprintf("/contacts/%s", id), nil)
		if err != nil {
			return toolErrorf("Failed to get contact: %v", err), nil
		}
		return toolText(string(data)), nil
	}
}
