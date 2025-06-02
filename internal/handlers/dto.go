package handlers

import (
	"time"

	"github.com/FacelessWayfarer/test-task-medods/internal/service/models"
)

type GenResponse struct {
	AccessToken           string    `json:"AccessToken,omitempty"`
	RefreshToken          string    `json:"RefreshToken,omitempty"`
	AccessTokenExpiresAt  time.Time `json:"AccessTokenExpiresAt,omitempty"`
	RefreshTokenExpiresAt time.Time `json:"RefreshTokenExpiresAt,omitempty"`
}

func ToGenResponse(tokens models.GeneratedTokens) GenResponse {
	return GenResponse{
		AccessToken:           tokens.AccessToken,
		RefreshToken:          tokens.RefreshToken,
		AccessTokenExpiresAt:  tokens.AccessTokenExpiresAt,
		RefreshTokenExpiresAt: tokens.RefreshTokenExpiresAt,
	}
}

type GenErrResponse struct {
	Error string `json:"Error,omitempty"`
}

type RefreshTokensRequest struct {
	AccessToken        string `json:"AccessToken"`
	Base64RefreshToken string `json:"Base64RefreshToken"`
}

type RefreshResponse struct {
	AccessToken           string    `json:"AccessToken,omitempty"`
	RefreshToken          string    `json:"RefreshToken,omitempty"`
	AccessTokenExpiresAt  time.Time `json:"AccessTokenExpiresAt,omitempty"`
	RefreshTokenExpiresAt time.Time `json:"RefreshTokenExpiresAt,omitempty"`
}

func ToRefreshResponse(tokens models.RefreshedTokens) RefreshResponse {
	return RefreshResponse{
		AccessToken:           tokens.AccessToken,
		RefreshToken:          tokens.RefreshToken,
		AccessTokenExpiresAt:  tokens.AccessTokenExpiresAt,
		RefreshTokenExpiresAt: tokens.RefreshTokenExpiresAt,
	}
}

type RefreshErrResponse struct {
	Error string `json:"Error,omitempty"`
}
