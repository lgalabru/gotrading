package liqui

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"gotrading/core"
	"gotrading/networking"
)

type Liqui struct {
}

func (b Liqui) GetSettings() func() (core.ExchangeSettings, error) {
	return func() (core.ExchangeSettings, error) {

		type Response struct {
			ServerTime    int                                  `json:"server_time"`
			PairsSettings map[string]core.CurrencyPairSettings `json:"pairs"`
		}

		response := Response{}
		settings := core.ExchangeSettings{}
		gatling := networking.SharedGatling()
		contents, err := gatling.GET("https://api.liqui.io/api/3/info")
		err = json.Unmarshal(contents[:], &response)

		if len(response.PairsSettings) == 0 {
			return settings, errors.New("info empty")
		}

		settings.IsCurrencyPairNormalized = true
		settings.AvailablePairs = make([]core.CurrencyPair, len(response.PairsSettings))
		settings.PairsSettings = make(map[core.CurrencyPair]core.CurrencyPairSettings, len(response.PairsSettings))

		i := 0
		for key := range response.PairsSettings {
			currs := strings.Split(strings.ToUpper(key), "_")
			base := core.Currency(currs[0])
			quote := core.Currency(currs[1])
			pair := core.CurrencyPair{Base: base, Quote: quote}
			settings.AvailablePairs[i] = pair
			settings.PairsSettings[pair] = response.PairsSettings[key]
			i++
		}
		return settings, err
	}
}

func (b Liqui) GetOrderbook() func(hit core.Hit) (core.Orderbook, error) {
	return func(hit core.Hit) (core.Orderbook, error) {

		type Response struct {
			Orderbook map[string]struct {
				Asks [][]float64 `json:"asks"`
				Bids [][]float64 `json:"bids"`
			}
		}

		response := Response{}
		endpoint := hit.Endpoint
		dst := &core.Orderbook{}
		curr := strings.ToLower(fmt.Sprintf("%s_%s", endpoint.From, endpoint.To))

		req := fmt.Sprintf("%s/%s/%s/%s?limit=3", "https://api.Liqui.io/api", "3", "depth", curr)

		start := time.Now()
		gatling := networking.SharedGatling()
		contents, err := gatling.GET(req)
		err = json.Unmarshal(contents, &response.Orderbook)
		if err != nil {
			log.Println(string(contents[:]))
		}
		end := time.Now()
		src := response.Orderbook[curr]

		if err == nil {
			dst.Bids = make([]core.Order, 0)
			dst.Asks = make([]core.Order, 0)
			dst.StartedLastUpdateAt = start
			dst.EndedLastUpdateAt = end

			for _, ask := range src.Asks {
				a := core.NewAsk(ask[0], ask[1])
				dst.Asks = append(dst.Asks, a)
			}
			for _, bid := range src.Bids {
				b := core.NewBid(bid[0], bid[1])
				dst.Bids = append(dst.Bids, b)
			}
		} else {
			fmt.Println("Error", endpoint.Description(), err)
		}
		return *dst, err
	}
}

func (b Liqui) GetPortfolio() func() (core.Portfolio, error) {
	return func() (core.Portfolio, error) {
		var p core.Portfolio
		var err error
		fmt.Println("Getting Portfolio from Liqui")
		return p, err
	}
}

func (b Liqui) PostOrder() func(order core.Order) (core.Order, error) {
	return func(order core.Order) (core.Order, error) {
		var o core.Order
		var err error
		fmt.Println("Posting Order on Liqui")

		// exchange := r.Path.Hits[i].Endpoint.Exchange
		// pair := strings.ToLower(string(r.Path.Hits[i].Endpoint.From)) + "_" + strings.ToLower(string(r.Path.Hits[i].Endpoint.To))
		// var orderType string
		// var amount float64

		// if o.TransactionType == core.Ask {
		// 	orderType = "sell"
		// 	amount = o.BaseVolumeIn
		// } else {
		// 	orderType = "buy"
		// 	amount = o.QuoteVolumeIn / o.Price
		// }
		// price := o.Price
		// // decimals := exec.chain.Path.Hits[i].Endpoint.Exchange.Liqui.Info.Pairs[pair].DecimalPlaces
		// decimals := 8
		// res, error := exchange.PostOrder(o)

		// // res, error := exchange.Trade(pair, orderType, toFixed(amount, decimals), price)
		// fmt.Println("Executing order:", pair, orderType, decimals, toFixed(amount, decimals), price, res, error)

		return o, err
	}
}

// func (g *Gatling) FetchUSD() float64 {
// 	cp := pair.NewCurrencyPair("BTC", "USDT")

// 	client := g.Clients[len(g.Clients)-1]

// 	type Orderbook struct {
// 		Asks [][]float64 `json:"asks"`
// 		Bids [][]float64 `json:"bids"`
// 	}
// 	type Response struct {
// 		Data map[string]Orderbook
// 	}

// 	response := Response{}
// 	curr := fmt.Sprintf("%s", cp.Display("_", false))

// 	req := fmt.Sprintf("%s/%s/%s/%s?limit=1", "https://api.Liqui.io/api", "3", "depth", curr)

// 	g.SendHTTPGetRequest(client, req, true, false, &response.Data)
// 	src := response.Data[curr]
// 	return src.Asks[0][0]
// }

// func (b *Binance) Deposit(client http.Client) (bool, error) {
// 	var err error
// 	return true, err
// }

// func (b *Binance) Withdraw(client http.Client) (bool, error) {
// 	var err error
// 	return true, err
// }
