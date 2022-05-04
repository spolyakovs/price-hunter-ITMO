package apiStore

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
	"github.com/spolyakovs/price-hunter-ITMO/internal/app/model"
)

// TODO: write APIEpicGames
// TODO: write APIGog

type APISteam struct {
	apiKey string
}

func NewAPISteam(apiKey string) *APISteam {
	return &APISteam{
		apiKey: apiKey,
	}
}

func (api *APISteam) GetGamesFull() ([]*model.Game, error) {
	type gamesListResponseListItem struct {
		AppIDs int `json:"appid"`
	}

	type gamesListResponseList struct {
		Apps []gamesListResponseListItem `json:"apps"`
	}

	type gamesListResponse struct {
		AppList gamesListResponseList `json:"applist"`
	}

	apiName := "Steam"
	methodName := "GetGamesFull"
	errWrapMessage := fmt.Sprintf(errAPIStoreMessageFormat, apiName, methodName)

	idsListURL := fmt.Sprintf("http://api.steampowered.com/ISteamApps/GetAppList/v0002/?key=%s&format=json", api.apiKey)

	resp, err := http.Get(idsListURL)
	if err != nil {
		errWrapped := errors.Wrap(err, errWrapMessage)
		return nil, errWrapped
	}

	gamesListResponseStruct := &gamesListResponse{}

	if err := json.NewDecoder(resp.Body).Decode(gamesListResponseStruct); err != nil {
		errWrapped := errors.Wrap(err, errWrapMessage)
		return nil, errWrapped
	}

	// TODO: get all apps info from gamesListResponseStruct
	// TODO: handle games names
	// TODO: add game tags if it hasn't been added before

	fmt.Printf("%+v\n", gamesListResponseStruct)

	return nil, nil
}
