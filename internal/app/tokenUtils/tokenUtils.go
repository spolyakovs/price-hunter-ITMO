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
	//Creating Access Token
	accessTokenClaims := jwt.MapClaims{}
	accessTokenClaims["authorized"] = true
	accessTokenClaims["access_uuid"] = tokenDetails.AccessUuid
	accessTokenClaims["user_id"] = userid
	accessTokenClaims["exp"] = tokenDetails.AtExpires
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims)
	tokenDetails.AccessToken, err = accessToken.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	if err != nil {
		return nil, err
	}

	//Creating Refresh Token
	refreshTokenClaims := jwt.MapClaims{}
	refreshTokenClaims["refresh_uuid"] = tokenDetails.RefreshUuid
	refreshTokenClaims["user_id"] = userid
	refreshTokenClaims["exp"] = tokenDetails.RtExpires
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaims)
	tokenDetails.RefreshToken, err = refreshToken.SignedString([]byte(os.Getenv("REFRESH_SECRET")))
	if err != nil {
		return nil, err
	}

	accessTokenExpires := time.Unix(tokenDetails.AtExpires, 0) //converting Unix to UTC(to Time object)
	refreshTokenExpires := time.Unix(tokenDetails.RtExpires, 0)
	now := time.Now()

	accessErr := redisStore.Set(tokenDetails.AccessUuid, strconv.Itoa(int(userid)), accessTokenExpires.Sub(now)).Err()
	if accessErr != nil {
		return nil, accessErr
	}
	refreshErr := redisStore.Set(tokenDetails.RefreshUuid, strconv.Itoa(int(userid)), refreshTokenExpires.Sub(now)).Err()
	if refreshErr != nil {
		return nil, refreshErr
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
		//Make sure that the token method conform to "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("ACCESS_SECRET")), nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}

func TokenValid(tokenString string) error {
	token, err := verifyToken(tokenString)
	if err != nil {
		return err
	}
	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		return fmt.Errorf("something wrong with token validation: %+v", token)
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
			return nil, fmt.Errorf("something wrong with access token uuid: %+v", claims)
		}
		userId, err := strconv.ParseUint(fmt.Sprintf("%.f", claims["user_id"]), 10, 64)
		if err != nil {
			return nil, err
		}
		return &AccessDetails{
			AccessUuid: accessUuid,
			UserId:     userId,
		}, nil
	}
	return nil, fmt.Errorf("something wrong with token claims: %+v", token.Claims)
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
			return nil, fmt.Errorf("something wrong with refresh token uuid: %+v", claims)
		}

		userId, err := strconv.ParseUint(fmt.Sprintf("%.f", claims["user_id"]), 10, 64)
		if err != nil {
			return nil, err
		}

		return &RefreshDetails{
			RefreshUuid: refreshUuid,
			UserId:      userId,
		}, nil
	}
	return nil, fmt.Errorf("something wrong with token claims: %+v", token.Claims)
}

// TODO: implement this into changeEmail and changePassword (don't forget about deleting tokens)

func FetchAuth(authD *AccessDetails) (uint64, error) {
	userid, err := redisStore.Get(authD.AccessUuid).Result()
	if err != nil {
		return 0, err
	}
	userID, _ := strconv.ParseUint(userid, 10, 64)
	return userID, nil
}

func DeleteAuth(uuid string) error {
	err := redisStore.Del(uuid).Err()
	if err != nil {
		return err
	}
	return nil
}
