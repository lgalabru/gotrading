package liqui

import (
	"errors"
	"fmt"
	"log"
	"net/url"
	"strconv"
	"strings"
	"time"

	"gotrading/exchanges"
	"gotrading/exchanges/ticker"

	"github.com/thrasher-/gocryptotrader/common"
	"github.com/thrasher-/gocryptotrader/config"
)

const (
	binanceAPIURL           = "https://api.binance.com/api/"
	binanceAPIVersion       = "1"
	binanceInfo             = "v1/exchangeInfo"
	binanceTicker           = "v1/ticker"
	binanceDepth            = "v1/depth"
	binanceTrades           = "v1/trades"
	binanceTicker24H        = "v1/ticker/24hr"
	binanceTickerPrice      = "v1/ticker/price"
	binanceTickerBookTicker = "v1/ticker/bookTicker"
	binanceCreateOrder      = "v3/order"
	binanceGetOrder         = "v3/order"
	binanceDeleteOrder      = "v3/order"
	binanceOpenOrders       = "v3/openOrders"
	binanceAllOrders        = "v3/allOrders"
	binanceAccountInfo      = "v3/account"
	binanceMyTrades         = "v3/myTrades"
)

// Liqui is the overarching type across the liqui package
type Binance struct {
	exchange.Base
	Ticker map[string]Ticker
	Info   Info
}

// OwvOvs24vGMBIIa0q1B0RuIRh7uRNUFdreN5XsxdQzCKzA0yBCDBDJwiMJCAlw4J
// Iif3t03a3JDYfqF05pPcrfoEOIHwByKOLggwDAGgdugIYdv4rvstbef9jSsa7pnA
// SetDefaults sets current default values for liqui
func (b *Binance) SetDefaults() {
	b.Name = "Binance"
	b.Enabled = false
	b.Fee = 0.25
	b.Verbose = false
	b.Websocket = false
	b.RESTPollingDelay = 10
	// b.Ticker = make(map[string]Ticker)
	b.RequestCurrencyPairFormat.Delimiter = "_"
	b.RequestCurrencyPairFormat.Uppercase = false
	b.RequestCurrencyPairFormat.Separator = "-"
	b.ConfigCurrencyPairFormat.Delimiter = "_"
	b.ConfigCurrencyPairFormat.Uppercase = true
	b.AssetTypes = []string{ticker.Spot}
}

// Setup sets exchange configuration parameters for liqui
func (b *Binance) Setup(exch config.ExchangeConfig) {
	if !exch.Enabled {
		b.SetEnabled(false)
	} else {
		b.Enabled = true
		b.AuthenticatedAPISupport = exch.AuthenticatedAPISupport
		b.SetAPIKeys(exch.APIKey, exch.APISecret, "", false)
		b.RESTPollingDelay = exch.RESTPollingDelay
		b.Verbose = exch.Verbose
		b.Websocket = exch.Websocket
		b.BaseCurrencies = common.SplitStrings(exch.BaseCurrencies, ",")
		b.AvailablePairs = common.SplitStrings(exch.AvailablePairs, ",")
		b.EnabledPairs = common.SplitStrings(exch.EnabledPairs, ",")
		err := b.SetCurrencyPairFormat()
		if err != nil {
			log.Fatal(err)
		}
		err = b.SetAssetTypes()
		if err != nil {
			log.Fatal(err)
		}
	}
}

// GetFee returns a fee for a specific currency
func (b *Binance) GetFee(currency string) (float64, error) {
	log.Println(b.Info.Pairs)
	val, ok := b.Info.Pairs[common.StringToLower(currency)]
	if !ok {
		return 0, errors.New("currency does not exist")
	}

	return val.Fee, nil
}

// GetAvailablePairs returns all available pairs
func (b *Binance) GetAvailablePairs(nonHidden bool) []string {
	var pairs []string
	for x, y := range b.Info.Pairs {
		if nonHidden && y.Hidden == 1 || x == "" {
			continue
		}
		pairs = append(pairs, common.StringToUpper(x))
	}
	return pairs
}

// GetInfo provides all the information about currently active pairs, such as
// the maximum number of digits after the decimal point, the minimum price, the
// maximum price, the minimum transaction size, whether the pair is hidden, the
// commission for each pair.
func (b *Binance) GetInfo() (Info, error) {
	resp := Info{}
	req := fmt.Sprintf("%s/%s/%s/", liquiAPIPublicURL, liquiAPIPublicVersion, liquiInfo)

	return resp, common.SendHTTPGetRequest(req, true, b.Verbose, &resp)
}

// GetTicker returns information about currently active pairs, such as: the
// maximum price, the minimum price, average price, trade volume, trade volume
// in currency, the last trade, Buy and Sell price. All information is provided
// over the past 24 hours.
//
// currencyPair - example "eth_btc"
func (b *Binance) GetTicker(currencyPair string) (map[string]Ticker, error) {
	type Response struct {
		Data map[string]Ticker
	}

	response := Response{}
	req := fmt.Sprintf("%s/%s/%s/%s", liquiAPIPublicURL, liquiAPIPublicVersion, liquiTicker, currencyPair)

	return response.Data,
		common.SendHTTPGetRequest(req, true, b.Verbose, &response.Data)
}

// GetDepth information about active orders on the pair. Additionally it accepts
// an optional GET-parameter limit, which indicates how many orders should be
// displayed (150 by default). Is set to less than 2000.
func (b *Binance) GetDepth(currencyPair string) (Orderbook, error) {
	type Response struct {
		Data map[string]Orderbook
	}

	response := Response{}
	req := fmt.Sprintf("%s/%s/%s/%s", liquiAPIPublicURL, liquiAPIPublicVersion, liquiDepth, currencyPair)

	return response.Data[currencyPair],
		common.SendHTTPGetRequest(req, true, b.Verbose, &response.Data)
}

