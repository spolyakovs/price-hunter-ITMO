package store

import "github.com/pkg/errors"

var (
	ErrUnknownSQL = errors.New("Something wrong with SQL request")
	ErrNotFound   = errors.New("Record not found")
	ErrCreate     = errors.New("Couldn't create DB record")
)
