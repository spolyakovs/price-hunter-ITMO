package apiserver

import (
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"

	"github.com/golang-jwt/jwt"
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
			server.log(err)
			server.error(writer, req, http.StatusBadRequest, errors.Cause(err))
			return
		}

		tokenDetails, tokenDetailsErr := tokenUtils.CreateTokens(user.ID)
		if tokenDetailsErr != nil {
			server.log(tokenDetailsErr)
			server.error(writer, req, http.StatusInternalServerError, tokenDetailsErr)
			return
		}

		tokens := map[string]string{
			"access_token":  tokenDetails.AccessToken,
			"refresh_token": tokenDetails.RefreshToken,
		}

		server.respond(writer, req, http.StatusCreated, tokens)
		return
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
			server.log(userErr)
			server.error(writer, req, http.StatusUnauthorized, errIncorrectEmailOrPassword)
			return
		}

		tokenDetails, tokenDetailsErr := tokenUtils.CreateTokens(user.ID)
		if tokenDetailsErr != nil {
			server.log(tokenDetailsErr)
			server.error(writer, req, http.StatusInternalServerError, tokenDetailsErr)
			return
		}

		tokens := map[string]string{
			"access_token":  tokenDetails.AccessToken,
			"refresh_token": tokenDetails.RefreshToken,
		}

		server.respond(writer, req, http.StatusOK, tokens)
		return
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
				return
			case tokenUtils.ErrTokenWrongFormat:
				server.error(writer, req, http.StatusBadRequest, tokenExtractErr)
				return
			}
		}

		tokenDetails, tokenDetailsErr := tokenUtils.ExtractTokenMetadata(tokenString)
		if tokenDetailsErr != nil {
			validErr, ok := tokenDetailsErr.(*jwt.ValidationError)
			if ok && validErr.Errors == jwt.ValidationErrorExpired {
				server.respond(writer, req, http.StatusOK, nil)
				return
			}

			tokenDetailsErr = errors.Wrap(tokenDetailsErr, "Couldn't get token metadata")
			server.log(errors.WithMessage(tokenDetailsErr, errTokenDamaged.Error()))
			server.error(writer, req, http.StatusBadRequest, errTokenDamaged)
			return
		}

		if err := tokenUtils.DeleteAuth(tokenDetails.Uuid); err != nil {
			server.log(errors.Wrap(err, "Couldn't delete token"))
			server.error(writer, req, http.StatusInternalServerError, errSomethingWentWrong)
			return
		}

		server.respond(writer, req, http.StatusOK, nil)
		return
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
				return
			case tokenUtils.ErrTokenWrongFormat:
				server.error(writer, req, http.StatusBadRequest, tokenExtractErr)
				return
			}
		}

		// Access token must be valid, but expired
		accessTokenDetails, accessTokenDetailsErr := tokenUtils.ExtractTokenMetadata(tokenString)
		if accessTokenDetailsErr != nil {
			validErr, ok := accessTokenDetailsErr.(*jwt.ValidationError)
			if !ok || validErr.Errors != jwt.ValidationErrorExpired {
				accessTokenDetailsErr = errors.Wrap(accessTokenDetailsErr, "Couldn't get access token metadata")
				server.log(errors.WithMessage(accessTokenDetailsErr, errTokenDamaged.Error()))
				server.error(writer, req, http.StatusBadRequest, errTokenDamaged)
				return
			}
		} else {
			if err := tokenUtils.DeleteAuth(accessTokenDetails.Uuid); err != nil {
				server.log(errors.Wrap(err, "Couldn't delete access token"))
				server.error(writer, req, http.StatusInternalServerError, errSomethingWentWrong)
				return
			}
		}

		// Extract, validate and delete refresh token
		if err := tokenUtils.IsValid(requestStruct.RefreshToken); err != nil {
			server.log(errors.WithMessage(err, errTokenExpiredOrDeleted.Error()))
			server.error(writer, req, http.StatusForbidden, errTokenExpiredOrDeleted)
			return
		}

		refreshTokenDetails, refreshTokenDetailsErr := tokenUtils.ExtractTokenMetadata(requestStruct.RefreshToken)
		if refreshTokenDetailsErr != nil {
			refreshTokenDetailsErr = errors.Wrap(refreshTokenDetailsErr, "Couldn't get refresh token metadata")
			server.log(errors.WithMessage(refreshTokenDetailsErr, errTokenDamaged.Error()))
			server.error(writer, req, http.StatusBadRequest, errTokenDamaged)
			return
		}

		if err := tokenUtils.DeleteAuth(refreshTokenDetails.Uuid); err != nil {
			server.log(errors.Wrap(err, "Couldn't delete token"))
			server.error(writer, req, http.StatusInternalServerError, errSomethingWentWrong)
			return
		}

		// Creating new pair of tokens
		tokenDetails, createErr := tokenUtils.CreateTokens(refreshTokenDetails.UserId)
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
		return
	}
}

func (server *server) handleUsersMe() http.HandlerFunc {
	return func(writer http.ResponseWriter, req *http.Request) {
		server.respond(writer, req, http.StatusOK, req.Context().Value(ctxKeyUser).(*model.User))
		return
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
			if errors.Cause(err) == model.ErrValidationFailed {
				server.error(writer, req, http.StatusBadRequest, err)
				return
			}
			server.error(writer, req, http.StatusInternalServerError, err)
			return
		}

		if err := tokenUtils.DeleteAllAuths(user.ID); err != nil {
			server.log(errors.Wrap(err, "Couldn't delete all tokens for user"))
			server.error(writer, req, http.StatusInternalServerError, errSomethingWentWrong)
			return
		}

		// Creating new pair of tokens
		tokenDetails, createErr := tokenUtils.CreateTokens(user.ID)
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
		return
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
			if errors.Cause(err) == model.ErrValidationFailed {
				server.error(writer, req, http.StatusBadRequest, err)
				return
			}
			server.error(writer, req, http.StatusInternalServerError, err)
			return
		}

		if err := tokenUtils.DeleteAllAuths(user.ID); err != nil {
			server.log(errors.Wrap(err, "Couldn't delete all tokens for user"))
			server.error(writer, req, http.StatusInternalServerError, errSomethingWentWrong)
			return
		}

		// Creating new pair of tokens
		tokenDetails, createErr := tokenUtils.CreateTokens(user.ID)
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
		return
	}
}
