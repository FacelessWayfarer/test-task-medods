package tokenrefresher

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/FacelessWayfarer/test-task-medods/internal/database"
	"github.com/FacelessWayfarer/test-task-medods/internal/tokens"
	"github.com/FacelessWayfarer/test-task-medods/pkg/logging"
	"github.com/FacelessWayfarer/test-task-medods/pkg/response"
	"github.com/go-chi/render"
)

type RefreshTokensReq struct {
	Base64RefreshToken string `json:"encoded_refresh_token"`
}
type Response struct {
	response.Response
	AccessToken          string    `json:"AccessToken,omitempty"`
	AccessTokenExpiresAt time.Time `json:"AccessTokenExpiresAt,omitempty"`
}

// interface of the database decleared where it is used
type UserGetter interface {
	GetUser(ctx context.Context, userID string) (database.User, error)
}
type SessionGetter interface {
	GetSession(ctx context.Context, sessionid string) (database.Session, error)
}

const (
	jwtAccessTokenDuration time.Duration = 1 * time.Minute
)

// New refershes access token with it's response token, declared in request body. Example of request body, format json: {"encoded_refresh_token":"encodedreshtokenhere"}
func New(ctx context.Context, userGetter UserGetter, sessionGetter SessionGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logging.L(ctx).Info("refreshing tokens")
		var req RefreshTokensReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			logging.L(ctx).Error("failed to decode request body", logging.ErrorField(err))
			render.JSON(w, r, response.Error("error decoding request body"))
			return
		}
		secretString := os.Getenv("JWT_SECRET")
		jwtMaker := tokens.NewJWT(secretString)

		refreshToken, err := base64.StdEncoding.DecodeString(req.Base64RefreshToken)
		if err != nil {
			logging.L(ctx).Error("failed to decode token", logging.ErrorField(err))
			render.JSON(w, r, response.Error("internal error"))
			return
		}

		refreshClaims, err := jwtMaker.VerifyToken(string(refreshToken))
		if err != nil {
			logging.L(ctx).Error("failed to verify token", logging.ErrorField(err))
			render.JSON(w, r, response.Error("error verifying token"))
			return
		}
		session, err := sessionGetter.GetSession(ctx, refreshClaims.RegisteredClaims.ID)
		if err != nil {
			logging.L(ctx).Error("failed get session", logging.ErrorField(err))
			render.JSON(w, r, response.Error("internal error"))
			return
		}

		if refreshClaims.ExpiresAt.Unix() < time.Now().Unix() || session.ExpiresAt.Unix() < time.Now().Unix() {
			logging.L(ctx).Error("token expired", logging.ErrorField(err))
			render.JSON(w, r, response.Error("token expired"))
			return
		}

		fullip := r.RemoteAddr
		stringIp := strings.Split(fullip, ":")
		ip := stringIp[0]

		if session.UserIp != ip {
			userid := fmt.Sprint(refreshClaims.UserID)
			user, err := userGetter.GetUser(ctx, userid)
			if err != nil {
				logging.L(ctx).Error("failed to find user", logging.ErrorField(err))

				render.JSON(w, r, response.Error("user not found"))

				return
			}
			logging.L(ctx).Info("user ip changed, sending email to", logging.StringField("user_email", user.Email))
			// Mock interface to send email
			type EmailSender interface {
				SendEmail(string) error
			}
			// Should be passed to hendler, same as other interfaces

			// var emailSender EmailSender
			// err = emailSender.SendEmail(user.Email)
			// if err != nil {
			// 	///
			// }
		}

		accessToken, accessTokenClaims, err := jwtMaker.CreateToken(refreshClaims.UserID, ip, jwtAccessTokenDuration)
		if err != nil {
			logging.L(ctx).Error("error creating access token", logging.ErrorField(err))

			render.JSON(w, r, response.Error("internal error"))

			return
		}

		logging.L(ctx).Info("successfuly refreshed jwt access token")

		w.Header().Set("Content-Type", "application/json")
		render.JSON(w, r, Response{
			Response:             response.OK(),
			AccessToken:          accessToken,
			AccessTokenExpiresAt: accessTokenClaims.ExpiresAt.Time,
		})
	}
}
