package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/stefanoamorelli/lever-mcp/internal/client"
)

func listPostingsTool() *mcp.Tool {
	return &mcp.Tool{
		Annotations: readOnlyAnnotations(),
		Name:        "list_postings",
		Description: "List all job postings in your Lever account. Supports filtering by state, team, location, and commitment.",
		InputSchema: objectSchema(mergeProperties(paginationProperties(), map[string]any{
			"state":           prop("string", "Filter by state: published, internal, closed, draft, pending, rejected"),
			"team":            prop("string", "Filter by team name"),
			"location":        prop("string", "Filter by location"),
			"commitment":      prop("string", "Filter by commitment (e.g. Full-time)"),
			"confidentiality": prop("string", "Filter: all, non-confidential, or confidential"),
		}), nil),
	}
}

func listPostingsHandler(c client.LeverClient) mcp.ToolHandler {
	return func(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args, err := extractArgs(req)
		if err != nil {
			return toolError(err.Error()), nil
		}

		params := url.Values{}
		setPagination(params, args)
		setString(params, "state", getString(args, "state", ""))
		setString(params, "team", getString(args, "team", ""))
		setString(params, "location", getString(args, "location", ""))
		setString(params, "commitment", getString(args, "commitment", ""))
		setString(params, "confidentiality", getString(args, "confidentiality", ""))

		data, err := c.Get(ctx, "/postings", params)
		if err != nil {
			return toolErrorf("Failed to list postings: %v", err), nil
		}
		return toolText(string(data)), nil
	}
}

func getPostingTool() *mcp.Tool {
	return &mcp.Tool{
		Annotations: readOnlyAnnotations(),
		Name:        "get_posting",
		Description: "Retrieve a single job posting by ID from Lever.",
		InputSchema: objectSchema(map[string]any{
			"posting_id": prop("string", "The posting ID"),
		}, []string{"posting_id"}),
	}
}

func getPostingHandler(c client.LeverClient) mcp.ToolHandler {
	return func(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args, err := extractArgs(req)
		if err != nil {
			return toolError(err.Error()), nil
		}

		id := getString(args, "posting_id", "")
		if id == "" {
			return toolError("posting_id is required"), nil
		}

		data, err := c.Get(ctx, fmt.Sprintf("/postings/%s", id), nil)
		if err != nil {
			return toolErrorf("Failed to get posting: %v", err), nil
		}
		return toolText(string(data)), nil
	}
}

func createPostingTool() *mcp.Tool {
	return &mcp.Tool{
		Annotations: createAnnotations(),
		Name:        "create_posting",
		Description: "Create a new job posting in Lever.",
		InputSchema: objectSchema(map[string]any{
			"perform_as":  prop("string", "User ID to perform this action as"),
			"text":        prop("string", "Job posting title"),
			"state":       prop("string", "Posting state: published, internal, closed, draft"),
			"team":        prop("string", "Team name"),
			"department":  prop("string", "Department name"),
			"location":    prop("string", "Location"),
			"commitment":  prop("string", "Commitment (e.g. Full-time, Part-time)"),
			"description": prop("string", "Job description (HTML supported)"),
		}, []string{"perform_as", "text"}),
	}
}

func createPostingHandler(c client.LeverClient) mcp.ToolHandler {
	return func(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args, err := extractArgs(req)
		if err != nil {
			return toolError(err.Error()), nil
		}

		performAs := getString(args, "perform_as", "")
		text := getString(args, "text", "")
		if performAs == "" || text == "" {
			return toolError("perform_as and text are required"), nil
		}

		params := url.Values{}
		params.Set("perform_as", performAs)

		body := map[string]any{"text": text}
		for _, key := range []string{"state", "team", "department", "location", "commitment", "description"} {
			if v := getString(args, key, ""); v != "" {
				body[key] = v
			}
		}

		bodyJSON, _ := json.Marshal(body)
		data, err := c.Post(ctx, "/postings", params, bodyJSON)
		if err != nil {
			return toolErrorf("Failed to create posting: %v", err), nil
		}
		return toolText(string(data)), nil
	}
}

