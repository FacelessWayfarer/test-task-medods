package service

import (
	"context"
	"log"
	"os"

	"github.com/FacelessWayfarer/test-task-medods/internal/database"
	"github.com/FacelessWayfarer/test-task-medods/internal/messenger"
	"github.com/FacelessWayfarer/test-task-medods/internal/service/models"
	"github.com/FacelessWayfarer/test-task-medods/pkg/tokens"
)

// JWT is created with secret string that is stored in ENV variable
const secretEnv = "JWT_SECRET"

//go:generate go run github.com/vektra/mockery/v2@latest --name=IService
type IService interface {
	GenerateTokens(ctx context.Context, userID, ip string) (*models.GeneratedTokens, error)
	UpdateTokens(ctx context.Context, req models.TokensToRefresh, ip string) (*models.RefreshedTokens, error)
}

type Service struct {
	userStorage    UserStorage
	sessionStorage SessionStorage
	tokenCreator   TokenCreator
	emailSender    EmailSender
	logger         Logger
}

func NewService(storage database.Database) *Service {
	secretString := os.Getenv(secretEnv)

	tokenCreator := tokens.NewJWT(secretString)

	return &Service{
		userStorage:    &storage,
		sessionStorage: &storage,
		tokenCreator:   tokenCreator,
		emailSender:    &messenger.Messenger{},
		logger:         log.New(os.Stdout, "test:", log.LstdFlags),
	}
}
