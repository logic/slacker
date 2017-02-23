// Copyright 2015, 2016, 2017 Ed Marshall. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the COPYING file.

package main

import (
	"fmt"
	"io"
	"log"
	"strings"
	"time"
)

import "github.com/BurntSushi/toml"

// Tokens represent an array of Slack connection tokens
type Tokens []string

func (t *Tokens) String() string {
	return fmt.Sprintf("%+v", *t)
}

// Set accepts an additional token and adds it to the list
func (t *Tokens) Set(v string) error {
	*t = append(*t, v)
	return nil
}

// Configuration represents the fields of a TOML configuration file.
type Configuration struct {
	File              string
	Tokens            Tokens
	ListenAddress     string
	AsyncResponse     bool
	HTTPClientTimeout time.Duration
}

// LoadConfig sets our configuration defaults, and loads a configuration from
// the TOML-formatted configuration file over the defaults.
func LoadConfig(config *Configuration, configStream io.Reader) error {
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

	return err
}
