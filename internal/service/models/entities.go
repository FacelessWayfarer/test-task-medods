package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID
	Email     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Session struct {
	ID           uuid.UUID
	UserID       uuid.UUID
	UserIP       string
	RefreshToken string
	CreatedAt    time.Time
	ExpiredAt    time.Time
}

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

type RefreshedTokens struct {
	AccessToken           string
	RefreshToken          string
	AccessTokenExpiresAt  time.Time
	RefreshTokenExpiresAt time.Time
}
