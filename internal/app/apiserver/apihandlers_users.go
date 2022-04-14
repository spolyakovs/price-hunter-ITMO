package apiserver

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/pkg/errors"

	"github.com/golang-jwt/jwt"
	"github.com/gorilla/mux"
	"github.com/spolyakovs/price-hunter-ITMO/internal/app/model"
	"github.com/spolyakovs/price-hunter-ITMO/internal/app/tokenUtils"
)

// TODO: add comments
// TODO: get app info 1 at a time for steam at least
// TODO: ONLY english letters in game name (debatable about description)
// TODO: restrincting number of games in games_list request (and add offset param)

func (server *server) handleRegistration() http.HandlerFunc {
	type request struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	return func(writer http.ResponseWriter, req *http.Request) {
		requestStruct := &request{}
		if err := json.NewDecoder(req.Body).Decode(requestStruct); err != nil {
			server.error(writer, req, http.StatusBadRequest, err)
			return
		}

		user := &model.User{
			Username: requestStruct.Username,
			Email:    requestStruct.Email,
			Password: requestStruct.Password,
		}

		if err := server.store.Users().Create(user); err != nil {
			fmt.Printf("DEBUG: %s\n", err.Error())
			server.error(writer, req, http.StatusBadRequest, errAlreadyRegistered)
			return
		}

		tokenDetails, tokenErr := tokenUtils.CreateTokens(user.ID)
		if tokenErr != nil {
			fmt.Printf("DEBUG: %s\n", tokenErr.Error())
			server.error(writer, req, http.StatusInternalServerError, tokenErr)
			return
		}

		tokens := map[string]string{
			"access_token":  tokenDetails.AccessToken,
			"refresh_token": tokenDetails.RefreshToken,
		}

		server.respond(writer, req, http.StatusCreated, tokens)
	}
}

func (server *server) handleLogin() http.HandlerFunc {
	type request struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	return func(writer http.ResponseWriter, req *http.Request) {
		requestStruct := &request{}
		if err := json.NewDecoder(req.Body).Decode(requestStruct); err != nil {
			server.error(writer, req, http.StatusBadRequest, err)
			return
		}

		user, userErr := server.store.Users().FindBy("username", requestStruct.Username)
		if userErr != nil || !user.ComparePassword(requestStruct.Password) {
			fmt.Printf("DEBUG: %s\n", userErr.Error())
			server.error(writer, req, http.StatusUnauthorized, errIncorrectEmailOrPassword)
			return
		}

		tokenDetails, tokenErr := tokenUtils.CreateTokens(user.ID)
		if tokenErr != nil {
			fmt.Printf("DEBUG: %s\n", tokenErr.Error())
			server.error(writer, req, http.StatusInternalServerError, tokenErr)
			return
		}

		tokens := map[string]string{
			"access_token":  tokenDetails.AccessToken,
			"refresh_token": tokenDetails.RefreshToken,
		}

		server.respond(writer, req, http.StatusOK, tokens)
	}
}

func (server *server) handleLogout() http.HandlerFunc {
	return func(writer http.ResponseWriter, req *http.Request) {
		tokenString, tokenExtractErr := tokenUtils.ExtractToken(req)
		if tokenExtractErr != nil {
			server.log(tokenExtractErr)
			switch errors.Cause(tokenExtractErr) {
			case tokenUtils.ErrTokenNotProvided:
				server.error(writer, req, http.StatusUnauthorized, tokenExtractErr)
			case tokenUtils.ErrTokenWrongFormat:
				server.error(writer, req, http.StatusBadRequest, tokenExtractErr)
			}
		}

		tokenData, tokenDataErr := tokenUtils.ExtractTokenMetadata(tokenString)
		if tokenDataErr != nil {
			validErr, ok := tokenDataErr.(*jwt.ValidationError)
			if ok && validErr.Errors == jwt.ValidationErrorExpired {
				server.respond(writer, req, http.StatusOK, nil)
				return
			}

			tokenDataErr = errors.Wrap(tokenDataErr, "Couldn't get token metadata")
			server.log(errors.WithMessage(tokenDataErr, errTokenDamaged.Error()))
			server.error(writer, req, http.StatusBadRequest, errTokenDamaged)
			return
		}

		if err := tokenUtils.DeleteAuth(tokenData.Uuid); err != nil {
			server.log(errors.Wrap(err, "Couldn't delete token"))
			server.error(writer, req, http.StatusInternalServerError, errSomethingWentWrong)
			return
		}

		server.respond(writer, req, http.StatusOK, nil)
	}
}

