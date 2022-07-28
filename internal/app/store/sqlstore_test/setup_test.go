package sqlstore_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/BurntSushi/toml"
	"github.com/spolyakovs/price-hunter-ITMO/internal/app/apiserver"
	"github.com/spolyakovs/price-hunter-ITMO/internal/app/store/sqlstore"
)

var st *sqlstore.Store

func TestMain(m *testing.M) {
	config := apiserver.NewConfig()

	if _, err := toml.DecodeFile("../../../../configs/local_test.toml", config); err != nil {
		fmt.Printf("Couldn't get config:\n\t%s", err.Error())
		return
	}

	db, err := apiserver.NewDB(*config)
	if err != nil {
		fmt.Printf("Couldn't initalize DB:\n\t%s", err.Error())
		return
	}

	st, err = sqlstore.New(db)
	if err != nil {
		fmt.Printf("Couldn't initalize SQLStore:\n\t%s", err.Error())
		return
	}

	if err := st.InsertTestData(); err != nil {
		fmt.Printf("Couldn't insert test data into DB:\n\t%s", err.Error())
		return
	}

	os.Exit(m.Run())
}
