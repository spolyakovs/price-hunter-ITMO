package apiStore

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/pkg/errors"
	"github.com/spolyakovs/price-hunter-ITMO/internal/app/model"
	"github.com/spolyakovs/price-hunter-ITMO/internal/app/store"
)

type APIEpicGames struct {
	store store.Store
}

func NewAPIEpicGames(st store.Store) *APIEpicGames {
	return &APIEpicGames{
		store: st,
	}
}

func expectedEpicGamesURL(gameName string) string {
	gameURL := strings.ToLower(gameName)
	gameURL = strings.Replace(gameURL, " ", "-", -1)
	return gameURL
}

// TODO: write APIEpicGames.GetGames()
func (api *APIEpicGames) GetGames() error {
	type responseDataCatalogStoreItemPriceTotal struct {
		FinalValue      int `json:"discountPrice"`
		InitialValue    int `json:"originalPrice"`
		DiscountPercent int `json:"discount"`
	}
	type responseDataCatalogStoreItemPrice struct {
		TotalPrice responseDataCatalogStoreItemPriceTotal `json:"totalPrice"`
	}
	type responseDataCatalogStoreItem struct {
		ProductSlug string                            `json:"productSlug,omitempty"`
		Title       string                            `json:"title"`
		Price       responseDataCatalogStoreItemPrice `json:"price"`
	}
	type responseDataCatalogStore struct {
		Elements []responseDataCatalogStoreItem `json:"elements"`
	}
	type responseDataCatalog struct {
		Store responseDataCatalogStore `json:"searchStore"`
	}
	type responseData struct {
		Catalog responseDataCatalog `json:"Catalog"`
	}
	type response struct {
		Data responseData `json:"data"`
	}

	apiName := "EpicGames"
	methodName := "GetGames"
	errWrapMessage := fmt.Sprintf(errAPIStoreMessageFormat, apiName, methodName)

	games, err := api.store.Games().FindAll()
	if err != nil {
		errWrapped := errors.Wrap(err, errWrapMessage)
		return errWrapped
	}

	for _, game := range games {
		url := fmt.Sprintf("https://www.epicgames.com/graphql?query="+
			"{Catalog {searchStore(keywords: \"%s\", country: \"RU\", locale: \"US\", count: 1) "+
			"{elements {"+
			"id productSlug namespace title description price(country: \"RU\") "+
			"{totalPrice{discountPrice originalPrice discount } } } } } }", expectedEpicGamesURL(game.Name))

		resp, err := http.Get(url)
		if err != nil {
			errWrapped := errors.Wrap(err, errWrapMessage)
			return errWrapped
		}

		fmt.Println(expectedEpicGamesURL(game.Name))

		responseStruct := &response{}

		if err := json.NewDecoder(resp.Body).Decode(responseStruct); err != nil {
			errWrapped := errors.Wrap(err, errWrapMessage)
			return errWrapped
		}

		gameDataRaw := responseStruct.Data.Catalog.Store.Elements[0]

		if gameDataRaw.Title != game.Name {
			continue
		}

		marketEpicGames, err := api.store.Markets().FindBy("name", "EpicGamesStore")
		if err != nil {
			errWrapped := errors.Wrap(err, errWrapMessage)

			return errWrapped
		}

		priceInitialFormatted := fmt.Sprintf("%d руб.", gameDataRaw.Price.TotalPrice.InitialValue/100)
		priceFinalFormatted := fmt.Sprintf("%d руб.", gameDataRaw.Price.TotalPrice.FinalValue/100)

		gameMarketPrice := &model.GameMarketPrice{
			InitialValueFormatted: priceInitialFormatted,
			FinalValueFormatted:   priceFinalFormatted,
			DiscountPercent:       gameDataRaw.Price.TotalPrice.DiscountPercent,
			MarketGameURL:         gameDataRaw.ProductSlug,
			Game:                  game,
			Market:                marketEpicGames,
		}

		gameMarketPriceFound, err := api.store.GameMarketPrices().FindByGameMarket(game, marketEpicGames)
		if err != nil {
			if errors.Cause(err) != store.ErrNotFound {
				errWrapped := errors.Wrap(err, errWrapMessage)
				return errWrapped
			}

			if err := api.store.GameMarketPrices().Create(gameMarketPrice); err != nil {
				errWrapped := errors.Wrap(err, errWrapMessage)
				return errWrapped
			}
		} else {
			gameMarketPrice.ID = gameMarketPriceFound.ID

			if err := api.store.GameMarketPrices().Update(gameMarketPrice); err != nil {
				errWrapped := errors.Wrap(err, errWrapMessage)
				return errWrapped
			}
		}
	}

	// TODO: example URL in Postman, need to choose correct one (example with "The Long Dark"), check what happens when game isn't in Store
	return nil
}
