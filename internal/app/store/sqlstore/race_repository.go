package sqlstore

import (
	"database/sql"
	"fmt"

	"github.com/spolyakovs/price-hunter-ITMO/internal/app/model"
	"github.com/spolyakovs/price-hunter-ITMO/internal/app/store"
)

type RaceRepository struct {
	store *Store
}

func (raceRepository *RaceRepository) Create(race *model.Race) error {
	createQuery := "INSERT INTO races (name, location, date) VALUES ($1, $2, $3) RETURNING id;"

	return raceRepository.store.db.Get(
		&race.ID,
		createQuery,
		race.Name, race.Location, race.Date,
	)
}

func (raceRepository *RaceRepository) Find(id int) (*model.Race, error) {
	return raceRepository.FindBy("id", id, "")
}

func (raceRepository *RaceRepository) FindBy(columnName string, value interface{}, condition string) (*model.Race, error) {
	race := &model.Race{}

	findQuery := fmt.Sprintf("SELECT * FROM races WHERE %s = $1%s LIMIT 1;", columnName, condition)
	if err := raceRepository.store.db.Get(
		race,
		findQuery,
		value,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrNotFound
		}

		return nil, err
	}

	return race, nil
}

func (raceRepository *RaceRepository) FindAllBy(columnName string, value interface{}, condition string) ([]*model.Race, error) {
	races := []*model.Race{}

	findQuery := fmt.Sprintf("SELECT * FROM races WHERE %s = $1%s;", columnName, condition)
	if err := raceRepository.store.db.Select(
		&races,
		findQuery,
		value,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrNotFound
		}

		return nil, err
	}

	return races, nil
}

func (raceRepository *RaceRepository) GetAll() ([]*model.Race, error) {
	races := []*model.Race{}

	findQuery := "SELECT * FROM races;"
	if err := raceRepository.store.db.Select(
		&races,
		findQuery,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrNotFound
		}

		return nil, err
	}

	return races, nil
}

func (raceRepository *RaceRepository) Update(race *model.Race) error {
	updateQuery := "UPDATE races " +
		"SET name = :name, " +
		"location = :location, " +
		"date = :date " +
		"WHERE id = :id;"

	countResult, countResultErr := raceRepository.store.db.NamedExec(
		updateQuery,
		race,
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

func (raceRepository *RaceRepository) Delete(id int) error {
	deleteQuery := "DELETE FROM races WHERE id = $1;"

	countResult, countResultErr := raceRepository.store.db.Exec(
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
