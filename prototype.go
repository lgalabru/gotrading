package main

import (
	"fmt"
	"net/http"
	"strings"

	"gotrading/exchanges"

	"github.com/spf13/viper"
)

func main() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	arbitrage := viper.GetStringMapString("strategies.arbitrage")

	exchangesEnabled := strings.Split(arbitrage["exchanges_enabled"], ",")

	factory := exchanges.Factory{}
	for _, name := range exchangesEnabled {
		exchange := factory.BuildExchange(name)
		exchange.GetPortfolio(http.Client{})
	}

}
