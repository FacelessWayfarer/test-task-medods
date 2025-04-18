package postgresql

import (
	"context"
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/FacelessWayfarer/test-task-medods/internal/config"
	database "github.com/FacelessWayfarer/test-task-medods/internal/database"
	"github.com/FacelessWayfarer/test-task-medods/internal/database/migrations"
	"github.com/jackc/pgx/v5/stdlib"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

type Database struct {
	DB *sql.DB
}

// Init sets up a connection to database and applies migrations
func Init(cfg *config.Conifg) *Database {
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
		log.Fatal("can't connect to database")
	}

	err = migrate(dbURL)
	if err != nil {
		log.Fatal("can't migrate database ", err)
	}

	return &Database{
		DB: db,
	}
}

func (db *Database) GetUser(ctx context.Context, userID string) (database.User, error) {
	const mark = "database.postgres.GetUser"
	var user database.User
	row := db.DB.QueryRowContext(ctx, "SELECT * FROM users WHERE id = $1", userID)
	err := row.Scan(&user.ID, &user.Email, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return database.User{}, database.ErrUserNotFound
		}
		return database.User{}, fmt.Errorf("%s:%w", mark, err)
	}

	return user, nil
}

func (db *Database) SaveSession(ctx context.Context, session *database.Session) error {
	const mark = "database.postgres.SaveToken"
	fmt.Println("database.postgres.SaveToken", session.UserIp)
	query := "INSERT INTO sessions (id, user_id, user_ip, refresh_token, created_at, expires_at) VALUES ($1,$2,$3,$4,$5,$6)"
	hashedToken := base64.StdEncoding.EncodeToString([]byte(session.Refresh_token))
	_, err := db.DB.ExecContext(ctx, query, session.ID, session.UserId, session.UserIp, hashedToken, time.Now().UTC(), session.ExpiresAt.UTC())
	if err != nil {
		return fmt.Errorf("%s:%w", mark, err)
	}
	return nil
}

func (db *Database) GetSession(ctx context.Context, sessionID string) (database.Session, error) {
	const mark = "database.postgres.GetUser"
	var session database.Session
	row := db.DB.QueryRowContext(ctx, "SELECT * FROM sessions WHERE id = $1", sessionID)
	err := row.Scan(&session.ID, &session.UserId, &session.UserIp, &session.Refresh_token, &session.CreatedAt, &session.ExpiresAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return database.Session{}, database.ErrSessionNotFound
		}
		return database.Session{}, fmt.Errorf("%s:%w", mark, err)
	}

	return session, nil
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
