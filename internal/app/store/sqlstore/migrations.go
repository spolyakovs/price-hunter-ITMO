package sqlstore

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/spolyakovs/price-hunter-ITMO/internal/app/store"
)

func (s *Store) createTables() error {
	errWrapMessage := "Creating tables error"

	tx, err := s.db.Beginx()
	if err != nil {
		errWrapped := errors.WithMessage(store.ErrUnknownSQL, err.Error())
		errWrapped = errors.Wrap(errWrapped, errWrapMessage)
		return errWrapped
	}

	if err := createTableUsers(tx); err != nil {
		errWrapped := errors.Wrap(err, errWrapMessage)
		return errWrapped
	}

	if err := createTablePublishers(tx); err != nil {
		errWrapped := errors.Wrap(err, errWrapMessage)
		return errWrapped
	}

	if err := createTableGames(tx); err != nil {
		errWrapped := errors.Wrap(err, errWrapMessage)
		return errWrapped
	}

	if err := createTableTags(tx); err != nil {
		errWrapped := errors.Wrap(err, errWrapMessage)
		return errWrapped
	}

	if err := createTableMarkets(tx); err != nil {
		errWrapped := errors.Wrap(err, errWrapMessage)
		return errWrapped
	}

	if err := createTableUserGameFavourites(tx); err != nil {
		errWrapped := errors.Wrap(err, errWrapMessage)
		return errWrapped
	}

	if err := createTableGameTags(tx); err != nil {
		errWrapped := errors.Wrap(err, errWrapMessage)
		return errWrapped
	}

	if err := createTableGameMarketPrices(tx); err != nil {
		errWrapped := errors.Wrap(err, errWrapMessage)
		return errWrapped
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		errWrapped := errors.WithMessage(store.ErrUnknownSQL, err.Error())
		errWrapped = errors.Wrap(errWrapped, errWrapMessage)
		return errWrapped
	}

	return nil
}

func createTableUsers(tx *sqlx.Tx) error {
	tableName := "Users"
	errWrapMessage := fmt.Sprintf(store.ErrCreateTablesMessageFormat, tableName)

	createTableUsersQuery := "CREATE TABLE IF NOT EXISTS users (" +
		"id bigserial NOT NULL PRIMARY KEY," +
		"username varchar NOT NULL UNIQUE," +
		"email varchar NOT NULL UNIQUE," +
		"encrypted_password varchar NOT NULL );"

	if _, err := tx.Exec(createTableUsersQuery); err != nil {
		errWrapped := errors.WithMessage(store.ErrUnknownSQL, err.Error())
		errWrapped = errors.Wrap(errWrapped, errWrapMessage)
		return errWrapped
	}

	return nil
}

func createTablePublishers(tx *sqlx.Tx) error {
	tableName := "Publishers"
	errWrapMessage := fmt.Sprintf(store.ErrCreateTablesMessageFormat, tableName)

	createTablePublishersQuery := "CREATE TABLE IF NOT EXISTS publishers (" +
		"id bigserial NOT NULL PRIMARY KEY," +
		"name varchar NOT NULL );"

	if _, err := tx.Exec(createTablePublishersQuery); err != nil {
		errWrapped := errors.WithMessage(store.ErrUnknownSQL, err.Error())
		errWrapped = errors.Wrap(errWrapped, errWrapMessage)
		return errWrapped
	}

	return nil
}

func createTableGames(tx *sqlx.Tx) error {
	tableName := "Users"
	errWrapMessage := fmt.Sprintf(store.ErrCreateTablesMessageFormat, tableName)

	createTableGamesQuery := "CREATE TABLE IF NOT EXISTS games (" +
		"id bigserial NOT NULL PRIMARY KEY," +
		"header_image_url varchar NOT NULL," +
		"name varchar NOT NULL," +
		"description varchar NOT NULL," +
		"publisher_id bigserial NOT NULL REFERENCES publishers (id) ON DELETE CASCADE );"

	if _, err := tx.Exec(createTableGamesQuery); err != nil {
		errWrapped := errors.WithMessage(store.ErrUnknownSQL, err.Error())
		errWrapped = errors.Wrap(errWrapped, errWrapMessage)
		return errWrapped
	}

	return nil
}

