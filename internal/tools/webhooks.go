package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/stefanoamorelli/lever-mcp/internal/client"
)

func listWebhooksTool() *mcp.Tool {
	return &mcp.Tool{
		Annotations: readOnlyAnnotations(),
		Name:        "list_webhooks",
		Description: "List all configured webhooks in your Lever account.",
		InputSchema: objectSchema(paginationProperties(), nil),
	}
}

func listWebhooksHandler(c client.LeverClient) mcp.ToolHandler {
	return func(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args, err := extractArgs(req)
		if err != nil {
			return toolError(err.Error()), nil
		}

		params := url.Values{}
		setPagination(params, args)

		data, err := c.Get(ctx, "/webhooks", params)
		if err != nil {
			return toolErrorf("Failed to list webhooks: %v", err), nil
		}
		return toolText(string(data)), nil
	}
}

func createWebhookTool() *mcp.Tool {
	return &mcp.Tool{
		Annotations: createAnnotations(),
		Name:        "create_webhook",
		Description: "Create a new webhook in Lever. Specify the URL and event to listen for.",
		InputSchema: objectSchema(map[string]any{
			"url":          prop("string", "The webhook callback URL"),
			"event":        prop("string", "The event to subscribe to (e.g. candidateStageChange, candidateHired, candidateArchived)"),
			"secret_token": prop("string", "Optional secret token for webhook signature verification"),
		}, []string{"url", "event"}),
	}
}

func createWebhookHandler(c client.LeverClient) mcp.ToolHandler {
	return func(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args, err := extractArgs(req)
		if err != nil {
			return toolError(err.Error()), nil
		}

		webhookURL := getString(args, "url", "")
		event := getString(args, "event", "")
		if webhookURL == "" || event == "" {
			return toolError("url and event are required"), nil
		}

		body := map[string]any{"url": webhookURL, "event": event}
		if v := getString(args, "secret_token", ""); v != "" {
			body["secretToken"] = v
		}

		bodyJSON, _ := json.Marshal(body)
		data, err := c.Post(ctx, "/webhooks", nil, bodyJSON)
		if err != nil {
			return toolErrorf("Failed to create webhook: %v", err), nil
		}
		return toolText(string(data)), nil
	}
}

func deleteWebhookTool() *mcp.Tool {
	return &mcp.Tool{
		Annotations: deleteAnnotations(),
		Name:        "delete_webhook",
		Description: "Delete a webhook from Lever.",
		InputSchema: objectSchema(map[string]any{
			"webhook_id": prop("string", "The webhook ID to delete"),
		}, []string{"webhook_id"}),
	}
}

func deleteWebhookHandler(c client.LeverClient) mcp.ToolHandler {
	return func(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args, err := extractArgs(req)
		if err != nil {
			return toolError(err.Error()), nil
		}

		id := getString(args, "webhook_id", "")
		if id == "" {
			return toolError("webhook_id is required"), nil
		}

		data, err := c.Delete(ctx, fmt.Sprintf("/webhooks/%s", id), nil)
		if err != nil {
			return toolErrorf("Failed to delete webhook: %v", err), nil
		}
		return toolText(string(data)), nil
	}
}
