package store

import "github.com/pkg/errors"

var (
	ErrUnknownSQL   = errors.New("Something wrong with SQL request")
	ErrUserUsername = errors.New("User with this username already exists")
	ErrUserEmail    = errors.New("User with this email already exists")
	ErrNotFound     = errors.New("Record not found")
)

const (
	ErrStoreMessage = "%s repository %s error"
)
