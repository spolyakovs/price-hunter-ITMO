package apiStore

import (
	"strings"

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
	// TODO: example URL in Postman, need to choose correct one (example with "The Long Dark"), check what happens when game isn't in Store
	return nil
}
