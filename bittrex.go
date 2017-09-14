// Package Bittrex is an implementation of the Biitrex API in Golang.
package bittrex

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	API_BASE    = "https://bittrex.com/api/" // Bittrex API endpoint
	API_VERSION = "v1.1"                     // Bittrex API version
)

// New returns an instantiated bittrex struct
func New(apiKey, apiSecret string) *Bittrex {
	client := NewClient(apiKey, apiSecret)
	return &Bittrex{client}
}

// NewWithCustomHttpClient returns an instantiated bittrex struct with custom http client
func NewWithCustomHttpClient(apiKey, apiSecret string, httpClient *http.Client) *Bittrex {
	client := NewClientWithCustomHttpConfig(apiKey, apiSecret, httpClient)
	return &Bittrex{client}
}

// NewWithCustomTimeout returns an instantiated bittrex struct with custom timeout
func NewWithCustomTimeout(apiKey, apiSecret string, timeout time.Duration) *Bittrex {
	client := NewClientWithCustomTimeout(apiKey, apiSecret, timeout)
	return &Bittrex{client}
}

// handleErr gets JSON response from Bittrex API en deal with error
func handleErr(r jsonResponse) error {
	if !r.Success {
		return errors.New(r.Message)
	}
	return nil
}

// bittrex represent a bittrex client
type Bittrex struct {
	client *client
}

// GetDistribution is used to get the distribution.
func (b *Bittrex) GetDistribution(market string) (*Distribution, error) {
	r, err := b.client.do("GET", "https://bittrex.com/Api/v2.0/pub/currency/GetBalanceDistribution?currencyName="+strings.ToUpper(market), "", false)
	if err != nil {
		return nil, err
	}

	var response jsonResponse
	if err = json.Unmarshal(r, &response); err != nil {
		return nil, err
	}

	if err = handleErr(response); err != nil {
		return nil, err
	}

	distribution := Distribution{}
	if err = json.Unmarshal(response.Result, &distribution); err != nil {
		return nil, err
	}

	return &distribution, nil
}

// GetCurrencies is used to get all supported currencies at Bittrex along with other meta data.
func (b *Bittrex) GetCurrencies() ([]*Currency, error) {
	r, err := b.client.do("GET", "public/getcurrencies", "", false)
	if err != nil {
		return nil, err
	}
	var response jsonResponse
	if err = json.Unmarshal(r, &response); err != nil {
		return nil, err
	}
	if err = handleErr(response); err != nil {
		return nil, err
	}
	currencies := []*Currency{}
	if err = json.Unmarshal(response.Result, &currencies); err != nil {
		return nil, err
	}

	return currencies, nil
}

// GetMarkets is used to get the open and available trading markets at Bittrex along with other meta data.
func (b *Bittrex) GetMarkets() ([]*Market, error) {
	r, err := b.client.do("GET", "public/getmarkets", "", false)
	if err != nil {
		return nil, err
	}
	var response jsonResponse
	if err = json.Unmarshal(r, &response); err != nil {
		return nil, err
	}
	if err = handleErr(response); err != nil {
		return nil, err
	}

	markets := []*Market{}
	if err = json.Unmarshal(response.Result, &markets); err != nil {
		return nil, err
	}

	return markets, nil
}

// GetTicker is used to get the current ticker values for a market.
func (b *Bittrex) GetTicker(market string) (*Ticker, error) {
	r, err := b.client.do("GET", "public/getticker?market="+strings.ToUpper(market), "", false)
	if err != nil {
		return nil, err
	}
	var response jsonResponse
	if err = json.Unmarshal(r, &response); err != nil {
		return nil, err
	}
	if err = handleErr(response); err != nil {
		return nil, err
	}

	ticker := Ticker{}
	if err = json.Unmarshal(response.Result, &ticker); err != nil {
		return nil, err
	}
	return &ticker, nil
}

// GetMarketSummaries is used to get the last 24 hour summary of all active exchanges
func (b *Bittrex) GetMarketSummaries() ([]*MarketSummary, error) {
	r, err := b.client.do("GET", "public/getmarketsummaries", "", false)
	if err != nil {
		return nil, err
	}
	var response jsonResponse
	if err = json.Unmarshal(r, &response); err != nil {
		return nil, err
	}
	if err = handleErr(response); err != nil {
		return nil, err
	}
	marketSummaries := []*MarketSummary{}
	if err = json.Unmarshal(response.Result, &marketSummaries); err != nil {
		return nil, err
	}
	return marketSummaries, nil
}

