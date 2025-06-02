package service

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/FacelessWayfarer/test-task-medods/internal/service/models"
	"github.com/FacelessWayfarer/test-task-medods/pkg/tokens"
)

//go:generate go run github.com/vektra/mockery/v2@latest --name=UserStorage
type UserStorage interface {
	GetUser(ctx context.Context, userID string) (*models.User, error)
}

//go:generate go run github.com/vektra/mockery/v2@latest --name=SessionStorage
type SessionStorage interface {
	GetSession(ctx context.Context, sessionID uuid.UUID) (*models.Session, error)
	SaveSession(ctx context.Context, session models.Session) error
}

//go:generate go run github.com/vektra/mockery/v2@latest --name=TokenCreator
type TokenCreator interface {
	CreateToken(userid uuid.UUID, ip string, duration time.Duration) (string, *tokens.UserClaims, error)
	VerifyToken(tokenString string) (*tokens.UserClaims, error)
	CreateRefreshTokenHash(refreshToken string) []byte
}

//go:generate go run github.com/vektra/mockery/v2@latest --name=Logger
type Logger interface {
	Printf(format string, v ...any)
	Println(v ...any)
}

// EmailSender is a mock interface to send email
//
//go:generate go run github.com/vektra/mockery/v2@latest --name=EmailSender
type EmailSender interface {
	SendEmail(string) error
}
