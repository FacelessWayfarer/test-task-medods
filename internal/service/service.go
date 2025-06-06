package service

import (
	"context"

	"github.com/FacelessWayfarer/test-task-medods/internal/messenger"
	"github.com/FacelessWayfarer/test-task-medods/internal/service/models"
	"github.com/FacelessWayfarer/test-task-medods/internal/storage"
	"github.com/FacelessWayfarer/test-task-medods/pkg/logger"
	"github.com/FacelessWayfarer/test-task-medods/pkg/tokens"
)

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

func NewService(logger logger.Logger, storage storage.IStorage, tokenCreator tokens.IJWTMaker, messenger messenger.IMessenger) *Service {
	return &Service{
		userStorage:    storage,
		sessionStorage: storage,
		tokenCreator:   tokenCreator,
		emailSender:    messenger,
		logger:         logger,
	}
}
