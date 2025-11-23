# Makefile Guide

This project includes a comprehensive Makefile for easy development and deployment.

## Quick Commands

### Most Common

```bash
# Development with hot reload
make dev              # or: make d

# Build the binary
make build            # or: make b

# Run the binary
make run              # or: make r

# Start with Docker Compose (recommended)
make docker-compose-up    # or: make dcu
```

### Quick References

```bash
make                  # Shows this help message
make help             # Full help documentation
```

## Development Commands

### Building & Running

```bash
# Build locally
make build
./relay

# Run directly (builds if needed)
make run

# Run with hot reload (auto-restarts on file changes)
make dev
make d  # shorthand
```

### Setup

```bash
# Initial project setup (downloads deps, installs air)
make setup

# Download dependencies
make mod-download

# Tidy modules
make mod-tidy

# Install air hot reload tool
make install-air
```

## Docker Commands

### Docker Compose (Recommended for Development)

```bash
# Start services with hot reload
make docker-compose-up
make dcu              # shorthand

# Stop services
make docker-compose-down
make dcd              # shorthand
```

### Docker Standalone

```bash
# Build production image
make docker-build

# Run in container
make docker-run

# Stop containers
make docker-stop
```

## Code Quality Commands

### Format & Lint

```bash
# Format code with gofmt
make fmt

# Run go vet
make vet

# Run linter (installs if needed)
make lint

# All checks: format, vet, tidy, build
make all
```

## Testing Commands

### Health & API Testing

```bash
# Check if service is running
make health-check
make hc               # shorthand

# Send test relay request
make test-relay
make tr               # shorthand
```

### Example Output

```bash
$ make test-relay
Sending test relay request...
{
  "success": false,
  "error": "CONTRACT_ADDRESS is required"
}
```

## Information Commands

```bash
# Show Go version
make version

# Show .env variables
make env
```

## Cleaning

```bash
# Remove build artifacts
make clean
```

## Workflow Examples

### Local Development (with hot reload)

```bash
# First time setup
make setup

# Edit .env with your config
nano .env

# Start developing
make dev
```

### Docker Development (recommended)

```bash
# Setup
make setup

# Edit .env
nano .env

# Start with hot reload
make docker-compose-up
```

### Production Deployment

```bash
# Build production image
make docker-build

# Run container
make docker-run
```

### Before Committing

```bash
# Run all checks
make all

# Or individual checks:
make fmt
make vet
make lint
make build
```

## Environment Variables

Control behavior with environment variables:

```bash
# Custom port
make health-check PORT=9000

# Custom image name
make docker-build DOCKER_IMAGE_NAME=my-relay DOCKER_IMAGE_TAG=v1.0
```

## All Available Commands

| Command | Description | Shorthand |
|---------|-------------|-----------|
| `help` | Show help message | - |
| `build` | Build binary | `b` |
| `run` | Run binary | `r` |
| `dev` | Run with hot reload | `d` |
| `docker-build` | Build Docker image | - |
| `docker-run` | Run Docker container | - |
| `docker-stop` | Stop Docker containers | - |
| `docker-compose-up` | Start with docker-compose | `dcu` |
| `docker-compose-down` | Stop docker-compose | `dcd` |
| `fmt` | Format code | - |
| `vet` | Run go vet | - |
| `lint` | Run linter | - |
| `install-air` | Install air tool | - |
| `mod-download` | Download dependencies | - |
| `mod-tidy` | Tidy modules | - |
| `health-check` | Check service health | `hc` |
| `test-relay` | Send test request | `tr` |
| `version` | Show Go version | - |
| `env` | Show .env variables | - |
| `clean` | Clean artifacts | - |
| `setup` | Initial setup | - |
| `all` | Full quality check + build | - |

## Tips

1. **Use `make dev` for local development** - It auto-reloads on file changes
2. **Use `make docker-compose-up` for clean environment** - Matches production
3. **Run `make all` before committing** - Ensures code quality
4. **Use shorthand commands** - `make d`, `make b`, `make r`, etc. are faster

## Troubleshooting

### `make: command not found`
- Ensure you're in the relay directory
- Check that Makefile exists: `ls -la Makefile`

### `.env file not found` error
- Run: `cp .env.example .env`
- Edit with your configuration

### `air not found`
- Run: `make install-air`
- Or manually: `go install github.com/cosmtrek/air@latest`

### Docker build fails
- Check docker is installed: `docker --version`
- Check docker daemon is running

### golangci-lint not found
- The Makefile will auto-install it
- Or manually: `go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest`
