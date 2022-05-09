package apiserver

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"

	"github.com/spolyakovs/price-hunter-ITMO/internal/app/model"
	"github.com/spolyakovs/price-hunter-ITMO/internal/app/store"
)

// TODO: handleTags (list of all tag names)

// TODO: add offset param
func (server *server) handleGames() http.HandlerFunc {
	type request struct {
		Query string   `json:"query,omitempty"`
		Tags  []string `json:"tags,omitempty"`
		Limit int      `json:"limit,omitempty"`
	}
	type responseItem struct {
		ID             uint64   `json:"id"`
		HeaderImageURL string   `json:"header_image"`
		Name           string   `json:"name"`
		Publisher      string   `json:"publisher"`
		ReleaseDate    string   `json:"release_date"`
		Tags           []string `json:"tags"`
	}

	return func(writer http.ResponseWriter, req *http.Request) {
		methodName := "Games"
		errWrapMessage := fmt.Sprintf(errHandlerMessageFormat, methodName)

		requestStruct := &request{
			Limit: 500,
		}
		if err := json.NewDecoder(req.Body).Decode(requestStruct); err != nil {
			errWrapped := errors.Wrap(err, errWrapMessage)
			server.log(errWrapped)
			server.error(writer, req, http.StatusBadRequest, errWrongRequestFormat)
			return
		}

		queryTags := []*model.Tag{}

		for _, tagName := range requestStruct.Tags {
			if !regexp.MustCompile("[a-z_]+").MatchString(tagName) {
				errWrapped := errors.Wrap(errWrongRequestFormat, errWrapMessage)
				errWrapped = errors.Wrap(errWrapped, fmt.Sprintf("TagName = %s", tagName))
				server.error(writer, req, http.StatusBadRequest, errWrongRequestFormat)
				return
			}

			tag, err := server.store.Tags().FindBy("name", tagName)
			if err != nil {
				errWrapped := errors.Wrap(err, errWrapMessage)
				errWrapped = errors.Wrap(errWrapped, fmt.Sprintf("TagName = %s", tagName))
				server.log(errWrapped)

				if errors.Cause(err) == store.ErrNotFound {
					server.error(writer, req, http.StatusBadRequest, errWrongRequestFormat)
				} else {
					server.error(writer, req, http.StatusInternalServerError, errSomethingWentWrong)
				}
				return
			}
			queryTags = append(queryTags, tag)
		}

		games, err := server.store.Games().FindAllByQueryTags(requestStruct.Query, queryTags)
		if err != nil {
			errWrapped := errors.Wrap(err, errWrapMessage)
			server.log(errWrapped)
			server.error(writer, req, http.StatusInternalServerError, errSomethingWentWrong)
			return
		}

		responseData := []responseItem{}

		for _, game := range games {
			if len(responseData) >= requestStruct.Limit {
				break
			}
			tags, err := server.store.Tags().FindAllByGame(game)
			if err != nil {
				errWrapped := errors.Wrap(err, errWrapMessage)
				server.log(errWrapped)
				server.error(writer, req, http.StatusInternalServerError, errSomethingWentWrong)
				return
			}

			tagNames := []string{}
			for _, tag := range tags {
				tagNames = append(tagNames, tag.Name)
			}

			responseItemStruct := responseItem{
				ID:             game.ID,
				HeaderImageURL: game.HeaderImageURL,
				Name:           game.Name,
				Publisher:      game.Publisher.Name,
				ReleaseDate:    game.ReleaseDate,
				Tags:           tagNames,
			}

			responseData = append(responseData, responseItemStruct)
		}

		server.respond(writer, req, http.StatusOK, responseData)
	}
}

