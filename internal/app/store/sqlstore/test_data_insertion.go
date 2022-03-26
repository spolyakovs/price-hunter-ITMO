package sqlstore

import (
	"fmt"

	"github.com/spolyakovs/price-hunter-ITMO/internal/app/model"
)

var (
	users       []*model.User
	teams       []*model.Team
	drivers     []*model.Driver
	races       []*model.Race
	contracts   []*model.TeamDriverContract
	raceResults []*model.RaceResult
)

func (store *Store) fillTables() error {
	if err := store.fillTableUsers(); err != nil {
		return err
	}

	if err := store.fillTableTeams(); err != nil {
		return err
	}

	if err := store.fillTableDrivers(); err != nil {
		return err
	}

	if err := store.fillTableRaces(); err != nil {
		return err
	}

	if err := store.fillTableTeamDriverContracts(); err != nil {
		return err
	}

	if err := store.fillTableRaceResults(); err != nil {
		return err
	}

	return nil
}

func (store *Store) fillTableUsers() error {
	for i := 1; i <= 5; i++ {
		users = append(users, &model.User{
			Username: fmt.Sprintf("test_username_%d", i),
			Email:    fmt.Sprintf("test_email_%d@example.org", i),
			Password: fmt.Sprintf("test_password_%d", i),
		})
	}

	for _, user := range users {
		if err := store.Users().Create(user); err != nil {
			return err
		}
	}

	return nil
}

func (store *Store) fillTableTeams() error {
	teams = append(teams, &model.Team{
		Name:               "Mercedes",
		EngineManufacturer: "Mercedes",
	})
	teams = append(teams, &model.Team{
		Name:               "Red Bull Racing",
		EngineManufacturer: "Honda",
	})
	teams = append(teams, &model.Team{
		Name:               "McLaren",
		EngineManufacturer: "Mercedes",
	})

	for _, team := range teams {
		if err := store.Teams().Create(team); err != nil {
			return err
		}
	}

	return nil
}

func (store *Store) fillTableDrivers() error {
	drivers = append(drivers, &model.Driver{
		FirstName: "Lewis",
		LastName:  "Hammilton",
		BirthDate: "1985-01-07",
	})
	drivers = append(drivers, &model.Driver{
		FirstName: "Valtteri",
		LastName:  "Bottas",
		BirthDate: "1985-08-28",
	})
	drivers = append(drivers, &model.Driver{
		FirstName: "Max",
		LastName:  "Verstappen",
		BirthDate: "1997-08-30",
	})
	drivers = append(drivers, &model.Driver{
		FirstName: "Sergio",
		LastName:  "Perez",
		BirthDate: "1990-01-26",
	})
	drivers = append(drivers, &model.Driver{
		FirstName: "Daniel",
		LastName:  "Ricciardo",
		BirthDate: "1989-07-01",
	})
	drivers = append(drivers, &model.Driver{
		FirstName: "Lando",
		LastName:  "Norris",
		BirthDate: "1999-11-13",
	})

	for _, driver := range drivers {
		if err := store.Drivers().Create(driver); err != nil {
			return err
		}
	}

	return nil
}

func (store *Store) fillTableRaces() error {
	races = append(races, &model.Race{
		Name:     "Italian Grand Prix 2021",
		Location: "Autodromo Enzo e Dino Ferrari",
		Date:     "2021-04-18",
	})
	races = append(races, &model.Race{
		Name:     "Potugal Grand Prix 2021",
		Location: "Autodromo Internacional do Algarve",
		Date:     "2021-05-02",
	})
	races = append(races, &model.Race{
		Name:     "Spanish Grand Prix 2021",
		Location: "Circuit de Barcelona-Catalunya",
		Date:     "2021-05-09",
	})

	for _, race := range races {
		if err := store.Races().Create(race); err != nil {
			return err
		}
	}

	return nil
}