// GetMarketSummary is used to get the last 24 hour summary for a given market
func (b *Bittrex) GetMarketSummary(market string) ([]*MarketSummary, error) {
	r, err := b.client.do("GET", fmt.Sprintf("public/getmarketsummary?market=%s", strings.ToUpper(market)), "", false)
	if err != nil {
		return nil, err
	}
	var response jsonResponse
	if err = json.Unmarshal(r, &response); err != nil {
		return nil, err
	}
	if err = handleErr(response); err != nil {
		return nil, err
	}

	marketSummary := []*MarketSummary{}
	if err = json.Unmarshal(response.Result, &marketSummary); err != nil {
		return nil, err
	}

	return marketSummary, nil
}

// GetOrderBook is used to get retrieve the orderbook for a given market
// market: a string literal for the market (ex: BTC-LTC)
// cat: buy, sell or both to identify the type of orderbook to return.
// depth: how deep of an order book to retrieve. Max is 100
func (b *Bittrex) GetOrderBook(market, cat string, depth int) (*OrderBook, error) {
	if cat != "buy" && cat != "sell" && cat != "both" {
		cat = "both"
	}
	if depth > 100 {
		depth = 100
	}
	if depth < 1 {
		depth = 1
	}
	r, err := b.client.do("GET", fmt.Sprintf("public/getorderbook?market=%s&type=%s&depth=%d", strings.ToUpper(market), cat, depth), "", false)
	if err != nil {
		return nil, err
	}
	var response jsonResponse
	if err = json.Unmarshal(r, &response); err != nil {
		return nil, err
	}
	if err = handleErr(response); err != nil {
		return nil, err
	}

	orderBook := OrderBook{}
	if cat == "buy" {
		err = json.Unmarshal(response.Result, &orderBook.Buy)
	} else if cat == "sell" {
		err = json.Unmarshal(response.Result, &orderBook.Sell)
	} else {
		err = json.Unmarshal(response.Result, &orderBook)
	}
	if err != nil {
		return nil, err
	}

	return &orderBook, nil
}

// GetOrderBookBuySell is used to get retrieve the buy or sell side of an orderbook for a given market
// market: a string literal for the market (ex: BTC-LTC)
// cat: buy or sell to identify the type of orderbook to return.
// depth: how deep of an order book to retrieve. Max is 100
func (b *Bittrex) GetOrderBookBuySell(market, cat string, depth int) ([]*Orderb, error) {
	if cat != "buy" && cat != "sell" {
		cat = "buy"
	}
	if depth > 100 {
		depth = 100
	}
	if depth < 1 {
		depth = 1
	}
	r, err := b.client.do("GET", fmt.Sprintf("public/getorderbook?market=%s&type=%s&depth=%d", strings.ToUpper(market), cat, depth), "", false)
	if err != nil {
		return nil, err
	}
	var response jsonResponse
	if err = json.Unmarshal(r, &response); err != nil {
		return nil, err
	}
	if err = handleErr(response); err != nil {
		return nil, err
	}

	orderb := []*Orderb{}
	if err = json.Unmarshal(response.Result, &orderb); err != nil {
		return nil, err
	}

	return orderb, nil
}

// GetMarketHistory is used to retrieve the latest trades that have occured for a specific market.
// market a string literal for the market (ex: BTC-LTC)
func (b *Bittrex) GetMarketHistory(market string) ([]*Trade, error) {
	r, err := b.client.do("GET", fmt.Sprintf("public/getmarkethistory?market=%s", strings.ToUpper(market)), "", false)
	if err != nil {
		return nil, err
	}
	var response jsonResponse
	if err = json.Unmarshal(r, &response); err != nil {
		return nil, err
	}
	if err = handleErr(response); err != nil {
		return nil, err
	}
	trades := []*Trade{}
	if err = json.Unmarshal(response.Result, &trades); err != nil {
		return nil, err
	}
	return trades, nil
}

// Market

// BuyLimit is used to place a limited buy order in a specific market.
func (b *Bittrex) BuyLimit(market string, quantity, rate float64) (uuid string, err error) {
	r, err := b.client.do("GET", "market/buylimit?market="+market+"&quantity="+strconv.FormatFloat(quantity, 'f', 8, 64)+"&rate="+strconv.FormatFloat(rate, 'f', 8, 64), "", true)
	if err != nil {
		return
	}
	var response jsonResponse
	if err = json.Unmarshal(r, &response); err != nil {
		return
	}
	if err = handleErr(response); err != nil {
		return
	}
	var u Uuid
	err = json.Unmarshal(response.Result, &u)
	uuid = u.Id
	return
}

// BuyMarket is used to place a market buy order in a spacific market.
func (b *Bittrex) BuyMarket(market string, quantity float64) (uuid string, err error) {
	r, err := b.client.do("GET", "market/buymarket?market="+market+"&quantity="+strconv.FormatFloat(quantity, 'f', 8, 64), "", true)
	if err != nil {
		return
	}
	var response jsonResponse
	if err = json.Unmarshal(r, &response); err != nil {
		return
	}
	if err = handleErr(response); err != nil {
		return
	}
	var u Uuid
	err = json.Unmarshal(response.Result, &u)
	uuid = u.Id
	return
}

