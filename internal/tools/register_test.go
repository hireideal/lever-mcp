package tools

import (
	"context"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/stefanoamorelli/lever-mcp/internal/testutil"
)

func listToolNames(t *testing.T, filter ToolFilter) []string {
	t.Helper()

	server := mcp.NewServer(&mcp.Implementation{
		Name:    "lever-mcp",
		Version: "0.1.0-test",
	}, nil)

	mock := &testutil.MockLeverClient{}
	RegisterAll(server, mock, filter)

	clientTransport, serverTransport := mcp.NewInMemoryTransports()

	_, err := server.Connect(context.Background(), serverTransport, nil)
	if err != nil {
		t.Fatalf("server.Connect failed: %v", err)
	}

	client := mcp.NewClient(&mcp.Implementation{
		Name:    "test-client",
		Version: "1.0.0",
	}, nil)

	session, err := client.Connect(context.Background(), clientTransport, nil)
	if err != nil {
		t.Fatalf("client.Connect failed: %v", err)
	}
	defer session.Close()

	toolResult, err := session.ListTools(context.Background(), nil)
	if err != nil {
		t.Fatalf("ListTools failed: %v", err)
	}

	var names []string
	for _, tool := range toolResult.Tools {
		names = append(names, tool.Name)
	}
	return names
}

func listTools(t *testing.T) []*mcp.Tool {
	t.Helper()

	server := mcp.NewServer(&mcp.Implementation{
		Name:    "lever-mcp",
		Version: "0.1.0-test",
	}, nil)

	mock := &testutil.MockLeverClient{}
	RegisterAll(server, mock, nil)

	clientTransport, serverTransport := mcp.NewInMemoryTransports()

	_, err := server.Connect(context.Background(), serverTransport, nil)
	if err != nil {
		t.Fatalf("server.Connect failed: %v", err)
	}

	client := mcp.NewClient(&mcp.Implementation{
		Name:    "test-client",
		Version: "1.0.0",
	}, nil)

	session, err := client.Connect(context.Background(), clientTransport, nil)
	if err != nil {
		t.Fatalf("client.Connect failed: %v", err)
	}
	defer session.Close()

	toolResult, err := session.ListTools(context.Background(), nil)
	if err != nil {
		t.Fatalf("ListTools failed: %v", err)
	}

	return toolResult.Tools
}

var expectedTools = []string{
	// Opportunities
	"list_opportunities",
	"get_opportunity",
	"create_opportunity",
	"list_deleted_opportunities",
	// Opportunity actions
	"archive_opportunity",
	"change_opportunity_stage",
	"add_opportunity_tags",
	"add_opportunity_sources",
	"add_opportunity_links",
	"remove_opportunity_tags",
	"remove_opportunity_sources",
	"remove_opportunity_links",
	// Notes
	"list_opportunity_notes",
	"create_opportunity_note",
	"delete_opportunity_note",
	// Offers
	"list_opportunity_offers",
	// Interviews
	"list_opportunity_interviews",
	"update_opportunity_interview",
	"delete_opportunity_interview",
	// Feedback
	"list_opportunity_feedback",
	"get_opportunity_feedback",
	"create_opportunity_feedback",
	"update_opportunity_feedback",
	"delete_opportunity_feedback",
	// Panels
	"create_opportunity_panel",
	"delete_opportunity_panel",
	// Referrals
	"list_opportunity_referrals",
	// Resumes
	"list_opportunity_resumes",
	// Files
	"list_opportunity_files",
	// Applications
	"list_opportunity_applications",
	"get_opportunity_application",
	// Forms
	"list_opportunity_forms",
	// Archive reasons
	"list_archive_reasons",
	"get_archive_reason",
	// Contacts
	"get_contact",
	// Postings
	"list_postings",
	"get_posting",
	"create_posting",
	"update_posting",
	"apply_to_posting",
	"get_posting_apply_questions",
	// Users
	"list_users",
	"get_user",
	"create_user",
	"deactivate_user",
	"reactivate_user",
	// Stages
	"list_stages",
	// Sources
	"list_sources",
	// Tags
	"list_tags",
	// Requisitions
	"list_requisitions",
	"get_requisition",
	// Feedback templates
	"list_feedback_templates",
	"create_feedback_template",
	"delete_feedback_template",
	// Form templates
	"list_form_templates",
	// Webhooks
	"list_webhooks",
	"create_webhook",
	"delete_webhook",
	// EEO
	"list_eeo_responses",
	// Audit events
	"list_audit_events",
}

func TestConformance_AllToolsRegistered(t *testing.T) {
	names := listToolNames(t, nil)

	if len(names) != len(expectedTools) {
		t.Errorf("expected %d tools, got %d", len(expectedTools), len(names))
		for _, name := range names {
			t.Logf("  registered: %s", name)
		}
	}

	registeredNames := make(map[string]bool)
	for _, name := range names {
		registeredNames[name] = true
	}

	for _, name := range expectedTools {
		if !registeredNames[name] {
			t.Errorf("expected tool %q to be registered", name)
		}
	}
}

