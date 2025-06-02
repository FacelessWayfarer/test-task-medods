package routes

import (
	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"

	_ "github.com/FacelessWayfarer/test-task-medods/docs"
	"github.com/FacelessWayfarer/test-task-medods/internal/database"
	"github.com/FacelessWayfarer/test-task-medods/internal/handlers"
	"github.com/FacelessWayfarer/test-task-medods/internal/service"
)

func New(storage *database.Database) *chi.Mux {
	router := chi.NewRouter()

	services := service.NewService(*storage)

	handler := handlers.NewHandler(*services)

	router.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"),
	))

	router.Get("/tokens/{user_id}", handler.GetTokens)

	router.Post("/tokens/", handler.RefreshTokens)

	return router
}