// GetTrades returns information about the last trades. Additionally it accepts
// an optional GET-parameter limit, which indicates how many orders should be
// displayed (150 by default). The maximum allowable value is 2000.
func (b *Binance) GetTrades(currencyPair string) ([]Trades, error) {
	type Response struct {
		Data map[string][]Trades
	}

	response := Response{}
	req := fmt.Sprintf("%s/%s/%s/%s", liquiAPIPublicURL, liquiAPIPublicVersion, liquiTrades, currencyPair)

	return response.Data[currencyPair],
		common.SendHTTPGetRequest(req, true, b.Verbose, &response.Data)
}

// GetAccountInfo returns information about the user’s current balance, API-key
// privileges, the number of open orders and Server Time. To use this method you
// need a privilege of the key info.
func (b *Binance) GetAccountInfo() (AccountInfo, error) {
	var result AccountInfo

	return result,
		b.SendAuthenticatedHTTPRequest(liquiAccountInfo, url.Values{}, &result)
}

// Trade creates orders on the exchange.
// to-do: convert orderid to int64
func (b *Binance) Trade(pair, orderType string, amount, price float64) (float64, error) {
	req := url.Values{}
	req.Add("pair", pair)
	req.Add("type", orderType)
	req.Add("amount", strconv.FormatFloat(amount, 'f', -1, 64))
	req.Add("rate", strconv.FormatFloat(price, 'f', -1, 64))

	var result Trade

	return result.OrderID, b.SendAuthenticatedHTTPRequest(liquiTrade, req, &result)
}

// GetActiveOrders returns the list of your active orders.
func (b *Binance) GetActiveOrders(pair string) (map[string]ActiveOrders, error) {
	req := url.Values{}
	req.Add("pair", pair)

	var result map[string]ActiveOrders
	return result, b.SendAuthenticatedHTTPRequest(liquiActiveOrders, req, &result)
}

// GetOrderInfo returns the information on particular order.
func (b *Binance) GetOrderInfo(OrderID int64) (map[string]OrderInfo, error) {
	req := url.Values{}
	req.Add("order_id", strconv.FormatInt(OrderID, 10))

	var result map[string]OrderInfo
	return result, b.SendAuthenticatedHTTPRequest(liquiOrderInfo, req, &result)
}

// CancelOrder method is used for order cancelation.
func (b *Binance) CancelOrder(OrderID int64) (bool, error) {
	req := url.Values{}
	req.Add("order_id", strconv.FormatInt(OrderID, 10))

	var result CancelOrder

	err := b.SendAuthenticatedHTTPRequest(liquiCancelOrder, req, &result)
	if err != nil {
		return false, err
	}

	return true, nil
}

// GetTradeHistory returns trade history
func (b *Binance) GetTradeHistory(vals url.Values, pair string) (map[string]TradeHistory, error) {
	if pair != "" {
		vals.Add("pair", pair)
	}

	var result map[string]TradeHistory
	return result, b.SendAuthenticatedHTTPRequest(liquiTradeHistory, vals, &result)
}

// WithdrawCoins is designed for cryptocurrency withdrawals.
// API mentions that this isn't active now, but will be soon - you must provide the first 8 characters of the key
// in your ticket to support.
func (b *Binance) WithdrawCoins(coin string, amount float64, address string) (WithdrawCoins, error) {
	req := url.Values{}
	req.Add("coinName", coin)
	req.Add("amount", strconv.FormatFloat(amount, 'f', -1, 64))
	req.Add("address", address)

	var result WithdrawCoins
	return result, b.SendAuthenticatedHTTPRequest(liquiWithdrawCoin, req, &result)
}

// SendAuthenticatedHTTPRequest sends an authenticated http request to liqui
func (b *Binance) SendAuthenticatedHTTPRequest(method string, values url.Values, result interface{}) (err error) {
	if !b.AuthenticatedAPISupport {
		return fmt.Errorf(exchange.WarningAuthenticatedRequestWithoutCredentialsSet, b.Name)
	}

	if b.Nonce.Get() == 0 {
		b.Nonce.Set(time.Now().Unix())
	} else {
		b.Nonce.Inc()
	}
	values.Set("nonce", b.Nonce.String())
	values.Set("method", method)

	encoded := values.Encode()
	hmac := common.GetHMAC(common.HashSHA512, []byte(encoded), []byte(b.APISecret))

	if b.Verbose {
		log.Printf("Sending POST request to %s calling method %s with params %s\n", liquiAPIPrivateURL, method, encoded)
	}

	headers := make(map[string]string)
	headers["Key"] = b.APIKey
	headers["Sign"] = common.HexEncodeToString(hmac)
	headers["Content-Type"] = "application/x-www-form-urlencoded"

	resp, err := common.SendHTTPRequest("POST", liquiAPIPrivateURL, headers, strings.NewReader(encoded))
	if err != nil {
		return err
	}

	response := Response{}

	err = common.JSONDecode([]byte(resp), &response)
	if err != nil {
		return err
	}

	if response.Success != 1 {
		return errors.New(response.Error)
	}

	jsonEncoded, err := common.JSONEncode(response.Return)
	if err != nil {
		return err
	}

	err = common.JSONDecode(jsonEncoded, &result)
	if err != nil {
		return err
	}

	return nil
}