func TestConformance_ToolsHaveDescriptions(t *testing.T) {
	tools := listTools(t)

	for _, tool := range tools {
		if tool.Description == "" {
			t.Errorf("tool %q has no description", tool.Name)
		}
		if tool.InputSchema == nil {
			t.Errorf("tool %q has no input schema", tool.Name)
		}
	}
}

func TestConformance_ToolAnnotations(t *testing.T) {
	tools := listTools(t)

	readOnlyTools := map[string]bool{
		"list_opportunities": true, "get_opportunity": true, "list_deleted_opportunities": true,
		"list_opportunity_notes": true, "list_opportunity_offers": true,
		"list_opportunity_interviews": true, "list_opportunity_feedback": true,
		"get_opportunity_feedback":   true,
		"list_opportunity_referrals": true, "list_opportunity_resumes": true,
		"list_opportunity_files": true, "list_opportunity_applications": true,
		"get_opportunity_application": true, "list_opportunity_forms": true,
		"list_archive_reasons": true, "get_archive_reason": true,
		"get_contact": true, "list_postings": true, "get_posting": true,
		"get_posting_apply_questions": true, "list_users": true, "get_user": true,
		"list_stages": true, "list_sources": true, "list_tags": true,
		"list_requisitions": true, "get_requisition": true,
		"list_feedback_templates": true, "list_form_templates": true,
		"list_webhooks":      true,
		"list_eeo_responses": true, "list_audit_events": true,
	}

	destructiveTools := map[string]bool{
		"delete_opportunity_note": true, "delete_opportunity_panel": true,
		"delete_opportunity_interview": true, "delete_opportunity_feedback": true,
		"delete_feedback_template": true, "delete_webhook": true,
	}

	for _, tool := range tools {
		if tool.Annotations == nil {
			t.Errorf("tool %q has no annotations", tool.Name)
			continue
		}

		if readOnlyTools[tool.Name] {
			if !tool.Annotations.ReadOnlyHint {
				t.Errorf("tool %q should have ReadOnlyHint=true", tool.Name)
			}
			if tool.Annotations.DestructiveHint == nil || *tool.Annotations.DestructiveHint {
				t.Errorf("tool %q should have DestructiveHint=false", tool.Name)
			}
		}

		if destructiveTools[tool.Name] {
			if tool.Annotations.ReadOnlyHint {
				t.Errorf("tool %q should have ReadOnlyHint=false", tool.Name)
			}
			if tool.Annotations.DestructiveHint == nil || !*tool.Annotations.DestructiveHint {
				t.Errorf("tool %q should have DestructiveHint=true", tool.Name)
			}
		}
	}
}

func TestRegisterAll_EnabledFilter(t *testing.T) {
	filter := NewToolFilter("list_opportunities,get_opportunity", "")
	names := listToolNames(t, filter)

	if len(names) != 2 {
		t.Fatalf("expected 2 tools, got %d: %v", len(names), names)
	}

	set := make(map[string]bool)
	for _, n := range names {
		set[n] = true
	}
	if !set["list_opportunities"] {
		t.Error("expected list_opportunities to be registered")
	}
	if !set["get_opportunity"] {
		t.Error("expected get_opportunity to be registered")
	}
}

func TestRegisterAll_DisabledFilter(t *testing.T) {
	filter := NewToolFilter("", "list_audit_events")
	names := listToolNames(t, filter)

	for _, n := range names {
		if n == "list_audit_events" {
			t.Error("list_audit_events should not be registered when disabled")
		}
	}

	if len(names) != len(expectedTools)-1 {
		t.Errorf("expected %d tools, got %d", len(expectedTools)-1, len(names))
	}
}

func TestNewToolFilter(t *testing.T) {
	t.Run("nil when neither set", func(t *testing.T) {
		f := NewToolFilter("", "")
		if f != nil {
			t.Error("expected nil filter")
		}
	})

	t.Run("enabled only", func(t *testing.T) {
		f := NewToolFilter("a,b", "")
		if !f("a") || !f("b") {
			t.Error("expected a and b to pass")
		}
		if f("c") {
			t.Error("expected c to be rejected")
		}
	})

	t.Run("disabled only", func(t *testing.T) {
		f := NewToolFilter("", "x")
		if f("x") {
			t.Error("expected x to be rejected")
		}
		if !f("y") {
			t.Error("expected y to pass")
		}
	})

	t.Run("enabled takes precedence", func(t *testing.T) {
		f := NewToolFilter("a", "b")
		if !f("a") {
			t.Error("expected a to pass")
		}
		if f("b") {
			t.Error("expected b to be rejected")
		}
	})

	t.Run("whitespace trimming", func(t *testing.T) {
		f := NewToolFilter(" a , b ", "")
		if !f("a") || !f("b") {
			t.Error("expected trimmed names to pass")
		}
	})
}
