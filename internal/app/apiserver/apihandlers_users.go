package apiserver

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/golang-jwt/jwt"
	"github.com/gorilla/mux"
	"github.com/spolyakovs/price-hunter-ITMO/internal/app/model"
	"github.com/spolyakovs/price-hunter-ITMO/internal/app/tokenUtils"
)

//TODO: add refresh_token request
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

		server.respond(writer, req, http.StatusCreated, user)
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

		saveErr := tokenUtils.CreateAuth(user.ID, tokenDetails)
		if saveErr != nil {
			fmt.Printf("DEBUG: %s\n", saveErr.Error())
			server.error(writer, req, http.StatusInternalServerError, saveErr)
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
		//TODO: think about errors
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

		saveErr := tokenUtils.CreateAuth(tokenData.UserId, tokenDetails)
		if saveErr != nil {
			fmt.Printf("DEBUG: %s\n", saveErr.Error())
			server.error(writer, req, http.StatusInternalServerError, saveErr)
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

// TODO: userUpdateEmail (StatusOk)
// TODO: userUpdatePassword (New token, старые ВСЕ удаляются)

func (server *server) handleUsersUpdate() http.HandlerFunc {
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
			ID:       req.Context().Value(ctxKeyUser).(*model.User).ID,
			Username: requestStruct.Username,
			Email:    requestStruct.Email,
			Password: requestStruct.Password,
		}

		if err := server.store.Users().Update(user); err != nil {
			server.error(writer, req, http.StatusInternalServerError, err)
			return
		}
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
