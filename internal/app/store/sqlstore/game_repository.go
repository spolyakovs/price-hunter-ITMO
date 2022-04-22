package sqlstore

import (
	"database/sql"
	"fmt"

	"github.com/pkg/errors"

	"github.com/spolyakovs/price-hunter-ITMO/internal/app/model"
	"github.com/spolyakovs/price-hunter-ITMO/internal/app/store"
)

type GameRepository struct {
	store *Store
}

func (gameRepository *GameRepository) Create(game *model.Game) error {
	repositoryName := "Game"
	methodName := "Create"
	errWrapMessage := fmt.Sprintf(store.ErrRepositoryMessageFormat, repositoryName, methodName)

	createQuery := "INSERT INTO games (header_image_url, name, description, publisher_id) VALUES ($1, $2, $3, $4) RETURNING id;"

	if err := gameRepository.store.db.Get(
		&game.ID,
		createQuery,
		game.HeaderImageURL,
		game.Name,
		game.Description,
		game.Publisher.ID,
	); err != nil {
		return errors.Wrap(errors.WithMessage(store.ErrUnknownSQL, err.Error()), errWrapMessage)
	}

	return nil
}

func (gameRepository *GameRepository) Find(id uint64) (*model.Game, error) {
	return gameRepository.FindBy("id", id)
}

func (gameRepository *GameRepository) FindBy(columnName string, value interface{}) (*model.Game, error) {
	repositoryName := "Game"
	methodName := "FindBy"
	errWrapMessage := fmt.Sprintf(store.ErrRepositoryMessageFormat, repositoryName, methodName)

	game := &model.Game{}
	findQuery := fmt.Sprintf("SELECT "+
		"games.id AS id, "+
		"games.header_image_url AS header_image_url, "+
		"games.name AS name, "+
		"games.description AS description, "+

		"publishers.id AS \"publisher.id\", "+
		"publishers.name AS \"publisher.name\" "+

		"FROM games "+

		"LEFT JOIN publishers "+
		"ON (games.publisher_id = publishers.id) "+

		"WHERE %s = $1 LIMIT 1;", columnName)

	if err := gameRepository.store.db.Get(
		game,
		findQuery,
		value,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.Wrap(store.ErrNotFound, errWrapMessage)
		}

		return nil, errors.Wrap(errors.WithMessage(store.ErrUnknownSQL, err.Error()), errWrapMessage)
	}

	return game, nil
}

// TODO: findAllBy (tags, ...)
func (gameRepository *GameRepository) FindAllBy(columnName string, value interface{}) ([]*model.Game, error) {
	repositoryName := "Game"
	methodName := "FindAllBy"
	errWrapMessage := fmt.Sprintf(store.ErrRepositoryMessageFormat, repositoryName, methodName)

	games := []*model.Game{}
	findQuery := fmt.Sprintf("SELECT "+
		"games.id AS id, "+
		"games.header_image_url AS header_image_url, "+
		"games.name AS name, "+
		"games.description AS description, "+

		"publishers.id AS \"publisher.id\", "+
		"publishers.name AS \"publisher.name\" "+

		"FROM games "+

		"LEFT JOIN publishers "+
		"ON (games.publisher_id = publishers.id) "+

		"WHERE %s = $1 LIMIT 1;", columnName)

	if err := gameRepository.store.db.Select(
		&games,
		findQuery,
		value,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.Wrap(store.ErrNotFound, errWrapMessage)
		}

		return nil, errors.Wrap(errors.WithMessage(store.ErrUnknownSQL, err.Error()), errWrapMessage)
	}

	return games, nil
}

func (gameRepository *GameRepository) Update(newGame *model.Game) error {
	repositoryName := "Game"
	methodName := "Update"
	errWrapMessage := fmt.Sprintf(store.ErrRepositoryMessageFormat, repositoryName, methodName)

	updateQuery := "UPDATE games " +
		"SET header_image_url = :header_image_url, " +
		"SET name = :name, " +
		"SET description = :description, " +
		"SET publisher_id = :publisher.id, " +
		"WHERE id = :id;"

	countResult, countResultErr := gameRepository.store.db.NamedExec(
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

func (gameRepository *GameRepository) Delete(id uint64) error {
	repositoryName := "Game"
	methodName := "Delete"
	errWrapMessage := fmt.Sprintf(store.ErrRepositoryMessageFormat, repositoryName, methodName)

	deleteQuery := "DELETE FROM games WHERE id = $1;"

	countResult, countResultErr := gameRepository.store.db.Exec(
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
