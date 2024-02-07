package domain

import "errors"

var (
	ErrUserNotFound        = errors.New("user with such credentials not found")
	ErrRefreshTokenExpired = errors.New("refresh token expired")
)
