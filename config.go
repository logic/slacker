// Copyright 2015 Ed Marshall. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the COPYING file.

package main

import "github.com/BurntSushi/toml"

type Configuration struct {
	Tokens        []string
	ListenAddress string
	AsyncResponse bool
}

// LoadConfig sets our configuration defaults, and loads a configuration from
// the TOML-formatted configuration file over the defaults.
func LoadConfig() (Configuration, error) {
	config := Configuration{nil, "0.0.0.0:8888", false}
	_, err := toml.DecodeFile("slack.toml", &config)
	return config, err
}
