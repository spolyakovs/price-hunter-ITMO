package sqlstore

import (
	"database/sql"
	"fmt"
	"strings"

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

	createQuery := "INSERT INTO games (header_image_url, name, description, release_date, publisher_id) VALUES ($1, $2, $3, TO_DATE('$4', 'dd.MM.YYYY'), $5) RETURNING id;"

	if err := gameRepository.store.db.Get(
		&game.ID,
		createQuery,
		game.HeaderImageURL,
		game.Name,
		game.Description,
		game.ReleaseDate,
		game.Publisher.ID,
	); err != nil {
		return errors.Wrap(errors.Wrap(store.ErrUnknownSQL, err.Error()), errWrapMessage)
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
		"publishers.id AS \"publisher.id\", "+
		"publishers.name AS \"publisher.name\" "+

		"games.id AS id, "+
		"games.header_image_url AS header_image_url, "+
		"games.name AS name, "+
		"TO_CHAR(games.release_date, 'dd.MM.YYYY') AS release_date, "+
		"games.description AS description, "+

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

		return nil, errors.Wrap(errors.Wrap(store.ErrUnknownSQL, err.Error()), errWrapMessage)
	}

	return game, nil
}

func (gameRepository *GameRepository) FindAllByUser(user *model.User) ([]*model.Game, error) {
	repositoryName := "Game"
	methodName := "FindAllByUser"
	errWrapMessage := fmt.Sprintf(store.ErrRepositoryMessageFormat, repositoryName, methodName)

	games := []*model.Game{}

	findQuery := "SELECT " +
		"publishers.id AS \"publisher.id\", " +
		"publishers.name AS \"publisher.name\" " +

		"games.id AS id, " +
		"games.header_image_url AS header_image_url, " +
		"games.name AS name, " +
		"TO_CHAR(games.release_date, 'dd.MM.YYYY') AS release_date, " +
		"games.description AS description, " +

		"FROM games " +

		"LEFT JOIN publishers " +
		"ON (games.publisher_id = publishers.id) " +

		"LEFT JOIN user_game_favourites " +
		"ON (games.id = user_game_favourites.game_id) " +
		"WHERE user_game_favourites.user_id = $1;"

	if err := gameRepository.store.db.Select(
		&games,
		findQuery,
		user.ID,
	); err != nil {
		if err == sql.ErrNoRows {
			return []*model.Game{}, nil
		}

		return nil, errors.Wrap(errors.Wrap(store.ErrUnknownSQL, err.Error()), errWrapMessage)
	}

	return games, nil
}

func (gameRepository *GameRepository) FindAllByQueryTags(query string, tags []*model.Tag) ([]*model.Game, error) {
	repositoryName := "Game"
	methodName := "FindAllByQueryTags"
	errWrapMessage := fmt.Sprintf(store.ErrRepositoryMessageFormat, repositoryName, methodName)

	query = strings.ReplaceAll(query, "%", "\\%")
	query = strings.ReplaceAll(query, "_", "\\_")
	query = "%" + query + "%"
	args := []interface{}{}
	args = append(args, query)
	games := []*model.Game{}

	findQuery := "SELECT " +
		"publishers.id AS \"publisher.id\", " +
		"publishers.name AS \"publisher.name\", " +

		"games.id AS id, " +
		"games.header_image_url AS header_image_url, " +
		"games.name AS name, " +
		"TO_CHAR(games.release_date, 'dd.MM.YYYY') AS release_date, " +
		"games.description AS description " +

		"FROM games " +

		"LEFT JOIN publishers " +
		"ON (games.publisher_id = publishers.id) " +

		"LEFT JOIN game_tags " +
		"ON (games.id = game_tags.game_id) " +
		"WHERE (games.name LIKE $1 OR publishers.name LIKE $1)"

	if len(tags) != 0 {
		findQuery += " AND id IN (" +
			"    SELECT DISTINCT game_id FROM game_tags WHERE tag_id in ($2)" +
			")"

		tagsString := ""
		for _, tag := range tags {
			if tagsString != "" {
				tagsString += ", "
			}
			tagsString += fmt.Sprint(tag.ID)
		}

		args = append(args, tagsString)
	}

	findQuery += ";"

	println(len(args))

	if err := gameRepository.store.db.Select(
		&games,
		findQuery,
		args...,
	); err != nil {
		if err == sql.ErrNoRows {
			return []*model.Game{}, nil
		}

		return nil, errors.Wrap(errors.Wrap(store.ErrUnknownSQL, err.Error()), errWrapMessage)
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
		"SET release_date = TO_DATE(':release_date', 'dd.MM.YYYY'), " +
		"SET publisher_id = :publisher.id, " +
		"WHERE id = :id;"

	countResult, err := gameRepository.store.db.NamedExec(
		updateQuery,
		newGame,
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

func (gameRepository *GameRepository) Delete(id uint64) error {
	repositoryName := "Game"
	methodName := "Delete"
	errWrapMessage := fmt.Sprintf(store.ErrRepositoryMessageFormat, repositoryName, methodName)

	deleteQuery := "DELETE FROM games WHERE id = $1;"

	countResult, err := gameRepository.store.db.Exec(
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