// SellLimit is used to place a limited sell order in a specific market.
func (b *Bittrex) SellLimit(market string, quantity, rate float64) (uuid string, err error) {
	r, err := b.client.do("GET", "market/selllimit?market="+market+"&quantity="+strconv.FormatFloat(quantity, 'f', 8, 64)+"&rate="+strconv.FormatFloat(rate, 'f', 8, 64), "", true)
	if err != nil {
		return
	}
	var response jsonResponse
	if err = json.Unmarshal(r, &response); err != nil {
		return
	}
	if err = handleErr(response); err != nil {
		return
	}
	var u Uuid
	err = json.Unmarshal(response.Result, &u)
	uuid = u.Id
	return
}

// SellMarket is used to place a market sell order in a specific market.
func (b *Bittrex) SellMarket(market string, quantity float64) (uuid string, err error) {
	r, err := b.client.do("GET", "market/sellmarket?market="+market+"&quantity="+strconv.FormatFloat(quantity, 'f', 8, 64), "", true)
	if err != nil {
		return
	}
	var response jsonResponse
	if err = json.Unmarshal(r, &response); err != nil {
		return
	}
	if err = handleErr(response); err != nil {
		return
	}
	var u Uuid
	err = json.Unmarshal(response.Result, &u)
	uuid = u.Id
	return
}

// CancelOrder is used to cancel a buy or sell order.
func (b *Bittrex) CancelOrder(orderID string) (err error) {
	r, err := b.client.do("GET", "market/cancel?uuid="+orderID, "", true)
	if err != nil {
		return
	}
	var response jsonResponse
	if err = json.Unmarshal(r, &response); err != nil {
		return
	}
	err = handleErr(response)
	return
}

// GetOpenOrders returns orders that you currently have opened.
// If market is set to "all", GetOpenOrders return all orders
// If market is set to a specific order, GetOpenOrders return orders for this market
func (b *Bittrex) GetOpenOrders(market string) ([]*OrderHistory, error) {
	ressource := "market/getopenorders"
	if market != "all" {
		ressource += "?market=" + strings.ToUpper(market)
	}
	r, err := b.client.do("GET", ressource, "", true)
	if err != nil {
		return nil, err
	}
	var response jsonResponse
	if err = json.Unmarshal(r, &response); err != nil {
		return nil, err
	}
	if err = json.Unmarshal(r, &response); err != nil {
		return nil, err
	}
	openOrders := []*OrderHistory{}
	if err = json.Unmarshal(response.Result, &openOrders); err != nil {
		return nil, err
	}
	return openOrders, nil
}

// Account

// GetBalances is used to retrieve all balances from your account
func (b *Bittrex) GetBalances() ([]*Balance, error) {
	r, err := b.client.do("GET", "account/getbalances", "", true)
	if err != nil {
		return nil, err
	}
	var response jsonResponse
	if err = json.Unmarshal(r, &response); err != nil {
		return nil, err
	}
	if err = handleErr(response); err != nil {
		return nil, err
	}
	balances := []*Balance{}
	if err = json.Unmarshal(response.Result, &balances); err != nil {
		return nil, err
	}
	return balances, nil
}

// Getbalance is used to retrieve the balance from your account for a specific currency.
// currency: a string literal for the currency (ex: LTC)
func (b *Bittrex) GetBalance(currency string) (*Balance, error) {
	r, err := b.client.do("GET", "account/getbalance?currency="+strings.ToUpper(currency), "", true)
	if err != nil {
		return nil, err
	}
	var response jsonResponse
	if err = json.Unmarshal(r, &response); err != nil {
		return nil, err
	}
	if err = handleErr(response); err != nil {
		return nil, err
	}
	balance := Balance{}
	if err = json.Unmarshal(response.Result, &balance); err != nil {
		return nil, err
	}
	return &balance, nil
}

// GetDepositAddress is sed to generate or retrieve an address for a specific currency.
// currency a string literal for the currency (ie. BTC)
func (b *Bittrex) GetDepositAddress(currency string) (*Address, error) {
	r, err := b.client.do("GET", "account/getdepositaddress?currency="+strings.ToUpper(currency), "", true)
	if err != nil {
		return nil, err
	}
	var response jsonResponse
	if err = json.Unmarshal(r, &response); err != nil {
		return nil, err
	}
	if err = handleErr(response); err != nil {
		return nil, err
	}
	address := Address{}
	if err = json.Unmarshal(response.Result, &address); err != nil {
		return nil, err
	}
	return &address, nil
}

