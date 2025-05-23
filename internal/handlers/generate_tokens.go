package handlers

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/google/uuid"

	"github.com/FacelessWayfarer/test-task-medods/internal/models"
	"github.com/FacelessWayfarer/test-task-medods/pkg/response"
)

const (
	jwtRefreshTokenDuration time.Duration = 24 * time.Hour
)

// ShowAccount godoc
// @Summary      GenerateTokens
// @Description  Generates access and refresh tokens
// @Tags         Tokens
// @Accept       json
// @Produce      json
// @Param        user_id path string true "user_id"
// @Success      200  {object}  GenResponse
// @Failure      400  {object}  response.Response
// @Router       /tokens/{user_id} [GET]
func (h *Handler) GenerateTokens(w http.ResponseWriter, r *http.Request) {
	h.logger.Println("Generating jwt access and refresh tokens")

	userID, ip, err := checkRequest(*r)
	if err != nil {
		h.logger.Printf("error checking request: %v", err)

		w.WriteHeader(http.StatusBadRequest)

		render.JSON(w, r, response.Error("invalid request"))

		return
	}

	ctx := r.Context()

	user, err := h.getUserFromStorage(ctx, userID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		render.JSON(w, r, response.Error("user not found"))

		return
	}

	base64RefreshToken, refreshTokenExpiresAt, err := h.createAndSaveRefreshToken(ctx, user.ID, ip)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		render.JSON(w, r, response.Error("internal error"))

		return
	}

	accessToken, accessTokenExpiresAt, err := h.createAccessToken(user.ID, ip)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		render.JSON(w, r, response.Error("internal error"))

		return
	}

	h.logger.Println("Successfully generated and saved jwt access and refresh tokens")

	w.Header().Set("Content-Type", "application/json")

	render.JSON(w, r, GenResponse{
		Response:              response.OK(),
		AccessToken:           accessToken,
		RefreshToken:          base64RefreshToken,
		AccessTokenExpiresAt:  *accessTokenExpiresAt,
		RefreshTokenExpiresAt: *refreshTokenExpiresAt,
	})

}

func checkRequest(r http.Request) (string, string, error) {
	userID := chi.URLParam(&r, "user_id")
	if userID == "" {

		return "", "", errors.New("empty user_id")
	}

	fullip := r.RemoteAddr

	stringIp := strings.Split(fullip, ":")

	ip := stringIp[0]

	return userID, ip, nil
}

func (h *Handler) getUserFromStorage(ctx context.Context, userID string) (*models.User, error) {
	user, err := h.userStorage.GetUser(ctx, userID)
	if err != nil {
		if errors.Is(err, models.ErrUserNotFound) {
			h.logger.Printf("failed to find user: %v", err)

			return nil, fmt.Errorf("could not find user: %w", err)
		}
		h.logger.Printf("failed to find user: %v", err)

		return nil, fmt.Errorf("could not find user: %w", err)
	}

	return user, nil
}

func (h *Handler) createAndSaveRefreshToken(ctx context.Context, userID uuid.UUID, ip string) (string, *time.Time, error) {
	refreshToken, refreshTokenClaims, err := h.tokenCreator.CreateToken(userID, ip, jwtRefreshTokenDuration)
	if err != nil {
		h.logger.Printf("error creating refresh token: %v", err)

		return "", nil, fmt.Errorf("could not create refresh token: %w", err)
	}

	refreshTokenHash := h.tokenCreator.CreateRefreshTokenHash(refreshToken)

	expiresAt := refreshTokenClaims.ExpiresAt.Time

	if err = h.sessionStorage.SaveSession(ctx, models.Session{ID: refreshTokenClaims.TokenID,
		UserIp:       refreshTokenClaims.UserIP,
		UserId:       refreshTokenClaims.UserID,
		RefreshToken: string(refreshTokenHash),
		ExpiredAt:    expiresAt,
	}); err != nil {
		h.logger.Printf("error saving session token: %v", err)

		return "", nil, fmt.Errorf("could not save refresh token: %w", err)
	}

	return base64.StdEncoding.EncodeToString([]byte(refreshToken)), &expiresAt, nil
}

func (h *Handler) createAccessToken(userID uuid.UUID, ip string) (string, *time.Time, error) {
	refreshToken, refreshTokenClaims, err := h.tokenCreator.CreateToken(userID, ip, jwtAccessTokenDuration)
	if err != nil {
		h.logger.Printf("error creating access token: %v", err)

		return "", nil, fmt.Errorf("could not create access token: %w", err)
	}

	return refreshToken, &refreshTokenClaims.ExpiresAt.Time, nil
}
