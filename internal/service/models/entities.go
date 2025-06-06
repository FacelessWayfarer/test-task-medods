package models

import (
	"time"

	"github.com/google/uuid"
)

// User is a struct representing "users" table in the database
type User struct {
	ID        uuid.UUID
	Email     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Session is a struct representing "sessions" table in the database
type Session struct {
	ID           uuid.UUID
	UserID       uuid.UUID
	UserIP       string
	RefreshToken string
	CreatedAt    time.Time
	ExpiredAt    time.Time
}

// GeneratedTokens is a response for a method GenerateTokens from service
type GeneratedTokens struct {
	AccessToken           string
	RefreshToken          string
	AccessTokenExpiresAt  time.Time
	RefreshTokenExpiresAt time.Time
}

type TokensToRefresh struct {
	AccessToken        string
	Base64RefreshToken string
}

// RefreshedTokens is a response for a method UpdateTokens from service
type RefreshedTokens struct {
	AccessToken           string
	RefreshToken          string
	AccessTokenExpiresAt  time.Time
	RefreshTokenExpiresAt time.Time
}
