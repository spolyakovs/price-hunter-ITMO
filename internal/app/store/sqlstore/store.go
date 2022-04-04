package sqlstore

import (
	"github.com/jmoiron/sqlx"
	"github.com/spolyakovs/price-hunter-ITMO/internal/app/store"
)

var (
	pointsByPlace = map[int]int{
		1:  25,
		2:  18,
		3:  15,
		4:  12,
		5:  10,
		6:  8,
		7:  6,
		8:  4,
		9:  2,
		10: 1,
	}
)

type Store struct {
	db                           *sqlx.DB
	userRepository               *UserRepository
	teamRepository               *TeamRepository
	driverRepository             *DriverRepository
	raceRepository               *RaceRepository
	teamDriverContractRepository *TeamDriverContractRepository
	raceResultRepository         *RaceResultRepository
}

func New(db *sqlx.DB) (*Store, error) {
	// points after 10th place are not awarded
	for i := 11; i <= 20; i++ {
		pointsByPlace[i] = 0
	}

	newStore := &Store{
		db: db,
	}

	if err := newStore.createTables(); err != nil {
		return nil, err
	}

	// if err := newStore.fillTables(); err != nil {
	// 	return nil, err
	// }

	return newStore, nil
}

func (st *Store) Users() store.UserRepository {
	if st.userRepository != nil {
		return st.userRepository
	}

	st.userRepository = &UserRepository{
		store: st,
	}

	return st.userRepository
}

func (st *Store) Teams() store.TeamRepository {
	if st.teamRepository != nil {
		return st.teamRepository
	}

	st.teamRepository = &TeamRepository{
		store: st,
	}

	return st.teamRepository
}

func (st *Store) Drivers() store.DriverRepository {
	if st.driverRepository != nil {
		return st.driverRepository
	}

	st.driverRepository = &DriverRepository{
		store: st,
	}

	return st.driverRepository
}

func (st *Store) Races() store.RaceRepository {
	if st.raceRepository != nil {
		return st.raceRepository
	}

	st.raceRepository = &RaceRepository{
		store: st,
	}

	return st.raceRepository
}

func (st *Store) TeamDriverContracts() store.TeamDriverContractRepository {
	if st.teamDriverContractRepository != nil {
		return st.teamDriverContractRepository
	}

	st.teamDriverContractRepository = &TeamDriverContractRepository{
		store: st,
	}

	return st.teamDriverContractRepository
}

func (st *Store) RaceResults() store.RaceResultRepository {
	if st.raceResultRepository != nil {
		return st.raceResultRepository
	}

	st.raceResultRepository = &RaceResultRepository{
		store: st,
	}

	return st.raceResultRepository
}
