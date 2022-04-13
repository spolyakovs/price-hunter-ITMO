package sqlstore

import (
	"database/sql"
	"fmt"

	"github.com/spolyakovs/price-hunter-ITMO/internal/app/model"
	"github.com/spolyakovs/price-hunter-ITMO/internal/app/store"
)

type TeamDriverContractRepository struct {
	store *Store
}

func (teamDriverContractRepository *TeamDriverContractRepository) Create(contract *model.TeamDriverContract) error {
	if contract.ToDate != "" {
		createQuery := "INSERT INTO team_driver_contracts (from_date, to_date, team_id, driver_id) VALUES ($1, $2, $3, $4) RETURNING id;"
		return teamDriverContractRepository.store.db.Get(
			&contract.ID,
			createQuery,
			contract.FromDate,
			contract.ToDate,
			contract.Team.ID,
			contract.Driver.ID,
		)
	} else {
		createQuery := "INSERT INTO team_driver_contracts (from_date, team_id, driver_id) VALUES ($1, $2, $3) RETURNING id;"
		return teamDriverContractRepository.store.db.Get(
			&contract.ID,
			createQuery,
			contract.FromDate,
			contract.Team.ID,
			contract.Driver.ID,
		)
	}
}
func (teamDriverContractRepository *TeamDriverContractRepository) Find(id int) (*model.TeamDriverContract, error) {
	return teamDriverContractRepository.FindBy("id", id, "")
}

func (teamDriverContractRepository *TeamDriverContractRepository) FindBy(columnName string, value interface{}, condition string) (*model.TeamDriverContract, error) {
	contract := &model.TeamDriverContract{}

	findQuery := fmt.Sprintf("SELECT "+
		"team_driver_contracts.id AS id, "+
		"team_driver_contracts.from_date AS from_date, "+
		"team_driver_contracts.to_date AS to_date, "+
		"teams.id AS \"team.id\", "+
		"teams.name AS \"team.name\", "+
		"teams.engine_manufacturer AS \"team.engine_manufacturer\", "+
		"drivers.id AS \"driver.id\", "+
		"drivers.first_name AS \"driver.first_name\", "+
		"drivers.last_name AS \"driver.last_name\", "+
		"drivers.birth_date AS \"driver.birth_date\" "+
		"FROM team_driver_contracts "+
		"LEFT JOIN teams "+
		"ON (team_driver_contracts.team_id = teams.id) "+
		"LEFT JOIN drivers "+
		"ON (team_driver_contracts.driver_id = drivers.id) "+
		"WHERE team_driver_contracts.%s = $1%s LIMIT 1;", columnName, condition)

	if err := teamDriverContractRepository.store.db.Get(
		contract,
		findQuery,
		value,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrNotFound
		}

		return nil, err
	}
	return contract, nil
}

func (teamDriverContractRepository *TeamDriverContractRepository) FindAllBy(columnName string, value interface{}, condition string) ([]*model.TeamDriverContract, error) {
	contracts := []*model.TeamDriverContract{}

	findQuery := fmt.Sprintf("SELECT "+
		"team_driver_contracts.id AS id, "+
		"team_driver_contracts.from_date AS from_date, "+
		"team_driver_contracts.to_date AS to_date, "+
		"teams.id AS \"team.id\", "+
		"teams.name AS \"team.name\", "+
		"teams.engine_manufacturer AS \"team.engine_manufacturer\", "+
		"drivers.id AS \"driver.id\", "+
		"drivers.first_name AS \"driver.first_name\", "+
		"drivers.last_name AS \"driver.last_name\", "+
		"drivers.birth_date AS \"driver.birth_date\" "+
		"FROM team_driver_contracts "+
		"LEFT JOIN teams "+
		"ON (team_driver_contracts.team_id = teams.id) "+
		"LEFT JOIN drivers "+
		"ON (team_driver_contracts.driver_id = drivers.id) "+
		"WHERE team_driver_contracts.%s = $1%s;", columnName, condition)

	if err := teamDriverContractRepository.store.db.Select(
		&contracts,
		findQuery,
		value,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrNotFound
		}

		return nil, err
	}
	return contracts, nil
}

func (teamDriverContractRepository *TeamDriverContractRepository) FindCurrentBy(columnName string,
	value interface{}) (*model.TeamDriverContract, error) {

	return teamDriverContractRepository.FindBy(columnName, value,
		"AND CURRENT_DATE BETWEEN team_driver_contracts.from_date AND team_driver_contracts.to_date")
}

func (teamDriverContractRepository *TeamDriverContractRepository) FindAllCurrentBy(columnName string, value interface{}) ([]*model.TeamDriverContract, error) {
	return teamDriverContractRepository.FindAllBy(columnName,
		value,
		"AND CURRENT_DATE BETWEEN team_driver_contracts.from_date AND team_driver_contracts.to_date")
}

func (teamDriverContractRepository *TeamDriverContractRepository) Update(contract *model.TeamDriverContract) error {
	updateTeamDriverContractQuery := "UPDATE team_driver_contracts " +
		"SET from_date = :from_date, " +
		"to_date = :to_date, " +
		"team_id = :team.id " +
		"driver_id = :driver.id, " +
		"WHERE id = :id;"

	countResult, countResultErr := teamDriverContractRepository.store.db.NamedExec(
		updateTeamDriverContractQuery,
		contract,
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

func (teamDriverContractRepository *TeamDriverContractRepository) Delete(id int) error {
	deleteTeamDriverContractQuery := "DELETE FROM team_driver_contracts WHERE id = $1;"

	countResult, countResultErr := teamDriverContractRepository.store.db.Exec(
		deleteTeamDriverContractQuery,
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
