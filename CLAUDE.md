# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

"Gopher Post" is an automated RSS aggregation and email delivery service where a gopher delivers your daily mail. The application automatically collects the latest articles from RSS feeds specified in feeds.yml and delivers them via email to subscribers using Korean-localized HTML templates. It's designed to run as a containerized service that processes feeds and delivers new posts to users' inboxes daily.

The project uses Go 1.24.0 and follows a standard Go module structure with Docker containerization for deployment.

## Project Structure

```
gopher-post/
├── cmd/
│   └── main.go              # Main application entry point with .env loading
├── internal/
│   ├── config/
│   │   └── config.go        # Configuration management (feeds.yml, recipients.yml, SMTP)
│   ├── feed/
│   │   └── parser.go        # RSS feed parsing using gofeed library
│   ├── email/
│   │   └── smtp.go          # SMTP email service using gopkg.in/mail.v2
│   └── template/
│       └── template.go      # HTML template rendering service
├── templates/
│   └── newsletter.html      # Korean HTML email template with Noto Sans KR
├── feeds.yml                # RSS feed configuration (name, url, rss)
├── recipients.yml           # Email recipients configuration
├── .env                     # Environment variables (SMTP credentials)
├── .env.example             # Example environment variables template
├── Dockerfile               # Multi-stage Docker build
├── go.mod                   # Go module with dependencies
├── go.sum                   # Go module checksums
├── README.md                # Korean project documentation
└── CLAUDE.md                # This file
```

## Common Commands

### Building and Running
- `go run cmd/main.go` - Run the application directly (requires .env file)
- `go build -o gopher-post cmd/main.go` - Build binary executable
- `go build cmd/main.go` - Build with default binary name

### Docker Commands
- `docker build -t gopher-post .` - Build Docker image
- `docker run --env-file .env gopher-post` - Run container with environment file
- `docker run -e SMTP_HOST=... -e SMTP_USERNAME=... gopher-post` - Run with inline env vars

### Development
- `go fmt ./...` - Format all Go files
- `go vet ./...` - Run Go vet to catch common mistakes
- `go mod tidy` - Clean up module dependencies
- `go mod download` - Download dependencies

### Testing
- `go test ./...` - Run all tests (currently no tests exist)
- `go test -v ./...` - Run tests with verbose output

## Architecture Notes

This is an automated RSS aggregation and email delivery service with the following components:

### Core Functionality
- **Feed Processing**: Reads RSS and Atom feeds from feeds.yml configuration
- **Content Collection**: Automatically collects latest articles from feeds (24-hour window)
- **Email Delivery**: Sends new posts to subscribers via email with Korean HTML templates
- **Environment Configuration**: Uses .env files and environment variables for SMTP settings
- **Scheduling**: Designed for daily automated runs via cron or container orchestration

### Key Components (implemented)
- **Feed Parser** (`internal/feed/parser.go`): Uses `gofeed` library for RSS and Atom parsing
- **Email Service** (`internal/email/smtp.go`): SMTP integration using `gopkg.in/mail.v2`
- **Template Engine** (`internal/template/template.go`): Korean HTML email template rendering
- **Configuration** (`internal/config/config.go`): YAML parsing and environment variable management
- **Main Application** (`cmd/main.go`): Orchestrates all services with .env loading

### Dependencies
- RSS/Atom parsing: `github.com/mmcdole/gofeed v1.2.1`
- Email/SMTP: `gopkg.in/mail.v2 v2.3.1`
- YAML parsing: `gopkg.in/yaml.v3 v3.0.1`
- Environment variables: `github.com/joho/godotenv v1.5.1`
- HTML template processing: Built-in Go `html/template`

### Configuration Files

#### feeds.yml Structure
```yaml
feeds:
  - name: "Feed Name"
    url: "https://website.com"
    rss: "https://website.com/feed.xml"
```

#### recipients.yml Structure
```yaml
recipients:
  - email: "user@example.com"
```

#### Required Environment Variables
- `SMTP_HOST` - SMTP server hostname (default: smtp.gmail.com)
- `SMTP_PORT` - SMTP server port (default: 587)
- `SMTP_USERNAME` - SMTP authentication username
- `SMTP_PASSWORD` - SMTP authentication password
- `FROM_EMAIL` - Sender email address

### Email Template Features
- Korean language interface
- Responsive HTML design with Noto Sans KR font
- Modern gradient styling with container-based layout
- Article cards with feed source and publication date
- "전체 글 읽기" (Read Full Article) links
- Footer with service description and unsubscribe links

### Deployment Strategy
- Docker containerization for consistent environments
- Environment-based configuration for secrets (SMTP credentials)
- Multi-stage Docker build for optimized container size
- Volume mounting for configuration files (feeds.yml, recipients.yml)
- Suitable for scheduled execution via cron or container orchestration

### Current Implementation Status
- ✅ RSS and Atom feed parsing with date filtering (24-hour window)
- ✅ SMTP email delivery with HTML templates
- ✅ Korean localization and responsive design
- ✅ Docker containerization
- ✅ Environment variable configuration
- ✅ YAML configuration file support
- ❌ Web interface for feed/recipient management
- ❌ Database storage for feeds and recipients
- ❌ Email subscription/unsubscription handling
- ❌ Automated testing suite