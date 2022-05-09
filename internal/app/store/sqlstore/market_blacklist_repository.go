package sqlstore

import (
	"database/sql"
	"fmt"

	"github.com/pkg/errors"

	"github.com/spolyakovs/price-hunter-ITMO/internal/app/model"
	"github.com/spolyakovs/price-hunter-ITMO/internal/app/store"
)

type MarketBlacklistItemRepository struct {
	store *Store
}

func (marketBlacklistItemRepository *MarketBlacklistItemRepository) Create(marketBlacklistItem *model.MarketBlacklistItem) error {
	repositoryName := "MarketBlacklistItem"
	methodName := "Create"
	errWrapMessage := fmt.Sprintf(store.ErrRepositoryMessageFormat, repositoryName, methodName)

	createQuery := "INSERT INTO market_blacklist (market_game_url, market_id) VALUES ($1, $2) RETURNING id;"

	if err := marketBlacklistItemRepository.store.db.Get(
		&marketBlacklistItem.ID,
		createQuery,
		marketBlacklistItem.MarketGameURL,
		marketBlacklistItem.Market.ID,
	); err != nil {
		return errors.Wrap(errors.Wrap(store.ErrUnknownSQL, err.Error()), errWrapMessage)
	}

	return nil
}

func (marketBlacklistItemRepository *MarketBlacklistItemRepository) CheckByURL(marketGameURL string) (bool, error) {
	repositoryName := "MarketBlacklistItem"
	methodName := "CheckByURL"
	errWrapMessage := fmt.Sprintf(store.ErrRepositoryMessageFormat, repositoryName, methodName)

	var result bool
	checkQuery := "SELECT EXISTS(SELECT " +
		"market_blacklist.id AS id, " +
		"market_blacklist.market_game_url AS market_game_url, " +

		"markets.id AS \"market.id\", " +
		"markets.name AS \"market.name\" " +

		"FROM market_blacklist " +

		"LEFT JOIN markets " +
		"ON (market_blacklist.market_id = markets.id) " +

		"WHERE market_blacklist.market_game_url = $1);"

	if err := marketBlacklistItemRepository.store.db.Get(
		&result,
		checkQuery,
		marketGameURL,
	); err != nil {
		if err == sql.ErrNoRows {
			return false, errors.Wrap(store.ErrNotFound, errWrapMessage)
		}

		return false, errors.Wrap(errors.Wrap(store.ErrUnknownSQL, err.Error()), errWrapMessage)
	}

	return result, nil
}

func (marketBlacklistItemRepository *MarketBlacklistItemRepository) Delete(id uint64) error {
	repositoryName := "MarketBlacklistItem"
	methodName := "Delete"
	errWrapMessage := fmt.Sprintf(store.ErrRepositoryMessageFormat, repositoryName, methodName)

	deleteQuery := "DELETE FROM market_blacklist WHERE id = $1;"

	countResult, err := marketBlacklistItemRepository.store.db.Exec(
		deleteQuery,
		id,
	)

	if err != nil {
		return errors.Wrap(errors.Wrap(store.ErrUnknownSQL, err.Error()), errWrapMessage)
	}

	count, err := countResult.RowsAffected()

	if err != nil {
		return errors.Wrap(errors.Wrap(store.ErrUnknownSQL, err.Error()), errWrapMessage)
	}

	if count == 0 {
		return errors.Wrap(store.ErrNotFound, errWrapMessage)
	}

	return nil
}
