package tools

import (
	"context"
	"encoding/json"
	"net/url"
	"testing"

	"github.com/stefanoamorelli/lever-mcp/internal/testutil"
)

func TestListNotesHandler(t *testing.T) {
	mock := &testutil.MockLeverClient{
		GetFunc: func(ctx context.Context, path string, params url.Values) (json.RawMessage, error) {
			if path != "/opportunities/opp-123/notes" {
				t.Errorf("expected path /opportunities/opp-123/notes, got %s", path)
			}
			return testutil.NoteListFixture(), nil
		},
	}

	handler := listNotesHandler(mock)
	result, err := handler(context.Background(), makeReq(map[string]any{"opportunity_id": "opp-123"}))
	if err != nil {
		t.Fatal(err)
	}
	if result.IsError {
		t.Fatal("expected success")
	}
}

func TestListNotesHandler_MissingID(t *testing.T) {
	mock := &testutil.MockLeverClient{}
	handler := listNotesHandler(mock)
	result, err := handler(context.Background(), makeReq(map[string]any{}))
	if err != nil {
		t.Fatal(err)
	}
	if !result.IsError {
		t.Fatal("expected error when opportunity_id missing")
	}
}

func TestCreateNoteHandler(t *testing.T) {
	mock := &testutil.MockLeverClient{
		PostFunc: func(ctx context.Context, path string, params url.Values, body json.RawMessage) (json.RawMessage, error) {
			if path != "/opportunities/opp-123/notes" {
				t.Errorf("expected path /opportunities/opp-123/notes, got %s", path)
			}
			if params.Get("perform_as") != "user-123" {
				t.Errorf("expected perform_as=user-123, got %s", params.Get("perform_as"))
			}
			var b map[string]any
			json.Unmarshal(body, &b)
			if b["value"] != "Test note" {
				t.Errorf("expected value=Test note, got %v", b["value"])
			}
			return testutil.NoteFixture(), nil
		},
	}

	handler := createNoteHandler(mock)
	result, err := handler(context.Background(), makeReq(map[string]any{
		"opportunity_id": "opp-123",
		"perform_as":     "user-123",
		"value":          "Test note",
	}))
	if err != nil {
		t.Fatal(err)
	}
	if result.IsError {
		t.Fatal("expected success")
	}
}

func TestDeleteNoteHandler(t *testing.T) {
	mock := &testutil.MockLeverClient{
		DeleteFunc: func(ctx context.Context, path string, params url.Values) (json.RawMessage, error) {
			if path != "/opportunities/opp-123/notes/note-456" {
				t.Errorf("expected path /opportunities/opp-123/notes/note-456, got %s", path)
			}
			return json.RawMessage(`{}`), nil
		},
	}

	handler := deleteNoteHandler(mock)
	result, err := handler(context.Background(), makeReq(map[string]any{
		"opportunity_id": "opp-123",
		"note_id":        "note-456",
	}))
	if err != nil {
		t.Fatal(err)
	}
	if result.IsError {
		t.Fatal("expected success")
	}
}
