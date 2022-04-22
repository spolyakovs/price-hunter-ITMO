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
	db                  *sqlx.DB
	userRepository      *UserRepository
	publisherRepository *PublisherRepository
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

func (st *Store) Publishers() store.PublisherRepository {
	if st.publisherRepository != nil {
		return st.publisherRepository
	}

	st.publisherRepository = &PublisherRepository{
		store: st,
	}

	return st.publisherRepository
}
