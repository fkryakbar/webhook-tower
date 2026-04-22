# Stage 1: Build
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Copy dependency manifests and download modules
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o webhook-tower ./cmd/webhook-tower/main.go

# Stage 2: Runtime
FROM alpine:latest

WORKDIR /app

# Install basic shell tools if needed (e.g., curl, git, bash)
RUN apk add --no-cache ca-certificates bash

# Copy the binary from builder
COPY --from=builder /app/webhook-tower .

# Expose the API port
EXPOSE 8080

# Run the application
ENTRYPOINT ["./webhook-tower"]
CMD ["--config", "/config/config.yaml"]
