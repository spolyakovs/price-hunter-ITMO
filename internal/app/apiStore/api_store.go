package apiStore

import "github.com/spolyakovs/price-hunter-ITMO/internal/app/model"

type APIStore interface {
	GetGamePrices() []*model.GameMarketPrice
}
