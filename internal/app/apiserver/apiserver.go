package apiserver

import (
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/spolyakovs/price-hunter-ITMO/internal/app/store/sqlstore"
)

// TODO: read about GoLang 1.18

func Start(config *Config) error {
	startLogger := logrus.New()

	startLogger.Info("Creating database")
	db, dbErr := newDB(config.DatabaseURL)
	if dbErr != nil {
		return dbErr
	}

	defer db.Close()

	startLogger.Info("Configuring store")
	store, storeErr := sqlstore.New(db)
	if storeErr != nil {
		return storeErr
	}

	// TODO: change cookies to JWT(???), prob need to use self written db table

	sessionStore := sessions.NewCookieStore([]byte(config.SessionKey))
	srv := newServer(store, sessionStore)
	startLogger.Info("Server started")

	return http.ListenAndServe(config.BindAddr, srv)
}

func newDB(dbURL string) (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres", dbURL)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
