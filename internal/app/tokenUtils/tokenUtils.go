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

type TokenDetails struct {
	AccessToken  string
	RefreshToken string
	AccessUuid   string
	RefreshUuid  string
	AtExpires    int64
	RtExpires    int64
}

func CreateToken(userid uint64) (*TokenDetails, error) {
	tokenDetails := &TokenDetails{}
	tokenDetails.AtExpires = time.Now().Add(time.Minute * 15).Unix()
	tokenDetails.AccessUuid = uuid.New().String()

	tokenDetails.RtExpires = time.Now().Add(time.Hour * 24 * 7).Unix()
	tokenDetails.RefreshUuid = uuid.New().String()

	var err error
	// Creating Access Token
	accessTokenClaims := jwt.MapClaims{}
	accessTokenClaims["authorized"] = true
	accessTokenClaims["access_uuid"] = tokenDetails.AccessUuid
	accessTokenClaims["user_id"] = userid
	accessTokenClaims["exp"] = tokenDetails.AtExpires
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims)
	tokenDetails.AccessToken, err = accessToken.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	if err != nil {
		return nil, errors.Wrap(ErrTokenCreate, err.Error())
	}

	// Creating Refresh Token
	refreshTokenClaims := jwt.MapClaims{}
	refreshTokenClaims["refresh_uuid"] = tokenDetails.RefreshUuid
	refreshTokenClaims["user_id"] = userid
	refreshTokenClaims["exp"] = tokenDetails.RtExpires
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaims)
	tokenDetails.RefreshToken, err = refreshToken.SignedString([]byte(os.Getenv("REFRESH_SECRET")))
	if err != nil {
		return nil, errors.Wrap(ErrTokenCreate, err.Error())
	}

	accessTokenExpires := time.Unix(tokenDetails.AtExpires, 0) //converting Unix to UTC(to Time object)
	refreshTokenExpires := time.Unix(tokenDetails.RtExpires, 0)
	now := time.Now()

	accessErr := redisStore.Set(tokenDetails.AccessUuid, strconv.Itoa(int(userid)), accessTokenExpires.Sub(now)).Err()
	if accessErr != nil {
		return nil, errors.Wrap(ErrTokenSave, accessErr.Error())
	}
	refreshErr := redisStore.Set(tokenDetails.RefreshUuid, strconv.Itoa(int(userid)), refreshTokenExpires.Sub(now)).Err()
	if refreshErr != nil {
		return nil, errors.Wrap(ErrTokenSave, refreshErr.Error())
	}

	return tokenDetails, nil
}

func ExtractToken(r *http.Request) string {
	bearToken := r.Header.Get("Authorization")
	strArr := strings.Split(bearToken, " ")
	if len(strArr) == 2 {
		return strArr[1]
	}
	return ""
}

func verifyToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Make sure that the token method conform to "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrTokenSigningMethod
		}
		return []byte(os.Getenv("ACCESS_SECRET")), nil
	})
	if err != nil {
		return nil, errors.Wrap(ErrTokenParse, err.Error())
	}
	return token, nil
}

func TokenValid(tokenString string) error {
	token, err := verifyToken(tokenString)
	if err != nil {
		return err
	}

	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		return ErrTokenValidation
	}

	return nil
}

type AccessDetails struct {
	AccessUuid string
	UserId     uint64
}

type RefreshDetails struct {
	RefreshUuid string
	UserId      uint64
}

func ExtractAccessTokenMetadata(tokenString string) (*AccessDetails, error) {
	token, err := verifyToken(tokenString)
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		accessUuid, ok := claims["access_uuid"].(string)
		if !ok {
			return nil, errors.WithMessage(ErrTokenUUID, fmt.Sprintf("Claims: %+v", claims))
		}
		userId, err := strconv.ParseUint(fmt.Sprintf("%.f", claims["user_id"]), 10, 64)
		if err != nil {
			return nil, errors.Wrap(ErrUintParse, err.Error())
		}
		return &AccessDetails{
			AccessUuid: accessUuid,
			UserId:     userId,
		}, nil
	}
	return nil, errors.WithMessage(ErrTokenClaims, fmt.Sprintf("Claims: %+v", token.Claims))
}

func ExtractRefreshTokenMetadata(tokenString string) (*RefreshDetails, error) {
	token, err := verifyToken(tokenString)
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		refreshUuid, ok := claims["refresh_uuid"].(string)
		if !ok {
			return nil, errors.WithMessage(ErrTokenUUID, fmt.Sprintf("Claims: %+v", claims))
		}

		userId, err := strconv.ParseUint(fmt.Sprintf("%.f", claims["user_id"]), 10, 64)
		if err != nil {
			return nil, errors.Wrap(ErrUintParse, err.Error())
		}

		return &RefreshDetails{
			RefreshUuid: refreshUuid,
			UserId:      userId,
		}, nil
	}
	return nil, errors.WithMessage(ErrTokenClaims, fmt.Sprintf("Claims: %+v", token.Claims))
}

// TODO: implement this into changeEmail and changePassword (don't forget about deleting tokens)

func FetchAuth(authDetails *AccessDetails) (uint64, error) {
	userid, err := redisStore.Get(authDetails.AccessUuid).Result()
	if err != nil {
		return 0, err
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
