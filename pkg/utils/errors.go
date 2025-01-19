package utils

import "errors"

var (
	ErrDisabledUser        = errors.New("user is disabled")
	ErrInvalidCredentials  = errors.New("invalid credentials")
	ErrUserHasPass         = errors.New("user already has password")
	ErrPasswordsDoNotMatch = errors.New("passwords do not match")
	ErrInvalidID           = errors.New("invalid id")
)
