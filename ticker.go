// Copyright 2015, 2016, 2017 Ed Marshall. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the COPYING file.

package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// TickerOpts represnts a set of parsed /ticker command options.
type TickerOpts struct {
	Symbol   string
	Period   string
	Interval int
	Type     string
	Log      bool
}

// ParseTickerCommand takes the /ticker command line and parses it into
// TickerOpts, returning an error if anything goes wrong.
func ParseTickerCommand(cmd string) (TickerOpts, error) {
	var opts TickerOpts
	var output bytes.Buffer
	flags := flag.NewFlagSet("/ticker", flag.ContinueOnError)
	flags.Usage = func() {
		fmt.Fprintln(&output, "usage: /ticker [flags] symbol")
		flags.PrintDefaults()
	}
	flags.SetOutput(&output)
	flags.StringVar(&opts.Period, "period", "1d", "period [xd|xY]")
	flags.IntVar(&opts.Interval, "interval", 60, "interval [seconds]")

	if err := flags.Parse(strings.Split(cmd, " ")); err != nil {
		fmt.Fprintln(&output, err)
		return opts, errors.New(output.String())
	}

	if value, err := strconv.Atoi(opts.Period[0 : len(opts.Period)-1]); value == 0 || err != nil {
		fmt.Fprintln(&output, "*Error:* period must be a positive number (followed by [d|Y])")
		flags.Usage()
		return opts, errors.New(output.String())
	}
	switch opts.Period[len(opts.Period)-1] {
	case 'd', 'Y':
		break
	default:
		fmt.Fprintln(&output, "*Error:* period must be one of 'd' (days) or 'Y' (years)")
		flags.Usage()
		return opts, errors.New(output.String())
	}

	if flags.NArg() != 1 || flags.Arg(0) == "" {
		if flags.NArg() <= 1 {
			fmt.Fprintln(&output, "*Error:* no ticker symbol specified")
		} else {
			fmt.Fprintln(&output, "*Error:* only one ticker symbol at a time")
		}
		flags.Usage()
		return opts, errors.New(output.String())
	}

	opts.Symbol = strings.ToUpper(flags.Arg(0))
	if regexp.MustCompile(`[^A-Z0-9.]+`).MatchString(opts.Symbol) {
		return opts, errors.New("*Error:* Invalid ticker symbol (letters, numbers, and '.' only)")
	}

	return opts, nil
}

// Ticker is the handler for the "/ticker" Slack slash command.
func Ticker(w http.ResponseWriter, req *http.Request) error {
	var payload map[string]interface{}

	opts, err := ParseTickerCommand(req.FormValue("text"))
	if err != nil {
		payload = map[string]interface{}{
			"response_type": "ephemeral",
			"text":          err.Error(),
		}
	} else {

		// We can either do responses in-line, if we think we can get it done
		// in time before the Slack timeout. However, if we think the response
		// will take too long, we can send the response asynchronously; the
		// downside, though, is that we have to display the original "/ticker
		// foo" command regardless of whether the lookup was successful,
		// because we have to decide whether to show it here.
		if Config.AsyncResponse {
			responseURL := req.FormValue("response_url")
			if len(responseURL) == 0 {
				return StatusError{http.StatusBadRequest,
					errors.New("No response URL supplied (Slack bug?)")}
			}
			go TickerPoster(opts, responseURL, req.Context())
			payload = map[string]interface{}{
				"response_type": "in_channel",
			}
		} else {
			payload = BuildTickerPayload(opts, req.Context())
			if _, ok := payload["attachments"]; !ok {
				// In the async case, we'd want to deliver this as
				// payload to the caller, but for immediate-response,
				// we might as well log this like a normal error and
				// return a more appropriate HTTP status code.
				return StatusError{http.StatusInternalServerError,
					errors.New(payload["text"].(string))}
			}
		}
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return StatusError{http.StatusInternalServerError,
			fmt.Errorf("Could not marshal response: %+v", payload)}
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonPayload)
	return nil
}

