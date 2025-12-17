# Quick Start Guide

This guide will help you get the microservice up and running in minutes.

## Prerequisites

- Go 1.25 or later
- Make

## Step 1: Run the Setup Script

The easiest way to get started is to run the automated setup script:

```bash
bash scripts/setup.sh
```

This script will:
- Check your Go installation
- Install all required tools (buf, protoc plugins, etc.)
- Download Go dependencies
- Generate code from proto files

## Step 2: Update Module Path

Before building, update the module path in the following files:

1. **go.mod**: Change `github.com/yourusername/go-microservice-template` to your actual module path
2. **buf.yaml**: Update the `name` field
3. **All .go files**: Update import paths to match your module

You can use `find` and `sed` to automate this:

```bash
# Replace with your actual module path
NEW_MODULE="github.com/yourname/yourproject"

# Update go.mod
sed -i '' "s|github.com/yourusername/go-microservice-template|$NEW_MODULE|g" go.mod

# Update all Go files
find . -name "*.go" -type f -exec sed -i '' "s|github.com/yourusername/go-microservice-template|$NEW_MODULE|g" {} +

# Update buf.yaml
sed -i '' "s|yourusername/go-microservice-template|yourname/yourproject|g" buf.yaml

# Update buf.gen.yaml
sed -i '' "s|github.com/yourusername/go-microservice-template|$NEW_MODULE|g" buf.gen.yaml
```

## Step 3: Build and Run

```bash
# Build the application
make build

# Run the application
make run
```

Or run directly in development mode:

```bash
make run-dev
```

## Step 4: Test the APIs

The server will start with:
- **gRPC** on port `9090`
- **HTTP/REST** on port `8080`
- **Swagger UI** at http://localhost:8080/swagger/

### Test REST API

Create a user:
```bash
curl -X POST http://localhost:8080/v1/users \
  -H "Content-Type: application/json" \
  -d '{
    "user": {
      "email": "alice@example.com",
      "display_name": "Alice Smith",
      "phone_number": "+1234567890"
    }
  }'
```

Get a user:
```bash
curl http://localhost:8080/v1/users/1
```

List users:
```bash
curl http://localhost:8080/v1/users
```

### Test gRPC API

Install grpcurl first:
```bash
go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest
```

List services:
```bash
make grpcurl-list
```

Create a user via gRPC:
```bash
grpcurl -plaintext -d '{
  "user": {
    "email": "bob@example.com",
    "display_name": "Bob Johnson"
  }
}' localhost:9090 api.v1.UserService/CreateUser
```

## Step 5: View API Documentation

Open your browser and navigate to:

```
http://localhost:8080/swagger/
```

You'll see the interactive Swagger UI with all available endpoints.

## Common Commands

```bash
# Install tools
make install-tools

# Generate proto files
make proto

# Build
make build

# Run
make run

# Run in dev mode
make run-dev

# Run tests
make test

# Format code
make fmt

# Clean build artifacts
make clean

# Show all available commands
make help
```

## Docker Quick Start

If you prefer to use Docker:

```bash
# Build Docker image
make docker-build

# Run with Docker
make docker-run
```

Or with docker-compose:

```bash
docker-compose up
```

## Next Steps

1. **Customize the User Service**: Modify `api/proto/v1/user.proto` to fit your needs
2. **Add More Services**: Create new proto files and implement services
3. **Add Database**: Integrate your preferred database (PostgreSQL, MongoDB, etc.)
4. **Add Middleware**: Implement authentication, rate limiting, etc.
5. **Deploy**: Use the Dockerfile to deploy to your cloud provider

## Troubleshooting

### "command not found: buf"

Run `make install-tools` to install all required tools.

### Import errors after generating proto files

Make sure you've updated all import paths to match your module name.

### Port already in use

If ports 8080 or 9090 are already in use, modify `config/config.yaml`:

```yaml
server:
  grpc_port: 9091  # Change this
  http_port: 8081  # Change this
  host: "0.0.0.0"
```

## Support

For more detailed information, see [README.md](README.md).
