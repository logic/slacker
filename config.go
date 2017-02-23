// Copyright 2015, 2016, 2017 Ed Marshall. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the COPYING file.

package main

import (
	"io"
	"log"
	"strings"
	"time"
)

import "github.com/BurntSushi/toml"

// Configuration represents the fields of a TOML configuration file.
type Configuration struct {
	Tokens            []string
	ListenAddress     string
	AsyncResponse     bool
	HTTPClientTimeout time.Duration
}

// LoadConfig sets our configuration defaults, and loads a configuration from
// the TOML-formatted configuration file over the defaults.
func LoadConfig(configStream io.Reader) (Configuration, error) {
	config := Configuration{nil, "0.0.0.0:8888", false, 0}
	_, err := toml.DecodeReader(configStream, &config)

	// Normalize timeout to seconds, because toml lacks duration support
	config.HTTPClientTimeout = config.HTTPClientTimeout * time.Second
	log.Println("Configuration loaded:")
	if len(config.Tokens) > 0 {
		log.Printf("  Tokens: [<hidden>%s]\n",
			strings.Repeat(", <hidden>", len(config.Tokens)-1))
	} else {
		log.Println("  No tokens defined (accepting all requests)")
	}
	log.Printf("  Listening on %s\n", config.ListenAddress)
	if config.AsyncResponse {
		log.Println("  Sending responses asynchronously")
	} else {
		log.Println("  Sending responses immediately")
	}
	log.Printf("  HTTP connections time out in %v\n", config.HTTPClientTimeout)

	return config, err
}
