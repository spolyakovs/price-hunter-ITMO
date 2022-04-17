package apiserver

import "github.com/pkg/errors"

var (
	errWrongRequestFormat = errors.New("Wrong request format")
	errSomethingWentWrong = errors.New("Oops, something went wrong")
)

const (
	errMiddlewareMessageFormat        = "Request middleware %s error"
	errHandlerMessageFormat           = "Request handler %s error"
	errUserExistsUsernameMessage      = "User with this username already exists"
	errUserExistsEmailMessage         = "User with this email already exists"
	errWrongUsernameOrPasswordMessage = "Wrong username or password"
	errWrongPasswordMessage           = "Wrong password"
)
