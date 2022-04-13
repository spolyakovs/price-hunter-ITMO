package apiserver

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/pkg/errors"
	"github.com/spolyakovs/price-hunter-ITMO/internal/app/tokenUtils"

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
		logger.Infof("started %s %s\n", req.Method, req.RequestURI)

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
		tokenString := tokenUtils.ExtractToken(req)
		tokenData, tokenDataErr := tokenUtils.ExtractAccessTokenMetadata(tokenString)
		if tokenDataErr != nil {
			// TODO: change error
			tokenDataErr = errors.Wrap(tokenDataErr, "Couldn't get token metadata")
			server.log(tokenDataErr)
			server.error(writer, req, http.StatusUnauthorized, tokenDataErr)
			return
		}
		userId, userIdErr := tokenUtils.FetchAuth(tokenData)
		if userIdErr != nil {
			// TODO: change error
			userIdErr = errors.Wrap(userIdErr, "Couldn't fetch auth")
			server.log(userIdErr)
			server.error(writer, req, http.StatusUnauthorized, userIdErr)
			return
		}

		user, err := server.store.Users().Find(userId)
		if err != nil {
			// TODO: probably this is how errors should be logged and sent
			err = errors.Wrap(err, "Couldn't find user")
			err = errors.Wrap(err, errNotAuthenticated.Error())
			server.log(err)
			server.error(writer, req, http.StatusUnauthorized, errNotAuthenticated)
			return
		}

		next.ServeHTTP(writer, req.WithContext(context.WithValue(req.Context(), ctxKeyUser, user)))
	})
}

type stackTracer interface {
	StackTrace() errors.StackTrace
}

func (server *server) log(err error) {
	fmt.Printf("%v\n", err)
}

func (server *server) error(writer http.ResponseWriter, req *http.Request, code int, err error) {
	server.respond(writer, req, code, map[string]string{"error": fmt.Sprintf("%v", err)})
}

func (server *server) respond(writer http.ResponseWriter, req *http.Request, code int, data interface{}) {
	writer.WriteHeader(code)
	if data != nil {
		json.NewEncoder(writer).Encode(data)
	}
}
