package apiserver

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
	"github.com/spolyakovs/price-hunter-ITMO/internal/app/model"
	"github.com/spolyakovs/price-hunter-ITMO/internal/app/store"
)

// TODO

func (server *server) handleFavourites() http.HandlerFunc {
	type responseItem struct {
		ID             uint64   `json:"id"`
		HeaderImageURL string   `json:"header_image"`
		Name           string   `json:"name"`
		Publisher      string   `json:"publisher"`
		ReleaseDate    string   `json:"release_date"`
		Tags           []string `json:"tags"`
	}

	return func(writer http.ResponseWriter, req *http.Request) {
		methodName := "Favourites"
		errWrapMessage := fmt.Sprintf(errHandlerMessageFormat, methodName)

		user := req.Context().Value(ctxKeyUser).(*model.User)

		games, err := server.store.Games().FindAllByUser(user)
		if err != nil {
			errWrapped := errors.Wrap(err, errWrapMessage)
			server.log(errWrapped)
			server.error(writer, req, http.StatusInternalServerError, errSomethingWentWrong)
			return
		}

		responseData := []responseItem{}

		for _, game := range games {
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

func (server *server) handleFavouritesAdd() http.HandlerFunc {
	type request struct {
		ID uint64 `json:"id"`
	}
	return func(writer http.ResponseWriter, req *http.Request) {
		methodName := "FavouritesAdd"
		errWrapMessage := fmt.Sprintf(errHandlerMessageFormat, methodName)

		requestStruct := &request{}
		if err := json.NewDecoder(req.Body).Decode(requestStruct); err != nil {
			errWrapped := errors.Wrap(err, errWrapMessage)
			server.log(errWrapped)
			server.error(writer, req, http.StatusBadRequest, errWrongRequestFormat)
			return
		}

		user := req.Context().Value(ctxKeyUser).(*model.User)

		game, err := server.store.Games().Find(requestStruct.ID)
		if err != nil {
			errWrapped := errors.Wrap(err, errWrapMessage)
			errWrapped = errors.Wrap(errWrapped, fmt.Sprintf("game.Name = %s", game.Name))
			server.log(errWrapped)

			if errors.Cause(err) == store.ErrNotFound {
				server.error(writer, req, http.StatusBadRequest, errWrongRequestFormat)
			} else {
				server.error(writer, req, http.StatusInternalServerError, errSomethingWentWrong)
			}
			return
		}

		if _, err := server.store.UserGameFavourites().FindByUserGame(user, game); err == nil {
			server.respond(writer, req, http.StatusOK, map[string]string{})
			return
		} else if errors.Cause(err) != store.ErrNotFound {
			errWrapped := errors.Wrap(err, errWrapMessage)
			errWrapped = errors.Wrap(errWrapped, fmt.Sprintf("user.Username = %s; game.Name = %s", user.Username, game.Name))
			server.log(errWrapped)

			server.error(writer, req, http.StatusInternalServerError, errSomethingWentWrong)
			return
		}

		newUserGameFavourite := &model.UserGameFavourite{
			User: user,
			Game: game,
		}

		if err := server.store.UserGameFavourites().Create(newUserGameFavourite); err != nil {
			errWrapped := errors.Wrap(err, errWrapMessage)
			errWrapped = errors.Wrap(errWrapped, fmt.Sprintf("user.Username = %s; game.Name = %s", user.Username, game.Name))
			server.log(errWrapped)

			server.error(writer, req, http.StatusInternalServerError, errSomethingWentWrong)
			return
		}

		server.respond(writer, req, http.StatusOK, map[string]string{})
	}
}

func (server *server) handleFavouritesRemove() http.HandlerFunc {
	type request struct {
		ID uint64 `json:"id"`
	}
	return func(writer http.ResponseWriter, req *http.Request) {
		methodName := "FavouritesRemove"
		errWrapMessage := fmt.Sprintf(errHandlerMessageFormat, methodName)

		requestStruct := &request{}
		if err := json.NewDecoder(req.Body).Decode(requestStruct); err != nil {
			errWrapped := errors.Wrap(err, errWrapMessage)
			server.log(errWrapped)
			server.error(writer, req, http.StatusBadRequest, errWrongRequestFormat)
			return
		}

		user := req.Context().Value(ctxKeyUser).(*model.User)

		game, err := server.store.Games().Find(requestStruct.ID)
		if err != nil {
			errWrapped := errors.Wrap(err, errWrapMessage)
			errWrapped = errors.Wrap(errWrapped, fmt.Sprintf("game.Name = %s", game.Name))
			server.log(errWrapped)

			if errors.Cause(err) == store.ErrNotFound {
				server.error(writer, req, http.StatusBadRequest, errWrongRequestFormat)
			} else {
				server.error(writer, req, http.StatusInternalServerError, errSomethingWentWrong)
			}
			return
		}

		userGameFavourite, err := server.store.UserGameFavourites().FindByUserGame(user, game)

		if err != nil {
			if errors.Cause(err) == store.ErrNotFound {
				server.respond(writer, req, http.StatusOK, map[string]string{})
				return
			}
			errWrapped := errors.Wrap(err, errWrapMessage)
			errWrapped = errors.Wrap(errWrapped, fmt.Sprintf("user.Username = %s; game.Name = %s", user.Username, game.Name))
			server.log(errWrapped)

			server.error(writer, req, http.StatusInternalServerError, errSomethingWentWrong)
			return
		}

		if err := server.store.UserGameFavourites().Delete(userGameFavourite.ID); err != nil {
			errWrapped := errors.Wrap(err, errWrapMessage)
			errWrapped = errors.Wrap(errWrapped, fmt.Sprintf("user.Username = %s; game.Name = %s", user.Username, game.Name))
			server.log(errWrapped)

			server.error(writer, req, http.StatusInternalServerError, errSomethingWentWrong)
			return
		}

		server.respond(writer, req, http.StatusOK, map[string]string{})
	}
}
