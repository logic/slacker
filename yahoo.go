// Copyright 2015, 2016, 2017 Ed Marshall. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the COPYING file.

package main

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"unsafe"
)

var apiYahooFinance = "https://download.finance.yahoo.com/d/quotes.csv"

// TickerOption represents a single piece of available information about a
// ticker symbol.
type TickerOption string

// TickerResults is an array of individual ticker lookup results, each
// formatted as a mapping of a TickerOption to the retreived string value.
type TickerResults []map[TickerOption]string

// The full list of TickerOption values
const (
	TOAsk                                            TickerOption = "a"
	TOAverageDailyVolume                             TickerOption = "a2"
	TOBid                                            TickerOption = "b"
	TOAskRealtime                                    TickerOption = "b2"
	TOBidRealtime                                    TickerOption = "b3"
	TOBookValue                                      TickerOption = "b4"
	TOChangeAndPercentChange                         TickerOption = "c"
	TOChange                                         TickerOption = "c1"
	TOCommission                                     TickerOption = "c3"
	TOCurrency                                       TickerOption = "c4"
	TOChangeRealtime                                 TickerOption = "c6"
	TOAfterHoursChangeRealtime                       TickerOption = "c8"
	TODividendShare                                  TickerOption = "d"
	TOLastTradeDate                                  TickerOption = "d1"
	TOTradeDate                                      TickerOption = "d2"
	TOEarningsShare                                  TickerOption = "e"
	TOErrorIndicationreturnedforsymbolchangedinvalid TickerOption = "e1"
	TOEPSEstimateCurrentYear                         TickerOption = "e7"
	TOEPSEstimateNextYear                            TickerOption = "e8"
	TOEPSEstimateNextQuarter                         TickerOption = "e9"
	TODaysLow                                        TickerOption = "g"
	TOHoldingsGainPercent                            TickerOption = "g1"
	TOAnnualizedGain                                 TickerOption = "g3"
	TOHoldingsGain                                   TickerOption = "g4"
	TOHoldingsGainPercentRealtime                    TickerOption = "g5"
	TOHoldingsGainRealtime                           TickerOption = "g6"
	TODaysHigh                                       TickerOption = "h"
	TOMoreInfo                                       TickerOption = "i"
	TOOrderBookRealtime                              TickerOption = "i5"
	TOYearLow                                        TickerOption = "j"
	TOMarketCapitalization                           TickerOption = "j1"
	TOMarketCapRealtime                              TickerOption = "j3"
	TOEBITDA                                         TickerOption = "j4"
	TOChangeFromYearLow                              TickerOption = "j5"
	TOPercentChangeFromYearLow                       TickerOption = "j6"
	TOYearHigh                                       TickerOption = "k"
	TOLastTradeRealtimeWithTime                      TickerOption = "k1"
	TOChangePercentRealtime                          TickerOption = "k2"
	TOChangeFromYearHigh                             TickerOption = "k4"
	TOPercebtChangeFromYearHigh                      TickerOption = "k5"
	TOLastTradeWithTime                              TickerOption = "l"
	TOLastTradePriceOnly                             TickerOption = "l1"
	TOHighLimit                                      TickerOption = "l2"
	TOLowLimit                                       TickerOption = "l3"
	TODaysRange                                      TickerOption = "m"
	TODaysRangeRealtime                              TickerOption = "m2"
	TOFiftydayMovingAverage                          TickerOption = "m3"
	TOTwoHundreddayMovingAverage                     TickerOption = "m4"
	TOChangeFromTwoHundreddayMovingAverage           TickerOption = "m5"
	TOPercentChangeFromTwoHundreddayMovingAverage    TickerOption = "m6"
	TOChangeFromFiftydayMovingAverage                TickerOption = "m7"
	TOPercentChangeFromFiftydayMovingAverage         TickerOption = "m8"
	TOName                                           TickerOption = "n"
	TONotes                                          TickerOption = "n4"
	TOOpen                                           TickerOption = "o"
	TOPreviousClose                                  TickerOption = "p"
	TOPricePaid                                      TickerOption = "p1"
	TOChangeinPercent                                TickerOption = "p2"
	TOPriceSales                                     TickerOption = "p5"
	TOPriceBook                                      TickerOption = "p6"
	TOExDividendDate                                 TickerOption = "q"
	TOPERatio                                        TickerOption = "r"
	TODividendPayDate                                TickerOption = "r1"
	TOPERatioRealtime                                TickerOption = "r2"
	TOPEGRatio                                       TickerOption = "r5"
	TOPriceEPSEstimateCurrentYear                    TickerOption = "r6"
	TOPriceEPSEstimateNextYear                       TickerOption = "r7"
	TOSymbol                                         TickerOption = "s"
	TOSharesOwned                                    TickerOption = "s1"
	TOShortRatio                                     TickerOption = "s7"
	TOLastTradeTime                                  TickerOption = "t1"
	TOTickerTrend                                    TickerOption = "t7"
	TOOneyrTargetPrice                               TickerOption = "t8"
	TOVolume                                         TickerOption = "v"
	TOHoldingsValue                                  TickerOption = "v1"
	TOHoldingsValueRealtime                          TickerOption = "v7"
	TOYearRange                                      TickerOption = "w"
	TODaysValueChange                                TickerOption = "w1"
	TODaysValueChangeRealtime                        TickerOption = "w4"
	TOStockExchange                                  TickerOption = "x"
	TODividendYield                                  TickerOption = "y"
)

func buildFParams(values []TickerOption) string {
	// ew ew ew ew
	foo := *(*[]string)((unsafe.Pointer(&values)))
	return strings.Join(foo, "")
}

// GetTickers asks Yahoo Finance for a complete rundown of information about
// a given stock symbol, and returns it as a YahooQuote, or returns an error
// if something goes wrong.
func GetTickers(symbols []string, values []TickerOption) (TickerResults, error) {
	query, err := url.Parse(apiYahooFinance)
	if err != nil {
		return nil, err
	}

	params := url.Values{
		"s": {strings.Join(symbols, ",")},
		"f": {buildFParams(values)},
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
	if resp.Header["Content-Type"][0] != "application/octet-stream" {
		return nil, fmt.Errorf("Yahoo Finance API returned `%s` content type",
			resp.Header["Content-Type"][0])
	}

	r := csv.NewReader(resp.Body)
	r.FieldsPerRecord = len(values)
	records, err := r.ReadAll()
	if err != nil {
		return nil, err
	}
	if len(records) != len(symbols) {
		return nil, fmt.Errorf("Yahoo Finance API returned %d results, expected %d",
			len(records), len(symbols))
	}

	results := make(TickerResults, len(records))
	for i, record := range records {
		result := make(map[TickerOption]string, len(values))
		for j, field := range record {
			result[values[j]] = field
		}
		results[i] = result
	}
	return results, nil
}
