package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/FacelessWayfarer/test-task-medods/internal/service/mocks"
	"github.com/FacelessWayfarer/test-task-medods/internal/service/models"
)

func TestService_RefreshTokens(t *testing.T) {
	type testCase struct {
		name               string
		ctx                context.Context
		IP                 string
		input              models.TokensToRefresh
		output             *models.RefreshedTokens
		mockError          error
		expectedStatusCode int
		mockSetup          func(tt testCase) *Handler
	}
	tests := []testCase{
		{
			name: "Happy_path",
			ctx:  context.Background(),
			IP:   "111.111.111.111",
			input: models.TokensToRefresh{AccessToken: "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJUb2tlbklEIjoiMjE1OGZlYjEtZTVkYS00MjNjLThmZWYtNDQ1ODNhYzVjNGFmIiwiVXNlcklEIjoiMTcxNmRhYWItNTg2OC00NzdlLTlmNTEtMGRmMmEwZTkyNWI3IiwiVXNlcklQIjoiMTcyLjE4LjAuMSIsImV4cCI6MTc0NzkyMzcxMCwiaWF0IjoxNzQ3OTIzNjUwfQ.Sn4tSHBBPvalKiU23ib1lImvPEdWNQrTYqshUYoSFXxKnLm2xGLXWJej0t6D2N8SZImd6Lv8PJBs1aTAuNNBxQ",
				Base64RefreshToken: "ZXlKaGJHY2lPaUpJVXpVeE1pSXNJblI1Y0NJNklrcFhWQ0o5LmV5SlViMnRsYmtsRUlqb2lNamM0TVRnM09XSXRaR0UxTkMwME5ETTRMVGxqWmpndFlXRXpNREUzWW1RMVptRmlJaXdpVlhObGNrbEVJam9pTVRjeE5tUmhZV0l0TlRnMk9DMDBOemRsTFRsbU5URXRNR1JtTW1Fd1pUa3lOV0kzSWl3aVZYTmxja2xRSWpvaU1UY3lMakU0TGpBdU1TSXNJbVY0Y0NJNk1UYzBPREF4TURBMU1Dd2lhV0YwSWpveE56UTNPVEl6TmpVd2ZRLmdXMURRMXNoemN2aGloOUtVWUR0X0IteWpCTmNqVHhLVFVReW9xa1JRdU5aLW04QUtxTFM1eGpoZUVwalIzWDNzS21yNHlwZHR6U2FqbE9ia3lvWDd3"},
			output: &models.RefreshedTokens{
				AccessToken:           "access",
				RefreshToken:          "refresh",
				AccessTokenExpiresAt:  time.Now().Add(jwtTestAccessTokenDuration),
				RefreshTokenExpiresAt: time.Now().Add(jwtTestRefreshTokenDuration),
			},
			mockError:          nil,
			expectedStatusCode: http.StatusOK,
			mockSetup: func(tt testCase) *Handler {
				service := mocks.NewIService(t)

				service.On("UpdateTokens", tt.ctx, tt.input, tt.IP).
					Return(tt.output, tt.mockError)

				return &Handler{
					service: service,
					logger:  log.New(os.Stdout, "test:", log.LstdFlags),
				}
			},
		},
		{
			name: "token_expired",
			ctx:  context.Background(),
			IP:   "111.111.111.111",
			input: models.TokensToRefresh{AccessToken: "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJUb2tlbklEIjoiMjE1OGZlYjEtZTVkYS00MjNjLThmZWYtNDQ1ODNhYzVjNGFmIiwiVXNlcklEIjoiMTcxNmRhYWItNTg2OC00NzdlLTlmNTEtMGRmMmEwZTkyNWI3IiwiVXNlcklQIjoiMTcyLjE4LjAuMSIsImV4cCI6MTc0NzkyMzcxMCwiaWF0IjoxNzQ3OTIzNjUwfQ.Sn4tSHBBPvalKiU23ib1lImvPEdWNQrTYqshUYoSFXxKnLm2xGLXWJej0t6D2N8SZImd6Lv8PJBs1aTAuNNBxQ",
				Base64RefreshToken: "ZXlKaGJHY2lPaUpJVXpVeE1pSXNJblI1Y0NJNklrcFhWQ0o5LmV5SlViMnRsYmtsRUlqb2lNamM0TVRnM09XSXRaR0UxTkMwME5ETTRMVGxqWmpndFlXRXpNREUzWW1RMVptRmlJaXdpVlhObGNrbEVJam9pTVRjeE5tUmhZV0l0TlRnMk9DMDBOemRsTFRsbU5URXRNR1JtTW1Fd1pUa3lOV0kzSWl3aVZYTmxja2xRSWpvaU1UY3lMakU0TGpBdU1TSXNJbVY0Y0NJNk1UYzBPREF4TURBMU1Dd2lhV0YwSWpveE56UTNPVEl6TmpVd2ZRLmdXMURRMXNoemN2aGloOUtVWUR0X0IteWpCTmNqVHhLVFVReW9xa1JRdU5aLW04QUtxTFM1eGpoZUVwalIzWDNzS21yNHlwZHR6U2FqbE9ia3lvWDd3"},
			mockError:          models.ErrTokenExpired,
			expectedStatusCode: http.StatusInternalServerError,
			mockSetup: func(tt testCase) *Handler {
				service := mocks.NewIService(t)

				service.On("UpdateTokens", tt.ctx, tt.input, tt.IP).
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

			var buf bytes.Buffer

			if err := json.NewEncoder(&buf).Encode(tt.input); err != nil {
				log.Fatal(err)
			}

			r := httptest.NewRequest("POST", "/tokens/", &buf)

			r.RemoteAddr = tt.IP

			w := httptest.NewRecorder()

			h.RefreshTokens(w, r)

			result := w.Result()
			if result.StatusCode != tt.expectedStatusCode {
				t.Errorf("got %v, want %v", result.StatusCode, tt.expectedStatusCode)
			}

			var resp RefreshResponse

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
