package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	v1 "github.com/ChyiYaqing/go-microservice-template/api/v1"
	"github.com/ChyiYaqing/go-microservice-template/internal/service"
	"github.com/ChyiYaqing/go-microservice-template/pkg/config"
	"github.com/ChyiYaqing/go-microservice-template/pkg/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

func main() {
	// Initialize logger
	log := logger.NewLogger()

	// Load configuration
	cfg := config.Default()
	if len(os.Args) > 1 {
		loadedCfg, err := config.Load(os.Args[1])
		if err != nil {
			log.Warn("Failed to load config file, using defaults: %v", err)
		} else {
			cfg = loadedCfg
		}
	}

	// Create context that listens for the interrupt signal
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Start gRPC server
	grpcServer := startGRPCServer(cfg, log)

	// Start HTTP server with grpc-gateway
	httpServer := startHTTPServer(ctx, cfg, log)

	log.Info("Server started successfully")
	log.Info("gRPC server listening on %s:%d", cfg.Server.Host, cfg.Server.GRPCPort)
	log.Info("HTTP server listening on %s:%d", cfg.Server.Host, cfg.Server.HTTPPort)
	log.Info("Swagger UI available at http://%s:%d/swagger/", cfg.Server.Host, cfg.Server.HTTPPort)

	// Wait for interrupt signal
	<-ctx.Done()
	log.Info("Shutting down servers...")

	// Graceful shutdown
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		log.Error("HTTP server shutdown error: %v", err)
	}

	grpcServer.GracefulStop()
	log.Info("Servers stopped")
}

func startGRPCServer(cfg *config.Config, log logger.Logger) *grpc.Server {
	// Create gRPC server
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(loggingInterceptor(log)),
	)

	// Register services
	userService := service.NewUserService()
	v1.RegisterUserServiceServer(grpcServer, userService)

	// Register reflection service for grpcurl
	reflection.Register(grpcServer)

	// Start listening
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.GRPCPort))
	if err != nil {
		log.Error("Failed to listen: %v", err)
		os.Exit(1)
	}

	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Error("Failed to serve gRPC: %v", err)
			os.Exit(1)
		}
	}()

	return grpcServer
}

func startHTTPServer(ctx context.Context, cfg *config.Config, log logger.Logger) *http.Server {
	// Create gRPC client connection
	conn, err := grpc.NewClient(
		fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.GRPCPort),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Error("Failed to create gRPC client: %v", err)
		os.Exit(1)
	}

	// Create gRPC-Gateway mux
	mux := runtime.NewServeMux(
		runtime.WithErrorHandler(customErrorHandler),
	)

	// Register service handlers
	if err := v1.RegisterUserServiceHandler(ctx, mux, conn); err != nil {
		log.Error("Failed to register gateway: %v", err)
		os.Exit(1)
	}

	// Create HTTP mux for additional routes
	httpMux := http.NewServeMux()

	// API routes
	httpMux.Handle("/", mux)

	// Swagger UI
	httpMux.HandleFunc("/swagger/", serveSwagger)
	httpMux.HandleFunc("/swagger/api.swagger.json", serveSwaggerJSON)

	// Health check
	httpMux.HandleFunc("/health", healthCheckHandler)

	// Create HTTP server
	httpServer := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.HTTPPort),
		Handler: corsMiddleware(loggingMiddleware(log, httpMux)),
	}

	go func() {
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("Failed to serve HTTP: %v", err)
			os.Exit(1)
		}
	}()

	return httpServer
}

// loggingInterceptor logs gRPC requests
func loggingInterceptor(log logger.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		start := time.Now()
		resp, err := handler(ctx, req)
		duration := time.Since(start)

		if err != nil {
			log.Error("gRPC %s failed: %v (duration: %v)", info.FullMethod, err, duration)
		} else {
			log.Info("gRPC %s succeeded (duration: %v)", info.FullMethod, duration)
		}

		return resp, err
	}
}

// loggingMiddleware logs HTTP requests
func loggingMiddleware(log logger.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		duration := time.Since(start)
		log.Info("HTTP %s %s (duration: %v)", r.Method, r.URL.Path, duration)
	})
}

// corsMiddleware adds CORS headers
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// customErrorHandler handles errors from gRPC-Gateway
func customErrorHandler(ctx context.Context, mux *runtime.ServeMux, marshaler runtime.Marshaler, w http.ResponseWriter, r *http.Request, err error) {
	runtime.DefaultHTTPErrorHandler(ctx, mux, marshaler, w, r, err)
}

// healthCheckHandler handles health check requests
func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok"}`))
}

// serveSwagger serves the Swagger UI
func serveSwagger(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "docs/swagger/index.html")
}

// serveSwaggerJSON serves the Swagger JSON
func serveSwaggerJSON(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "docs/swagger/api.swagger.json")
}
