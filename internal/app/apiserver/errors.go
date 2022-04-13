package apiserver

import "github.com/pkg/errors"

var (
	errWrongPathValue           = errors.New("Incorrect path value")
	errAlreadyRegistered        = errors.New("This user already exists")
	errUserDoesNotExist         = errors.New("This user does not exist")
	errIncorrectEmailOrPassword = errors.New("Incorrect email or password")
	errNotAuthenticated         = errors.New("Not authenticated")
)
