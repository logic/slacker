// Copyright 2015 Ed Marshall. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the COPYING file.

package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"
)

// Ticker is the handler for the "/ticker" Slack slash command.
func Ticker(w http.ResponseWriter, req *http.Request) error {
	symbol := strings.ToUpper(req.FormValue("text"))
	if len(symbol) == 0 {
		return StatusError{http.StatusBadRequest,
			errors.New("Usage: /ticker [symbol]")}
	}

	if regexp.MustCompile(`[^A-Z0-9.]+`).MatchString(symbol) {
		return StatusError{http.StatusBadRequest,
			errors.New("Invalid ticker symbol (letters, numbers, and '.' only)")}
	}

	// We can either do responses in-line, if we think we can get it done
	// in time before the Slack timeout. However, if we think the response
	// will take too long, we can send the response asynchronously; the
	// downside, though, is that we have to display the original "/ticker
	// foo" command regardless of whether the lookup was successful,
	// because we have to decide whether to show it here.
	var payload map[string]interface{}
	if Config.AsyncResponse {
		responseUrl := req.FormValue("response_url")
		if len(responseUrl) == 0 {
			return StatusError{http.StatusBadRequest,
				errors.New("No response URL supplied (Slack bug?)")}
		}
		go TickerPoster(symbol, responseUrl)
		payload = map[string]interface{}{
			"response_type": "in_channel",
		}
	} else {
		payload = BuildTickerPayload(symbol)
		if _, ok := payload["attachments"]; !ok {
			// In the async case, we'd want to deliver this as payload to
			// the caller, but for immediate-response, we might as well
			// log this like a normal error and return a more appropriate
			// HTTP status code.
			return StatusError{http.StatusInternalServerError,
				errors.New(payload["text"].(string))}
		}
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Couldn't marshal payload: %v", payload)
		return StatusError{http.StatusInternalServerError,
			errors.New("Could not marshal response")}
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonPayload)
	return nil
}

// BuildTickerPayload formats the requested ticker symbol information into
// a JSON payload for rendering to the user in Slack.
func BuildTickerPayload(symbol string) map[string]interface{} {
	payload := map[string]interface{}{}
	quote, err := GetTicker(symbol)
	if err != nil {
		payload["text"] = fmt.Sprintf("Ticker symbol lookup failed for `%s`: %s",
			symbol, err.Error())
	} else if err != nil || quote == nil {
		payload["text"] = fmt.Sprintf("Unknown ticker symbol `%s`", symbol)
	} else {
		emoji := ":chart_with_upwards_trend:"
		color := "good"
		if quote.PercentChange[0] == '-' {
			emoji = ":chart_with_downwards_trend:"
			color = "danger"
		}
		payload["attachments"] = []map[string]interface{}{{
			"fallback": fmt.Sprintf(
				"%s (%s): %s (%s from previous close of %s) as of %s %s",
				quote.Symbol, quote.Name, quote.LastTradePriceOnly,
				quote.PercentChange, quote.PreviousClose, quote.LastTradeTime,
				quote.LastTradeDate),
			"pretext": fmt.Sprintf("%s *<https://finance.yahoo.com/q?s=%s|%s - %s>*",
				emoji, quote.Symbol, quote.Symbol, quote.Name),
			"text": fmt.Sprintf("*%s* _(%s from previous close of %s)_\n%s %s",
				quote.LastTradePriceOnly, quote.PercentChange, quote.PreviousClose,
				quote.LastTradeTime, quote.LastTradeDate),
			"color": color,
			// The "fresh" parameter is non-standard, but is used
			// to defeat any caching here.
			"image_url": fmt.Sprintf(
				"https://chart.finance.yahoo.com/t?s=%s&width=400&height=185&fresh=%d",
				quote.Symbol, time.Now().Unix()),
			"mrkdwn_in": []string{"text", "pretext"},
		}}
		payload["response_type"] = "in_channel"
	}
	return payload
}

// TickerPoster, as a goroutine, collects and formats the requested ticker
// symbol information, and posts it back to Slack asynchronously.
func TickerPoster(symbol string, response_url string) {
	payload := BuildTickerPayload(symbol)
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Couldn't marshal payload: %v", payload)
		return
	}
	resp, err := http.Post(response_url, "application/json", bytes.NewReader(jsonPayload))
	if err != nil {
		log.Printf("Failed to post response to '%s': %s\n", response_url, err.Error())
		return
	}
	if resp.StatusCode != http.StatusOK {
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Printf("Couldn't read from '%s': %s\n", response_url, err.Error())
			return
		}
		log.Printf("Got %d from %s: %s\n", resp.StatusCode, response_url, string(body))
	}
}
