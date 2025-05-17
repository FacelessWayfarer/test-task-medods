package tokens

import (
	"errors"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type UserClaims struct {
	TokenID uuid.UUID
	UserID  uuid.UUID
	UserIP  string
	jwt.RegisteredClaims
}

func NewUserClaims(id uuid.UUID, ip string, duration time.Duration) (*UserClaims, error) {
	tokenID, err := uuid.NewRandom()
	if err != nil {

		log.Printf("error generating tokenID: %e", err)

		return nil, errors.New("internal error")
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
