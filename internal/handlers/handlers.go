package handlers

import (
	"context"
	"github.com/FacelessWayfarer/test-task-medods/internal/database"
	"log"
	"os"
	"time"

	"github.com/google/uuid"

	"github.com/FacelessWayfarer/test-task-medods/internal/models"
	"github.com/FacelessWayfarer/test-task-medods/pkg/response"
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

type Handler struct {
	userStorage    UserStorage
	sessionStorage SessionStorage
	tokenCreator   TokenCreator
	logger         Logger
	emailSender    EmailSender
}

func NewHandler(db *database.Database, tokenCreator TokenCreator) *Handler {
	return &Handler{
		userStorage:    db,
		sessionStorage: db,
		logger:         log.New(os.Stdout, "client:", log.LstdFlags),
		tokenCreator:   tokenCreator,
		emailSender:    Messenger{},
	}
}

type Messenger struct {
}

func (m Messenger) SendEmail(msg string) error {
	_ = msg
	return nil
}

type GenResponse struct {
	response.Response
	AccessToken           string    `json:"AccessToken,omitempty"`
	RefreshToken          string    `json:"RefreshToken,omitempty"`
	AccessTokenExpiresAt  time.Time `json:"AccessTokenExpiresAt,omitempty"`
	RefreshTokenExpiresAt time.Time `json:"RefreshTokenExpiresAt,omitempty"`
}

type RefreshTokensReq struct {
	AccessToken        string `json:"AccessToken"`
	Base64RefreshToken string `json:"Encoded_refresh_token"`
}

type RefreshResponse struct {
	response.Response
	AccessToken           string    `json:"AccessToken,omitempty"`
	RefreshToken          string    `json:"RefreshToken,omitempty"`
	AccessTokenExpiresAt  time.Time `json:"AccessTokenExpiresAt,omitempty"`
	RefreshTokenExpiresAt time.Time `json:"RefreshTokenExpiresAt,omitempty"`
}
