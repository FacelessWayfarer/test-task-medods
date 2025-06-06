package handlers

import (
	"time"

	"github.com/FacelessWayfarer/test-task-medods/internal/service/models"
)

type GenResponse struct {
	AccessToken           string     `json:"access_token,omitempty"`
	RefreshToken          string     `json:"refresh_token,omitempty"`
	AccessTokenExpiresAt  *time.Time `json:"access_token_expires_at,omitempty"`
	RefreshTokenExpiresAt *time.Time `json:"refresh_token_expires_at,omitempty"`
}

func ToGenResponse(tokens models.GeneratedTokens) GenResponse {
	return GenResponse{
		AccessToken:           tokens.AccessToken,
		RefreshToken:          tokens.RefreshToken,
		AccessTokenExpiresAt:  &tokens.AccessTokenExpiresAt,
		RefreshTokenExpiresAt: &tokens.RefreshTokenExpiresAt,
	}
}

type GenErrResponse struct {
	Error string `json:"Error"`
}

type RefreshTokensRequest struct {
	AccessToken        string `json:"access_token"`
	Base64RefreshToken string `json:"base_64_refresh_token"`
}

func RequestToRequest(request RefreshTokensRequest) models.TokensToRefresh {
	return models.TokensToRefresh{
		AccessToken:        request.AccessToken,
		Base64RefreshToken: request.Base64RefreshToken,
	}
}

type RefreshResponse struct {
	AccessToken           string     `json:"access_token,omitempty"`
	RefreshToken          string     `json:"refresh_token,omitempty"`
	AccessTokenExpiresAt  *time.Time `json:"access_token_expires_at,omitempty"`
	RefreshTokenExpiresAt *time.Time `json:"refresh_token_expires_at,omitempty"`
}

func ToRefreshResponse(tokens models.RefreshedTokens) RefreshResponse {
	return RefreshResponse{
		AccessToken:           tokens.AccessToken,
		RefreshToken:          tokens.RefreshToken,
		AccessTokenExpiresAt:  &tokens.AccessTokenExpiresAt,
		RefreshTokenExpiresAt: &tokens.RefreshTokenExpiresAt,
	}
}

type RefreshErrResponse struct {
	Error string `json:"Error"`
}
