// Copyright 2015 Ed Marshall. All rights reserved.
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

// SlashCommands represents a slack command and a handler for it
type SlashCommands map[string]ErrorHandler

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
