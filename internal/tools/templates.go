package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/stefanoamorelli/lever-mcp/internal/client"
)

func listFeedbackTemplatesTool() *mcp.Tool {
	return &mcp.Tool{
		Annotations: readOnlyAnnotations(),
		Name:        "list_feedback_templates",
		Description: "List all feedback templates in your Lever account.",
		InputSchema: objectSchema(paginationProperties(), nil),
	}
}

func listFeedbackTemplatesHandler(c client.LeverClient) mcp.ToolHandler {
	return func(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args, err := extractArgs(req)
		if err != nil {
			return toolError(err.Error()), nil
		}

		params := url.Values{}
		setPagination(params, args)

		data, err := c.Get(ctx, "/feedback_templates", params)
		if err != nil {
			return toolErrorf("Failed to list feedback templates: %v", err), nil
		}
		return toolText(string(data)), nil
	}
}

func createFeedbackTemplateTool() *mcp.Tool {
	return &mcp.Tool{
		Annotations: createAnnotations(),
		Name:        "create_feedback_template",
		Description: "Create a new feedback template in Lever.",
		InputSchema: objectSchema(map[string]any{
			"name":         prop("string", "Template name"),
			"instructions": prop("string", "Template instructions"),
			"group_uid":    prop("string", "Group UID for the template"),
			"fields":       map[string]any{"type": "array", "items": map[string]any{"type": "object"}, "description": "Array of field definitions"},
		}, []string{"name", "group_uid", "fields"}),
	}
}

func createFeedbackTemplateHandler(c client.LeverClient) mcp.ToolHandler {
	return func(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args, err := extractArgs(req)
		if err != nil {
			return toolError(err.Error()), nil
		}

		name := getString(args, "name", "")
		groupUID := getString(args, "group_uid", "")
		if name == "" || groupUID == "" {
			return toolError("name and group_uid are required"), nil
		}

		body := map[string]any{"name": name, "groupUid": groupUID}
		if v := getString(args, "instructions", ""); v != "" {
			body["instructions"] = v
		}
		if v, ok := args["fields"]; ok {
			body["fields"] = v
		}

		bodyJSON, _ := json.Marshal(body)
		data, err := c.Post(ctx, "/feedback_templates", nil, bodyJSON)
		if err != nil {
			return toolErrorf("Failed to create feedback template: %v", err), nil
		}
		return toolText(string(data)), nil
	}
}

func deleteFeedbackTemplateTool() *mcp.Tool {
	return &mcp.Tool{
		Annotations: deleteAnnotations(),
		Name:        "delete_feedback_template",
		Description: "Delete a feedback template from Lever. Only templates created via the API can be deleted.",
		InputSchema: objectSchema(map[string]any{
			"feedback_template_id": prop("string", "The feedback template ID"),
		}, []string{"feedback_template_id"}),
	}
}

func deleteFeedbackTemplateHandler(c client.LeverClient) mcp.ToolHandler {
	return func(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args, err := extractArgs(req)
		if err != nil {
			return toolError(err.Error()), nil
		}

		id := getString(args, "feedback_template_id", "")
		if id == "" {
			return toolError("feedback_template_id is required"), nil
		}

		data, err := c.Delete(ctx, fmt.Sprintf("/feedback_templates/%s", id), nil)
		if err != nil {
			return toolErrorf("Failed to delete feedback template: %v", err), nil
		}
		return toolText(string(data)), nil
	}
}

func listFormTemplatesTool() *mcp.Tool {
	return &mcp.Tool{
		Annotations: readOnlyAnnotations(),
		Name:        "list_form_templates",
		Description: "List all active profile form templates in your Lever account.",
		InputSchema: objectSchema(paginationProperties(), nil),
	}
}

func listFormTemplatesHandler(c client.LeverClient) mcp.ToolHandler {
	return func(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args, err := extractArgs(req)
		if err != nil {
			return toolError(err.Error()), nil
		}

		params := url.Values{}
		setPagination(params, args)

		data, err := c.Get(ctx, "/form_templates", params)
		if err != nil {
			return toolErrorf("Failed to list form templates: %v", err), nil
		}
		return toolText(string(data)), nil
	}
}
