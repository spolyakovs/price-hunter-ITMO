package apiStore

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/spolyakovs/price-hunter-ITMO/internal/app/model"
	"github.com/spolyakovs/price-hunter-ITMO/internal/app/store"
)

// TODO: write APIGog

type APISteam struct {
	apiKey string
	store  store.Store
}

func NewAPISteam(apiKey string, st store.Store) *APISteam {
	return &APISteam{
		apiKey: apiKey,
		store:  st,
	}
}

func (api *APISteam) GetGames() error {
	type responseListItem struct {
		AppID int    `json:"appid"`
		Name  string `json:"name"`
	}

	type responseList struct {
		Apps []responseListItem `json:"apps"`
	}

	type response struct {
		AppList responseList `json:"applist"`
	}

	apiName := "Steam"
	methodName := "GetGames"
	errWrapMessage := fmt.Sprintf(errAPIStoreMessageFormat, apiName, methodName)

	url := fmt.Sprintf("http://api.steampowered.com/ISteamApps/GetAppList/v2/?key=%s&format=json", api.apiKey)

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
	defer resp.Body.Close()

	appIDs := []string{}

	gamesToUpdate := make(map[string]*model.Game)

	maxAppsToLoad := 500

	if games, err := api.store.Games().FindAll(); err != nil {
		errWrapped := errors.Wrap(err, errWrapMessage)
		return errWrapped
	} else if len(games) >= maxAppsToLoad {
		fmt.Printf("Already downloaded %d games\n", len(games))
		return nil
	}

	for _, app := range responseStruct.AppList.Apps {
		gameNameClean := cleanGameName(app.Name)
		appID := strconv.Itoa(app.AppID)

		if checkGameName(gameNameClean) {
			blacklisted, err := api.store.MarketBlacklist().CheckByURL(appID)

			if err != nil {
				errWrapped := errors.Wrap(err, errWrapMessage)
				return errWrapped
			}

			if blacklisted {
				continue
			}

			gameFound, err := api.store.Games().FindBy("name", gameNameClean)

			if err != nil {
				if errors.Cause(err) != store.ErrNotFound {
					errWrapped := errors.Wrap(err, errWrapMessage)
					return errWrapped
				}

				appIDs = append(appIDs, appID)
			} else {
				gamesToUpdate[appID] = gameFound
			}
		}
	}

	marketSteam, err := api.store.Markets().FindBy("name", "Steam")
	if err != nil {
		errWrapped := errors.Wrap(err, errWrapMessage)
		return errWrapped
	}

	if err := api.UpdateGameMarketPrices(gamesToUpdate, marketSteam); err != nil {
		errWrapped := errors.Wrap(err, errWrapMessage)
		return errWrapped
	}

	maxGameCount := 2000
	if maxGameCount > len(appIDs) {
		maxGameCount = len(appIDs)
	}
	counter := 0
	fmt.Println("Getting GameInfo from Steam")

	for _, appID := range appIDs[:maxGameCount] {
		if err := api.getSteamGameInfo(appID, marketSteam); err != nil {
			errWrapped := errors.Wrap(err, errWrapMessage)
			return errWrapped
		}
		fmt.Printf("%d/%d\r", counter, maxGameCount)
		counter += 1
	}

	fmt.Printf("Successfully got GameInfo from Steam for all %d games\n", maxGameCount)

	return nil
}