func createTableTags(tx *sqlx.Tx) error {
	tableName := "Users"
	errWrapMessage := fmt.Sprintf(store.ErrCreateTablesMessageFormat, tableName)

	createTableTagsQuery := "CREATE TABLE IF NOT EXISTS tags (" +
		"id bigserial NOT NULL PRIMARY KEY," +
		"name varchar NOT NULL );"

	if _, err := tx.Exec(createTableTagsQuery); err != nil {
		errWrapped := errors.WithMessage(store.ErrUnknownSQL, err.Error())
		errWrapped = errors.Wrap(errWrapped, errWrapMessage)
		return errWrapped
	}

	return nil
}

func createTableMarkets(tx *sqlx.Tx) error {
	tableName := "Users"
	errWrapMessage := fmt.Sprintf(store.ErrCreateTablesMessageFormat, tableName)

	createTableMarketsQuery := "CREATE TABLE IF NOT EXISTS markets (" +
		"id bigserial NOT NULL PRIMARY KEY," +
		"name varchar NOT NULL );"

	if _, err := tx.Exec(createTableMarketsQuery); err != nil {
		errWrapped := errors.WithMessage(store.ErrUnknownSQL, err.Error())
		errWrapped = errors.Wrap(errWrapped, errWrapMessage)
		return errWrapped
	}

	return nil
}

func createTableUserGameFavourites(tx *sqlx.Tx) error {
	tableName := "Users"
	errWrapMessage := fmt.Sprintf(store.ErrCreateTablesMessageFormat, tableName)

	createTableGameFavouritesQuery := "CREATE TABLE IF NOT EXISTS user_game_favourites (" +
		"id bigserial NOT NULL PRIMARY KEY," +
		"user_id bigserial NOT NULL REFERENCES users (id) ON DELETE CASCADE," +
		"game_id bigserial NOT NULL REFERENCES games (id) ON DELETE CASCADE );"

	if _, err := tx.Exec(createTableGameFavouritesQuery); err != nil {
		errWrapped := errors.WithMessage(store.ErrUnknownSQL, err.Error())
		errWrapped = errors.Wrap(errWrapped, errWrapMessage)
		return errWrapped
	}

	return nil
}

func createTableGameTags(tx *sqlx.Tx) error {
	tableName := "Users"
	errWrapMessage := fmt.Sprintf(store.ErrCreateTablesMessageFormat, tableName)

	createTableGameTagsQuery := "CREATE TABLE IF NOT EXISTS game_tags (" +
		"id bigserial NOT NULL PRIMARY KEY," +
		"game_id bigserial NOT NULL REFERENCES games (id) ON DELETE CASCADE," +
		"tag_id bigserial NOT NULL REFERENCES tags (id) ON DELETE CASCADE );"

	if _, err := tx.Exec(createTableGameTagsQuery); err != nil {
		errWrapped := errors.WithMessage(store.ErrUnknownSQL, err.Error())
		errWrapped = errors.Wrap(errWrapped, errWrapMessage)
		return errWrapped
	}

	return nil
}

func createTableGameMarketPrices(tx *sqlx.Tx) error {
	tableName := "Users"
	errWrapMessage := fmt.Sprintf(store.ErrCreateTablesMessageFormat, tableName)

	createTableGameMarketPricesQuery := "CREATE TABLE IF NOT EXISTS game_market_prices (" +
		"id bigserial NOT NULL PRIMARY KEY," +
		"initial_value_formatted varchar NOT NULL," +
		"final_value_formatted varchar NOT NULL," +
		"discount_percent integer NOT NULL," +
		"game_id bigserial NOT NULL REFERENCES games (id) ON DELETE CASCADE," +
		"market_id bigserial NOT NULL REFERENCES markets (id) ON DELETE CASCADE );"

	if _, err := tx.Exec(createTableGameMarketPricesQuery); err != nil {
		errWrapped := errors.WithMessage(store.ErrUnknownSQL, err.Error())
		errWrapped = errors.Wrap(errWrapped, errWrapMessage)
		return errWrapped
	}

	return nil
}
