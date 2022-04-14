package tokenUtils

import "github.com/pkg/errors"

var (
	ErrTokenCreate        = errors.New("Couldn't create token")
	ErrTokenSave          = errors.New("Couldn't save token to redis db")
	ErrTokenSigningMethod = errors.New("Unexpected signing method")
	ErrTokenWrongFormat   = errors.New("Wrong authentication header format")
	ErrTokenNotProvided   = errors.New("You need to authenticate")
	ErrTokenDoesNotExist  = errors.New("You need to authenticate")
	ErrTokenValidation    = errors.New("Something wrong with token validation")
	ErrTokenUUID          = errors.New("Something wrong with uuid")
	ErrTokenClaims        = errors.New("Something wrong with claims")
	ErrTokenParse         = errors.New("Couldn't parse token")
	ErrTokenDelete        = errors.New("Couldn't delete token")
	ErrUintParse          = errors.New("Couldn't parse uint")
	ErrRedis              = errors.New("Something wrong in redis")
)

const (
	ErrRedisNilMessage = "redis: nil"
)
