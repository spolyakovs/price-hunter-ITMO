package apiserver

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/spolyakovs/price-hunter-ITMO/internal/app/model"
)

func (server *server) handleDriversCreate() http.HandlerFunc {
	type request struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		BirthDate string `json:"birth_date"`
	}

	return func(writer http.ResponseWriter, req *http.Request) {
		requestStruct := &request{}
		if err := json.NewDecoder(req.Body).Decode(requestStruct); err != nil {
			server.error(writer, req, http.StatusBadRequest, err)
			return
		}

		driver := &model.Driver{
			FirstName: requestStruct.FirstName,
			LastName:  requestStruct.LastName,
			BirthDate: requestStruct.BirthDate,
		}

		if err := server.store.Drivers().Create(driver); err != nil {
			server.error(writer, req, http.StatusBadRequest, err)
			return
		}

		server.respond(writer, req, http.StatusCreated, driver)
	}
}

func (server *server) handleDriversUpdate() http.HandlerFunc {
	type request struct {
		ID        int    `json:"id"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		BirthDate string `json:"birth_date"`
	}

	return func(writer http.ResponseWriter, req *http.Request) {
		requestStruct := &request{}
		if err := json.NewDecoder(req.Body).Decode(requestStruct); err != nil {
			server.error(writer, req, http.StatusBadRequest, err)
			return
		}

		driver := &model.Driver{
			ID:        requestStruct.ID,
			FirstName: requestStruct.FirstName,
			LastName:  requestStruct.LastName,
			BirthDate: requestStruct.BirthDate,
		}

		if err := server.store.Drivers().Update(driver); err != nil {
			server.error(writer, req, http.StatusInternalServerError, err)
			return
		}
		server.respond(writer, req, http.StatusOK, driver)
	}
}

func (server *server) handleDriversDelete() http.HandlerFunc {
	type request struct {
		ID int `json:"id"`
	}

	return func(writer http.ResponseWriter, req *http.Request) {
		requestStruct := &request{}
		if err := json.NewDecoder(req.Body).Decode(requestStruct); err != nil {
			server.error(writer, req, http.StatusBadRequest, err)
			return
		}

		if err := server.store.Drivers().Delete(requestStruct.ID); err != nil {
			server.error(writer, req, http.StatusBadRequest, err)
			return
		}

		server.respond(writer, req, http.StatusOK, nil)
	}
}

func (server *server) handleDriversGet() http.HandlerFunc {
	return func(writer http.ResponseWriter, req *http.Request) {
		drivers, err := server.store.Drivers().GetAll()
		if err != nil {
			server.error(writer, req, http.StatusNotFound, err)
			return
		}

		server.respond(writer, req, http.StatusOK, drivers)
	}
}

func (server *server) handleDriversGetByID() http.HandlerFunc {
	return func(writer http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)

		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			fmt.Printf("DEBUG: %s\n", err.Error())
			server.error(writer, req, http.StatusBadRequest, errWrongPathValue)
			return
		}

		driver, err := server.store.Drivers().Find(id)
		if err != nil {
			server.error(writer, req, http.StatusNotFound, err)
			return
		}

		server.respond(writer, req, http.StatusOK, driver)
	}
}

func (server *server) handleDriversGetTeamByID() http.HandlerFunc {
	return func(writer http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)

		driver_id, err := strconv.Atoi(vars["id"])
		if err != nil {
			fmt.Printf("DEBUG: %s\n", err.Error())
			server.error(writer, req, http.StatusBadRequest, errWrongPathValue)
			return
		}

		contract, err := server.store.TeamDriverContracts().FindCurrentBy("driver_id", driver_id)
		if err != nil {
			server.error(writer, req, http.StatusNotFound, err)
			return
		}

		server.respond(writer, req, http.StatusOK, contract)
	}
}

func (server *server) handleDriversGetCareerByID() http.HandlerFunc {
	return func(writer http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)

		driver_id, err := strconv.Atoi(vars["id"])
		if err != nil {
			fmt.Printf("DEBUG: %s\n", err.Error())
			server.error(writer, req, http.StatusBadRequest, errWrongPathValue)
			return
		}

		raceResults, err := server.store.RaceResults().FindAllBy("driver_id", driver_id)
		if err != nil {
			server.error(writer, req, http.StatusNotFound, err)
			return
		}

		server.respond(writer, req, http.StatusOK, raceResults)
	}
}
