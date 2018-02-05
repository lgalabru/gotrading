package binance

import (
	"errors"
	"fmt"
	"gotrading/core"
	"gotrading/networking"
	"strconv"

	jsoniter "github.com/json-iterator/go"
)

const (
	hostURL           = "https://api.binance.com/api/v1"
	exchangeInfo      = "exchangeInfo"
	liquiTicker       = "ticker"
	liquiDepth        = "depth"
	liquiTrades       = "trades"
	liquiGetInfo      = "getInfo"
	liquiTrade        = "Trade"
	liquiActiveOrders = "ActiveOrders"
	liquiOrderInfo    = "OrderInfo"
	liquiCancelOrder  = "CancelOrder"
	liquiTradeHistory = "TradeHistory"
	liquiWithdrawCoin = "WithdrawCoin"
)

type Binance struct {
}

func (b Binance) GetSettings() func() (core.ExchangeSettings, error) {
	return func() (core.ExchangeSettings, error) {

		type pairSettings struct {
			Symbol              string              `json:"symbol"`
			Status              string              `json:"status"`
			BaseAsset           string              `json:"baseAsset"`
			QuoteAsset          string              `json:"quoteAsset"`
			BaseAssetPrecision  int                 `json:"baseAssetPrecision"`
			QuoteAssetPrecision int                 `json:"quoteAssetPrecision"`
			Filters             []map[string]string `json:"filters"`
		}

		type Response struct {
			ServerTime int64          `json:"serverTime"`
			Symbols    []pairSettings `json:"symbols"`
		}

		response := Response{}
		settings := core.ExchangeSettings{}
		gatling := networking.SharedGatling()

		url := fmt.Sprintf("%s/%s", hostURL, exchangeInfo)

		contents, err, _, _ := gatling.GET(url)
		var json = jsoniter.ConfigCompatibleWithStandardLibrary
		err = json.Unmarshal(contents[:], &response)

		if len(response.Symbols) == 0 {
			return settings, errors.New("info empty")
		}

		settings.IsCurrencyPairNormalized = true
		settings.AvailablePairs = make([]core.CurrencyPair, len(response.Symbols))
		settings.PairsSettings = make(map[core.CurrencyPair]core.CurrencyPairSettings, len(response.Symbols))

		for i, sym := range response.Symbols {
			base := core.Currency(sym.BaseAsset)
			quote := core.Currency(sym.QuoteAsset)
			pair := core.CurrencyPair{Base: base, Quote: quote}
			settings.AvailablePairs[i] = pair

			cps := core.CurrencyPairSettings{}
			cps.BasePrecision = sym.BaseAssetPrecision
			cps.QuotePrecision = sym.QuoteAssetPrecision
			cps.MinAmount, _ = strconv.ParseFloat(sym.Filters[1]["minQty"], 64)
			cps.MaxAmount, _ = strconv.ParseFloat(sym.Filters[1]["maxQty"], 64)
			cps.MinPrice, _ = strconv.ParseFloat(sym.Filters[0]["minPrice"], 64)
			cps.MaxPrice, _ = strconv.ParseFloat(sym.Filters[0]["maxPrice"], 64)
			settings.PairsSettings[pair] = cps
		}
		return settings, err
	}
}

func (b Binance) GetOrderbook() func(hit core.Hit) (core.Orderbook, error) {
	return func(hit core.Hit) (core.Orderbook, error) {
		var ob core.Orderbook
		var err error
		fmt.Println("Getting Orderbooks from Binance")
		return ob, err
	}
}

func (b Binance) GetPortfolio() func(settings core.ExchangeSettings) (core.Portfolio, error) {
	return func(settings core.ExchangeSettings) (core.Portfolio, error) {
		var p core.Portfolio
		var err error
		fmt.Println("Getting Portfolio from Binance")
		return p, err
	}
}

func (b Binance) PostOrder() func(order core.Order, settings core.ExchangeSettings) (core.OrderDispatched, error) {
	return func(order core.Order, settings core.ExchangeSettings) (core.OrderDispatched, error) {
		var o core.OrderDispatched
		var err error
		fmt.Println("Posting Order on Binance")
		return o, err
	}
}

// func (b *Binance) Deposit(client http.Client) (bool, error) {
// 	var err error
// 	return true, err
// }

// func (b *Binance) Withdraw(client http.Client) (bool, error) {
// 	var err error
// 	return true, err
// }
