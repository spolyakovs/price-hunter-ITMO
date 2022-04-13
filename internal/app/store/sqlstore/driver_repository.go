package sqlstore

import (
	"database/sql"
	"fmt"

	"github.com/spolyakovs/price-hunter-ITMO/internal/app/model"
	"github.com/spolyakovs/price-hunter-ITMO/internal/app/store"
)

type DriverRepository struct {
	store *Store
}

func (driverRepository *DriverRepository) Create(driver *model.Driver) error {
	createQuery := "INSERT INTO drivers (first_name, last_name, birth_date) VALUES ($1, $2, $3) RETURNING id;"

	return driverRepository.store.db.Get(
		&driver.ID,
		createQuery,
		driver.FirstName, driver.LastName, driver.BirthDate,
	)
}

func (driverRepository *DriverRepository) Find(id int) (*model.Driver, error) {
	return driverRepository.FindBy("id", id)
}

func (driverRepository *DriverRepository) FindBy(columnName string, value interface{}) (*model.Driver, error) {
	driver := &model.Driver{}

	findQuery := fmt.Sprintf("SELECT * FROM drivers WHERE %s = $1 LIMIT 1;", columnName)
	if err := driverRepository.store.db.Get(
		driver,
		findQuery,
		value,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrNotFound
		}

		return nil, err
	}

	return driver, nil
}

func (driverRepository *DriverRepository) FindAllBy(columnName string, value interface{}) ([]*model.Driver, error) {
	drivers := []*model.Driver{}

	findQuery := fmt.Sprintf("SELECT * FROM drivers WHERE %s = $1;", columnName)
	if err := driverRepository.store.db.Select(
		&drivers,
		findQuery,
		value,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrNotFound
		}

		return nil, err
	}

	return drivers, nil
}

func (driverRepository *DriverRepository) GetAll() ([]*model.Driver, error) {
	drivers := []*model.Driver{}

	findQuery := "SELECT * FROM drivers;"
	if err := driverRepository.store.db.Select(
		&drivers,
		findQuery,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrNotFound
		}

		return nil, err
	}

	return drivers, nil
}

func (driverRepository *DriverRepository) Update(driver *model.Driver) error {
	updateQuery := "UPDATE drivers " +
		"SET first_name = :first_name, " +
		"last_name = :last_name, " +
		"birth_date = :birth_date " +
		"WHERE id = :id;"

	countResult, countResultErr := driverRepository.store.db.NamedExec(
		updateQuery,
		driver,
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

func (driverRepository *DriverRepository) Delete(id int) error {
	deleteQuery := "DELETE FROM drivers WHERE id = $1;"

	countResult, countResultErr := driverRepository.store.db.Exec(
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
