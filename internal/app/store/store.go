package store

type Store interface {
	Users() UserRepository
	Publishers() PublisherRepository
	Games() GameRepository
	Tags() TagRepository
	Markets() MarketRepository
	UserGameFavourites() UserGameFavouriteRepository
	GameTags() GameTagRepository
	GameMarketPrices() GameMarketPriceRepository
	MarketBlacklist() MarketBlacklistItemRepository
}
