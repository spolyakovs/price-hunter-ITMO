package tokenutils_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/BurntSushi/toml"
	"github.com/spolyakovs/price-hunter-ITMO/internal/app/apiserver"
	"github.com/spolyakovs/price-hunter-ITMO/internal/app/tokenutils"
)

func TestMain(m *testing.M) {
	config := apiserver.NewConfig()

	if _, err := toml.DecodeFile("../../../configs/local_test.toml", config); err != nil {
		fmt.Printf("Couldn't get config:\n\t%s", err.Error())
		return
	}

	os.Setenv("TOKEN_SECRET", config.TokenSecret)

	if err := tokenutils.SetupRedis(config.RedisAddr); err != nil {
		fmt.Printf("Couldn't setup Redis:\n\t%s", err.Error())
		return
	}

	os.Exit(m.Run())
}