// BuildTickerPayload formats the requested ticker symbol information into
// a JSON payload for rendering to the user in Slack.
func BuildTickerPayload(opts TickerOpts, ctx context.Context) map[string]interface{} {
	payload := map[string]interface{}{}
	quotes, err := GetTickers([]string{opts.Symbol})
	if err != nil {
		log.Printf("[%d] Error: %s\n", RequestID(ctx), err)
		payload["text"] = fmt.Sprintf("An error occurred looking up _%s_", opts.Symbol)
	} else if (quotes == nil) || len(quotes) == 0 {
		payload["text"] = fmt.Sprintf("Unknown ticker symbol _%s_", opts.Symbol)
	} else {
		quote := quotes[0]
		if err != nil {
			payload["text"] = err.Error()
		} else {
			var emoji string
			var color string
			if quote.RegularMarketChange < 0 {
				emoji = ":chart_with_downwards_trend:"
				color = "danger"
			} else if quote.RegularMarketChange > 0 {
				emoji = ":chart_with_upwards_trend:"
				color = "good"
			} else {
				emoji = ":bar_chart:"
				color = "warning"
			}

			var name string
			if len(quote.LongName) != 0 {
				name = fmt.Sprintf("%s - %s", quote.Symbol, quote.LongName)
			} else {
				name = quote.Symbol
			}

			var upDown string
			if quote.RegularMarketChange < 0 {
				upDown = fmt.Sprintf("down %0.2f%%",
					quote.RegularMarketChangePercent*(-1))
			} else if quote.RegularMarketChange > 0 {
				upDown = fmt.Sprintf("up %0.2f%%",
					quote.RegularMarketChangePercent)
			} else {
				upDown = "unchanged"
			}

			price := quote.RegularMarketPrice

			change := fmt.Sprintf("_(%s from previous close of $%0.2f)_ ",
				upDown, quote.RegularMarketPreviousClose)

			asOf := time.Unix(quote.RegularMarketTime, 0).Format(time.RFC822)

			payload["attachments"] = []map[string]interface{}{{
				"fallback": fmt.Sprintf("%s: $%0.2f %sas of %s",
					name, price, change, asOf),
				"pretext": fmt.Sprintf("%s *<https://finance.yahoo.com/q?s=%s|%s>*",
					emoji, quote.Symbol, name),
				"text":  fmt.Sprintf("*$%0.2f* %s\n%s", price, change, asOf),
				"color": color,
				// The "fresh" parameter is non-standard, but is used
				// to defeat any caching here.
				"image_url": fmt.Sprintf(
					"https://finance.google.com/finance/getchart?q=%s&p=%s&i=%d&fresh=%d",
					quote.Symbol,
					opts.Period, opts.Interval, time.Now().Unix()),
				"mrkdwn_in": []string{"text", "pretext"},
			}}
			payload["response_type"] = "in_channel"
			log.Printf("[%d] %s $%0.2f (%s)\n", RequestID(ctx),
				quote.Symbol, price, change)
		}
	}
	return payload
}

// TickerPoster (as a goroutine) collects and formats the requested ticker
// symbol information, and posts it back to Slack asynchronously.
func TickerPoster(opts TickerOpts, responseURL string, ctx context.Context) {
	payload := BuildTickerPayload(opts, ctx)
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Couldn't marshal payload: %v", payload)
		return
	}
	resp, err := http.Post(responseURL, "application/json", bytes.NewReader(jsonPayload))
	if err != nil {
		log.Printf("[%d] POST failed for '%s': %s\n", RequestID(ctx),
			responseURL, err.Error())
		return
	}
	if resp.StatusCode != http.StatusOK {
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Printf("[%d] Couldn't read from '%s': %s\n",
				RequestID(ctx), responseURL, err.Error())
			return
		}
		log.Printf("[%d] Got %d from %s: %s\n", RequestID(ctx),
			resp.StatusCode, responseURL, string(body))
	}
}
