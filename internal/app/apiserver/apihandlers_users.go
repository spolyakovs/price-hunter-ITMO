package apiserver

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/spolyakovs/price-hunter-ITMO/internal/app/model"
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
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	return func(writer http.ResponseWriter, req *http.Request) {
		requestStruct := &request{}
		if err := json.NewDecoder(req.Body).Decode(requestStruct); err != nil {
			server.error(writer, req, http.StatusBadRequest, err)
			return
		}

		user, err := server.store.Users().FindBy("email", requestStruct.Email)
		if err != nil || !user.ComparePassword(requestStruct.Password) {
			fmt.Printf("DEBUG: %s\n", err.Error())
			server.error(writer, req, http.StatusUnauthorized, errIncorrectEmailOrPassword)
			return
		}

		session, err := server.sessionStore.Get(req, sessionName)
		if err != nil {
			server.error(writer, req, http.StatusInternalServerError, err)
			return
		}

		session.Values["user_id"] = user.ID
		if err := server.sessionStore.Save(req, writer, session); err != nil {
			server.error(writer, req, http.StatusInternalServerError, err)
			return
		}

		server.respond(writer, req, http.StatusOK, nil)
	}
}

func (server *server) handleUsersMe() http.HandlerFunc {
	return func(writer http.ResponseWriter, req *http.Request) {
		server.respond(writer, req, http.StatusOK, req.Context().Value(ctxKeyUser).(*model.User))
	}
}

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

		session, err := server.sessionStore.Get(req, sessionName)
		if err != nil {
			server.error(writer, req, http.StatusInternalServerError, err)
			return
		}
		session.Values["user_id"] = nil

		server.respond(writer, req, http.StatusOK, nil)
	}
}

func (server *server) handleUsersGetByID() http.HandlerFunc {
	return func(writer http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)

		id, err := strconv.Atoi(vars["id"])
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
