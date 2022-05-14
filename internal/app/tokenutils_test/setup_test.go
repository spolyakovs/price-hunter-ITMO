package tokenutils_test

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/spolyakovs/price-hunter-ITMO/internal/app/apiserver"
	"github.com/spolyakovs/price-hunter-ITMO/internal/app/tokenutils"
)

func setupRedis() error {
	config := apiserver.NewConfig()

	if _, err := toml.DecodeFile("../../../configs/local_test.toml", config); err != nil {
		return fmt.Errorf("Couldn't get config:\n\t%s", err.Error())
	}

	os.Setenv("TOKEN_SECRET", config.TokenSecret)

	if err := tokenutils.SetupRedis(config.RedisAddr); err != nil {
		return fmt.Errorf("Couldn't setup Redis:\n\t%s", err.Error())
	}

	return nil
}
