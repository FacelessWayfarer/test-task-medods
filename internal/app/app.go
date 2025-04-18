package app

import (
	"context"
	"fmt"
	"net/http"

	"github.com/FacelessWayfarer/test-task-medods/internal/config"
	database "github.com/FacelessWayfarer/test-task-medods/internal/database/postgres"
	tokengenerator "github.com/FacelessWayfarer/test-task-medods/internal/handlers/token-generator"
	tokenrefresher "github.com/FacelessWayfarer/test-task-medods/internal/handlers/token-refresher"

	"github.com/go-chi/chi/v5"
)

func NewHTTPServer(ctx context.Context, cfg *config.Conifg) *http.Server {
	Storage := database.Init(cfg)

	router := chi.NewRouter()

	router.Get("/{user_id}", tokengenerator.New(ctx, Storage, Storage))

	router.Post("/", tokenrefresher.New(ctx, Storage, Storage))

	return &http.Server{
		Addr:         fmt.Sprintf("%s:%d", cfg.HTTP.IP, cfg.HTTP.Port),
		Handler:      router,
		ReadTimeout:  cfg.HTTP.ReadTimeout,
		WriteTimeout: cfg.HTTP.WriteTimeout,
	}
}