func (server *server) handleRefreshToken() http.HandlerFunc {
	type request struct {
		RefreshToken string `json:"refresh_token"`
	}

	return func(writer http.ResponseWriter, req *http.Request) {
		requestStruct := &request{}
		if err := json.NewDecoder(req.Body).Decode(requestStruct); err != nil {
			server.error(writer, req, http.StatusBadRequest, err)
			return
		}

		// Extract, validate and delete access token
		tokenString, tokenExtractErr := tokenUtils.ExtractToken(req)
		if tokenExtractErr != nil {
			server.log(tokenExtractErr)
			switch errors.Cause(tokenExtractErr) {
			case tokenUtils.ErrTokenNotProvided:
				server.error(writer, req, http.StatusUnauthorized, tokenExtractErr)
			case tokenUtils.ErrTokenWrongFormat:
				server.error(writer, req, http.StatusBadRequest, tokenExtractErr)
			}
		}

		authTokenData, authTokenDataErr := tokenUtils.ExtractTokenMetadata(tokenString)
		if authTokenDataErr != nil {
			validErr, ok := authTokenDataErr.(*jwt.ValidationError)
			if !ok || validErr.Errors != jwt.ValidationErrorExpired {
				authTokenDataErr = errors.Wrap(authTokenDataErr, "Couldn't get access token metadata")
				server.log(errors.WithMessage(authTokenDataErr, errTokenDamaged.Error()))
				server.error(writer, req, http.StatusBadRequest, errTokenDamaged)
			}
		}

		if err := tokenUtils.DeleteAuth(authTokenData.Uuid); err != nil {
			server.log(errors.Wrap(err, "Couldn't delete access token"))
			server.error(writer, req, http.StatusInternalServerError, errSomethingWentWrong)
			return
		}

		// Extract, validate and delete refresh token
		if err := tokenUtils.IsValid(requestStruct.RefreshToken); err != nil {
			server.log(errors.WithMessage(err, errTokenExpiredOrDeleted.Error()))
			server.error(writer, req, http.StatusForbidden, errTokenExpiredOrDeleted)
			return
		}

		refreshTokenData, refreshTokenDataErr := tokenUtils.ExtractTokenMetadata(requestStruct.RefreshToken)
		if refreshTokenDataErr != nil {
			refreshTokenDataErr = errors.Wrap(refreshTokenDataErr, "Couldn't get token metadata")
			server.log(errors.WithMessage(refreshTokenDataErr, errTokenDamaged.Error()))
			server.error(writer, req, http.StatusBadRequest, errTokenDamaged)
			return
		}

		if err := tokenUtils.DeleteAuth(refreshTokenData.Uuid); err != nil {
			server.log(errors.Wrap(err, "Couldn't delete token"))
			server.error(writer, req, http.StatusInternalServerError, errSomethingWentWrong)
			return
		}

		// Creating new pair of tokens
		tokenDetails, createErr := tokenUtils.CreateTokens(refreshTokenData.UserId)
		if createErr != nil {
			server.log(errors.Wrap(createErr, "Couldn't create token"))
			server.error(writer, req, http.StatusInternalServerError, errSomethingWentWrong)
			return
		}

		tokens := map[string]string{
			"access_token":  tokenDetails.AccessToken,
			"refresh_token": tokenDetails.RefreshToken,
		}

		server.respond(writer, req, http.StatusOK, tokens)
	}
}

