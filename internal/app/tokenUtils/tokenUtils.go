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

func CreateTokens(userid uint64) (*TokenPairDetails, error) {
	methodName := "CreateTokens"
	errorMethodMessage := fmt.Sprintf(errTokenUtilsMessage, methodName)

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
		return nil, errors.Wrap(errors.WithMessage(errors.Wrap(ErrInternal, errTokenCreateMessage), err.Error()), errorMethodMessage)
	}

	// Creating Refresh Token
	refreshTokenClaims := jwt.MapClaims{}
	refreshTokenClaims["uuid"] = tokensDetails.RefreshUuid
	refreshTokenClaims["user_id"] = userid
	refreshTokenClaims["exp"] = tokensDetails.RtExpires
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaims)
	tokensDetails.RefreshToken, err = refreshToken.SignedString([]byte(os.Getenv("TOKEN_SECRET")))
	if err != nil {
		return nil, errors.Wrap(errors.WithMessage(errors.Wrap(ErrInternal, errTokenCreateMessage), err.Error()), errorMethodMessage)
	}

	accessTokenExpires := time.Unix(tokensDetails.AtExpires, 0) //converting Unix to UTC(to Time object)
	refreshTokenExpires := time.Unix(tokensDetails.RtExpires, 0)
	now := time.Now()

	accessErr := redisStore.Set(tokensDetails.AccessUuid, strconv.Itoa(int(userid)), accessTokenExpires.Sub(now)).Err()
	if accessErr != nil {
		return nil, errors.Wrap(errors.WithMessage(errors.Wrap(ErrInternal, errTokenSaveMessage), accessErr.Error()), errorMethodMessage)
	}
	refreshErr := redisStore.Set(tokensDetails.RefreshUuid, strconv.Itoa(int(userid)), refreshTokenExpires.Sub(now)).Err()
	if refreshErr != nil {
		return nil, errors.Wrap(errors.WithMessage(errors.Wrap(ErrInternal, errTokenSaveMessage), refreshErr.Error()), errorMethodMessage)
	}

	return tokensDetails, nil
}

func ExtractToken(r *http.Request) (string, error) {
	methodName := "ExtractToken"
	errorMethodMessage := fmt.Sprintf(errTokenUtilsMessage, methodName)

	bearToken := r.Header.Get("Authorization")
	if bearToken == "" {
		return "", errors.Wrap(ErrTokenNotProvided, errorMethodMessage)
	}
	strArr := strings.Split(bearToken, " ")
	if len(strArr) == 2 {
		return strArr[1], nil
	}
	return "", errors.Wrap(ErrTokenWrongFormat, errorMethodMessage)
}

func verifyToken(tokenString string) (*jwt.Token, error) {
	methodName := "VerifyToken"
	errorMethodMessage := fmt.Sprintf(errTokenUtilsMessage, methodName)

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Make sure that the token method conform to "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.Wrap(errors.Wrap(ErrTokenDamaged, errTokenSigningMethodMessage), errorMethodMessage)
		}
		return []byte(os.Getenv("TOKEN_SECRET")), nil
	})
	if err != nil {
		if err.Error() == errTokenExpiredMessage {
			return nil, errors.Wrap(ErrTokenExpiredOrDeleted, errorMethodMessage)
		}
		return nil, errors.Wrap(errors.Wrap(errors.Wrap(ErrTokenDamaged, errTokenParseMessage), err.Error()), errorMethodMessage)
	}
	return token, nil
}

func IsValid(tokenString string) error {
	methodName := "IsValid"
	errorMethodMessage := fmt.Sprintf(errTokenUtilsMessage, methodName)

	token, err := verifyToken(tokenString)
	if err != nil {
		return errors.Wrap(err, errorMethodMessage)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return errors.Wrap(errors.WithMessage(errors.Wrap(ErrTokenDamaged, errTokenClaimsMessage), fmt.Sprintf("Claims: %+v", token.Claims)), errorMethodMessage)
	}

	if !token.Valid {
		return errors.Wrap(errors.Wrap(ErrTokenDamaged, errTokenValidationMessage), errorMethodMessage)
	}

	uuid, ok := claims["uuid"].(string)
	if !ok {
		return errors.Wrap(errors.WithMessage(errors.Wrap(ErrTokenDamaged, errTokenUUIDMessage), fmt.Sprintf("Claims: %+v", claims)), errorMethodMessage)
	}

	count, err := redisStore.Exists(uuid).Result()
	if err != nil {
		return errors.Wrap(errors.WithMessage(errors.Wrap(ErrInternal, errRedisMessage), err.Error()), errorMethodMessage)
	}

	if count == 0 {
		return errors.Wrap(ErrTokenExpiredOrDeleted, errorMethodMessage)
	}

	return nil
}

