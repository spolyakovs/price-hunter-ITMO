package apiserver

import (
	"fmt"
	"net/http"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/spolyakovs/price-hunter-ITMO/internal/app/apistore"
	"github.com/spolyakovs/price-hunter-ITMO/internal/app/store"
	"github.com/spolyakovs/price-hunter-ITMO/internal/app/store/sqlstore"
	"github.com/spolyakovs/price-hunter-ITMO/internal/app/tokenutils"
	"golang.org/x/sync/errgroup"
)

// Migrate from wrapping errors to logging (Lexa vk)

func Start(config *Config) error {
	startLogger := logrus.New()

	startLogger.Info("Creating database")
	db, err := NewDB(*config)
	if err != nil {
		return err
	}

	defer db.Close()

	startLogger.Info("Configuring store")
	store, err := sqlstore.New(db)
	if err != nil {
		return err
	}

	os.Setenv("TOKEN_SECRET", config.TokenSecret)

	startLogger.Info("Configuring Redis")
	if err := tokenutils.SetupRedis(config.RedisAddr); err != nil {
		return err
	}

	startLogger.Info("Updating games info")
	if err := updateGames(*config, store); err != nil {
		return err
	}

	srv := newServer(store)
	startLogger.Info("Server started")

	return http.ListenAndServe(config.BindAddr, srv)
}

func updateGames(config Config, st store.Store) error {
	apiSteam := *apistore.NewAPISteam(config.SteamAPIKey, st)
	apiEpicGames := *apistore.NewAPIEpicGames(st)
	apiGOG := *apistore.NewAPIGOG(st)

	if err := apiSteam.GetGames(); err != nil {
		return err
	}
	// if err := apiEpicGames.GetGames(); err != nil {
	// 	return err
	// }
	// if err := apiGOG.GetGames(); err != nil {
	// 	return err
	// }

	g := new(errgroup.Group)

	g.Go(func() error {
		return apiEpicGames.GetGames()
	})
	g.Go(func() error {
		return apiGOG.GetGames()
	})

	if err := g.Wait(); err != nil {
		return err
	}

	// if err := apiSteam.UpdateGameMarketPrices(); err != nil {
	// 	return err
	// }

	return nil
}

func NewDB(config Config) (*sqlx.DB, error) {
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
