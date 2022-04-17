package tokenUtils

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type TokenPairDetails struct {
	AccessToken  string
	RefreshToken string
	AccessUuid   string
	RefreshUuid  string
	AtExpires    int64
	RtExpires    int64
}

// TODO: refactor errors like here (errWrapped)

func CreateTokens(userid uint64) (*TokenPairDetails, error) {
	methodName := "CreateTokens"
	errWrapMessage := fmt.Sprintf(errTokenUtilsMessageFormat, methodName)

	tokensDetails := &TokenPairDetails{}
	tokensDetails.AtExpires = time.Now().Add(time.Minute * 15).Unix()
	tokensDetails.AccessUuid = uuid.New().String()

	tokensDetails.RtExpires = time.Now().Add(time.Hour * 24 * 7).Unix()
	tokensDetails.RefreshUuid = uuid.New().String()

	var err error
	// Creating Access Token
	accessTokenClaims := jwt.MapClaims{}
	accessTokenClaims["authorized"] = true
	accessTokenClaims["uuid"] = tokensDetails.AccessUuid
	accessTokenClaims["user_id"] = userid
	accessTokenClaims["exp"] = tokensDetails.AtExpires
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims)
	tokensDetails.AccessToken, err = accessToken.SignedString([]byte(os.Getenv("TOKEN_SECRET")))
	if err != nil {
		errWrapped := errors.Wrap(ErrInternal, errTokenCreateMessage)
		errWrapped = errors.WithMessage(errWrapped, err.Error())
		errWrapped = errors.Wrap(errWrapped, errWrapMessage)
		return nil, errWrapped
	}

	// Creating Refresh Token
	refreshTokenClaims := jwt.MapClaims{}
	refreshTokenClaims["uuid"] = tokensDetails.RefreshUuid
	refreshTokenClaims["user_id"] = userid
	refreshTokenClaims["exp"] = tokensDetails.RtExpires
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaims)
	tokensDetails.RefreshToken, err = refreshToken.SignedString([]byte(os.Getenv("TOKEN_SECRET")))
	if err != nil {
		errWrapped := errors.Wrap(ErrInternal, errTokenCreateMessage)
		errWrapped = errors.WithMessage(errWrapped, err.Error())
		errWrapped = errors.Wrap(errWrapped, errWrapMessage)
		return nil, errWrapped
	}
	// Converting Unix to UTC(to Time object)
	accessTokenExpires := time.Unix(tokensDetails.AtExpires, 0)
	refreshTokenExpires := time.Unix(tokensDetails.RtExpires, 0)
	now := time.Now()

	// Saving Access token
	if err := redisStore.Set(tokensDetails.AccessUuid, strconv.Itoa(int(userid)), accessTokenExpires.Sub(now)).Err(); err != nil {
		errWrapped := errors.Wrap(ErrInternal, errTokenSaveMessage)
		errWrapped = errors.WithMessage(errWrapped, err.Error())
		errWrapped = errors.Wrap(errWrapped, errWrapMessage)
		return nil, errWrapped
	}

	// Saving Refresh token
	if err := redisStore.Set(tokensDetails.RefreshUuid, strconv.Itoa(int(userid)), refreshTokenExpires.Sub(now)).Err(); err != nil {
		errWrapped := errors.Wrap(ErrInternal, errTokenSaveMessage)
		errWrapped = errors.WithMessage(errWrapped, err.Error())
		errWrapped = errors.Wrap(errWrapped, errWrapMessage)
		return nil, errWrapped
	}

	return tokensDetails, nil
}

func ExtractToken(r *http.Request) (string, error) {
	methodName := "ExtractToken"
	errWrapMessage := fmt.Sprintf(errTokenUtilsMessageFormat, methodName)

	bearToken := r.Header.Get("Authorization")
	if bearToken == "" {
		errWrapped := errors.Wrap(ErrTokenNotProvided, errWrapMessage)
		return "", errWrapped
	}
	strArr := strings.Split(bearToken, " ")
	if len(strArr) != 2 {
		errWrapped := errors.Wrap(ErrTokenWrongFormat, errWrapMessage)
		return "", errWrapped
	}

	return strArr[1], nil
}

