# Webhook Tower 🗼

Webhook Tower is a lightweight, headless webhook gateway and automation engine written in Go. It allows you to trigger local shell commands and scripts based on incoming HTTP webhooks with flexible matching logic and secure authentication.

## ✨ Features

- **🚀 Performance:** Built with Go and Gin-gonic for a minimal resource footprint.
- **🛠 Flexible Routing:** Define multiple routes and handlers in a single YAML configuration.
- **🎯 Precise Matching:** Filter webhooks by Path, HTTP Method, Headers, and JSON Payload values.
- **🔐 Secure:** Built-in API Key authentication (via Headers or Query Parameters).
- **📝 Variable Injection:** Inject webhook payload data directly into your commands using Go templates.
- **⚡ Hybrid Execution:** Run commands synchronously (wait for output) or asynchronously (fire-and-forget).
- **🐳 Docker Ready:** Optimized multi-stage Docker build.

---

## 🚀 Quick Start

### Prerequisites
- [Go 1.24+](https://golang.org/dl/) (if running locally)
- [Docker](https://www.docker.com/) (if deploying containerized)

### Installation

1. **Clone the repository:**
   ```bash
   git clone https://github.com/youruser/webhook-tower.git
   cd webhook-tower
   ```

2. **Setup your configuration:**
   Copy the example config and modify it to your needs.
   ```bash
   cp config.example.yaml config.yaml
   ```

3. **Run the application:**
   ```bash
   go run cmd/webhook-tower/main.go --config config.yaml
   ```

---

## ⚙️ Configuration Guide

The application is driven by a `config.yaml` file. Here is a breakdown of the structure:

```yaml
routes:
  - path: "/deploy"           # The API endpoint path
    method: "POST"            # HTTP method (POST, GET, etc.)
    api_key: "secret-123"     # Optional: X-API-Key header or api_key query param
    headers:                  # Optional: Required headers for matching
      - key: "Content-Type"
        value: "application/json"
    rules:                    # Optional: Payload condition matching
      - field: "ref"          # Field path in JSON (gjson syntax)
        operator: "=="        # Supported operators: ==
        value: "refs/heads/main"
    command:
      execute: "ls -la"       # The command to run
      async: false            # If true, returns 202 immediately
```

### Variable Injection
You can use Go template syntax `{{.fieldname}}` in the `execute` string. Webhook Tower will automatically parse the incoming JSON payload and inject the values.

**Example:**
Payload: `{"branch": "production"}`
Command: `echo deploying branch {{.branch}}`
Result: `echo deploying branch production`

---

## 🐳 Docker Deployment

### Building the Image
```bash
docker build -t webhook-tower .
```

### Running with Docker
Since Webhook Tower is designed to execute commands on the **host system**, you need to be mindful of the environment. If you want to trigger scripts on the host, you can mount the host's scripts folder or use a remote trigger.

```bash
docker run -d \
  -p 8080:8080 \
  -v $(pwd)/config.yaml:/config/config.yaml \
  webhook-tower
```

---

## 🛠 Development

### Running Tests
```bash
# Run all tests
go test -v ./...

# Check coverage
go test -cover ./...
```

### Project Structure
- `cmd/webhook-tower/`: Main entry point.
- `internal/config/`: Configuration parsing and models.
- `internal/executor/`: Shell command execution logic.
- `internal/router/`: Gin server and request matching logic.

---

## 📄 License
This project is licensed under the MIT License - see the LICENSE file for details.
