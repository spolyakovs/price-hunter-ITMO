package apiserver

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/golang-jwt/jwt"
	"github.com/gorilla/mux"
	"github.com/spolyakovs/price-hunter-ITMO/internal/app/model"
	"github.com/spolyakovs/price-hunter-ITMO/internal/app/tokenUtils"
)

//TODO: restrincting number of games in games_list request (and add offset param)

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

		tokenDetails, tokenErr := tokenUtils.CreateToken(user.ID)
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

		tokenDetails, tokenErr := tokenUtils.CreateToken(user.ID)
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
		tokenString := tokenUtils.ExtractToken(req)
		tokenData, tokenDataErr := tokenUtils.ExtractAccessTokenMetadata(tokenString)
		if tokenDataErr != nil {
			validErr, ok := tokenDataErr.(*jwt.ValidationError)
			if ok && validErr.Errors == jwt.ValidationErrorExpired {
				server.respond(writer, req, http.StatusOK, nil)
				return
			}
			// TODO: change error
			fmt.Printf("DEBUG: %s\n", tokenDataErr.Error())
			server.error(writer, req, http.StatusUnauthorized, tokenDataErr)
			return
		}

		delErr := tokenUtils.DeleteAuth(tokenData.AccessUuid)
		if delErr != nil {
			// TODO: change error
			fmt.Printf("DEBUG: %s\n", delErr.Error())
			server.error(writer, req, http.StatusInternalServerError, delErr)
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

		tokenData, tokenDataErr := tokenUtils.ExtractRefreshTokenMetadata(requestStruct.RefreshToken)
		if tokenDataErr != nil {
			// TODO: change error
			fmt.Printf("DEBUG: %s\n", tokenDataErr.Error())
			server.error(writer, req, http.StatusUnauthorized, tokenDataErr)
			return
		}

		delErr := tokenUtils.DeleteAuth(tokenData.RefreshUuid)
		if delErr != nil {
			fmt.Printf("DEBUG: %s\n", delErr.Error())
			server.error(writer, req, http.StatusInternalServerError, delErr)
			return
		}

		tokenDetails, createErr := tokenUtils.CreateToken(tokenData.UserId)
		if createErr != nil {
			fmt.Printf("DEBUG: %s\n", createErr.Error())
			server.error(writer, req, http.StatusInternalServerError, createErr)
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
			server.error(writer, req, http.StatusNotFound, errUserDoesNotExist)
			return
		}

		server.respond(writer, req, http.StatusOK, user)
	}
}
