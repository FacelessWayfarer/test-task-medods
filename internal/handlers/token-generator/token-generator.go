package tokengenerator

import (
	"context"
	"encoding/base64"
	"errors"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/FacelessWayfarer/test-task-medods/internal/database"
	"github.com/FacelessWayfarer/test-task-medods/internal/tokens"
	"github.com/FacelessWayfarer/test-task-medods/pkg/logging"
	"github.com/FacelessWayfarer/test-task-medods/pkg/response"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type Response struct {
	response.Response
	AccessToken           string    `json:"AccessToken,omitempty"`
	RefreshToken          string    `json:"RefreshToken,omitempty"`
	AccessTokenExpiresAt  time.Time `json:"AccessTokenExpiresAt,omitempty"`
	RefreshTokenExpiresAt time.Time `json:"RefreshTokenExpiresAt,omitempty"`
}

// interface of the database decleared where it is used
type UserGetter interface {
	GetUser(ctx context.Context, userID string) (database.User, error)
}
type SessionSaver interface {
	SaveSession(ctx context.Context, session *database.Session) error
}

// ttl of the tokens
const (
	jwtAccessTokenDuration  time.Duration = 1 * time.Minute
	jwtRefreshTokenDuration time.Duration = 24 * time.Hour
)

// New generates a pair of access and refresh JWT tokens for user with id in URL parametr {user_id}
func New(ctx context.Context, userGetter UserGetter, sessionSaver SessionSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logging.L(ctx).Info("generating jwt access and refresh tokens")

		userID := chi.URLParam(r, "user_id")
		if userID == "" {
			logging.L(ctx).Info("empty user_id parameter")

			render.JSON(w, r, response.Error("invalid request"))

			return
		}

		user, err := userGetter.GetUser(ctx, userID)
		if err != nil {
			if errors.Is(err, database.ErrUserNotFound) {
				logging.L(ctx).Error("user not found", logging.ErrorField(err))

				render.JSON(w, r, response.Error("user not found"))

				return
			}
			logging.L(ctx).Error("failed to find user", logging.ErrorField(err))

			render.JSON(w, r, response.Error("user not found"))

			return
		}

		//JWT is created with secret string that is stored in ENV variable JWT_SECRET
		secretString := os.Getenv("JWT_SECRET")
		jwtMaker := tokens.NewJWT(secretString)

		//getting ip addres from request
		fullip := r.RemoteAddr
		stringIp := strings.Split(fullip, ":")
		ip := stringIp[0]

		refreshToken, refreshTokenClaims, err := jwtMaker.CreateToken(user.ID, ip, jwtRefreshTokenDuration)
		if err != nil {
			logging.L(ctx).Error("error creating refresh token", logging.ErrorField(err))

			render.JSON(w, r, response.Error("internal error"))

			return
		}

		refreshTokenHash := tokens.CreateRefreshTokenHash(refreshToken)

		accessToken, accessTokenClaims, err := jwtMaker.CreateToken(user.ID, ip, jwtAccessTokenDuration)
		if err != nil {
			logging.L(ctx).Error("error creating access token", logging.ErrorField(err))

			render.JSON(w, r, response.Error("internal error"))

			return
		}

		err = sessionSaver.SaveSession(ctx, &database.Session{ID: refreshTokenClaims.RegisteredClaims.ID,
			UserIp:        refreshTokenClaims.UserIP,
			UserId:        refreshTokenClaims.UserID,
			Refresh_token: string(refreshTokenHash),
			ExpiresAt:     refreshTokenClaims.ExpiresAt.Time,
		})
		if err != nil {

			logging.L(ctx).Error("error saving session token", logging.ErrorField(err))

			render.JSON(w, r, response.Error("internal error"))

			return
		}
		base64RefreshToken := base64.StdEncoding.EncodeToString([]byte(refreshToken))

		logging.L(ctx).Info("successfuly generated and saved jwt access and refresh tokens")

		w.Header().Set("Content-Type", "application/json")
		render.JSON(w, r, Response{
			Response:              response.OK(),
			AccessToken:           accessToken,
			RefreshToken:          base64RefreshToken,
			AccessTokenExpiresAt:  accessTokenClaims.RegisteredClaims.ExpiresAt.Time,
			RefreshTokenExpiresAt: refreshTokenClaims.RegisteredClaims.ExpiresAt.Time,
		})
	}
}
