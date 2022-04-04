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
	private.HandleFunc("/users/me", server.handleUsersMe())
	private.HandleFunc("/users/update", server.handleUsersUpdate())
	private.HandleFunc("/users/delete", server.handleUsersDelete())
	private.HandleFunc("/users/{id:[0-9]+}", server.handleUsersGetByID())

	private.HandleFunc("/teams", server.handleTeamsGet())
	private.HandleFunc("/teams/create", server.handleTeamsCreate())
	private.HandleFunc("/teams/update", server.handleTeamsUpdate())
	private.HandleFunc("/teams/delete", server.handleTeamsDelete())
	private.HandleFunc("/teams/{id:[0-9]+}", server.handleTeamsGetByID())
	private.HandleFunc("/teams/{id:[0-9]+}/drivers", server.handleTeamsGetDriversByID())

	private.HandleFunc("/drivers", server.handleDriversGet())
	private.HandleFunc("/drivers/create", server.handleDriversCreate())
	private.HandleFunc("/drivers/update", server.handleDriversUpdate())
	private.HandleFunc("/drivers/delete", server.handleDriversDelete())
	private.HandleFunc("/drivers/{id:[0-9]+}", server.handleDriversGetByID())
	private.HandleFunc("/drivers/{id:[0-9]+}/team", server.handleDriversGetTeamByID())
	private.HandleFunc("/drivers/{id:[0-9]+}/career", server.handleDriversGetCareerByID())

	private.HandleFunc("/races", server.handleRacesGet())
	private.HandleFunc("/races/create", server.handleRacesCreate())
	private.HandleFunc("/races/update", server.handleRacesUpdate())
	private.HandleFunc("/races/delete", server.handleRacesDelete())
	private.HandleFunc("/races/{id:[0-9]+}", server.handleRacesGetByID())
}
