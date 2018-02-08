package binance

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"gotrading/core"
	"gotrading/networking"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	jsoniter "github.com/json-iterator/go"
)

const (
	hostURL      = "https://api.binance.com/api"
	exchangeInfo = "exchangeInfo"
	depth        = "depth"
	account      = "account"
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

		url := fmt.Sprintf("%s/v1/%s", hostURL, exchangeInfo)

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
		var err error

		type Response struct {
			Asks [][2]string `json:"asks"`
			Bids [][2]string `json:"bids"`
		}

		response := Response{}
		endpoint := hit.Endpoint
		dst := &core.Orderbook{}
		dst.CurrencyPair = core.CurrencyPair{Base: endpoint.From, Quote: endpoint.To}
		curr := strings.ToUpper(fmt.Sprintf("%s%s", endpoint.From, endpoint.To))

		depthValue := 5
		req := fmt.Sprintf("%s/v1/%s?symbol=%s&limit=%d", hostURL, depth, curr, depthValue)

		gatling := networking.SharedGatling()
		contents, err, start, end := gatling.GET(req)
		var json = jsoniter.ConfigCompatibleWithStandardLibrary
		err = json.Unmarshal(contents, &response)
		if err != nil {
			log.Println(string(contents[:]))
		}

		if err == nil {
			dst.Bids = make([]core.Order, depthValue)
			dst.Asks = make([]core.Order, depthValue)
			dst.StartedLastUpdateAt = start
			dst.EndedLastUpdateAt = end

			for i, ask := range response.Asks {
				p, _ := strconv.ParseFloat(ask[0], 64)
				v, _ := strconv.ParseFloat(ask[1], 64)
				dst.Asks[i] = core.NewAsk(p, v)
			}
			for i, bid := range response.Bids {
				p, _ := strconv.ParseFloat(bid[0], 64)
				v, _ := strconv.ParseFloat(bid[1], 64)
				dst.Bids[i] = core.NewBid(p, v)
			}
		} else {
			fmt.Println("Error", endpoint.Description(), err)
		}
		return *dst, err
	}
}

func (b Binance) GetPortfolio() func(settings core.ExchangeSettings) (core.Portfolio, error) {
	return func(settings core.ExchangeSettings) (core.Portfolio, error) {
		portfolio := core.Portfolio{}

		timestamp := time.Now().Unix() * 1000

		values := url.Values{}
		values.Set("timestamp", fmt.Sprintf("%d", timestamp))

		mac := hmac.New(sha256.New, []byte(settings.APISecret))
		_, err := mac.Write([]byte(values.Encode()))
		if err != nil {
			return portfolio, err
		}
		signature := hex.EncodeToString(mac.Sum(nil))

		url := fmt.Sprintf("%s/v3/%s?%s&signature=%s", hostURL, account, values.Encode(), signature)
		req, err := http.NewRequest("GET", url, nil)

		if err != nil {
			return portfolio, err
		}

		req.Header.Add("X-MBX-APIKEY", settings.APIKey)

		gatling := networking.SharedGatling()
		contents, err, _, _ := gatling.Send(req)

		if err != nil {
			log.Println(err)
		}

		type Balamce struct {
			Asset  string `json:"asset"`
			Free   string `json:"free"`
			Locked string `json:"locked"`
		}

		type Response struct {
			MakerCommission  int       `json:"makerCommission"`
			TakerCommission  int       `json:"takerCommission"`
			BuyerCommission  int       `json:"buyerCommission"`
			SellerCommission int       `json:"sellerCommission"`
			CanTrade         bool      `json:"canTrade"`
			CanWithdraw      bool      `json:"canWithdraw"`
			UpdateTime       int       `json:"updateTime"`
			Balances         []Balamce `json:"balances"`
		}

		response := Response{}
		var json = jsoniter.ConfigCompatibleWithStandardLibrary
		err = json.Unmarshal(contents, &response)
		if err != nil {
			log.Println(err)
		} else {
			log.Println(response)
		}

		balances := response.Balances
		for _, balance := range balances {
			curr := core.Currency(balance.Asset)
			position, _ := strconv.ParseFloat(balance.Free, 64)
			portfolio.UpdatePosition(settings.Name, curr, position)
		}

		curr := core.Currency("BTC")
		position := 100.0
		portfolio.UpdatePosition(settings.Name, curr, position)

		return portfolio, err
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
