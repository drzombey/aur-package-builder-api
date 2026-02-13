# AUR Package Builder API 🏗️

[![Go version](https://img.shields.io/badge/Go-1.24-blue.svg)](https://golang.org/doc/go1.24)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://github.com/drzombey/aur-package-builder-api/blob/master/LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/drzombey/aur-package-builder-api)](https://goreportcard.com/report/github.com/drzombey/aur-package-builder-api)
[![Build Status](https://github.com/drzombey/aur-package-builder-api/workflows/Build%20go%20project/badge.svg)](https://github.com/drzombey/aur-package-builder-api/actions)

> A Go API service that automatically builds Arch Linux User Repository (AUR) packages in isolated Docker containers.

## 🏛️ Architecture Overview

The AUR Package Builder API is designed with a clean architecture approach:

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   HTTP Client   │───▶│   REST API      │───▶│   Handler       │
└─────────────────┘    └─────────────────┘    └─────────────────┘
                                                        │
                       ┌─────────────────┐    ┌─────────────────┐
                       │  Builder Service│◀───│ Repository Layer│
                       └─────────────────┘    └─────────────────┘
                                │                       │
                       ┌─────────────────┐    ┌─────────────────┐
                       │ Docker Service  │    │    MongoDB      │
                       └─────────────────┘    └─────────────────┘
                                │
                       ┌─────────────────┐
                       │ Storage Provider│
                       │   (S3/Local)    │
                       └─────────────────┘
```

**Key Components:**
- **Handler Layer**: HTTP request handling with dependency injection
- **Repository Layer**: Data persistence abstraction (MongoDB)
- **Builder Service**: Docker-based package building orchestration  
- **Storage Provider**: File storage abstraction (S3, local filesystem)
- **Docker Service**: Container management for isolated builds

## 🚀 Features

- **Automated AUR Package Building**: Build packages directly from AUR sources
- **Container Isolation**: Each build runs in a clean Docker environment
- **Multiple Storage Backends**: Support for S3 and local filesystem storage
- **Package Management**: Track built packages with MongoDB
- **OTLP Tracing**: Comprehensive distributed tracing support
- **RESTful API**: Clean HTTP API for package operations
- **Configuration Management**: Environment-based configuration

## 📋 Prerequisites

- **Go 1.24+**
- **Docker** (for package building)
- **MongoDB** (for package metadata)
- **Optional**: S3-compatible storage (for package storage)

## ⚙️ Configuration

The application uses environment variables with the `AUR_BUILDER_` prefix:

```bash
# Server Configuration
AUR_BUILDER_WEBSERVER_PORT=8080
AUR_BUILDER_DEBUG=false
AUR_BUILDER_LOG_LEVEL=info

# Storage Configuration  
AUR_BUILDER_STORAGE_PROVIDER=filesystem  # or 's3'
AUR_BUILDER_PACKAGE_PATH=/var/packages

# Database Configuration
AUR_BUILDER_DATABASE_URI=mongodb://localhost:27017
AUR_BUILDER_DATABASE_NAME=aur_builder

# Docker Configuration
AUR_BUILDER_DOCKER_CONTAINER_IMAGE=aurbuilder:latest
AUR_BUILDER_DOCKER_USERNAME=your_registry_user
AUR_BUILDER_DOCKER_PASSWORD=your_registry_password
AUR_BUILDER_DOCKER_AUTH=false

# S3 Configuration (if using S3 storage)
AUR_BUILDER_BUCKET_ENDPOINT=s3.amazonaws.com
AUR_BUILDER_BUCKET_ACCESS_KEY_ID=your_access_key
AUR_BUILDER_BUCKET_SECRET_ACCESS_KEY=your_secret_key
AUR_BUILDER_BUCKET_NAME=aur-packages
AUR_BUILDER_BUCKET_REGION=us-east-1
AUR_BUILDER_BUCKET_USE_SSL=true

# OTLP Tracing Configuration
AUR_BUILDER_OTLP_ENDPOINT=http://localhost:4318
```

## 🛠️ Installation & Setup

### Development Setup

1. **Clone the repository**
```bash
git clone https://github.com/drzombey/aur-package-builder-api
cd aur-package-builder-api
```

2. **Install dependencies**
```bash
go mod tidy
```

3. **Setup configuration**
```bash
cd cmd/api && cp config.sample.yml config.yml
# Edit config.yml with your settings
```

4. **Start MongoDB (using Docker)**
```bash
docker run -d --name mongodb -p 27017:27017 mongo:latest
```

### Native Execution

```bash
cd cmd/api && go run main.go
```

### Development with Air (Hot Reload)

Install [Air](https://github.com/cosmtrek/air) for hot reloading:
```bash
go install github.com/cosmtrek/air@latest
```

Then run:
```bash
air
```

### Docker Setup

**Build the image:**
```bash
docker build -t aur-package-builder-api .
```

**Run with Docker Compose:**
```bash
docker-compose up
```

The `docker-compose.yml` includes:
- API service
- MongoDB database
- Volume mounts for package storage

## 🔌 API Endpoints

### Package Management

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/packages` | Get all built packages |
| `GET` | `/aur-package?packageName={name}` | Search AUR packages |
| `POST` | `/build` | Build a new package |

### Example Usage

**Search for a package:**
```bash
curl "http://localhost:8080/aur-package?packageName=yay"
```

**Build a package:**
```bash
curl -X POST http://localhost:8080/build \
  -H "Content-Type: application/json" \
  -d '{
    "ID": 12345,
    "Name": "yay",
    "PackageBase": "yay", 
    "Version": "12.0.5-1"
  }'
```

**List built packages:**
```bash
curl http://localhost:8080/packages
```

## 🧪 Testing

Run all tests:
```bash
go test ./...
```

Run tests with coverage:
```bash
go test -v -cover ./...
```

Run specific package tests:
```bash
go test -v ./cmd/api/handler
```

**Note**: Some integration tests require Docker to be running.

## 🏗️ Building

**Build for current platform:**
```bash
go build -o aur-package-builder-api cmd/api/main.go
```

**Cross-compilation example:**
```bash
GOOS=linux GOARCH=amd64 go build -o aur-package-builder-api-linux-amd64 cmd/api/main.go
```

## 🚀 Deployment

### Docker Deployment

1. **Build production image:**
```bash
docker build -t aur-package-builder-api:latest .
```

2. **Run with environment variables:**
```bash
docker run -d \
  --name aur-builder \
  -p 8080:8080 \
  -e AUR_BUILDER_DATABASE_URI=mongodb://your-mongo-host:27017 \
  -e AUR_BUILDER_STORAGE_PROVIDER=s3 \
  -e AUR_BUILDER_BUCKET_NAME=your-s3-bucket \
  -v /var/run/docker.sock:/var/run/docker.sock \
  aur-package-builder-api:latest
```

### Binary Deployment

1. Download release binary from GitHub releases
2. Create configuration file or set environment variables
3. Run: `./aur-package-builder-api`

## 📊 Monitoring & Tracing

The service supports OpenTelemetry (OTLP) tracing for observability:

- **Jaeger**: For trace collection and visualization
- **Prometheus**: For metrics (via OTEL)
- **Grafana**: For dashboards and monitoring

Configure OTLP endpoint:
```bash
AUR_BUILDER_OTLP_ENDPOINT=http://your-otlp-collector:4318
```

## 🤝 Contributing

We welcome contributions! Here's how to get started:

1. **Fork the repository**
2. **Create a feature branch**: `git checkout -b feature/amazing-feature`
3. **Make your changes**
4. **Add tests** for your changes
5. **Run tests**: `go test ./...`
6. **Run linting**: `golangci-lint run`
7. **Commit changes**: `git commit -m 'Add amazing feature'`
8. **Push to branch**: `git push origin feature/amazing-feature`
9. **Open a Pull Request**

### Development Guidelines

- **Code Style**: Follow Go conventions and use `gofmt`
- **Testing**: Write tests for new functionality
- **Documentation**: Update README and code comments
- **Commit Messages**: Use conventional commit format
- **Breaking Changes**: Document in PR description

### Setting up Development Environment

```bash
# Install development tools
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
go install github.com/cosmtrek/air@latest

# Run linter
golangci-lint run

# Run tests with coverage
go test -v -cover ./...
```

## 📄 License

This project is [MIT](https://github.com/drzombey/aur-package-builder-api/blob/master/LICENSE) licensed.

## 👥 Authors

**Tim Lange**
- Github: [@drzombey](https://github.com/drzombey)

**Niklas Dörfler** (Contributor)
- Github: [@niklasdoerfler](https://github.com/niklasdoerfler)

## 🙏 Acknowledgments

- [Arch Linux AUR](https://aur.archlinux.org/) for the package repository
- [Gin Web Framework](https://gin-gonic.com/) for HTTP handling
- [Docker](https://www.docker.com/) for containerization
- [MongoDB](https://www.mongodb.com/) for data persistence

---

⭐ **Star this repository if it helped you!**

> **Need help?** Check [issues](https://github.com/drzombey/aur-package-builder-api/issues) or create a new one.