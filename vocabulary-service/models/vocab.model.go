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

type Vocabulary struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserID    uint      `json:"user_id" gorm:"not null"`
	Word      string    `json:"word" gorm:"not null"`
	Meaning   string    `json:"meaning" gorm:"not null"`
	Example   string    `json:"example"`
	Date      time.Time `json:"date" gorm:"type:date;not null"`
	Status    string    `json:"status" gorm:"default:'review_needed'"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	User      User      `json:"user" gorm:"foreignKey:UserID"`
}

type VocabRequest struct {
	Word    string `json:"word" binding:"required"`
	Meaning string `json:"meaning" binding:"required"`
	Example string `json:"example"`
	Date    string `json:"date" binding:"required"`
	Status  string `json:"status"`
}

type VocabUpdateRequest struct {
	Word    string `json:"word"`
	Meaning string `json:"meaning"`
	Example string `json:"example"`
	Status  string `json:"status"`
}