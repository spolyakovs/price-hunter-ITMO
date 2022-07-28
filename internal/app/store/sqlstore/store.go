package sqlstore

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/spolyakovs/price-hunter-ITMO/internal/app/store"
)

type Store struct {
	db                            *sqlx.DB
	userRepository                *UserRepository
	publisherRepository           *PublisherRepository
	gameRepository                *GameRepository
	tagRepository                 *TagRepository
	marketRepository              *MarketRepository
	userGameFavouriteRepository   *UserGameFavouriteRepository
	gameTagRepository             *GameTagRepository
	gameMarketPriceRepository     *GameMarketPriceRepository
	marketBlacklistItemRepository *MarketBlacklistItemRepository
}

func New(db *sqlx.DB) (*Store, error) {
	newStore := &Store{
		db: db,
	}

	if err := newStore.createTables(); err != nil {
		return nil, err
	}

	if err := newStore.insertDataMarkets(); err != nil {
		return nil, err
	}

	if err := newStore.InsertTestData(); err != nil {
		fmt.Printf("Couldn't insert test data into DB:\n\t%s", err.Error())
		return nil, err
	}

	return newStore, nil
}

func (st *Store) Users() store.UserRepository {
	if st.userRepository != nil {
		return st.userRepository
	}

	st.userRepository = &UserRepository{
		store: st,
	}

	return st.userRepository
}

func (st *Store) Publishers() store.PublisherRepository {
	if st.publisherRepository != nil {
		return st.publisherRepository
	}

	st.publisherRepository = &PublisherRepository{
		store: st,
	}

	return st.publisherRepository
}

func (st *Store) Games() store.GameRepository {
	if st.gameRepository != nil {
		return st.gameRepository
	}

	st.gameRepository = &GameRepository{
		store: st,
	}

	return st.gameRepository
}

func (st *Store) Tags() store.TagRepository {
	if st.tagRepository != nil {
		return st.tagRepository
	}

	st.tagRepository = &TagRepository{
		store: st,
	}

	return st.tagRepository
}

func (st *Store) Markets() store.MarketRepository {
	if st.marketRepository != nil {
		return st.marketRepository
	}

	st.marketRepository = &MarketRepository{
		store: st,
	}

	return st.marketRepository
}

func (st *Store) UserGameFavourites() store.UserGameFavouriteRepository {
	if st.userGameFavouriteRepository != nil {
		return st.userGameFavouriteRepository
	}

	st.userGameFavouriteRepository = &UserGameFavouriteRepository{
		store: st,
	}

	return st.userGameFavouriteRepository
}

func (st *Store) GameTags() store.GameTagRepository {
	if st.gameTagRepository != nil {
		return st.gameTagRepository
	}

	st.gameTagRepository = &GameTagRepository{
		store: st,
	}

	return st.gameTagRepository
}

func (st *Store) GameMarketPrices() store.GameMarketPriceRepository {
	if st.gameMarketPriceRepository != nil {
		return st.gameMarketPriceRepository
	}

	st.gameMarketPriceRepository = &GameMarketPriceRepository{
		store: st,
	}

	return st.gameMarketPriceRepository
}

func (st *Store) MarketBlacklist() store.MarketBlacklistItemRepository {
	if st.marketBlacklistItemRepository != nil {
		return st.marketBlacklistItemRepository
	}

	st.marketBlacklistItemRepository = &MarketBlacklistItemRepository{
		store: st,
	}

	return st.marketBlacklistItemRepository
}
