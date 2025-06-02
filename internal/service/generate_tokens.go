package service

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/FacelessWayfarer/test-task-medods/internal/service/models"
)

const (
	jwtRefreshTokenDuration time.Duration = 24 * time.Hour
)

func (s *Service) GenerateTokens(ctx context.Context, userID, ip string) (*models.GeneratedTokens, error) {
	s.logger.Println("Generating jwt access and refresh tokens")

	user, err := s.getUserFromStorage(ctx, userID)
	if err != nil {
		return nil, err
	}

	base64RefreshToken, refreshTokenExpiresAt, err := s.createAndSaveRefreshToken(ctx, user.ID, ip)
	if err != nil {
		return nil, err
	}

	accessToken, accessTokenExpiresAt, err := s.createAccessToken(user.ID, ip)
	if err != nil {
		return nil, err
	}

	return &models.GeneratedTokens{
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  *accessTokenExpiresAt,
		RefreshToken:          base64RefreshToken,
		RefreshTokenExpiresAt: *refreshTokenExpiresAt,
	}, nil
}

func (s *Service) getUserFromStorage(ctx context.Context, userID string) (*models.User, error) {
	user, err := s.userStorage.GetUser(ctx, userID)
	if err != nil {
		if errors.Is(err, models.ErrUserNotFound) {
			s.logger.Printf("failed to find user: %v", err)

			return nil, models.ErrUserNotFound
		}

		s.logger.Printf("failed to find user: %v", err)

		return nil, fmt.Errorf("could not find user: %v", err)
	}

	return user, nil
}

func (s *Service) createAndSaveRefreshToken(ctx context.Context, userID uuid.UUID, ip string) (string, *time.Time, error) {
	refreshToken, refreshTokenClaims, err := s.tokenCreator.CreateToken(userID, ip, jwtRefreshTokenDuration)
	if err != nil {
		s.logger.Printf("error creating refresh token: %v", err)

		return "", nil, fmt.Errorf("could not create refresh token: %w", err)
	}

	refreshTokenHash := s.tokenCreator.CreateRefreshTokenHash(refreshToken)

	expiresAt := refreshTokenClaims.ExpiresAt.Time

	if err = s.sessionStorage.SaveSession(ctx, models.Session{ID: refreshTokenClaims.TokenID,
		UserIP:       refreshTokenClaims.UserIP,
		UserID:       refreshTokenClaims.UserID,
		RefreshToken: string(refreshTokenHash),
		ExpiredAt:    expiresAt,
	}); err != nil {
		s.logger.Printf("error saving session token: %v", err)

		return "", nil, fmt.Errorf("could not save refresh token: %w", err)
	}

	return base64.StdEncoding.EncodeToString([]byte(refreshToken)), &expiresAt, nil
}

func (s *Service) createAccessToken(userID uuid.UUID, ip string) (string, *time.Time, error) {
	accessToken, accessTokenClaims, err := s.tokenCreator.CreateToken(userID, ip, jwtAccessTokenDuration)
	if err != nil {
		s.logger.Printf("error creating access token: %v", err)

		return "", nil, fmt.Errorf("could not create access token: %w", err)
	}

	return accessToken, &accessTokenClaims.ExpiresAt.Time, nil
}
