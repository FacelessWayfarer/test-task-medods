package service

import (
	"context"
	"encoding/base64"
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

func TestService_RefreshTokens(t *testing.T) {
	type testCase struct {
		name              string
		ctx               context.Context
		mockError         error
		input             models.TokensToRefresh
		userIP            string
		mockUser          *models.User
		mockAccessClaims  *tokens.UserClaims
		mockRefreshClaims *tokens.UserClaims
		mockSession       *models.Session
		expectedError     error
		mockSetup         func(tt testCase) *Service
	}
	tests := []testCase{
		{
			name: "happy_path",
			ctx:  context.Background(),
			input: models.TokensToRefresh{AccessToken: "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJUb2tlbklEIjoiMjE1OGZlYjEtZTVkYS00MjNjLThmZWYtNDQ1ODNhYzVjNGFmIiwiVXNlcklEIjoiMTcxNmRhYWItNTg2OC00NzdlLTlmNTEtMGRmMmEwZTkyNWI3IiwiVXNlcklQIjoiMTcyLjE4LjAuMSIsImV4cCI6MTc0NzkyMzcxMCwiaWF0IjoxNzQ3OTIzNjUwfQ.Sn4tSHBBPvalKiU23ib1lImvPEdWNQrTYqshUYoSFXxKnLm2xGLXWJej0t6D2N8SZImd6Lv8PJBs1aTAuNNBxQ",
				Base64RefreshToken: "ZXlKaGJHY2lPaUpJVXpVeE1pSXNJblI1Y0NJNklrcFhWQ0o5LmV5SlViMnRsYmtsRUlqb2lNamM0TVRnM09XSXRaR0UxTkMwME5ETTRMVGxqWmpndFlXRXpNREUzWW1RMVptRmlJaXdpVlhObGNrbEVJam9pTVRjeE5tUmhZV0l0TlRnMk9DMDBOemRsTFRsbU5URXRNR1JtTW1Fd1pUa3lOV0kzSWl3aVZYTmxja2xRSWpvaU1UY3lMakU0TGpBdU1TSXNJbVY0Y0NJNk1UYzBPREF4TURBMU1Dd2lhV0YwSWpveE56UTNPVEl6TmpVd2ZRLmdXMURRMXNoemN2aGloOUtVWUR0X0IteWpCTmNqVHhLVFVReW9xa1JRdU5aLW04QUtxTFM1eGpoZUVwalIzWDNzS21yNHlwZHR6U2FqbE9ia3lvWDd3"},
			userIP: "8.8.8.8",
			mockUser: &models.User{
				ID:    uuid.MustParse("736f6d65-5f72-616e-646f-6d5f75756964"),
				Email: "someone@gmail.com",
			},
			mockAccessClaims: &tokens.UserClaims{
				UserID:           uuid.MustParse("736f6d65-5f72-616e-646f-6d5f75756964"),
				RegisteredClaims: jwt.RegisteredClaims{ExpiresAt: &jwt.NumericDate{Time: time.Now().Add(time.Minute)}}},
			mockRefreshClaims: &tokens.UserClaims{
				TokenID:          uuid.MustParse("555f6d65-5f72-616e-646f-6d5f75756964"),
				UserID:           uuid.MustParse("736f6d65-5f72-616e-646f-6d5f75756964"),
				UserIP:           "8.8.8.8",
				RegisteredClaims: jwt.RegisteredClaims{ExpiresAt: &jwt.NumericDate{Time: time.Now().Add(time.Hour)}}},
			mockSession: &models.Session{
				ID:           uuid.MustParse("555f6d65-5f72-616e-646f-6d5f75756964"),
				UserIP:       "8.8.8.8",
				UserID:       uuid.MustParse("736f6d65-5f72-616e-646f-6d5f75756964"),
				RefreshToken: string([]byte{}),
				ExpiredAt:    time.Now().Add(time.Hour),
			},
			expectedError: nil,
			mockSetup: func(tt testCase) *Service {
				userStorage := mocks.NewUserStorage(t)

				userStorage.On("GetUser", tt.ctx, fmt.Sprint(tt.mockRefreshClaims.UserID)).
					Return(tt.mockUser, nil)

				sessionStorage := mocks.NewSessionStorage(t)

				sessionStorage.On("GetSession", tt.ctx, tt.mockRefreshClaims.TokenID).
					Return(tt.mockSession, nil)

				sessionStorage.On("SaveSession", tt.ctx, *tt.mockSession).
					Return(nil)

				tokenCreator := mocks.NewTokenCreator(t)

				tokenCreator.On("CreateToken", tt.mockUser.ID, tt.userIP, jwtRefreshTokenDuration).
					Return("refreshToken", tt.mockRefreshClaims, nil)

				tokenCreator.On("CreateToken", tt.mockUser.ID, tt.userIP, jwtAccessTokenDuration).
					Return("accessToken", tt.mockAccessClaims, nil)

				tokenCreator.On("CreateRefreshTokenHash", "refreshToken").
					Return([]byte{})

				refreshToken, err := base64.StdEncoding.DecodeString(tt.input.Base64RefreshToken)
				if err != nil {
					t.Fatal(err)
				}

				tokenCreator.On("VerifyToken", string(refreshToken)).
					Return(tt.mockRefreshClaims, nil)

				tokenCreator.On("VerifyToken", tt.input.AccessToken).
					Return(tt.mockAccessClaims, nil)

				return &Service{
					userStorage:    userStorage,
					sessionStorage: sessionStorage,
					tokenCreator:   tokenCreator,
					logger:         log.New(os.Stdout, "test:", log.LstdFlags),
				}
			},
		},
		{
			name:      "error_token_expired",
			ctx:       context.Background(),
			mockError: models.ErrTokenExpired,
			input: models.TokensToRefresh{AccessToken: "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJUb2tlbklEIjoiMjE1OGZlYjEtZTVkYS00MjNjLThmZWYtNDQ1ODNhYzVjNGFmIiwiVXNlcklEIjoiMTcxNmRhYWItNTg2OC00NzdlLTlmNTEtMGRmMmEwZTkyNWI3IiwiVXNlcklQIjoiMTcyLjE4LjAuMSIsImV4cCI6MTc0NzkyMzcxMCwiaWF0IjoxNzQ3OTIzNjUwfQ.Sn4tSHBBPvalKiU23ib1lImvPEdWNQrTYqshUYoSFXxKnLm2xGLXWJej0t6D2N8SZImd6Lv8PJBs1aTAuNNBxQ",
				Base64RefreshToken: "ZXlKaGJHY2lPaUpJVXpVeE1pSXNJblI1Y0NJNklrcFhWQ0o5LmV5SlViMnRsYmtsRUlqb2lNamM0TVRnM09XSXRaR0UxTkMwME5ETTRMVGxqWmpndFlXRXpNREUzWW1RMVptRmlJaXdpVlhObGNrbEVJam9pTVRjeE5tUmhZV0l0TlRnMk9DMDBOemRsTFRsbU5URXRNR1JtTW1Fd1pUa3lOV0kzSWl3aVZYTmxja2xRSWpvaU1UY3lMakU0TGpBdU1TSXNJbVY0Y0NJNk1UYzBPREF4TURBMU1Dd2lhV0YwSWpveE56UTNPVEl6TmpVd2ZRLmdXMURRMXNoemN2aGloOUtVWUR0X0IteWpCTmNqVHhLVFVReW9xa1JRdU5aLW04QUtxTFM1eGpoZUVwalIzWDNzS21yNHlwZHR6U2FqbE9ia3lvWDd3"},
			mockAccessClaims: &tokens.UserClaims{UserID: uuid.MustParse("736f6d65-5f72-616e-646f-6d5f75756964"),
				RegisteredClaims: jwt.RegisteredClaims{ExpiresAt: &jwt.NumericDate{Time: time.Now().Add(-time.Minute)}}},
			mockRefreshClaims: &tokens.UserClaims{UserID: uuid.MustParse("736f6d65-5f72-616e-646f-6d5f75756964"),
				RegisteredClaims: jwt.RegisteredClaims{ExpiresAt: &jwt.NumericDate{Time: time.Now().Add(-time.Hour)}}},
			mockSession:   &models.Session{ExpiredAt: time.Now().Add(-time.Hour)},
			expectedError: models.ErrTokenExpired,
			mockSetup: func(tt testCase) *Service {
				sessionStorage := mocks.NewSessionStorage(t)

				sessionStorage.On("GetSession", tt.ctx, mock.Anything).
					Return(tt.mockSession, nil)

				tokenCreator := mocks.NewTokenCreator(t)

				tokenCreator.On("VerifyToken", tt.input.AccessToken).
					Return(tt.mockAccessClaims, nil)

				refreshToken, err := base64.StdEncoding.DecodeString(tt.input.Base64RefreshToken)
				if err != nil {
					t.Fatal(err)
				}

				tokenCreator.On("VerifyToken", string(refreshToken)).
					Return(tt.mockRefreshClaims, nil)

				return &Service{
					sessionStorage: sessionStorage,
					tokenCreator:   tokenCreator,
					logger:         log.New(os.Stdout, "test:", log.LstdFlags),
				}
			},
		},
		{
			name:      "error_invalid_tokens",
			ctx:       context.Background(),
			mockError: models.ErrInvalidTokens,
			input: models.TokensToRefresh{AccessToken: "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJUb2tlbklEIjoiMjE1OGZlYjEtZTVkYS00MjNjLThmZWYtNDQ1ODNhYzVjNGFmIiwiVXNlcklEIjoiMTcxNmRhYWItNTg2OC00NzdlLTlmNTEtMGRmMmEwZTkyNWI3IiwiVXNlcklQIjoiMTcyLjE4LjAuMSIsImV4cCI6MTc0NzkyMzcxMCwiaWF0IjoxNzQ3OTIzNjUwfQ.Sn4tSHBBPvalKiU23ib1lImvPEdWNQrTYqshUYoSFXxKnLm2xGLXWJej0t6D2N8SZImd6Lv8PJBs1aTAuNNBxQ",
				Base64RefreshToken: "ZXlKaGJHY2lPaUpJVXpVeE1pSXNJblI1Y0NJNklrcFhWQ0o5LmV5SlViMnRsYmtsRUlqb2lNamM0TVRnM09XSXRaR0UxTkMwME5ETTRMVGxqWmpndFlXRXpNREUzWW1RMVptRmlJaXdpVlhObGNrbEVJam9pTVRjeE5tUmhZV0l0TlRnMk9DMDBOemRsTFRsbU5URXRNR1JtTW1Fd1pUa3lOV0kzSWl3aVZYTmxja2xRSWpvaU1UY3lMakU0TGpBdU1TSXNJbVY0Y0NJNk1UYzBPREF4TURBMU1Dd2lhV0YwSWpveE56UTNPVEl6TmpVd2ZRLmdXMURRMXNoemN2aGloOUtVWUR0X0IteWpCTmNqVHhLVFVReW9xa1JRdU5aLW04QUtxTFM1eGpoZUVwalIzWDNzS21yNHlwZHR6U2FqbE9ia3lvWDd3"},
			mockAccessClaims: &tokens.UserClaims{UserID: uuid.MustParse("736f6d65-5f72-616e-646f-6d5f75756333"),
				RegisteredClaims: jwt.RegisteredClaims{ExpiresAt: &jwt.NumericDate{Time: time.Now().Add(time.Minute)}}},
			mockRefreshClaims: &tokens.UserClaims{UserID: uuid.MustParse("736f6d65-5f72-616e-646f-6d5f75756964"),
				RegisteredClaims: jwt.RegisteredClaims{ExpiresAt: &jwt.NumericDate{Time: time.Now().Add(time.Hour)}}},
			expectedError: models.ErrInvalidTokens,
			mockSetup: func(tt testCase) *Service {
				tokenCreator := mocks.NewTokenCreator(t)

				tokenCreator.On("VerifyToken", tt.input.AccessToken).
					Return(tt.mockAccessClaims, nil)

				refreshToken, err := base64.StdEncoding.DecodeString(tt.input.Base64RefreshToken)
				if err != nil {
					t.Fatal(err)
				}

				tokenCreator.On("VerifyToken", string(refreshToken)).
					Return(tt.mockRefreshClaims, nil)

				return &Service{
					tokenCreator: tokenCreator,
					logger:       log.New(os.Stdout, "test:", log.LstdFlags),
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := tt.mockSetup(tt)

			_, err := s.UpdateTokens(tt.ctx, tt.input, tt.userIP)
			if !errors.Is(err, tt.expectedError) {
				t.Errorf("got %v, want %v", err, tt.expectedError)
			}
		})
	}
}
