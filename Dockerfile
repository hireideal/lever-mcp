FROM golang:1.25-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /lever-mcp ./cmd/lever-mcp

FROM alpine:3.21
RUN apk --no-cache add ca-certificates
COPY --from=builder /lever-mcp /lever-mcp
EXPOSE 3000
CMD ["/lever-mcp"]
