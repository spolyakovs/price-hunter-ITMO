package apiserver

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/spolyakovs/price-hunter-ITMO/internal/app/model"
)

func (server *server) handleTeamsCreate() http.HandlerFunc {
	type request struct {
		Name               string `json:"name"`
		EngineManufacturer string `json:"engine_manufacturer"`
	}

	return func(writer http.ResponseWriter, req *http.Request) {
		requestStruct := &request{}
		if err := json.NewDecoder(req.Body).Decode(requestStruct); err != nil {
			server.error(writer, req, http.StatusBadRequest, err)
			return
		}

		team := &model.Team{
			Name:               requestStruct.Name,
			EngineManufacturer: requestStruct.EngineManufacturer,
		}

		if err := server.store.Teams().Create(team); err != nil {
			server.error(writer, req, http.StatusBadRequest, err)
			return
		}

		server.respond(writer, req, http.StatusCreated, team)
	}
}

func (server *server) handleTeamsUpdate() http.HandlerFunc {
	type request struct {
		ID                 int    `json:"id"`
		Name               string `json:"name"`
		EngineManufacturer string `json:"engine_manufacturer"`
	}

	return func(writer http.ResponseWriter, req *http.Request) {
		requestStruct := &request{}
		if err := json.NewDecoder(req.Body).Decode(requestStruct); err != nil {
			server.error(writer, req, http.StatusBadRequest, err)
			return
		}

		team := &model.Team{
			ID:                 requestStruct.ID,
			Name:               requestStruct.Name,
			EngineManufacturer: requestStruct.EngineManufacturer,
		}

		if err := server.store.Teams().Update(team); err != nil {
			server.error(writer, req, http.StatusInternalServerError, err)
			return
		}
		server.respond(writer, req, http.StatusOK, team)
	}
}

func (server *server) handleTeamsDelete() http.HandlerFunc {
	type request struct {
		ID int `json:"id"`
	}

	return func(writer http.ResponseWriter, req *http.Request) {
		requestStruct := &request{}
		if err := json.NewDecoder(req.Body).Decode(requestStruct); err != nil {
			server.error(writer, req, http.StatusBadRequest, err)
			return
		}

		if err := server.store.Teams().Delete(requestStruct.ID); err != nil {
			server.error(writer, req, http.StatusBadRequest, err)
			return
		}

		server.respond(writer, req, http.StatusOK, nil)
	}
}

func (server *server) handleTeamsGet() http.HandlerFunc {
	return func(writer http.ResponseWriter, req *http.Request) {
		teams, err := server.store.Teams().GetAll()
		if err != nil {
			server.error(writer, req, http.StatusNotFound, err)
			return
		}

		server.respond(writer, req, http.StatusOK, teams)
	}
}

func (server *server) handleTeamsGetByID() http.HandlerFunc {
	return func(writer http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)

		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			fmt.Printf("DEBUG: %s\n", err.Error())
			server.error(writer, req, http.StatusBadRequest, errWrongPathValue)
			return
		}

		team, err := server.store.Teams().Find(id)
		if err != nil {
			server.error(writer, req, http.StatusNotFound, err)
			return
		}

		server.respond(writer, req, http.StatusOK, team)
	}
}

func (server *server) handleTeamsGetDriversByID() http.HandlerFunc {
	return func(writer http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)

		team_id, err := strconv.Atoi(vars["id"])
		if err != nil {
			fmt.Printf("DEBUG: %s\n", err.Error())
			server.error(writer, req, http.StatusBadRequest, errWrongPathValue)
			return
		}

		contracts, err := server.store.TeamDriverContracts().FindAllCurrentBy("team_id", team_id)
		if err != nil {
			server.error(writer, req, http.StatusNotFound, err)
			return
		}

		server.respond(writer, req, http.StatusOK, contracts)
	}
}
