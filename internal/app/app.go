package app

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/FacelessWayfarer/test-task-medods/internal/config"
	"github.com/FacelessWayfarer/test-task-medods/internal/database"
	"github.com/FacelessWayfarer/test-task-medods/internal/routes"
)

func Run() error {
	defer func() {
		if panicErr := recover(); panicErr != nil {
			log.Println(context.Background(), panicErr)
		}
	}()

	cfg := config.SetConfig()

	log.Println("Running Application")

	Server := newHTTPServer(cfg)
	if err := Server.ListenAndServe(); err != nil {
		return err
	}

	return nil
}

func newHTTPServer(cfg *config.Conifg) *http.Server {
	storage := database.New(cfg)

	router := routes.New(storage)

	return &http.Server{
		Addr:         fmt.Sprintf("%s:%d", cfg.HTTP.IP, cfg.HTTP.Port),
		Handler:      router,
		ReadTimeout:  cfg.HTTP.ReadTimeout,
		WriteTimeout: cfg.HTTP.WriteTimeout,
	}
}
