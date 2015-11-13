// Copyright 2015 Ed Marshall. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the COPYING file.

package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

// YahooTiming is used by several elements of Query.Diagnostics, so we
// pull it out for simlicity's sake.
type YahooTiming struct {
	ExecutionStartTime int `json:"execution-start-time,string"`
	ExecutionStopTime  int `json:"execution-stop-time,string"`
	ExecutionTime      int `json:"execution-time,string"`
}

// YahooResponse is the full JSON document received back from the Yahoo
// Finance API.
//
// Note that this is not entirely correct; if Query.Count > 1, Query.Results
// will be an array of Quotes (rather than a single Quote). Similarly, if
// only a single Query.Diagnostics.Url is available, the array wrapper will
// also be eliminated. For our purposes (requesting single quotes), this is
// sufficient.
type YahooResponse struct {
	Query struct {
		Created     time.Time
		Count       int
		Diagnostics struct {
			BuildVersion string `json:"build-version"`
			Cache        struct {
				*YahooTiming
				Content string
				Method  string
				Type    string
			}
			Javascript struct {
				*YahooTiming
				InstructionsUsed int    `json:"instructions-used,string"`
				TableName        string `json:"table-name"`
			}
			PubliclyCallable bool `json:"publiclyCallable,string"`
			Query            map[string]*string
			ServiceTime      string `json:"service-time"`
			Url              []struct {
				*YahooTiming
				Content string
			}
			UserTime string `json:"user-time"`
		}
		Lang    string
		Results struct {
			Quote YahooQuote
		}
	}
}

// YahooQuote is a direct translation of the JSON object returned for a
// single symbol quote result.
type YahooQuote struct {
	AfterHoursChangeRealtime                       string
	AnnualizedGain                                 string
	AskRealtime                                    string
	Ask                                            string
	AverageDailyVolume                             string
	BidRealtime                                    string
	Bid                                            string
	BookValue                                      string
	ChangeFromFiftydayMovingAverage                string
	ChangeFromTwoHundreddayMovingAverage           string
	ChangeFromYearHigh                             string
	ChangeFromYearLow                              string
	ChangeinPercent                                string
	Change_PercentChange                           string
	ChangePercentRealtime                          string
	ChangeRealtime                                 string
	Change                                         string
	Commission                                     string
	Currency                                       string
	DaysHigh                                       string
	DaysLow                                        string
	DaysRangeRealtime                              string
	DaysRange                                      string
	DaysValueChangeRealtime                        string
	DaysValueChange                                string
	DividendPayDate                                string
	DividendShare                                  string
	DividendYield                                  string
	EarningsShare                                  string
	EBITDA                                         string
	EPSEstimateCurrentYear                         string
	EPSEstimateNextQuarter                         string
	EPSEstimateNextYear                            string
	ErrorIndicationreturnedforsymbolchangedinvalid string
	ExDividendDate                                 string
	FiftydayMovingAverage                          string
	HighLimit                                      string
	HoldingsGainPercentRealtime                    string
	HoldingsGainPercent                            string
	HoldingsGainRealtime                           string
	HoldingsGain                                   string
	HoldingsValueRealtime                          string
	HoldingsValue                                  string
	LastTradeDate                                  string
	LastTradePriceOnly                             string
	LastTradeRealtimeWithTime                      string
	LastTradeTime                                  string
	LastTradeWithTime                              string
	LowLimit                                       string
	MarketCapitalization                           string
	MarketCapRealtime                              string
	MoreInfo                                       string
	Name                                           string
	Notes                                          string
	OneyrTargetPrice                               string
	Open                                           string
	OrderBookRealtime                              string
	PEGRatio                                       string
	PERatioRealtime                                string
	PERatio                                        string
	PercebtChangeFromYearHigh                      string
	PercentChangeFromFiftydayMovingAverage         string
	PercentChangeFromTwoHundreddayMovingAverage    string
	PercentChangeFromYearLow                       string
	PercentChange                                  string
	PreviousClose                                  string
	PriceBook                                      string
	PriceEPSEstimateCurrentYear                    string
	PriceEPSEstimateNextYear                       string
	PricePaid                                      string
	PriceSales                                     string
	SharesOwned                                    string
	ShortRatio                                     string
	StockExchange                                  string
	Symbol                                         string `json:"Symbol"`
	TickerTrend                                    string
	TradeDate                                      string
	TwoHundreddayMovingAverage                     string
	Volume                                         string
	YearHigh                                       string
	YearLow                                        string
	YearRange                                      string
	_Symbol                                        string `json:"symbol"`
}

// GetTicker asks Yahoo Finance for a complete rundown of information about
// a given stock symbol, and returns it as a YahooQuote, or returns an error
// if something goes wrong.
func GetTicker(symbol string) (*YahooQuote, error) {
	params := url.Values{
		"q": {
			fmt.Sprintf("select * from yahoo.finance.quotes where symbol=\"%s\"",
				symbol),
		},
		"format":      {"json"},
		"diagnostics": {"true"},
		"env":         {"http://datatables.org/alltables.env"},
		"callback":    {""},
	}

	query, err := url.Parse("https://query.yahooapis.com/v1/public/yql")
	if err != nil {
		return nil, err
	}
	query.RawQuery = params.Encode()
	client := http.Client{Timeout: Config.HttpClientTimeout}
	resp, err := client.Get(query.String())
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	dec := json.NewDecoder(resp.Body)
	var yr YahooResponse
	if err := dec.Decode(&yr); err != nil {
		return nil, err
	}

	if len(yr.Query.Results.Quote.Name) == 0 {
		return nil, nil
	} else {
		return &yr.Query.Results.Quote, nil
	}
}
