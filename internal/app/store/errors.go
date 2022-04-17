package store

import "github.com/pkg/errors"

var (
	ErrUnknownSQL = errors.New("Something wrong with SQL request")
	ErrNotFound   = errors.New("Record not found")
)

const (
	ErrStoreMessageFormat = "%s repository %s error"
)
