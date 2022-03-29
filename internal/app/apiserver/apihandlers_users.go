package apiserver

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/spolyakovs/price-hunter-ITMO/internal/app/model"
	"github.com/spolyakovs/price-hunter-ITMO/internal/app/tokenUtils"
	"net/http"
	"strconv"
)

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

		ts, tokenErr := tokenUtils.CreateToken(user.ID)
		if tokenErr != nil {
			fmt.Printf("DEBUG: %s\n", tokenErr.Error())
			server.error(writer, req, http.StatusInternalServerError, tokenErr)
			return
		}

		saveErr := tokenUtils.CreateAuth(user.ID, ts)
		if saveErr != nil {
			fmt.Printf("DEBUG: %s\n", saveErr.Error())
			server.error(writer, req, http.StatusInternalServerError, saveErr)
			return
		}

		tokens := map[string]string{
			"access_token":  ts.AccessToken,
			"refresh_token": ts.RefreshToken,
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
