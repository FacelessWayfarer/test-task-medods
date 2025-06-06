package tokens

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var ErrGeneratingTokenID = errors.New("internal error")

type UserClaims struct {
	TokenID uuid.UUID
	UserID  uuid.UUID
	UserIP  string
	jwt.RegisteredClaims
}

func NewUserClaims(id uuid.UUID, ip string, duration time.Duration) (*UserClaims, error) {
	tokenID, err := uuid.NewRandom()
	if err != nil {
		return nil, ErrGeneratingTokenID
	}

	return &UserClaims{
		TokenID: tokenID,
		UserID:  id,
		UserIP:  ip,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
		},
	}, nil
}
