package store

import "github.com/pkg/errors"

var (
	ErrUnknownSQL = errors.New("Something wrong with SQL request")
	ErrNotFound   = errors.New("Record not found")
)

const (
	ErrRepositoryMessageFormat = "%s repository %s error"
	ErrCreateTablesMessageFormat = "Creating %s table error"
)
