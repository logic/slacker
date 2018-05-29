// Copyright 2015, 2016, 2017 Ed Marshall. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the COPYING file.

package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
)

// Commands holds our registry of recognized /-commands
var Commands = map[string]ErrorHandler{}

// SlackDispatcher does basic validation, then routes Slack slash commands
// to an appropriate handler.
func SlackDispatcher(w http.ResponseWriter, req *http.Request) error {
	if req.Method != "POST" {
		return StatusError{http.StatusBadRequest,
			errors.New("Only POST is supported")}
	}
	if len(Config.Tokens) > 0 {
		token := req.FormValue("token")
		found := false
		for _, t := range Config.Tokens {
			if token == t {
				found = true
				break
			}
		}
		if !found {
			return StatusError{http.StatusBadRequest,
				fmt.Errorf("Token '%s' is invalid", token)}
		}
	}

	// All commands take a short-circuit "-version" argument.
	if req.FormValue("text") == "-version" {
		fmt.Fprint(w, versionString())
		return nil
	}

	command := strings.ToLower(req.FormValue("command"))
	if handler, ok := Commands[command]; ok {
		log.Printf("[%d] %s@%s:%s %s %s",
			RequestID(req.Context()),
			req.FormValue("user_name"),
			req.FormValue("team_domain"),
			req.FormValue("channel_name"),
			command,
			req.FormValue("text"))
		return handler(w, req)
	}
	return StatusError{http.StatusBadRequest,
		fmt.Errorf("Command '%s' is invalid", command)}
}