// Withdraw is used to withdraw funds from your account.
// address string the address where to send the funds.
// currency string literal for the currency (ie. BTC)
// quantity float the quantity of coins to withdraw
func (b *Bittrex) Withdraw(address, currency string, quantity float64) (withdrawUuid string, err error) {
	r, err := b.client.do("GET", "account/withdraw?currency="+strings.ToUpper(currency)+"&quantity="+strconv.FormatFloat(quantity, 'f', 8, 64)+"&address="+address, "", true)
	if err != nil {
		return
	}
	var response jsonResponse
	if err = json.Unmarshal(r, &response); err != nil {
		return
	}
	if err = handleErr(response); err != nil {
		return
	}
	var u Uuid
	err = json.Unmarshal(response.Result, &u)
	withdrawUuid = u.Id
	return
}

// GetOrderHistory used to retrieve your order history.
// market string literal for the market (ie. BTC-LTC). If set to "all", will return for all market
func (b *Bittrex) GetOrderHistory(market string) ([]*OrderHistory, error) {
	ressource := "account/getorderhistory"
	if market != "all" {
		ressource += "?market=" + market
	}
	r, err := b.client.do("GET", ressource, "", true)
	if err != nil {
		return nil, err
	}
	var response jsonResponse
	if err = json.Unmarshal(r, &response); err != nil {
		return nil, err
	}
	if err = handleErr(response); err != nil {
		return nil, err
	}
	orders := []*OrderHistory{}
	if err = json.Unmarshal(response.Result, &orders); err != nil {
		return nil, err
	}
	return orders, nil
}

// GetWithdrawalHistory is used to retrieve your withdrawal history
// currency string a string literal for the currency (ie. BTC). If set to "all", will return for all currencies
func (b *Bittrex) GetWithdrawalHistory(currency string) ([]*Withdrawal, error) {
	ressource := "account/getwithdrawalhistory"
	if currency != "all" {
		ressource += "currency=" + currency
	}
	r, err := b.client.do("GET", ressource, "", true)
	if err != nil {
		return nil, err
	}
	var response jsonResponse
	if err = json.Unmarshal(r, &response); err != nil {
		return nil, err
	}
	if err = handleErr(response); err != nil {
		return nil, err
	}
	withdrawals := []*Withdrawal{}
	if err = json.Unmarshal(response.Result, &withdrawals); err != nil {
		return nil, err
	}
	return withdrawals, nil
}

// GetDepositHistory is used to retrieve your deposit history
// currency string a string literal for the currency (ie. BTC). If set to "all", will return for all currencies
func (b *Bittrex) GetDepositHistory(currency string) ([]*Deposit, error) {
	ressource := "account/getdeposithistory"
	if currency != "all" {
		ressource += "currency=" + currency
	}
	r, err := b.client.do("GET", ressource, "", true)
	if err != nil {
		return nil, err
	}
	var response jsonResponse
	if err = json.Unmarshal(r, &response); err != nil {
		return nil, err
	}
	if err = handleErr(response); err != nil {
		return nil, err
	}
	deposits := []*Deposit{}
	if err = json.Unmarshal(response.Result, &deposits); err != nil {
		return nil, err
	}
	return deposits, nil
}

func (b *Bittrex) GetOrder(order_uuid string) (*Order, error) {

	ressource := "account/getorder?uuid=" + order_uuid

	r, err := b.client.do("GET", ressource, "", true)
	if err != nil {
		return nil, err
	}
	var response jsonResponse
	if err = json.Unmarshal(r, &response); err != nil {
		return nil, err
	}
	if err = json.Unmarshal(r, &response); err != nil {
		return nil, err
	}

	order := Order{}
	if err = json.Unmarshal(response.Result, &order); err != nil {
		return nil, err
	}
	return &order, nil
}

// GetTicks is used to get ticks history values for a market.
func (b *Bittrex) GetTicks(market string, interval Interval) ([]*Candle, error) {
	_, ok := CANDLE_INTERVALS[interval]
	if !ok {
		return nil, errors.New("wrong interval")
	}

	endpoint := fmt.Sprintf(
		"https://bittrex.com/Api/v2.0/pub/market/GetTicks?tickInterval=%s&marketName=%s&_=%d",
		interval, strings.ToUpper(market), rand.Int(),
	)
	r, err := b.client.do("GET", endpoint, "", false)
	if err != nil {
		return nil, fmt.Errorf("could not get market ticks: %v", err)
	}

	var response jsonResponse
	if err := json.Unmarshal(r, &response); err != nil {
		return nil, err
	}

	if err := handleErr(response); err != nil {
		return nil, err
	}

	candles := []*Candle{}
	if err := json.Unmarshal(response.Result, &candles); err != nil {
		return nil, fmt.Errorf("could not unmarshal candles: %v", err)
	}

	return candles, nil
}
