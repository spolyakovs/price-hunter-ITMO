package sqlstore

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/spolyakovs/price-hunter-ITMO/internal/app/model"
	"github.com/spolyakovs/price-hunter-ITMO/internal/app/store"
)

var markets []*model.Market

var tableNamesToTruncate = []string{
	"game_market_prices",
	"game_tags",
	"user_game_favourites",
	// "markets", (do not truncate)
	"tags",
	"games",
	"publishers",
	"users",
}

func (st *Store) ClearTables() error {
	errWrapMessage := "Clearing tables to insert test data"

	for _, tableName := range tableNamesToTruncate {
		truncateQuery := fmt.Sprintf("TRUNCATE %s RESTART IDENTITY CASCADE;", tableName)
		if _, err := st.db.Exec(truncateQuery); err != nil {
			return errors.Wrap(err, errWrapMessage)
		}
	}

	return nil
}

func (st *Store) insertDataMarkets() error {
	tableName := "Markets"
	errWrapMessage := fmt.Sprintf(store.ErrTestDataInsertionMessageFormat, tableName)

	markets = append(markets, &model.Market{
		Name: "Steam",
	})

	markets = append(markets, &model.Market{
		Name: "EpicGamesStore",
	})

	markets = append(markets, &model.Market{
		Name: "GOG.com",
	})

	for _, market := range markets {
		if err := st.Markets().Create(market); err != nil {
			return errors.Wrap(err, errWrapMessage)
		}
	}

	return nil
}

func (st *Store) InsertTestData() error {
	FillTestData()

	if err := st.ClearTables(); err != nil {
		return err
	}

	if err := st.insertTestDataUsers(); err != nil {
		return err
	}

	if err := st.insertTestDataPublishers(); err != nil {
		return err
	}

	if err := st.insertTestDataGames(); err != nil {
		return err
	}

	if err := st.insertTestDataTags(); err != nil {
		return err
	}

	if err := st.insertTestDataMarkets(); err != nil {
		return err
	}

	if err := st.insertTestDataUserGameFavourites(); err != nil {
		return err
	}

	if err := st.insertTestDataGameTags(); err != nil {
		return err
	}

	if err := st.insertTestDataGameMarketPrices(); err != nil {
		return err
	}

	return nil
}

func (st *Store) insertTestDataUsers() error {
	tableName := "Users"
	errWrapMessage := fmt.Sprintf(store.ErrTestDataInsertionMessageFormat, tableName)

	for _, user := range TestUsers {
		if err := st.Users().Create(user); err != nil {
			return errors.Wrap(err, errWrapMessage)
		}
	}

	return nil
}

func (st *Store) insertTestDataPublishers() error {
	tableName := "Publishers"
	errWrapMessage := fmt.Sprintf(store.ErrTestDataInsertionMessageFormat, tableName)

	for _, publisher := range TestPublishers {
		if err := st.Publishers().Create(publisher); err != nil {
			return errors.Wrap(err, errWrapMessage)
		}
	}

	return nil
}

func (st *Store) insertTestDataGames() error {
	tableName := "Games"
	errWrapMessage := fmt.Sprintf(store.ErrTestDataInsertionMessageFormat, tableName)

	for _, game := range TestGames {
		if err := st.Games().Create(game); err != nil {
			return errors.Wrap(err, errWrapMessage)
		}
	}

	return nil
}

func (st *Store) insertTestDataTags() error {
	tableName := "Tags"
	errWrapMessage := fmt.Sprintf(store.ErrTestDataInsertionMessageFormat, tableName)

	for _, tag := range TestTags {
		if err := st.Tags().Create(tag); err != nil {
			return errors.Wrap(err, errWrapMessage)
		}
	}

	return nil
}

func (st *Store) insertTestDataMarkets() error {
	tableName := "Markets"
	errWrapMessage := fmt.Sprintf(store.ErrTestDataInsertionMessageFormat, tableName)

	markets = append(markets, &model.Market{
		Name: "Steam",
	})

	markets = append(markets, &model.Market{
		Name: "EpicGamesStore",
	})

	markets = append(markets, &model.Market{
		Name: "GOG.com",
	})

	for _, market := range markets {
		if err := st.Markets().Create(market); err != nil {
			return errors.Wrap(err, errWrapMessage)
		}
	}

	return nil
}

func (st *Store) insertTestDataUserGameFavourites() error {
	tableName := "UserGameFavourites"
	errWrapMessage := fmt.Sprintf(store.ErrTestDataInsertionMessageFormat, tableName)

	for _, userGameFavourite := range TestUserGameFavourites {
		if err := st.UserGameFavourites().Create(userGameFavourite); err != nil {
			return errors.Wrap(err, errWrapMessage)
		}
	}

	return nil
}

func (st *Store) insertTestDataGameTags() error {
	tableName := "GameTags"
	errWrapMessage := fmt.Sprintf(store.ErrTestDataInsertionMessageFormat, tableName)

	for _, gameTag := range TestGameTags {
		if err := st.GameTags().Create(gameTag); err != nil {
			return errors.Wrap(err, errWrapMessage)
		}
	}

	return nil
}

func (st *Store) insertTestDataGameMarketPrices() error {
	tableName := "GameMarketPrices"
	errWrapMessage := fmt.Sprintf(store.ErrTestDataInsertionMessageFormat, tableName)

	for _, gameMarketPrice := range TestGameMarketPrices {
		if err := st.GameMarketPrices().Create(gameMarketPrice); err != nil {
			return errors.Wrap(err, errWrapMessage)
		}
	}

	return nil
}
