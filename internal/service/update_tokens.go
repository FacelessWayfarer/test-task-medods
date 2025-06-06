package service

import (
	"context"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/FacelessWayfarer/test-task-medods/internal/service/models"
	"github.com/FacelessWayfarer/test-task-medods/pkg/tokens"
)

const (
	jwtAccessTokenDuration time.Duration = 1 * time.Minute
)

func (s *Service) UpdateTokens(ctx context.Context, req models.TokensToRefresh, ip string) (*models.RefreshedTokens, error) {
	s.logger.Println("refreshing tokens")

	refreshClaims, session, err := s.verifyTokenAndFetchSession(ctx, req)
	if err != nil {
		return nil, err
	}

	if err = s.verifyExpiration(refreshClaims.ExpiresAt.Unix(), session.ExpiredAt.Unix()); err != nil {
		return nil, err
	}

	user, err := s.getUserFromStorage(ctx, fmt.Sprint(refreshClaims.UserID))
	if err != nil {
		return nil, err
	}

	if session.UserIP != ip {
		if err = s.sendEmail(user.Email); err != nil {
			return nil, err
		}
	}

	base64RefreshToken, refreshTokenExpiresAt, err := s.createAndSaveRefreshToken(ctx, user.ID, ip)
	if err != nil {
		return nil, err
	}

	accessToken, accessTokenExpiresAt, err := s.createAccessToken(user.ID, ip)
	if err != nil {
		return nil, err
	}

	return &models.RefreshedTokens{
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  *accessTokenExpiresAt,
		RefreshToken:          base64RefreshToken,
		RefreshTokenExpiresAt: *refreshTokenExpiresAt,
	}, nil
}

func (s *Service) verifyTokenAndFetchSession(ctx context.Context, req models.TokensToRefresh) (*tokens.UserClaims, *models.Session, error) {
	refreshToken, err := base64.StdEncoding.DecodeString(req.Base64RefreshToken)
	if err != nil {
		s.logger.Printf("failed to decode refresh token from request: %v", err)

		return nil, nil, fmt.Errorf("could not decode refresh token: %v", err)
	}

	refreshClaims, err := s.tokenCreator.VerifyToken(string(refreshToken))
	if err != nil {
		s.logger.Printf("failed to verify token: %v", err)

		return nil, nil, fmt.Errorf("could not verify token: %v", err)
	}

	accessClaims, err := s.tokenCreator.VerifyToken(req.AccessToken)
	if err != nil {
		s.logger.Printf("failed to verify token: %v", err)

		return nil, nil, fmt.Errorf("could not verify token: %v", err)
	}

	if refreshClaims.UserID != accessClaims.UserID {
		s.logger.Println("failed to verify token: invalid tokens")

		return nil, nil, models.ErrInvalidTokens
	}

	session, err := s.sessionStorage.GetSession(ctx, refreshClaims.TokenID)
	if err != nil {
		s.logger.Printf("failed to fetch session: %v", err)

		return nil, nil, fmt.Errorf("could not fetch session: %v", err)
	}

	return refreshClaims, session, nil
}

func (s *Service) sendEmail(userEmail string) error {
	s.logger.Printf("user ip changed, sending email to: %v", userEmail)

	return s.emailSender.SendEmail(userEmail)
}

func (s *Service) verifyExpiration(time1 int64, time2 int64) error {
	if time1 < time.Now().Unix() || time2 < time.Now().Unix() {
		s.logger.Println("token expired")

		return models.ErrTokenExpired
	}

	return nil
}
