package sqlstore

import (
	"database/sql"
	"fmt"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/pkg/errors"

	"github.com/spolyakovs/price-hunter-ITMO/internal/app/model"
	"github.com/spolyakovs/price-hunter-ITMO/internal/app/store"
)

type UserRepository struct {
	store *Store
}

func (userRepository *UserRepository) Create(user *model.User) error {
	repositoryName := "User"
	methodName := "Create"
	errWrapMessage := fmt.Sprintf(store.ErrRepositoryMessageFormat, repositoryName, methodName)

	if err := user.Validate(); err != nil {
		return errors.Wrap(err, errWrapMessage)
	}

	if err := user.BeforeCreate(); err != nil {
		return errors.Wrap(err, errWrapMessage)
	}

	createQuery := "INSERT INTO users (username, email, encrypted_password) VALUES ($1, $2, $3) " +
		"ON CONFLICT(username) DO UPDATE SET username = EXCLUDED.username RETURNING id;"

	if err := userRepository.store.db.Get(
		&user.ID,
		createQuery,
		user.Username, user.Email, user.EncryptedPassword,
	); err != nil {
		return errors.Wrap(errors.Wrap(store.ErrUnknownSQL, err.Error()), errWrapMessage)
	}

	return nil
}

func (userRepository *UserRepository) Find(id uint64) (*model.User, error) {
	return userRepository.FindBy("id", id)
}

func (userRepository *UserRepository) FindBy(columnName string, value interface{}) (*model.User, error) {
	repositoryName := "User"
	methodName := "Find"
	errWrapMessage := fmt.Sprintf(store.ErrRepositoryMessageFormat, repositoryName, methodName)

	user := &model.User{}
	findQuery := fmt.Sprintf("SELECT * FROM users WHERE %s = $1 LIMIT 1;", columnName)

	if err := userRepository.store.db.Get(
		user,
		findQuery,
		value,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.Wrap(store.ErrNotFound, errWrapMessage)
		}

		return nil, errors.Wrap(errors.Wrap(store.ErrUnknownSQL, err.Error()), errWrapMessage)
	}

	return user, nil
}

func (userRepository *UserRepository) UpdateEmail(newEmail string, userId uint64) error {
	repositoryName := "User"
	methodName := "UpdateEmail"
	errWrapMessage := fmt.Sprintf(store.ErrRepositoryMessageFormat, repositoryName, methodName)

	if err := validation.Validate(&newEmail, model.ValidationRulesEmail...); err != nil {
		return errors.Wrap(errors.Wrap(model.ErrValidationFailed, err.Error()), errWrapMessage)
	}

	updateEmailQuery := "UPDATE users " +
		"SET email = $1 " +
		"WHERE id = $2;"

	countResult, err := userRepository.store.db.Exec(
		updateEmailQuery,
		newEmail,
		userId,
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

func (userRepository *UserRepository) UpdatePassword(newPassword string, userId uint64) error {
	repositoryName := "User"
	methodName := "UpdatePassword"
	errWrapMessage := fmt.Sprintf(store.ErrRepositoryMessageFormat, repositoryName, methodName)

	if err := validation.Validate(&newPassword, model.ValidationRulesPassword...); err != nil {
		return errors.Wrap(errors.Wrap(model.ErrValidationFailed, err.Error()), errWrapMessage)
	}

	newPasswordEncrypted, err := model.EncryptString(newPassword)
	if err != nil {
		return errors.Wrap(err, errWrapMessage)
	}

	updatePasswordQuery := "UPDATE users " +
		"SET encrypted_password = $1 " +
		"WHERE id = $2;"
	countResult, err := userRepository.store.db.Exec(
		updatePasswordQuery,
		newPasswordEncrypted,
		userId,
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

func (userRepository *UserRepository) Delete(id uint64) error {
	repositoryName := "User"
	methodName := "Delete"
	errWrapMessage := fmt.Sprintf(store.ErrRepositoryMessageFormat, repositoryName, methodName)

	deleteQuery := "DELETE FROM users WHERE id = $1;"

	countResult, err := userRepository.store.db.Exec(
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
