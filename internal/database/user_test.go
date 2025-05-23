package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/FacelessWayfarer/test-task-medods/internal/models"
)

func TestDatabase_GetUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		fmt.Println("expected no error, but got:", err)

		return
	}

	defer db.Close()

	storage := &Database{DB: db}

	type testCase struct {
		name      string
		mockUser  models.User
		mockErr   error
		mockSetup func(tt testCase)
	}
	tests := []testCase{
		{
			name: "Happy path",
			mockUser: models.User{
				ID:        uuid.New(),
				Email:     "email",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			mockErr: nil,
			mockSetup: func(tt testCase) {
				rows := sqlmock.NewRows([]string{"id", "email", "created_at", "updated_at"}).
					AddRow(fmt.Sprint(tt.mockUser.ID), tt.mockUser.Email, tt.mockUser.CreatedAt, tt.mockUser.UpdatedAt)

				mock.ExpectQuery("SELECT id").WithArgs(fmt.Sprint(tt.mockUser.ID)).
					WillReturnRows(rows)

			},
		}, {
			name: "Error user not found",
			mockUser: models.User{
				ID: uuid.New(),
			},
			mockErr: models.ErrUserNotFound,
			mockSetup: func(tt testCase) {
				mock.ExpectQuery("SELECT id").WithArgs(fmt.Sprint(tt.mockUser.ID)).
					WillReturnError(sql.ErrNoRows)

			},
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup(tt)

			var user *models.User

			user, err = storage.GetUser(context.Background(), fmt.Sprint(tt.mockUser.ID))
			if !errors.Is(err, tt.mockErr) {
				t.Errorf("expected error: %v, got: %v", tt.mockErr, err)
			}

			if user != nil {
				assert.Equal(t, tt.mockUser, *user)
			}

			if err = mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}
