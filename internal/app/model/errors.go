package model

import "github.com/pkg/errors"

var (
	ErrValidationFailed = errors.New("Wrong data format")
	ErrEncryptionFailed = errors.New("Couldn't encrypt password")
)
