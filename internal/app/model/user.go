package model

import (
	validation "github.com/go-ozzo/ozzo-validation"
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
	return validation.ValidateStruct(
		user,
		validation.Field(&user.Username, validation.By(ValidateUsername)),
		validation.Field(&user.Email, validation.By(ValidateEmail)),
		validation.Field(&user.Password, validation.By(ValidatePassword)),
	)
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
		return "", err
	}

	return string(b), nil
}