func (server *server) handleUsersMe() http.HandlerFunc {
	return func(writer http.ResponseWriter, req *http.Request) {
		server.respond(writer, req, http.StatusOK, req.Context().Value(ctxKeyUser).(*model.User))
	}
}

func (server *server) handleUsersDeleteAllAuth() http.HandlerFunc {
	return func(writer http.ResponseWriter, req *http.Request) {
		if err := tokenUtils.DeleteAllAuths(req.Context().Value(ctxKeyUser).(*model.User).ID); err != nil {
			server.log(errors.Wrap(err, "Couldn't delete all tokens for user"))
			server.error(writer, req, http.StatusInternalServerError, errSomethingWentWrong)
		}
		server.respond(writer, req, http.StatusOK, nil)
	}
}

func (server *server) handleUsersChangeEmail() http.HandlerFunc {
	type request struct {
		NewEmail string `json:"new_email"`
	}

	return func(writer http.ResponseWriter, req *http.Request) {
		requestStruct := &request{}
		if err := json.NewDecoder(req.Body).Decode(requestStruct); err != nil {
			server.error(writer, req, http.StatusBadRequest, err)
			return
		}

		user := req.Context().Value(ctxKeyUser).(*model.User)

		if err := server.store.Users().UpdateEmail(requestStruct.NewEmail, user.ID); err != nil {
			server.error(writer, req, http.StatusInternalServerError, err)
			return
		}
		//TODO: delete old tokens
		server.respond(writer, req, http.StatusOK, nil)
	}
}

func (server *server) handleUsersChangePassword() http.HandlerFunc {
	type request struct {
		CurrentPassword string `json:"current_password"`
		NewPassword     string `json:"new_password"`
	}

	return func(writer http.ResponseWriter, req *http.Request) {
		requestStruct := &request{}
		if err := json.NewDecoder(req.Body).Decode(requestStruct); err != nil {
			server.error(writer, req, http.StatusBadRequest, err)
			return
		}

		user := req.Context().Value(ctxKeyUser).(*model.User)

		if user == nil {
			server.error(writer, req, http.StatusInternalServerError, errors.New("No user after authentication"))
			return
		}

		if !user.ComparePassword(requestStruct.CurrentPassword) {
			server.respond(writer, req, http.StatusOK, map[string]string{"error": "Wrong password"})
			return
		}

		if err := server.store.Users().UpdatePassword(requestStruct.NewPassword, user.ID); err != nil {
			server.error(writer, req, http.StatusInternalServerError, err)
			return
		}

		//TODO: delete old tokens
		server.respond(writer, req, http.StatusOK, user)
	}
}

func (server *server) handleUsersDelete() http.HandlerFunc {
	return func(writer http.ResponseWriter, req *http.Request) {
		user := req.Context().Value(ctxKeyUser).(*model.User)
		if err := server.store.Users().Delete(user.ID); err != nil {
			server.error(writer, req, http.StatusInternalServerError, err)
			return
		}

		//session, err := server.sessionStore.Get(req, sessionName)
		//if err != nil {
		//	server.error(writer, req, http.StatusInternalServerError, err)
		//	return
		//}
		//session.Values["user_id"] = nil

		server.respond(writer, req, http.StatusOK, nil)
	}
}

func (server *server) handleUsersGetByID() http.HandlerFunc {
	return func(writer http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)

		id, err := strconv.ParseUint(vars["id"], 10, 64)
		if err != nil {
			fmt.Printf("DEBUG: %s\n", err.Error())
			server.error(writer, req, http.StatusBadRequest, errWrongPathValue)
			return
		}

		user, err := server.store.Users().Find(id)
		if err != nil {
			fmt.Printf("DEBUG: %s\n", err.Error())
			server.error(writer, req, http.StatusNotFound, errSomethingWentWrong)
			return
		}

		server.respond(writer, req, http.StatusOK, user)
	}
}
