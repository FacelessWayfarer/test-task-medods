package handlers

import (
	"context"
	"net/http"

	"github.com/FacelessWayfarer/test-task-medods/internal/service"
	"github.com/FacelessWayfarer/test-task-medods/internal/service/models"
	"github.com/FacelessWayfarer/test-task-medods/pkg/logger"
)

type Logger interface {
	Printf(format string, v ...any)
	Println(v ...any)
}

type Service interface {
	GenerateTokens(ctx context.Context, userID, ip string) (*models.GeneratedTokens, error)
	UpdateTokens(ctx context.Context, req models.TokensToRefresh, ip string) (*models.RefreshedTokens, error)
}

type IHandler interface {
	GetTokens(w http.ResponseWriter, r *http.Request)
	RefreshTokens(w http.ResponseWriter, r *http.Request)
}

type Handler struct {
	service Service
	logger  Logger
}

func NewHandler(logger logger.Logger, service service.IService) *Handler {
	return &Handler{
		service: service,
		logger:  logger,
	}
}
