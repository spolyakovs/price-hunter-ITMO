package store

import (
	"github.com/spolyakovs/price-hunter-ITMO/internal/app/model"
)

type UserRepository interface {
	Create(*model.User) error
	Find(uint64) (*model.User, error)
	FindBy(string, interface{}) (*model.User, error)
	UpdateEmail(string, uint64) error
	UpdatePassword(string, uint64) error
	Delete(uint64) error
}

type PublisherRepository interface {
	Create(*model.Publisher) error
	Find(uint64) (*model.Publisher, error)
	FindBy(string, interface{}) (*model.Publisher, error)
	Update(*model.Publisher) error
	Delete(uint64) error
}
