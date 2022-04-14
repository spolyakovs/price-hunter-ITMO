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

		if err := tokenUtils.IsValid(tokenString); err != nil {
			server.log(errors.WithMessage(err, errTokenExpiredOrDeleted.Error()))
			server.error(writer, req, http.StatusForbidden, errTokenExpiredOrDeleted)
			return
		}

		tokenData, tokenDataErr := tokenUtils.ExtractTokenMetadata(tokenString)
		if tokenDataErr != nil {
			tokenDataErr = errors.Wrap(tokenDataErr, "Couldn't get token metadata")
			server.log(errors.WithMessage(tokenDataErr, errTokenDamaged.Error()))
			server.error(writer, req, http.StatusBadRequest, errTokenDamaged)
			return
		}

		userId, userIdErr := tokenUtils.FetchAuth(tokenData)
		if userIdErr != nil {
			server.log(errors.Wrap(userIdErr, "Couldn't fetch auth"))
			server.error(writer, req, http.StatusForbidden, errTokenExpiredOrDeleted)
			return
		}

		user, err := server.store.Users().Find(userId)
		if err != nil {
			server.log(errors.Wrap(err, "Couldn't find user"))
			server.error(writer, req, http.StatusInternalServerError, errSomethingWentWrong)
			return
		}

		next.ServeHTTP(writer, req.WithContext(context.WithValue(req.Context(), ctxKeyUser, user)))
	})
}

type stackTracer interface {
	StackTrace() errors.StackTrace
}

func (server *server) log(err error) {
	// TODO: create log file, log not only error but request as well
	// TODO: think about different levels of logging (400+ status, 500+ status)
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
