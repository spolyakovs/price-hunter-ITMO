package apiserver

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/spolyakovs/price-hunter-ITMO/internal/app/model"
)

func (server *server) handleRacesCreate() http.HandlerFunc {
	type request struct {
		Name     string `json:"name"`
		Location string `json:"location"`
		Date     string `json:"date"`
	}

	return func(writer http.ResponseWriter, req *http.Request) {
		requestStruct := &request{}
		if err := json.NewDecoder(req.Body).Decode(requestStruct); err != nil {
			server.error(writer, req, http.StatusBadRequest, err)
			return
		}

		race := &model.Race{
			Name:     requestStruct.Name,
			Location: requestStruct.Location,
			Date:     requestStruct.Date,
		}

		if err := server.store.Races().Create(race); err != nil {
			server.error(writer, req, http.StatusBadRequest, err)
			return
		}

		server.respond(writer, req, http.StatusCreated, race)
	}
}

func (server *server) handleRacesUpdate() http.HandlerFunc {
	type request struct {
		ID       int    `json:"id"`
		Name     string `json:"name"`
		Location string `json:"location"`
		Date     string `json:"date"`
	}

	return func(writer http.ResponseWriter, req *http.Request) {
		requestStruct := &request{}
		if err := json.NewDecoder(req.Body).Decode(requestStruct); err != nil {
			server.error(writer, req, http.StatusBadRequest, err)
			return
		}

		race := &model.Race{
			ID:       requestStruct.ID,
			Name:     requestStruct.Name,
			Location: requestStruct.Location,
			Date:     requestStruct.Date,
		}

		if err := server.store.Races().Update(race); err != nil {
			server.error(writer, req, http.StatusInternalServerError, err)
			return
		}
		server.respond(writer, req, http.StatusOK, race)
	}
}

func (server *server) handleRacesDelete() http.HandlerFunc {
	type request struct {
		ID int `json:"id"`
	}

	return func(writer http.ResponseWriter, req *http.Request) {
		requestStruct := &request{}
		if err := json.NewDecoder(req.Body).Decode(requestStruct); err != nil {
			server.error(writer, req, http.StatusBadRequest, err)
			return
		}

		if err := server.store.Races().Delete(requestStruct.ID); err != nil {
			server.error(writer, req, http.StatusBadRequest, err)
			return
		}

		server.respond(writer, req, http.StatusOK, nil)
	}
}

func (server *server) handleRacesGet() http.HandlerFunc {
	return func(writer http.ResponseWriter, req *http.Request) {
		races, err := server.store.Races().GetAll()
		if err != nil {
			server.error(writer, req, http.StatusNotFound, err)
			return
		}

		server.respond(writer, req, http.StatusOK, races)
	}
}

func (server *server) handleRacesGetByID() http.HandlerFunc {
	type response struct {
		Race    *model.Race
		Results []*model.RaceResult
	}

	return func(writer http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)

		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			fmt.Printf("DEBUG: %s\n", err.Error())
			server.error(writer, req, http.StatusBadRequest, errWrongPathValue)
			return
		}

		race, err := server.store.Races().Find(id)
		if err != nil {
			server.error(writer, req, http.StatusNotFound, err)
			return
		}

		raceResults, err := server.store.RaceResults().FindAllBy("race_id", id)
		if err != nil {
			server.error(writer, req, http.StatusNotFound, err)
			return
		}

		resp := response{
			Race:    race,
			Results: raceResults,
		}

		server.respond(writer, req, http.StatusOK, resp)
	}
}
