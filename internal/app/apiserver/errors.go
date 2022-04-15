package apiserver

import "github.com/pkg/errors"

var (
	errWrongPathValue           = errors.New("Incorrect path value")
	errAlreadyRegistered        = errors.New("This user already exists")
	errIncorrectEmailOrPassword = errors.New("Incorrect email or password")
	errTokenExpiredOrDeleted    = errors.New("Token expired or has been deleted")
	errTokenDamaged             = errors.New("Token has been damaged")
	errSomethingWentWrong       = errors.New("Oops, something went wrong")
	errUnknownRepository        = errors.New("Something wrong with repository request")
)
