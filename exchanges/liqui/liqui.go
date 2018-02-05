package liqui

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"math"

	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"gotrading/core"
	"gotrading/networking"

	"github.com/json-iterator/go"
)

const (
	liquiAPIPublicURL  = "https://api.Liqui.io/api/3"
	liquiAPIPrivateURL = "https://api.Liqui.io/tapi"
	liquiInfo          = "info"
	liquiTicker        = "ticker"
	liquiDepth         = "depth"
	liquiTrades        = "trades"
	liquiGetInfo       = "getInfo"
	liquiTrade         = "Trade"
	liquiActiveOrders  = "ActiveOrders"
	liquiOrderInfo     = "OrderInfo"
	liquiCancelOrder   = "CancelOrder"
	liquiTradeHistory  = "TradeHistory"
	liquiWithdrawCoin  = "WithdrawCoin"
)

type Liqui struct {
}

func (b Liqui) GetSettings() func() (core.ExchangeSettings, error) {
	return func() (core.ExchangeSettings, error) {

		type pairsSettings struct {
			DecimalPlaces int     `json:"decimal_places"`
			MinPrice      float64 `json:"min_price"`
			MaxPrice      float64 `json:"max_price"`
			MinAmount     float64 `json:"min_amount"`
			MaxAmount     float64 `json:"max_amount"`
			MinTotal      float64 `json:"min_total"`
			Fee           float64 `json:"fee"`
		}

		type Response struct {
			ServerTime    int                      `json:"server_time"`
			PairsSettings map[string]pairsSettings `json:"pairs"`
		}

		response := Response{}
		settings := core.ExchangeSettings{}
		gatling := networking.SharedGatling()

		url := fmt.Sprintf("%s/%s", liquiAPIPublicURL, liquiInfo)

		contents, err, _, _ := gatling.GET(url)
		var json = jsoniter.ConfigCompatibleWithStandardLibrary
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

			cps := core.CurrencyPairSettings{}
			cps.BasePrecision = response.PairsSettings[key].DecimalPlaces
			cps.QuotePrecision = response.PairsSettings[key].DecimalPlaces
			cps.MinAmount = response.PairsSettings[key].MinAmount
			cps.MaxAmount = response.PairsSettings[key].MaxAmount
			cps.MinPrice = response.PairsSettings[key].MinPrice
			cps.MaxPrice = response.PairsSettings[key].MaxPrice
			settings.PairsSettings[pair] = cps
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
		dst.CurrencyPair = core.CurrencyPair{Base: endpoint.From, Quote: endpoint.To}
		curr := strings.ToLower(fmt.Sprintf("%s_%s", endpoint.From, endpoint.To))

		depth := 3
		req := fmt.Sprintf("%s/%s/%s?limit=%d", liquiAPIPublicURL, liquiDepth, curr, depth)

		gatling := networking.SharedGatling()
		contents, err, start, end := gatling.GET(req)
		var json = jsoniter.ConfigCompatibleWithStandardLibrary
		err = json.Unmarshal(contents, &response.Orderbook)
		if err != nil {
			log.Println(string(contents[:]))
		}
		src := response.Orderbook[curr]

		if err == nil {
			dst.Bids = make([]core.Order, depth)
			dst.Asks = make([]core.Order, depth)
			dst.StartedLastUpdateAt = start
			dst.EndedLastUpdateAt = end

			for i, ask := range src.Asks {
				dst.Asks[i] = core.NewAsk(ask[0], ask[1])
			}
			for i, bid := range src.Bids {
				dst.Bids[i] = core.NewBid(bid[0], bid[1])
			}
		} else {
			fmt.Println("Error", endpoint.Description(), err)
		}
		return *dst, err
	}
}

