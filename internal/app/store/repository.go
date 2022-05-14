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

type GameRepository interface {
	Create(*model.Game) error
	Find(uint64) (*model.Game, error)
	FindBy(string, interface{}) (*model.Game, error)
	FindAll() ([]*model.Game, error)
	FindAllByUser(*model.User) ([]*model.Game, error)
	FindAllByQueryTags(string, []*model.Tag) ([]*model.Game, error)
	Update(*model.Game) error
	Delete(uint64) error
}

type TagRepository interface {
	Create(*model.Tag) error
	Find(uint64) (*model.Tag, error)
	FindBy(string, interface{}) (*model.Tag, error)
	FindAllByGame(*model.Game) ([]*model.Tag, error)
	Update(*model.Tag) error
	Delete(uint64) error
}

type MarketRepository interface {
	Create(*model.Market) error
	Find(uint64) (*model.Market, error)
	FindBy(string, interface{}) (*model.Market, error)
	Update(*model.Market) error
	Delete(uint64) error
}

type UserGameFavouriteRepository interface {
	Create(*model.UserGameFavourite) error
	Find(uint64) (*model.UserGameFavourite, error)
	// FindBy(string, interface{}) (*model.UserGameFavourite, error)
	FindByUserGame(*model.User, *model.Game) (*model.UserGameFavourite, error)
	// Update(*model.UserGameFavourite) error
	Delete(uint64) error
}

type GameTagRepository interface {
	Create(*model.GameTag) error
	Find(uint64) (*model.GameTag, error)
	FindBy(string, interface{}) (*model.GameTag, error)
	Update(*model.GameTag) error
	Delete(uint64) error
}

type GameMarketPriceRepository interface {
	Create(*model.GameMarketPrice) error
	Find(uint64) (*model.GameMarketPrice, error)
	FindBy(string, interface{}) (*model.GameMarketPrice, error)
	FindByGameMarket(*model.Game, *model.Market) (*model.GameMarketPrice, error)
	FindAllByGame(*model.Game) ([]*model.GameMarketPrice, error)
	Update(*model.GameMarketPrice) error
	Delete(uint64) error
}

type MarketBlacklistItemRepository interface {
	Create(*model.MarketBlacklistItem) error
	CheckByURL(string) (bool, error)
	Delete(uint64) error
}
