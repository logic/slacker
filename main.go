// Copyright 2015, 2016, 2017 Ed Marshall. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the COPYING file.

package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
)

// Config is our global (read-only) configuration state.
var Config Configuration

var version = "development version"
var timestamp = "unknown"

func versionString() string {
	return fmt.Sprintf("slacker %s (build date %s)", version, timestamp)
}

func parseCli() Configuration {
	var c Configuration
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "%s: [flags]\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.StringVar(&c.File, "config-file", "slack.toml",
		"Configuration file to load")
	flag.Var(&c.Tokens, "token",
		"Token to accept (can be specified multiple times)")
	flag.StringVar(&c.ListenAddress, "listen-address", "0.0.0.0:8000",
		"Address and port to listen on")
	flag.BoolVar(&c.AsyncResponse, "async-response", true,
		"Whether to respond to requests asynchronously")
	flag.DurationVar(&c.HTTPClientTimeout, "http-client-timeout", 0,
		"Time to wait before cancelling an external request")
	ver := flag.Bool("version", false, "Display current version")
	flag.Parse()
	if *ver {
		fmt.Println(versionString())
		os.Exit(0)
	}
	if f, err := os.Open(c.File); err != nil {
		log.Printf("Warning: could not load configuration file %s\n",
			c.File)
	} else {
		if err = LoadConfig(&c, f); err != nil {
			log.Fatal("Could not decode configuration: ", err)
		}
		f.Close()
	}
	return c
}

func main() {
	Config = parseCli()
	log.Printf("%s started\n", versionString())

	http.Handle("/cmd", RequestIDMiddleware(ErrorHandler(SlackDispatcher)))
	if err := http.ListenAndServe(Config.ListenAddress, nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
