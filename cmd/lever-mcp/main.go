// HIREIDEAL HARDENED main.go — replaces cmd/lever-mcp/main.go
// Changes from upstream:
//   1. Added AUTH_TOKEN bearer auth middleware (timing-safe)
//   2. Added token bucket rate limiter (8 req/sec)
//   3. Version tagged as 0.1.0-hireideal
//   4. Health endpoint left unauthenticated (Railway health checks)
//   5. MCP handler registered on /, /mcp, and /sse for client compatibility

package main

import (
	"crypto/subtle"
	"encoding/json"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/stefanoamorelli/lever-mcp/internal/client"
	"github.com/stefanoamorelli/lever-mcp/internal/tools"
)

const version = "0.1.0-hireideal"

// rateLimiter implements a token bucket rate limiter.
type rateLimiter struct {
	mu         sync.Mutex
	tokens     float64
	maxTokens  float64
	refillRate float64 // tokens per second
	lastRefill time.Time
}

func newRateLimiter(maxPerSecond float64) *rateLimiter {
	return &rateLimiter{
		tokens:     maxPerSecond,
		maxTokens:  maxPerSecond,
		refillRate: maxPerSecond,
		lastRefill: time.Now(),
	}
}

func (rl *rateLimiter) Allow() bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	now := time.Now()
	elapsed := now.Sub(rl.lastRefill).Seconds()
	rl.tokens += elapsed * rl.refillRate
	if rl.tokens > rl.maxTokens {
		rl.tokens = rl.maxTokens
	}
	rl.lastRefill = now
	if rl.tokens >= 1 {
		rl.tokens--
		return true
	}
	return false
}

func authMiddleware(token string, next http.Handler) http.Handler {
	tokenBytes := []byte(token)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if !strings.HasPrefix(auth, "Bearer ") {
			http.Error(w, `{"error":"missing or invalid Authorization header"}`, http.StatusUnauthorized)
			return
		}
		provided := []byte(strings.TrimPrefix(auth, "Bearer "))
		if subtle.ConstantTimeCompare(provided, tokenBytes) != 1 {
			http.Error(w, `{"error":"invalid token"}`, http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func rateLimitMiddleware(limiter *rateLimiter, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !limiter.Allow() {
			w.Header().Set("Retry-After", "1")
			http.Error(w, `{"error":"rate limit exceeded"}`, http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func main() {
	apiKey := os.Getenv("LEVER_API_KEY")
	if apiKey == "" {
		log.Fatal("LEVER_API_KEY environment variable is required")
	}
	authToken := os.Getenv("AUTH_TOKEN")
	if authToken == "" {
		log.Fatal("AUTH_TOKEN environment variable is required \u2014 generate one with: openssl rand -hex 32")
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	listenAddr := os.Getenv("LEVER_LISTEN_ADDR")
	if listenAddr == "" {
		listenAddr = "0.0.0.0"
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

	limiter := newRateLimiter(8)
	protectedMCP := rateLimitMiddleware(limiter, authMiddleware(authToken, mcpHandler))

	mux := http.NewServeMux()
	// Register MCP handler on multiple paths for client compatibility
	// (Cowork, Claude Desktop, and other MCP clients may use different paths)
	mux.Handle("/mcp", protectedMCP)
	mux.Handle("/sse", protectedMCP)
	mux.Handle("/", protectedMCP)
	// Health endpoint is unauthenticated for Railway health checks
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "ok",
			"version": version,
		})
	})

	addr := net.JoinHostPort(listenAddr, port)
	log.Printf("lever-mcp %s listening on %s (auth: enabled, rate-limit: 8 req/sec)", version, addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}
