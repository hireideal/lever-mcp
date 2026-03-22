# lever-mcp

MCP server for the Lever ATS API -- 59 tools over Streamable HTTP.

![CI](https://github.com/stefanoamorelli/lever-mcp/actions/workflows/ci.yml/badge.svg)
![Go](https://img.shields.io/badge/Go-1.24-00ADD8?logo=go)
![License](https://img.shields.io/badge/License-MIT-blue)
![MCP](https://img.shields.io/badge/MCP-Streamable_HTTP-purple)

## Quick Start

| Variable | Required | Default | Description |
|---|---|---|---|
| `LEVER_API_KEY` | Yes | -- | Lever API key |
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

- `POST /mcp` -- MCP Streamable HTTP
- `GET /health` -- health check

## Tool Filtering

Restrict which tools are registered using comma-separated env vars:

```bash
# allowlist -- only these tools are registered
LEVER_ENABLED_TOOLS=list_opportunities,get_opportunity

# blocklist -- everything except these
LEVER_DISABLED_TOOLS=delete_opportunity_note,delete_webhook
```

If both are set, `LEVER_ENABLED_TOOLS` takes precedence.

## Available Tools

Legend: :eyes: read | :pencil2: create | :arrows_counterclockwise: update | :wastebasket: destructive

| Category | Tool | Annotation |
|---|---|---|
| **Opportunities** | `list_opportunities` | :eyes: |
| | `get_opportunity` | :eyes: |
| | `create_opportunity` | :pencil2: |
| | `list_deleted_opportunities` | :eyes: |
| **Opportunity Actions** | `archive_opportunity` | :arrows_counterclockwise: |
| | `change_opportunity_stage` | :arrows_counterclockwise: |
| | `add_opportunity_tags` | :pencil2: |
| | `add_opportunity_sources` | :pencil2: |
| | `add_opportunity_links` | :pencil2: |
| | `remove_opportunity_tags` | :arrows_counterclockwise: |
| | `remove_opportunity_sources` | :arrows_counterclockwise: |
| | `remove_opportunity_links` | :arrows_counterclockwise: |
| **Notes** | `list_opportunity_notes` | :eyes: |
| | `create_opportunity_note` | :pencil2: |
| | `delete_opportunity_note` | :wastebasket: |
| **Offers** | `list_opportunity_offers` | :eyes: |
| **Interviews** | `list_opportunity_interviews` | :eyes: |
| | `update_opportunity_interview` | :arrows_counterclockwise: |
| | `delete_opportunity_interview` | :wastebasket: |
| **Feedback** | `list_opportunity_feedback` | :eyes: |
| | `get_opportunity_feedback` | :eyes: |
| | `create_opportunity_feedback` | :pencil2: |
| | `update_opportunity_feedback` | :arrows_counterclockwise: |
| | `delete_opportunity_feedback` | :wastebasket: |
| **Panels** | `create_opportunity_panel` | :pencil2: |
| | `delete_opportunity_panel` | :wastebasket: |
| **Referrals** | `list_opportunity_referrals` | :eyes: |
| **Resumes** | `list_opportunity_resumes` | :eyes: |
| **Files** | `list_opportunity_files` | :eyes: |
| **Applications** | `list_opportunity_applications` | :eyes: |
| | `get_opportunity_application` | :eyes: |
| **Forms** | `list_opportunity_forms` | :eyes: |
| **Archive Reasons** | `list_archive_reasons` | :eyes: |
| | `get_archive_reason` | :eyes: |
| **Contacts** | `get_contact` | :eyes: |
| **Postings** | `list_postings` | :eyes: |
| | `get_posting` | :eyes: |
| | `create_posting` | :pencil2: |
| | `update_posting` | :arrows_counterclockwise: |
| | `apply_to_posting` | :pencil2: |
| | `get_posting_apply_questions` | :eyes: |
| **Users** | `list_users` | :eyes: |
| | `get_user` | :eyes: |
| | `create_user` | :pencil2: |
| | `deactivate_user` | :arrows_counterclockwise: |
| | `reactivate_user` | :arrows_counterclockwise: |
| **Stages** | `list_stages` | :eyes: |
| **Sources** | `list_sources` | :eyes: |
| **Tags** | `list_tags` | :eyes: |
| **Requisitions** | `list_requisitions` | :eyes: |
| | `get_requisition` | :eyes: |
| **Templates** | `list_feedback_templates` | :eyes: |
| | `create_feedback_template` | :pencil2: |
| | `delete_feedback_template` | :wastebasket: |
| | `list_form_templates` | :eyes: |
| **Webhooks** | `list_webhooks` | :eyes: |
| | `create_webhook` | :pencil2: |
| | `delete_webhook` | :wastebasket: |
| **EEO** | `list_eeo_responses` | :eyes: |
| **Audit Events** | `list_audit_events` | :eyes: |

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

MIT
