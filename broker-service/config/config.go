package config

import (
	"fmt"
	"log"
	"os"

	pb "github.com/vocal-tracker/broker-service/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Config struct {
	AuthServiceConn    *grpc.ClientConn
	AuthServiceClient  pb.AuthServiceClient
	VocabServiceConn   *grpc.ClientConn
	VocabServiceClient pb.VocabularyServiceClient
}

func NewConfig() *Config {
	// Get service hostnames from environment variables (default to localhost for development)
	authHost := getEnv("AUTH_SERVICE_HOST", "localhost")
	vocabHost := getEnv("VOCAB_SERVICE_HOST", "localhost")

	// Connect to auth-service (runs on port 50051)
	authAddr := fmt.Sprintf("%s:50051", authHost)
	authConn, err := grpc.NewClient(authAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("Failed to connect to auth service at %s: %v", authAddr, err)
		return nil
	}

	// Connect to vocabulary-service (runs on port 50052)
	vocabAddr := fmt.Sprintf("%s:50052", vocabHost)
	vocabConn, err := grpc.NewClient(vocabAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("Failed to connect to vocabulary service: %v", err)
		authConn.Close()
		return nil
	}

	authClient := pb.NewAuthServiceClient(authConn)
	vocabClient := pb.NewVocabularyServiceClient(vocabConn)

	return &Config{
		AuthServiceConn:    authConn,
		AuthServiceClient:  authClient,
		VocabServiceConn:   vocabConn,
		VocabServiceClient: vocabClient,
	}
}

func (c *Config) Close() {
	if c.AuthServiceConn != nil {
		c.AuthServiceConn.Close()
	}
	if c.VocabServiceConn != nil {
		c.VocabServiceConn.Close()
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
