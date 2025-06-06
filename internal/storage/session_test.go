package storage //nolint:testpackage

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"encoding/base64"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/FacelessWayfarer/test-task-medods/internal/service/models"
)

// Define a custom type AnyTime to match any value to pass as args in Mock methods
type AnyTime struct{}

func (a AnyTime) Match(_ driver.Value) bool {
	return true
}

func TestDatabase_SaveSession(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		fmt.Println("expected no error, but got:", err)

		return
	}

	defer db.Close()

	storage := &Storage{db: db}

	type testCase struct {
		name        string
		mockSession models.Session
		mockErr     error
		mockSetup   func(tt testCase)
	}
	tests := []testCase{
		{
			name: "Happy_path",
			mockSession: models.Session{
				ID:           uuid.New(),
				UserID:       uuid.New(),
				UserIP:       "ip",
				RefreshToken: "token",
				CreatedAt:    time.Now(),
				ExpiredAt:    time.Now().Add(time.Hour),
			},
			mockErr: nil,
			mockSetup: func(tt testCase) {
				mock.ExpectExec("INSERT INTO sessions").
					WithArgs(tt.mockSession.ID,
						tt.mockSession.UserID,
						tt.mockSession.UserIP,
						base64.StdEncoding.EncodeToString([]byte(tt.mockSession.RefreshToken)),
						AnyTime{},
						AnyTime{}).
					WillReturnResult(sqlmock.NewResult(1, 1)).
					WillReturnError(tt.mockErr)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup(tt)

			if err = storage.SaveSession(context.Background(), tt.mockSession); !errors.Is(err, tt.mockErr) {
				t.Errorf("expected error: %v", err)
			}

			if err = mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestDatabase_GetSession(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		fmt.Println("expected no error, but got:", err)

		return
	}

	defer db.Close()

	storage := &Storage{db: db}

	type testCase struct {
		name        string
		mockSession models.Session
		mockErr     error
		mockSetup   func(tt testCase)
	}
	tests := []testCase{
		{
			name: "Happy_path",
			mockSession: models.Session{
				ID:           uuid.New(),
				UserID:       uuid.New(),
				UserIP:       "ip",
				RefreshToken: "token",
				CreatedAt:    time.Now(),
				ExpiredAt:    time.Now().Add(time.Hour),
			},
			mockErr: nil,
			mockSetup: func(tt testCase) {
				rows := sqlmock.NewRows([]string{"id", "user_id", "user_ip", "refresh_token", "created_at", "expired_at"}).
					AddRow(tt.mockSession.ID,
						tt.mockSession.UserID,
						tt.mockSession.UserIP,
						tt.mockSession.RefreshToken,
						tt.mockSession.CreatedAt,
						tt.mockSession.ExpiredAt,
					)

				mock.ExpectQuery("SELECT id").WithArgs(fmt.Sprint(tt.mockSession.ID)).
					WillReturnRows(rows)
			},
		},
		{
			name: "Error_session not found",
			mockSession: models.Session{
				ID: uuid.New(),
			},
			mockErr: models.ErrSessionNotFound,
			mockSetup: func(tt testCase) {
				mock.ExpectQuery("SELECT id").WithArgs(fmt.Sprint(tt.mockSession.ID)).
					WillReturnError(sql.ErrNoRows)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup(tt)

			var session *models.Session

			session, err = storage.GetSession(context.Background(), tt.mockSession.ID)
			if !errors.Is(err, tt.mockErr) {
				t.Errorf("expected error: %v, got: %v", tt.mockErr, err)
			}

			if session != nil {
				assert.Equal(t, tt.mockSession, *session)
			}

			if err = mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}
