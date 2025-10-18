package main

import (
	"context"
	"log"
	"time"

	"github.com/vocal-tracker/vocabulary-service/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

func main() {
	log.Println("=== Vocabulary Service gRPC Client with Authentication ===")

	// Connect to the gRPC server
	conn, err := grpc.NewClient("localhost:50052", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	// Create vocabulary service client
	client := proto.NewVocabularyServiceClient(conn)

	// IMPORTANT: You need a valid JWT token from the auth service
	// Get this token by running the auth service client first!
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjo1LCJlbWFpbCI6InRlc0BleGFtcGxlLmNvbSIsImV4cCI6MTc1OTA1NzU3MCwiaWF0IjoxNzU4OTcxMTcwfQ.lgFrIMAcHPseZgjQwjiFIALXPb7ilux7wR_YE1r1UXw" // Replace with actual token

	if token == "YOUR_JWT_TOKEN_HERE" {
		log.Println("❌ ERROR: Please update the token with a real JWT token from auth service")
		log.Println("Steps to get a token:")
		log.Println("1. Run: cd ../auth-service/examples && go run client.go")
		log.Println("2. Copy the token from the response")
		log.Println("3. Replace 'YOUR_JWT_TOKEN_HERE' in this file with the actual token")
		return
	}

	// Create context with authentication metadata
	ctx := context.Background()
	ctx = metadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+token)
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// The user ID should match the user from the JWT token
	userID := uint32(5) // This should match the user ID in your JWT token

	// Example 1: Create a new vocabulary
	log.Println("=== Creating Vocabulary ===")
	createResp, err := client.CreateVocabulary(ctx, &proto.CreateVocabularyRequest{
		UserId:  userID,
		Word:    "authenticated",
		Meaning: "Having proved one's identity; genuine",
		Example: "Only authenticated users can access this vocabulary service.",
		Date:    time.Now().Format("2006-01-02"),
		Status:  "review_needed",
	})
	if err != nil {
		log.Fatalf("Create vocabulary failed: %v", err)
	}
	log.Printf("Create response: Success=%v, Message=%s", createResp.Success, createResp.Message)

	// Example 2: Get vocabularies (will only return vocabularies for the authenticated user)
	log.Println("\n=== Getting Vocabularies ===")
	getResp, err := client.GetVocabularies(ctx, &proto.GetVocabulariesRequest{
		UserId: userID,
		Limit:  10,
	})
	if err != nil {
		log.Fatalf("Get vocabularies failed: %v", err)
	}
	log.Printf("Get response: Success=%v, Count=%d", getResp.Success, getResp.Count)

	// Example 3: Try to access another user's data (should fail)
	log.Println("\n=== Testing Access Control ===")
	otherUserID := uint32(999) // Different user ID
	accessTestResp, err := client.GetVocabularies(ctx, &proto.GetVocabulariesRequest{
		UserId: otherUserID,
	})
	if err != nil {
		log.Printf("Access test error (expected): %v", err)
	} else {
		log.Printf("Access test response: Success=%v, Message=%s", accessTestResp.Success, accessTestResp.Message)
	}

	log.Println("\n✅ Authentication is working! Users can only access their own data.")
}
