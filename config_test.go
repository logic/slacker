// Copyright 2015, 2016, 2017 Ed Marshall. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the COPYING file.

package main

import (
	"strings"
	"testing"
	"time"
)

func TestConfigListenAddress(t *testing.T) {
	if _, err := LoadConfig(strings.NewReader("Tokens = [\"a\", \"b\"]\n")); err != nil {
		t.Error("Error parsing TOML configuration:", err)
	}
}

func TestConfigTokens(t *testing.T) {
	if _, err := LoadConfig(strings.NewReader("ListenAddress = \"1.2.3.4\"\n")); err != nil {
		t.Error("Error parsing TOML configuration:", err)
	}
}

func TestConfigAsyncResponse(t *testing.T) {
	if _, err := LoadConfig(strings.NewReader("AsyncResponse = true\n")); err != nil {
		t.Error("Error parsing TOML configuration:", err)
	}
}

func TestConfigHTTPClientTimeout(t *testing.T) {
	config, err := LoadConfig(strings.NewReader("HTTPClientTimeout = 10\n"))
	if err != nil {
		t.Error("Error parsing TOML configuration:", err)
	}
	if config.HTTPClientTimeout != 10*time.Second {
		t.Errorf("Timeout incorrectly converted: 10 -> %d", config.HTTPClientTimeout)
	}
}
