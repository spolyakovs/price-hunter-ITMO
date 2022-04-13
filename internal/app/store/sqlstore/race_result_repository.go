package sqlstore

import (
	"database/sql"
	"fmt"

	"github.com/spolyakovs/price-hunter-ITMO/internal/app/model"
	"github.com/spolyakovs/price-hunter-ITMO/internal/app/store"
)

type RaceResultRepository struct {
	store *Store
}

func (raceResultRepository *RaceResultRepository) Create(raceResult *model.RaceResult) error {
	createQuery := "INSERT INTO race_results (place, points, race_id, driver_id) VALUES ($1, $2, $3, $4) RETURNING id;"

	return raceResultRepository.store.db.Get(
		&raceResult.ID,
		createQuery,
		raceResult.Place,
		raceResult.Points,
		raceResult.Race.ID,
		raceResult.Driver.ID,
	)
}

func (raceResultRepository *RaceResultRepository) Find(id int) (*model.RaceResult, error) {
	return raceResultRepository.FindBy("id", id)
}

func (raceResultRepository *RaceResultRepository) FindBy(columnName string, value interface{}) (*model.RaceResult, error) {
	raceResult := &model.RaceResult{}

	findQuery := fmt.Sprintf("SELECT "+
		"race_results.id AS id, "+
		"race_results.place AS place, "+
		"race_results.points AS points, "+
		"races.id AS \"race.id\", "+
		"races.name AS \"race.name\", "+
		"races.location AS \"race.location\", "+
		"races.date AS \"race.date\", "+
		"drivers.id AS \"driver.id\", "+
		"drivers.first_name AS \"driver.first_name\", "+
		"drivers.last_name AS \"driver.last_name\", "+
		"drivers.birth_date AS \"driver.birth_date\" "+
		"FROM race_results "+
		"LEFT JOIN races "+
		"ON (race_results.race_id = races.id) "+
		"LEFT JOIN drivers "+
		"ON (race_results.driver_id = drivers.id) "+
		"WHERE race_results.%s = $1 LIMIT 1;", columnName)

	if err := raceResultRepository.store.db.Get(
		raceResult,
		findQuery,
		value,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrNotFound
		}

		return nil, err
	}
	return raceResult, nil
}

func (raceResultRepository *RaceResultRepository) FindAllBy(columnName string, value interface{}) ([]*model.RaceResult, error) {
	raceResults := []*model.RaceResult{}

	findQuery := fmt.Sprintf("SELECT "+
		"race_results.id AS id, "+
		"race_results.place AS place, "+
		"race_results.points AS points, "+
		"races.id AS \"race.id\", "+
		"races.name AS \"race.name\", "+
		"races.location AS \"race.location\", "+
		"races.date AS \"race.date\", "+
		"drivers.id AS \"driver.id\", "+
		"drivers.first_name AS \"driver.first_name\", "+
		"drivers.last_name AS \"driver.last_name\", "+
		"drivers.birth_date AS \"driver.birth_date\" "+
		"FROM race_results "+
		"LEFT JOIN races "+
		"ON (race_results.race_id = races.id) "+
		"LEFT JOIN drivers "+
		"ON (race_results.driver_id = drivers.id) "+
		"WHERE race_results.%s = $1;", columnName)

	if err := raceResultRepository.store.db.Select(
		&raceResults,
		findQuery,
		value,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrNotFound
		}

		return nil, err
	}
	return raceResults, nil
}

func (raceResultRepository *RaceResultRepository) Update(raceResult *model.RaceResult) error {
	updateQuery := "UPDATE race_results " +
		"SET place = :place, " +
		"points = :points, " +
		"race_id = :race.id " +
		"driver_id = :driver.id, " +
		"WHERE id = :id;"

	countResult, countResultErr := raceResultRepository.store.db.NamedExec(
		updateQuery,
		raceResult,
	)

	if countResultErr != nil {
		return countResultErr
	}

	count, countErr := countResult.RowsAffected()

	if countErr != nil {
		return countErr
	}

	if count == 0 {
		return store.ErrNotFound
	}

	return nil
}

func (raceResultRepository *RaceResultRepository) Delete(id int) error {
	deleteQuery := "DELETE FROM race_results WHERE id = $1;"

	countResult, countResultErr := raceResultRepository.store.db.Exec(
		deleteQuery,
		id,
	)

	if countResultErr != nil {
		return countResultErr
	}

	count, countErr := countResult.RowsAffected()

	if countErr != nil {
		return countErr
	}

	if count == 0 {
		return store.ErrNotFound
	}

	return nil
}
