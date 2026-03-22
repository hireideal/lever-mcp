package main

import (
	"encoding/json"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/stefanoamorelli/lever-mcp/internal/client"
	"github.com/stefanoamorelli/lever-mcp/internal/tools"
)

const version = "0.1.0"

func main() {
	apiKey := os.Getenv("LEVER_API_KEY")
	if apiKey == "" {
		log.Fatal("LEVER_API_KEY environment variable is required")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	baseURL := os.Getenv("LEVER_BASE_URL")
	var opts []client.Option
	if baseURL != "" {
		opts = append(opts, client.WithBaseURL(baseURL))
	}

	c := client.New(apiKey, opts...)
	server := mcp.NewServer(&mcp.Implementation{
		Name:    "lever-mcp",
		Version: version,
	}, nil)

	filter := tools.NewToolFilter(
		os.Getenv("LEVER_ENABLED_TOOLS"),
		os.Getenv("LEVER_DISABLED_TOOLS"),
	)
	tools.RegisterAll(server, c, filter)

	mcpHandler := mcp.NewStreamableHTTPHandler(func(r *http.Request) *mcp.Server {
		return server
	}, nil)

	mux := http.NewServeMux()
	mux.Handle("/mcp", mcpHandler)
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "ok",
			"version": version,
		})
	})

	addr := net.JoinHostPort("", port)
	log.Printf("lever-mcp %s listening on %s", version, addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}
