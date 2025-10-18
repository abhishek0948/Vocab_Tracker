package services

import (
	"context"
	"errors"
	"time"

	"github.com/vocal-tracker/auth-service/database"
	"github.com/vocal-tracker/auth-service/middleware"
	"github.com/vocal-tracker/auth-service/models"
	"github.com/vocal-tracker/auth-service/proto"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthServiceImpl struct {
	proto.UnimplementedAuthServiceServer
}

func NewAuthService() *AuthServiceImpl {
	return &AuthServiceImpl{}
}

// Register implements the Register RPC method
func (s *AuthServiceImpl) Register(ctx context.Context, req *proto.RegisterRequest) (*proto.AuthResponse, error) {
	// Check if user already exists
	var existingUser models.User
	if err := database.DB.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		return &proto.AuthResponse{
			Success: false,
			Message: "User already exists",
		}, nil
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return &proto.AuthResponse{
			Success: false,
			Message: "Failed to hash password",
		}, err
	}

	// Create user
	user := models.User{
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		CreatedAt:    time.Now(),
	}

	if err := database.DB.Create(&user).Error; err != nil {
		return &proto.AuthResponse{
			Success: false,
			Message: "Failed to create user",
		}, err
	}

	// Generate token
	token, err := middleware.GenerateToken(user.ID, user.Email)
	if err != nil {
		return &proto.AuthResponse{
			Success: false,
			Message: "Failed to generate token",
		}, err
	}

	return &proto.AuthResponse{
		Success: true,
		Message: "User registered successfully",
		Token:   token,
		User: &proto.User{
			Id:        uint32(user.ID),
			Email:     user.Email,
			CreatedAt: user.CreatedAt.Format(time.RFC3339),
		},
	}, nil
}

// Login implements the Login RPC method
func (s *AuthServiceImpl) Login(ctx context.Context, req *proto.LoginRequest) (*proto.AuthResponse, error) {
	var user models.User
	if err := database.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &proto.AuthResponse{
				Success: false,
				Message: "Invalid credentials",
			}, nil
		}
		return &proto.AuthResponse{
			Success: false,
			Message: "Database error",
		}, err
	}

	// Check password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return &proto.AuthResponse{
			Success: false,
			Message: "Invalid credentials",
		}, nil
	}

	// Generate token
	token, err := middleware.GenerateToken(user.ID, user.Email)
	if err != nil {
		return &proto.AuthResponse{
			Success: false,
			Message: "Failed to generate token",
		}, err
	}

	return &proto.AuthResponse{
		Success: true,
		Message: "Login successful",
		Token:   token,
		User: &proto.User{
			Id:        uint32(user.ID),
			Email:     user.Email,
			CreatedAt: user.CreatedAt.Format(time.RFC3339),
		},
	}, nil
}

// ValidateToken implements the ValidateToken RPC method
func (s *AuthServiceImpl) ValidateToken(ctx context.Context, req *proto.ValidateTokenRequest) (*proto.ValidateTokenResponse, error) {
	userID, email, err := middleware.ValidateToken(req.Token)
	if err != nil {
		return &proto.ValidateTokenResponse{
			Valid:   false,
			Message: "Invalid token",
		}, nil
	}

	return &proto.ValidateTokenResponse{
		Valid:   true,
		Message: "Token is valid",
		UserId:  uint32(userID),
		Email:   email,
	}, nil
}

// GetProfile implements the GetProfile RPC method
func (s *AuthServiceImpl) GetProfile(ctx context.Context, req *proto.GetProfileRequest) (*proto.UserResponse, error) {
	var user models.User
	if err := database.DB.First(&user, req.UserId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &proto.UserResponse{
				Success: false,
				Message: "User not found",
			}, nil
		}
		return &proto.UserResponse{
			Success: false,
			Message: "Database error",
		}, err
	}

	return &proto.UserResponse{
		Success: true,
		Message: "Profile retrieved successfully",
		User: &proto.User{
			Id:        uint32(user.ID),
			Email:     user.Email,
			CreatedAt: user.CreatedAt.Format(time.RFC3339),
		},
	}, nil
}
