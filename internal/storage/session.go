package storage

import (
	"context"
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"

	"github.com/FacelessWayfarer/test-task-medods/internal/service/models"
)

func (s *Storage) SaveSession(ctx context.Context, session models.Session) error {
	const mark = "database.SaveSession"

	query := `INSERT INTO sessions (id, user_id, user_ip, refresh_token, created_at, expired_at) VALUES ($1,$2,$3,$4,$5,$6);`

	hashedToken := base64.StdEncoding.EncodeToString([]byte(session.RefreshToken))

	_, err := s.db.ExecContext(ctx, query, session.ID, session.UserID, session.UserIP, hashedToken, time.Now(), session.ExpiredAt)

	if err != nil {
		return fmt.Errorf("%s:%w", mark, err)
	}

	return nil
}

func (s *Storage) GetSession(ctx context.Context, sessionID uuid.UUID) (*models.Session, error) {
	const mark = "database.GetSession"

	var session models.Session

	row := s.db.QueryRowContext(ctx, `SELECT id, user_id, user_ip, refresh_token, created_at, expired_at FROM sessions WHERE id = $1;`, sessionID)

	if err := row.Scan(&session.ID, &session.UserID, &session.UserIP, &session.RefreshToken, &session.CreatedAt, &session.ExpiredAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrSessionNotFound
		}

		return nil, fmt.Errorf("%s:%w", mark, err)
	}

	return &session, nil
}
