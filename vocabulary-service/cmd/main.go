package main

import (
	"fmt"
	"log"
	"net"

	"github.com/vocal-tracker/vocabulary-service/config"
	"github.com/vocal-tracker/vocabulary-service/database"
	"github.com/vocal-tracker/vocabulary-service/middleware"
	"github.com/vocal-tracker/vocabulary-service/proto"
	"github.com/vocal-tracker/vocabulary-service/services"

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

	// Create auth interceptor
	authInterceptor := middleware.NewAuthInterceptor()

	// Create gRPC server with authentication interceptor
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(authInterceptor.UnaryInterceptor),
	)

	// Register vocabulary service
	vocabService := services.NewVocabularyService()
	proto.RegisterVocabularyServiceServer(grpcServer, vocabService)

	// Start listening on port 50052 (different from auth service)
	port := ":50052"
	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal("Failed to listen:", err)
	}

	fmt.Printf("Vocabulary service listening on port %s\n", port)

	// Start the gRPC server
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatal("Failed to serve:", err)
	}
}
