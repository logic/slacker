// Copyright 2015, 2016, 2017 Ed Marshall. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the COPYING file.

package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

var apiYahooFinance = "https://query1.finance.yahoo.com/v7/finance/quote"

// APIError represents an error response
type APIError struct {
	Code        string `json:"code"`        // "argument-error"
	Description string `json:"description"` // "Missing value for the \"symbols\" argument"
}

// APIResult represents a successful query
type APIResult struct {
	Ask                               float64 `json:"ask,omitempty"`                               // 0.0,
	AskSize                           int     `json:"askSize,omitempty"`                           // 0,
	AverageDailyVolume10Day           int64   `json:"averageDailyVolume10Day,omitempty"`           // 14627912,
	AverageDailyVolume3Month          int64   `json:"averageDailyVolume3Month,omitempty"`          // 16268415,
	Bid                               float64 `json:"bid,omitempty"`                               // 0.0,
	BidSize                           int     `json:"bidSize,omitempty"`                           // 0,
	BookValue                         float64 `json:"bookValue,omitempty"`                         // 6.452,
	Currency                          string  `json:"currency,omitempty"`                          // "USD",
	EarningsTimestamp                 int64   `json:"earningsTimestamp,omitempty"`                 // 1509021000,
	EarningsTimestampEnd              int64   `json:"earningsTimestampEnd,omitempty"`              // 1518442200,
	EarningsTimestampStart            int64   `json:"earningsTimestampStart,omitempty"`            // 1518010200,
	EpsForward                        float64 `json:"epsForward,omitempty"`                        // 0.45,
	EpsTrailingTwelveMonths           float64 `json:"epsTrailingTwelveMonths,omitempty"`           // -0.627,
	Exchange                          string  `json:"exchange,omitempty"`                          // "NYQ",
	ExchangeDataDelayedBy             int     `json:"exchangeDataDelayedBy,omitempty"`             // 0,
	ExchangeTimezoneName              string  `json:"exchangeTimezoneName,omitempty"`              // "America/New_York",
	ExchangeTimezoneShortName         string  `json:"exchangeTimezoneShortName,omitempty"`         // "EST",
	FiftyDayAverage                   float64 `json:"fiftyDayAverage,omitempty"`                   // 19.315556,
	FiftyDayAverageChange             float64 `json:"fiftyDayAverageChange,omitempty"`             // 2.954445,
	FiftyDayAverageChangePercent      float64 `json:"fiftyDayAverageChangePercent,omitempty"`      // 0.15295677,
	FiftyTwoWeekHigh                  float64 `json:"fiftyTwoWeekHigh,omitempty"`                  // 22.4,
	FiftyTwoWeekHighChange            float64 `json:"fiftyTwoWeekHighChange,omitempty"`            // -0.12999916,
	FiftyTwoWeekHighChangePercent     float64 `json:"fiftyTwoWeekHighChangePercent,omitempty"`     // -0.005803534,
	FiftyTwoWeekLow                   float64 `json:"fiftyTwoWeekLow,omitempty"`                   // 14.12,
	FiftyTwoWeekLowChange             float64 `json:"fiftyTwoWeekLowChange,omitempty"`             // 8.150001,
	FiftyTwoWeekLowChangePercent      float64 `json:"fiftyTwoWeekLowChangePercent,omitempty"`      // 0.5771955,
	FinancialCurrency                 string  `json:"financialCurrency,omitempty"`                 // "USD",
	ForwardPE                         float64 `json:"forwardPE,omitempty"`                         // 49.48889,
	FullExchangeName                  string  `json:"fullExchangeName,omitempty"`                  // "NYSE",
	GmtOffSetMilliseconds             int64   `json:"gmtOffSetMilliseconds,omitempty"`             // -18000000,
	Language                          string  `json:"language,omitempty"`                          // "en-US",
	LongName                          string  `json:"longName,omitempty"`                          // "Twitter, Inc.",
	Market                            string  `json:"market,omitempty"`                            // "us_market",
	MarketCap                         int64   `json:"marketCap,omitempty"`                         // 16471872512,
	MarketState                       string  `json:"marketState,omitempty"`                       // "CLOSED",
	MessageBoardId                    string  `json:"messageBoardId,omitempty"`                    // "finmb_35962803",
	PostMarketChange                  float64 `json:"postMarketChange,omitempty"`                  // -0.010000229,
	PostMarketChangePercent           float64 `json:"postMarketChangePercent,omitempty"`           // -0.044904485,
	PostMarketPrice                   float64 `json:"postMarketPrice,omitempty"`                   // 22.26,
	PostMarketTime                    int64   `json:"postMarketTime,omitempty"`                    // 1511398614,
	PriceHint                         int     `json:"priceHint,omitempty"`                         // 2,
	PriceToBook                       float64 `json:"priceToBook,omitempty"`                       // 3.451643,
	QuoteSourceName                   string  `json:"quoteSourceName,omitempty"`                   // "Delayed Quote",
	QuoteType                         string  `json:"quoteType,omitempty"`                         // "EQUITY",
	RegularMarketChange               float64 `json:"regularMarketChange,omitempty"`               // 0.3900013,
	RegularMarketChangePercent        float64 `json:"regularMarketChangePercent,omitempty"`        // 1.7824557,
	RegularMarketDayHigh              float64 `json:"regularMarketDayHigh,omitempty"`              // 22.4,
	RegularMarketDayLow               float64 `json:"regularMarketDayLow,omitempty"`               // 21.8,
	RegularMarketOpen                 float64 `json:"regularMarketOpen,omitempty"`                 // 21.9,
	RegularMarketPreviousClose        float64 `json:"regularMarketPreviousClose,omitempty"`        // 21.88,
	RegularMarketPrice                float64 `json:"regularMarketPrice,omitempty"`                // 22.27,
	RegularMarketTime                 int64   `json:"regularMarketTime,omitempty"`                 // 1511384466,
	RegularMarketVolume               int64   `json:"regularMarketVolume,omitempty"`               // 21161825,
	SharesOutstanding                 int64   `json:"sharesOutstanding,omitempty"`                 // 739644032,
	ShortName                         string  `json:"shortName,omitempty"`                         // "Twitter, Inc.",
	SourceInterval                    int     `json:"sourceInterval,omitempty"`                    // 15,
	Symbol                            string  `json:"symbol,omitempty"`                            // "TWTR",
	Tradeable                         bool    `json:"tradeable,omitempty"`                         // true,
	TwoHundredDayAverage              float64 `json:"twoHundredDayAverage,omitempty"`              // 18.075928,
	TwoHundredDayAverageChange        float64 `json:"twoHundredDayAverageChange,omitempty"`        // 4.1940727,
	TwoHundredDayAverageChangePercent float64 `json:"twoHundredDayAverageChangePercent,omitempty"` // 0.23202531
}