func (api *APISteam) UpdateGameMarketPrices(gamesToUpdate map[string]*model.Game, marketSteam *model.Market) error {
	type responseAppDataPrice struct {
		InitialFormatted string `json:"initial_formatted,omitempty"`
		FinalFormatted   string `json:"final_formatted,omitempty"`
		DiscountPercent  int    `json:"discount_percent,omitempty"`
	}

	type responseAppData struct {
		PriceOverview responseAppDataPrice `json:"price_overview,omitempty"`
	}

	type responseApp struct {
		Success bool            `json:"success"`
		Data    responseAppData `json:"data,omitempty"`
	}

	apiName := "Steam"
	methodName := "UpdateGameMarketPrices"
	errWrapMessage := fmt.Sprintf(errAPIStoreMessageFormat, apiName, methodName)

	appIDs := make([]string, 0, len(gamesToUpdate))
	for appID := range gamesToUpdate {
		appIDs = append(appIDs, appID)
	}

	maxGamesCount := 100
	offset := 0

	counter := 0
	fmt.Println("UpdatingPrices from Steam")

	for {
		var currentAppIDs []string

		if len(appIDs)-offset > maxGamesCount {
			currentAppIDs = appIDs[offset : offset+maxGamesCount]
			offset += maxGamesCount
		} else if len(appIDs)-offset > 0 {
			currentAppIDs = appIDs[offset:]
			offset = len(appIDs)
		} else {
			break
		}

		url := fmt.Sprintf("http://store.steampowered.com/api/appdetails?appids=%s&finters=price_overview&cc=ru&l=en", strings.Join(currentAppIDs, ","))

		resp, err := http.Get(url)
		if err != nil {
			errWrapped := errors.Wrap(err, errWrapMessage)
			return errWrapped
		}

		responseStruct := make(map[string]responseApp)

		if err := json.NewDecoder(resp.Body).Decode(&responseStruct); err != nil {
			errWrapped := errors.Wrap(err, errWrapMessage)
			return errWrapped
		}
		defer resp.Body.Close()

		for _, appID := range currentAppIDs {
			gameInfoRaw := responseStruct[appID]

			gameMarketPrice := &model.GameMarketPrice{
				InitialValueFormatted: gameInfoRaw.Data.PriceOverview.InitialFormatted,
				FinalValueFormatted:   gameInfoRaw.Data.PriceOverview.FinalFormatted,
				DiscountPercent:       gameInfoRaw.Data.PriceOverview.DiscountPercent,
				MarketGameURL:         appID,
				Game:                  gamesToUpdate[appID],
				Market:                marketSteam,
			}

			gameMarketPriceFound, err := api.store.GameMarketPrices().FindByGameMarket(gamesToUpdate[appID], marketSteam)
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

			fmt.Printf("%d/%d\r", counter, len(gamesToUpdate))
			counter += 1
		}
	}

	fmt.Printf("Successfully updated prices from Steam for all %d games\n", len(gamesToUpdate))

	return nil
}

