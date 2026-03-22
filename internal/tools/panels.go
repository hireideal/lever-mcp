package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/stefanoamorelli/lever-mcp/internal/client"
)

func createPanelTool() *mcp.Tool {
	return &mcp.Tool{
		Annotations: createAnnotations(),
		Name:        "create_opportunity_panel",
		Description: "Create an interview panel for an opportunity in Lever. The panel must include a timezone and at least one interview.",
		InputSchema: objectSchema(map[string]any{
			"opportunity_id": prop("string", "The opportunity ID"),
			"perform_as":     prop("string", "User ID to perform this action as"),
			"timezone":       prop("string", "Timezone for the panel (e.g. America/New_York)"),
			"interviews":     map[string]any{"type": "array", "items": map[string]any{"type": "object"}, "description": "Array of interview objects with date, duration, interviewers, etc."},
			"note":           prop("string", "Optional note for the panel"),
		}, []string{"opportunity_id", "perform_as", "timezone", "interviews"}),
	}
}

func createPanelHandler(c client.LeverClient) mcp.ToolHandler {
	return func(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args, err := extractArgs(req)
		if err != nil {
			return toolError(err.Error()), nil
		}

		id := getString(args, "opportunity_id", "")
		performAs := getString(args, "perform_as", "")
		timezone := getString(args, "timezone", "")
		if id == "" || performAs == "" || timezone == "" {
			return toolError("opportunity_id, perform_as, and timezone are required"), nil
		}

		params := url.Values{}
		params.Set("perform_as", performAs)

		body := map[string]any{"timezone": timezone}
		if v, ok := args["interviews"]; ok {
			body["interviews"] = v
		}
		if v := getString(args, "note", ""); v != "" {
			body["note"] = v
		}

		bodyJSON, _ := json.Marshal(body)
		data, err := c.Post(ctx, fmt.Sprintf("/opportunities/%s/panels", id), params, bodyJSON)
		if err != nil {
			return toolErrorf("Failed to create panel: %v", err), nil
		}
		return toolText(string(data)), nil
	}
}

func deletePanelTool() *mcp.Tool {
	return &mcp.Tool{
		Annotations: deleteAnnotations(),
		Name:        "delete_opportunity_panel",
		Description: "Delete an interview panel from an opportunity in Lever. Only panels where externallyManaged is true can be deleted.",
		InputSchema: objectSchema(map[string]any{
			"opportunity_id": prop("string", "The opportunity ID"),
			"panel_id":       prop("string", "The panel ID to delete"),
		}, []string{"opportunity_id", "panel_id"}),
	}
}

func deletePanelHandler(c client.LeverClient) mcp.ToolHandler {
	return func(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args, err := extractArgs(req)
		if err != nil {
			return toolError(err.Error()), nil
		}

		id := getString(args, "opportunity_id", "")
		panelID := getString(args, "panel_id", "")
		if id == "" || panelID == "" {
			return toolError("opportunity_id and panel_id are required"), nil
		}

		data, err := c.Delete(ctx, fmt.Sprintf("/opportunities/%s/panels/%s", id, panelID), nil)
		if err != nil {
			return toolErrorf("Failed to delete panel: %v", err), nil
		}
		return toolText(string(data)), nil
	}
}
