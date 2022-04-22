package sqlstore

import (
	"fmt"

	"github.com/spolyakovs/price-hunter-ITMO/internal/app/model"
)

var (
	users []*model.User
)

func (store *Store) fillTables() error {
	if err := store.fillTableUsers(); err != nil {
		return err
	}

	return nil
}

func (store *Store) fillTableUsers() error {
	for i := 1; i <= 5; i++ {
		users = append(users, &model.User{
			Username: fmt.Sprintf("test_username_%d", i),
			Email:    fmt.Sprintf("test_email_%d@example.org", i),
			Password: fmt.Sprintf("test_password_%d", i),
		})
	}

	for _, user := range users {
		if err := store.Users().Create(user); err != nil {
			return err
		}
	}

	return nil
}