func verifyToken(tokenString string) (*jwt.Token, error) {
	methodName := "VerifyToken"
	errWrapMessage := fmt.Sprintf(errTokenUtilsMessageFormat, methodName)

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Make sure that the token method conform to "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			errWrapped := errors.Wrap(ErrTokenDamaged, errTokenSigningMethodMessage)
			errWrapped = errors.Wrap(errWrapped, errWrapMessage)
			return nil, errWrapped
		}
		return []byte(os.Getenv("TOKEN_SECRET")), nil
	})
	if err != nil {
		if err.Error() == errTokenExpiredMessage {
			errWrapped := errors.Wrap(ErrTokenExpiredOrDeleted, errWrapMessage)
			return nil, errWrapped
		}
		errWrapped := errors.Wrap(ErrTokenDamaged, errTokenParseMessage)
		errWrapped = errors.WithMessage(errWrapped, err.Error())
		errWrapped = errors.Wrap(errWrapped, errWrapMessage)
		return nil, errWrapped
	}
	return token, nil
}

// Probably unnecessary cause it doubles part of ExtractTokenMetadata
func IsValid(tokenString string) error {
	methodName := "IsValid"
	errWrapMessage := fmt.Sprintf(errTokenUtilsMessageFormat, methodName)

	token, err := verifyToken(tokenString)
	if err != nil {
		errWrapped := errors.Wrap(err, errWrapMessage)
		return errWrapped
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		errWrapped := errors.Wrap(ErrTokenDamaged, errTokenClaimsMessage)
		errWrapped = errors.WithMessage(errWrapped, fmt.Sprintf("Claims: %+v", token.Claims))
		errWrapped = errors.Wrap(errWrapped, errWrapMessage)
		return errWrapped
	}

	if !token.Valid {
		errWrapped := errors.Wrap(ErrTokenDamaged, errTokenValidationMessage)
		errWrapped = errors.Wrap(errWrapped, errWrapMessage)
		return errWrapped
	}

	uuid, ok := claims["uuid"].(string)
	if !ok {
		errWrapped := errors.Wrap(ErrTokenDamaged, errTokenUUIDMessage)
		errWrapped = errors.WithMessage(errWrapped, fmt.Sprintf("Claims: %+v", claims))
		errWrapped = errors.Wrap(errWrapped, errWrapMessage)
		return errWrapped
	}

	count, err := redisStore.Exists(uuid).Result()
	if err != nil {
		errWrapped := errors.Wrap(ErrInternal, errRedisMessage)
		errWrapped = errors.WithMessage(errWrapped, err.Error())
		errWrapped = errors.Wrap(errWrapped, errWrapMessage)
		return errWrapped
	}

	if count == 0 {
		errWrapped := errors.Wrap(ErrTokenExpiredOrDeleted, errWrapMessage)
		return errWrapped
	}

	return nil
}

type TokenDetails struct {
	Uuid   string
	UserId uint64
}