func updatePostingTool() *mcp.Tool {
	return &mcp.Tool{
		Annotations: updateAnnotations(),
		Name:        "update_posting",
		Description: "Update an existing job posting in Lever. This is a full replacement - all fields must be provided.",
		InputSchema: objectSchema(map[string]any{
			"posting_id":  prop("string", "The posting ID to update"),
			"perform_as":  prop("string", "User ID to perform this action as"),
			"text":        prop("string", "Job posting title"),
			"state":       prop("string", "Posting state: published, internal, closed, draft"),
			"team":        prop("string", "Team name"),
			"department":  prop("string", "Department name"),
			"location":    prop("string", "Location"),
			"commitment":  prop("string", "Commitment (e.g. Full-time, Part-time)"),
			"description": prop("string", "Job description (HTML supported)"),
		}, []string{"posting_id", "perform_as"}),
	}
}

func updatePostingHandler(c client.LeverClient) mcp.ToolHandler {
	return func(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args, err := extractArgs(req)
		if err != nil {
			return toolError(err.Error()), nil
		}

		id := getString(args, "posting_id", "")
		performAs := getString(args, "perform_as", "")
		if id == "" || performAs == "" {
			return toolError("posting_id and perform_as are required"), nil
		}

		params := url.Values{}
		params.Set("perform_as", performAs)

		body := make(map[string]any)
		for _, key := range []string{"text", "state", "team", "department", "location", "commitment", "description"} {
			if v := getString(args, key, ""); v != "" {
				body[key] = v
			}
		}

		bodyJSON, _ := json.Marshal(body)
		data, err := c.Put(ctx, fmt.Sprintf("/postings/%s", id), params, bodyJSON)
		if err != nil {
			return toolErrorf("Failed to update posting: %v", err), nil
		}
		return toolText(string(data)), nil
	}
}

func applyToPostingTool() *mcp.Tool {
	return &mcp.Tool{
		Annotations: createAnnotations(),
		Name:        "apply_to_posting",
		Description: "Submit an application to a published or internal posting in Lever programmatically.",
		InputSchema: objectSchema(map[string]any{
			"posting_id": prop("string", "The posting ID to apply to"),
			"name":       prop("string", "Applicant full name"),
			"email":      prop("string", "Applicant email address"),
			"phone":      prop("string", "Applicant phone number"),
			"org":        prop("string", "Applicant organization/company"),
			"urls":       map[string]any{"type": "array", "items": map[string]any{"type": "string"}, "description": "Applicant URLs (LinkedIn, portfolio, etc.)"},
			"comments":   prop("string", "Additional comments"),
			"origin":     prop("string", "Application origin"),
			"source":     prop("string", "Application source"),
		}, []string{"posting_id", "name", "email"}),
	}
}

func applyToPostingHandler(c client.LeverClient) mcp.ToolHandler {
	return func(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args, err := extractArgs(req)
		if err != nil {
			return toolError(err.Error()), nil
		}

		id := getString(args, "posting_id", "")
		name := getString(args, "name", "")
		email := getString(args, "email", "")
		if id == "" || name == "" || email == "" {
			return toolError("posting_id, name, and email are required"), nil
		}

		body := map[string]any{"name": name, "email": email}
		for _, key := range []string{"phone", "org", "comments", "origin", "source"} {
			if v := getString(args, key, ""); v != "" {
				body[key] = v
			}
		}
		if urls := getStringSlice(args, "urls"); len(urls) > 0 {
			body["urls"] = urls
		}

		bodyJSON, _ := json.Marshal(body)
		data, err := c.Post(ctx, fmt.Sprintf("/postings/%s/apply", id), nil, bodyJSON)
		if err != nil {
			return toolErrorf("Failed to apply to posting: %v", err), nil
		}
		return toolText(string(data)), nil
	}
}

func getPostingApplyTool() *mcp.Tool {
	return &mcp.Tool{
		Annotations: readOnlyAnnotations(),
		Name:        "get_posting_apply_questions",
		Description: "Retrieve the application form questions for a specific posting in Lever.",
		InputSchema: objectSchema(map[string]any{
			"posting_id": prop("string", "The posting ID"),
		}, []string{"posting_id"}),
	}
}

func getPostingApplyHandler(c client.LeverClient) mcp.ToolHandler {
	return func(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args, err := extractArgs(req)
		if err != nil {
			return toolError(err.Error()), nil
		}

		id := getString(args, "posting_id", "")
		if id == "" {
			return toolError("posting_id is required"), nil
		}

		data, err := c.Get(ctx, fmt.Sprintf("/postings/%s/apply", id), nil)
		if err != nil {
			return toolErrorf("Failed to get posting apply questions: %v", err), nil
		}
		return toolText(string(data)), nil
	}
}
