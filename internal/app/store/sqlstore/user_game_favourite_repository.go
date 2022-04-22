package sqlstore

import (
	"database/sql"
	"fmt"

	"github.com/pkg/errors"

	"github.com/spolyakovs/price-hunter-ITMO/internal/app/model"
	"github.com/spolyakovs/price-hunter-ITMO/internal/app/store"
)

type UserGameFavouriteRepository struct {
	store *Store
}

func (userGameFavouriteRepository *UserGameFavouriteRepository) Create(userGameFavourite *model.UserGameFavourite) error {
	repositoryName := "UserGameFavourite"
	methodName := "Create"
	errWrapMessage := fmt.Sprintf(store.ErrRepositoryMessageFormat, repositoryName, methodName)

	createQuery := "INSERT INTO user_game_favourites (game_id, user_id) VALUES ($1, $2) RETURNING id;"

	if err := userGameFavouriteRepository.store.db.Get(
		&userGameFavourite.ID,
		createQuery,
		userGameFavourite.Game.ID,
		userGameFavourite.User.ID,
	); err != nil {
		return errors.Wrap(errors.WithMessage(store.ErrUnknownSQL, err.Error()), errWrapMessage)
	}

	return nil
}

func (userGameFavouriteRepository *UserGameFavouriteRepository) Find(id uint64) (*model.UserGameFavourite, error) {
	return userGameFavouriteRepository.FindBy("id", id)
}

// TODO: test especially this (userGameFavourite -> game -> publisher)
func (userGameFavouriteRepository *UserGameFavouriteRepository) FindBy(columnName string, value interface{}) (*model.UserGameFavourite, error) {
	repositoryName := "UserGameFavourite"
	methodName := "FindBy"
	errWrapMessage := fmt.Sprintf(store.ErrRepositoryMessageFormat, repositoryName, methodName)

	userGameFavourite := &model.UserGameFavourite{}
	findQuery := fmt.Sprintf("SELECT "+
		"user_game_favourites.id AS id, "+

		"games.id AS \"game.id\", "+
		"games.header_image_url AS \"game.header_image_url\", "+
		"games.name AS \"game.name\", "+
		"games.description AS \"game.description\", "+

		"publishers.id AS \"game.publisher.id\", "+
		"publishers.name AS \"game.publisher.name\", "+

		"users.id AS \"user.id\", "+
		"users.username AS \"user.username\", "+
		"users.email AS \"user.email\" "+

		"FROM games "+

		"LEFT JOIN games "+
		"ON (user_game_favourites.game_id = games.id) "+

		"LEFT JOIN publishers "+
		"ON (games.publisher_id = publishers.id) "+

		"LEFT JOIN users "+
		"ON (user_game_favourites.user_id = users.id) "+

		"WHERE %s = $1 LIMIT 1;", columnName)

	if err := userGameFavouriteRepository.store.db.Get(
		userGameFavourite,
		findQuery,
		value,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.Wrap(store.ErrNotFound, errWrapMessage)
		}

		return nil, errors.Wrap(errors.WithMessage(store.ErrUnknownSQL, err.Error()), errWrapMessage)
	}

	return userGameFavourite, nil
}

func (userGameFavouriteRepository *UserGameFavouriteRepository) FindAllBy(columnName string, value interface{}) ([]*model.UserGameFavourite, error) {
	repositoryName := "UserGameFavourite"
	methodName := "FindAllBy"
	errWrapMessage := fmt.Sprintf(store.ErrRepositoryMessageFormat, repositoryName, methodName)

	userGameFavourites := []*model.UserGameFavourite{}
	findQuery := fmt.Sprintf("SELECT "+
		"user_game_favourites.id AS id, "+

		"games.id AS \"game.id\", "+
		"games.header_image_url AS \"game.header_image_url\", "+
		"games.name AS \"game.name\", "+
		"games.description AS \"game.description\", "+

		"publishers.id AS \"game.publisher.id\", "+
		"publishers.name AS \"game.publisher.name\", "+

		"users.id AS \"user.id\", "+
		"users.username AS \"user.username\", "+
		"users.email AS \"user.email\" "+

		"FROM games "+

		"LEFT JOIN games "+
		"ON (user_game_favourites.game_id = games.id) "+

		"LEFT JOIN publishers "+
		"ON (games.publisher_id = publishers.id) "+

		"LEFT JOIN users "+
		"ON (user_game_favourites.user_id = users.id) "+

		"WHERE %s = $1 LIMIT 1;", columnName)

	if err := userGameFavouriteRepository.store.db.Select(
		&userGameFavourites,
		findQuery,
		value,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.Wrap(store.ErrNotFound, errWrapMessage)
		}

		return nil, errors.Wrap(errors.WithMessage(store.ErrUnknownSQL, err.Error()), errWrapMessage)
	}

	return userGameFavourites, nil
}

func (userGameFavouriteRepository *UserGameFavouriteRepository) Update(newGame *model.UserGameFavourite) error {
	repositoryName := "UserGameFavourite"
	methodName := "Update"
	errWrapMessage := fmt.Sprintf(store.ErrRepositoryMessageFormat, repositoryName, methodName)

	updateQuery := "UPDATE user_game_favourites " +
		"SET game_id = :game.id, " +
		"SET user_id = :user.id, " +
		"WHERE id = :id;"

	countResult, countResultErr := userGameFavouriteRepository.store.db.NamedExec(
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

func (userGameFavouriteRepository *UserGameFavouriteRepository) Delete(id uint64) error {
	repositoryName := "UserGameFavourite"
	methodName := "Delete"
	errWrapMessage := fmt.Sprintf(store.ErrRepositoryMessageFormat, repositoryName, methodName)

	deleteQuery := "DELETE FROM user_game_favourites WHERE id = $1;"

	countResult, countResultErr := userGameFavouriteRepository.store.db.Exec(
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
