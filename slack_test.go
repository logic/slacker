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
				"token":   {"valid-token"},
				"command": {"/valid-command"},
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
				"token":   {"invalid-token"},
				"command": {"/valid-command"},
			},
			true,
		},
		{
			"POST",
			url.Values{
				"token":   {"valid-token"},
				"command": {"/invalid-command"},
			},
			true,
		},
	}
	Config = Configuration{
		Tokens:            []string{"valid-token"},
		ListenAddress:     "127.0.0.1",
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
		err := SlackDispatcher(w, req)
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