func ExtractTokenMetadata(tokenString string) (*TokenDetails, error) {
	methodName := "ExtractTokenMetadata"
	errWrapMessage := fmt.Sprintf(errTokenUtilsMessageFormat, methodName)

	token, err := verifyToken(tokenString)
	if err != nil {
		errWrapped := errors.Wrap(err, errWrapMessage)
		return nil, errWrapped
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		errWrapped := errors.Wrap(ErrTokenDamaged, errTokenClaimsMessage)
		errWrapped = errors.WithMessage(errWrapped, fmt.Sprintf("Claims: %+v", token.Claims))
		errWrapped = errors.Wrap(errWrapped, errWrapMessage)
		return nil, errWrapped
	}

	if !token.Valid {
		errWrapped := errors.Wrap(ErrTokenDamaged, errTokenValidationMessage)
		errWrapped = errors.Wrap(errWrapped, errWrapMessage)
		return nil, errWrapped
	}

	uuid, ok := claims["uuid"].(string)
	if !ok {
		errWrapped := errors.Wrap(ErrTokenDamaged, errTokenUUIDMessage)
		errWrapped = errors.WithMessage(errWrapped, fmt.Sprintf("Claims: %+v", claims))
		errWrapped = errors.Wrap(errWrapped, errWrapMessage)
		return nil, errWrapped
	}

	count, err := redisStore.Exists(uuid).Result()
	if err != nil {
		errWrapped := errors.Wrap(ErrInternal, errRedisMessage)
		errWrapped = errors.WithMessage(errWrapped, err.Error())
		errWrapped = errors.Wrap(errWrapped, errWrapMessage)
		return nil, errWrapped
	}

	if count == 0 {
		errWrapped := errors.Wrap(ErrTokenExpiredOrDeleted, errWrapMessage)
		return nil, errWrapped
	}

	userId, err := strconv.ParseUint(fmt.Sprintf("%.f", claims["user_id"]), 10, 64)
	if err != nil {
		errWrapped := errors.Wrap(ErrTokenDamaged, errUintParseMessage)
		errWrapped = errors.WithMessage(errWrapped, err.Error())
		errWrapped = errors.Wrap(errWrapped, errWrapMessage)
		return nil, errWrapped
	}

	return &TokenDetails{
		Uuid:   uuid,
		UserId: userId,
	}, nil
}

func FetchAuth(authDetails *TokenDetails) (uint64, error) {
	methodName := "FetchAuth"
	errWrapMessage := fmt.Sprintf(errTokenUtilsMessageFormat, methodName)

	userIDRaw, err := redisStore.Get(authDetails.Uuid).Result()
	if err != nil {
		if err.Error() == errRedisNilMessage {
			errWrapped := errors.Wrap(ErrTokenExpiredOrDeleted, errWrapMessage)
			return 0, errWrapped
		}
		errWrapped := errors.Wrap(ErrInternal, errRedisMessage)
		errWrapped = errors.WithMessage(errWrapped, err.Error())
		errWrapped = errors.Wrap(errWrapped, errWrapMessage)
		return 0, errWrapped
	}
	userID, _ := strconv.ParseUint(userIDRaw, 10, 64)
	return userID, nil
}

func DeleteAuth(uuid string) error {
	methodName := "DeleteAuth"
	errWrapMessage := fmt.Sprintf(errTokenUtilsMessageFormat, methodName)

	err := redisStore.Del(uuid).Err()
	if err != nil {
		errWrapped := errors.Wrap(ErrInternal, errTokenDeleteMessage)
		errWrapped = errors.WithMessage(errWrapped, err.Error())
		errWrapped = errors.Wrap(errWrapped, errWrapMessage)
		return errWrapped
	}
	return nil
}

func DeleteAllAuths(userid uint64) error {
	methodName := "DeleteAllAuths"
	errWrapMessage := fmt.Sprintf(errTokenUtilsMessageFormat, methodName)

	var cursor uint64
	for {
		var keys []string
		var err error
		keys, cursor, err = redisStore.Scan(cursor, "*", 0).Result()
		if err != nil {
			errWrapped := errors.Wrap(ErrInternal, errRedisMessage)
			errWrapped = errors.WithMessage(errWrapped, err.Error())
			errWrapped = errors.Wrap(errWrapped, errWrapMessage)
			return errWrapped
		}

		for _, key := range keys {
			result, err := redisStore.Get(key).Result()
			if err != nil {
				errWrapped := errors.Wrap(ErrInternal, errRedisMessage)
				errWrapped = errors.WithMessage(errWrapped, err.Error())
				errWrapped = errors.Wrap(errWrapped, errWrapMessage)
				return errWrapped
			}

			if result == strconv.Itoa(int(userid)) {
				if delErr := DeleteAuth(key); delErr != nil {
					errWrapped := errors.Wrap(delErr, errWrapMessage)
					return errWrapped
				}
			}
		}

		if cursor == 0 { // no more keys
			break
		}
	}
	return nil
}
