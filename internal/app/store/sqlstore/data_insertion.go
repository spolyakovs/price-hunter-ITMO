package sqlstore

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/spolyakovs/price-hunter-ITMO/internal/app/model"
	"github.com/spolyakovs/price-hunter-ITMO/internal/app/store"
)

var tableNames = []string{
	"game_market_prices",
	"game_tags",
	"user_game_favourites",
	// "markets",
	"tags",
	"games",
	"publishers",
	"users",
}

func (st *Store) ClearTables() error {
	errWrapMessage := "Clearing tables to insert test data"

	for _, tableName := range tableNames {
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

	var markets []*model.Market

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
