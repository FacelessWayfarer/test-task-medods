package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Session struct {
	ID            string    `json:"id"`
	UserId        uuid.UUID `json:"user_id"`
	UserIp        string    `json:"user_ip"`
	Refresh_token string    `json:"refresh_token"`
	CreatedAt     time.Time `json:"created_at"`
	ExpiresAt     time.Time `json:"updated_at"`
}
