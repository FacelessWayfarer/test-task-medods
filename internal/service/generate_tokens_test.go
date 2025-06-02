package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"

	"github.com/FacelessWayfarer/test-task-medods/internal/service/mocks"
	"github.com/FacelessWayfarer/test-task-medods/internal/service/models"
	"github.com/FacelessWayfarer/test-task-medods/pkg/tokens"
)

func TestService_GenerateTokens(t *testing.T) {
	type testCase struct {
		name          string
		mockError     error
		ctx           context.Context
		userID        uuid.UUID
		userIP        string
		mockUser      *models.User
		mockClaims    *tokens.UserClaims
		expectedError error
		mockSetup     func(tt testCase) *Service
	}
	tests := []testCase{
		{
			name:      "Happy_path",
			mockError: nil,
			ctx:       context.Background(),
			userID:    uuid.MustParse("736f6d65-5f72-616e-646f-6d5f75756964"),
			userIP:    "199.166.0.1",
			mockUser: &models.User{
				ID: uuid.MustParse("736f6d65-5f72-616e-646f-6d5f75756964"),
			},
			mockClaims: &tokens.UserClaims{UserID: uuid.MustParse("736f6d65-5f72-616e-646f-6d5f75756964"),
				UserIP:           "199.166.0.1",
				RegisteredClaims: jwt.RegisteredClaims{ExpiresAt: &jwt.NumericDate{Time: time.Now()}}},
			expectedError: nil,
			mockSetup: func(tt testCase) *Service {
				userStorage := mocks.NewUserStorage(t)

				userStorage.On("GetUser", tt.ctx, tt.userID.String()).
					Return(tt.mockUser, tt.mockError)

				sessionStorage := mocks.NewSessionStorage(t)

				sessionStorage.On("SaveSession", tt.ctx, models.Session{
					UserIP:    tt.userIP,
					UserID:    tt.userID,
					ExpiredAt: tt.mockClaims.ExpiresAt.Time,
				}).
					Return(tt.mockError)

				tokenCreator := mocks.NewTokenCreator(t)

				tokenCreator.On("CreateToken", tt.userID, tt.userIP, jwtAccessTokenDuration).
					Return("", tt.mockClaims, tt.mockError)

				tokenCreator.On("CreateToken", tt.userID, tt.userIP, jwtRefreshTokenDuration).
					Return("", tt.mockClaims, tt.mockError)

				tokenCreator.On("CreateRefreshTokenHash", mock.Anything).
					Return([]byte{})

				return &Service{
					userStorage:    userStorage,
					sessionStorage: sessionStorage,
					tokenCreator:   tokenCreator,
					logger:         log.New(os.Stdout, "test:", log.LstdFlags),
				}
			},
		},
		{
			name:          "No_user_in_storage",
			mockError:     models.ErrUserNotFound,
			ctx:           context.Background(),
			userID:        uuid.MustParse("111f6d65-1111-1111-1111-6d5f75756964"),
			mockUser:      &models.User{},
			mockClaims:    &tokens.UserClaims{RegisteredClaims: jwt.RegisteredClaims{ExpiresAt: &jwt.NumericDate{Time: time.Now()}}},
			expectedError: models.ErrUserNotFound,
			mockSetup: func(tt testCase) *Service {
				userStorage := mocks.NewUserStorage(t)

				userStorage.On("GetUser", tt.ctx, tt.userID.String()).
					Return(nil, tt.mockError)

				return &Service{
					userStorage: userStorage,
					logger:      log.New(os.Stdout, "test:", log.LstdFlags),
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := tt.mockSetup(tt)

			resp, err := s.GenerateTokens(tt.ctx, fmt.Sprint(tt.userID), tt.userIP)
			if !errors.Is(err, tt.expectedError) {
				t.Errorf("got %v err, want %v", err, tt.expectedError)
			}

			_ = resp
		})
	}
}
