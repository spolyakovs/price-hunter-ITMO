package apiserver

import (
	"net/http"
)

// TODO

func (server *server) handleFavourites() http.HandlerFunc {
	return func(writer http.ResponseWriter, req *http.Request) {
		// methodName := "Favourites"
		// errWrapMessage := fmt.Sprintf(errHandlerMessageFormat, methodName)
		//
		// tokenDetails, tokenDetailsErr := tokenUtils.CreateTokens(user.ID)
		// if tokenDetailsErr != nil {
		// 	errWrapped := errors.Wrap(tokenDetailsErr, errWrapMessage)
		// 	server.log(errWrapped)
		// 	server.error(writer, req, http.StatusInternalServerError, errSomethingWentWrong)
		// 	return
		// }
		//
		// tokens := map[string]string{
		// 	"access_token":  tokenDetails.AccessToken,
		// 	"refresh_token": tokenDetails.RefreshToken,
		// }

		server.respond(writer, req, http.StatusCreated, nil)
	}
}

func (server *server) handleFavouritesAdd() http.HandlerFunc {
	return func(writer http.ResponseWriter, req *http.Request) {
		// methodName := "FavouritesAdd"
		// errWrapMessage := fmt.Sprintf(errHandlerMessageFormat, methodName)
		//

		server.respond(writer, req, http.StatusCreated, nil)
	}
}

func (server *server) handleFavouritesRemove() http.HandlerFunc {
	return func(writer http.ResponseWriter, req *http.Request) {
		// methodName := "FavouritesRemove"
		// errWrapMessage := fmt.Sprintf(errHandlerMessageFormat, methodName)
		//

		server.respond(writer, req, http.StatusCreated, nil)
	}
}
