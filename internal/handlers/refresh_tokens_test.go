package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/mock"

	"github.com/FacelessWayfarer/test-task-medods/internal/handlers/mocks"
	"github.com/FacelessWayfarer/test-task-medods/internal/models"
	"github.com/FacelessWayfarer/test-task-medods/pkg/tokens"
)

func TestHandler_RefreshTokens(t *testing.T) {
	type testCase struct {
		name               string
		mockError          error
		input              *RefreshTokensReq
		mockUser           *models.User
		mockClaims         *tokens.UserClaims
		mockSession        *models.Session
		expectedStatusCode int
		mockSetup          func(tt testCase) *Handler
	}
	tests := []testCase{
		{
			name: "Happy path",
			input: &RefreshTokensReq{AccessToken: "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJUb2tlbklEIjoiMjE1OGZlYjEtZTVkYS00MjNjLThmZWYtNDQ1ODNhYzVjNGFmIiwiVXNlcklEIjoiMTcxNmRhYWItNTg2OC00NzdlLTlmNTEtMGRmMmEwZTkyNWI3IiwiVXNlcklQIjoiMTcyLjE4LjAuMSIsImV4cCI6MTc0NzkyMzcxMCwiaWF0IjoxNzQ3OTIzNjUwfQ.Sn4tSHBBPvalKiU23ib1lImvPEdWNQrTYqshUYoSFXxKnLm2xGLXWJej0t6D2N8SZImd6Lv8PJBs1aTAuNNBxQ",
				Base64RefreshToken: "ZXlKaGJHY2lPaUpJVXpVeE1pSXNJblI1Y0NJNklrcFhWQ0o5LmV5SlViMnRsYmtsRUlqb2lNamM0TVRnM09XSXRaR0UxTkMwME5ETTRMVGxqWmpndFlXRXpNREUzWW1RMVptRmlJaXdpVlhObGNrbEVJam9pTVRjeE5tUmhZV0l0TlRnMk9DMDBOemRsTFRsbU5URXRNR1JtTW1Fd1pUa3lOV0kzSWl3aVZYTmxja2xRSWpvaU1UY3lMakU0TGpBdU1TSXNJbVY0Y0NJNk1UYzBPREF4TURBMU1Dd2lhV0YwSWpveE56UTNPVEl6TmpVd2ZRLmdXMURRMXNoemN2aGloOUtVWUR0X0IteWpCTmNqVHhLVFVReW9xa1JRdU5aLW04QUtxTFM1eGpoZUVwalIzWDNzS21yNHlwZHR6U2FqbE9ia3lvWDd3"},
			mockUser:           &models.User{},
			mockClaims:         &tokens.UserClaims{RegisteredClaims: jwt.RegisteredClaims{ExpiresAt: &jwt.NumericDate{Time: time.Now()}}},
			mockSession:        &models.Session{ExpiredAt: time.Now().Add(time.Hour)},
			expectedStatusCode: http.StatusOK,
			mockSetup: func(tt testCase) *Handler {
				userStorage := mocks.NewUserStorage(t)

				userStorage.On("GetUser", mock.Anything, mock.Anything).
					Return(tt.mockUser, nil)

				sessionStorage := mocks.NewSessionStorage(t)

				sessionStorage.On("GetSession", mock.Anything, mock.Anything).
					Return(tt.mockSession, nil)

				sessionStorage.On("SaveSession", mock.Anything, mock.Anything).
					Return(nil)

				tokenCreator := mocks.NewTokenCreator(t)

				tokenCreator.On("CreateToken", mock.Anything, mock.Anything, mock.Anything).
					Return("", tt.mockClaims, nil)

				tokenCreator.On("CreateRefreshTokenHash", mock.Anything).
					Return([]byte{})

				tokenCreator.On("VerifyToken", mock.Anything).
					Return(tt.mockClaims, nil)

				emailSender := mocks.NewEmailSender(t)

				emailSender.On("SendEmail", mock.Anything).
					Return(nil)

				return &Handler{
					userStorage:    userStorage,
					sessionStorage: sessionStorage,
					tokenCreator:   tokenCreator,
					logger:         log.New(os.Stdout, "test:", log.LstdFlags),
					emailSender:    emailSender,
				}
			},
		},
		{
			name:      "token expired",
			mockError: errors.New("token expired"),
			input: &RefreshTokensReq{AccessToken: "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJUb2tlbklEIjoiMjE1OGZlYjEtZTVkYS00MjNjLThmZWYtNDQ1ODNhYzVjNGFmIiwiVXNlcklEIjoiMTcxNmRhYWItNTg2OC00NzdlLTlmNTEtMGRmMmEwZTkyNWI3IiwiVXNlcklQIjoiMTcyLjE4LjAuMSIsImV4cCI6MTc0NzkyMzcxMCwiaWF0IjoxNzQ3OTIzNjUwfQ.Sn4tSHBBPvalKiU23ib1lImvPEdWNQrTYqshUYoSFXxKnLm2xGLXWJej0t6D2N8SZImd6Lv8PJBs1aTAuNNBxQ",
				Base64RefreshToken: "ZXlKaGJHY2lPaUpJVXpVeE1pSXNJblI1Y0NJNklrcFhWQ0o5LmV5SlViMnRsYmtsRUlqb2lNamM0TVRnM09XSXRaR0UxTkMwME5ETTRMVGxqWmpndFlXRXpNREUzWW1RMVptRmlJaXdpVlhObGNrbEVJam9pTVRjeE5tUmhZV0l0TlRnMk9DMDBOemRsTFRsbU5URXRNR1JtTW1Fd1pUa3lOV0kzSWl3aVZYTmxja2xRSWpvaU1UY3lMakU0TGpBdU1TSXNJbVY0Y0NJNk1UYzBPREF4TURBMU1Dd2lhV0YwSWpveE56UTNPVEl6TmpVd2ZRLmdXMURRMXNoemN2aGloOUtVWUR0X0IteWpCTmNqVHhLVFVReW9xa1JRdU5aLW04QUtxTFM1eGpoZUVwalIzWDNzS21yNHlwZHR6U2FqbE9ia3lvWDd3"},
			mockUser:           &models.User{},
			mockClaims:         &tokens.UserClaims{RegisteredClaims: jwt.RegisteredClaims{ExpiresAt: &jwt.NumericDate{Time: time.Now()}}},
			mockSession:        &models.Session{ExpiredAt: time.Now().Add(-time.Hour)},
			expectedStatusCode: http.StatusBadRequest,
			mockSetup: func(tt testCase) *Handler {
				sessionStorage := mocks.NewSessionStorage(t)

				sessionStorage.On("GetSession", mock.Anything, mock.Anything).
					Return(tt.mockSession, nil)

				tokenCreator := mocks.NewTokenCreator(t)

				tokenCreator.On("VerifyToken", mock.Anything).
					Return(tt.mockClaims, nil)

				return &Handler{
					sessionStorage: sessionStorage,
					tokenCreator:   tokenCreator,
					logger:         log.New(os.Stdout, "test:", log.LstdFlags),
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := tt.mockSetup(tt)

			var buf bytes.Buffer

			if err := json.NewEncoder(&buf).Encode(tt.input); err != nil {
				log.Fatal(err)
			}

			r := httptest.NewRequest("POST", "/tokens/", &buf)

			w := httptest.NewRecorder()

			h.RefreshTokens(w, r)

			fmt.Println(w.Body.String())

			resp := w.Result()
			if resp.StatusCode != tt.expectedStatusCode {
				t.Errorf("got %v, want %v", resp.StatusCode, tt.expectedStatusCode)
			}
		})
	}
}
