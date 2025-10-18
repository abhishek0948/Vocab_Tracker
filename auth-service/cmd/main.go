package main

import (
	"fmt"
	"log"
	"net"

	"github.com/vocal-tracker/auth-service/config"
	"github.com/vocal-tracker/auth-service/database"
	"github.com/vocal-tracker/auth-service/proto"
	"github.com/vocal-tracker/auth-service/services"

	"google.golang.org/grpc"
)

func main() {
	// Load configuration
	cfg := config.GetConfig()

	// Initialize database
	if err := database.InitDB(cfg); err != nil {
		log.Fatal("Failed to initialize database:", err)
	}

	// Run database migrations
	if err := database.Migrate(); err != nil {
		log.Fatal("Failed to run migrations:", err)
	}

	// Create gRPC server
	grpcServer := grpc.NewServer()

	// Register auth service
	authService := services.NewAuthService()
	proto.RegisterAuthServiceServer(grpcServer, authService)

	// Start listening on port 50051
	port := ":50051"
	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal("Failed to listen:", err)
	}

	fmt.Printf("Auth service listening on port %s\n", port)

	// Start the gRPC server
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatal("Failed to serve:", err)
	}
}
