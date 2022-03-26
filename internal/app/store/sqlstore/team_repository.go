package sqlstore

import (
	"database/sql"
	"fmt"

	"github.com/spolyakovs/price-hunter-ITMO/internal/app/model"
	"github.com/spolyakovs/price-hunter-ITMO/internal/app/store"
)

type TeamRepository struct {
	store *Store
}

func (teamRepository *TeamRepository) Create(team *model.Team) error {
	createQuery := "INSERT INTO teams (name, engine_manufacturer) VALUES ($1, $2) RETURNING id"
	return teamRepository.store.db.Get(
		&team.ID,
		createQuery,
		team.Name, team.EngineManufacturer,
	)
}

func (teamRepository *TeamRepository) Find(id int) (*model.Team, error) {
	return teamRepository.FindBy("id", id)
}

func (teamRepository *TeamRepository) FindBy(columnName string, value interface{}) (*model.Team, error) {
	team := &model.Team{}

	findQuery := fmt.Sprintf("SELECT * FROM teams WHERE %s = $1 LIMIT 1;", columnName)
	if err := teamRepository.store.db.Get(
		team,
		findQuery,
		value,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}

		return nil, err
	}

	return team, nil
}

func (teamRepository *TeamRepository) FindAllBy(columnName string, value interface{}) ([]*model.Team, error) {
	teams := []*model.Team{}

	findQuery := fmt.Sprintf("SELECT * FROM teams WHERE %s = $1;", columnName)
	if err := teamRepository.store.db.Select(
		&teams,
		findQuery,
		value,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}

		return nil, err
	}

	return teams, nil
}

func (teamRepository *TeamRepository) GetAll() ([]*model.Team, error) {
	teams := []*model.Team{}

	findQuery := "SELECT * FROM teams;"
	if err := teamRepository.store.db.Select(
		&teams,
		findQuery,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}

		return nil, err
	}

	return teams, nil
}

func (teamRepository *TeamRepository) Update(team *model.Team) error {
	updateQuery := "UPDATE teams " +
		`SET name = :name, ` +
		`engine_manufacturer = :engine_manufacturer ` +
		`WHERE id = :id;`

	countResult, countResultErr := teamRepository.store.db.NamedExec(
		updateQuery,
		team,
	)

	if countResultErr != nil {
		return countResultErr
	}

	count, countErr := countResult.RowsAffected()

	if countErr != nil {
		return countErr
	}

	if count == 0 {
		return store.ErrRecordNotFound
	}

	return nil
}

func (teamRepository *TeamRepository) Delete(id int) error {
	deleteQuery := "DELETE FROM teams WHERE id = $1"

	countResult, countResultErr := teamRepository.store.db.Exec(
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
		return store.ErrRecordNotFound
	}

	return nil
}
