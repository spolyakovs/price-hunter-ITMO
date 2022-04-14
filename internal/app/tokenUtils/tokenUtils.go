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

type TokensDetails struct {
	AccessToken  string
	RefreshToken string
	AccessUuid   string
	RefreshUuid  string
	AtExpires    int64
	RtExpires    int64
}

func CreateTokens(userid uint64) (*TokensDetails, error) {
	tokensDetails := &TokensDetails{}
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
		return nil, errors.Wrap(ErrTokenCreate, err.Error())
	}

	// Creating Refresh Token
	refreshTokenClaims := jwt.MapClaims{}
	refreshTokenClaims["uuid"] = tokensDetails.RefreshUuid
	refreshTokenClaims["user_id"] = userid
	refreshTokenClaims["exp"] = tokensDetails.RtExpires
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaims)
	tokensDetails.RefreshToken, err = refreshToken.SignedString([]byte(os.Getenv("TOKEN_SECRET")))
	if err != nil {
		return nil, errors.Wrap(ErrTokenCreate, err.Error())
	}

	accessTokenExpires := time.Unix(tokensDetails.AtExpires, 0) //converting Unix to UTC(to Time object)
	refreshTokenExpires := time.Unix(tokensDetails.RtExpires, 0)
	now := time.Now()

	accessErr := redisStore.Set(tokensDetails.AccessUuid, strconv.Itoa(int(userid)), accessTokenExpires.Sub(now)).Err()
	if accessErr != nil {
		return nil, errors.Wrap(ErrTokenSave, accessErr.Error())
	}
	refreshErr := redisStore.Set(tokensDetails.RefreshUuid, strconv.Itoa(int(userid)), refreshTokenExpires.Sub(now)).Err()
	if refreshErr != nil {
		return nil, errors.Wrap(ErrTokenSave, refreshErr.Error())
	}

	return tokensDetails, nil
}

func ExtractToken(r *http.Request) (string, error) {
	bearToken := r.Header.Get("Authorization")
	if bearToken == "" {
		return "", ErrTokenNotProvided
	}
	strArr := strings.Split(bearToken, " ")
	if len(strArr) == 2 {
		return strArr[1], nil
	}
	return "", ErrTokenWrongFormat
}

func verifyToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Make sure that the token method conform to "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrTokenSigningMethod
		}
		return []byte(os.Getenv("TOKEN_SECRET")), nil
	})
	if err != nil {
		return nil, errors.Wrap(ErrTokenParse, err.Error())
	}
	return token, nil
}

func IsValid(tokenString string) error {
	token, err := verifyToken(tokenString)
	if err != nil {
		return err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		uuid, ok := claims["uuid"].(string)
		if !ok {
			return errors.WithMessage(ErrTokenUUID, fmt.Sprintf("Claims: %+v", claims))
		}

		count, err := redisStore.Exists(uuid).Result()
		if err != nil {
			return errors.Wrap(ErrRedis, err.Error())
		}
		if count == 0 {
			return ErrTokenDoesNotExist
		}
	}

	return nil
}

type SingleTokenDetails struct {
	Uuid   string
	UserId uint64
}

func ExtractTokenMetadata(tokenString string) (*SingleTokenDetails, error) {
	token, err := verifyToken(tokenString)
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		uuid, ok := claims["uuid"].(string)
		if !ok {
			return nil, errors.WithMessage(ErrTokenUUID, fmt.Sprintf("Claims: %+v", claims))
		}

		userId, err := strconv.ParseUint(fmt.Sprintf("%.f", claims["user_id"]), 10, 64)
		if err != nil {
			return nil, errors.Wrap(ErrUintParse, err.Error())
		}

		return &SingleTokenDetails{
			Uuid:   uuid,
			UserId: userId,
		}, nil
	}
	return nil, errors.WithMessage(ErrTokenClaims, fmt.Sprintf("Claims: %+v", token.Claims))
}

// TODO: implement this into changeEmail and changePassword (don't forget about deleting tokens)

func FetchAuth(authDetails *SingleTokenDetails) (uint64, error) {
	userid, err := redisStore.Get(authDetails.Uuid).Result()
	if err != nil {
		if err.Error() == ErrRedisNilMessage {
			return 0, ErrTokenDoesNotExist
		}
		return 0, errors.Wrap(ErrRedis, err.Error())
	}
	userID, _ := strconv.ParseUint(userid, 10, 64)
	return userID, nil
}

func DeleteAuth(uuid string) error {
	err := redisStore.Del(uuid).Err()
	if err != nil {
		return errors.Wrap(ErrTokenDelete, err.Error())
	}
	return nil
}

func DeleteAllAuths(userid uint64) error {
	var cursor uint64
	for {
		var keys []string
		var err error
		keys, cursor, err = redisStore.Scan(cursor, "*", 0).Result()
		if err != nil {
			return errors.Wrap(ErrRedis, err.Error())
		}

		for _, key := range keys {
			result, err := redisStore.Get(key).Result()
			if err != nil {
				return errors.Wrap(ErrRedis, err.Error())
			}

			if result == strconv.Itoa(int(userid)) {
				if delErr := DeleteAuth(key); delErr != nil {
					return delErr
				}
			}
		}

		if cursor == 0 { // no more keys
			break
		}
	}
	return nil
}
