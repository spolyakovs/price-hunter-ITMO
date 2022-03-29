package store

import (
	"github.com/spolyakovs/price-hunter-ITMO/internal/app/model"
)

type UserRepository interface {
	Create(*model.User) error
	Find(uint64) (*model.User, error)
	FindBy(string, interface{}) (*model.User, error)
	Update(*model.User) error
	Delete(uint64) error
}

type TeamRepository interface {
	Create(*model.Team) error
	Find(int) (*model.Team, error)
	FindBy(string, interface{}) (*model.Team, error)
	FindAllBy(string, interface{}) ([]*model.Team, error)
	GetAll() ([]*model.Team, error)
	Update(*model.Team) error
	Delete(int) error
}

type DriverRepository interface {
	Create(*model.Driver) error
	Find(int) (*model.Driver, error)
	FindBy(string, interface{}) (*model.Driver, error)
	FindAllBy(string, interface{}) ([]*model.Driver, error)
	GetAll() ([]*model.Driver, error)
	Update(*model.Driver) error
	Delete(int) error
}

type RaceRepository interface {
	Create(*model.Race) error
	Find(int) (*model.Race, error)
	FindBy(string, interface{}, string) (*model.Race, error)
	FindAllBy(string, interface{}, string) ([]*model.Race, error)
	GetAll() ([]*model.Race, error)
	Update(*model.Race) error
	Delete(int) error
}

type TeamDriverContractRepository interface {
	Create(*model.TeamDriverContract) error
	Find(int) (*model.TeamDriverContract, error)
	FindBy(string, interface{}, string) (*model.TeamDriverContract, error)
	FindAllBy(string, interface{}, string) ([]*model.TeamDriverContract, error)
	FindCurrentBy(string, interface{}) (*model.TeamDriverContract, error)
	FindAllCurrentBy(string, interface{}) ([]*model.TeamDriverContract, error)
	Update(*model.TeamDriverContract) error
	Delete(int) error
}

type RaceResultRepository interface {
	Create(*model.RaceResult) error
	Find(int) (*model.RaceResult, error)
	FindBy(string, interface{}) (*model.RaceResult, error)
	FindAllBy(string, interface{}) ([]*model.RaceResult, error)
	Update(*model.RaceResult) error
	Delete(int) error
}
