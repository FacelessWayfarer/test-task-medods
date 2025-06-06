package storage

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/stdlib"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"

	"github.com/FacelessWayfarer/test-task-medods/internal/migrations"
	"github.com/FacelessWayfarer/test-task-medods/internal/migrations/fixtures"
	"github.com/FacelessWayfarer/test-task-medods/internal/service/models"
)

type IStorage interface {
	SaveSession(ctx context.Context, session models.Session) error
	GetSession(ctx context.Context, sessionID uuid.UUID) (*models.Session, error)
	GetUser(ctx context.Context, userID string) (*models.User, error)
}
type Storage struct {
	db *sql.DB
}

// New sets up a connection to database
func New(db *sql.DB) *Storage {
	return &Storage{db: db}
}

func Migrate(dbURL string) error {
	stdlib.GetDefaultDriver()

	db, err := goose.OpenDBWithDriver("pgx", dbURL)
	if err != nil {
		return err
	}

	goose.SetBaseFS(&migrations.Content)

	err = goose.SetDialect("postgres")
	if err != nil {
		return fmt.Errorf("failed to migrate database: %v", err)
	}

	err = goose.Up(db, ".")
	if err != nil {
		return fmt.Errorf("failed to migrate database: %v", err)
	}

	err = db.Close()
	if err != nil {
		return fmt.Errorf("failed to close database after applying migrations: %v", err)
	}

	return nil
}

func MigrateFixtures(dbURL string) error {
	stdlib.GetDefaultDriver()

	db, err := goose.OpenDBWithDriver("pgx", dbURL)
	if err != nil {
		return err
	}

	goose.SetBaseFS(&fixtures.Content)

	err = goose.SetDialect("postgres")
	if err != nil {
		return fmt.Errorf("failed to migrate fixtures to database: %v", err)
	}

	err = goose.Up(db, ".")
	if err != nil {
		return fmt.Errorf("failed to migrate fixtures to database: %v", err)
	}

	err = db.Close()
	if err != nil {
		return fmt.Errorf("failed to close database after applying migrations: %v", err)
	}

	return nil
}
