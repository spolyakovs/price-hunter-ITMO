package sqlstore

import (
	"database/sql"
	"fmt"

	"github.com/pkg/errors"

	"github.com/spolyakovs/price-hunter-ITMO/internal/app/model"
	"github.com/spolyakovs/price-hunter-ITMO/internal/app/store"
)

type GameTagRepository struct {
	store *Store
}

func (gameTagRepository *GameTagRepository) Create(gameTag *model.GameTag) error {
	repositoryName := "GameTag"
	methodName := "Create"
	errWrapMessage := fmt.Sprintf(store.ErrRepositoryMessageFormat, repositoryName, methodName)

	createQuery := "INSERT INTO game_tags (game_id, tag_id) VALUES ($1, $2) RETURNING id;"

	if err := gameTagRepository.store.db.Get(
		&gameTag.ID,
		createQuery,
		gameTag.Game.ID,
		gameTag.Tag.ID,
	); err != nil {
		return errors.Wrap(errors.WithMessage(store.ErrUnknownSQL, err.Error()), errWrapMessage)
	}

	return nil
}

func (gameTagRepository *GameTagRepository) Find(id uint64) (*model.GameTag, error) {
	return gameTagRepository.FindBy("id", id)
}

// TODO: test especially this (gameTag -> game -> publisher)
func (gameTagRepository *GameTagRepository) FindBy(columnName string, value interface{}) (*model.GameTag, error) {
	repositoryName := "GameTag"
	methodName := "Find"
	errWrapMessage := fmt.Sprintf(store.ErrRepositoryMessageFormat, repositoryName, methodName)

	gameTag := &model.GameTag{}
	findQuery := fmt.Sprintf("SELECT "+
		"game_tags.id AS \"id\", "+

		"games.id AS \"game.id\", "+
		"games.header_image_url AS \"game.header_image_url\", "+
		"games.name AS \"game.name\", "+
		"games.description AS \"game.description\", "+

		"publishers.id AS \"game.publisher.id\", "+
		"publishers.name AS \"game.publisher.name\" "+

		"FROM games "+

		"LEFT JOIN games "+
		"ON (game_tags.game_id = games.id) "+

		"LEFT JOIN publishers "+
		"ON (games.publisher_id = publishers.id) "+

		"WHERE %s = $1 LIMIT 1;", columnName)

	if err := gameTagRepository.store.db.Get(
		gameTag,
		findQuery,
		value,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.Wrap(store.ErrNotFound, errWrapMessage)
		}

		return nil, errors.Wrap(errors.WithMessage(store.ErrUnknownSQL, err.Error()), errWrapMessage)
	}

	return gameTag, nil
}

func (gameTagRepository *GameTagRepository) FindAllBy(columnName string, value interface{}) ([]*model.GameTag, error) {
	repositoryName := "GameTag"
	methodName := "FindAllBy"
	errWrapMessage := fmt.Sprintf(store.ErrRepositoryMessageFormat, repositoryName, methodName)

	gameTags := []*model.GameTag{}
	findQuery := fmt.Sprintf("SELECT "+
		"game_tags.id AS \"id\", "+

		"games.id AS \"game.id\", "+
		"games.header_image_url AS \"game.header_image_url\", "+
		"games.name AS \"game.name\", "+
		"games.description AS \"game.description\", "+

		"publishers.id AS \"game.publisher.id\", "+
		"publishers.name AS \"game.publisher.name\" "+

		"tags.id AS \"tag.id\", "+
		"tags.name AS \"tag.name\" "+

		"FROM games "+

		"LEFT JOIN games "+
		"ON (game_tags.game_id = games.id) "+

		"LEFT JOIN publishers "+
		"ON (games.publisher_id = publishers.id) "+

		"LEFT JOIN tags "+
		"ON (game_tags.tag_id = tags.id) "+

		"WHERE %s = $1 LIMIT 1;", columnName)

	if err := gameTagRepository.store.db.Select(
		&gameTags,
		findQuery,
		value,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.Wrap(store.ErrNotFound, errWrapMessage)
		}

		return nil, errors.Wrap(errors.WithMessage(store.ErrUnknownSQL, err.Error()), errWrapMessage)
	}

	return gameTags, nil
}

func (gameTagRepository *GameTagRepository) Update(newGame *model.GameTag) error {
	repositoryName := "GameTag"
	methodName := "Update"
	errWrapMessage := fmt.Sprintf(store.ErrRepositoryMessageFormat, repositoryName, methodName)

	updateQuery := "UPDATE game_tags " +
		"SET game_id = :game.id, " +
		"SET tag_id = :tag.id, " +
		"WHERE id = :id;"

	countResult, countResultErr := gameTagRepository.store.db.NamedExec(
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

func (gameTagRepository *GameTagRepository) Delete(id uint64) error {
	repositoryName := "GameTag"
	methodName := "Delete"
	errWrapMessage := fmt.Sprintf(store.ErrRepositoryMessageFormat, repositoryName, methodName)

	deleteQuery := "DELETE FROM game_tags WHERE id = $1;"

	countResult, countResultErr := gameTagRepository.store.db.Exec(
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