func (api *APISteam) getSteamGameInfo(appID string, marketSteam *model.Market) error {
	type responseAppDataGenre struct {
		Description string `json:"description"`
	}

	type responseAppDataReleaseDate struct {
		ComingSoon bool   `json:"coming_soon"`
		Date       string `json:"date"`
	}

	type responseAppDataPrice struct {
		InitialFormatted string `json:"initial_formatted,omitempty"`
		FinalFormatted   string `json:"final_formatted,omitempty"`
		DiscountPercent  int    `json:"discount_percent,omitempty"`
	}

	type responseAppData struct {
		Type          string                     `json:"type"`
		Name          string                     `json:"name"`
		HeaderImage   string                     `json:"header_image"`
		Genres        []responseAppDataGenre     `json:"genres"`
		ReleaseDate   responseAppDataReleaseDate `json:"release_date"`
		Description   string                     `json:"short_description"`
		Publishers    []string                   `jsin:"publishers"`
		PriceOverview responseAppDataPrice       `json:"price_overview,omitempty"`
	}

	type responseApp struct {
		Success bool            `json:"success"`
		Data    responseAppData `json:"data,omitempty"`
	}

	apiName := "Steam"
	methodName := "getGamesInfo"
	errWrapMessage := fmt.Sprintf(errAPIStoreMessageFormat, apiName, methodName)

	url := fmt.Sprintf("http://store.steampowered.com/api/appdetails?appids=%s&cc=ru&l=en", appID)

	resp, err := http.Get(url)
	if err != nil {
		errWrapped := errors.Wrap(err, errWrapMessage)
		errWrapped = errors.Wrap(err, fmt.Sprintf("AppID: %s", appID))
		return errWrapped
	}

	responseStruct := make(map[string]responseApp)

	if err := json.NewDecoder(resp.Body).Decode(&responseStruct); err != nil {
		errWrapped := errors.Wrap(err, errWrapMessage)
		errWrapped = errors.Wrap(err, fmt.Sprintf("AppID: %s", appID))
		return errWrapped
	}
	defer resp.Body.Close()

	gameInfoRaw := responseStruct[appID]
	marketBlacklist := &model.MarketBlacklistItem{
		MarketGameURL: appID,
		Market:        marketSteam,
	}

	if !gameInfoRaw.Success {
		if err := api.store.MarketBlacklist().Create(marketBlacklist); err != nil {
			errWrapped := errors.Wrap(err, errWrapMessage)
			errWrapped = errors.Wrap(err, fmt.Sprintf("AppID: %s", appID))
			return errWrapped
		}
		return nil
	}

	if gameInfoRaw.Data.Type != "game" {
		if err := api.store.MarketBlacklist().Create(marketBlacklist); err != nil {
			errWrapped := errors.Wrap(err, errWrapMessage)
			errWrapped = errors.Wrap(err, fmt.Sprintf("AppID: %s", appID))
			return errWrapped
		}
		return nil
	}

	if gameInfoRaw.Data.ReleaseDate.ComingSoon {
		if err := api.store.MarketBlacklist().Create(marketBlacklist); err != nil {
			errWrapped := errors.Wrap(err, errWrapMessage)
			errWrapped = errors.Wrap(err, fmt.Sprintf("AppID: %s", appID))
			return errWrapped
		}
		return nil
	}

	inputDateLayout := "2 Jan, 2006"
	outputDateLayout := "02.01.2006"

	releaseDateClean, errParse := time.Parse(inputDateLayout, gameInfoRaw.Data.ReleaseDate.Date)
	if errParse != nil {
		if err := api.store.MarketBlacklist().Create(marketBlacklist); err != nil {
			errWrapped := errors.Wrap(err, errWrapMessage)
			errWrapped = errors.Wrap(err, fmt.Sprintf("AppID: %s", appID))
			return errWrapped
		}
		return nil
	}

	publisher := &model.Publisher{
		Name: gameInfoRaw.Data.Publishers[0],
	}

	if publisherFound, err := api.store.Publishers().FindBy("name", publisher.Name); err != nil {
		if errors.Cause(err) != store.ErrNotFound {
			errWrapped := errors.Wrap(err, errWrapMessage)
			return errWrapped
		}

		if err := api.store.Publishers().Create(publisher); err != nil {
			errWrapped := errors.Wrap(err, errWrapMessage)
			return errWrapped
		}
	} else {
		publisher = publisherFound
	}

	for _, tag := range gameInfoRaw.Data.Genres {
		tagNameClean := cleanTagName(tag.Description)

		if _, err := api.store.Tags().FindBy("name", tagNameClean); err == nil {
			continue
		} else if errors.Cause(err) != store.ErrNotFound {
			errWrapped := errors.Wrap(err, errWrapMessage)
			return errWrapped
		}

		newTag := &model.Tag{
			Name: tagNameClean,
		}

		if err := api.store.Tags().Create(newTag); err != nil {
			errWrapped := errors.Wrap(err, errWrapMessage)
			return errWrapped
		}
	}

	game := &model.Game{
		HeaderImageURL: gameInfoRaw.Data.HeaderImage,
		Name:           cleanGameName(gameInfoRaw.Data.Name),
		Description:    gameInfoRaw.Data.Description,
		ReleaseDate:    releaseDateClean.Format(outputDateLayout),
		Publisher:      publisher,
	}

	if err := api.store.Games().Create(game); err != nil {
		errWrapped := errors.Wrap(err, errWrapMessage)
		return errWrapped
	}

	gameMarketPrice := &model.GameMarketPrice{
		InitialValueFormatted: gameInfoRaw.Data.PriceOverview.InitialFormatted,
		FinalValueFormatted:   gameInfoRaw.Data.PriceOverview.FinalFormatted,
		DiscountPercent:       gameInfoRaw.Data.PriceOverview.DiscountPercent,
		MarketGameURL:         appID,
		Game:                  game,
		Market:                marketSteam,
	}

	gameMarketPriceFound, err := api.store.GameMarketPrices().FindByGameMarket(game, marketSteam)
	if err != nil {
		if errors.Cause(err) != store.ErrNotFound {
			errWrapped := errors.Wrap(err, errWrapMessage)
			return errWrapped
		}

		if err := api.store.GameMarketPrices().Create(gameMarketPrice); err != nil {
			errWrapped := errors.Wrap(err, errWrapMessage)
			return errWrapped
		}

		return nil
	}

	gameMarketPrice.ID = gameMarketPriceFound.ID

	if err := api.store.GameMarketPrices().Update(gameMarketPrice); err != nil {
		errWrapped := errors.Wrap(err, errWrapMessage)
		return errWrapped
	}

	return nil
}
