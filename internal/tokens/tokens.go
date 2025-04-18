package tokens

import (
	"crypto/sha256"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JWTMaker struct {
	privateKey string
}

func NewJWT(privateKey string) *JWTMaker {
	return &JWTMaker{
		privateKey: privateKey,
	}
}

func (maker *JWTMaker) CreateToken(userid uuid.UUID, ip string, duration time.Duration) (string, *UserClaims, error) {
	claims, err := NewUserClaims(userid, ip, duration)
	if err != nil {
		return "", nil, err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	tokenString, err := token.SignedString([]byte(maker.privateKey))

	if err != nil {
		return "", nil, err
	}

	return tokenString, claims, nil
}
func CreateRefreshTokenHash(refreshToken string) []byte {
	refreshTokenArray := sha256.Sum256([]byte(refreshToken))
	refreshTokenHash := refreshTokenArray[:]
	return refreshTokenHash
}

func (maker *JWTMaker) VerifyToken(tokenString string) (*UserClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(t *jwt.Token) (interface{}, error) {
		_, ok := t.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, fmt.Errorf("error parsing token")
		}
		return []byte(maker.privateKey), nil
	})
	if err != nil {
		return nil, fmt.Errorf("invalid token")
	}
	claims, ok := token.Claims.(*UserClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token")
	}
	return claims, nil
}
