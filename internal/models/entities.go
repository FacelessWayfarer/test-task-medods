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
	UserId       uuid.UUID
	UserIp       string
	RefreshToken string
	CreatedAt    time.Time
	ExpiredAt    time.Time
}
