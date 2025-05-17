package handlers

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/render"

	"github.com/FacelessWayfarer/test-task-medods/internal/models"
	"github.com/FacelessWayfarer/test-task-medods/pkg/response"
	"github.com/FacelessWayfarer/test-task-medods/pkg/tokens"
)

const (
	jwtAccessTokenDuration time.Duration = 1 * time.Minute
)

// ShowAccount godoc
// @Summary      RefreshTokens
// @Description  Refreshes access and refresh tokens
// @Tags         Tokens
// @Accept       json
// @Produce      json
// @Param        tokens body RefreshTokensReq true "access and refresh tokens"
// @Success      200  {object}  RefreshResponse
// @Failure      400  {object}  response.Response
// @Router       /tokens/ [POST]
func (h *Handler) RefreshTokens(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.logger.Println("refreshing tokens")

		req, err := decodeRefreshRequestBody(r.Body)
		if err != nil {
			h.logger.Printf("error checking request: %v", err)

			render.JSON(w, r, response.Error("invalid request"))

			return
		}

		refreshClaims, session, err := h.verifyTokenAndFetchSession(ctx, req.Base64RefreshToken)
		if err != nil {
			render.JSON(w, r, response.Error("internal error"))

			return
		}

		if err = h.checkExpiration(refreshClaims.ExpiresAt.Unix(), session.ExpiredAt.Unix()); err != nil {
		}

		stringIp := strings.Split(r.RemoteAddr, ":")

		ip := stringIp[0]

		user, err := h.getUserFromStorage(ctx, fmt.Sprint(refreshClaims.UserID))
		if err != nil {
			render.JSON(w, r, response.Error("internal error"))

			return
		}

		if session.UserIp != ip {
			if err = h.sendEmail(user.Email); err != nil {
				render.JSON(w, r, response.Error("internal error"))

				return
			}
		}

		base64RefreshToken, refreshTokenExpiresAt, err := h.createAndSaveRefreshToken(ctx, user.ID, ip)
		if err != nil {
			render.JSON(w, r, response.Error("internal error"))

			return
		}

		accessToken, accessTokenExpiresAt, err := h.createAccessToken(user.ID, ip)
		if err != nil {
			render.JSON(w, r, response.Error("internal error"))

			return
		}

		h.logger.Println("successfully refreshed jwt access token")

		w.Header().Set("Content-Type", "application/json")

		render.JSON(w, r, RefreshResponse{
			Response:              response.OK(),
			AccessToken:           accessToken,
			RefreshToken:          base64RefreshToken,
			AccessTokenExpiresAt:  *accessTokenExpiresAt,
			RefreshTokenExpiresAt: *refreshTokenExpiresAt,
		})
	}
}

func decodeRefreshRequestBody(body io.ReadCloser) (*RefreshTokensReq, error) {
	var req RefreshTokensReq

	if err := json.NewDecoder(body).Decode(&req); err != nil {

		return nil, fmt.Errorf("could not decode request body: %v", err)
	}

	return &req, nil
}

func (h *Handler) verifyTokenAndFetchSession(ctx context.Context, requstToken string) (*tokens.UserClaims, *models.Session, error) {
	refreshToken, err := base64.StdEncoding.DecodeString(requstToken)
	if err != nil {
		h.logger.Printf("failed to decode refresh token from request: %v", err)

		return nil, nil, fmt.Errorf("could not decode refresh token: %v", err)
	}

	refreshClaims, err := h.tokenCreator.VerifyToken(string(refreshToken))
	if err != nil {
		h.logger.Printf("failed to verify token: %v", err)

		return nil, nil, fmt.Errorf("could not verify token: %v", err)
	}
	session, err := h.sessionStorage.GetSession(ctx, refreshClaims.TokenID)
	if err != nil {
		h.logger.Printf("failed to fetch session: %v", err)

		return nil, nil, fmt.Errorf("could not fetch session: %v", err)
	}

	return refreshClaims, session, nil
}

func (h *Handler) sendEmail(userEmail string) error {
	h.logger.Printf("user ip changed, sending email to: %v", userEmail)

	if err := h.emailSender.SendEmail(userEmail); err != nil {
		return fmt.Errorf("could not send email: %v", err)
	}

	return nil
}

func (h *Handler) checkExpiration(time1 int64, time2 int64) error {
	if time1 < time.Now().Unix() || time2 < time.Now().Unix() {
		h.logger.Println("token expired")

		return errors.New("token expired")
	}

	return nil
}
