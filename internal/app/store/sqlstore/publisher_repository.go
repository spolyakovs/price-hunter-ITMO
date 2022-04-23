package sqlstore

import (
	"database/sql"
	"fmt"

	"github.com/pkg/errors"

	"github.com/spolyakovs/price-hunter-ITMO/internal/app/model"
	"github.com/spolyakovs/price-hunter-ITMO/internal/app/store"
)

type PublisherRepository struct {
	store *Store
}

func (publisherRepository *PublisherRepository) Create(publisher *model.Publisher) error {
	repositoryName := "Publisher"
	methodName := "Create"
	errWrapMessage := fmt.Sprintf(store.ErrRepositoryMessageFormat, repositoryName, methodName)

	createQuery := "INSERT INTO publishers (name) VALUES ($1) RETURNING id;"

	if err := publisherRepository.store.db.Get(
		&publisher.ID,
		createQuery,
		publisher.Name,
	); err != nil {
		return errors.Wrap(errors.Wrap(store.ErrUnknownSQL, err.Error()), errWrapMessage)
	}

	return nil
}

func (publisherRepository *PublisherRepository) Find(id uint64) (*model.Publisher, error) {
	return publisherRepository.FindBy("id", id)
}

func (publisherRepository *PublisherRepository) FindBy(columnName string, value interface{}) (*model.Publisher, error) {
	repositoryName := "Publisher"
	methodName := "FindBy"
	errWrapMessage := fmt.Sprintf(store.ErrRepositoryMessageFormat, repositoryName, methodName)

	publisher := &model.Publisher{}
	findQuery := fmt.Sprintf("SELECT * FROM publishers WHERE %s = $1 LIMIT 1;", columnName)

	if err := publisherRepository.store.db.Get(
		publisher,
		findQuery,
		value,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.Wrap(store.ErrNotFound, errWrapMessage)
		}

		return nil, errors.Wrap(errors.Wrap(store.ErrUnknownSQL, err.Error()), errWrapMessage)
	}

	return publisher, nil
}

func (publisherRepository *PublisherRepository) Update(newPublisher *model.Publisher) error {
	repositoryName := "Publisher"
	methodName := "Update"
	errWrapMessage := fmt.Sprintf(store.ErrRepositoryMessageFormat, repositoryName, methodName)

	updateQuery := "UPDATE publishers " +
		"SET name = :name " +
		"WHERE id = :id;"

	countResult, err := publisherRepository.store.db.NamedExec(
		updateQuery,
		newPublisher,
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

func (publisherRepository *PublisherRepository) Delete(id uint64) error {
	repositoryName := "Publisher"
	methodName := "Delete"
	errWrapMessage := fmt.Sprintf(store.ErrRepositoryMessageFormat, repositoryName, methodName)

	deleteQuery := "DELETE FROM publishers WHERE id = $1;"

	countResult, err := publisherRepository.store.db.Exec(
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
