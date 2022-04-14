package model

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID                uint64 `json:"id" db:"id,omitempty"`
	Username          string `json:"username" db:"username"`
	Email             string `json:"email" db:"email"`
	EncryptedPassword string `json:"-" db:"encrypted_password"`
	Password          string `json:"-"`
}

func (user *User) Validate() error {
	if err := validation.ValidateStruct(
		user,
		validation.Field(&user.Username, ValidationRulesUsername...),
		validation.Field(&user.Email, ValidationRulesEmail...),
		validation.Field(&user.Password, ValidationRulesPassword...),
	); err != nil {
		return errors.Wrap(ErrValidationFailed, err.Error())
	}

	return nil
}

func (user *User) BeforeCreate() error {
	if len(user.Password) > 0 {
		enc, err := encryptString(user.Password)
		if err != nil {
			return err
		}

		user.EncryptedPassword = enc
		user.Sanitize()
	}

	return nil
}

func (user *User) Sanitize() {
	user.Password = ""
}

func (user *User) ComparePassword(password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(user.EncryptedPassword), []byte(password)) == nil
}

func encryptString(s string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(s), bcrypt.MinCost)
	if err != nil {
		return "", errors.Wrap(ErrEncryptionFailed, err.Error())
	}

	return string(b), nil
}
