package liqui

import (
	"fmt"
	"strings"
	"time"

	"gotrading/core"
	"gotrading/networking"
)

type Liqui struct {
}

func (b Liqui) GetOrderbook() func(hit core.Hit) (core.Orderbook, error) {
	return func(hit core.Hit) (core.Orderbook, error) {

		fmt.Println("FETCHING ORDERBOOK", hit)
		endpoint := hit.Endpoint
		exch := endpoint.Exchange

		type Orderbook struct {
			Asks [][]float64 `json:"asks"`
			Bids [][]float64 `json:"bids"`
		}
		type Response struct {
			Data map[string]Orderbook
		}

		response := Response{}
		curr := strings.ToLower(fmt.Sprintf("%s_%s", endpoint.From, endpoint.To))

		req := fmt.Sprintf("%s/%s/%s/%s?limit=5", "https://api.Liqui.io/api", "3", "depth", curr)

		t1 := time.Now()

		gatling := networking.SharedGatling()
		err := gatling.GET(req, &response.Data)
		t2 := time.Now()
		src := response.Data[curr]

		dst := &core.Orderbook{}

		if err == nil {
			dst.CurrencyPair = core.CurrencyPair{endpoint.From, endpoint.To}
			dst.Bids = make([]core.Order, 0)
			dst.Asks = make([]core.Order, 0)
			dst.StartedLastUpdateAt = t1
			dst.EndedLastUpdateAt = t2

			for _, ask := range src.Asks {
				if exch.IsCurrencyPairNormalized == true {
					dst.Asks = append(dst.Asks, core.NewAsk(dst.CurrencyPair, ask[0], ask[1]))
				} else {
				}
			}
			for _, bid := range src.Bids {
				if exch.IsCurrencyPairNormalized == true {
					dst.Bids = append(dst.Bids, core.NewBid(dst.CurrencyPair, bid[0], bid[1]))
				} else {
				}
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
