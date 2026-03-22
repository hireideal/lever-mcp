package tools

import (
	"net/url"
	"testing"
)

func TestProp(t *testing.T) {
	p := prop("string", "test description")
	if p["type"] != "string" {
		t.Errorf("expected type=string, got %v", p["type"])
	}
	if p["description"] != "test description" {
		t.Errorf("expected description=test description, got %v", p["description"])
	}
}

func TestObjectSchema(t *testing.T) {
	schema := objectSchema(map[string]any{"name": prop("string", "name")}, []string{"name"})
	if schema["type"] != "object" {
		t.Error("expected type=object")
	}
	req := schema["required"].([]string)
	if len(req) != 1 || req[0] != "name" {
		t.Error("expected required=[name]")
	}
}

func TestObjectSchema_NoRequired(t *testing.T) {
	schema := objectSchema(map[string]any{"name": prop("string", "name")}, nil)
	if _, ok := schema["required"]; ok {
		t.Error("expected no required field")
	}
}

func TestGetString(t *testing.T) {
	args := map[string]any{"name": "test", "count": 42}
	if v := getString(args, "name", ""); v != "test" {
		t.Errorf("expected test, got %s", v)
	}
	if v := getString(args, "missing", "default"); v != "default" {
		t.Errorf("expected default, got %s", v)
	}
	if v := getString(args, "count", "default"); v != "default" {
		t.Errorf("expected default for non-string, got %s", v)
	}
}

func TestGetInt(t *testing.T) {
	args := map[string]any{"count": float64(42)}
	v, ok := getInt(args, "count")
	if !ok || v != 42 {
		t.Errorf("expected 42, got %d", v)
	}
	_, ok = getInt(args, "missing")
	if ok {
		t.Error("expected not ok for missing key")
	}
}

func TestGetBool(t *testing.T) {
	args := map[string]any{"active": true}
	v, ok := getBool(args, "active")
	if !ok || !v {
		t.Error("expected true")
	}
	_, ok = getBool(args, "missing")
	if ok {
		t.Error("expected not ok for missing key")
	}
}

func TestSetString(t *testing.T) {
	params := url.Values{}
	setString(params, "key", "value")
	if params.Get("key") != "value" {
		t.Error("expected value to be set")
	}
	setString(params, "empty", "")
	if params.Get("empty") != "" {
		t.Error("expected empty value to not be set")
	}
}

func TestMergeProperties(t *testing.T) {
	a := map[string]any{"x": 1}
	b := map[string]any{"y": 2}
	merged := mergeProperties(a, b)
	if len(merged) != 2 {
		t.Errorf("expected 2 properties, got %d", len(merged))
	}
}

func TestGetStringSlice(t *testing.T) {
	args := map[string]any{
		"tags":   []any{"a", "b", "c"},
		"single": "x",
	}
	tags := getStringSlice(args, "tags")
	if len(tags) != 3 {
		t.Errorf("expected 3 tags, got %d", len(tags))
	}
	single := getStringSlice(args, "single")
	if len(single) != 1 || single[0] != "x" {
		t.Error("expected single string to be wrapped in slice")
	}
	missing := getStringSlice(args, "missing")
	if missing != nil {
		t.Error("expected nil for missing key")
	}
}

func TestBoolPtr(t *testing.T) {
	tr := boolPtr(true)
	if !*tr {
		t.Error("expected true")
	}
	fl := boolPtr(false)
	if *fl {
		t.Error("expected false")
	}
}

func TestReadOnlyAnnotations(t *testing.T) {
	a := readOnlyAnnotations()
	if !a.ReadOnlyHint {
		t.Error("expected ReadOnlyHint=true")
	}
	if a.DestructiveHint == nil || *a.DestructiveHint {
		t.Error("expected DestructiveHint=false")
	}
}

func TestCreateAnnotations(t *testing.T) {
	a := createAnnotations()
	if a.ReadOnlyHint {
		t.Error("expected ReadOnlyHint=false")
	}
	if a.DestructiveHint == nil || *a.DestructiveHint {
		t.Error("expected DestructiveHint=false")
	}
}

func TestDeleteAnnotations(t *testing.T) {
	a := deleteAnnotations()
	if a.ReadOnlyHint {
		t.Error("expected ReadOnlyHint=false")
	}
	if a.DestructiveHint == nil || !*a.DestructiveHint {
		t.Error("expected DestructiveHint=true")
	}
	if !a.IdempotentHint {
		t.Error("expected IdempotentHint=true")
	}
}

func TestUpdateAnnotations(t *testing.T) {
	a := updateAnnotations()
	if a.ReadOnlyHint {
		t.Error("expected ReadOnlyHint=false")
	}
	if a.DestructiveHint == nil || *a.DestructiveHint {
		t.Error("expected DestructiveHint=false")
	}
	if !a.IdempotentHint {
		t.Error("expected IdempotentHint=true")
	}
}