type TokenDetails struct {
	Uuid   string
	UserId uint64
}

func ExtractTokenMetadata(tokenString string) (*TokenDetails, error) {
	methodName := "ExtractTokenMetadata"
	errorMethodMessage := fmt.Sprintf(errTokenUtilsMessage, methodName)

	token, err := verifyToken(tokenString)
	if err != nil {
		return nil, errors.Wrap(err, errorMethodMessage)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.Wrap(errors.WithMessage(errors.Wrap(ErrTokenDamaged, errTokenClaimsMessage), fmt.Sprintf("Claims: %+v", token.Claims)), errorMethodMessage)
	}

	if !token.Valid {
		return nil, errors.Wrap(errors.Wrap(ErrTokenDamaged, errTokenValidationMessage), errorMethodMessage)
	}

	uuid, ok := claims["uuid"].(string)
	if !ok {
		return nil, errors.Wrap(errors.WithMessage(errors.Wrap(ErrTokenDamaged, errTokenUUIDMessage), fmt.Sprintf("Claims: %+v", claims)), errorMethodMessage)
	}

	userId, err := strconv.ParseUint(fmt.Sprintf("%.f", claims["user_id"]), 10, 64)
	if err != nil {
		return nil, errors.Wrap(errors.WithMessage(errors.Wrap(ErrTokenDamaged, errUintParseMessage), err.Error()), errorMethodMessage)
	}

	return &TokenDetails{
		Uuid:   uuid,
		UserId: userId,
	}, nil
}

// TODO: implement this into changeEmail and changePassword (don't forget about deleting tokens)

func FetchAuth(authDetails *TokenDetails) (uint64, error) {
	methodName := "FetchAuth"
	errorMethodMessage := fmt.Sprintf(errTokenUtilsMessage, methodName)

	userid, err := redisStore.Get(authDetails.Uuid).Result()
	if err != nil {
		if err.Error() == errRedisNilMessage {
			return 0, errors.Wrap(ErrTokenExpiredOrDeleted, errorMethodMessage)
		}
		return 0, errors.Wrap(errors.WithMessage(errors.Wrap(ErrInternal, errRedisMessage), err.Error()), errorMethodMessage)
	}
	userID, _ := strconv.ParseUint(userid, 10, 64)
	return userID, nil
}

func DeleteAuth(uuid string) error {
	methodName := "DeleteAuth"
	errorMethodMessage := fmt.Sprintf(errTokenUtilsMessage, methodName)

	err := redisStore.Del(uuid).Err()
	if err != nil {
		return errors.Wrap(errors.WithMessage(errors.Wrap(ErrInternal, errTokenDeleteMessage), err.Error()), errorMethodMessage)
	}
	return nil
}

func DeleteAllAuths(userid uint64) error {
	methodName := "DeleteAllAuths"
	errorMethodMessage := fmt.Sprintf(errTokenUtilsMessage, methodName)

	var cursor uint64
	for {
		var keys []string
		var err error
		keys, cursor, err = redisStore.Scan(cursor, "*", 0).Result()
		if err != nil {
			return errors.Wrap(errors.WithMessage(errors.Wrap(ErrInternal, errRedisMessage), err.Error()), errorMethodMessage)
		}

		for _, key := range keys {
			result, err := redisStore.Get(key).Result()
			if err != nil {
				return errors.Wrap(errors.WithMessage(errors.Wrap(ErrInternal, errRedisMessage), err.Error()), errorMethodMessage)
			}

			if result == strconv.Itoa(int(userid)) {
				if delErr := DeleteAuth(key); delErr != nil {
					return errors.Wrap(delErr, errorMethodMessage)
				}
			}
		}

		if cursor == 0 { // no more keys
			break
		}
	}
	return nil
}
