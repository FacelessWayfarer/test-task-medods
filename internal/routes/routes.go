package routes

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"

	_ "github.com/FacelessWayfarer/test-task-medods/docs" //nolint:revive
	"github.com/FacelessWayfarer/test-task-medods/internal/handlers"
)

type IRoutes interface {
	ServeHTTP(http.ResponseWriter, *http.Request)
}

// New -.
// @title           Test task API
// @version         1.0
// @description     This is a test task
// @host      localhost:8080
// @BasePath  /
func New(handler handlers.IHandler) *chi.Mux {
	router := chi.NewRouter()

	router.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"),
	))

	router.Get("/tokens/{user_id}", handler.GetTokens)

	router.Post("/tokens/", handler.RefreshTokens)

	return router
}
