package app

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/FacelessWayfarer/test-task-medods/internal/config"
	"github.com/FacelessWayfarer/test-task-medods/internal/handlers"
	"github.com/FacelessWayfarer/test-task-medods/internal/messenger"
	"github.com/FacelessWayfarer/test-task-medods/internal/routes"
	"github.com/FacelessWayfarer/test-task-medods/internal/service"
	"github.com/FacelessWayfarer/test-task-medods/internal/storage"
	"github.com/FacelessWayfarer/test-task-medods/pkg/tokens"
)

// JWT is created with secret string that is stored in ENV variable
const secretEnv = "JWT_SECRET"

func Run() error {
	defer func() {
		if panicErr := recover(); panicErr != nil {
			log.Println(context.Background(), panicErr)
		}
	}()

	logger := log.New(os.Stdout, "app:", log.LstdFlags)

	cfg, err := config.SetConfig(logger)
	if err != nil {
		return err
	}

	dbURL := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.PostgreSQL.Username,
		cfg.PostgreSQL.Password,
		cfg.PostgreSQL.Host,
		cfg.PostgreSQL.Port,
		cfg.PostgreSQL.Database,
	)

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return fmt.Errorf("can't connect to database: %v", err)
	}

	if err = storage.Migrate(dbURL); err != nil {
		return err
	}

	if err = storage.MigrateFixtures(dbURL); err != nil {
		return err
	}

	database := storage.New(db)

	tokenCreator := tokens.NewJWT(os.Getenv(secretEnv))

	emailSender := messenger.New()

	services := service.NewService(logger, database, tokenCreator, emailSender)

	handler := handlers.NewHandler(logger, services)

	router := routes.New(handler)

	if err = newHTTPServer(*cfg, router).ListenAndServe(); err != nil {
		return err
	}

	return nil
}

func newHTTPServer(cfg config.Conifg, router routes.IRoutes) *http.Server {
	return &http.Server{
		Addr:         fmt.Sprintf("%s:%d", cfg.HTTP.IP, cfg.HTTP.Port),
		Handler:      router,
		ReadTimeout:  cfg.HTTP.ReadTimeout,
		WriteTimeout: cfg.HTTP.WriteTimeout,
	}
}
