package apiserver

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pkg/errors"

	"github.com/spolyakovs/price-hunter-ITMO/internal/app/model"
	"github.com/spolyakovs/price-hunter-ITMO/internal/app/store"
	"github.com/spolyakovs/price-hunter-ITMO/internal/app/tokenutils"
)

// TODO: add comments
// TODO: think about moving all "response" errors to apiserver/errors, and in other packages will be internal errors (check like "switch error.Cause(err)")

func (server *server) handleRegistration() http.HandlerFunc {
	type request struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	return func(writer http.ResponseWriter, req *http.Request) {
		methodName := "Registration"
		errWrapMessage := fmt.Sprintf(errHandlerMessageFormat, methodName)

		requestStruct := &request{}
		if err := json.NewDecoder(req.Body).Decode(requestStruct); err != nil {
			errWrapped := errors.Wrap(err, errWrapMessage)
			server.log(errWrapped)
			server.error(writer, req, http.StatusBadRequest, errWrongRequestFormat)
			return
		}

		user := &model.User{
			Username: requestStruct.Username,
			Email:    requestStruct.Email,
			Password: requestStruct.Password,
		}

		responseError := map[string]string{}

		if _, err := server.store.Users().FindBy("username", user.Username); err == nil {
			responseError["username"] = errUserExistsUsernameMessage
		} else {
			if errors.Cause(err) != store.ErrNotFound {
				errWrapped := errors.Wrap(err, errWrapMessage)
				server.log(errWrapped)
				server.error(writer, req, http.StatusInternalServerError, errSomethingWentWrong)
				return
			}
		}

		if _, err := server.store.Users().FindBy("email", user.Email); err == nil {
			responseError["email"] = errUserExistsEmailMessage
		} else {
			if errors.Cause(err) != store.ErrNotFound {
				errWrapped := errors.Wrap(err, errWrapMessage)
				server.log(errWrapped)
				server.error(writer, req, http.StatusInternalServerError, errSomethingWentWrong)
				return
			}
		}

		_, usernameErrorExists := responseError["username"]
		_, emailErrorExists := responseError["email"]
		if usernameErrorExists || emailErrorExists {
			response := map[string]map[string]string{}
			response["error"] = responseError

			server.respond(writer, req, http.StatusOK, response)
			return
		}

		if err := server.store.Users().Create(user); err != nil {
			errWrapped := errors.Wrap(err, errWrapMessage)

			switch errors.Cause(errWrapped) {
			case model.ErrValidationFailed:
				server.error(writer, req, http.StatusBadRequest, errWrapped)
			default:
				server.log(errWrapped)
				server.error(writer, req, http.StatusInternalServerError, errSomethingWentWrong)
			}

			return
		}

		server.respond(writer, req, http.StatusOK, map[string]string{})
	}
}

