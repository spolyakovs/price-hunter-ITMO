package sqlstore

import (
	"database/sql"
	"fmt"

	"github.com/pkg/errors"

	"github.com/spolyakovs/price-hunter-ITMO/internal/app/model"
	"github.com/spolyakovs/price-hunter-ITMO/internal/app/store"
)

type TagRepository struct {
	store *Store
}

func (tagRepository *TagRepository) Create(tag *model.Tag) error {
	repositoryName := "Tag"
	methodName := "Create"
	errWrapMessage := fmt.Sprintf(store.ErrRepositoryMessageFormat, repositoryName, methodName)

	createQuery := "INSERT INTO tags (name) VALUES ($1) " +
		"ON CONFLICT(name) DO UPDATE SET name = EXCLUDED.name RETURNING id;"

	if err := tagRepository.store.db.Get(
		&tag.ID,
		createQuery,
		tag.Name,
	); err != nil {
		return errors.Wrap(errors.Wrap(store.ErrUnknownSQL, err.Error()), errWrapMessage)
	}

	return nil
}

func (tagRepository *TagRepository) Find(id uint64) (*model.Tag, error) {
	return tagRepository.FindBy("id", id)
}

func (tagRepository *TagRepository) FindBy(columnName string, value interface{}) (*model.Tag, error) {
	repositoryName := "Tag"
	methodName := "FindBy"
	errWrapMessage := fmt.Sprintf(store.ErrRepositoryMessageFormat, repositoryName, methodName)

	tag := &model.Tag{}
	findQuery := fmt.Sprintf("SELECT * FROM tags WHERE %s = $1 LIMIT 1;", columnName)

	if err := tagRepository.store.db.Get(
		tag,
		findQuery,
		value,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.Wrap(store.ErrNotFound, errWrapMessage)
		}

		return nil, errors.Wrap(errors.Wrap(store.ErrUnknownSQL, err.Error()), errWrapMessage)
	}

	return tag, nil
}

func (tagRepository *TagRepository) FindAllByGame(game *model.Game) ([]*model.Tag, error) {
	repositoryName := "Tag"
	methodName := "FindAllByGame"
	errWrapMessage := fmt.Sprintf(store.ErrRepositoryMessageFormat, repositoryName, methodName)

	tags := []*model.Tag{}
	findQuery := "SELECT " +
		"tags.id AS id, " +
		"tags.name AS name " +

		"FROM tags " +

		"LEFT JOIN game_tags " +
		"ON (tags.id = game_tags.tag_id) " +

		"WHERE game_tags.game_id = $1;"

	if err := tagRepository.store.db.Select(
		&tags,
		findQuery,
		game.ID,
	); err != nil {
		if err == sql.ErrNoRows {
			return []*model.Tag{}, nil
		}

		return nil, errors.Wrap(errors.Wrap(store.ErrUnknownSQL, err.Error()), errWrapMessage)
	}

	return tags, nil
}

func (tagRepository *TagRepository) Update(newTag *model.Tag) error {
	repositoryName := "Tag"
	methodName := "Update"
	errWrapMessage := fmt.Sprintf(store.ErrRepositoryMessageFormat, repositoryName, methodName)

	updateQuery := "UPDATE tags " +
		"SET name = :name " +
		"WHERE id = :id;"

	countResult, err := tagRepository.store.db.NamedExec(
		updateQuery,
		newTag,
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

func (tagRepository *TagRepository) Delete(id uint64) error {
	repositoryName := "Tag"
	methodName := "Delete"
	errWrapMessage := fmt.Sprintf(store.ErrRepositoryMessageFormat, repositoryName, methodName)

	deleteQuery := "DELETE FROM tags WHERE id = $1;"

	countResult, err := tagRepository.store.db.Exec(
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
