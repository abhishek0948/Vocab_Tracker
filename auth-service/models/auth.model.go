package models

import (
	"time"
)

type User struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	Email        string    `json:"email" gorm:"unique;not null"`
	PasswordHash string    `json:"-" gorm:"column:password_hash;not null"`
	CreatedAt    time.Time `json:"created_at"`
}

type AuthRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}