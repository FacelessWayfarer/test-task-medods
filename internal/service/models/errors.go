package models

import "errors"

var (
	ErrUserNotFound    = errors.New("user not found")
	ErrSessionNotFound = errors.New("session not found")
	ErrInvalidTokens   = errors.New("invalid tokens")
	ErrTokenExpired    = errors.New("token expired")
)
