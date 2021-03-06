package tokenutils

import "github.com/pkg/errors"

var (
	ErrTokenWrongFormat      = errors.New("Wrong authentication header format")
	ErrTokenNotProvided      = errors.New("Authentication token wasn't provided")
	ErrTokenExpiredOrDeleted = errors.New("Authentication token expired or has been deleted")
	ErrTokenDamaged          = errors.New("Authentication token has been damaged")
	ErrInternal              = errors.New("Internal error")
)

const (
	errTokenUtilsMessageFormat   = "TokenUtils %s error"
	errTokenDeleteMessage        = "Couldn't delete token"
	errTokenCreateMessage        = "Couldn't create token"
	errTokenSaveMessage          = "Couldn't save token to redis db"
	errTokenValidationMessage    = "Something wrong with token validation"
	errTokenUUIDMessage          = "Something wrong with uuid"
	errTokenClaimsMessage        = "Something wrong with claims"
	errTokenParseMessage         = "Couldn't parse token"
	errTokenSigningMethodMessage = "Unexpected signing method"
	errTokenExpiredMessage       = "Token is expired"
	errUintParseMessage          = "Couldn't parse uint"
	errRedisMessage              = "Something wrong in redis"
	errRedisNilMessage           = "redis: nil"
)