func (server *server) handleGamesGetByID() http.HandlerFunc {
	type responsePricesItem struct {
		InitialFormatted string `json:"initial_formatted"`
		FinalFormatted   string `json:"final_formatted"`
		DiscountPercent  int    `json:"discount_percent"`
		MarketGameURL    string `json:"uri_string"`
	}
	type response struct {
		ID             uint64                        `json:"id"`
		HeaderImageURL string                        `json:"header_image"`
		Name           string                        `json:"name"`
		Publisher      string                        `json:"publisher"`
		Description    string                        `json:"description"`
		ReleaseDate    string                        `json:"release_date"`
		IsFavourite    bool                          `json:"is_favourite"`
		Tags           []string                      `json:"tags"`
		Prices         map[string]responsePricesItem `json:"prices"`
	}

	return func(writer http.ResponseWriter, req *http.Request) {
		methodName := "GamesGetByID"
		errWrapMessage := fmt.Sprintf(errHandlerMessageFormat, methodName)

		vars := mux.Vars(req)

		id, err := strconv.ParseUint(vars["id"], 10, 64)
		if err != nil {
			errWrapped := errors.Wrap(err, errWrapMessage)
			server.log(errWrapped)
			server.error(writer, req, http.StatusBadRequest, errWrongRequestFormat)
			return
		}

		game, err := server.store.Games().Find(id)
		if err != nil {
			errWrapped := errors.Wrap(err, errWrapMessage)
			errWrapped = errors.Wrap(errWrapped, fmt.Sprintf("ID = %d", id))
			server.log(errWrapped)

			if errors.Cause(err) == store.ErrNotFound {
				server.error(writer, req, http.StatusBadRequest, errWrongRequestFormat)
			} else {
				server.error(writer, req, http.StatusInternalServerError, errSomethingWentWrong)
			}
			return
		}

		tags, err := server.store.Tags().FindAllByGame(game)
		if err != nil {
			errWrapped := errors.Wrap(err, errWrapMessage)
			server.log(errWrapped)
			server.error(writer, req, http.StatusInternalServerError, errSomethingWentWrong)
			return
		}

		tagNames := []string{}
		for _, tag := range tags {
			tagNames = append(tagNames, tag.Name)
		}

		user := req.Context().Value(ctxKeyUser).(*model.User)
		isFavourite := false

		if _, err := server.store.UserGameFavourites().FindByUserGame(user, game); err == nil {
			isFavourite = true
		} else if errors.Cause(err) != store.ErrNotFound {
			errWrapped := errors.Wrap(err, errWrapMessage)
			server.log(errWrapped)
			server.error(writer, req, http.StatusInternalServerError, errSomethingWentWrong)
			return
		}

		responseStruct := response{
			ID:             game.ID,
			HeaderImageURL: game.HeaderImageURL,
			Name:           game.Name,
			Publisher:      game.Publisher.Name,
			ReleaseDate:    game.ReleaseDate,
			Description:    game.Description,
			IsFavourite:    isFavourite,
			Tags:           tagNames,
			Prices:         make(map[string]responsePricesItem),
		}

		gameMarketPrices, err := server.store.GameMarketPrices().FindAllByGame(game)
		if err != nil {
			errWrapped := errors.Wrap(err, errWrapMessage)
			server.log(errWrapped)
			server.error(writer, req, http.StatusInternalServerError, errSomethingWentWrong)
			return
		}

		for _, gameMarketPrice := range gameMarketPrices {
			responsePricesItemStruct := responsePricesItem{
				InitialFormatted: gameMarketPrice.InitialValueFormatted,
				FinalFormatted:   gameMarketPrice.FinalValueFormatted,
				DiscountPercent:  gameMarketPrice.DiscountPercent,
				MarketGameURL:    gameMarketPrice.MarketGameURL,
			}

			// TODO: just use market name lowercased
			switch gameMarketPrice.Market.Name {
			case "Steam":
				responseStruct.Prices["steam"] = responsePricesItemStruct
			case "EpicGamesStore":
				responseStruct.Prices["egs"] = responsePricesItemStruct
			case "GOG.com":
				responseStruct.Prices["gog"] = responsePricesItemStruct
			default:
				server.error(writer, req, http.StatusInternalServerError, errSomethingWentWrong)
			}
		}

		server.respond(writer, req, http.StatusOK, responseStruct)
	}
}