func (server *server) handleLogin() http.HandlerFunc {
	type request struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	return func(writer http.ResponseWriter, req *http.Request) {
		methodName := "Login"
		errWrapMessage := fmt.Sprintf(errHandlerMessageFormat, methodName)

		requestStruct := &request{}
		if err := json.NewDecoder(req.Body).Decode(requestStruct); err != nil {
			errWrapped := errors.Wrap(err, errWrapMessage)
			server.log(errWrapped)
			server.error(writer, req, http.StatusBadRequest, errWrongRequestFormat)
			return
		}

		user, err := server.store.Users().FindBy("username", requestStruct.Username)
		if err != nil {
			errWrapped := errors.Wrap(err, errWrapMessage)

			switch errors.Cause(err) {
			case store.ErrNotFound:
				server.respond(writer, req, http.StatusOK, map[string]string{"error": errWrongUsernameOrPasswordMessage})
			default:
				server.log(errWrapped)
				server.error(writer, req, http.StatusInternalServerError, errSomethingWentWrong)
			}

			return
		}

		if !user.ComparePassword(requestStruct.Password) {
			server.respond(writer, req, http.StatusOK, map[string]string{"error": errWrongUsernameOrPasswordMessage})
			return
		}

		tokenDetails, err := tokenutils.CreateTokens(user.ID)
		if err != nil {
			errWrapped := errors.Wrap(err, errWrapMessage)
			server.log(errWrapped)
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

func (server *server) handleLogout() http.HandlerFunc {
	return func(writer http.ResponseWriter, req *http.Request) {
		methodName := "Logout"
		errWrapMessage := fmt.Sprintf(errHandlerMessageFormat, methodName)

		tokenString, err := tokenutils.ExtractToken(req)
		if err != nil {
			errWrapped := errors.Wrap(err, errWrapMessage)

			switch errors.Cause(err) {
			case tokenutils.ErrTokenNotProvided:
				server.error(writer, req, http.StatusUnauthorized, errWrapped)
			case tokenutils.ErrTokenWrongFormat:
				server.error(writer, req, http.StatusBadRequest, errWrapped)
			default:
				// Mostly TokenUtils.ErrInternal, probably something with Redis
				server.log(errWrapped)
				server.error(writer, req, http.StatusInternalServerError, errSomethingWentWrong)
			}

			return
		}

		tokenDetails, err := tokenutils.ExtractTokenMetadata(tokenString)
		if err != nil {
			errWrapped := errors.Wrap(err, errWrapMessage)

			switch errors.Cause(err) {
			case tokenutils.ErrTokenDamaged:
				server.error(writer, req, http.StatusBadRequest, errWrapped)
			case tokenutils.ErrTokenExpiredOrDeleted:
				server.respond(writer, req, http.StatusOK, map[string]string{})
			default:
				// Mostly TokenUtils.ErrInternal, probably something with Redis
				server.log(errWrapped)
				server.error(writer, req, http.StatusInternalServerError, errSomethingWentWrong)
			}
			return
		}

		if err := tokenutils.DeleteAuth(tokenDetails.Uuid); err != nil {
			errWrapped := errors.Wrap(err, errWrapMessage)
			server.log(errWrapped)
			server.error(writer, req, http.StatusInternalServerError, errSomethingWentWrong)
			return
		}

		server.respond(writer, req, http.StatusOK, map[string]string{})
	}
}

func (server *server) handleRefreshToken() http.HandlerFunc {
	type request struct {
		RefreshToken string `json:"refresh_token"`
	}

	return func(writer http.ResponseWriter, req *http.Request) {
		methodName := "RefreshToken"
		errWrapMessage := fmt.Sprintf(errHandlerMessageFormat, methodName)

		requestStruct := &request{}
		if err := json.NewDecoder(req.Body).Decode(requestStruct); err != nil {
			errWrapped := errors.Wrap(err, errWrapMessage)
			server.log(errWrapped)
			server.error(writer, req, http.StatusBadRequest, errWrongRequestFormat)
			return
		}

		// Extract, validate and delete access token
		tokenString, err := tokenutils.ExtractToken(req)
		if err != nil {
			errWrapped := errors.Wrap(err, errWrapMessage)

			switch errors.Cause(err) {
			case tokenutils.ErrTokenNotProvided:
				server.error(writer, req, http.StatusUnauthorized, errWrapped)
			case tokenutils.ErrTokenWrongFormat:
				server.error(writer, req, http.StatusBadRequest, errWrapped)
			default:
				// Mostly TokenUtils.ErrInternal, probably something with Redis
				server.log(errWrapped)
				server.error(writer, req, http.StatusInternalServerError, errSomethingWentWrong)
			}

			return
		}

		// Access token must be valid, but expired
		accessTokenDetails, err := tokenutils.ExtractTokenMetadata(tokenString)
		if err != nil {
			errWrapped := errors.Wrap(err, errWrapMessage)

			switch errors.Cause(err) {
			case tokenutils.ErrTokenDamaged:
				server.error(writer, req, http.StatusBadRequest, errWrapped)
				return
			case tokenutils.ErrTokenExpiredOrDeleted:
				break
			default:
				// Mostly TokenUtils.ErrInternal, probably something with Redis
				server.log(errWrapped)
				server.error(writer, req, http.StatusInternalServerError, errSomethingWentWrong)
				return
			}
		} else {
			if err := tokenutils.DeleteAuth(accessTokenDetails.Uuid); err != nil {
				errWrapped := errors.Wrap(err, errWrapMessage)
				server.log(errWrapped)
				server.error(writer, req, http.StatusInternalServerError, errSomethingWentWrong)
				return
			}
		}

		refreshTokenDetails, err := tokenutils.ExtractTokenMetadata(requestStruct.RefreshToken)
		if err != nil {
			errWrapped := errors.Wrap(err, errWrapMessage)

			switch errors.Cause(err) {
			case tokenutils.ErrTokenDamaged:
				server.error(writer, req, http.StatusBadRequest, errWrapped)
			case tokenutils.ErrTokenExpiredOrDeleted:
				server.error(writer, req, http.StatusForbidden, errWrapped)
			default:
				// Mostly TokenUtils.ErrInternal, probably something with Redis
				server.log(errWrapped)
				server.error(writer, req, http.StatusInternalServerError, errSomethingWentWrong)
			}

			return
		}

		if err := tokenutils.DeleteAuth(refreshTokenDetails.Uuid); err != nil {
			errWrapped := errors.Wrap(err, errWrapMessage)
			server.log(errWrapped)
			server.error(writer, req, http.StatusInternalServerError, errSomethingWentWrong)
			return
		}

		// Creating new pair of tokens
		tokenDetails, err := tokenutils.CreateTokens(refreshTokenDetails.UserId)
		if err != nil {
			errWrapped := errors.Wrap(err, errWrapMessage)
			server.log(errWrapped)
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

func (server *server) handleUsersChangeEmail() http.HandlerFunc {
	type request struct {
		NewEmail string `json:"new_email"`
	}

	return func(writer http.ResponseWriter, req *http.Request) {
		methodName := "UserChangeEmail"
		errWrapMessage := fmt.Sprintf(errHandlerMessageFormat, methodName)

		requestStruct := &request{}
		if err := json.NewDecoder(req.Body).Decode(requestStruct); err != nil {
			errWrapped := errors.Wrap(err, errWrapMessage)
			server.log(errWrapped)
			server.error(writer, req, http.StatusBadRequest, errWrongRequestFormat)
			return
		}

		if _, err := server.store.Users().FindBy("email", requestStruct.NewEmail); err == nil {
			server.respond(writer, req, http.StatusOK, map[string]string{"error": errUserExistsEmailMessage})
			return
		} else {
			if errors.Cause(err) != store.ErrNotFound {
				errWrapped := errors.Wrap(err, errWrapMessage)
				server.log(errWrapped)
				server.error(writer, req, http.StatusInternalServerError, errSomethingWentWrong)
				return
			}
		}

		user := req.Context().Value(ctxKeyUser).(*model.User)

		if err := server.store.Users().UpdateEmail(requestStruct.NewEmail, user.ID); err != nil {
			errWrapped := errors.Wrap(err, errWrapMessage)

			switch errors.Cause(errWrapped) {
			case model.ErrValidationFailed:
				server.error(writer, req, http.StatusBadRequest, errWrapped)
			default:
				server.log(errWrapped)
				server.error(writer, req, http.StatusInternalServerError, errSomethingWentWrong)
			}

			return
		}

		if err := tokenutils.DeleteAllAuths(user.ID); err != nil {
			errWrapped := errors.Wrap(err, errWrapMessage)
			server.log(errWrapped)
			server.error(writer, req, http.StatusInternalServerError, errSomethingWentWrong)
			return
		}

		// Creating new pair of tokens
		tokenDetails, err := tokenutils.CreateTokens(user.ID)
		if err != nil {
			errWrapped := errors.Wrap(err, errWrapMessage)
			server.log(errWrapped)
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

func (server *server) handleUsersChangePassword() http.HandlerFunc {
	type request struct {
		CurrentPassword string `json:"current_password"`
		NewPassword     string `json:"new_password"`
	}

	return func(writer http.ResponseWriter, req *http.Request) {
		methodName := "UserChangePassword"
		errWrapMessage := fmt.Sprintf(errHandlerMessageFormat, methodName)

		requestStruct := &request{}
		if err := json.NewDecoder(req.Body).Decode(requestStruct); err != nil {
			server.error(writer, req, http.StatusBadRequest, err)
			return
		}

		user := req.Context().Value(ctxKeyUser).(*model.User)

		if !user.ComparePassword(requestStruct.CurrentPassword) {
			server.respond(writer, req, http.StatusOK, map[string]string{"error": errWrongPasswordMessage})
			return
		}

		if err := server.store.Users().UpdatePassword(requestStruct.NewPassword, user.ID); err != nil {
			errWrapped := errors.Wrap(err, errWrapMessage)

			switch errors.Cause(errWrapped) {
			case model.ErrValidationFailed:
				server.error(writer, req, http.StatusBadRequest, errWrapped)
			default:
				server.log(errWrapped)
				server.error(writer, req, http.StatusInternalServerError, errSomethingWentWrong)
			}

			return
		}

		if err := tokenutils.DeleteAllAuths(user.ID); err != nil {
			errWrapped := errors.Wrap(err, errWrapMessage)
			server.log(errWrapped)
			server.error(writer, req, http.StatusInternalServerError, errSomethingWentWrong)
			return
		}

		// Creating new pair of tokens
		tokenDetails, err := tokenutils.CreateTokens(user.ID)
		if err != nil {
			errWrapped := errors.Wrap(err, errWrapMessage)
			server.log(errWrapped)
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
