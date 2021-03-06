package main

import (
	"flag"
	"log"

	"github.com/BurntSushi/toml"
	"github.com/spolyakovs/price-hunter-ITMO/internal/app/apiserver"
)

var (
	configPath string
)

func init() {
	flag.StringVar(&configPath, "config-path", "configs/local.toml", "path to config file")
}

// TODO: think about multiple ports (with groupcache, simple example in bookmarks)

func main() {
	flag.Parse()

	config := apiserver.NewConfig()
	_, err := toml.DecodeFile(configPath, config)
	if err != nil {
		log.Fatal(err)
	}

	if err := apiserver.Start(config); err != nil {
		log.Fatal(err)
	}
}
