package apiserver

import (
	"net/http"
)

func (server *server) handleGames() http.HandlerFunc {
	return func(writer http.ResponseWriter, req *http.Request) {
		// methodName := "Games"
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

func (server *server) handleGamesGetByID() http.HandlerFunc {
	return func(writer http.ResponseWriter, req *http.Request) {
		// methodName := "GamesGetByID"
		// errWrapMessage := fmt.Sprintf(errHandlerMessageFormat, methodName)
		//

		server.respond(writer, req, http.StatusCreated, nil)
	}
}
