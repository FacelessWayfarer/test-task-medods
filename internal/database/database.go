package database

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/stdlib"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"

	"github.com/FacelessWayfarer/test-task-medods/internal/config"
	"github.com/FacelessWayfarer/test-task-medods/internal/database/fixtures"
	"github.com/FacelessWayfarer/test-task-medods/internal/database/migrations"
)

type Database struct {
	DB *sql.DB
}

// New sets up a connection to models and applies migrations
func New(cfg *config.Conifg) *Database {
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
		log.Fatal("can't connect to models")
	}

	if err = migrate(dbURL); err != nil {
		log.Fatal("can't migrate models", err)
	}

	if err = migrateFixtures(dbURL); err != nil {
		log.Fatal("can't migrate data", err)
	}

	return &Database{
		DB: db,
	}
}

func migrate(dbURL string) error {
	stdlib.GetDefaultDriver()

	db, err := goose.OpenDBWithDriver("pgx", dbURL)
	if err != nil {
		return err
	}

	goose.SetBaseFS(&migrations.Content)

	err = goose.SetDialect("postgres")
	if err != nil {
		return err
	}

	err = goose.Up(db, ".")
	if err != nil {
		return err
	}

	err = db.Close()
	if err != nil {
		return err
	}

	return nil
}

func migrateFixtures(dbURL string) error {
	stdlib.GetDefaultDriver()

	db, err := goose.OpenDBWithDriver("pgx", dbURL)
	if err != nil {
		return err
	}

	goose.SetBaseFS(&fixtures.Content)

	err = goose.SetDialect("postgres")
	if err != nil {
		return err
	}

	err = goose.Up(db, ".")
	if err != nil {
		return err
	}

	err = db.Close()
	if err != nil {
		return err
	}

	return nil
}