// APIMessage represents a standard API response
type APIMessage struct {
	Error   *APIError   `json:"error,omitempty"`
	Results []APIResult `json:"result,omitempty"`
}

// APIEnvelope is the wrapping envelope around a Yahoo Finance API response
type APIEnvelope map[string]APIMessage

// GetTickers asks Yahoo Finance for a complete rundown of information about
// a given stock symbol, and returns it as a YahooQuote, or returns an error
// if something goes wrong.
func GetTickers(symbols []string) ([]APIResult, error) {
	query, err := url.Parse(apiYahooFinance)
	if err != nil {
		return nil, err
	}

	params := url.Values{
		"symbols": {strings.Join(symbols, ",")},
	}
	query.RawQuery = params.Encode()
	client := http.Client{Timeout: Config.HTTPClientTimeout}
	resp, err := client.Get(query.String())
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Yahoo Finance API returned %d status",
			resp.StatusCode)
	}
	if !strings.HasPrefix(resp.Header["Content-Type"][0], "application/json") {
		return nil, fmt.Errorf("Yahoo Finance API returned `%s` content type",
			resp.Header["Content-Type"][0])
	}

	var envelope APIEnvelope
	var results []APIResult
	dec := json.NewDecoder(resp.Body)
	for {
		if err := dec.Decode(&envelope); err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}
		for _, msg := range envelope {
			if msg.Error != nil {
				return nil, fmt.Errorf("Yahoo Finance API returned `%s` error: %s",
					msg.Error.Code, msg.Error.Description)
			}
			for _, result := range msg.Results {
				results = append(results, result)
			}
		}
	}
	return results, nil
}
