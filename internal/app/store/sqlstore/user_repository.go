package sqlstore

import (
	"database/sql"
	"fmt"

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

	return userRepository.store.db.Get(
		&user.ID,
		createQuery,
		user.Username, user.Email, user.EncryptedPassword,
	)
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
			return nil, store.ErrRecordNotFound
		}

		return nil, err
	}

	return user, nil
}

func (userRepository *UserRepository) Update(user *model.User) error {
	if err := user.Validate(); err != nil {
		return err
	}

	if err := user.BeforeCreate(); err != nil {
		return err
	}

	updateQuery := "UPDATE users " +
		`SET username = :username, ` +
		`email = :email, ` +
		`encrypted_password = :encrypted_password ` +
		`WHERE id = :id;`
	countResult, countResultErr := userRepository.store.db.NamedExec(
		updateQuery,
		user,
	)

	if countResultErr != nil {
		return countResultErr
	}

	count, countErr := countResult.RowsAffected()

	if countErr != nil {
		return countErr
	}

	if count == 0 {
		return store.ErrRecordNotFound
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
		return countResultErr
	}

	count, countErr := countResult.RowsAffected()

	if countErr != nil {
		return countErr
	}

	if count == 0 {
		return store.ErrRecordNotFound
	}

	return nil
}
