package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/vocal-tracker/vocabulary-service/config"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type Claims struct {
	UserID uint   `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

func GenerateToken(userID uint, email string) (string, error) {
	cfg := config.GetConfig()

	claims := &Claims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(cfg.JWTSecret))
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Bearer token required"})
			c.Abort()
			return
		}

		cfg := config.GetConfig()
		token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(cfg.JWTSecret), nil
		})

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		if claims, ok := token.Claims.(*Claims); ok && token.Valid {
			c.Set("userID", claims.UserID)
			c.Set("email", claims.Email)
			c.Next()
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
		}
	}
}

// ValidateToken validates a JWT token and returns user info (for auth service compatibility)
func ValidateToken(tokenString string) (uint, string, error) {
	cfg := config.GetConfig()

	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(cfg.JWTSecret), nil
	})

	if err != nil {
		return 0, "", err
	}

	if !token.Valid {
		return 0, "", errors.New("token is not valid")
	}

	return claims.UserID, claims.Email, nil
}

// gRPC Authentication Interceptor

// AuthInterceptor is a gRPC interceptor that validates JWT tokens
type AuthInterceptor struct {
	jwtSecret string
}

// NewAuthInterceptor creates a new auth interceptor
func NewAuthInterceptor() *AuthInterceptor {
	cfg := config.GetConfig()
	return &AuthInterceptor{
		jwtSecret: cfg.JWTSecret,
	}
}

// UnaryInterceptor validates JWT tokens for unary RPC calls
func (a *AuthInterceptor) UnaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	// Skip authentication for health checks or other public methods
	if isPublicMethod(info.FullMethod) {
		return handler(ctx, req)
	}

	// Extract token from metadata
	token, err := extractTokenFromMetadata(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "missing or invalid token: %v", err)
	}

	// Validate token locally
	userID, email, err := ValidateToken(token)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid token: %v", err)
	}

	// Add user info to context
	ctx = context.WithValue(ctx, "userID", uint32(userID))
	ctx = context.WithValue(ctx, "email", email)

	return handler(ctx, req)
}

// extractTokenFromMetadata extracts JWT token from gRPC metadata
func extractTokenFromMetadata(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", errors.New("missing metadata")
	}

	authHeader := md.Get("authorization")
	if len(authHeader) == 0 {
		return "", errors.New("missing authorization header")
	}

	// Expected format: "Bearer <token>"
	authValue := authHeader[0]
	if !strings.HasPrefix(authValue, "Bearer ") {
		return "", errors.New("invalid authorization header format")
	}

	token := strings.TrimPrefix(authValue, "Bearer ")
	if token == "" {
		return "", errors.New("empty token")
	}

	return token, nil
}

// isPublicMethod checks if a method should skip authentication
func isPublicMethod(fullMethod string) bool {
	publicMethods := []string{
		// Add any public methods here (like health checks)
		// "/grpc.health.v1.Health/Check",
	}

	for _, method := range publicMethods {
		if fullMethod == method {
			return true
		}
	}
	return false
}

// GetUserIDFromContext extracts user ID from context (for use in service methods)
func GetUserIDFromContext(ctx context.Context) (uint32, error) {
	userID, ok := ctx.Value("userID").(uint32)
	if !ok {
		return 0, errors.New("user ID not found in context")
	}
	return userID, nil
}

// GetEmailFromContext extracts email from context (for use in service methods)
func GetEmailFromContext(ctx context.Context) (string, error) {
	email, ok := ctx.Value("email").(string)
	if !ok {
		return "", errors.New("email not found in context")
	}
	return email, nil
}
