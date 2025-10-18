package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/vocal-tracker/broker-service/config"
	pb "github.com/vocal-tracker/broker-service/proto"
	"google.golang.org/grpc/metadata"
)

type AuthMiddleware struct {
	cfg *config.Config
}

type AuthenticatedUser struct {
	UserID uint32
	Email  string
}

func NewAuthMiddleware(cfg *config.Config) *AuthMiddleware {
	return &AuthMiddleware{cfg: cfg}
}

// RequireAuth middleware validates JWT token and extracts user info
func (am *AuthMiddleware) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers for protected routes too
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		// Handle preflight OPTIONS request
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Extract token from Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header is required", http.StatusUnauthorized)
			return
		}

		// Check if token has Bearer prefix
		if !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Authorization header must start with Bearer", http.StatusUnauthorized)
			return
		}

		// Extract token
		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == "" {
			http.Error(w, "Token is required", http.StatusUnauthorized)
			return
		}

		// Validate token with auth service
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		validateReq := &pb.ValidateTokenRequest{
			Token: token,
		}

		resp, err := am.cfg.AuthServiceClient.ValidateToken(ctx, validateReq)
		if err != nil {
			http.Error(w, "Failed to validate token", http.StatusInternalServerError)
			return
		}

		if !resp.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Add user info to request context
		user := &AuthenticatedUser{
			UserID: resp.UserId,
			Email:  resp.Email,
		}

		// Create new context with user info
		ctx = context.WithValue(r.Context(), "user", user)
		r = r.WithContext(ctx)

		// Call next handler
		next.ServeHTTP(w, r)
	})
}

// GetUserFromContext extracts authenticated user from request context
func GetUserFromContext(ctx context.Context) (*AuthenticatedUser, bool) {
	user, ok := ctx.Value("user").(*AuthenticatedUser)
	return user, ok
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// WriteErrorResponse writes an error response as JSON
func WriteErrorResponse(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(ErrorResponse{
		Success: false,
		Message: message,
	})
}

// CreateAuthenticatedContext creates a gRPC context with the JWT token from the HTTP request
func CreateAuthenticatedContext(ctx context.Context, r *http.Request) context.Context {
	// Extract token from Authorization header
	authHeader := r.Header.Get("Authorization")
	if authHeader != "" {
		// Add authorization header to gRPC metadata
		md := metadata.New(map[string]string{
			"authorization": authHeader,
		})
		return metadata.NewOutgoingContext(ctx, md)
	}
	return ctx
}
