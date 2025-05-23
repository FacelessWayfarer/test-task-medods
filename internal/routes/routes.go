package routes

import (
	"os"

	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"

	_ "github.com/FacelessWayfarer/test-task-medods/docs"
	"github.com/FacelessWayfarer/test-task-medods/internal/database"
	"github.com/FacelessWayfarer/test-task-medods/internal/handlers"
	"github.com/FacelessWayfarer/test-task-medods/pkg/tokens"
)

// JWT is created with secret string that is stored in ENV variable
const secretEnv = "JWT_SECRET"

func New(storage *database.Database) *chi.Mux {
	router := chi.NewRouter()

	secretString := os.Getenv(secretEnv)

	tokenCreator := tokens.NewJWT(secretString)

	handler := handlers.NewHandler(storage, tokenCreator)

	router.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"), //The url pointing to API definition
	))

	router.Get("/tokens/{user_id}", handler.GenerateTokens)

	router.Post("/tokens/", handler.RefreshTokens)

	return router
}
