package handlers

import (
	"context"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/mock"

	"github.com/FacelessWayfarer/test-task-medods/internal/handlers/mocks"
	"github.com/FacelessWayfarer/test-task-medods/internal/models"
	"github.com/FacelessWayfarer/test-task-medods/pkg/tokens"
)

func TestHandler_GenerateTokens(t *testing.T) {
	type testCase struct {
		name               string
		mockError          error
		userID             string
		mockUser           *models.User
		mockClaims         *tokens.UserClaims
		expectedStatusCode int
		mockSetup          func(tt testCase) *Handler
	}
	tests := []testCase{
		{
			name:               "Happy path",
			mockError:          nil,
			userID:             "some_user_id",
			mockUser:           &models.User{},
			mockClaims:         &tokens.UserClaims{RegisteredClaims: jwt.RegisteredClaims{ExpiresAt: &jwt.NumericDate{Time: time.Now()}}},
			expectedStatusCode: http.StatusOK,
			mockSetup: func(tt testCase) *Handler {
				userStorage := mocks.NewUserStorage(t)

				userStorage.On("GetUser", mock.Anything, tt.userID).
					Return(tt.mockUser, tt.mockError)

				sessionStorage := mocks.NewSessionStorage(t)

				sessionStorage.On("SaveSession", mock.Anything, mock.Anything).
					Return(tt.mockError)

				tokenCreator := mocks.NewTokenCreator(t)

				tokenCreator.On("CreateToken", mock.Anything, mock.Anything, mock.Anything).
					Return("", tt.mockClaims, tt.mockError)

				tokenCreator.On("CreateRefreshTokenHash", mock.Anything).
					Return([]byte{})

				return &Handler{
					userStorage:    userStorage,
					sessionStorage: sessionStorage,
					tokenCreator:   tokenCreator,
					logger:         log.New(os.Stdout, "test:", log.LstdFlags),
				}
			},
		},
		{
			name:               "No user in storage",
			mockError:          models.ErrUserNotFound,
			userID:             "wrong_user_id",
			mockUser:           &models.User{},
			mockClaims:         &tokens.UserClaims{RegisteredClaims: jwt.RegisteredClaims{ExpiresAt: &jwt.NumericDate{Time: time.Now()}}},
			expectedStatusCode: http.StatusBadRequest,
			mockSetup: func(tt testCase) *Handler {
				userStorage := mocks.NewUserStorage(t)

				userStorage.On("GetUser", mock.Anything, tt.userID).
					Return(nil, tt.mockError)

				return &Handler{
					userStorage: userStorage,
					logger:      log.New(os.Stdout, "test:", log.LstdFlags),
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := tt.mockSetup(tt)

			r := httptest.NewRequest("GET", "/tokens/{user_id}", nil)

			r = addChiURLParam(r, "user_id", tt.userID)

			w := httptest.NewRecorder()

			h.GenerateTokens(w, r)

			resp := w.Result()
			if resp.StatusCode != tt.expectedStatusCode {
				t.Errorf("got %v status code, want %v", resp.StatusCode, tt.expectedStatusCode)
			}

		})
	}
}

func addChiURLParam(r *http.Request, key, value string) *http.Request {
	ctx := chi.NewRouteContext()

	ctx.URLParams.Add(key, value)

	return r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, ctx))
}
