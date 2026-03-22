package tools

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func toolText(text string) *mcp.CallToolResult {
	return &mcp.CallToolResult{
		Content: []mcp.Content{&mcp.TextContent{Text: text}},
	}
}

func toolError(msg string) *mcp.CallToolResult {
	return &mcp.CallToolResult{
		Content: []mcp.Content{&mcp.TextContent{Text: msg}},
		IsError: true,
	}
}

func toolErrorf(format string, args ...any) *mcp.CallToolResult {
	return toolError(fmt.Sprintf(format, args...))
}

func extractArgs(req *mcp.CallToolRequest) (map[string]any, error) {
	if len(req.Params.Arguments) == 0 {
		return map[string]any{}, nil
	}
	var args map[string]any
	if err := json.Unmarshal(req.Params.Arguments, &args); err != nil {
		return nil, fmt.Errorf("invalid arguments: %w", err)
	}
	return args, nil
}

func getString(args map[string]any, key, defaultVal string) string {
	v, ok := args[key]
	if !ok {
		return defaultVal
	}
	s, ok := v.(string)
	if !ok {
		return defaultVal
	}
	return s
}

func getInt(args map[string]any, key string) (int, bool) {
	v, ok := args[key]
	if !ok {
		return 0, false
	}
	switch n := v.(type) {
	case float64:
		return int(n), true
	case json.Number:
		i, err := n.Int64()
		if err != nil {
			return 0, false
		}
		return int(i), true
	}
	return 0, false
}

func getBool(args map[string]any, key string) (bool, bool) {
	v, ok := args[key]
	if !ok {
		return false, false
	}
	b, ok := v.(bool)
	return b, ok
}

func setString(params url.Values, key, value string) {
	if value != "" {
		params.Set(key, value)
	}
}

func setInt(params url.Values, key string, args map[string]any) {
	if v, ok := getInt(args, key); ok {
		params.Set(key, strconv.Itoa(v))
	}
}

func setBool(params url.Values, key string, args map[string]any) {
	if v, ok := getBool(args, key); ok {
		params.Set(key, strconv.FormatBool(v))
	}
}

func prop(typ, description string) map[string]any {
	return map[string]any{"type": typ, "description": description}
}

func objectSchema(properties map[string]any, required []string) map[string]any {
	schema := map[string]any{"type": "object", "properties": properties}
	if len(required) > 0 {
		schema["required"] = required
	}
	return schema
}

func paginationProperties() map[string]any {
	return map[string]any{
		"limit":  prop("integer", "Number of results to return (1-100, default 50)"),
		"offset": prop("string", "Pagination offset token from previous response"),
	}
}

func setPagination(params url.Values, args map[string]any) {
	setInt(params, "limit", args)
	setString(params, "offset", getString(args, "offset", ""))
}

func mergeProperties(maps ...map[string]any) map[string]any {
	result := make(map[string]any)
	for _, m := range maps {
		for k, v := range m {
			result[k] = v
		}
	}
	return result
}

func getStringSlice(args map[string]any, key string) []string {
	v, ok := args[key]
	if !ok {
		return nil
	}
	switch val := v.(type) {
	case []any:
		result := make([]string, 0, len(val))
		for _, item := range val {
			if s, ok := item.(string); ok {
				result = append(result, s)
			}
		}
		return result
	case string:
		return []string{val}
	}
	return nil
}

func boolPtr(v bool) *bool { return &v }

func readOnlyAnnotations() *mcp.ToolAnnotations {
	return &mcp.ToolAnnotations{
		ReadOnlyHint:    true,
		DestructiveHint: boolPtr(false),
		OpenWorldHint:   boolPtr(true),
		Title:           "",
	}
}

func createAnnotations() *mcp.ToolAnnotations {
	return &mcp.ToolAnnotations{
		ReadOnlyHint:    false,
		DestructiveHint: boolPtr(false),
		IdempotentHint:  false,
		OpenWorldHint:   boolPtr(true),
	}
}

func updateAnnotations() *mcp.ToolAnnotations {
	return &mcp.ToolAnnotations{
		ReadOnlyHint:    false,
		DestructiveHint: boolPtr(false),
		IdempotentHint:  true,
		OpenWorldHint:   boolPtr(true),
	}
}

func deleteAnnotations() *mcp.ToolAnnotations {
	return &mcp.ToolAnnotations{
		ReadOnlyHint:    false,
		DestructiveHint: boolPtr(true),
		IdempotentHint:  true,
		OpenWorldHint:   boolPtr(true),
	}
}
