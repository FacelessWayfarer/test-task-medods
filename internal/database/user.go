package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/lib/pq"

	"github.com/FacelessWayfarer/test-task-medods/internal/service/models"
)

func (db *Database) GetUser(ctx context.Context, userID string) (*models.User, error) {
	const mark = "database.GetUser"

	var user models.User

	row := db.DB.QueryRowContext(ctx, `SELECT id, email, created_at, updated_at FROM users WHERE id = $1;`, userID)

	if err := row.Scan(&user.ID, &user.Email, &user.CreatedAt, &user.UpdatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrUserNotFound
		}

		return nil, fmt.Errorf("%s:%w", mark, err)
	}

	return &user, nil
}
