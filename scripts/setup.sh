#!/bin/bash

set -e

echo "🚀 Setting up Go Microservice Template..."

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed. Please install Go 1.25 or later."
    exit 1
fi

# Check Go version
GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
echo "✓ Go version: $GO_VERSION"

# Install required tools
echo "📦 Installing required tools..."

echo "  - Installing buf..."
go install github.com/bufbuild/buf/cmd/buf@latest

echo "  - Installing protoc-gen-go..."
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest

echo "  - Installing protoc-gen-go-grpc..."
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

echo "  - Installing protoc-gen-grpc-gateway..."
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest

echo "  - Installing protoc-gen-openapiv2..."
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest

echo "  - Installing golangci-lint (optional)..."
if ! command -v golangci-lint &> /dev/null; then
    go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
fi

echo "  - Installing grpcurl (optional)..."
if ! command -v grpcurl &> /dev/null; then
    go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest
fi

echo "✓ Tools installed"

# Download dependencies
echo "📥 Downloading Go dependencies..."
go mod download
go mod tidy
echo "✓ Dependencies downloaded"

# Generate proto files
echo "🔨 Generating code from proto files..."
buf lint
buf generate
echo "✓ Proto files generated"

echo ""
echo "✅ Setup complete!"
echo ""
echo "Next steps:"
echo "  1. Update module name in go.mod and all import paths"
echo "  2. Run 'make build' to build the application"
echo "  3. Run 'make run' to start the server"
echo "  4. Visit http://localhost:8080/swagger/ for API documentation"
echo ""
echo "For more information, see README.md"
