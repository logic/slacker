// Copyright 2015, 2016, 2017 Ed Marshall. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the COPYING file.

package main

import (
	"strings"
	"testing"
	"time"
)

func TestTokens(t *testing.T) {
	var tk Tokens
	if tk.String() != "[]" {
		t.Error("expected [], got", tk.String())
	}
	tk.Set("token1")
	if tk.String() != "[token1]" {
		t.Error("expected [token1], got", tk.String())
	}
	tk.Set("token2")
	if tk.String() != "[token1 token2]" {
		t.Error("expected [token1 token2], got", tk.String())
	}
}

func TestConfigListenAddress(t *testing.T) {
	var c Configuration
	if err := LoadConfig(&c, strings.NewReader("Tokens = [\"a\", \"b\"]\n")); err != nil {
		t.Error("Error parsing TOML configuration:", err)
	}
}

func TestConfigTokens(t *testing.T) {
	var c Configuration
	if err := LoadConfig(&c, strings.NewReader("ListenAddress = \"1.2.3.4\"\n")); err != nil {
		t.Error("Error parsing TOML configuration:", err)
	}
}

func TestConfigAsyncResponse(t *testing.T) {
	var c Configuration
	if err := LoadConfig(&c, strings.NewReader("AsyncResponse = true\n")); err != nil {
		t.Error("Error parsing TOML configuration:", err)
	}
}

func TestConfigHTTPClientTimeout(t *testing.T) {
	var c Configuration
	if err := LoadConfig(&c, strings.NewReader("HTTPClientTimeout = 10\n")); err != nil {
		t.Error("Error parsing TOML configuration:", err)
	}
	if c.HTTPClientTimeout != 10*time.Second {
		t.Errorf("Timeout incorrectly converted: 10 -> %d", c.HTTPClientTimeout)
	}
}