func (store *Store) fillTableTeamDriverContracts() error {
	contracts = append(contracts, &model.TeamDriverContract{
		FromDate: "2013-01-01",
		ToDate:   "2023-12-31",
		Team:     teams[0],
		Driver:   drivers[0],
	})
	contracts = append(contracts, &model.TeamDriverContract{
		FromDate: "2017-01-01",
		ToDate:   "2021-12-31",
		Team:     teams[0],
		Driver:   drivers[1],
	})
	contracts = append(contracts, &model.TeamDriverContract{
		FromDate: "2016-01-01",
		ToDate:   "2023-12-31",
		Team:     teams[1],
		Driver:   drivers[2],
	})
	contracts = append(contracts, &model.TeamDriverContract{
		FromDate: "2013-01-01",
		ToDate:   "2013-12-31",
		Team:     teams[2],
		Driver:   drivers[3],
	})
	contracts = append(contracts, &model.TeamDriverContract{
		FromDate: "2021-01-01",
		ToDate:   "2022-12-31",
		Team:     teams[1],
		Driver:   drivers[3],
	})
	contracts = append(contracts, &model.TeamDriverContract{
		FromDate: "2014-01-01",
		ToDate:   "2018-12-31",
		Team:     teams[1],
		Driver:   drivers[4],
	})
	contracts = append(contracts, &model.TeamDriverContract{
		FromDate: "2021-01-01",
		ToDate:   "2023-12-31",
		Team:     teams[2],
		Driver:   drivers[4],
	})
	contracts = append(contracts, &model.TeamDriverContract{
		FromDate: "2019-01-01",
		Team:     teams[2],
		Driver:   drivers[5],
	})

	for _, contract := range contracts {
		if err := store.TeamDriverContracts().Create(contract); err != nil {
			return err
		}
	}

	return nil
}

func (store *Store) fillTableRaceResults() error {
	raceResults = append(raceResults, &model.RaceResult{
		Place:  2,
		Race:   races[0],
		Driver: drivers[0],
	})
	raceResults = append(raceResults, &model.RaceResult{
		Place:  18,
		Race:   races[0],
		Driver: drivers[1],
	})
	raceResults = append(raceResults, &model.RaceResult{
		Place:  1,
		Race:   races[0],
		Driver: drivers[2],
	})
	raceResults = append(raceResults, &model.RaceResult{
		Place:  11,
		Race:   races[0],
		Driver: drivers[3],
	})
	raceResults = append(raceResults, &model.RaceResult{
		Place:  6,
		Race:   races[0],
		Driver: drivers[4],
	})
	raceResults = append(raceResults, &model.RaceResult{
		Place:  3,
		Race:   races[0],
		Driver: drivers[5],
	})

	raceResults = append(raceResults, &model.RaceResult{
		Place:  1,
		Race:   races[1],
		Driver: drivers[0],
	})
	raceResults = append(raceResults, &model.RaceResult{
		Place:  3,
		Race:   races[1],
		Driver: drivers[1],
	})
	raceResults = append(raceResults, &model.RaceResult{
		Place:  2,
		Race:   races[1],
		Driver: drivers[2],
	})
	raceResults = append(raceResults, &model.RaceResult{
		Place:  4,
		Race:   races[1],
		Driver: drivers[3],
	})
	raceResults = append(raceResults, &model.RaceResult{
		Place:  9,
		Race:   races[1],
		Driver: drivers[4],
	})
	raceResults = append(raceResults, &model.RaceResult{
		Place:  5,
		Race:   races[1],
		Driver: drivers[5],
	})

	raceResults = append(raceResults, &model.RaceResult{
		Place:  1,
		Race:   races[2],
		Driver: drivers[0],
	})
	raceResults = append(raceResults, &model.RaceResult{
		Place:  3,
		Race:   races[2],
		Driver: drivers[1],
	})
	raceResults = append(raceResults, &model.RaceResult{
		Place:  2,
		Race:   races[2],
		Driver: drivers[2],
	})
	raceResults = append(raceResults, &model.RaceResult{
		Place:  5,
		Race:   races[2],
		Driver: drivers[3],
	})
	raceResults = append(raceResults, &model.RaceResult{
		Place:  6,
		Race:   races[2],
		Driver: drivers[4],
	})
	raceResults = append(raceResults, &model.RaceResult{
		Place:  8,
		Race:   races[2],
		Driver: drivers[5],
	})

	for _, raceResult := range raceResults {
		raceResult.Points = pointsByPlace[raceResult.Place]
		if err := store.RaceResults().Create(raceResult); err != nil {
			return err
		}
	}

	return nil
}
