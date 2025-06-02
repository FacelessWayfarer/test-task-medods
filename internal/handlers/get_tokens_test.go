package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/mock"

	"github.com/FacelessWayfarer/test-task-medods/internal/service/mocks"
	"github.com/FacelessWayfarer/test-task-medods/internal/service/models"
)

// jwtTestRefreshTokenDuration and jwtTestAccessTokenDuration are used here only for testing
const (
	jwtTestRefreshTokenDuration time.Duration = 24 * time.Hour
	jwtTestAccessTokenDuration  time.Duration = 1 * time.Minute
)

func TestHandler_GetTokens(t *testing.T) {
	type testCase struct {
		name               string
		IP                 string
		userID             string
		output             *models.GeneratedTokens
		mockError          error
		expectedError      error
		expectedStatusCode int
		mockSetup          func(tt testCase) *Handler
	}
	tests := []testCase{
		{
			name:   "Happy_path",
			IP:     "8.8.8.8",
			userID: "some_user_id",
			output: &models.GeneratedTokens{
				AccessToken:           "access",
				RefreshToken:          "refresh",
				AccessTokenExpiresAt:  time.Now().Add(jwtTestAccessTokenDuration),
				RefreshTokenExpiresAt: time.Now().Add(jwtTestRefreshTokenDuration),
			},
			mockError:          nil,
			expectedError:      nil,
			expectedStatusCode: http.StatusOK,
			mockSetup: func(tt testCase) *Handler {
				service := mocks.NewIService(t)

				// It's impossible mock request context as it uses chi URL parameters
				service.On("GenerateTokens", mock.Anything, tt.userID, tt.IP).
					Return(tt.output, tt.mockError)

				return &Handler{
					service: service,
					logger:  log.New(os.Stdout, "test:", log.LstdFlags),
				}
			},
		},
		{
			name:               "No_user_in_storage",
			IP:                 "9.9.9.9",
			userID:             "wrong_user_id",
			output:             nil,
			mockError:          models.ErrUserNotFound,
			expectedError:      models.ErrUserNotFound,
			expectedStatusCode: http.StatusInternalServerError,
			mockSetup: func(tt testCase) *Handler {
				service := mocks.NewIService(t)

				// It's impossible to mock request context as it uses chi URL parameters
				service.On("GenerateTokens", mock.Anything, tt.userID, tt.IP).
					Return(tt.output, tt.mockError)

				return &Handler{
					service: service,
					logger:  log.New(os.Stdout, "test:", log.LstdFlags),
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := tt.mockSetup(tt)

			r := httptest.NewRequest("GET", "/tokens/{user_id}", nil)

			r = addChiURLParam(r, "user_id", tt.userID)

			r.RemoteAddr = tt.IP

			w := httptest.NewRecorder()

			h.GetTokens(w, r)

			result := w.Result()
			if result.StatusCode != tt.expectedStatusCode {
				t.Errorf("got %v status code, want %v", result.StatusCode, tt.expectedStatusCode)
			}

			var resp GenResponse

			if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
				t.Errorf("failed to decode response body: %v", err)
			}

			fmt.Println(resp)

			if result.StatusCode == http.StatusOK && resp.AccessTokenExpiresAt.Unix() < time.Now().Unix() {
				t.Errorf("wrong expiration date of AccessToken: %v", resp.AccessTokenExpiresAt)
			}
		})
	}
}

func addChiURLParam(r *http.Request, key, value string) *http.Request {
	ctx := chi.NewRouteContext()

	ctx.URLParams.Add(key, value)

	return r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, ctx))
}
