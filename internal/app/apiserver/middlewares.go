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

// TODO: fix errors in middlewares and handlers like in tokenUtils, model.User and sqlstore.UserRepository

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
		methodName := "AuthenticateUser"
		errWrapMessage := fmt.Sprintf(errMiddlewareMessageFormat, methodName)

		tokenString, tokenExtractErr := tokenUtils.ExtractToken(req)
		if tokenExtractErr != nil {
			errWrapped := errors.Wrap(tokenExtractErr, errWrapMessage)

			switch errors.Cause(tokenExtractErr) {
			case tokenUtils.ErrTokenNotProvided:
				server.error(writer, req, http.StatusUnauthorized, errWrapped)
			case tokenUtils.ErrTokenWrongFormat:
				server.error(writer, req, http.StatusBadRequest, errWrapped)
			default:
				// Mostly TokenUtils.ErrInternal, probably something with Redis
				server.log(errWrapped)
				server.error(writer, req, http.StatusInternalServerError, errSomethingWentWrong)
			}

			return
		}

		// Already checks in tokenUtils.ExtractTokenMetadata

		// if err := tokenUtils.IsValid(tokenString); err != nil {
		// 	err = errors.Wrap(err, errWrapMessage)
		//
		// 	switch errors.Cause(err) {
		// 	case tokenUtils.ErrTokenDamaged:
		// 		server.error(writer, req, http.StatusBadRequest, err)
		// 	case tokenUtils.ErrTokenExpiredOrDeleted:
		// 		server.error(writer, req, http.StatusForbidden, err)
		// 	default:
		// 		// Mostly TokenUtils.ErrInternal, probably something with Redis
		// 		server.log(err)
		// 		server.error(writer, req, http.StatusInternalServerError, errSomethingWentWrong)
		// 	}
		//
		// 	return
		// }

		tokenDetails, tokenDetailsErr := tokenUtils.ExtractTokenMetadata(tokenString)
		if tokenDetailsErr != nil {
			errWrapped := errors.Wrap(tokenDetailsErr, errWrapMessage)

			switch errors.Cause(tokenDetailsErr) {
			case tokenUtils.ErrTokenDamaged:
				server.error(writer, req, http.StatusBadRequest, errWrapped)
			case tokenUtils.ErrTokenExpiredOrDeleted:
				server.error(writer, req, http.StatusForbidden, errWrapped)
			default:
				// Mostly TokenUtils.ErrInternal, probably something with Redis
				server.log(errWrapped)
				server.error(writer, req, http.StatusInternalServerError, errSomethingWentWrong)
			}
			return
		}

		userId, userIdErr := tokenUtils.FetchAuth(tokenDetails)
		if userIdErr != nil {
			errWrapped := errors.Wrap(userIdErr, errWrapMessage)

			switch errors.Cause(userIdErr) {
			case tokenUtils.ErrTokenExpiredOrDeleted:
				server.error(writer, req, http.StatusForbidden, errWrapped)
			default:
				// Mostly TokenUtils.ErrInternal, probably something with Redis
				server.log(errWrapped)
				server.error(writer, req, http.StatusInternalServerError, errSomethingWentWrong)
				return
			}
		}

		user, userErr := server.store.Users().Find(userId)
		if userErr != nil {
			errWrapped := errors.Wrap(userErr, errWrapMessage)

			server.log(errWrapped)
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
	server.log(err)
	server.respond(writer, req, code, map[string]string{"error": errors.Cause(err).Error()})
}

func (server *server) respond(writer http.ResponseWriter, req *http.Request, code int, data interface{}) {
	writer.WriteHeader(code)
	if data != nil {
		json.NewEncoder(writer).Encode(data)
	}
}
