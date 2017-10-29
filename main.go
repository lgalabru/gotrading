package main

import (
	"log"

	"github.com/thrasher-/gocryptotrader/config"
)

func main() {

	cfg := &config.Cfg
	err := cfg.LoadConfig("config.dat")
	if err != nil {
		log.Fatal(err)
	}
}
