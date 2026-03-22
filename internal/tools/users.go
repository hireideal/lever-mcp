package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/stefanoamorelli/lever-mcp/internal/client"
)

func listUsersTool() *mcp.Tool {
	return &mcp.Tool{
		Annotations: readOnlyAnnotations(),
		Name:        "list_users",
		Description: "List all users in your Lever account. Supports filtering by email and including deactivated users.",
		InputSchema: objectSchema(mergeProperties(paginationProperties(), map[string]any{
			"email":               prop("string", "Filter by email address"),
			"include_deactivated": prop("boolean", "Include deactivated users in results"),
		}), nil),
	}
}

func listUsersHandler(c client.LeverClient) mcp.ToolHandler {
	return func(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args, err := extractArgs(req)
		if err != nil {
			return toolError(err.Error()), nil
		}

		params := url.Values{}
		setPagination(params, args)
		setString(params, "email", getString(args, "email", ""))
		setBool(params, "includeDeactivated", args)

		data, err := c.Get(ctx, "/users", params)
		if err != nil {
			return toolErrorf("Failed to list users: %v", err), nil
		}
		return toolText(string(data)), nil
	}
}

func getUserTool() *mcp.Tool {
	return &mcp.Tool{
		Annotations: readOnlyAnnotations(),
		Name:        "get_user",
		Description: "Retrieve a single user by ID from Lever.",
		InputSchema: objectSchema(map[string]any{
			"user_id": prop("string", "The user ID"),
		}, []string{"user_id"}),
	}
}

func getUserHandler(c client.LeverClient) mcp.ToolHandler {
	return func(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args, err := extractArgs(req)
		if err != nil {
			return toolError(err.Error()), nil
		}

		id := getString(args, "user_id", "")
		if id == "" {
			return toolError("user_id is required"), nil
		}

		data, err := c.Get(ctx, fmt.Sprintf("/users/%s", id), nil)
		if err != nil {
			return toolErrorf("Failed to get user: %v", err), nil
		}
		return toolText(string(data)), nil
	}
}

func createUserTool() *mcp.Tool {
	return &mcp.Tool{
		Annotations: createAnnotations(),
		Name:        "create_user",
		Description: "Create a new user in your Lever account.",
		InputSchema: objectSchema(map[string]any{
			"name":        prop("string", "User full name"),
			"email":       prop("string", "User email address"),
			"access_role": prop("string", "Access role: super admin, admin, team lead, member, limited member"),
		}, []string{"name", "email"}),
	}
}

func createUserHandler(c client.LeverClient) mcp.ToolHandler {
	return func(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args, err := extractArgs(req)
		if err != nil {
			return toolError(err.Error()), nil
		}

		name := getString(args, "name", "")
		email := getString(args, "email", "")
		if name == "" || email == "" {
			return toolError("name and email are required"), nil
		}

		body := map[string]any{"name": name, "email": email}
		if v := getString(args, "access_role", ""); v != "" {
			body["accessRole"] = v
		}

		bodyJSON, _ := json.Marshal(body)
		data, err := c.Post(ctx, "/users", nil, bodyJSON)
		if err != nil {
			return toolErrorf("Failed to create user: %v", err), nil
		}
		return toolText(string(data)), nil
	}
}

func deactivateUserTool() *mcp.Tool {
	return &mcp.Tool{
		Annotations: updateAnnotations(),
		Name:        "deactivate_user",
		Description: "Deactivate a user in your Lever account.",
		InputSchema: objectSchema(map[string]any{
			"user_id": prop("string", "The user ID to deactivate"),
		}, []string{"user_id"}),
	}
}

func deactivateUserHandler(c client.LeverClient) mcp.ToolHandler {
	return func(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args, err := extractArgs(req)
		if err != nil {
			return toolError(err.Error()), nil
		}

		id := getString(args, "user_id", "")
		if id == "" {
			return toolError("user_id is required"), nil
		}

		data, err := c.Post(ctx, fmt.Sprintf("/users/%s/deactivate", id), nil, nil)
		if err != nil {
			return toolErrorf("Failed to deactivate user: %v", err), nil
		}
		return toolText(string(data)), nil
	}
}

func reactivateUserTool() *mcp.Tool {
	return &mcp.Tool{
		Annotations: updateAnnotations(),
		Name:        "reactivate_user",
		Description: "Reactivate a previously deactivated user in your Lever account.",
		InputSchema: objectSchema(map[string]any{
			"user_id": prop("string", "The user ID to reactivate"),
		}, []string{"user_id"}),
	}
}

func reactivateUserHandler(c client.LeverClient) mcp.ToolHandler {
	return func(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args, err := extractArgs(req)
		if err != nil {
			return toolError(err.Error()), nil
		}

		id := getString(args, "user_id", "")
		if id == "" {
			return toolError("user_id is required"), nil
		}

		data, err := c.Post(ctx, fmt.Sprintf("/users/%s/reactivate", id), nil, nil)
		if err != nil {
			return toolErrorf("Failed to reactivate user: %v", err), nil
		}
		return toolText(string(data)), nil
	}
}
