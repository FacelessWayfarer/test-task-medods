package tokens

import (
	"errors"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type UserClaims struct {
	UserID uuid.UUID `json:"user_id"`
	UserIP string    `json:"user_ip"`
	jwt.RegisteredClaims
}

func NewUserClaims(id uuid.UUID, ip string, duration time.Duration) (*UserClaims, error) {
	tokenID, err := uuid.NewRandom()
	if err != nil {
		log.Printf("error generating tokenID: %e", err)
		return nil, errors.New("internal error")
	}
	return &UserClaims{
		UserID: id,
		UserIP: ip,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        tokenID.String(),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
		},
	}, nil
}
