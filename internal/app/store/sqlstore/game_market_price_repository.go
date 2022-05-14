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

	createQuery := "INSERT INTO game_market_prices (initial_value_formatted, final_value_formatted, discount_percent, market_game_url, game_id, market_id) " +
		"VALUES ($1, $2, $3, $4, $5, $6) RETURNING id;"

	if err := gameMarketPriceRepository.store.db.Get(
		&gameMarketPrice.ID,
		createQuery,
		gameMarketPrice.InitialValueFormatted,
		gameMarketPrice.FinalValueFormatted,
		gameMarketPrice.DiscountPercent,
		gameMarketPrice.MarketGameURL,
		gameMarketPrice.Game.ID,
		gameMarketPrice.Market.ID,
	); err != nil {
		return errors.Wrap(errors.Wrap(store.ErrUnknownSQL, err.Error()), errWrapMessage)
	}

	return nil
}

func (gameMarketPriceRepository *GameMarketPriceRepository) Find(id uint64) (*model.GameMarketPrice, error) {
	return gameMarketPriceRepository.FindBy("id", id)
}

// TODO: move this func as Find(id)
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
		"game_market_prices.market_game_url AS market_game_url, "+

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

		return nil, errors.Wrap(errors.Wrap(store.ErrUnknownSQL, err.Error()), errWrapMessage)
	}

	return gameMarketPrice, nil
}

func (gameMarketPriceRepository *GameMarketPriceRepository) FindByGameMarket(game *model.Game, market *model.Market) (*model.GameMarketPrice, error) {
	repositoryName := "GameMarketPrice"
	methodName := "FindByGameMarket"
	errWrapMessage := fmt.Sprintf(store.ErrRepositoryMessageFormat, repositoryName, methodName)

	gameMarketPrice := &model.GameMarketPrice{}
	findQuery := "SELECT " +
		"game_market_prices.id AS id, " +
		"game_market_prices.initial_value_formatted AS initial_value_formatted, " +
		"game_market_prices.final_value_formatted AS final_value_formatted, " +
		"game_market_prices.discount_percent AS discount_percent, " +
		"game_market_prices.market_game_url AS market_game_url, " +

		"games.id AS \"game.id\", " +
		"games.header_image_url AS \"game.header_image_url\", " +
		"games.name AS \"game.name\", " +
		"games.description AS \"game.description\", " +

		"publishers.id AS \"game.publisher.id\", " +
		"publishers.name AS \"game.publisher.name\", " +

		"markets.id AS \"market.id\", " +
		"markets.name AS \"market.name\" " +

		"FROM game_market_prices " +

		"LEFT JOIN games " +
		"ON (game_market_prices.game_id = games.id) " +

		"LEFT JOIN publishers " +
		"ON (games.publisher_id = publishers.id) " +

		"LEFT JOIN markets " +
		"ON (game_market_prices.market_id = markets.id) " +

		"WHERE game_market_prices.game_id = $1 AND game_market_prices.market_id = $2 LIMIT 1;"

	if err := gameMarketPriceRepository.store.db.Get(
		gameMarketPrice,
		findQuery,
		game.ID,
		market.ID,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.Wrap(store.ErrNotFound, errWrapMessage)
		}

		return nil, errors.Wrap(errors.Wrap(store.ErrUnknownSQL, err.Error()), errWrapMessage)
	}

	return gameMarketPrice, nil
}

func (gameMarketPriceRepository *GameMarketPriceRepository) FindAllByGame(game *model.Game) ([]*model.GameMarketPrice, error) {
	repositoryName := "GameMarketPrice"
	methodName := "FindAllByGame"
	errWrapMessage := fmt.Sprintf(store.ErrRepositoryMessageFormat, repositoryName, methodName)

	gameMarketPrices := []*model.GameMarketPrice{}
	findQuery := "SELECT " +
		"game_market_prices.id AS id, " +
		"game_market_prices.initial_value_formatted AS initial_value_formatted, " +
		"game_market_prices.final_value_formatted AS final_value_formatted, " +
		"game_market_prices.discount_percent AS discount_percent, " +
		"game_market_prices.market_game_url AS market_game_url, " +

		"games.id AS \"game.id\", " +
		"games.header_image_url AS \"game.header_image_url\", " +
		"games.name AS \"game.name\", " +
		"games.description AS \"game.description\", " +

		"publishers.id AS \"game.publisher.id\", " +
		"publishers.name AS \"game.publisher.name\", " +

		"markets.id AS \"market.id\", " +
		"markets.name AS \"market.name\" " +

		"FROM game_market_prices " +

		"LEFT JOIN games " +
		"ON (game_market_prices.game_id = games.id) " +

		"LEFT JOIN publishers " +
		"ON (games.publisher_id = publishers.id) " +

		"LEFT JOIN markets " +
		"ON (game_market_prices.market_id = markets.id) " +

		"WHERE game_market_prices.game_id = $1;"

	if err := gameMarketPriceRepository.store.db.Select(
		&gameMarketPrices,
		findQuery,
		game.ID,
	); err != nil {
		if err == sql.ErrNoRows {
			return []*model.GameMarketPrice{}, nil
		}

		return nil, errors.Wrap(errors.Wrap(store.ErrUnknownSQL, err.Error()), errWrapMessage)
	}

	return gameMarketPrices, nil
}

func (gameMarketPriceRepository *GameMarketPriceRepository) Update(newGameMarket *model.GameMarketPrice) error {
	repositoryName := "GameMarketPrice"
	methodName := "Update"
	errWrapMessage := fmt.Sprintf(store.ErrRepositoryMessageFormat, repositoryName, methodName)

	updateQuery := "UPDATE game_market_prices " +
		"SET initial_value_formatted = :initial_value_formatted, " +
		"final_value_formatted = :final_value_formatted, " +
		"discount_percent = :discount_percent, " +
		"market_game_url = :market_game_url, " +
		"game_id = :game.id, " +
		"market_id = :market.id " +
		"WHERE id = :id;"

	countResult, err := gameMarketPriceRepository.store.db.NamedExec(
		updateQuery,
		newGameMarket,
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

func (gameMarketPriceRepository *GameMarketPriceRepository) Delete(id uint64) error {
	repositoryName := "GameMarketPrice"
	methodName := "Delete"
	errWrapMessage := fmt.Sprintf(store.ErrRepositoryMessageFormat, repositoryName, methodName)

	deleteQuery := "DELETE FROM game_market_prices WHERE id = $1;"

	countResult, err := gameMarketPriceRepository.store.db.Exec(
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
