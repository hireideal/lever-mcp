package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/stefanoamorelli/lever-mcp/internal/client"
)

func listFeedbackTool() *mcp.Tool {
	return &mcp.Tool{
		Annotations: readOnlyAnnotations(),
		Name:        "list_opportunity_feedback",
		Description: "List all feedback forms for a specific opportunity in Lever.",
		InputSchema: objectSchema(mergeProperties(paginationProperties(), map[string]any{
			"opportunity_id": prop("string", "The opportunity ID"),
		}), []string{"opportunity_id"}),
	}
}

func listFeedbackHandler(c client.LeverClient) mcp.ToolHandler {
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

		data, err := c.Get(ctx, fmt.Sprintf("/opportunities/%s/feedback", id), params)
		if err != nil {
			return toolErrorf("Failed to list feedback: %v", err), nil
		}
		return toolText(string(data)), nil
	}
}

func createFeedbackTool() *mcp.Tool {
	return &mcp.Tool{
		Annotations: createAnnotations(),
		Name:        "create_opportunity_feedback",
		Description: "Create a feedback form for a specific opportunity in Lever. Requires a base template ID.",
		InputSchema: objectSchema(map[string]any{
			"opportunity_id":   prop("string", "The opportunity ID"),
			"perform_as":       prop("string", "User ID to perform this action as"),
			"base_template_id": prop("string", "Base feedback template ID"),
			"panel_id":         prop("string", "Panel ID to associate feedback with"),
			"interview_id":     prop("string", "Interview ID to associate feedback with"),
			"completed_at":     prop("integer", "Timestamp when feedback was completed (Unix ms)"),
		}, []string{"opportunity_id", "perform_as", "base_template_id"}),
	}
}

func createFeedbackHandler(c client.LeverClient) mcp.ToolHandler {
	return func(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args, err := extractArgs(req)
		if err != nil {
			return toolError(err.Error()), nil
		}

		id := getString(args, "opportunity_id", "")
		performAs := getString(args, "perform_as", "")
		baseTemplateID := getString(args, "base_template_id", "")
		if id == "" || performAs == "" || baseTemplateID == "" {
			return toolError("opportunity_id, perform_as, and base_template_id are required"), nil
		}

		params := url.Values{}
		params.Set("perform_as", performAs)

		body := map[string]any{"baseTemplateId": baseTemplateID}
		if v := getString(args, "panel_id", ""); v != "" {
			body["panelId"] = v
		}
		if v := getString(args, "interview_id", ""); v != "" {
			body["interviewId"] = v
		}
		if v, ok := getInt(args, "completed_at"); ok {
			body["completedAt"] = v
		}

		bodyJSON, _ := json.Marshal(body)
		data, err := c.Post(ctx, fmt.Sprintf("/opportunities/%s/feedback", id), params, bodyJSON)
		if err != nil {
			return toolErrorf("Failed to create feedback: %v", err), nil
		}
		return toolText(string(data)), nil
	}
}

func getFeedbackTool() *mcp.Tool {
	return &mcp.Tool{
		Annotations: readOnlyAnnotations(),
		Name:        "get_opportunity_feedback",
		Description: "Retrieve a single feedback form for an opportunity in Lever.",
		InputSchema: objectSchema(map[string]any{
			"opportunity_id": prop("string", "The opportunity ID"),
			"feedback_id":    prop("string", "The feedback form ID"),
		}, []string{"opportunity_id", "feedback_id"}),
	}
}

func getFeedbackHandler(c client.LeverClient) mcp.ToolHandler {
	return func(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args, err := extractArgs(req)
		if err != nil {
			return toolError(err.Error()), nil
		}

		id := getString(args, "opportunity_id", "")
		feedbackID := getString(args, "feedback_id", "")
		if id == "" || feedbackID == "" {
			return toolError("opportunity_id and feedback_id are required"), nil
		}

		data, err := c.Get(ctx, fmt.Sprintf("/opportunities/%s/feedback/%s", id, feedbackID), nil)
		if err != nil {
			return toolErrorf("Failed to get feedback: %v", err), nil
		}
		return toolText(string(data)), nil
	}
}

func updateFeedbackTool() *mcp.Tool {
	return &mcp.Tool{
		Annotations: updateAnnotations(),
		Name:        "update_opportunity_feedback",
		Description: "Update a feedback form for an opportunity in Lever.",
		InputSchema: objectSchema(map[string]any{
			"opportunity_id": prop("string", "The opportunity ID"),
			"feedback_id":    prop("string", "The feedback form ID"),
			"perform_as":     prop("string", "User ID to perform this action as"),
			"completed_at":   prop("integer", "Timestamp when feedback was completed (Unix ms)"),
		}, []string{"opportunity_id", "feedback_id", "perform_as"}),
	}
}

func updateFeedbackHandler(c client.LeverClient) mcp.ToolHandler {
	return func(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args, err := extractArgs(req)
		if err != nil {
			return toolError(err.Error()), nil
		}

		id := getString(args, "opportunity_id", "")
		feedbackID := getString(args, "feedback_id", "")
		performAs := getString(args, "perform_as", "")
		if id == "" || feedbackID == "" || performAs == "" {
			return toolError("opportunity_id, feedback_id, and perform_as are required"), nil
		}

		params := url.Values{}
		params.Set("perform_as", performAs)

		body := make(map[string]any)
		if v, ok := getInt(args, "completed_at"); ok {
			body["completedAt"] = v
		}

		bodyJSON, _ := json.Marshal(body)
		data, err := c.Put(ctx, fmt.Sprintf("/opportunities/%s/feedback/%s", id, feedbackID), params, bodyJSON)
		if err != nil {
			return toolErrorf("Failed to update feedback: %v", err), nil
		}
		return toolText(string(data)), nil
	}
}

func deleteFeedbackTool() *mcp.Tool {
	return &mcp.Tool{
		Annotations: deleteAnnotations(),
		Name:        "delete_opportunity_feedback",
		Description: "Delete a feedback form from an opportunity in Lever.",
		InputSchema: objectSchema(map[string]any{
			"opportunity_id": prop("string", "The opportunity ID"),
			"feedback_id":    prop("string", "The feedback form ID"),
		}, []string{"opportunity_id", "feedback_id"}),
	}
}

func deleteFeedbackHandler(c client.LeverClient) mcp.ToolHandler {
	return func(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args, err := extractArgs(req)
		if err != nil {
			return toolError(err.Error()), nil
		}

		id := getString(args, "opportunity_id", "")
		feedbackID := getString(args, "feedback_id", "")
		if id == "" || feedbackID == "" {
			return toolError("opportunity_id and feedback_id are required"), nil
		}

		data, err := c.Delete(ctx, fmt.Sprintf("/opportunities/%s/feedback/%s", id, feedbackID), nil)
		if err != nil {
			return toolErrorf("Failed to delete feedback: %v", err), nil
		}
		return toolText(string(data)), nil
	}
}
