package apiserver

import (
	"github.com/spolyakovs/price-hunter-ITMO/internal/app/tokenUtils"
	"net/http"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/spolyakovs/price-hunter-ITMO/internal/app/store/sqlstore"
)

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

	// TODO: migrate to JWT

	os.Setenv("ACCESS_SECRET", config.AccessSecret)
	os.Setenv("REFRESH_SECRET", config.RefreshSecret)

	redisErr := tokenUtils.SetupRedis()
	if redisErr != nil {
		return redisErr
	}

	srv := newServer(store)
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
