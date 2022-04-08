package sqlstore

import "github.com/jmoiron/sqlx"

func (store *Store) createTables() error {
	tx, err := store.db.Beginx()
	if err != nil {
		return err
	}

	if err := createTableUsers(tx); err != nil {
		return err
	}

	if err := createTableTeams(tx); err != nil {
		return err
	}

	if err := createTableDrivers(tx); err != nil {
		return err
	}

	if err := createTableRaces(tx); err != nil {
		return err
	}

	if err := createTableTeamDriverContracts(tx); err != nil {
		return err
	}

	if err := createTableRaceResults(tx); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func createTableUsers(tx *sqlx.Tx) error {
	createTableUsersQuery := "CREATE TABLE IF NOT EXISTS users (" +
		"id bigserial NOT NULL PRIMARY KEY," +
		"username varchar NOT NULL UNIQUE," +
		"email varchar NOT NULL UNIQUE," +
		"encrypted_password varchar NOT NULL );"

	if _, err := tx.Exec(createTableUsersQuery); err != nil {
		return err
	}

	return nil
}

func createTableTeams(tx *sqlx.Tx) error {
	createTableTeamsQuery := "CREATE TABLE IF NOT EXISTS teams (" +
		"id bigserial NOT NULL PRIMARY KEY," +
		"name varchar NOT NULL UNIQUE," +
		"engine_manufacturer varchar NOT NULL );"

	if _, err := tx.Exec(createTableTeamsQuery); err != nil {
		return err
	}

	return nil
}

func createTableDrivers(tx *sqlx.Tx) error {
	createTableDriversQuery := "CREATE TABLE IF NOT EXISTS drivers (" +
		"id bigserial NOT NULL PRIMARY KEY," +
		"first_name varchar NOT NULL," +
		"last_name varchar NOT NULL," +
		"birth_date date NOT NULL );"

	if _, err := tx.Exec(createTableDriversQuery); err != nil {
		return err
	}

	return nil
}

func createTableRaces(tx *sqlx.Tx) error {
	createTableRacesQuery := "CREATE TABLE IF NOT EXISTS races (" +
		"id bigserial NOT NULL PRIMARY KEY," +
		"name varchar NOT NULL UNIQUE," +
		"location varchar NOT NULL," +
		"date date NOT NULL );"

	if _, err := tx.Exec(createTableRacesQuery); err != nil {
		return err
	}

	return nil
}

func createTableTeamDriverContracts(tx *sqlx.Tx) error {
	createTableTeamDriverContractsQuery := "CREATE TABLE IF NOT EXISTS team_driver_contracts (" +
		"id bigserial NOT NULL PRIMARY KEY," +
		"from_date date NOT NULL," +
		"to_date date," +
		"team_id bigserial NOT NULL REFERENCES teams (id) ON DELETE CASCADE," +
		"driver_id bigserial NOT NULL REFERENCES drivers (id) ON DELETE CASCADE );"

	if _, err := tx.Exec(createTableTeamDriverContractsQuery); err != nil {
		return err
	}

	return nil
}

func createTableRaceResults(tx *sqlx.Tx) error {
	createTableRaceResultsQuery := "CREATE TABLE IF NOT EXISTS race_results (" +
		"id bigserial NOT NULL PRIMARY KEY," +
		"place integer NOT NULL," +
		"points integer NOT NULL DEFAULT 0," +
		"race_id bigserial NOT NULL REFERENCES races (id) ON DELETE CASCADE," +
		"driver_id bigserial NOT NULL REFERENCES drivers (id) ON DELETE CASCADE );"

	if _, err := tx.Exec(createTableRaceResultsQuery); err != nil {
		return err
	}

	return nil
}
