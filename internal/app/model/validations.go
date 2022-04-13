package model

import (
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/pkg/errors"
)

func ValidatePassword(value interface{}) error {
	err := validation.Validate(
		&value,
		validation.Required,
		validation.Length(8, 0),
		validation.Match(regexp.MustCompile(".*[A-Z].*")),
		validation.Match(regexp.MustCompile(".*[0-9].*")),
		validation.Match(regexp.MustCompile(".*[-!#$%&*+\\/=?^_{|}~].*")),
	)

	return errors.Wrap(ErrValidationFailed, err.Error())
}

func ValidateEmail(value interface{}) error {
	err := validation.Validate(
		&value,
		validation.Required,
		is.Email,
	)

	return errors.Wrap(ErrValidationFailed, err.Error())
}

func ValidateUsername(value interface{}) error {
	err := validation.Validate(
		&value,
		validation.Required,
		validation.Length(6, 30),
		validation.Match(regexp.MustCompile("^[a-zA-Z0-9_-]{6,30}$")),
	)

	return errors.Wrap(ErrValidationFailed, err.Error())
}
