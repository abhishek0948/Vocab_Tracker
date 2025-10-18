package routes

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/vocal-tracker/broker-service/config"
	pb "github.com/vocal-tracker/broker-service/proto"
)

type AuthHandler struct {
	cfg *config.Config
}

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Token   string `json:"token,omitempty"`
	User    *User  `json:"user,omitempty"`
}

type User struct {
	ID        uint32 `json:"id"`
	Email     string `json:"email"`
	CreatedAt string `json:"created_at"`
}

func NewAuthHandler(cfg *config.Config) *AuthHandler {
	return &AuthHandler{cfg: cfg}
}

func (a *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest

	// Parse request body
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate request
	if req.Email == "" || req.Password == "" {
		http.Error(w, "Email and password are required", http.StatusBadRequest)
		return
	}

	// Create gRPC request
	grpcReq := &pb.RegisterRequest{
		Email:    req.Email,
		Password: req.Password,
	}

	// Set timeout context
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Call auth service
	resp, err := a.cfg.AuthServiceClient.Register(ctx, grpcReq)
	if err != nil {
		log.Printf("Failed to register user: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Convert response
	authResp := &AuthResponse{
		Success: resp.Success,
		Message: resp.Message,
		Token:   resp.Token,
	}

	if resp.User != nil {
		authResp.User = &User{
			ID:        resp.User.Id,
			Email:     resp.User.Email,
			CreatedAt: resp.User.CreatedAt,
		}
	}

	// Send response
	w.Header().Set("Content-Type", "application/json")
	if !resp.Success {
		w.WriteHeader(http.StatusBadRequest)
	}
	json.NewEncoder(w).Encode(authResp)
}

func (a *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest

	// Parse request body
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate request
	if req.Email == "" || req.Password == "" {
		http.Error(w, "Email and password are required", http.StatusBadRequest)
		return
	}

	// Create gRPC request
	grpcReq := &pb.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
	}

	// Set timeout context
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Call auth service
	resp, err := a.cfg.AuthServiceClient.Login(ctx, grpcReq)
	if err != nil {
		log.Printf("Failed to login user: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Convert response
	authResp := &AuthResponse{
		Success: resp.Success,
		Message: resp.Message,
		Token:   resp.Token,
	}

	if resp.User != nil {
		authResp.User = &User{
			ID:        resp.User.Id,
			Email:     resp.User.Email,
			CreatedAt: resp.User.CreatedAt,
		}
	}

	// Send response
	w.Header().Set("Content-Type", "application/json")
	if !resp.Success {
		w.WriteHeader(http.StatusUnauthorized)
	}
	json.NewEncoder(w).Encode(authResp)
}
