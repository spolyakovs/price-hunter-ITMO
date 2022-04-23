package store

// TODO: change tables structure for Price Hunter app

type Store interface {
	Users() UserRepository
	Publishers() PublisherRepository
	Games() GameRepository
	Tags() TagRepository
	Markets() MarketRepository
	UserGameFavourites() UserGameFavouriteRepository
	GameTags() GameTagRepository
	GameMarketPrices() GameMarketPriceRepository
}
