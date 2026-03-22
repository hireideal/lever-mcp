package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/stefanoamorelli/lever-mcp/internal/client"
)

func listNotesTool() *mcp.Tool {
	return &mcp.Tool{
		Annotations: readOnlyAnnotations(),
		Name:        "list_opportunity_notes",
		Description: "List all notes for a specific opportunity in Lever.",
		InputSchema: objectSchema(mergeProperties(paginationProperties(), map[string]any{
			"opportunity_id": prop("string", "The opportunity ID"),
		}), []string{"opportunity_id"}),
	}
}

func listNotesHandler(c client.LeverClient) mcp.ToolHandler {
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

		data, err := c.Get(ctx, fmt.Sprintf("/opportunities/%s/notes", id), params)
		if err != nil {
			return toolErrorf("Failed to list notes: %v", err), nil
		}
		return toolText(string(data)), nil
	}
}

func createNoteTool() *mcp.Tool {
	return &mcp.Tool{
		Annotations: createAnnotations(),
		Name:        "create_opportunity_note",
		Description: "Create a note on an opportunity in Lever.",
		InputSchema: objectSchema(map[string]any{
			"opportunity_id":   prop("string", "The opportunity ID"),
			"perform_as":       prop("string", "User ID to perform this action as"),
			"value":            prop("string", "The note text content"),
			"notify_followers": prop("boolean", "Whether to notify opportunity followers"),
		}, []string{"opportunity_id", "perform_as", "value"}),
	}
}

func createNoteHandler(c client.LeverClient) mcp.ToolHandler {
	return func(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args, err := extractArgs(req)
		if err != nil {
			return toolError(err.Error()), nil
		}

		id := getString(args, "opportunity_id", "")
		performAs := getString(args, "perform_as", "")
		value := getString(args, "value", "")
		if id == "" || performAs == "" || value == "" {
			return toolError("opportunity_id, perform_as, and value are required"), nil
		}

		params := url.Values{}
		params.Set("perform_as", performAs)

		body := map[string]any{"value": value}
		if v, ok := getBool(args, "notify_followers"); ok {
			body["notifyFollowers"] = v
		}

		bodyJSON, _ := json.Marshal(body)
		data, err := c.Post(ctx, fmt.Sprintf("/opportunities/%s/notes", id), params, bodyJSON)
		if err != nil {
			return toolErrorf("Failed to create note: %v", err), nil
		}
		return toolText(string(data)), nil
	}
}

func deleteNoteTool() *mcp.Tool {
	return &mcp.Tool{
		Annotations: deleteAnnotations(),
		Name:        "delete_opportunity_note",
		Description: "Delete a note from an opportunity in Lever.",
		InputSchema: objectSchema(map[string]any{
			"opportunity_id": prop("string", "The opportunity ID"),
			"note_id":        prop("string", "The note ID to delete"),
		}, []string{"opportunity_id", "note_id"}),
	}
}

func deleteNoteHandler(c client.LeverClient) mcp.ToolHandler {
	return func(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args, err := extractArgs(req)
		if err != nil {
			return toolError(err.Error()), nil
		}

		id := getString(args, "opportunity_id", "")
		noteID := getString(args, "note_id", "")
		if id == "" || noteID == "" {
			return toolError("opportunity_id and note_id are required"), nil
		}

		data, err := c.Delete(ctx, fmt.Sprintf("/opportunities/%s/notes/%s", id, noteID), nil)
		if err != nil {
			return toolErrorf("Failed to delete note: %v", err), nil
		}
		return toolText(string(data)), nil
	}
}
