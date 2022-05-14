package apistore

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/pkg/errors"
	"github.com/spolyakovs/price-hunter-ITMO/internal/app/model"
	"github.com/spolyakovs/price-hunter-ITMO/internal/app/store"
)

type APIGOG struct {
	store store.Store
}

func NewAPIGOG(st store.Store) *APIGOG {
	return &APIGOG{
		store: st,
	}
}

func expectedGOGURL(gameName string) string {
	var reOther = regexp.MustCompile(`[^a-z ]`)

	gameURL := strings.ToLower(gameName)
	gameURL = reOther.ReplaceAllString(gameURL, "")
	gameURL = strings.Replace(gameURL, " ", "-", -1)

	return gameURL
}

// TODO: write APIGOG.GetGames()
func (api *APIGOG) GetGames() error {

	type responseProductPrice struct {
		FinalValue      string `json:"finalAmount"`
		InitialValue    string `json:"baseAmount"`
		DiscountPercent int    `json:"discount"`
	}
	type responseProduct struct {
		ID    int                  `json:"id"`
		Title string               `json:"title"`
		URL   string               `json:"url"`
		Price responseProductPrice `json:"price"`
	}
	type response struct {
		Products []responseProduct `json:"products"`
	}

	apiName := "GOG"
	methodName := "GetGames"
	errWrapMessage := fmt.Sprintf(errAPIStoreMessageFormat, apiName, methodName)

	games, err := api.store.Games().FindAll()
	if err != nil {
		errWrapped := errors.Wrap(err, errWrapMessage)
		return errWrapped
	}

	counter := 0
	fmt.Println("Getting prices from GOG")

	for _, game := range games {
		url := fmt.Sprintf("https://embed.gog.com/games/ajax/filtered?search=%s&language=en", game.Name)

		url = strings.Replace(url, " ", "%20", -1)

		resp, err := http.Get(url)
		if err != nil {
			errWrapped := errors.Wrap(err, errWrapMessage)
			return errWrapped
		}

		responseStruct := &response{}

		if err := json.NewDecoder(resp.Body).Decode(responseStruct); err != nil {
			errWrapped := errors.Wrap(err, errWrapMessage)
			return errWrapped
		}

		marketGOG, err := api.store.Markets().FindBy("name", "GOG.com")
		if err != nil {
			errWrapped := errors.Wrap(err, errWrapMessage)

			return errWrapped
		}

		for _, gameDataRaw := range responseStruct.Products {
			if gameDataRaw.Title != game.Name {
				continue
			}

			priceFinalFormatted := fmt.Sprintf("%s руб.", gameDataRaw.Price.FinalValue)

			priceInitialFormatted := fmt.Sprintf("%s руб.", gameDataRaw.Price.InitialValue)
			if priceFinalFormatted == priceInitialFormatted {
				priceInitialFormatted = ""
			}

			splitURL := strings.Split(gameDataRaw.URL, "/")

			marketGameURL := splitURL[len(splitURL)-1]

			gameMarketPrice := &model.GameMarketPrice{
				InitialValueFormatted: priceInitialFormatted,
				FinalValueFormatted:   priceFinalFormatted,
				DiscountPercent:       gameDataRaw.Price.DiscountPercent,
				MarketGameURL:         marketGameURL,
				Game:                  game,
				Market:                marketGOG,
			}

			gameMarketPriceFound, err := api.store.GameMarketPrices().FindByGameMarket(game, marketGOG)
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

			counter += 1
			break
		}
	}

	fmt.Printf("Successfully got prices from GOG for all %d games\n", counter)

	return nil
}
