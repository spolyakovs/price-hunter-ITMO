package apiserver

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/spolyakovs/price-hunter-ITMO/internal/app/tokenUtils"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

func (server *server) setRequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
		id := uuid.New().String()
		writer.Header().Set("X-Request-ID", id)
		next.ServeHTTP(writer, req.WithContext(context.WithValue(req.Context(), ctxKeyRequestID, id)))
	})
}

func (server *server) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
		logger := server.logger.WithFields(logrus.Fields{
			"remote_addr": req.RemoteAddr,
			"request_id":  req.Context().Value(ctxKeyRequestID),
		})
		logger.Infof("started %s %s", req.Method, req.RequestURI)

		start := time.Now()
		rw := &responseWriter{writer, http.StatusOK}
		next.ServeHTTP(rw, req)

		var level logrus.Level
		switch {
		case rw.code >= 500:
			level = logrus.ErrorLevel
		case rw.code >= 400:
			level = logrus.WarnLevel
		default:
			level = logrus.InfoLevel
		}
		logger.Logf(
			level,
			"completed with %d %s in %v",
			rw.code,
			http.StatusText(rw.code),
			time.Now().Sub(start),
		)
	})
}

func (server *server) authenticateUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
		tokenAuth, tokenAuthErr := tokenUtils.ExtractTokenMetadata(req)
		if tokenAuthErr != nil {
			// TODO: change error
			fmt.Printf("DEBUG: %s\n", tokenAuthErr.Error())
			server.error(writer, req, http.StatusUnauthorized, tokenAuthErr)
			return
		}
		userId, userIdErr := tokenUtils.FetchAuth(tokenAuth)
		if userIdErr != nil {
			// TODO: change error
			fmt.Printf("DEBUG: %s\n", userIdErr.Error())
			server.error(writer, req, http.StatusUnauthorized, userIdErr)
			return
		}

		user, err := server.store.Users().Find(userId)
		if err != nil {
			fmt.Printf("DEBUG: %s\n", err.Error())
			server.error(writer, req, http.StatusUnauthorized, errNotAuthenticated)
			return
		}

		next.ServeHTTP(writer, req.WithContext(context.WithValue(req.Context(), ctxKeyUser, user)))
	})
}

func (server *server) error(writer http.ResponseWriter, req *http.Request, code int, err error) {
	server.respond(writer, req, code, map[string]string{"error": err.Error()})
}

func (server *server) respond(writer http.ResponseWriter, req *http.Request, code int, data interface{}) {
	writer.WriteHeader(code)
	if data != nil {
		json.NewEncoder(writer).Encode(data)
	}
}
