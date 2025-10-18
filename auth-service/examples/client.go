package main

import (
	"context"
	"log"
	"time"

	"github.com/vocal-tracker/auth-service/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// Connect to the gRPC server
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	// Create auth service client
	client := proto.NewAuthServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// Example: Register a new user
	registerResp, err := client.Register(ctx, &proto.RegisterRequest{
		Email:    "tes@example.com",
		Password: "password123",
	})
	if err != nil {
		log.Fatalf("Register failed: %v", err)
	}

	log.Printf("Register response: %v", registerResp)

	// Example: Login
	loginResp, err := client.Login(ctx, &proto.LoginRequest{
		Email:    "tes@example.com",
		Password: "password123",
	})
	if err != nil {
		log.Fatalf("Login failed: %v", err)
	}

	log.Printf("Login response: %v", loginResp)

	// Example: Validate token
	if loginResp.Success && loginResp.Token != "" {
		validateResp, err := client.ValidateToken(ctx, &proto.ValidateTokenRequest{
			Token: loginResp.Token,
		})
		if err != nil {
			log.Fatalf("Validate token failed: %v", err)
		}

		log.Printf("Validate token response: %v", validateResp)
	}
}
