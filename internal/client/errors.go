package client

import "errors"

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrExpiredToken = errors.New("expired token")
	ErrServerError  = errors.New("server error")
)
