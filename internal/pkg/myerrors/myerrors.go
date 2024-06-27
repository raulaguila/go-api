package myerrors

import "errors"

var (
	ErrDisabledUser       = errors.New("user is disabled")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserHasNoPhoto     = errors.New("user has no photo")
	ErrUserHasPass        = errors.New("user already has password")
)
