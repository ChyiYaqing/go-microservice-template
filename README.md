# Go Microservice Template

A production-ready Go microservice template featuring gRPC and RESTful APIs with Google API design compliance, complete with Swagger documentation.

## Features

- **Dual API Support**: Both gRPC and RESTful APIs through grpc-gateway
- **Google API Design Compliance**: Following [Google API Design Guide](https://cloud.google.com/apis/design)
- **Swagger/OpenAPI Documentation**: Auto-generated API documentation with Swagger UI
- **Protocol Buffers**: Using buf for proto management
- **Modern Go**: Built with Go 1.25
- **Production Ready**: Includes logging, configuration management, and graceful shutdown
- **Developer Friendly**: Comprehensive Makefile with common tasks

## Architecture

```
.
├── api/
│   └── proto/v1/          # Protocol buffer definitions
├── cmd/
│   └── server/            # Application entry point
├── internal/
│   ├── service/           # Business logic implementation
│   └── handler/           # Request handlers (if needed)
├── pkg/
│   ├── config/            # Configuration management
│   └── logger/            # Logging utilities
├── docs/
│   └── swagger/           # Swagger documentation
├── config/                # Configuration files
├── scripts/               # Utility scripts
├── buf.yaml               # Buf configuration
├── buf.gen.yaml           # Code generation configuration
└── Makefile               # Build automation
```

## Prerequisites

- Go 1.25 or later
- Make
- Docker (optional, for containerization)

## Quick Start

### 1. Install Required Tools

```bash
make install-tools
```

This will install:
- buf (for proto management)
- protoc-gen-go
- protoc-gen-go-grpc
- protoc-gen-grpc-gateway
- protoc-gen-openapiv2

### 2. Initialize Project

```bash
# Update module name in go.mod, buf.yaml, and source files
# Replace "github.com/yourusername/go-microservice-template" with your module path

# Install dependencies
make init
```

### 3. Generate Code from Proto Files

```bash
make proto
```

### 4. Build and Run

```bash
# Build the application
make build

# Run the application
make run

# Or run in development mode (without building)
make run-dev
```

The service will start with:
- gRPC server on port **9090**
- HTTP server on port **8080**
- Swagger UI at http://localhost:8080/swagger/

## API Documentation

### Swagger UI

Once the server is running, access the interactive API documentation at:

```
http://localhost:8080/swagger/
```

### gRPC Endpoints

The service exposes the following gRPC methods:

- `CreateUser` - Create a new user
- `GetUser` - Retrieve a user by ID
- `ListUsers` - List users with pagination
- `UpdateUser` - Update user information
- `DeleteUser` - Delete a user
- `BatchGetUsers` - Retrieve multiple users

### RESTful API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/v1/users` | Create a new user |
| GET | `/v1/users/{id}` | Get a user by ID |
| GET | `/v1/users` | List users |
| PATCH | `/v1/{user.name=users/*}` | Update a user |
| DELETE | `/v1/users/{id}` | Delete a user |
| GET | `/v1/users:batchGet` | Batch get users |

## Usage Examples

### Creating a User (RESTful API)

```bash
curl -X POST http://localhost:8080/v1/users \
  -H "Content-Type: application/json" \
  -d '{
    "user": {
      "email": "user@example.com",
      "display_name": "John Doe",
      "phone_number": "+1234567890"
    }
  }'
```

### Getting a User (RESTful API)

```bash
curl http://localhost:8080/v1/users/1
```

### Listing Users (RESTful API)

```bash
curl "http://localhost:8080/v1/users?page_size=10"
```

### Using gRPC (with grpcurl)

```bash
# List available services
grpcurl -plaintext localhost:9090 list

# Describe the UserService
grpcurl -plaintext localhost:9090 describe api.v1.UserService

# Create a user
grpcurl -plaintext -d '{
  "user": {
    "email": "user@example.com",
    "display_name": "John Doe"
  }
}' localhost:9090 api.v1.UserService/CreateUser

# Get a user
grpcurl -plaintext -d '{
  "name": "users/1"
}' localhost:9090 api.v1.UserService/GetUser
```

## Development

### Adding a New Service

1. Define your service in a new proto file under `api/proto/v1/`
2. Add Google API annotations for RESTful API mapping
3. Generate code: `make proto`
4. Implement the service in `internal/service/`
5. Register the service in `cmd/server/main.go`

### Running Tests

```bash
# Run all tests
make test

# Run tests with coverage report
make test-coverage
```

### Code Quality

```bash
# Format code
make fmt

# Run go vet
make vet

# Run linters (requires golangci-lint)
make lint

# Run all checks
make check
```

## Configuration

Configuration is managed through YAML files. The default configuration is in `config/config.yaml`:

```yaml
server:
  grpc_port: 9090
  http_port: 8080
  host: "0.0.0.0"

log:
  level: "info"
  format: "json"
```

You can create environment-specific configs (e.g., `config/production.yaml`) and pass them when starting the server:

```bash
./bin/go-microservice-template config/production.yaml
```

## Docker Support

### Build Docker Image

```bash
make docker-build
```

### Run in Docker

```bash
make docker-run
```

Or manually:

```bash
docker build -t go-microservice-template .
docker run -p 8080:8080 -p 9090:9090 go-microservice-template
```

## Google API Design Compliance

This template follows the [Google API Design Guide](https://cloud.google.com/apis/design) with:

- **Resource-oriented design**: Resources have standard methods (Create, Get, List, Update, Delete)
- **Standard fields**: Using `name`, `create_time`, `update_time` fields
- **Standard methods**: Following naming conventions (CreateUser, GetUser, etc.)
- **Pagination**: Using `page_size` and `page_token` for list methods
- **Field masks**: Supporting partial updates with `update_mask`
- **Batch operations**: Supporting batch get operations
- **RESTful mapping**: Proper HTTP verb and URL mapping through `google.api.http`

## Project Structure Best Practices

- **`api/`**: API definitions (proto files) that define the service contract
- **`cmd/`**: Application entry points
- **`internal/`**: Private application code that shouldn't be imported by other projects
- **`pkg/`**: Public libraries that can be imported by other projects
- **`docs/`**: Documentation and generated docs (Swagger)
- **`config/`**: Configuration files for different environments
- **`scripts/`**: Build and deployment scripts

## Makefile Commands

Run `make help` to see all available commands:

```bash
make help
```

Common commands:
- `make install-tools` - Install required development tools
- `make init` - Initialize project dependencies
- `make proto` - Generate code from proto files
- `make build` - Build the application
- `make run` - Build and run the application
- `make run-dev` - Run in development mode
- `make test` - Run tests
- `make clean` - Clean build artifacts
- `make docker-build` - Build Docker image

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## References

- [Google API Design Guide](https://cloud.google.com/apis/design)
- [gRPC](https://grpc.io/)
- [grpc-gateway](https://github.com/grpc-ecosystem/grpc-gateway)
- [Protocol Buffers](https://developers.google.com/protocol-buffers)
- [Buf](https://buf.build/)
