package sqlstore

import (
	"database/sql"
	"fmt"

	"github.com/pkg/errors"

	"github.com/spolyakovs/price-hunter-ITMO/internal/app/model"
	"github.com/spolyakovs/price-hunter-ITMO/internal/app/store"
)

type GameMarketPriceRepository struct {
	store *Store
}

func (gameMarketPriceRepository *GameMarketPriceRepository) Create(gameMarketPrice *model.GameMarketPrice) error {
	repositoryName := "GameMarketPrice"
	methodName := "Create"
	errWrapMessage := fmt.Sprintf(store.ErrRepositoryMessageFormat, repositoryName, methodName)

	createQuery := "INSERT INTO game_market_prices (initial_value_formatted, final_value_formatted, discount_percent, game_id, market_id) VALUES ($1, $2, $3, $4, $5) RETURNING id;"

	if err := gameMarketPriceRepository.store.db.Get(
		&gameMarketPrice.ID,
		createQuery,
		gameMarketPrice.InitialValueFormatted,
		gameMarketPrice.FinalValueFormatted,
		gameMarketPrice.DiscountPercent,
		gameMarketPrice.Game.ID,
		gameMarketPrice.Market.ID,
	); err != nil {
		return errors.Wrap(errors.WithMessage(store.ErrUnknownSQL, err.Error()), errWrapMessage)
	}

	return nil
}

func (gameMarketPriceRepository *GameMarketPriceRepository) Find(id uint64) (*model.GameMarketPrice, error) {
	return gameMarketPriceRepository.FindBy("id", id)
}

// TODO: test especially this (gameMarketPrice -> game -> publisher)
func (gameMarketPriceRepository *GameMarketPriceRepository) FindBy(columnName string, value interface{}) (*model.GameMarketPrice, error) {
	repositoryName := "GameMarketPrice"
	methodName := "FindBy"
	errWrapMessage := fmt.Sprintf(store.ErrRepositoryMessageFormat, repositoryName, methodName)

	gameMarketPrice := &model.GameMarketPrice{}
	findQuery := fmt.Sprintf("SELECT "+
		"game_market_prices.id AS id, "+
		"game_market_prices.initial_value_formatted AS initial_value_formatted, "+
		"game_market_prices.final_value_formatted AS final_value_formatted, "+
		"game_market_prices.discount_percent AS discount_percent, "+

		"games.id AS \"game.id\", "+
		"games.header_image_url AS \"game.header_image_url\", "+
		"games.name AS \"game.name\", "+
		"games.description AS \"game.description\", "+

		"publishers.id AS \"game.publisher.id\", "+
		"publishers.name AS \"game.publisher.name\" "+

		"markets.id AS \"market.id\", "+
		"markets.name AS \"market.name\" "+

		"FROM games "+

		"LEFT JOIN games "+
		"ON (game_market_prices.game_id = games.id) "+

		"LEFT JOIN publishers "+
		"ON (games.publisher_id = publishers.id) "+

		"LEFT JOIN markets "+
		"ON (game_market_prices.market_id = markets.id) "+

		"WHERE %s = $1 LIMIT 1;", columnName)

	if err := gameMarketPriceRepository.store.db.Get(
		gameMarketPrice,
		findQuery,
		value,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.Wrap(store.ErrNotFound, errWrapMessage)
		}

		return nil, errors.Wrap(errors.WithMessage(store.ErrUnknownSQL, err.Error()), errWrapMessage)
	}

	return gameMarketPrice, nil
}

func (gameMarketPriceRepository *GameMarketPriceRepository) FindAllBy(columnName string, value interface{}) ([]*model.GameMarketPrice, error) {
	repositoryName := "GameMarketPrice"
	methodName := "FindAllBy"
	errWrapMessage := fmt.Sprintf(store.ErrRepositoryMessageFormat, repositoryName, methodName)

	gameMarketPrices := []*model.GameMarketPrice{}
	findQuery := fmt.Sprintf("SELECT "+
		"game_market_prices.id AS id, "+
		"game_market_prices.initial_value_formatted AS initial_value_formatted, "+
		"game_market_prices.final_value_formatted AS final_value_formatted, "+
		"game_market_prices.discount_percent AS discount_percent, "+

		"games.id AS \"game.id\", "+
		"games.header_image_url AS \"game.header_image_url\", "+
		"games.name AS \"game.name\", "+
		"games.description AS \"game.description\", "+

		"publishers.id AS \"game.publisher.id\", "+
		"publishers.name AS \"game.publisher.name\" "+

		"markets.id AS \"market.id\", "+
		"markets.name AS \"market.name\" "+

		"FROM games "+

		"LEFT JOIN games "+
		"ON (game_market_prices.game_id = games.id) "+

		"LEFT JOIN publishers "+
		"ON (games.publisher_id = publishers.id) "+

		"LEFT JOIN markets "+
		"ON (game_market_prices.market_id = markets.id) "+

		"WHERE %s = $1 LIMIT 1;", columnName)

	if err := gameMarketPriceRepository.store.db.Select(
		&gameMarketPrices,
		findQuery,
		value,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.Wrap(store.ErrNotFound, errWrapMessage)
		}

		return nil, errors.Wrap(errors.WithMessage(store.ErrUnknownSQL, err.Error()), errWrapMessage)
	}

	return gameMarketPrices, nil
}

func (gameMarketPriceRepository *GameMarketPriceRepository) Update(newGame *model.GameMarketPrice) error {
	repositoryName := "GameMarketPrice"
	methodName := "Update"
	errWrapMessage := fmt.Sprintf(store.ErrRepositoryMessageFormat, repositoryName, methodName)

	updateQuery := "UPDATE game_market_prices " +
		"SET initial_value_formatted = :initial_value_formatted, " +
		"SET final_value_formatted = :final_value_formatted, " +
		"SET discount_percent = :discount_percent, " +
		"SET game_id = :game.id, " +
		"SET market_id = :market.id, " +
		"WHERE id = :id;"

	countResult, countResultErr := gameMarketPriceRepository.store.db.NamedExec(
		updateQuery,
		newGame,
	)

	if countResultErr != nil {
		return errors.Wrap(errors.WithMessage(store.ErrUnknownSQL, countResultErr.Error()), errWrapMessage)
	}

	count, countErr := countResult.RowsAffected()

	if countErr != nil {
		return errors.Wrap(errors.WithMessage(store.ErrUnknownSQL, countErr.Error()), errWrapMessage)
	}

	if count == 0 {
		return errors.Wrap(store.ErrNotFound, errWrapMessage)
	}

	return nil
}

func (gameMarketPriceRepository *GameMarketPriceRepository) Delete(id uint64) error {
	repositoryName := "GameMarketPrice"
	methodName := "Delete"
	errWrapMessage := fmt.Sprintf(store.ErrRepositoryMessageFormat, repositoryName, methodName)

	deleteQuery := "DELETE FROM game_market_prices WHERE id = $1;"

	countResult, countResultErr := gameMarketPriceRepository.store.db.Exec(
		deleteQuery,
		id,
	)

	if countResultErr != nil {
		return errors.Wrap(errors.WithMessage(store.ErrUnknownSQL, countResultErr.Error()), errWrapMessage)
	}

	count, countErr := countResult.RowsAffected()

	if countErr != nil {
		return errors.Wrap(errors.WithMessage(store.ErrUnknownSQL, countErr.Error()), errWrapMessage)
	}

	if count == 0 {
		return errors.Wrap(store.ErrNotFound, errWrapMessage)
	}

	return nil
}
