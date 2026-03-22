# Changelog

## 0.1.0

Initial release of lever-mcp.

- 59 MCP tools covering the full Lever ATS API (opportunities, postings, users, feedback, interviews, webhooks, and more)
- Streamable HTTP transport on `/mcp`
- Health check endpoint on `/health`
- Tool filtering via `LEVER_ENABLED_TOOLS` and `LEVER_DISABLED_TOOLS` env vars
- Configurable base URL for sandbox and EU regions
- Cross-platform binaries (linux, darwin, windows on amd64 and arm64)
