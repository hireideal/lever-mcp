package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/stefanoamorelli/lever-mcp/internal/client"
)

func listInterviewsTool() *mcp.Tool {
	return &mcp.Tool{
		Annotations: readOnlyAnnotations(),
		Name:        "list_opportunity_interviews",
		Description: "List all interviews for a specific opportunity in Lever.",
		InputSchema: objectSchema(mergeProperties(paginationProperties(), map[string]any{
			"opportunity_id": prop("string", "The opportunity ID"),
		}), []string{"opportunity_id"}),
	}
}

func listInterviewsHandler(c client.LeverClient) mcp.ToolHandler {
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

		data, err := c.Get(ctx, fmt.Sprintf("/opportunities/%s/interviews", id), params)
		if err != nil {
			return toolErrorf("Failed to list interviews: %v", err), nil
		}
		return toolText(string(data)), nil
	}
}

func updateInterviewTool() *mcp.Tool {
	return &mcp.Tool{
		Annotations: updateAnnotations(),
		Name:        "update_opportunity_interview",
		Description: "Update an interview for an opportunity in Lever. Only interviews within externally managed panels can be updated. Requires the entire interview object.",
		InputSchema: objectSchema(map[string]any{
			"opportunity_id": prop("string", "The opportunity ID"),
			"interview_id":   prop("string", "The interview ID to update"),
			"perform_as":     prop("string", "User ID to perform this action as"),
			"panel_id":       prop("string", "Panel ID the interview belongs to"),
			"date":           prop("integer", "Interview date (Unix timestamp in ms)"),
			"duration":       prop("integer", "Interview duration in minutes"),
			"location":       prop("string", "Interview location"),
			"subject":        prop("string", "Interview subject"),
			"note":           prop("string", "Interview note"),
			"interviewers":   map[string]any{"type": "array", "items": map[string]any{"type": "object"}, "description": "Array of interviewer objects [{id, name, email}]"},
		}, []string{"opportunity_id", "interview_id", "perform_as"}),
	}
}

func updateInterviewHandler(c client.LeverClient) mcp.ToolHandler {
	return func(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args, err := extractArgs(req)
		if err != nil {
			return toolError(err.Error()), nil
		}

		id := getString(args, "opportunity_id", "")
		interviewID := getString(args, "interview_id", "")
		performAs := getString(args, "perform_as", "")
		if id == "" || interviewID == "" || performAs == "" {
			return toolError("opportunity_id, interview_id, and perform_as are required"), nil
		}

		params := url.Values{}
		params.Set("perform_as", performAs)

		body := make(map[string]any)
		for _, key := range []string{"panel_id", "location", "subject", "note"} {
			if v := getString(args, key, ""); v != "" {
				body[key] = v
			}
		}
		if v, ok := getInt(args, "date"); ok {
			body["date"] = v
		}
		if v, ok := getInt(args, "duration"); ok {
			body["duration"] = v
		}
		if v, ok := args["interviewers"]; ok {
			body["interviewers"] = v
		}

		bodyJSON, _ := json.Marshal(body)
		data, err := c.Put(ctx, fmt.Sprintf("/opportunities/%s/interviews/%s", id, interviewID), params, bodyJSON)
		if err != nil {
			return toolErrorf("Failed to update interview: %v", err), nil
		}
		return toolText(string(data)), nil
	}
}

func deleteInterviewTool() *mcp.Tool {
	return &mcp.Tool{
		Annotations: deleteAnnotations(),
		Name:        "delete_opportunity_interview",
		Description: "Delete an interview from an opportunity in Lever. Only interviews within externally managed panels can be deleted.",
		InputSchema: objectSchema(map[string]any{
			"opportunity_id": prop("string", "The opportunity ID"),
			"interview_id":   prop("string", "The interview ID to delete"),
		}, []string{"opportunity_id", "interview_id"}),
	}
}

func deleteInterviewHandler(c client.LeverClient) mcp.ToolHandler {
	return func(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args, err := extractArgs(req)
		if err != nil {
			return toolError(err.Error()), nil
		}

		id := getString(args, "opportunity_id", "")
		interviewID := getString(args, "interview_id", "")
		if id == "" || interviewID == "" {
			return toolError("opportunity_id and interview_id are required"), nil
		}

		data, err := c.Delete(ctx, fmt.Sprintf("/opportunities/%s/interviews/%s", id, interviewID), nil)
		if err != nil {
			return toolErrorf("Failed to delete interview: %v", err), nil
		}
		return toolText(string(data)), nil
	}
}
