<h1 align="center">lever-mcp</h1>

<p align="center">
  Community-developed MCP server for the <a href="https://hire.lever.co/developer/documentation">Lever ATS API</a>, exposing 59 tools over Streamable HTTP.
</p>

<p align="center">
  <img src="https://github.com/stefanoamorelli/lever-mcp/actions/workflows/ci.yml/badge.svg" alt="CI">
  <img src="https://img.shields.io/badge/Go-1.24-00ADD8?logo=go" alt="Go">
  <img src="https://img.shields.io/badge/License-AGPL--3.0-blue" alt="License">
  <img src="https://img.shields.io/badge/MCP-Streamable_HTTP-purple" alt="MCP">
</p>

> [!IMPORTANT]
> This project is **experimental** and provided as-is. It is **not** endorsed by, affiliated with, or supported by [Lever](https://www.lever.co) or [Employ, Inc.](https://www.employinc.com) "Lever" is a trademark of [Employ, Inc.](https://www.employinc.com) All product names, logos, and brands are property of their respective owners. Use at your own risk.

## Quick Start

| Variable | Required | Default | Description |
|---|---|---|---|
| `LEVER_API_KEY` | Yes | | Lever API key |
| `PORT` | No | `3000` | HTTP listen port |
| `LEVER_BASE_URL` | No | `https://api.lever.co/v1` | Override for sandbox or EU regions |

## Build & Run

```bash
# build
make build          # or: go build -o lever-mcp ./cmd/lever-mcp

# run
export LEVER_API_KEY=your-key
./lever-mcp         # listening on :3000
```

The server exposes two endpoints:

- `POST /mcp` (MCP Streamable HTTP)
- `GET /health` (health check)

## Tool Filtering

Restrict which tools are registered using comma-separated env vars:

```bash
# allowlist: only these tools are registered
LEVER_ENABLED_TOOLS=list_opportunities,get_opportunity

# blocklist: everything except these
LEVER_DISABLED_TOOLS=delete_opportunity_note,delete_webhook
```

If both are set, `LEVER_ENABLED_TOOLS` takes precedence.

## Available Tools (59)

<sub>:eyes: read · :pencil2: create · :arrows_counterclockwise: update · :wastebasket: destructive</sub>

<details>
<summary><strong>Opportunities</strong> (4 tools)</summary>

| Tool | Annotation | API Docs |
|---|---|---|
| `list_opportunities` | :eyes: | [docs](https://hire.lever.co/developer/documentation#list-all-opportunities) |
| `get_opportunity` | :eyes: | [docs](https://hire.lever.co/developer/documentation#retrieve-a-single-opportunity) |
| `create_opportunity` | :pencil2: | [docs](https://hire.lever.co/developer/documentation#create-an-opportunity) |
| `list_deleted_opportunities` | :eyes: | [docs](https://hire.lever.co/developer/documentation#list-deleted-opportunities) |

</details>

<details>
<summary><strong>Opportunity Actions</strong> (8 tools)</summary>

| Tool | Annotation | API Docs |
|---|---|---|
| `archive_opportunity` | :arrows_counterclockwise: | [docs](https://hire.lever.co/developer/documentation#update-opportunity-archived-state) |
| `change_opportunity_stage` | :arrows_counterclockwise: | [docs](https://hire.lever.co/developer/documentation#update-opportunity-stage) |
| `add_opportunity_tags` | :pencil2: | [docs](https://hire.lever.co/developer/documentation#update-opportunity-tags) |
| `add_opportunity_sources` | :pencil2: | [docs](https://hire.lever.co/developer/documentation#update-opportunity-sources) |
| `add_opportunity_links` | :pencil2: | [docs](https://hire.lever.co/developer/documentation#update-contact-links-by-opportunity) |
| `remove_opportunity_tags` | :arrows_counterclockwise: | [docs](https://hire.lever.co/developer/documentation#update-opportunity-tags) |
| `remove_opportunity_sources` | :arrows_counterclockwise: | [docs](https://hire.lever.co/developer/documentation#update-opportunity-sources) |
| `remove_opportunity_links` | :arrows_counterclockwise: | [docs](https://hire.lever.co/developer/documentation#update-contact-links-by-opportunity) |

</details>

<details>
<summary><strong>Notes</strong> (3 tools)</summary>

| Tool | Annotation | API Docs |
|---|---|---|
| `list_opportunity_notes` | :eyes: | [docs](https://hire.lever.co/developer/documentation#list-all-notes) |
| `create_opportunity_note` | :pencil2: | [docs](https://hire.lever.co/developer/documentation#create-a-note) |
| `delete_opportunity_note` | :wastebasket: | [docs](https://hire.lever.co/developer/documentation#delete-a-note) |

</details>

<details>
<summary><strong>Offers</strong> (1 tool)</summary>

| Tool | Annotation | API Docs |
|---|---|---|
| `list_opportunity_offers` | :eyes: | [docs](https://hire.lever.co/developer/documentation#list-all-offers) |

</details>

<details>
<summary><strong>Interviews</strong> (3 tools)</summary>

| Tool | Annotation | API Docs |
|---|---|---|
| `list_opportunity_interviews` | :eyes: | [docs](https://hire.lever.co/developer/documentation#list-all-interviews) |
| `update_opportunity_interview` | :arrows_counterclockwise: | [docs](https://hire.lever.co/developer/documentation#update-an-interview) |
| `delete_opportunity_interview` | :wastebasket: | [docs](https://hire.lever.co/developer/documentation#delete-an-interview) |

</details>

<details>
<summary><strong>Feedback</strong> (5 tools)</summary>

| Tool | Annotation | API Docs |
|---|---|---|
| `list_opportunity_feedback` | :eyes: | [docs](https://hire.lever.co/developer/documentation#list-all-feedback) |
| `get_opportunity_feedback` | :eyes: | [docs](https://hire.lever.co/developer/documentation#retrieve-a-feedback-form) |
| `create_opportunity_feedback` | :pencil2: | [docs](https://hire.lever.co/developer/documentation#create-a-feedback-form) |
| `update_opportunity_feedback` | :arrows_counterclockwise: | [docs](https://hire.lever.co/developer/documentation#update-feedback) |
| `delete_opportunity_feedback` | :wastebasket: | [docs](https://hire.lever.co/developer/documentation#delete-feedback) |

</details>

<details>
<summary><strong>Panels</strong> (2 tools)</summary>

| Tool | Annotation | API Docs |
|---|---|---|
| `create_opportunity_panel` | :pencil2: | [docs](https://hire.lever.co/developer/documentation#create-a-panel) |
| `delete_opportunity_panel` | :wastebasket: | [docs](https://hire.lever.co/developer/documentation#delete-a-panel) |

</details>

<details>
<summary><strong>Referrals, Resumes, Files</strong> (3 tools)</summary>

| Tool | Annotation | API Docs |
|---|---|---|
| `list_opportunity_referrals` | :eyes: | [docs](https://hire.lever.co/developer/documentation#list-all-referrals) |
| `list_opportunity_resumes` | :eyes: | [docs](https://hire.lever.co/developer/documentation#list-all-resumes) |
| `list_opportunity_files` | :eyes: | [docs](https://hire.lever.co/developer/documentation#list-all-files) |

</details>

<details>
<summary><strong>Applications</strong> (2 tools)</summary>

| Tool | Annotation | API Docs |
|---|---|---|
| `list_opportunity_applications` | :eyes: | [docs](https://hire.lever.co/developer/documentation#list-all-applications) |
| `get_opportunity_application` | :eyes: | [docs](https://hire.lever.co/developer/documentation#retrieve-a-single-application) |

</details>

<details>
<summary><strong>Forms</strong> (1 tool)</summary>

| Tool | Annotation | API Docs |
|---|---|---|
| `list_opportunity_forms` | :eyes: | [docs](https://hire.lever.co/developer/documentation#list-all-profile-forms) |

</details>

<details>
<summary><strong>Archive Reasons</strong> (2 tools)</summary>

| Tool | Annotation | API Docs |
|---|---|---|
| `list_archive_reasons` | :eyes: | [docs](https://hire.lever.co/developer/documentation#list-all-archive-reasons) |
| `get_archive_reason` | :eyes: | [docs](https://hire.lever.co/developer/documentation#retrieve-a-single-archive-reason) |

</details>

<details>
<summary><strong>Contacts</strong> (1 tool)</summary>

| Tool | Annotation | API Docs |
|---|---|---|
| `get_contact` | :eyes: | [docs](https://hire.lever.co/developer/documentation#retrieve-a-single-contact) |

</details>

<details>
<summary><strong>Postings</strong> (6 tools)</summary>

| Tool | Annotation | API Docs |
|---|---|---|
| `list_postings` | :eyes: | [docs](https://hire.lever.co/developer/documentation#list-all-postings) |
| `get_posting` | :eyes: | [docs](https://hire.lever.co/developer/documentation#retrieve-a-single-posting) |
| `create_posting` | :pencil2: | [docs](https://hire.lever.co/developer/documentation#create-a-posting) |
| `update_posting` | :arrows_counterclockwise: | [docs](https://hire.lever.co/developer/documentation#update-a-posting) |
| `apply_to_posting` | :pencil2: | [docs](https://hire.lever.co/developer/documentation#apply-to-a-posting) |
| `get_posting_apply_questions` | :eyes: | [docs](https://hire.lever.co/developer/documentation#retrieve-posting-application-questions) |

</details>

<details>
<summary><strong>Users</strong> (5 tools)</summary>

| Tool | Annotation | API Docs |
|---|---|---|
| `list_users` | :eyes: | [docs](https://hire.lever.co/developer/documentation#list-all-users) |
| `get_user` | :eyes: | [docs](https://hire.lever.co/developer/documentation#retrieve-a-single-user) |
| `create_user` | :pencil2: | [docs](https://hire.lever.co/developer/documentation#create-a-user) |
| `deactivate_user` | :arrows_counterclockwise: | [docs](https://hire.lever.co/developer/documentation#deactivate-a-user) |
| `reactivate_user` | :arrows_counterclockwise: | [docs](https://hire.lever.co/developer/documentation#reactivate-a-user) |

</details>

<details>
<summary><strong>Stages, Sources, Tags</strong> (3 tools)</summary>

| Tool | Annotation | API Docs |
|---|---|---|
| `list_stages` | :eyes: | [docs](https://hire.lever.co/developer/documentation#list-all-stages) |
| `list_sources` | :eyes: | [docs](https://hire.lever.co/developer/documentation#list-all-sources) |
| `list_tags` | :eyes: | [docs](https://hire.lever.co/developer/documentation#list-all-tags) |

</details>

<details>
<summary><strong>Requisitions</strong> (2 tools)</summary>

| Tool | Annotation | API Docs |
|---|---|---|
| `list_requisitions` | :eyes: | [docs](https://hire.lever.co/developer/documentation#list-all-requisitions) |
| `get_requisition` | :eyes: | [docs](https://hire.lever.co/developer/documentation#retrieve-a-single-requisition) |

</details>

<details>
<summary><strong>Templates</strong> (4 tools)</summary>

| Tool | Annotation | API Docs |
|---|---|---|
| `list_feedback_templates` | :eyes: | [docs](https://hire.lever.co/developer/documentation#list-all-feedback-templates) |
| `create_feedback_template` | :pencil2: | [docs](https://hire.lever.co/developer/documentation#create-a-feedback-template) |
| `delete_feedback_template` | :wastebasket: | [docs](https://hire.lever.co/developer/documentation#delete-a-feedback-template) |
| `list_form_templates` | :eyes: | [docs](https://hire.lever.co/developer/documentation#list-all-profile-form-templates) |

</details>

<details>
<summary><strong>Webhooks</strong> (3 tools)</summary>

| Tool | Annotation | API Docs |
|---|---|---|
| `list_webhooks` | :eyes: | [docs](https://hire.lever.co/developer/documentation#list-webhooks) |
| `create_webhook` | :pencil2: | [docs](https://hire.lever.co/developer/documentation#create-a-webhook) |
| `delete_webhook` | :wastebasket: | [docs](https://hire.lever.co/developer/documentation#delete-a-webhook) |

</details>

<details>
<summary><strong>EEO</strong> (1 tool)</summary>

| Tool | Annotation | API Docs |
|---|---|---|
| `list_eeo_responses` | :eyes: | [docs](https://hire.lever.co/developer/documentation#eeo) |

</details>

<details>
<summary><strong>Audit Events</strong> (1 tool)</summary>

| Tool | Annotation | API Docs |
|---|---|---|
| `list_audit_events` | :eyes: | [docs](https://hire.lever.co/developer/documentation#list-all-audit-events) |

</details>

## Client Configuration

### Claude Desktop

```json
{
  "mcpServers": {
    "lever": {
      "type": "streamable-http",
      "url": "http://localhost:3000/mcp"
    }
  }
}
```

### Generic MCP Client

Any client that supports Streamable HTTP transport can connect to `http://<host>:3000/mcp`.

## Development

```bash
make test                                              # run all tests
make build                                             # compile binary
go test ./internal/tools/ -run TestConformance         # annotation & schema conformance
```

## API Reference

[Lever API Documentation](https://hire.lever.co/developer/documentation)

## License

[AGPL-3.0](LICENSE)

Copyright (c) 2026 [Stefano Amorelli](https://amorelli.tech) (stefano@amorelli.tech)
