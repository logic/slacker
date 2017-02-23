// Copyright 2015, 2016, 2017 Ed Marshall. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the COPYING file.

package main

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"
)

func TestSlackDispatcher(t *testing.T) {
	var tests = []struct {
		method string
		form   url.Values
		fail   bool
	}{
		{
			"POST",
			url.Values{
				"token":        {"valid-token"},
				"command":      {"/valid-command"},
				"text":         {"arguments"},
				"user_name":    {"user"},
				"channel_name": {"channel"},
				"team_domain":  {"domain"},
			},
			false,
		},
		{
			"GET",
			nil,
			true,
		},
		{
			"POST",
			url.Values{
				"token":        {"invalid-token"},
				"command":      {"/valid-command"},
				"text":         {"arguments"},
				"user_name":    {"user"},
				"channel_name": {"channel"},
				"team_domain":  {"domain"},
			},
			true,
		},
		{
			"POST",
			url.Values{
				"token":        {"valid-token"},
				"command":      {"/invalid-command"},
				"text":         {"arguments"},
				"user_name":    {"user"},
				"channel_name": {"channel"},
				"team_domain":  {"domain"},
			},
			true,
		},
	}
	Config = Configuration{
		Tokens:            []string{"valid-token"},
		ListenAddress:     "127.0.0.1:8080",
		AsyncResponse:     false,
		HTTPClientTimeout: time.Duration(10) * time.Second,
	}
	handler := ErrorHandler(func(w http.ResponseWriter, r *http.Request) error {
		return nil
	})
	Commands = SlashCommands{
		"/valid-command": handler,
	}
	uri := "http://slacker.logic.github.io/cmd"

	for _, test := range tests {
		req := httptest.NewRequest(test.method, uri, strings.NewReader(test.form.Encode()))
		if test.method == "POST" {
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		}
		w := httptest.NewRecorder()
		ctx := NewContext(req.Context(), req)
		err := SlackDispatcher(w, req.WithContext(ctx))
		if test.fail {
			if err == nil {
				t.Error(err)
			}
		} else {
			if err != nil {
				t.Error("SlackDispatcher failed:", err)
			}
			if w.Code != 200 {
				t.Error("SlackDispatcher failure code:", w.Code)
			}
		}
	}
}
