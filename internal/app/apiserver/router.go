package apiserver

import "github.com/gorilla/handlers"

func (server *server) configureRouter() {
	server.router.Use(server.setRequestID)
	server.router.Use(server.logRequest)
	server.router.Use(handlers.CORS(handlers.AllowedOrigins([]string{"*"})))
	server.router.HandleFunc("/registration", server.handleRegistration()).Methods("POST")
	server.router.HandleFunc("/login", server.handleLogin()).Methods("POST")
	server.router.HandleFunc("/logout", server.handleLogout()).Methods("POST")
	server.router.HandleFunc("/token/refresh", server.handleRefreshToken()).Methods("POST")

	private := server.router.PathPrefix("/private").Subrouter()
	private.Use(server.authenticateUser)
	private.HandleFunc("/me", server.handleUsersMe()).Methods("GET")
	private.HandleFunc("/change/email", server.handleUsersChangeEmail()).Methods("POST")
	private.HandleFunc("/change/password", server.handleUsersChangePassword()).Methods("POST")

	private.HandleFunc("/games", server.handleGames()).Methods("POST")
	private.HandleFunc("/games/{id:[0-9]+}", server.handleGamesGetByID()).Methods("GET")

	private.HandleFunc("/favourites", server.handleFavourites()).Methods("GET")
	private.HandleFunc("/favourites/add", server.handleFavouritesAdd()).Methods("POST")
	private.HandleFunc("/favourites/remove", server.handleFavouritesRemove()).Methods("POST")
}
