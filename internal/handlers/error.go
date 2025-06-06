package handlers

import "errors"

var (
	ErrEmptyUserID  = errors.New("empty user_id")
	ErrTokenExpired = errors.New("token expired")
)
