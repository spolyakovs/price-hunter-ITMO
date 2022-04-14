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
	if err := user.Validate(); err != nil {
		return err
	}

	if err := user.BeforeCreate(); err != nil {
		return err
	}

	createQuery := "INSERT INTO users (username, email, encrypted_password) VALUES ($1, $2, $3) RETURNING id;"

	err := userRepository.store.db.Get(
		&user.ID,
		createQuery,
		user.Username, user.Email, user.EncryptedPassword,
	)

	return errors.Wrap(store.ErrCreate, err.Error())
}

func (userRepository *UserRepository) Find(id uint64) (*model.User, error) {
	return userRepository.FindBy("id", id)
}

func (userRepository *UserRepository) FindBy(columnName string, value interface{}) (*model.User, error) {
	user := &model.User{}

	findQuery := fmt.Sprintf("SELECT * FROM users WHERE %s = $1 LIMIT 1;", columnName)
	if err := userRepository.store.db.Get(
		user,
		findQuery,
		value,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrNotFound
		}

		return nil, errors.Wrap(store.ErrUnknownSQL, err.Error())
	}

	return user, nil
}

func (userRepository *UserRepository) UpdateEmail(newEmail string, userId uint64) error {
	if err := validation.Validate(&newEmail, model.ValidationRulesEmail...); err != nil {
		return err
	}

	updateEmailQuery := "UPDATE users " +
		`SET email = $1 ` +
		`WHERE id = $2;`

	countResult, countResultErr := userRepository.store.db.Exec(
		updateEmailQuery,
		newEmail,
		userId,
	)

	if countResultErr != nil {
		return errors.Wrap(store.ErrUnknownSQL, countResultErr.Error())
	}

	count, countErr := countResult.RowsAffected()

	if countErr != nil {
		return errors.Wrap(store.ErrUnknownSQL, countErr.Error())
	}

	if count == 0 {
		return store.ErrNotFound
	}

	return nil
}

func (userRepository *UserRepository) UpdatePassword(newPassword string, userId uint64) error {
	if err := validation.Validate(&newPassword, model.ValidationRulesPassword...); err != nil {
		return nil
	}

	updatePasswordQuery := "UPDATE users " +
		`SET encrypted_password = $1 ` +
		`WHERE id = $2;`
	countResult, countResultErr := userRepository.store.db.Exec(
		updatePasswordQuery,
		newPassword,
		userId,
	)

	if countResultErr != nil {
		return errors.Wrap(store.ErrUnknownSQL, countResultErr.Error())
	}

	count, countErr := countResult.RowsAffected()

	if countErr != nil {
		return errors.Wrap(store.ErrUnknownSQL, countErr.Error())
	}

	if count == 0 {
		return store.ErrNotFound
	}

	return nil
}

func (userRepository *UserRepository) Delete(id uint64) error {
	deleteQuery := "DELETE FROM users WHERE id = $1;"

	countResult, countResultErr := userRepository.store.db.Exec(
		deleteQuery,
		id,
	)

	if countResultErr != nil {
		return errors.Wrap(store.ErrUnknownSQL, countResultErr.Error())
	}

	count, countErr := countResult.RowsAffected()

	if countErr != nil {
		return errors.Wrap(store.ErrUnknownSQL, countErr.Error())
	}

	if count == 0 {
		return store.ErrNotFound
	}

	return nil
}