func (b Liqui) PostOrder() func(order core.Order, settings core.ExchangeSettings) (core.OrderDispatched, error) {
	return func(order core.Order, settings core.ExchangeSettings) (core.OrderDispatched, error) {
		var err error

		do := core.OrderDispatched{}
		do.Order = &order

		endpoint := order.Hit.Endpoint
		from := string(endpoint.From)
		to := string(endpoint.To)
		remotePair := strings.ToLower(from + "_" + to)
		pair := core.CurrencyPair{endpoint.From, endpoint.To}

		var orderType string
		var amount float64
		decimals := float64(settings.PairsSettings[pair].BasePrecision)

		if order.TransactionType == core.Ask {
			orderType = "sell"
			amount = order.BaseVolumeIn
		} else {
			orderType = "buy"
			amount = order.QuoteVolumeIn / order.Price
		}

		rate := order.Price
		amount = math.Trunc(amount*math.Pow(10, decimals)) / math.Pow(10, decimals)

		nonce := int(settings.Nonce.GetInc())

		values := url.Values{}
		values.Set("method", "Trade")
		values.Set("nonce", strconv.Itoa(nonce))
		values.Set("pair", remotePair)
		values.Set("type", orderType)
		values.Set("rate", strconv.FormatFloat(rate, 'f', -1, 64))
		values.Set("amount", strconv.FormatFloat(amount, 'f', int(decimals), 64))
		encoded := values.Encode()
		fmt.Println("Executing order:", encoded)

		h := hmac.New(sha512.New, []byte(settings.APISecret))
		h.Write([]byte(encoded))
		hmac := hex.EncodeToString(h.Sum(nil))

		headers := make(map[string]string)
		headers["Key"] = settings.APIKey
		headers["Sign"] = hmac
		headers["Content-Type"] = "application/x-www-form-urlencoded"

		req, err := http.NewRequest("POST", liquiAPIPrivateURL, strings.NewReader(encoded))

		if err != nil {
			return do, err
		}
		for k, v := range headers {
			req.Header.Add(k, v)
		}
		gatling := networking.SharedGatling()
		contents, err, start, end := gatling.Send(req)

		do.SentAt = start
		do.ConfirmedAt = end

		type Return struct {
			Received    float64            `json:"received"`
			Remains     float64            `json:"remains"`
			OrderID     int                `json:"order_id"`
			InitOrderID int                `json:"init_order_id"`
			Funds       map[string]float64 `json:"funds"`
		}

		type Response struct {
			Return Return `json:"return"`
		}

		response := Response{}
		var json = jsoniter.ConfigCompatibleWithStandardLibrary
		err = json.Unmarshal(contents, &response)
		if err != nil {
			log.Println(string(contents[:]))
		}
		fmt.Println(string(contents[:]))

		order.Progress = response.Return.Received / amount
		funds := response.Return.Funds

		manager := core.SharedPortfolioManager()
		state := manager.ForkCurrentState()
		for curr := range funds {
			state.UpdatePosition(settings.Name, core.Currency(strings.ToUpper(curr)), funds[curr])
		}
		manager.PushState(state)

		return do, err
	}
}

func (b Liqui) GetPortfolio() func(settings core.ExchangeSettings) (core.Portfolio, error) {
	return func(settings core.ExchangeSettings) (core.Portfolio, error) {
		portfolio := core.Portfolio{}
		var err error
		fmt.Println("Getting Portfolio from Liqui")

		nonce := int(settings.Nonce.GetInc())

		values := url.Values{}
		values.Set("method", liquiGetInfo)
		values.Set("nonce", strconv.Itoa(nonce))
		encoded := values.Encode()

		h := hmac.New(sha512.New, []byte(settings.APISecret))
		h.Write([]byte(encoded))
		hmac := hex.EncodeToString(h.Sum(nil))

		headers := make(map[string]string)
		headers["Key"] = settings.APIKey
		headers["Sign"] = hmac
		headers["Content-Type"] = "application/x-www-form-urlencoded"

		req, err := http.NewRequest("POST", liquiAPIPrivateURL, strings.NewReader(encoded))

		if err != nil {
			return portfolio, err
		}
		for k, v := range headers {
			req.Header.Add(k, v)
		}
		gatling := networking.SharedGatling()
		contents, err, _, _ := gatling.Send(req)

		if err != nil {
			log.Println(err)
		} else {
			log.Println(contents)
		}

		type Return struct {
			Funds map[string]float64 `json:"funds"`
		}

		type Response struct {
			Sucess int    `json:"success"`
			Error  string `json:"error"`
			Return Return `json:"return"`
		}

		response := Response{}
		var json = jsoniter.ConfigCompatibleWithStandardLibrary
		err = json.Unmarshal(contents, &response)
		if err != nil {
			log.Println(err)
		} else {
			log.Println(response)
		}

		funds := response.Return.Funds
		fmt.Println("->", funds)
		for curr := range funds {
			portfolio.UpdatePosition(settings.Name, core.Currency(strings.ToUpper(curr)), funds[curr])
		}

		return portfolio, err
	}
}
