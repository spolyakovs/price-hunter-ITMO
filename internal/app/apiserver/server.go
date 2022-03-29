package apiserver

import (
	"net/http"

	"github.com/sirupsen/logrus"

	"github.com/gorilla/mux"
	"github.com/spolyakovs/price-hunter-ITMO/internal/app/store"
)

const (
	ctxKeyUser ctxKey = iota
	ctxKeyRequestID
)

type ctxKey int8

type server struct {
	router     *mux.Router
	logger     *logrus.Logger
	store      store.Store
	sessionKey []byte
}

func newServer(store store.Store) *server {
	server := &server{
		router: mux.NewRouter(),
		logger: logrus.New(),
		store:  store,
	}

	server.configureRouter()

	return server
}

func (server *server) ServeHTTP(writer http.ResponseWriter, req *http.Request) {
	server.router.ServeHTTP(writer, req)
}
