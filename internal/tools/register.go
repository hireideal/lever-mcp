package tools

import (
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/stefanoamorelli/lever-mcp/internal/client"
)

type ToolFilter func(toolName string) bool

func NewToolFilter(enabled, disabled string) ToolFilter {
	if enabled != "" {
		set := parseCSV(enabled)
		return func(name string) bool { _, ok := set[name]; return ok }
	}
	if disabled != "" {
		set := parseCSV(disabled)
		return func(name string) bool { _, ok := set[name]; return !ok }
	}
	return nil
}

func parseCSV(s string) map[string]struct{} {
	set := make(map[string]struct{})
	for _, part := range strings.Split(s, ",") {
		if v := strings.TrimSpace(part); v != "" {
			set[v] = struct{}{}
		}
	}
	return set
}

type toolEntry struct {
	tool    *mcp.Tool
	handler mcp.ToolHandler
}

func RegisterAll(s *mcp.Server, c client.LeverClient, filter ToolFilter) {
	entries := []toolEntry{
		// Opportunities
		{listOpportunitiesTool(), listOpportunitiesHandler(c)},
		{getOpportunityTool(), getOpportunityHandler(c)},
		{createOpportunityTool(), createOpportunityHandler(c)},
		{listDeletedOpportunitiesTool(), listDeletedOpportunitiesHandler(c)},
		// Opportunity actions
		{archiveOpportunityTool(), archiveOpportunityHandler(c)},
		{changeOpportunityStageTool(), changeOpportunityStageHandler(c)},
		{addOpportunityTagsTool(), addOpportunityTagsHandler(c)},
		{addOpportunitySourcesTool(), addOpportunitySourcesHandler(c)},
		{addOpportunityLinksTool(), addOpportunityLinksHandler(c)},
		{removeOpportunityTagsTool(), removeOpportunityTagsHandler(c)},
		{removeOpportunitySourcesTool(), removeOpportunitySourcesHandler(c)},
		{removeOpportunityLinksTool(), removeOpportunityLinksHandler(c)},
		// Notes
		{listNotesTool(), listNotesHandler(c)},
		{createNoteTool(), createNoteHandler(c)},
		{deleteNoteTool(), deleteNoteHandler(c)},
		// Offers
		{listOffersTool(), listOffersHandler(c)},
		// Interviews
		{listInterviewsTool(), listInterviewsHandler(c)},
		{updateInterviewTool(), updateInterviewHandler(c)},
		{deleteInterviewTool(), deleteInterviewHandler(c)},
		// Feedback
		{listFeedbackTool(), listFeedbackHandler(c)},
		{getFeedbackTool(), getFeedbackHandler(c)},
		{createFeedbackTool(), createFeedbackHandler(c)},
		{updateFeedbackTool(), updateFeedbackHandler(c)},
		{deleteFeedbackTool(), deleteFeedbackHandler(c)},
		// Panels
		{createPanelTool(), createPanelHandler(c)},
		{deletePanelTool(), deletePanelHandler(c)},
		// Referrals
		{listReferralsTool(), listReferralsHandler(c)},
		// Resumes
		{listResumesTool(), listResumesHandler(c)},
		// Files
		{listFilesTool(), listFilesHandler(c)},
		// Applications
		{listApplicationsTool(), listApplicationsHandler(c)},
		{getApplicationTool(), getApplicationHandler(c)},
		// Forms
		{listFormsTool(), listFormsHandler(c)},
		// Archive reasons
		{listArchiveReasonsTool(), listArchiveReasonsHandler(c)},
		{getArchiveReasonTool(), getArchiveReasonHandler(c)},
		// Contacts
		{getContactTool(), getContactHandler(c)},
		// Postings
		{listPostingsTool(), listPostingsHandler(c)},
		{getPostingTool(), getPostingHandler(c)},
		{createPostingTool(), createPostingHandler(c)},
		{updatePostingTool(), updatePostingHandler(c)},
		{applyToPostingTool(), applyToPostingHandler(c)},
		{getPostingApplyTool(), getPostingApplyHandler(c)},
		// Users
		{listUsersTool(), listUsersHandler(c)},
		{getUserTool(), getUserHandler(c)},
		{createUserTool(), createUserHandler(c)},
		{deactivateUserTool(), deactivateUserHandler(c)},
		{reactivateUserTool(), reactivateUserHandler(c)},
		// Stages
		{listStagesTool(), listStagesHandler(c)},
		// Sources
		{listSourcesTool(), listSourcesHandler(c)},
		// Tags
		{listTagsTool(), listTagsHandler(c)},
		// Requisitions
		{listRequisitionsTool(), listRequisitionsHandler(c)},
		{getRequisitionTool(), getRequisitionHandler(c)},
		// Feedback templates
		{listFeedbackTemplatesTool(), listFeedbackTemplatesHandler(c)},
		{createFeedbackTemplateTool(), createFeedbackTemplateHandler(c)},
		{deleteFeedbackTemplateTool(), deleteFeedbackTemplateHandler(c)},
		// Form templates
		{listFormTemplatesTool(), listFormTemplatesHandler(c)},
		// Webhooks
		{listWebhooksTool(), listWebhooksHandler(c)},
		{createWebhookTool(), createWebhookHandler(c)},
		{deleteWebhookTool(), deleteWebhookHandler(c)},
		// EEO
		{listEEOResponsesTool(), listEEOResponsesHandler(c)},
		// Audit events
		{listAuditEventsTool(), listAuditEventsHandler(c)},
	}

	for _, e := range entries {
		if filter != nil && !filter(e.tool.Name) {
			continue
		}
		s.AddTool(e.tool, e.handler)
	}
}
