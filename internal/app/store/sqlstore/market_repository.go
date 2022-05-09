package sqlstore

import (
	"database/sql"
	"fmt"

	"github.com/pkg/errors"

	"github.com/spolyakovs/price-hunter-ITMO/internal/app/model"
	"github.com/spolyakovs/price-hunter-ITMO/internal/app/store"
)

type MarketRepository struct {
	store *Store
}

func (marketRepository *MarketRepository) Create(market *model.Market) error {
	repositoryName := "Market"
	methodName := "Create"
	errWrapMessage := fmt.Sprintf(store.ErrRepositoryMessageFormat, repositoryName, methodName)

	createQuery := "INSERT INTO markets (name) VALUES ($1) " +
		"ON CONFLICT(name) DO UPDATE SET name = EXCLUDED.name RETURNING id;"

	if err := marketRepository.store.db.Get(
		&market.ID,
		createQuery,
		market.Name,
	); err != nil {
		return errors.Wrap(errors.Wrap(store.ErrUnknownSQL, err.Error()), errWrapMessage)
	}

	return nil
}

func (marketRepository *MarketRepository) Find(id uint64) (*model.Market, error) {
	return marketRepository.FindBy("id", id)
}

func (marketRepository *MarketRepository) FindBy(columnName string, value interface{}) (*model.Market, error) {
	repositoryName := "Market"
	methodName := "FindBy"
	errWrapMessage := fmt.Sprintf(store.ErrRepositoryMessageFormat, repositoryName, methodName)

	market := &model.Market{}
	findQuery := fmt.Sprintf("SELECT * FROM markets WHERE %s = $1 LIMIT 1;", columnName)

	if err := marketRepository.store.db.Get(
		market,
		findQuery,
		value,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.Wrap(store.ErrNotFound, errWrapMessage)
		}

		return nil, errors.Wrap(errors.Wrap(store.ErrUnknownSQL, err.Error()), errWrapMessage)
	}

	return market, nil
}

func (marketRepository *MarketRepository) Update(newMarket *model.Market) error {
	repositoryName := "Market"
	methodName := "Update"
	errWrapMessage := fmt.Sprintf(store.ErrRepositoryMessageFormat, repositoryName, methodName)

	updateQuery := "UPDATE markets " +
		"SET name = :name " +
		"WHERE id = :id;"

	countResult, err := marketRepository.store.db.NamedExec(
		updateQuery,
		newMarket,
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

func (marketRepository *MarketRepository) Delete(id uint64) error {
	repositoryName := "Market"
	methodName := "Delete"
	errWrapMessage := fmt.Sprintf(store.ErrRepositoryMessageFormat, repositoryName, methodName)

	deleteQuery := "DELETE FROM markets WHERE id = $1;"

	countResult, err := marketRepository.store.db.Exec(
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
