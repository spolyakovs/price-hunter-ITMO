package apiserver

import (
	"fmt"
	"net/http"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/spolyakovs/price-hunter-ITMO/internal/app/store/sqlstore"
	"github.com/spolyakovs/price-hunter-ITMO/internal/app/tokenUtils"
)

// Migrate from wrapping errors to logging (Lexa vk)

func Start(config *Config) error {
	startLogger := logrus.New()

	startLogger.Info("Creating database")
	db, dbErr := newDB(*config)
	if dbErr != nil {
		return dbErr
	}

	defer db.Close()

	startLogger.Info("Configuring store")
	store, storeErr := sqlstore.New(db)
	if storeErr != nil {
		return storeErr
	}

	os.Setenv("TOKEN_SECRET", config.TokenSecret)

	redisErr := tokenUtils.SetupRedis(config.RedisAddr)
	if redisErr != nil {
		return redisErr
	}

	srv := newServer(store)
	startLogger.Info("Server started")

	return http.ListenAndServe(config.BindAddr, srv)
}

func newDB(config Config) (*sqlx.DB, error) {
	dbURL := fmt.Sprintf("host=%s dbname=%s user=%s password=%s sslmode=%s",
		config.DatabaseHost, config.DatabaseDBName, config.DatabaseUser, config.DatabasePassword, config.DatabaseSSLMode)
	db, err := sqlx.Open("postgres", dbURL)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
