// Copyright 2015, 2016, 2017 Ed Marshall. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the COPYING file.

package main

import (
	"log"
	"net/http"
	"os"
)

// Config is our global (read-only) configuration state.
var Config Configuration

// Commands are any Slack commands that we recognize, and their handlers
var Commands SlashCommands

var version string = "development version"
var timestamp string = "unknown"

func main() {
	log.Printf("Starting slacker %s (build date %s)\n", version, timestamp)
	f, err := os.Open("slack.toml")
	if err != nil {
		log.Fatal("Could not open slack.toml: ", err)
	}
	if Config, err = LoadConfig(f); err != nil {
		log.Fatal("Could not decode configuration: ", err)
	}
	f.Close()

	Commands = SlashCommands{
		"/ticker": Ticker,
	}

	http.Handle("/cmd", RequestIDMiddleware(ErrorHandler(SlackDispatcher)))
	if err = http.ListenAndServe(Config.ListenAddress, nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
