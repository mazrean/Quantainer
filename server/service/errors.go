package service

import "errors"

var (
	ErrInvalidAuthStateOrCode = errors.New("invalid auth state or code")
	ErrInvaliClientID         = errors.New("invalid client id")
	ErrOIDCSessionExpired     = errors.New("session access token expired")
	ErrForbidden              = errors.New("forbidden")
	ErrInvalidFormat          = errors.New("invalid format")
	ErrNoFile                 = errors.New("no file")
	ErrNoResource             = errors.New("no resource")
)
