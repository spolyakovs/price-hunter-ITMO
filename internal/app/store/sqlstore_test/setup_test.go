package sqlstore_test

import (
	"fmt"

	"github.com/BurntSushi/toml"
	"github.com/spolyakovs/price-hunter-ITMO/internal/app/apiserver"
	"github.com/spolyakovs/price-hunter-ITMO/internal/app/store/sqlstore"
)

func setupStore() (*sqlstore.Store, error) {
	config := apiserver.NewConfig()

	if _, err := toml.DecodeFile("../../../../configs/local_test.toml", config); err != nil {
		return nil, fmt.Errorf("Couldn't get config:\n\t%s", err.Error())
	}

	db, err := apiserver.NewDB(*config)
	if err != nil {
		return nil, fmt.Errorf("Couldn't initalize DB:\n\t%s", err.Error())
	}

	store, err := sqlstore.New(db)
	if err != nil {
		return nil, fmt.Errorf("Couldn't initalize SQLStore:\n\t%s", err.Error())
	}

	if err := insertTestData(store); err != nil {
		return nil, fmt.Errorf("Couldn't insert test data into DB:\n\t%s", err.Error())
	}

	return store, nil
}
