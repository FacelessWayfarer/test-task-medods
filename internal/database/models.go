package database

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `db:"id"`
	Email     string    `db:"email"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type Session struct {
	ID            string    `db:"id"`
	UserId        uuid.UUID `db:"user_id"`
	UserIp        string    `db:"user_ip"`
	Refresh_token string    `db:"refresh_token"`
	CreatedAt     time.Time `db:"created_at"`
	ExpiresAt     time.Time `db:"updated_at"`
}
