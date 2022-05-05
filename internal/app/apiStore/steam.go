package apiStore

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/pkg/errors"
	"github.com/spolyakovs/price-hunter-ITMO/internal/app/model"
)

// TODO: write APIEpicGames
// TODO: write APIGog

// TODO: ONLY english letters in game name
// TODO: may be restrincting number of games in games_list request (and add offset param)

type APISteam struct {
	apiKey string
}

func NewAPISteam(apiKey string) *APISteam {
	return &APISteam{
		apiKey: apiKey,
	}
}

func (api *APISteam) GetGames() ([]*model.Game, error) {
	type gamesListResponseListItem struct {
		AppID int `json:"appid"`
	}

	type gamesListResponseList struct {
		Apps []gamesListResponseListItem `json:"apps"`
	}

	type gamesListResponse struct {
		AppList gamesListResponseList `json:"applist"`
	}

	apiName := "Steam"
	methodName := "GetGames"
	errWrapMessage := fmt.Sprintf(errAPIStoreMessageFormat, apiName, methodName)

	// idsListURL := fmt.Sprintf("http://api.steampowered.com/ISteamApps/GetAppList/v2/?key=%s&format=json", api.apiKey)
	//
	// resp, err := http.Get(idsListURL)
	// if err != nil {
	// 	errWrapped := errors.Wrap(err, errWrapMessage)
	// 	return nil, errWrapped
	// }
	//
	// gamesListResponseStruct := &gamesListResponse{}
	//
	// if err := json.NewDecoder(resp.Body).Decode(gamesListResponseStruct); err != nil {
	// 	errWrapped := errors.Wrap(err, errWrapMessage)
	// 	return nil, errWrapped
	// }

	appIDs := []string{"730"}

	// for _, app := range gamesListResponseStruct.AppList.Apps {
	// 	appIDs = append(appIDs, strconv.Itoa(app.AppID))
	// }

	games := []model.Game{}

	for _, appID := range appIDs[:1] {
		game, err := getSteamGamesInfo(appID)
		if err != nil {
			if errors.Cause(err) == ErrAppSkiped {
				continue
			}
			errWrapped := errors.Wrap(err, errWrapMessage)
			return nil, errWrapped
		}

		games = append(games, *game)
	}

	fmt.Printf("%+v\n", games)

	return nil, nil
}

func getSteamGamesInfo(appID string) (*model.Game, error) {
	type gameInfoResponseAppDataGenre struct {
		Description string `json:"description"`
	}

	type gameInfoResponseAppDataReleaseDate struct {
		ComingSoon bool   `json:"coming_soon"`
		Date       string `json:"date"`
	}

	type gameInfoResponseAppData struct {
		Type        string                             `json:"type"`
		Name        string                             `json:"name"`
		HeaderImage string                             `json:"header_image"`
		Genres      []gameInfoResponseAppDataGenre     `json:"genres"`
		ReleaseDate gameInfoResponseAppDataReleaseDate `json:"release_date"`
		Description string                             `json:"short_description"`
		Publishers  []string                           `jsin:"publishers"`
	}

	type gameInfoResponseApp struct {
		Success bool                    `json:"success"`
		Data    gameInfoResponseAppData `json:"data,omitempty"`
	}

	apiName := "Steam"
	methodName := "getGamesInfo"
	errWrapMessage := fmt.Sprintf(errAPIStoreMessageFormat, apiName, methodName)

	gameInfoURL := fmt.Sprintf("http://store.steampowered.com/api/appdetails?appids=%s&cc=ru&l=en", appID)

	resp, err := http.Get(gameInfoURL)
	if err != nil {
		errWrapped := errors.Wrap(err, errWrapMessage)
		errWrapped = errors.Wrap(err, fmt.Sprintf("AppID: %s", appID))
		return nil, errWrapped
	}

	gameInfoResponseStruct := make(map[string]gameInfoResponseApp)

	if err := json.NewDecoder(resp.Body).Decode(&gameInfoResponseStruct); err != nil {
		errWrapped := errors.Wrap(err, errWrapMessage)
		errWrapped = errors.Wrap(err, fmt.Sprintf("AppID: %s", appID))
		return nil, errWrapped
	}

	gameInfoRaw := gameInfoResponseStruct[appID]

	if !gameInfoRaw.Success {
		errWrapped := errors.Wrap(err, errWrapMessage)
		errWrapped = errors.Wrap(ErrAppSkiped, fmt.Sprintf("AppID: %s", appID))
		return nil, errWrapped
	}

	if gameInfoRaw.Data.Type != "game" {
		errWrapped := errors.Wrap(err, errWrapMessage)
		errWrapped = errors.Wrap(ErrAppSkiped, fmt.Sprintf("AppID: %s", appID))
		return nil, errWrapped
	}

	if gameInfoRaw.Data.ReleaseDate.ComingSoon {
		errWrapped := errors.Wrap(err, errWrapMessage)
		errWrapped = errors.Wrap(ErrAppSkiped, fmt.Sprintf("AppID: %s", appID))
		return nil, errWrapped
	}

	inputDateLayout := "2 Jan, 2006"
	outputDateLayout := "02.01.2006"

	releaseDateClean, err := time.Parse(inputDateLayout, gameInfoRaw.Data.ReleaseDate.Date)
	if err != nil {
		errWrapped := errors.Wrap(err, errWrapMessage)
		errWrapped = errors.Wrap(err, fmt.Sprintf("AppID: %s", appID))
		return nil, errWrapped
	}

	// TODO: handle games names
	// TODO: add game publisher if it hasn't been added before, don't forget about getting ids
	// TODO: add game tags if it hasn't been added before, don't forget about getting ids
	// TODO: add game

	publisher := model.Publisher{
		Name: gameInfoRaw.Data.Publishers[0],
	}

	game := model.Game{
		HeaderImageURL: gameInfoRaw.Data.HeaderImage,
		Name:           gameInfoRaw.Data.Name,
		Description:    gameInfoRaw.Data.Description,
		ReleaseDate:    releaseDateClean.Format(outputDateLayout),
		Publisher:      &publisher,
	}

	return &game, nil
}
